// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0

package nonrtric

import (
	"github.com/onosproject/onos-lib-go/pkg/logging"

	"github.com/onosproject/onos-a1t/pkg/store"
)

var log = logging.GetLogger("test-manager")

type Manager struct {
	restserver  RestServer
	policyStore store.Store
	eijobsStore store.Store
	controller  Controller
}

func NewManager(baseURL, nearRTRicBaseURL string) (*Manager, error) {

	policyStore := store.NewStore()
	eijobsStore := store.NewStore()
	controller := NewController(nearRTRicBaseURL, policyStore, eijobsStore)

	rest, err := NewRestServer(baseURL, controller)
	if err != nil {
		log.Info(err)
	}

	mngr := &Manager{
		restserver:  rest,
		policyStore: policyStore,
		eijobsStore: eijobsStore,
		controller:  controller,
	}

	return mngr, nil
}

func (m *Manager) start() error {
	go m.restserver.Start()
	return nil
}

func (m *Manager) Run() {
	err := m.start()
	if err != nil {
		log.Errorf("Error when starting A1T test-manager: %v", err)
	}
}

func (m *Manager) GetController() Controller {
	return m.controller
}
