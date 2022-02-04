// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0

package controller

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/onosproject/onos-a1t/pkg/rnib"
	"github.com/onosproject/onos-a1t/pkg/store"
	"github.com/onosproject/onos-a1t/pkg/stream"
	"github.com/onosproject/onos-a1t/pkg/utils"
	"github.com/onosproject/onos-api/go/onos/a1t/a1"
	"github.com/onosproject/onos-lib-go/pkg/errors"
	"net/http"
	"time"

	policyschemas "github.com/onosproject/onos-a1-dm/go/policy_schemas"
	"github.com/onosproject/onos-lib-go/pkg/logging"
)

var logA1P = logging.GetLogger("controller", "a1p")

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
	HandleGetPolicytypesPolicyTypeId(ctx context.Context, policyTypeID string) (map[string]interface{}, error)
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
	logA1P.Info("Start watching subscription store at a1p controller")
	ch := make(chan store.Event)
	go a.subStoreListener(ctx, ch)
	err := a.subscriptionStore.Watch(ctx, ch)
	if err != nil {
		close(ch)
		logA1P.Error(err)
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
				logA1P.Warn(err)
			}
		//case store.Updated:
		//	err = a.createEventSubStoreHandler(ctx, entry)
		//	if err != nil {
		//		logA1P.Warn(err)
		//	}
		case store.Deleted:
			err = a.deleteEventSubStoreHandler(ctx, entry)
			if err != nil {
				logA1P.Warn(err)
			}
		}
	}
}

func (a *a1pController) createEventSubStoreHandler(ctx context.Context, entry *store.Entry) error {
	logA1P.Infof("Subscription store entry %v was just created or updated", *entry)
	key := entry.Key.(store.SubscriptionKey)
	targetXAppID := string(key.TargetXAppID)
	msgCh := make(chan *stream.SBStreamMessage)
	nbID := stream.ID{
		SrcEndpointID:  stream.GetEndpointIDWithTargetXAppID(targetXAppID, stream.PolicyManagement),
		DestEndpointID: "a1p-controller",
	}
	sbID := stream.ID{
		DestEndpointID: stream.GetEndpointIDWithTargetXAppID(targetXAppID, stream.PolicyManagement),
		SrcEndpointID:  "a1p-controller",
	}

	a.streamBroker.AddStream(ctx, nbID)
	a.streamBroker.AddStream(ctx, sbID)

	go func(msgCh chan *stream.SBStreamMessage) {
		for msg := range msgCh {
			err := a.dispatchReceivedMsg(ctx, msg)
			if err != nil {
				logA1P.Error(err)
			}
		}
	}(msgCh)

	watcherID := uuid.New()
	logA1P.Infof("New watcher %v added", watcherID)
	err := a.streamBroker.Watch(nbID, msgCh, watcherID)
	if err != nil {
		logA1P.Error(err)
		return err
	}
	return nil
}

func (a *a1pController) deleteEventSubStoreHandler(ctx context.Context, entry *store.Entry) error {
	logA1P.Infof("Subscription store entry %v was just deleted", *entry)
	// nothing to do with it - stream delete process should be running in southbound manager
	return nil
}

