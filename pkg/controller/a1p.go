// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0

package controller

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	policyschemas "github.com/onosproject/onos-a1-dm/go/policy_schemas"
	policystatusv2 "github.com/onosproject/onos-a1-dm/go/policy_status/v2"
	"github.com/onosproject/onos-a1t/pkg/rnib"
	"github.com/onosproject/onos-a1t/pkg/store"
	"github.com/onosproject/onos-a1t/pkg/stream"
	"github.com/onosproject/onos-a1t/pkg/utils"
	"github.com/onosproject/onos-api/go/onos/a1t/a1"
	"github.com/onosproject/onos-lib-go/pkg/errors"
)

func NewA1PController(subscriptionStore store.Store, rnibClient rnib.TopoClient, streamBroker stream.Broker) A1PController {
	return &a1pController{
		subscriptionStore: subscriptionStore,
		rnibClient:        rnibClient,
		streamBroker:      streamBroker,
	}
}

type A1PController interface {
	HandlePolicyCreate(ctx context.Context, policyID, policyTypeID string, params map[string]string, policyObject map[string]interface{}) error
	HandlePolicyDelete(ctx context.Context, policyID, policyTypeID string) error
	HandlePolicyUpdate(ctx context.Context, policyID, policyTypeID string, params map[string]string, policyObject map[string]interface{}) error
	HandleGetPolicyTypes(ctx context.Context) []string
	HandleGetPolicytypesPolicyTypeId(ctx context.Context, policyTypeID string) (map[string]interface{}, map[string]interface{}, error)
	HandleGetPolicytypesPolicyTypeIdPolicies(ctx context.Context, policyTypeID string) ([]string, error)
	HandleGetPolicy(ctx context.Context, policyID, policyTypeID string) (map[string]interface{}, error)
	HandleGetPolicyStatus(ctx context.Context, policyID, policyTypeID string) (map[string]interface{}, error)
	Receiver(ctx context.Context) error
}

type a1pController struct {
	subscriptionStore store.Store
	rnibClient        rnib.TopoClient
	streamBroker      stream.Broker
}

func (a *a1pController) Receiver(ctx context.Context) error {
	return a.watchSubStore(ctx)
}

func (a *a1pController) watchSubStore(ctx context.Context) error {
	log.Info("Start watching subscription store at a1p controller")
	ch := make(chan store.Event)
	go a.subStoreListener(ctx, ch)
	err := a.subscriptionStore.Watch(ctx, ch)
	if err != nil {
		close(ch)
		log.Error(err)
		return err
	}
	return nil
}

func (a *a1pController) subStoreListener(ctx context.Context, ch chan store.Event) {
	var err error
	for e := range ch {
		entry := e.Value.(*store.Entry)
		switch e.Type {
		case store.Created:
			err = a.createEventSubStoreHandler(ctx, entry)
			if err != nil {
				log.Warn(err)
			}
		case store.Deleted:
			err = a.deleteEventSubStoreHandler(ctx, entry)
			if err != nil {
				log.Warn(err)
			}
		}
	}
}

func (a *a1pController) createEventSubStoreHandler(ctx context.Context, entry *store.Entry) error {
	log.Infof("Subscription store entry %v was just created or updated", *entry)
	key := entry.Key.(store.SubscriptionKey)
	targetXAppID := string(key.TargetXAppID)
	msgCh := make(chan *stream.SBStreamMessage)
	sbID, nbID := stream.GetStreamID(stream.A1PController, stream.GetEndpointIDWithTargetXAppID(targetXAppID, stream.PolicyManagement))
	a.streamBroker.AddStream(ctx, nbID)
	a.streamBroker.AddStream(ctx, sbID)

	go func(msgCh chan *stream.SBStreamMessage) {
		for msg := range msgCh {
			err := a.dispatchReceivedMsg(ctx, msg)
			if err != nil {
				log.Error(err)
			}
		}
	}(msgCh)

	watcherID := uuid.New()
	log.Infof("New watcher %v added", watcherID)
	err := a.streamBroker.Watch(nbID, msgCh, watcherID)
	if err != nil {
		log.Error(err)
		return err
	}

	return nil
}

func (a *a1pController) deleteEventSubStoreHandler(ctx context.Context, entry *store.Entry) error {
	log.Infof("Subscription store entry %v was just deleted", *entry)
	// nothing to do with it - stream delete process should be running in southbound manager
	// for the future, if necessary, it should have
	return nil
}

