// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0

package subscription

import (
	"context"
	"github.com/onosproject/onos-a1t/pkg/rnib"
	"github.com/onosproject/onos-a1t/pkg/store"

	topoapi "github.com/onosproject/onos-api/go/onos/topo"
	"github.com/onosproject/onos-lib-go/pkg/logging"
)

var log = logging.GetLogger("subscription")

type Manager struct {
	subscriptionStore store.Store
	policiesStore     store.Store
	eiJobsStore       store.Store
	rnibClient        rnib.TopoClient
}

func NewSubscriptionManager(subscriptionStore store.Store, policiesStore store.Store, eiJobsStore store.Store) (*Manager, error) {
	rnibClient, err := rnib.NewClient()
	if err != nil {
		return &Manager{}, err
	}

	return &Manager{
		subscriptionStore: subscriptionStore,
		policiesStore:     policiesStore,
		eiJobsStore:       eiJobsStore,
		rnibClient:        rnibClient,
	}, nil
}

func (sm *Manager) Start() error {
	log.Info("Start SubscriptionManager")

	go func() {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		err := sm.watchXAppChanges(ctx)
		if err != nil {
			return
		}
	}()
	return nil
}

func (sm *Manager) watchXAppChanges(ctx context.Context) error {
	ch := make(chan topoapi.Event)
	err := sm.rnibClient.WatchTopoXapps(ctx, ch)
	if err != nil {
		log.Warn(err)
		return err
	}

	for topoEvent := range ch {
		log.Debugf("Received topo event: %v", topoEvent)
		if topoEvent.Object.GetEntity().GetKindID() == topoapi.XAPP {
			if topoEvent.Type == topoapi.EventType_ADDED || topoEvent.Type == topoapi.EventType_NONE {
				log.Info("xApp topo object added")
				err = sm.createSubscription(ctx, topoEvent.Object)
				// todo: add health check logic
				if err != nil {
					log.Error(err)
				}
			} else if topoEvent.Type == topoapi.EventType_REMOVED {
				log.Info("xApp topo object removed")
				err = sm.deleteSubscription(ctx, topoEvent.Object)
				if err != nil {
					log.Error(err)
				}
			} else if topoEvent.Type == topoapi.EventType_UPDATED {
				log.Info("xApp topo object updated")
				err = sm.updateSubscription(ctx, topoEvent.Object)
				if err != nil {
					log.Error(err)
				}
			}
		}
	}

	return nil
}

func (sm *Manager) createSubscription(ctx context.Context, topoObject topoapi.Object) error {
	// add relation
	err := sm.rnibClient.AddA1TXappRelation(ctx, topoObject.GetID())
	if err != nil {
		return err
	}
	// get aspect and store it to local store
	xAppInfo, err := sm.rnibClient.GetXappAspects(ctx, topoObject.GetID())
	if err != nil {
		return err
	}

	subKey := store.SubscriptionKey{
		TargetXAppID: topoObject.GetID(),
	}
	subValue := &store.SubscriptionValue{
		A1ServiceCapabilities: make([]*store.A1ServiceType, 0),
	}

	// get endpoint information
	for _, i := range xAppInfo.GetInterfaces() {
		if i.GetType() == topoapi.Interface_INTERFACE_A1_XAPP {
			subValue.A1EndpointIP = i.GetIP()
			subValue.A1EndpointPort = i.GetPort()
		}
	}

	// get capabilities
	for _, p := range xAppInfo.GetA1PolicyTypes() {
		serviceTypeDef := &store.A1ServiceType{
			A1Service: store.PolicyManagement,
			TypeID:    string(p.GetID()),
		}
		subValue.A1ServiceCapabilities = append(subValue.A1ServiceCapabilities, serviceTypeDef)
	}

	// todo: have to be clarified but now by default it added the EI capability; at this moment, it's fine
	eiServiceTypeDef := &store.A1ServiceType{
		A1Service: store.EnrichmentInformation,
		TypeID:    "",
	}
	subValue.A1ServiceCapabilities = append(subValue.A1ServiceCapabilities, eiServiceTypeDef)

	_, err = sm.subscriptionStore.Put(ctx, subKey, subValue)
	if err != nil {
		return err
	}

	// add entry on a1pm and a1ei stores
	a1Key := store.A1Key{
		TargetXAppID: topoObject.GetID(),
	}
	a1PmValue := &store.A1PMValue{
		A1PolicyObjects: make(map[store.A1PolicyObjectID]store.A1ServiceType),
	}
	a1EiValue := &store.A1EIValue{
		A1EIJobObjects: make(map[store.A1EIJobObjectID]store.A1ServiceType),
	}

	_, err = sm.policiesStore.Put(ctx, a1Key, a1PmValue)
	if err != nil {
		return err
	}
	_, err = sm.eiJobsStore.Put(ctx, a1Key, a1EiValue)
	if err != nil {
		return err
	}

	return nil
}

func (sm *Manager) deleteSubscription(ctx context.Context, topoObject topoapi.Object) error {
	subKey := store.SubscriptionKey{
		TargetXAppID: topoObject.GetID(),
	}

	err := sm.subscriptionStore.Delete(ctx, subKey)
	if err != nil {
		return err
	}

	// delete entry on a1pm and a1ei stores
	a1Key := store.A1Key{
		TargetXAppID: topoObject.GetID(),
	}
	err = sm.policiesStore.Delete(ctx, a1Key)
	if err != nil {
		return err
	}
	err = sm.eiJobsStore.Delete(ctx, a1Key)
	if err != nil {
		return err
	}

	return nil
}

func (sm *Manager) updateSubscription(ctx context.Context, topoObject topoapi.Object) error {
	xAppInfo, err := sm.rnibClient.GetXappAspects(ctx, topoObject.GetID())
	if err != nil {
		return err
	}

	// if topo has updated to have no interfaces, this means that xApp is terminated; subscription should be removed
	if len(xAppInfo.GetInterfaces()) == 0 {
		err = sm.deleteSubscription(ctx, topoObject)
		if err != nil {
			return err
		}
	}

	subKey := store.SubscriptionKey{
		TargetXAppID: topoObject.GetID(),
	}
	subValue := &store.SubscriptionValue{
		A1ServiceCapabilities: make([]*store.A1ServiceType, 0),
	}

	// get endpoint information
	for _, i := range xAppInfo.GetInterfaces() {
		if i.GetType() == topoapi.Interface_INTERFACE_A1_XAPP {
			subValue.A1EndpointIP = i.GetIP()
			subValue.A1EndpointPort = i.GetPort()
		}
	}

	// get capabilities
	for _, p := range xAppInfo.GetA1PolicyTypes() {
		serviceTypeDef := &store.A1ServiceType{
			A1Service: store.PolicyManagement,
			TypeID:    string(p.GetID()),
		}
		subValue.A1ServiceCapabilities = append(subValue.A1ServiceCapabilities, serviceTypeDef)
	}

	// todo: have to be clarified but now by default it added the EI capability; at this moment, it's fine
	eiServiceTypeDef := &store.A1ServiceType{
		A1Service: store.EnrichmentInformation,
		TypeID:    "",
	}
	subValue.A1ServiceCapabilities = append(subValue.A1ServiceCapabilities, eiServiceTypeDef)

	_, err = sm.subscriptionStore.Update(ctx, subKey, subValue)
	if err != nil {
		return err
	}

	return nil
}
