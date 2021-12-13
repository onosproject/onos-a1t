// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-Only-1.0

package controller

import (
	a1pstore "github.com/onosproject/onos-a1t/pkg/store/a1p"
	substore "github.com/onosproject/onos-a1t/pkg/store/subscription"
	// a1psbi "github.com/onosproject/onos-a1t/pkg/southbound/a1p"
	// a1pnbi "github.com/onosproject/onos-a1t/pkg/northbound/a1ap/policy_management"
)

type A1PController interface {
	HandlePolicyCreate() error
	HandlePolicyDelete() error
	HandlePolicyUpdate() error
	HandleGetPolicyTypes() []string
	HandleGetPolicies() []string
	HandleGetPolicy() []string
}

type a1pController struct {
	policiesStore     a1pstore.Store
	subscriptionStore substore.Store
}

func NewA1PController(subscriptionStore substore.Store, policiesStore a1pstore.Store) A1PController {
	return &a1pController{
		policiesStore:     policiesStore,
		subscriptionStore: subscriptionStore,
	}
}

func (a1p *a1pController) HandlePolicyCreate() error {
	// a1psbi.CreatePolicy()
	return nil
}

func (a1p *a1pController) HandlePolicyDelete() error {
	return nil
}

func (a1p *a1pController) HandlePolicyUpdate() error {
	return nil
}

func (a1p *a1pController) HandleGetPolicyTypes() []string {
	// subscriptionStore.Entries()
	return []string{}
}

func (a1p *a1pController) HandleGetPolicies() []string {
	return []string{}
}

func (a1p *a1pController) HandleGetPolicy() []string {
	return []string{}
}