func (a *a1pController) dispatchReceivedMsg(ctx context.Context, sbMessage *stream.SBStreamMessage) error {
	log.Infof("Received msg: %v", sbMessage)
	if sbMessage.A1SBIRPCType == stream.PolicyStatus && sbMessage.A1SBIMessageType == stream.PolicyStatusMessage {
		log.Infof("Received status msg: %v", sbMessage)
		msg := sbMessage.Payload.(*a1.PolicyStatusMessage)
		uri := msg.NotificationDestination
		payload := msg.Message.Payload
		ack := &a1.PolicyAckMessage{
			PolicyType: msg.PolicyType,
			PolicyId:   msg.PolicyId,
			Message: &a1.AckMessage{
				Header: msg.Message.Header,
			},
			NotificationDestination: msg.NotificationDestination,
		}
		sbID, _ := stream.GetStreamID(stream.A1PController, stream.GetEndpointIDWithTargetXAppID(sbMessage.TargetXAppID, stream.PolicyManagement))
		resp, err := http.Post(uri, "application/json", bytes.NewBuffer(payload))
		if err != nil {
			ack.Message.Result = &a1.Result{
				Success: false,
				Reason:  err.Error(),
			}

		} else {
			ack.Message.Result = &a1.Result{
				Success: true,
			}
		}
		log.Infof("PolicyStatus forwarding Resp: %v", resp)
		ackSbMessage := stream.NewSBStreamMessage(sbMessage.TargetXAppID, stream.PolicyAckMessage, sbMessage.A1SBIRPCType, stream.PolicyManagement, ack)
		err = a.streamBroker.Send(sbID, ackSbMessage)
		if err != nil {
			return err
		}
	}
	return nil
}

func (a *a1pController) HandlePolicyCreate(ctx context.Context, policyID, policyTypeID string, params map[string]string, policyObject map[string]interface{}) error {
	targetXAppIDs, err := a.rnibClient.GetXAppIDsForPolicyTypeID(ctx, policyTypeID)
	if err != nil {
		log.Error(err)
		return err
	}

	log.Infof("targetXAppIDs %v for policyTypeID %v", targetXAppIDs, policyTypeID)

	var resErr error = nil

	for _, targetXAppID := range targetXAppIDs {
		requestID := uuid.New().String()
		obj, err := json.Marshal(policyObject)
		if err != nil {
			log.Error(err)
			return err
		}
		reqMsg := &a1.PolicyRequestMessage{
			PolicyId: policyID,
			PolicyType: &a1.PolicyType{
				Id: policyTypeID,
			},
			Message: &a1.RequestMessage{
				Header: &a1.Header{
					RequestId:   requestID,
					AppId:       targetXAppID,
					Encoding:    a1.Encoding_PROTO,
					PayloadType: a1.PayloadType_POLICY,
				},
				Payload: obj,
			},
		}
		if callbackURI, ok := params[utils.NotificationDestination]; ok {
			reqMsg.NotificationDestination = callbackURI
		}
		sbID, nbID := stream.GetStreamID(stream.A1PController, stream.GetEndpointIDWithTargetXAppID(targetXAppID, stream.PolicyManagement))
		sbMessage := stream.NewSBStreamMessage(targetXAppID, stream.PolicyRequestMessage, stream.PolicySetup, stream.PolicyManagement, reqMsg)
		respCh := make(chan *stream.SBStreamMessage)

		watcherID := uuid.New()
		log.Infof("New watcher %v added", watcherID)
		outputCh := make(chan interface{}, 1)
		//eCh := make(chan error, 1)

		go waitRespMsgWithTimer(nbID, watcherID, requestID, respCh, outputCh, TimeoutTimer, a.streamBroker)

		err = a.streamBroker.Watch(nbID, respCh, watcherID)
		if err != nil {
			log.Error(err)

			a.streamBroker.DeleteWatcher(nbID, watcherID)
			return err
		}

		err = a.streamBroker.Send(sbID, sbMessage)
		if err != nil {
			log.Error(err)

			a.streamBroker.DeleteWatcher(nbID, watcherID)
			return err
		}

		err, _ = checkOutput(<-outputCh)
		if err != nil {
			resErr = err
		}
		close(outputCh)
	}
	if resErr != nil {
		log.Error(resErr)
	}

	return resErr
}

