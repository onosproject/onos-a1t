// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-Only-1.0

package southbound

import (
	"context"
	sbclient "github.com/onosproject/onos-a1t/pkg/southbound/client"
	"github.com/onosproject/onos-a1t/pkg/store"
	"github.com/onosproject/onos-a1t/pkg/stream"
	"github.com/onosproject/onos-lib-go/pkg/logging"
	"sync"
)

var log = logging.GetLogger("southbound")

func NewSouthboundManager(broker stream.Broker, subStore store.Store) Manager {
	return &manager{
		streamBroker: broker,
		a1pClients:   make(map[string]sbclient.Client),
		a1eiClients:  make(map[string]sbclient.Client),
		subStore:     subStore,
	}
}

type Manager interface {
	Run(ctx context.Context) error
	Close(xAppID string, a1Service stream.A1Service)
}

type manager struct {
	streamBroker stream.Broker
	a1pClients   map[string]sbclient.Client
	a1eiClients  map[string]sbclient.Client
	subStore     store.Store
	clientMu     sync.RWMutex
}

func (m *manager) Close(xAppID string, a1Service stream.A1Service) {
	switch a1Service {
	case stream.PolicyManagement:
		log.Infof("Closing A1 policy management southbound client for xApp ID %v", xAppID)
		m.clientMu.Lock()
		delete(m.a1pClients, xAppID)
		m.clientMu.Unlock()
	case stream.EnrichmentInformation:
		log.Infof("Closing A1 EI southbound client for xApp ID %v", xAppID)
		m.clientMu.Lock()
		delete(m.a1eiClients, xAppID)
		m.clientMu.Unlock()
	}
}

func (m *manager) Run(ctx context.Context) error {
	log.Info("Run southbound manager")
	return m.watchSubStore(ctx)
}

func (m *manager) watchSubStore(ctx context.Context) error {
	log.Info("Start watching subscription store at southbound manager")
	ch := make(chan store.Event)
	go m.subStoreListener(ctx, ch)
	err := m.subStore.Watch(ctx, ch)
	if err != nil {
		close(ch)
		log.Error(err)
		return err
	}
	return nil
}

func (m *manager) subStoreListener(ctx context.Context, ch chan store.Event) {
	var err error
	for e := range ch {
		entry := e.Value.(*store.Entry)
		switch e.Type {
		case store.Created:
			err = m.createEventSubStoreHandler(ctx, entry)
			if err != nil {
				log.Warn(err)
			}
		case store.Updated:
			err = m.createEventSubStoreHandler(ctx, entry)
			if err != nil {
				log.Warn(err)
			}
		case store.Deleted:
			err = m.deleteEventSubStoreHandler(ctx, entry)
			if err != nil {
				log.Warn(err)
			}
		}
	}
}

func (m *manager) createEventSubStoreHandler(ctx context.Context, entry *store.Entry) error {
	log.Infof("Subscription store entry %v was just created or updated", *entry)
	key := entry.Key.(store.SubscriptionKey)
	value := entry.Value.(*store.SubscriptionValue)
	m.clientMu.Lock()
	defer m.clientMu.Unlock()
	// todo: currently, a1ei is default session. If not for the future, it should be optional
	if _, ok := m.a1eiClients[string(key.TargetXAppID)]; !ok {
		// call NewClient function
		a1eiClient, err := sbclient.NewA1EIClient(ctx, string(key.TargetXAppID), value.A1EndpointIP, value.A1EndpointPort, m.streamBroker)
		if err != nil {
			return err
		}
		// store the created client to the map
		m.a1eiClients[string(key.TargetXAppID)] = a1eiClient
		go func() {
			err = a1eiClient.Run(ctx)
			if err != nil {
				log.Warn(err)
			}
		}()
	}
	for _, c := range value.A1ServiceCapabilities {
		switch c.A1Service {
		case store.PolicyManagement:
			if _, ok := m.a1pClients[string(key.TargetXAppID)]; ok {
				break
			}
			// call NewClient function
			a1pClient, err := sbclient.NewA1PClient(ctx, string(key.TargetXAppID), value.A1EndpointIP, value.A1EndpointPort, m.streamBroker)
			if err != nil {
				return err
			}
			// store the created client to the map
			m.a1pClients[string(key.TargetXAppID)] = a1pClient
			go func() {
				err = a1pClient.Run(ctx)
				if err != nil {
					log.Warn(err)
				}
			}()
		}
	}
	return nil
}

func (m *manager) deleteEventSubStoreHandler(ctx context.Context, entry *store.Entry) error {
	log.Infof("Subscription store entry %v was just deleted", *entry)
	key := entry.Key.(store.SubscriptionKey)
	if _, ok := m.a1eiClients[string(key.TargetXAppID)]; ok {
		m.Close(string(key.TargetXAppID), stream.EnrichmentInformation)
	}
	if _, ok := m.a1pClients[string(key.TargetXAppID)]; ok {
		m.Close(string(key.TargetXAppID), stream.PolicyManagement)
	}

	return nil
}
