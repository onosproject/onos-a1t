// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-Only-1.0

package a1ei

import (
	"context"
	"fmt"
	"sync"

	"github.com/google/uuid"
	"github.com/onosproject/onos-lib-go/pkg/errors"
	"github.com/onosproject/onos-lib-go/pkg/logging"

	"github.com/onosproject/onos-a1t/pkg/store"
)

var log = logging.GetLogger("store", "a1ei")

// Store A1 EI jobs store interface
type Store interface {
	Put(ctx context.Context, key Key, value interface{}) (*Entry, error)

	Get(ctx context.Context, key Key) (*Entry, error)

	Delete(ctx context.Context, key Key) error

	Entries(ctx context.Context, ch chan<- *Entry) error

	Watch(ctx context.Context, ch chan<- store.Event) error
}

type EIJobTarget struct {
	Address           string
	EIJobStatusObject map[string]string
}

type Value struct {
	NotificationDestination string
	EIJobObject             map[string]string
	EIJobStatusObjects      map[string]string
	Targets                 map[string]EIJobTarget
	EIJobStatus             bool
}

type Key struct {
	EIJobID   string
	EIJobtype string
}

type Entry struct {
	Key   Key
	Value interface{}
}

type a1eistore struct {
	measurements map[Key]*Entry
	mu           sync.RWMutex
	watchers     *store.Watchers
}

// NewStore creates new store
func NewStore() Store {
	watchers := store.NewWatchers()
	return &a1eistore{
		measurements: make(map[Key]*Entry),
		watchers:     watchers,
	}
}

func (s *a1eistore) Entries(ctx context.Context, ch chan<- *Entry) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if len(s.measurements) == 0 {
		close(ch)
		return fmt.Errorf("no measurements entries stored")
	}

	for _, entry := range s.measurements {
		ch <- entry
	}

	close(ch)
	return nil
}

func (s *a1eistore) Delete(ctx context.Context, key Key) error {
	// TODO check the key and make sure it is not empty
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.measurements, key)
	return nil

}

func (s *a1eistore) Put(ctx context.Context, key Key, value interface{}) (*Entry, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	entry := &Entry{
		Key:   key,
		Value: value,
	}
	s.measurements[key] = entry
	s.watchers.Send(store.Event{
		Key:   key,
		Value: entry,
		Type:  store.Created,
	})
	return entry, nil

}

func (s *a1eistore) Get(ctx context.Context, key Key) (*Entry, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if v, ok := s.measurements[key]; ok {
		return v, nil
	}
	return nil, errors.New(errors.NotFound, "the measurement entry does not exist")
}

func (s *a1eistore) Watch(ctx context.Context, ch chan<- store.Event) error {
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

var _ Store = &a1eistore{}