func (a *a1pController) HandlePolicyDelete(ctx context.Context, policyID, policyTypeID string) error {
	targetXAppIDs, err := a.rnibClient.GetXAppIDsForPolicyTypeID(ctx, policyTypeID)
	if err != nil {
		log.Error(err)
		return err
	}

	log.Infof("targetXAppIDs %v for policyTypeID %v", targetXAppIDs, policyTypeID)

	var resErr error = nil

	for _, targetXAppID := range targetXAppIDs {
		requestID := uuid.New().String()
		reqMsg := &a1.PolicyRequestMessage{
			PolicyId: policyID,
			PolicyType: &a1.PolicyType{
				Id: policyTypeID,
			},
			Message: &a1.RequestMessage{
				Header: &a1.Header{
					RequestId:   requestID,
					AppId:       targetXAppID,
					Encoding:    a1.Encoding_PROTO,
					PayloadType: a1.PayloadType_POLICY,
				},
			},
		}
		sbID, nbID := stream.GetStreamID(stream.A1PController, stream.GetEndpointIDWithTargetXAppID(targetXAppID, stream.PolicyManagement))
		sbMessage := stream.NewSBStreamMessage(targetXAppID, stream.PolicyResultMessage, stream.PolicyDelete, stream.PolicyManagement, reqMsg)
		respCh := make(chan *stream.SBStreamMessage)

		watcherID := uuid.New()
		log.Infof("New watcher %v added", watcherID)
		outputCh := make(chan interface{}, 1)

		go waitRespMsgWithTimer(nbID, watcherID, requestID, respCh, outputCh, TimeoutTimer, a.streamBroker)

		err = a.streamBroker.Watch(nbID, respCh, watcherID)
		if err != nil {
			log.Error(err)

			a.streamBroker.DeleteWatcher(nbID, watcherID)
			return err
		}

		err = a.streamBroker.Send(sbID, sbMessage)
		if err != nil {
			log.Error(err)

			a.streamBroker.DeleteWatcher(nbID, watcherID)
			return err
		}

		err, _ = checkOutput(<-outputCh)
		if err != nil {
			resErr = err
		}
		close(outputCh)
	}

	if resErr != nil {
		log.Error(resErr)
	}

	return resErr
}

func (a *a1pController) HandlePolicyUpdate(ctx context.Context, policyID, policyTypeID string, params map[string]string, policyObject map[string]interface{}) error {
	targetXAppIDs, err := a.rnibClient.GetXAppIDsForPolicyTypeID(ctx, policyTypeID)

	if err != nil {
		log.Error(err)
		return err
	}

	log.Infof("targetXAppIDs %v for policyTypeID %v", targetXAppIDs, policyTypeID)

	var resErr error = nil
	for _, targetXAppID := range targetXAppIDs {
		requestID := uuid.New().String()
		obj, err := json.Marshal(policyObject)
		if err != nil {
			log.Error(err)
			return err
		}
		reqMsg := &a1.PolicyRequestMessage{
			PolicyId: policyID,
			PolicyType: &a1.PolicyType{
				Id: policyTypeID,
			},
			Message: &a1.RequestMessage{
				Header: &a1.Header{
					RequestId:   requestID,
					AppId:       targetXAppID,
					Encoding:    a1.Encoding_PROTO,
					PayloadType: a1.PayloadType_POLICY,
				},
				Payload: obj,
			},
		}
		if callbackURI, ok := params[utils.NotificationDestination]; ok {
			reqMsg.NotificationDestination = callbackURI
		}
		sbID, nbID := stream.GetStreamID(stream.A1PController, stream.GetEndpointIDWithTargetXAppID(targetXAppID, stream.PolicyManagement))
		sbMessage := stream.NewSBStreamMessage(targetXAppID, stream.PolicyRequestMessage, stream.PolicyUpdate, stream.PolicyManagement, reqMsg)
		respCh := make(chan *stream.SBStreamMessage)

		watcherID := uuid.New()
		log.Infof("New watcher %v added", watcherID)
		outputCh := make(chan interface{}, 1)

		go waitRespMsgWithTimer(nbID, watcherID, requestID, respCh, outputCh, TimeoutTimer, a.streamBroker)

		err = a.streamBroker.Watch(nbID, respCh, watcherID)
		if err != nil {
			log.Error(err)
			a.streamBroker.DeleteWatcher(nbID, watcherID)
			return err
		}

		err = a.streamBroker.Send(sbID, sbMessage)
		if err != nil {
			log.Error(err)
			a.streamBroker.DeleteWatcher(nbID, watcherID)
			return err
		}

		err, _ = checkOutput(<-outputCh)
		if err != nil {
			resErr = err
		}
		close(outputCh)
	}

	if resErr != nil {
		log.Error(resErr)
	}

	return resErr
}

