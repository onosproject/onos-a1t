// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-Only-1.0

package controller

import (
	a1eistore "github.com/onosproject/onos-a1t/pkg/store/a1ei"
	substore "github.com/onosproject/onos-a1t/pkg/store/subscription"
	// a1eisbi "github.com/onosproject/onos-a1t/pkg/southbound/a1ei"
	// a1einbi "github.com/onosproject/onos-a1t/pkg/northbound/a1ei/enrichment_information"
)

type A1EIController interface {
	HandleEIJobCreate() error
	HandleEIJobDelete() error
	HandleEIJobUpdate() error
	HandleGetEIJobTypes() []string
	HandleGetEIJobs() []string
	HandleGetEIJob() []string
}

type a1eiController struct {
	eijobsStore       a1eistore.Store
	subscriptionStore substore.Store
}

func NewA1EIController(subscriptionStore substore.Store, eijobsStore a1eistore.Store) A1EIController {
	return &a1eiController{
		eijobsStore:       eijobsStore,
		subscriptionStore: subscriptionStore,
	}
}

func (a1p *a1eiController) HandleEIJobCreate() error {
	// a1einbi.CreateEIJob()
	return nil
}

func (a1p *a1eiController) HandleEIJobDelete() error {
	return nil
}

func (a1p *a1eiController) HandleEIJobUpdate() error {
	return nil
}

func (a1p *a1eiController) HandleGetEIJobTypes() []string {
	// a1einbi.GetEIJobTypes()
	return []string{}
}

func (a1p *a1eiController) HandleGetEIJobs() []string {
	return []string{}
}

func (a1p *a1eiController) HandleGetEIJob() []string {
	return []string{}
}
