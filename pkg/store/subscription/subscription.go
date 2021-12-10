// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-Only-1.0

package subscription

import (
	"context"
	"fmt"
	"sync"

	"github.com/google/uuid"
	"github.com/onosproject/onos-lib-go/pkg/errors"
	"github.com/onosproject/onos-lib-go/pkg/logging"

	"github.com/onosproject/onos-a1t/pkg/store"
)

var log = logging.GetLogger("store", "subscription")

type Store interface {
	Put(ctx context.Context, key Key, value interface{}) (*Entry, error)

	// Get gets a metric store entry based on a given key
	Get(ctx context.Context, key Key) (*Entry, error)

	// Delete deletes an entry based on a given key
	Delete(ctx context.Context, key Key) error

	// Entries list all of the metric store entries
	Entries(ctx context.Context, ch chan<- *Entry) error

	// Watch measurement store changes
	Watch(ctx context.Context, ch chan<- store.Event) error
}

type Client struct {
	Address  string
	CertPath string
	KeyPath  string
}

type Subscription struct {
	Types []string
	ID    string
}

type Key struct {
	TargetID string
}

type Value struct {
	Client       Client
	Subscription Subscription
}

type Entry struct {
	Key   Key
	Value interface{}
}

func NewSubscriptionKey(targetID string) Key {
	return Key{
		TargetID: targetID,
	}
}

type subscriptionstore struct {
	subscriptions map[Key]*Entry
	mu            sync.RWMutex
	watchers      *store.Watchers
}

// NewStore creates new store
func NewStore() Store {
	watchers := store.NewWatchers()
	return &subscriptionstore{
		subscriptions: make(map[Key]*Entry),
		watchers:      watchers,
	}
}

func (s *subscriptionstore) Entries(ctx context.Context, ch chan<- *Entry) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if len(s.subscriptions) == 0 {
		close(ch)
		return fmt.Errorf("no subscriptions entries stored")
	}

	for _, entry := range s.subscriptions {
		ch <- entry
	}

	close(ch)
	return nil
}

func (s *subscriptionstore) Delete(ctx context.Context, key Key) error {
	// TODO check the key and make sure it is not empty
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.subscriptions, key)
	return nil

}

func (s *subscriptionstore) Put(ctx context.Context, key Key, value interface{}) (*Entry, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	entry := &Entry{
		Key:   key,
		Value: value,
	}
	s.subscriptions[key] = entry
	s.watchers.Send(store.Event{
		Key:   key,
		Value: entry,
		Type:  store.Created,
	})
	return entry, nil

}

func (s *subscriptionstore) Get(ctx context.Context, key Key) (*Entry, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if v, ok := s.subscriptions[key]; ok {
		return v, nil
	}
	return nil, errors.New(errors.NotFound, "the measurement entry does not exist")
}

func (s *subscriptionstore) Watch(ctx context.Context, ch chan<- store.Event) error {
	id := uuid.New()
	err := s.watchers.AddWatcher(id, ch)
	if err != nil {
		log.Error(err)
		close(ch)
		return err
	}
	go func() {
		<-ctx.Done()
		err = s.watchers.RemoveWatcher(id)
		if err != nil {
			log.Error(err)
		}
		close(ch)
	}()
	return nil
}

var _ Store = &subscriptionstore{}
