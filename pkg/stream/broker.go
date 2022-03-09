// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0

package stream

import (
	"context"
	"sync"

	"github.com/google/uuid"
	"github.com/onosproject/onos-lib-go/pkg/errors"
	"github.com/onosproject/onos-lib-go/pkg/logging"
)

var log = logging.GetLogger()

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
	log.Info("Print streams:")
	for k, v := range b.streams {
		log.Infof("stream key: %v, value: %v", k, v)
	}
	log.Info("Print watchers")
	for k, v := range b.watchers {
		log.Infof("watcher key: %v, value: %v", k, v)
	}
}

func (b *broker) AddStream(ctx context.Context, id ID) {
	log.Infof("Creating stream for %v", id)
	b.mu.Lock()
	defer b.mu.Unlock()
	_, ok := b.streams[id]
	if ok {
		log.Warnf("Stream for %v already exists", id)
		return
	}
	stream := NewDirectionalStream(id)
	b.streams[id] = stream
	b.watchers[id] = make(map[uuid.UUID]chan *SBStreamMessage)

	go func(m *sync.RWMutex) {
		for {
			msg, err := stream.Recv(ctx)
			if err != nil {
				log.Warnf("Forwarding channel closed: %v", err)
				return
			}
			m.Lock()
			log.Infof("watchers: %v", b.watchers)
			for _, v := range b.watchers[id] {
				log.Infof("Send %v to watcher %v", msg, v)
				select {
				case v <- msg:
					log.Infof("Sent %v to watcher %v", msg, v)
				default:
					log.Infof("Failed to send %v on %v", msg, v)
				}
			}
			m.Unlock()
		}
	}(&b.mu)
}

func (b *broker) Close(id ID) {
	log.Infof("Closing stream id %v", id)
	b.mu.Lock()
	defer b.mu.Unlock()
	stream, ok := b.streams[id]
	if !ok {
		log.Warnf("Stream for SID %v not found", id)
		return
	}
	stream.Close()
	delete(b.streams, id)
	delete(b.watchers, id)
}

func (b *broker) Send(id ID, message *SBStreamMessage) error {
	log.Infof("Sending message id: %v", id)
	b.mu.RLock()
	defer b.mu.RUnlock()
	log.Infof("Start Sending message id: %v", id)
	return b.streams[id].Send(message)
}

func (b *broker) Watch(id ID, ch chan *SBStreamMessage, watcherID uuid.UUID) error {
	log.Infof("Watching message id: %v", id)
	log.Infof("Add watcher ID %v: %v", watcherID, id)
	b.mu.Lock()
	defer b.mu.Unlock()
	if _, ok := b.streams[id]; !ok {
		return errors.NewNotFound("stream ID %v not found", id)
	}
	b.watchers[id][watcherID] = ch
	return nil
}

func (b *broker) DeleteWatcher(id ID, watcherID uuid.UUID) {
	log.Infof("deleting watcher ID %v: watcher ID %v", id, watcherID)
	b.mu.Lock()
	defer b.mu.Unlock()
	log.Infof("Delete watcherID: %v, watchers", watcherID, b.watchers)
	close(b.watchers[id][watcherID])
	delete(b.watchers[id], watcherID)
	log.Infof("Deleted watcherID: %v, watchers", watcherID, b.watchers)
}