func (a *a1pController) dispatchReceivedMsg(ctx context.Context, sbMessage *stream.SBStreamMessage) error {
	logA1P.Infof("Received msg: %v", sbMessage)
	if sbMessage.A1SBIRPCType == stream.PolicyStatus && sbMessage.A1SBIMessageType == stream.PolicyStatusMessage {
		logA1P.Infof("Received status msg: %v", sbMessage)
		msg := sbMessage.Payload.(*a1.PolicyStatusMessage)
		uri := msg.NotificationDestination
		payload := msg.Message.Payload

		ackSbMessage := &stream.SBStreamMessage{
			A1SBIMessageType: stream.PolicyAckMessage,
			TargetXAppID:     sbMessage.TargetXAppID,
			A1Service:        stream.PolicyManagement,
			A1SBIRPCType:     sbMessage.A1SBIRPCType,
		}

		ack := &a1.PolicyAckMessage{
			PolicyType: msg.PolicyType,
			PolicyId:   msg.PolicyId,
			Message: &a1.AckMessage{
				Header: msg.Message.Header,
			},
			NotificationDestination: msg.NotificationDestination,
		}

		sbID := stream.ID{
			SrcEndpointID:  "a1p-controller",
			DestEndpointID: stream.GetEndpointIDWithTargetXAppID(sbMessage.TargetXAppID, stream.PolicyManagement),
		}

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
		logA1P.Infof("PolicyStatus forwarding Resp: %v", resp)
		ackSbMessage.Payload = ack
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
		logA1P.Error(err)
		return err
	}

	logA1P.Infof("targetXAppIDs %v for policyTypeID %v", targetXAppIDs, policyTypeID)

	var resErr error = nil

	for _, targetXAppID := range targetXAppIDs {
		requestID := uuid.New().String()
		obj, err := json.Marshal(policyObject)
		if err != nil {
			logA1P.Error(err)
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

		sbMessage := &stream.SBStreamMessage{
			A1SBIMessageType: stream.PolicyRequestMessage,
			A1SBIRPCType:     stream.PolicySetup,
			A1Service:        stream.PolicyManagement,
			Payload:          reqMsg,
		}

		sbID := stream.ID{
			SrcEndpointID:  "a1p-controller",
			DestEndpointID: stream.GetEndpointIDWithTargetXAppID(targetXAppID, stream.PolicyManagement),
		}

		nbID := stream.ID{
			DestEndpointID: "a1p-controller",
			SrcEndpointID:  stream.GetEndpointIDWithTargetXAppID(targetXAppID, stream.PolicyManagement),
		}

		respCh := make(chan *stream.SBStreamMessage)
		timerCh := make(chan bool, 1)

		go func(timer time.Duration, ch chan bool) {
			time.Sleep(timer)
			timerCh <- true
			close(timerCh)
		}(TimeoutTimer, timerCh)

		watcherID := uuid.New()
		logA1P.Infof("New watcher %v added", watcherID)
		eCh := make(chan error, 1)

		go func(id stream.ID, wID uuid.UUID, eCh chan error) {
			for {
				select {
				case msg := <-respCh:
					logA1P.Infof("Message %v received", msg)
					switch result := msg.Payload.(type) {
					case *a1.PolicyResultMessage:
						if result.Message.Header.RequestId == requestID {
							logA1P.Infof("same request ID matched: Message %v", msg)
							if !result.Message.Result.Success {
								logA1P.Error(fmt.Errorf(result.Message.Result.Reason))
								eCh <- fmt.Errorf(result.Message.Result.Reason)
								a.streamBroker.DeleteWatcher(nbID, watcherID)
								return
							}
							eCh <- nil

							a.streamBroker.DeleteWatcher(nbID, watcherID)
							return
						}
					}
				case <-timerCh:
					logA1P.Error(errors.NewTimeout("Could not receive PolicyResultMessage in time (timer: %v)", TimeoutTimer))
					eCh <- errors.NewTimeout("Could not receive PolicyResultMessage in time (timer: %v)", TimeoutTimer)
					a.streamBroker.DeleteWatcher(nbID, watcherID)

					return
				}
			}
		}(nbID, watcherID, eCh)

		err = a.streamBroker.Watch(nbID, respCh, watcherID)
		if err != nil {
			logA1P.Error(err)

			a.streamBroker.DeleteWatcher(nbID, watcherID)
			return err
		}

		err = a.streamBroker.Send(sbID, sbMessage)
		if err != nil {
			logA1P.Error(err)

			a.streamBroker.DeleteWatcher(nbID, watcherID)
			return err
		}

		err = <-eCh
		if err != nil {
			resErr = err
		}
		close(eCh)
	}
	if resErr != nil {
		logA1P.Error(resErr)
	}

	return resErr
}

func (a *a1pController) HandlePolicyDelete(ctx context.Context, policyID, policyTypeID string) error {
	targetXAppIDs, err := a.rnibClient.GetXAppIDsForPolicyTypeID(ctx, policyTypeID)
	if err != nil {
		logA1P.Error(err)
		return err
	}

	logA1P.Infof("targetXAppIDs %v for policyTypeID %v", targetXAppIDs, policyTypeID)

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

		sbMessage := &stream.SBStreamMessage{
			A1SBIMessageType: stream.PolicyRequestMessage,
			A1SBIRPCType:     stream.PolicyDelete,
			A1Service:        stream.PolicyManagement,
			Payload:          reqMsg,
		}

		sbID := stream.ID{
			SrcEndpointID:  "a1p-controller",
			DestEndpointID: stream.GetEndpointIDWithTargetXAppID(targetXAppID, stream.PolicyManagement),
		}

		nbID := stream.ID{
			DestEndpointID: "a1p-controller",
			SrcEndpointID:  stream.GetEndpointIDWithTargetXAppID(targetXAppID, stream.PolicyManagement),
		}

		respCh := make(chan *stream.SBStreamMessage)
		timerCh := make(chan bool, 1)

		go func(timer time.Duration, ch chan bool) {
			time.Sleep(timer)
			timerCh <- true
			close(timerCh)
		}(TimeoutTimer, timerCh)

		watcherID := uuid.New()
		logA1P.Infof("New watcher %v added", watcherID)
		eCh := make(chan error, 1)

		go func(id stream.ID, wID uuid.UUID, eCh chan error) {
			for {
				select {
				case msg := <-respCh:
					logA1P.Infof("Message %v received", msg)
					switch result := msg.Payload.(type) {
					case *a1.PolicyResultMessage:
						if result.Message.Header.RequestId == requestID {
							logA1P.Infof("same request ID matched: Message %v", msg)
							if !result.Message.Result.Success {
								logA1P.Error(fmt.Errorf(result.Message.Result.Reason))
								eCh <- fmt.Errorf(result.Message.Result.Reason)

								a.streamBroker.DeleteWatcher(nbID, watcherID)
								return
							}
							eCh <- nil

							a.streamBroker.DeleteWatcher(nbID, watcherID)
							return
						}
					}
				case <-timerCh:
					logA1P.Error(errors.NewTimeout("Could not receive PolicyResultMessage in time (timer: %v)", TimeoutTimer))
					eCh <- errors.NewTimeout("Could not receive PolicyResultMessage in time (timer: %v)", TimeoutTimer)

					a.streamBroker.DeleteWatcher(nbID, watcherID)
					return
				}
			}
		}(nbID, watcherID, eCh)

		err = a.streamBroker.Watch(nbID, respCh, watcherID)
		if err != nil {
			logA1P.Error(err)

			a.streamBroker.DeleteWatcher(nbID, watcherID)
			return err
		}

		err = a.streamBroker.Send(sbID, sbMessage)
		if err != nil {
			logA1P.Error(err)

			a.streamBroker.DeleteWatcher(nbID, watcherID)
			return err
		}

		err = <-eCh
		if err != nil {
			resErr = err
		}
		close(eCh)
	}

	if resErr != nil {
		logA1P.Error(resErr)
	}

	return resErr
}

