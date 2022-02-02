// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-Only-1.0

package controller

import (
	"context"
	"github.com/onosproject/onos-a1t/pkg/rnib"
	"github.com/onosproject/onos-a1t/pkg/store"
	"github.com/onosproject/onos-a1t/pkg/stream"
	"time"
)

var TimeoutTimer = 5 * time.Second

type Broker interface {
	A1PController() A1PController
	A1EIController() A1EIController
	Run(ctx context.Context) error
}

func NewBroker(nonRTRICURL string, subscriptionStore store.Store, eijobsStore store.Store, rnibClient rnib.TopoClient, streamBroker stream.Broker) Broker {
	return &broker{
		a1pController:  NewA1PController(subscriptionStore, rnibClient, streamBroker),
		a1eiController: NewA1EIController(nonRTRICURL, subscriptionStore, eijobsStore, rnibClient, streamBroker),
		rnibClient:     rnibClient,
	}
}

type broker struct {
	a1pController  A1PController
	a1eiController A1EIController
	rnibClient     rnib.TopoClient
}

func (b *broker) Run(ctx context.Context) error {
	err := b.a1pController.Receiver(ctx)
	if err != nil {
		return err
	}
	//err = b.a1eiController.Receiver(ctx)
	//if err != nil {
	//	return err
	//}
	return nil
}

func (b *broker) A1PController() A1PController {
	return b.a1pController
}

func (b *broker) A1EIController() A1EIController {
	return b.a1eiController
}
