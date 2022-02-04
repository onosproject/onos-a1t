// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0

package stream

import (
	"context"
	"github.com/onosproject/onos-lib-go/pkg/errors"
	"github.com/onosproject/onos-lib-go/pkg/logging"
	"io"
	"sync"
	"time"
)

var logStream = logging.GetLogger("stream")

var SendTimeout = time.Second * 5

type Reader interface {
	Recv(ctx context.Context) (*SBStreamMessage, error)
}

type Writer interface {
	Close()
	Send(message *SBStreamMessage) error
}

type EndpointID string

type ID struct {
	SrcEndpointID  EndpointID
	DestEndpointID EndpointID
}

type IO interface {
	GetID() ID
	GetSrcEndpointID() EndpointID
	GetDestEndpointID() EndpointID
}

type Stream interface {
	IO
	Reader
	Writer
}

type directionalStreamIO struct {
	id ID
}

func NewDirectionalStreamWriter(ch chan *SBStreamMessage) Writer {
	return &directionalStreamWriter{
		ch:     ch,
		closed: false,
	}
}

type directionalStreamWriter struct {
	ch     chan *SBStreamMessage
	closed bool
	mu     sync.RWMutex
}

func (d *directionalStreamWriter) Close() {
	logStream.Infof("Deleting stream")
	d.mu.Lock()
	defer d.mu.Unlock()
	d.closed = true
	close(d.ch)
}

func (d *directionalStreamWriter) Send(message *SBStreamMessage) error {
	d.mu.RLock()
	defer d.mu.RUnlock()
	if d.closed {
		return io.EOF
	}
	select {
	case d.ch <- message:
		return nil
	case <-time.After(SendTimeout):
		return errors.NewTimeout("Failed to send message before send timer expired")
	}
}

func NewDirectionalStreamReader(ch chan *SBStreamMessage) Reader {
	return &directionalStreamReader{
		ch: ch,
	}
}

type directionalStreamReader struct {
	ch chan *SBStreamMessage
}

func (d *directionalStreamReader) Recv(ctx context.Context) (*SBStreamMessage, error) {
	select {
	case message, ok := <-d.ch:
		if !ok {
			return nil, io.EOF
		}
		return message, nil
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

func NewDirectionalStream(id ID) Stream {
	ch := make(chan *SBStreamMessage)
	return &directionalStream{
		IO: &directionalStreamIO{
			id: id,
		},
		Reader: NewDirectionalStreamReader(ch),
		Writer: NewDirectionalStreamWriter(ch),
	}
}

type directionalStream struct {
	IO
	Reader
	Writer
}

func (d *directionalStreamIO) GetID() ID {
	return d.id
}

func (d *directionalStreamIO) GetSrcEndpointID() EndpointID {
	return d.id.SrcEndpointID
}

func (d *directionalStreamIO) GetDestEndpointID() EndpointID {
	return d.id.DestEndpointID
}

var _ Stream = &directionalStream{}