func (a *a1pController) HandlePolicyUpdate(ctx context.Context, policyID, policyTypeID string, params map[string]string, policyObject map[string]interface{}) error {
	targetXAppIDs, err := a.rnibClient.GetXAppIDsForPolicyTypeID(ctx, policyTypeID)

	if err != nil {
		logA1P.Error(err)
		return err
	}

	logA1P.Infof("targetXAppIDs %v for policyTypeID %v", targetXAppIDs, policyTypeID)

	var resErr error = nil
	for _, targetXAppID := range targetXAppIDs {
		requestID := uuid.New().String()
		obj, err := json.Marshal(policyObject)
		if err != nil {
			logA1P.Error(err)
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

		sbMessage := &stream.SBStreamMessage{
			A1SBIMessageType: stream.PolicyRequestMessage,
			A1SBIRPCType:     stream.PolicyUpdate,
			A1Service:        stream.PolicyManagement,
			Payload:          reqMsg,
		}

		sbID := stream.ID{
			SrcEndpointID:  "a1p-controller",
			DestEndpointID: stream.GetEndpointIDWithTargetXAppID(targetXAppID, stream.PolicyManagement),
		}
		nbID := stream.ID{
			DestEndpointID: "a1p-controller",
			SrcEndpointID:  stream.GetEndpointIDWithTargetXAppID(targetXAppID, stream.PolicyManagement),
		}
		respCh := make(chan *stream.SBStreamMessage)
		timerCh := make(chan bool, 1)

		go func(timer time.Duration, ch chan bool) {
			time.Sleep(timer)
			timerCh <- true
			close(timerCh)
		}(TimeoutTimer, timerCh)

		watcherID := uuid.New()
		logA1P.Infof("New watcher %v added", watcherID)
		eCh := make(chan error, 1)

		go func(id stream.ID, wID uuid.UUID, eCh chan error) {
			for {
				select {
				case msg := <-respCh:
					logA1P.Infof("Message %v received", msg)
					switch result := msg.Payload.(type) {
					case *a1.PolicyResultMessage:
						if result.Message.Header.RequestId == requestID {
							logA1P.Infof("same request ID matched: Message %v", msg)
							if !result.Message.Result.Success {
								logA1P.Error(fmt.Errorf(result.Message.Result.Reason))
								eCh <- fmt.Errorf(result.Message.Result.Reason)

								a.streamBroker.DeleteWatcher(nbID, watcherID)
								return
							}
							eCh <- nil

							a.streamBroker.DeleteWatcher(nbID, watcherID)
							return
						}
					}
				case <-timerCh:
					logA1P.Error(errors.NewTimeout("Could not receive PolicyResultMessage in time (timer: %v)", TimeoutTimer))
					eCh <- errors.NewTimeout("Could not receive PolicyResultMessage in time (timer: %v)", TimeoutTimer)

					a.streamBroker.DeleteWatcher(nbID, watcherID)
					return
				}
			}
		}(nbID, watcherID, eCh)

		err = a.streamBroker.Watch(nbID, respCh, watcherID)
		if err != nil {
			logA1P.Error(err)
			a.streamBroker.DeleteWatcher(nbID, watcherID)
			return err
		}

		err = a.streamBroker.Send(sbID, sbMessage)
		if err != nil {
			logA1P.Error(err)
			a.streamBroker.DeleteWatcher(nbID, watcherID)
			return err
		}

		err = <-eCh
		if err != nil {
			resErr = err
		}
		close(eCh)
	}

	if resErr != nil {
		logA1P.Error(resErr)
	}

	return resErr
}

func (a *a1pController) HandleGetPolicyTypes(ctx context.Context) []string {
	results := make([]string, 0)
	policyTypes, err := a.rnibClient.GetPolicyTypes(ctx)
	if err != nil {
		logA1P.Error(err)
		return nil
	}

	for k := range policyTypes {
		results = append(results, string(k))
	}
	return results
}

func (a *a1pController) HandleGetPolicytypesPolicyTypeId(ctx context.Context, policyTypeID string) (map[string]interface{}, error) {
	schema, ok := policyschemas.PolicySchemas[policyTypeID]
	if !ok {
		return nil, errors.NewNotSupported("PolicyTypeID %v is not supported - is it defined in onos-a1-dm?")
	}

	policyTypes, err := a.rnibClient.GetPolicyTypes(ctx)
	if err != nil {
		logA1P.Error(err)
	}
	for k := range policyTypes {
		if string(k) == policyTypeID {
			return utils.ConvertStringFormatJsonToMap(schema), nil
		}
	}
	return nil, errors.NewNotFound("Policy Type ID %v not found", policyTypeID)
}

func (a *a1pController) HandleGetPolicytypesPolicyTypeIdPolicies(ctx context.Context, policyTypeID string) ([]string, error) {
	targetXAppIDs, err := a.rnibClient.GetXAppIDsForPolicyTypeID(ctx, policyTypeID)
	if err != nil {
		logA1P.Error(err)
		return nil, err
	}

	logA1P.Infof("targetXAppIDs %v for policyTypeID %v", targetXAppIDs, policyTypeID)

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

		sbMessage := &stream.SBStreamMessage{
			A1SBIMessageType: stream.PolicyRequestMessage,
			A1SBIRPCType:     stream.PolicyQuery,
			A1Service:        stream.PolicyManagement,
			Payload:          reqMsg,
		}

		sbID := stream.ID{
			SrcEndpointID:  "a1p-controller",
			DestEndpointID: stream.GetEndpointIDWithTargetXAppID(targetXAppID, stream.PolicyManagement),
		}

		nbID := stream.ID{
			DestEndpointID: "a1p-controller",
			SrcEndpointID:  stream.GetEndpointIDWithTargetXAppID(targetXAppID, stream.PolicyManagement),
		}

		respCh := make(chan *stream.SBStreamMessage)
		timerCh := make(chan bool, 1)

		go func(timer time.Duration, ch chan bool) {
			time.Sleep(timer)
			timerCh <- true
			close(timerCh)
		}(TimeoutTimer, timerCh)

		watcherID := uuid.New()
		logA1P.Infof("New watcher %v added", watcherID)
		eCh := make(chan interface{}, 1)

		go func(id stream.ID, wID uuid.UUID, eCh chan interface{}) {
			for {
				select {
				case msg := <-respCh:
					logA1P.Infof("Message %v received", msg)
					switch result := msg.Payload.(type) {
					case *a1.PolicyResultMessage:
						if result.Message.Header.RequestId == requestID {
							logA1P.Infof("same request ID matched: Message %v", msg)
							if !result.Message.Result.Success {
								logA1P.Error(fmt.Errorf(result.Message.Result.Reason))
								eCh <- fmt.Errorf(result.Message.Result.Reason)

								a.streamBroker.DeleteWatcher(nbID, watcherID)
								return
							}
							var obj []string
							err = json.Unmarshal(result.Message.Payload, &obj)
							if err != nil {
								logA1P.Error(err)
								eCh <- err

								a.streamBroker.DeleteWatcher(nbID, watcherID)
								return
							}
							eCh <- obj

							a.streamBroker.DeleteWatcher(nbID, watcherID)
							return
						}
					}
				case <-timerCh:
					logA1P.Error(errors.NewTimeout("Could not receive PolicyResultMessage in time (timer: %v)", TimeoutTimer))
					eCh <- errors.NewTimeout("Could not receive PolicyResultMessage in time (timer: %v)", TimeoutTimer)

					a.streamBroker.DeleteWatcher(nbID, watcherID)
					return
				}
			}
		}(nbID, watcherID, eCh)
		err = a.streamBroker.Watch(nbID, respCh, watcherID)
		if err != nil {
			logA1P.Error(err)
			a.streamBroker.DeleteWatcher(nbID, watcherID)
			return nil, err
		}

		err = a.streamBroker.Send(sbID, sbMessage)
		if err != nil {
			logA1P.Error(err)
			a.streamBroker.DeleteWatcher(nbID, watcherID)
			return nil, err
		}

		o := <-eCh
		switch o := o.(type) {
		case error:
			resErr = o
		case []string:
			objs = append(objs, o)
		}
		close(eCh)
	}

	if resErr != nil {
		logA1P.Error(resErr)
		return nil, resErr
	}

	if ok, err := utils.PolicyObjListValidate(objs); !ok {
		logA1P.Error(err)
		return nil, err
	}

	return objs[0], nil
}

