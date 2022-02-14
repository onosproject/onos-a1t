// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0

package controller

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/onosproject/onos-a1t/pkg/stream"
	"github.com/onosproject/onos-api/go/onos/a1t/a1"
	"github.com/onosproject/onos-lib-go/pkg/errors"
	"reflect"
	"time"
)

func waitRespMsgWithTimer(id stream.ID, watcherID uuid.UUID, reqID string, respCh chan *stream.SBStreamMessage, outputCh chan interface{}, timeout time.Duration, streamBroker stream.Broker) {
	defer streamBroker.DeleteWatcher(id, watcherID)
	ctx, _ := context.WithTimeout(context.Background(), timeout)
	for {
		select {
		case resp := <-respCh:
			switch msg := resp.Payload.(type) {
			case *a1.PolicyResultMessage:
				if msg.Message.Header.RequestId == reqID {
					outputCh <- resp
					return
				}
			}
		case <-ctx.Done():
			outputCh <- errors.NewTimeout("Could not receive PolicyResultMessage in time (timer: %v)", TimeoutTimer)
			return
		}
	}
}

func checkOutput(output interface{}) (error, interface{}) {
	switch o := output.(type) {
	case error:
		return o, nil
	case *stream.SBStreamMessage:
		switch resp := o.Payload.(type) {
		case *a1.PolicyResultMessage:
			if !resp.Message.Result.Success {
				return fmt.Errorf(resp.Message.Result.Reason), nil
			}
			return nil, resp
		default:
			return errors.NewNotSupported("the response message %v should not come into A1T", reflect.TypeOf(resp)), nil
		}
	default:
		return errors.NewNotSupported("the response message %v should not come into A1T", reflect.TypeOf(o)), nil
	}
}
