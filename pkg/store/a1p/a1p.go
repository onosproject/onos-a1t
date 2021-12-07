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

// A1PStore a1 policies store interface
type Store interface {
	Put(ctx context.Context, key A1PKey, value interface{}) (*A1PEntry, error)

	Get(ctx context.Context, key A1PKey) (*A1PEntry, error)

	Delete(ctx context.Context, key A1PKey) error

	Entries(ctx context.Context, ch chan<- *A1PEntry) error

	Watch(ctx context.Context, ch chan<- store.Event) error
}

type PolicyTarget struct {
	TargetID           string
	PolicyStatusObject map[string]interface{}
}

type A1PKey struct {
	PolicyId     string
	PolicyTypeId string
}

type A1PValue struct {
	NotificationDestination string
	PolicyObject            map[string]interface{}
	PolicyStatusObjects     map[string]interface{}
	Targets                 map[string]PolicyTarget
}

type A1PEntry struct {
	Key   A1PKey
	Value interface{}
}

func NewA1PKey(policyID, policyTypeId string) *A1PKey {
	return &A1PKey{
		PolicyId:     policyID,
		PolicyTypeId: policyTypeId,
	}
}

type a1pStore struct {
	policies map[A1PKey]*A1PEntry
	mu       sync.RWMutex
	watchers *store.Watchers
}

// NewStore creates new store for A1P
func NewStore() Store {
	watchers := store.NewWatchers()
	return &a1pStore{
		policies: make(map[A1PKey]*A1PEntry),
		watchers: watchers,
	}
}

func (s *a1pStore) Entries(ctx context.Context, ch chan<- *A1PEntry) error {
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

func (s *a1pStore) Delete(ctx context.Context, key A1PKey) error {
	// TODO check the key and make sure it is not empty
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.policies, key)
	return nil

}

func (s *a1pStore) Put(ctx context.Context, key A1PKey, value interface{}) (*A1PEntry, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	entry := &A1PEntry{
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

func (s *a1pStore) Get(ctx context.Context, key A1PKey) (*A1PEntry, error) {
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