func (a *a1pController) HandleGetPolicyTypes(ctx context.Context) []string {
	results := make([]string, 0)
	policyTypes, err := a.rnibClient.GetPolicyTypes(ctx)
	if err != nil {
		log.Error(err)
		return nil
	}

	for k := range policyTypes {
		results = append(results, string(k))
	}
	return results
}

func (a *a1pController) HandleGetPolicytypesPolicyTypeId(ctx context.Context, policyTypeID string) (map[string]interface{}, map[string]interface{}, error) {
	schema, ok := policyschemas.PolicySchemas[policyTypeID]
	if !ok {
		return nil, nil, errors.NewNotSupported("PolicyTypeID %v is not supported - is it defined in onos-a1-dm?")
	}

	policyTypes, err := a.rnibClient.GetPolicyTypes(ctx)
	if err != nil {
		log.Error(err)
	}
	for k := range policyTypes {
		if string(k) == policyTypeID {
			typeSchema := utils.ConvertStringFormatJsonToMap(schema)
			statusSchema := utils.ConvertStringFormatJsonToMap(policystatusv2.RawSchema)
			return typeSchema, statusSchema, nil
		}
	}
	return nil, nil, errors.NewNotFound("Policy Type ID %v not found", policyTypeID)
}

func (a *a1pController) HandleGetPolicytypesPolicyTypeIdPolicies(ctx context.Context, policyTypeID string) ([]string, error) {
	targetXAppIDs, err := a.rnibClient.GetXAppIDsForPolicyTypeID(ctx, policyTypeID)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	log.Infof("targetXAppIDs %v for policyTypeID %v", targetXAppIDs, policyTypeID)

	objs := make([][]string, 0)

	var resErr error = nil

	for _, targetXAppID := range targetXAppIDs {
		requestID := uuid.New().String()
		reqMsg := &a1.PolicyRequestMessage{
			PolicyType: &a1.PolicyType{
				Id: policyTypeID,
			},
			Message: &a1.RequestMessage{
				Header: &a1.Header{
					RequestId:   requestID,
					AppId:       targetXAppID,
					Encoding:    a1.Encoding_PROTO,
					PayloadType: a1.PayloadType_POLICY,
				},
			},
		}
		sbID, nbID := stream.GetStreamID(stream.A1PController, stream.GetEndpointIDWithTargetXAppID(targetXAppID, stream.PolicyManagement))
		sbMessage := stream.NewSBStreamMessage(targetXAppID, stream.PolicyRequestMessage, stream.PolicyQuery, stream.PolicyManagement, reqMsg)
		respCh := make(chan *stream.SBStreamMessage)

		watcherID := uuid.New()
		log.Infof("New watcher %v added", watcherID)
		outputCh := make(chan interface{}, 1)

		go waitRespMsgWithTimer(nbID, watcherID, requestID, respCh, outputCh, TimeoutTimer, a.streamBroker)

		err = a.streamBroker.Watch(nbID, respCh, watcherID)
		if err != nil {
			log.Error(err)
			a.streamBroker.DeleteWatcher(nbID, watcherID)
			return nil, err
		}

		err = a.streamBroker.Send(sbID, sbMessage)
		if err != nil {
			log.Error(err)
			a.streamBroker.DeleteWatcher(nbID, watcherID)
			return nil, err
		}

		err, resp := checkOutput(<-outputCh)
		if err != nil {
			resErr = err
		}
		close(outputCh)

		var obj []string
		if err == nil {
			err = json.Unmarshal(resp.(*a1.PolicyResultMessage).Message.Payload, &obj)
			if err != nil {
				resErr = err
			} else {
				objs = append(objs, obj)
			}
		}
	}

	if resErr != nil {
		log.Error(resErr)
		return nil, resErr
	}

	if ok, err := utils.PolicyObjListValidate(objs); !ok {
		log.Error(err)
		return nil, err
	}

	return objs[0], nil
}

