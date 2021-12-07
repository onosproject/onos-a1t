// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-Only-1.0

package subscription

import (
	"context"

	topoapi "github.com/onosproject/onos-api/go/onos/topo"
	"github.com/onosproject/onos-lib-go/pkg/logging"

	controller "github.com/onosproject/onos-a1t/pkg/controller"
	substore "github.com/onosproject/onos-a1t/pkg/store/subscription"
)

var log = logging.GetLogger("subscription")

type SubscriptionManager struct {
	broker            controller.Broker
	subscriptionStore substore.Store
	rnibClient        Client
}

func NewSubscriptionManager(broker controller.Broker, subscriptionStore substore.Store) (*SubscriptionManager, error) {
	rnibClient, err := NewClient()
	if err != nil {
		return &SubscriptionManager{}, err
	}

	return &SubscriptionManager{
		broker:            broker,
		subscriptionStore: subscriptionStore,
		rnibClient:        rnibClient,
	}, nil
}

func (sm *SubscriptionManager) Start() error {
	log.Info("Start SubscriptionManager")

	go func() {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		err := sm.watchXappChanges(ctx)
		if err != nil {
			return
		}
	}()
	return nil
}

func (sm *SubscriptionManager) watchXappChanges(ctx context.Context) error {
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
				log.Info("xApp Added")
				//TODO Create xapp subscription and get xApp Aspects to get "interests" in a1p and/or a1ei

			} else if topoEvent.Type == topoapi.EventType_REMOVED {
				log.Info("xApp Removed")
				//TODO Get all xapp subscriptions and delete them (handle status of a1p and a1ei)
			}
		}
	}

	return nil
}

func (sm *SubscriptionManager) createSubscription(ctx context.Context, xappinfo topoapi.XAppInfo) error {

	subKey := substore.NewSubscriptionKey(xappinfo.String())

	policyTypes := []string{}
	a1polTypes := xappinfo.GetA1PolicyTypes()

	for _, pt := range a1polTypes {
		policyTypes = append(policyTypes, pt.String())
	}

	subValue := substore.SubscriptionValue{
		Client: substore.Client{},
		Subscription: substore.Subscription{
			Types: policyTypes,
		},
	}

	_, err := sm.subscriptionStore.Put(ctx, subKey, subValue)
	if err != nil {
		log.Warn(err)
		return err
	}

	return nil
}

func (sm *SubscriptionManager) updateSubscription() {}

func (sm *SubscriptionManager) deleteSubscription() {}