func (a *a1pController) HandleGetPolicy(ctx context.Context, policyID, policyTypeID string) (map[string]interface{}, error) {
	targetXAppIDs, err := a.rnibClient.GetXAppIDsForPolicyTypeID(ctx, policyTypeID)
	if err != nil {
		logA1P.Error(err)
		return nil, err
	}

	logA1P.Infof("targetXAppIDs %v for policyTypeID %v", targetXAppIDs, policyTypeID)

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

		sbMessage := &stream.SBStreamMessage{
			A1SBIMessageType: stream.PolicyRequestMessage,
			A1SBIRPCType:     stream.PolicyQuery,
			A1Service:        stream.PolicyManagement,
			Payload:          reqMsg,
		}

		sbID := stream.ID{
			SrcEndpointID:  "a1p-controller",
			DestEndpointID: stream.GetEndpointIDWithTargetXAppID(targetXAppID, stream.PolicyManagement),
		}

		nbID := stream.ID{
			DestEndpointID: "a1p-controller",
			SrcEndpointID:  stream.GetEndpointIDWithTargetXAppID(targetXAppID, stream.PolicyManagement),
		}

		respCh := make(chan *stream.SBStreamMessage)
		timerCh := make(chan bool, 1)

		go func(timer time.Duration, ch chan bool) {
			time.Sleep(timer)
			timerCh <- true
			close(timerCh)
		}(TimeoutTimer, timerCh)

		watcherID := uuid.New()
		logA1P.Infof("New watcher %v added", watcherID)
		eCh := make(chan interface{}, 1)

		go func(id stream.ID, wID uuid.UUID, eCh chan interface{}) {
			for {
				select {
				case msg := <-respCh:
					logA1P.Infof("Message %v received", msg)
					switch result := msg.Payload.(type) {
					case *a1.PolicyResultMessage:
						if result.Message.Header.RequestId == requestID {
							logA1P.Infof("same request ID matched: Message %v", msg)
							if !result.Message.Result.Success {
								logA1P.Error(fmt.Errorf(result.Message.Result.Reason))
								eCh <- fmt.Errorf(result.Message.Result.Reason)

								a.streamBroker.DeleteWatcher(nbID, watcherID)
								return
							}
							var obj map[string]interface{}
							err = json.Unmarshal(result.Message.Payload, &obj)
							if err != nil {
								logA1P.Error(err)
								eCh <- err

								a.streamBroker.DeleteWatcher(nbID, watcherID)
								return
							}
							eCh <- obj

							a.streamBroker.DeleteWatcher(nbID, watcherID)
							return
						}
					}
				case <-timerCh:
					logA1P.Error(errors.NewTimeout("Could not receive PolicyResultMessage in time (timer: %v)", TimeoutTimer))
					eCh <- errors.NewTimeout("Could not receive PolicyResultMessage in time (timer: %v)", TimeoutTimer)

					a.streamBroker.DeleteWatcher(nbID, watcherID)
					return
				}
			}
		}(nbID, watcherID, eCh)

		err = a.streamBroker.Watch(nbID, respCh, watcherID)
		if err != nil {
			logA1P.Error(err)

			a.streamBroker.DeleteWatcher(nbID, watcherID)
			return nil, err
		}

		err = a.streamBroker.Send(sbID, sbMessage)
		if err != nil {
			logA1P.Error(err)

			a.streamBroker.DeleteWatcher(nbID, watcherID)
			return nil, err
		}

		o := <-eCh
		switch o := o.(type) {
		case error:
			resErr = o
		case map[string]interface{}:
			objs = append(objs, o)
		}
		close(eCh)
	}

	if resErr != nil {
		logA1P.Error(resErr)
		return nil, resErr
	}

	if ok, err := utils.PolicyObjListValidate(objs); !ok {
		logA1P.Error(err)
		return nil, err
	}

	return objs[0], nil
}