func (a *a1pController) HandleGetPolicy(ctx context.Context, policyID, policyTypeID string) (map[string]interface{}, error) {
	targetXAppIDs, err := a.rnibClient.GetXAppIDsForPolicyTypeID(ctx, policyTypeID)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	log.Infof("targetXAppIDs %v for policyTypeID %v", targetXAppIDs, policyTypeID)

	objs := make([]map[string]interface{}, 0)

	var resErr error = nil

	for _, targetXAppID := range targetXAppIDs {
		requestID := uuid.New().String()
		reqMsg := &a1.PolicyRequestMessage{
			PolicyId: policyID,
			PolicyType: &a1.PolicyType{
				Id: policyTypeID,
			},
			Message: &a1.RequestMessage{
				Header: &a1.Header{
					RequestId:   requestID,
					AppId:       targetXAppID,
					Encoding:    a1.Encoding_PROTO,
					PayloadType: a1.PayloadType_POLICY,
				},
			},
		}
		sbID, nbID := stream.GetStreamID(stream.A1PController, stream.GetEndpointIDWithTargetXAppID(targetXAppID, stream.PolicyManagement))
		sbMessage := stream.NewSBStreamMessage(targetXAppID, stream.PolicyRequestMessage, stream.PolicyQuery, stream.PolicyManagement, reqMsg)
		respCh := make(chan *stream.SBStreamMessage)

		watcherID := uuid.New()
		log.Infof("New watcher %v added", watcherID)
		outputCh := make(chan interface{}, 1)

		go waitRespMsgWithTimer(nbID, watcherID, requestID, respCh, outputCh, TimeoutTimer, a.streamBroker)

		err = a.streamBroker.Watch(nbID, respCh, watcherID)
		if err != nil {
			log.Error(err)

			a.streamBroker.DeleteWatcher(nbID, watcherID)
			return nil, err
		}

		err = a.streamBroker.Send(sbID, sbMessage)
		if err != nil {
			log.Error(err)

			a.streamBroker.DeleteWatcher(nbID, watcherID)
			return nil, err
		}

		err, resp := checkOutput(<-outputCh)
		if err != nil {
			resErr = err
		}
		close(outputCh)

		var obj map[string]interface{}
		if err == nil {
			err = json.Unmarshal(resp.(*a1.PolicyResultMessage).Message.Payload, &obj)
			if err != nil {
				resErr = err
			} else {
				objs = append(objs, obj)
			}
		}
	}

	if resErr != nil {
		log.Error(resErr)
		return nil, resErr
	}

	if ok, err := utils.PolicyObjListValidate(objs); !ok {
		log.Error(err)
		return nil, err
	}

	return objs[0], nil
}

func (a *a1pController) HandleGetPolicyStatus(ctx context.Context, policyID, policyTypeID string) (map[string]interface{}, error) {
	targetXAppIDs, err := a.rnibClient.GetXAppIDsForPolicyTypeID(ctx, policyTypeID)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	log.Infof("targetXAppIDs %v for policyTypeID %v", targetXAppIDs, policyTypeID)

	objs := make([]map[string]interface{}, 0)

	var resErr error = nil

	for _, targetXAppID := range targetXAppIDs {
		requestID := uuid.New().String()
		reqMsg := &a1.PolicyRequestMessage{
			PolicyId: policyID,
			PolicyType: &a1.PolicyType{
				Id: policyTypeID,
			},
			Message: &a1.RequestMessage{
				Header: &a1.Header{
					RequestId:   requestID,
					AppId:       targetXAppID,
					Encoding:    a1.Encoding_PROTO,
					PayloadType: a1.PayloadType_STATUS,
				},
			},
		}
		sbID, nbID := stream.GetStreamID(stream.A1PController, stream.GetEndpointIDWithTargetXAppID(targetXAppID, stream.PolicyManagement))
		sbMessage := stream.NewSBStreamMessage(targetXAppID, stream.PolicyRequestMessage, stream.PolicyQuery, stream.PolicyManagement, reqMsg)
		respCh := make(chan *stream.SBStreamMessage)

		watcherID := uuid.New()
		log.Infof("New watcher %v added", watcherID)
		outputCh := make(chan interface{}, 1)

		go waitRespMsgWithTimer(nbID, watcherID, requestID, respCh, outputCh, TimeoutTimer, a.streamBroker)

		err = a.streamBroker.Watch(nbID, respCh, watcherID)
		if err != nil {
			log.Error(err)

			a.streamBroker.DeleteWatcher(nbID, watcherID)
			return nil, err
		}

		err = a.streamBroker.Send(sbID, sbMessage)
		if err != nil {
			log.Error(err)

			a.streamBroker.DeleteWatcher(nbID, watcherID)
			return nil, err
		}

		err, resp := checkOutput(<-outputCh)
		if err != nil {
			resErr = err
		}
		close(outputCh)

		var obj map[string]interface{}
		if err == nil {
			err = json.Unmarshal(resp.(*a1.PolicyResultMessage).Message.Payload, &obj)
			if err != nil {
				resErr = err
			} else {
				objs = append(objs, obj)
			}
		}
	}

	if resErr != nil {
		log.Error(resErr)
		return nil, resErr
	}

	if ok, err := utils.PolicyObjListValidate(objs); !ok {
		log.Error(err)
		return nil, err
	}

	return objs[0], nil
}
