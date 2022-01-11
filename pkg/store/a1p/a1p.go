// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-Only-1.0

package a1p

import (
	"context"
	"fmt"
	"sync"

	"github.com/google/uuid"
	"github.com/onosproject/onos-lib-go/pkg/errors"
	"github.com/onosproject/onos-lib-go/pkg/logging"

	"github.com/onosproject/onos-a1t/pkg/store"
)

var log = logging.GetLogger("store", "a1p")

// Store a1 policies store interface
type Store interface {
	Put(ctx context.Context, key Key, value interface{}) (*Entry, error)

	Get(ctx context.Context, key Key) (*Entry, error)

	Delete(ctx context.Context, key Key) error

	Entries(ctx context.Context, ch chan<- *Entry) error

	Watch(ctx context.Context, ch chan<- store.Event) error
}

type PolicyTarget struct {
	Address            string
	PolicyStatusObject map[string]string
}

type Key struct {
	PolicyId     string
	PolicyTypeId string
}

type Value struct {
	NotificationDestination string
	PolicyObject            map[string]interface{}
	PolicyStatusObjects     map[string]string
	Targets                 map[string]PolicyTarget
	PolicyStatus            bool
}

type Entry struct {
	Key   Key
	Value interface{}
}

type a1pStore struct {
	policies map[Key]*Entry
	mu       sync.RWMutex
	watchers *store.Watchers
}

// NewStore creates new store for A1P
func NewStore() Store {
	watchers := store.NewWatchers()
	return &a1pStore{
		policies: make(map[Key]*Entry),
		watchers: watchers,
	}
}

func (s *a1pStore) Entries(ctx context.Context, ch chan<- *Entry) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if len(s.policies) == 0 {
		close(ch)
		return fmt.Errorf("no policies entries stored")
	}

	for _, entry := range s.policies {
		ch <- entry
	}

	close(ch)
	return nil
}

func (s *a1pStore) Delete(ctx context.Context, key Key) error {
	// TODO check the key and make sure it is not empty
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.policies, key)
	return nil

}

func (s *a1pStore) Put(ctx context.Context, key Key, value interface{}) (*Entry, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	entry := &Entry{
		Key:   key,
		Value: value,
	}
	s.policies[key] = entry
	s.watchers.Send(store.Event{
		Key:   key,
		Value: entry,
		Type:  store.Created,
	})
	return entry, nil

}

func (s *a1pStore) Get(ctx context.Context, key Key) (*Entry, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if v, ok := s.policies[key]; ok {
		return v, nil
	}
	return nil, errors.New(errors.NotFound, "the policy entry does not exist")
}

func (s *a1pStore) Watch(ctx context.Context, ch chan<- store.Event) error {
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

var _ Store = &a1pStore{}