func (a *a1pController) HandleGetPolicyStatus(ctx context.Context, policyID, policyTypeID string) (map[string]interface{}, error) {
	targetXAppIDs, err := a.rnibClient.GetXAppIDsForPolicyTypeID(ctx, policyTypeID)
	if err != nil {
		logA1P.Error(err)
		return nil, err
	}

	logA1P.Infof("targetXAppIDs %v for policyTypeID %v", targetXAppIDs, policyTypeID)

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

		sbMessage := &stream.SBStreamMessage{
			A1SBIMessageType: stream.PolicyRequestMessage,
			A1SBIRPCType:     stream.PolicyQuery,
			A1Service:        stream.PolicyManagement,
			Payload:          reqMsg,
		}

		sbID := stream.ID{
			SrcEndpointID:  "a1p-controller",
			DestEndpointID: stream.GetEndpointIDWithTargetXAppID(targetXAppID, stream.PolicyManagement),
		}

		nbID := stream.ID{
			DestEndpointID: "a1p-controller",
			SrcEndpointID:  stream.GetEndpointIDWithTargetXAppID(targetXAppID, stream.PolicyManagement),
		}

		respCh := make(chan *stream.SBStreamMessage)
		timerCh := make(chan bool, 1)

		go func(timer time.Duration, ch chan bool) {
			time.Sleep(timer)
			timerCh <- true
			close(timerCh)
		}(TimeoutTimer, timerCh)

		watcherID := uuid.New()
		logA1P.Infof("New watcher %v added", watcherID)
		eCh := make(chan interface{}, 1)

		go func(id stream.ID, wID uuid.UUID, eCh chan interface{}) {
			for {
				select {
				case msg := <-respCh:
					logA1P.Infof("Message %v received", msg)
					switch result := msg.Payload.(type) {
					case *a1.PolicyResultMessage:
						if result.Message.Header.RequestId == requestID {
							logA1P.Infof("same request ID matched: Message %v", msg)
							if !result.Message.Result.Success {
								logA1P.Error(fmt.Errorf(result.Message.Result.Reason))
								eCh <- fmt.Errorf(result.Message.Result.Reason)

								a.streamBroker.DeleteWatcher(nbID, watcherID)
								return
							}
							var obj map[string]interface{}
							err = json.Unmarshal(result.Message.Payload, &obj)
							if err != nil {
								logA1P.Error(err)
								eCh <- err

								a.streamBroker.DeleteWatcher(nbID, watcherID)
								return
							}
							eCh <- obj

							a.streamBroker.DeleteWatcher(nbID, watcherID)
							return
						}
					}
				case <-timerCh:
					logA1P.Error(errors.NewTimeout("Could not receive PolicyResultMessage in time (timer: %v)", TimeoutTimer))
					eCh <- errors.NewTimeout("Could not receive PolicyResultMessage in time (timer: %v)", TimeoutTimer)

					a.streamBroker.DeleteWatcher(nbID, watcherID)
					return
				}
			}
		}(nbID, watcherID, eCh)

		err = a.streamBroker.Watch(nbID, respCh, watcherID)
		if err != nil {
			logA1P.Error(err)

			a.streamBroker.DeleteWatcher(nbID, watcherID)
			return nil, err
		}

		err = a.streamBroker.Send(sbID, sbMessage)
		if err != nil {
			logA1P.Error(err)

			a.streamBroker.DeleteWatcher(nbID, watcherID)
			return nil, err
		}

		o := <-eCh
		switch o := o.(type) {
		case error:
			resErr = o
		case map[string]interface{}:
			objs = append(objs, o)
		}
		close(eCh)
	}

	if resErr != nil {
		logA1P.Error(resErr)
		return nil, resErr
	}

	if ok, err := utils.PolicyObjListValidate(objs); !ok {
		logA1P.Error(err)
		return nil, err
	}

	return objs[0], nil
}
