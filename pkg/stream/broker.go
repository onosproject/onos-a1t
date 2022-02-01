// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-Only-1.0

package stream

import (
	"context"
	"github.com/google/uuid"
	"github.com/onosproject/onos-lib-go/pkg/errors"
	"github.com/onosproject/onos-lib-go/pkg/logging"
	"sync"
)

var logBroker = logging.GetLogger("stream", "broker")

type Broker interface {
	Close(id ID)
	AddStream(ctx context.Context, id ID)
	Send(id ID, message *SBStreamMessage) error
	Watch(id ID, ch chan *SBStreamMessage, watcherID uuid.UUID) error
	DeleteWatcher(id ID, watcherID uuid.UUID)
	Print()
}

func NewBroker() Broker {
	streams := make(map[ID]Stream)
	watchers := make(map[ID]map[uuid.UUID]chan *SBStreamMessage)
	return &broker{
		streams:  streams,
		watchers: watchers,
	}
}

type broker struct {
	streams  map[ID]Stream
	watchers map[ID]map[uuid.UUID]chan *SBStreamMessage
	mu       sync.RWMutex
}

func (b *broker) Print() {
	b.mu.Lock()
	defer b.mu.Unlock()
	logBroker.Info("Print streams:")
	for k, v := range b.streams {
		logBroker.Infof("stream key: %v, value: %v", k, v)
	}
	logBroker.Info("Print watchers")
	for k, v := range b.watchers {
		logBroker.Infof("watcher key: %v, value: %v", k, v)
	}
}

func (b *broker) AddStream(ctx context.Context, id ID) {
	logBroker.Infof("Creating stream for %v", id)
	b.mu.Lock()
	defer b.mu.Unlock()
	_, ok := b.streams[id]
	if ok {
		logBroker.Warnf("Stream for %v already exists", id)
	}
	stream := NewDirectionalStream(id)
	b.streams[id] = stream
	b.watchers[id] = make(map[uuid.UUID]chan *SBStreamMessage)

	go func() {
		for {
			msg, err := stream.Recv(ctx)
			if err != nil {
				logBroker.Warnf("Forwarding channel closed: %v", err)
				return
			}
			b.mu.Lock()
			logBroker.Info("watchers: %v", b.watchers)
			for _, v := range b.watchers[id] {
				v <- msg
			}
			b.mu.Unlock()
		}
	}()
}

func (b *broker) Close(id ID) {
	logBroker.Infof("Closing stream id %v", id)
	b.mu.Lock()
	defer b.mu.Unlock()
	stream, ok := b.streams[id]
	if !ok {
		logBroker.Warnf("Stream for SID %v not found", id)
		return
	}
	stream.Close()
	delete(b.streams, id)
	delete(b.watchers, id)
}

func (b *broker) Send(id ID, message *SBStreamMessage) error {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.streams[id].Send(message)
}

func (b *broker) Watch(id ID, ch chan *SBStreamMessage, watcherID uuid.UUID) error {
	b.mu.Lock()
	defer b.mu.Unlock()
	if _, ok := b.streams[id]; !ok {
		return errors.NewNotFound("stream ID %v not found", id)
	}
	b.watchers[id][watcherID] = ch
	return nil
}

func (b *broker) DeleteWatcher(id ID, watcherID uuid.UUID) {
	b.mu.Lock()
	defer b.mu.Unlock()
	logBroker.Infof("Delete watcherID: %v, watchers", watcherID, b.watchers)
	delete(b.watchers[id], watcherID)
	logBroker.Infof("Deleted watcherID: %v, watchers", watcherID, b.watchers)
}
