// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-Only-1.0

package controller

import (
	a1eistore "github.com/onosproject/onos-a1t/pkg/store/a1ei"
	a1pstore "github.com/onosproject/onos-a1t/pkg/store/a1p"
	substore "github.com/onosproject/onos-a1t/pkg/store/subscription"
)

type Broker interface {
	A1PController() A1PController
	A1EIController() A1EIController
}

func NewBroker(subscriptionStore substore.Store, policiesStore a1pstore.Store, eijobsStore a1eistore.Store) Broker {
	return &broker{
		a1pController:  NewA1PController(subscriptionStore, policiesStore),
		a1eiController: NewA1EIController(subscriptionStore, eijobsStore),
	}
}

type broker struct {
	a1pController  A1PController
	a1eiController A1EIController
}

func (b *broker) A1PController() A1PController {
	return b.a1pController
}

func (b *broker) A1EIController() A1EIController {
	return b.a1eiController
}
