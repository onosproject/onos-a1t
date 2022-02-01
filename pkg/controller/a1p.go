// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-Only-1.0

package controller

import (
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
	HandleGetPolicyStatus(ctx context.Context, policyID, policyTypeID string) (bool, error)
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
		case store.Updated:
			err = a.createEventSubStoreHandler(ctx, entry)
			if err != nil {
				logA1P.Warn(err)
			}
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

	go func(msgCh chan *stream.SBStreamMessage) {
		for msg := range msgCh {
			a.dispatchReceivedMsg(ctx, msg)
		}
	}(msgCh)

	watcherID := uuid.New()
	err := a.streamBroker.Watch(nbID, msgCh, watcherID)
	if err != nil {
		return err
	}
	return nil
}

func (a *a1pController) deleteEventSubStoreHandler(ctx context.Context, entry *store.Entry) error {
	logA1P.Infof("Subscription store entry %v was just deleted", *entry)
	return nil
}

func (a *a1pController) dispatchReceivedMsg(ctx context.Context, sbMessage *stream.SBStreamMessage) error {
	logA1P.Infof("Received msg: %v", sbMessage)
	return a.watchSubStore(ctx)
}

func (a *a1pController) HandlePolicyCreate(ctx context.Context, policyID, policyTypeID string, params map[string]string, policyObject map[string]interface{}) error {
	targetXAppIDs, err := a.rnibClient.GetXAppIDsForPolicyTypeID(ctx, policyTypeID)
	if err != nil {
		logA1P.Error(err)
		return err
	}

	logA1P.Infof("targetXAppIDs %v for policyTypeID %v", targetXAppIDs, policyTypeID)

	err = nil

	for _, targetXAppID := range targetXAppIDs {
		requestID := uuid.New().String()
		obj, err := json.Marshal(policyObject)
		reqMsg := &a1.PolicyRequestMessage{
			PolicyId: policyID,
			PolicyType: &a1.PolicyType{
				Id: policyTypeID,
			},
			Message: &a1.RequestMessage{
				Header: &a1.Header{
					RequestId: requestID,
					AppId:     targetXAppID,
					Encoding:  a1.Encoding_PROTO,
				},
				Payload: obj,
			},
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

		watcherID := uuid.New()
		err = a.streamBroker.Watch(nbID, respCh, watcherID)
		if err != nil {
			return err
		}

		go func(timer time.Duration, ch chan bool) {
			time.Sleep(timer)
			timerCh <- true
			close(timerCh)
		}(TimeoutTimer, timerCh)

		err = a.streamBroker.Send(sbID, sbMessage)
		if err != nil {
			return err
		}

	RowLoop:
		for {
			select {
			case msg := <-respCh:
				result := msg.Payload.(*a1.PolicyResultMessage)
				logA1P.Infof("Message %v received", msg)
				if result.Message.Header.RequestId == requestID {
					logA1P.Infof("same request ID matched: Message %v", msg)
					if !result.Message.Result.Success {
						err = fmt.Errorf(result.Message.Result.Reason)
					}
					break RowLoop
				}
			case <-timerCh:
				logA1P.Error(errors.NewTimeout("Could not receive PolicyResultMessage in time (timer: %v)", TimeoutTimer))
				break RowLoop
			}
		}
		a.streamBroker.DeleteWatcher(nbID, watcherID)
		close(respCh)
	}

	return err
}

func (a *a1pController) HandlePolicyDelete(ctx context.Context, policyID, policyTypeID string) error {
	targetXAppIDs, err := a.rnibClient.GetXAppIDsForPolicyTypeID(ctx, policyTypeID)
	if err != nil {
		logA1P.Error(err)
		return err
	}

	logA1P.Infof("targetXAppIDs %v for policyTypeID %v", targetXAppIDs, policyTypeID)

	err = nil

	for _, targetXAppID := range targetXAppIDs {
		requestID := uuid.New().String()
		reqMsg := &a1.PolicyRequestMessage{
			PolicyId: policyID,
			PolicyType: &a1.PolicyType{
				Id: policyTypeID,
			},
			Message: &a1.RequestMessage{
				Header: &a1.Header{
					RequestId: requestID,
					AppId:     targetXAppID,
					Encoding:  a1.Encoding_PROTO,
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

		watcherID := uuid.New()
		err = a.streamBroker.Watch(nbID, respCh, watcherID)
		if err != nil {
			return err
		}

		go func(timer time.Duration, ch chan bool) {
			time.Sleep(timer)
			timerCh <- true
			close(timerCh)
		}(TimeoutTimer, timerCh)

		err = a.streamBroker.Send(sbID, sbMessage)
		if err != nil {
			return err
		}

	RowLoop:
		for {
			select {
			case msg := <-respCh:
				result := msg.Payload.(*a1.PolicyResultMessage)
				logA1P.Infof("Message %v received", msg)
				if result.Message.Header.RequestId == requestID {
					logA1P.Infof("same request ID matched: Message %v", msg)
					if !result.Message.Result.Success {
						err = fmt.Errorf(result.Message.Result.Reason)
					}
					break RowLoop
				}
			case <-timerCh:
				logA1P.Error(errors.NewTimeout("Could not receive PolicyResultMessage in time (timer: %v)", TimeoutTimer))
				break RowLoop
			}
		}
		a.streamBroker.DeleteWatcher(nbID, watcherID)
		close(respCh)
	}

	return err
}

func (a *a1pController) HandlePolicyUpdate(ctx context.Context, policyID, policyTypeID string, params map[string]string, policyObject map[string]interface{}) error {
	targetXAppIDs, err := a.rnibClient.GetXAppIDsForPolicyTypeID(ctx, policyTypeID)
	if err != nil {
		logA1P.Error(err)
		return err
	}

	logA1P.Infof("targetXAppIDs %v for policyTypeID %v", targetXAppIDs, policyTypeID)

	err = nil

	for _, targetXAppID := range targetXAppIDs {
		requestID := uuid.New().String()
		obj, err := json.Marshal(policyObject)
		reqMsg := &a1.PolicyRequestMessage{
			PolicyId: policyID,
			PolicyType: &a1.PolicyType{
				Id: policyTypeID,
			},
			Message: &a1.RequestMessage{
				Header: &a1.Header{
					RequestId: requestID,
					AppId:     targetXAppID,
					Encoding:  a1.Encoding_PROTO,
				},
				Payload: obj,
			},
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

		watcherID := uuid.New()
		err = a.streamBroker.Watch(nbID, respCh, watcherID)
		if err != nil {
			return err
		}

		go func(timer time.Duration, ch chan bool) {
			time.Sleep(timer)
			timerCh <- true
			close(timerCh)
		}(TimeoutTimer, timerCh)

		err = a.streamBroker.Send(sbID, sbMessage)
		if err != nil {
			return err
		}

	RowLoop:
		for {
			select {
			case msg := <-respCh:
				result := msg.Payload.(*a1.PolicyResultMessage)
				logA1P.Infof("Message %v received", msg)
				if result.Message.Header.RequestId == requestID {
					logA1P.Infof("same request ID matched: Message %v", msg)
					if !result.Message.Result.Success {
						err = fmt.Errorf(result.Message.Result.Reason)
					}
					break RowLoop
				}
			case <-timerCh:
				logA1P.Error(errors.NewTimeout("Could not receive PolicyResultMessage in time (timer: %v)", TimeoutTimer))
				break RowLoop
			}
		}
		a.streamBroker.DeleteWatcher(nbID, watcherID)
		close(respCh)
	}

	return err
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

	for _, targetXAppID := range targetXAppIDs {
		requestID := uuid.New().String()
		reqMsg := &a1.PolicyRequestMessage{
			PolicyType: &a1.PolicyType{
				Id: policyTypeID,
			},
			Message: &a1.RequestMessage{
				Header: &a1.Header{
					RequestId: requestID,
					AppId:     targetXAppID,
					Encoding:  a1.Encoding_PROTO,
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

		watcherID := uuid.New()
		err = a.streamBroker.Watch(nbID, respCh, watcherID)
		if err != nil {
			return nil, err
		}

		go func(timer time.Duration, ch chan bool) {
			time.Sleep(timer)
			timerCh <- true
			close(timerCh)
		}(TimeoutTimer, timerCh)

		err = a.streamBroker.Send(sbID, sbMessage)
		if err != nil {
			return nil, err
		}

	RowLoop:
		for {
			select {
			case msg := <-respCh:
				result := msg.Payload.(*a1.PolicyResultMessage)
				logA1P.Infof("Message %v received", msg)
				if result.Message.Header.RequestId == requestID {
					logA1P.Infof("same request ID matched: Message %v", msg)
					var obj []string
					err = json.Unmarshal(result.Message.Payload, &obj)
					if err != nil {
						logA1P.Error(err)
					}
					objs = append(objs, obj)
					break RowLoop
				}
			case <-timerCh:
				logA1P.Error(errors.NewTimeout("Could not receive PolicyResultMessage in time (timer: %v)", TimeoutTimer))
				break RowLoop
			}
		}
		a.streamBroker.DeleteWatcher(nbID, watcherID)
		close(respCh)
	}

	if ok, err := utils.PolicyObjListValidate(objs); !ok {
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

	err = nil

	for _, targetXAppID := range targetXAppIDs {
		requestID := uuid.New().String()
		reqMsg := &a1.PolicyRequestMessage{
			PolicyId: policyID,
			PolicyType: &a1.PolicyType{
				Id: policyTypeID,
			},
			Message: &a1.RequestMessage{
				Header: &a1.Header{
					RequestId: requestID,
					AppId:     targetXAppID,
					Encoding:  a1.Encoding_PROTO,
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

		watcherID := uuid.New()
		err = a.streamBroker.Watch(nbID, respCh, watcherID)
		if err != nil {
			return nil, err
		}

		go func(timer time.Duration, ch chan bool) {
			time.Sleep(timer)
			timerCh <- true
			close(timerCh)
		}(TimeoutTimer, timerCh)

		err = a.streamBroker.Send(sbID, sbMessage)
		if err != nil {
			return nil, err
		}

	RowLoop:
		for {
			select {
			case msg := <-respCh:
				result := msg.Payload.(*a1.PolicyResultMessage)
				logA1P.Infof("Message %v received", msg)
				if result.Message.Header.RequestId == requestID {
					logA1P.Infof("same request ID matched: Message %v", msg)
					var obj map[string]interface{}
					err = json.Unmarshal(result.Message.Payload, &obj)
					objs = append(objs, obj)
					break RowLoop
				}
			case <-timerCh:
				logA1P.Error(errors.NewTimeout("Could not receive PolicyResultMessage in time (timer: %v)", TimeoutTimer))
				break RowLoop
			}
		}
		a.streamBroker.DeleteWatcher(nbID, watcherID)
		close(respCh)
	}

	if ok, err := utils.PolicyObjListValidate(objs); !ok {
		return nil, err
	}

	return objs[0], err
}

func (a *a1pController) HandleGetPolicyStatus(ctx context.Context, policyID, policyTypeID string) (bool, error) {
	panic("implement me")
}

//func (a1p *a1pController) HandleGetPolicyTypes(ctx context.Context) []string {
//	//policyTypes := []string{}
//	//
//	//tmpSubs := getSubscriptionPolicyTypes(ctx, a1p)
//	//
//	//for k := range tmpSubs {
//	//	policyTypes = append(policyTypes, k)
//	//}
//	//sort.Strings(policyTypes)
//	//
//	//return policyTypes
//	return nil
//}
//
//func (a1p *a1pController) HandleGetPoliciesTypeID(ctx context.Context, policyTypeID string) ([]store.A1PMValue, error) {
//	//policyEntries := []*a1pstore.Value{}
//	//policychEntries := make(chan *a1pstore.Entry)
//	//done := make(chan bool)
//	//
//	//go func(ch chan *a1pstore.Entry) {
//	//	for entry := range policychEntries {
//	//		value, ok := entry.Value.(a1pstore.Value)
//	//		if ok {
//	//			policyEntries = append(policyEntries, &value)
//	//		}
//	//	}
//	//	done <- true
//	//}(policychEntries)
//	//
//	//err := a1pstore.GetPoliciesByTypeID(ctx, a1p.policiesStore, policyTypeID, policychEntries)
//	//if err != nil {
//	//	close(policychEntries)
//	//	return policyEntries, err
//	//}
//	//
//	//<-done
//	//return policyEntries, nil
//
//	return nil, nil
//}
//
//func (a1p *a1pController) HandleGetPolicy(ctx context.Context, policyID, policyTypeID string) ([]store.A1PMValue, error) {
//	//a1pEntry, err := a1pstore.GetPolicyByID(ctx, a1p.policiesStore, policyID, policyTypeID)
//	//if err != nil {
//	//	return nil, err
//	//}
//	//
//	//a1pEntryValue := a1pEntry.Value.(a1pstore.Value)
//	//return &a1pEntryValue, nil
//	return nil, nil
//}
//
//func (a1p *a1pController) HandleGetPolicyStatus(ctx context.Context, policyID, policyTypeID string) (bool, error) {
//	//a1pEntry, err := a1pstore.GetPolicyByID(ctx, a1p.policiesStore, policyID, policyTypeID)
//	//if err != nil {
//	//	return false, err
//	//}
//	//
//	//a1pEntryValue := a1pEntry.Value.(a1pstore.Value)
//	//a1pEntryValueStatus := a1pEntryValue.PolicyStatus
//	//return a1pEntryValueStatus, nil
//	return false, nil
//}
//
//func (a1p *a1pController) HandlePolicyCreate(ctx context.Context, policyID, policyTypeID string, params map[string]string, policyObject map[string]interface{}) error {
//	//policyTypes := getSubscriptionPolicyTypes(ctx, a1p)
//	//
//	//if _, ok := policyTypes[policyTypeID]; !ok {
//	//	return fmt.Errorf("policyTypeID does not exist")
//	//}
//	//
//	//ch := make(chan *substore.Entry)
//	//err := substore.SubscriptionsByTypeID(ctx, a1p.subscriptionStore, substore.POLICY, policyTypeID, ch)
//	//if err != nil {
//	//	return err
//	//}
//	//
//	//policyTargets := make(map[string]a1pstore.PolicyTarget)
//	//
//	//policyStatus := true
//	//
//	//for subEntry := range ch {
//	//	subValue := subEntry.Value.(substore.Value)
//	//	subAddress := subValue.Client.Address
//	//
//	//	policyStatusTarget := a1psbi.CreatePolicy(ctx, subAddress, "", "", policyID, policyTypeID, policyObject)
//	//	if policyStatusTarget != nil {
//	//		policyStatus = false
//	//	}
//	//
//	//	policyTarget := a1pstore.PolicyTarget{
//	//		Address:            subAddress,
//	//		PolicyStatusObject: map[string]string{"status": policyStatusTarget.Error()},
//	//	}
//	//	policyTargets[subAddress] = policyTarget
//	//}
//	//
//	//a1pKey := a1pstore.Key{
//	//	PolicyId:     policyID,
//	//	PolicyTypeId: policyTypeID,
//	//}
//	//a1pValue := a1pstore.Value{
//	//	NotificationDestination: params["notificationDestination"],
//	//	PolicyObject:            policyObject,
//	//	Targets:                 policyTargets,
//	//	PolicyStatus:            policyStatus,
//	//}
//	//
//	//_, err = a1p.policiesStore.Put(ctx, a1pKey, a1pValue)
//	//if err != nil {
//	//	return err
//	//}
//
//	return nil
//}
//
//func (a1p *a1pController) HandlePolicyDelete(ctx context.Context, policyID, policyTypeID string) error {
//
//	//a1pEntry, err := a1pstore.GetPolicyByID(ctx, a1p.policiesStore, policyID, policyTypeID)
//	//if err != nil {
//	//	return err
//	//}
//	//
//	//a1pValue := a1pEntry.Value.(a1pstore.Value)
//	//
//	//for _, targetValue := range a1pValue.Targets {
//	//	err := a1psbi.DeletePolicy(ctx, targetValue.Address, "", "", policyID, policyTypeID)
//	//	if err != nil {
//	//		log.Warn(err)
//	//	}
//	//}
//	//
//	//err = a1p.policiesStore.Delete(ctx, a1pEntry.Key)
//	//if err != nil {
//	//	return err
//	//}
//
//	return nil
//}
//
//func (a1p *a1pController) HandlePolicyUpdate() error {
//	return nil
//}
//
//func getSubscriptionPolicyTypes(ctx context.Context, a1p *a1pController) map[string]struct{} {
//	//var exists = struct{}{}
//	//tmpSubs := make(map[string]struct{})
//	//ch := make(chan *substore.Entry)
//	//
//	//err := substore.SubscriptionsByType(ctx, a1p.subscriptionStore, substore.POLICY, ch)
//	//if err != nil {
//	//	return tmpSubs
//	//}
//	//
//	//for subEntry := range ch {
//	//	subValue := subEntry.Value.(substore.Value)
//	//	for _, sub := range subValue.Subscriptions {
//	//		if _, ok := tmpSubs[sub.TypeID]; !ok {
//	//			tmpSubs[sub.TypeID] = exists
//	//		}
//	//	}
//	//}
//	//
//	//return tmpSubs
//	return nil
//}
