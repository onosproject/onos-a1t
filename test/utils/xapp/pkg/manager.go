// SPDX-FileCopyrightText: 2022-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0

package xapp

import (
	"context"

	"github.com/onosproject/onos-api/go/onos/topo"
	"github.com/onosproject/onos-lib-go/pkg/logging"
	"github.com/onosproject/onos-lib-go/pkg/northbound"
)

var log = logging.GetLogger("a1t-xapp")

// Config is a manager configuration
type Config struct {
	CAPath      string
	KeyPath     string
	CertPath    string
	E2tEndpoint string
	GRPCPort    int
	RicActionID int32
	ConfigPath  string
}

// NewManager generates the new a1txapp manager
func NewManager(config Config) (*Manager, error) {

	a1PolicyTypes := make([]*topo.A1PolicyType, 0)
	a1Policy := &topo.A1PolicyType{
		Name:        "ORAN_TrafficSteeringPreference",
		Version:     "2.0.0",
		ID:          "ORAN_TrafficSteeringPreference_2.0.0",
		Description: "O-RAN traffic steering",
	}
	a1PolicyTypes = append(a1PolicyTypes, a1Policy)

	a1Manager, err := NewNBIManager(config.CAPath, config.KeyPath, config.CertPath, config.GRPCPort, "onos-a1txapp", a1PolicyTypes)
	if err != nil {
		log.Warn(err)
	}

	manager := &Manager{
		config:    config,
		a1Manager: *a1Manager,
	}
	return manager, nil
}

// Manager is an abstract struct for manager
type Manager struct {
	config    Config
	a1Manager NBIManager
}

// Run runs a1txapp manager
func (m *Manager) Run() {
	err := m.start()
	if err != nil {
		log.Errorf("Error when starting a1txapp: %v", err)
	}
}

// Close closes manager
func (m *Manager) Close() {
	log.Info("closing Manager")
	m.a1Manager.Close(context.Background())
}

func (m *Manager) start() error {
	err := m.startNorthboundServer()
	if err != nil {
		log.Warn(err)
		return err
	}

	m.a1Manager.Start()

	return nil
}

func (m *Manager) startNorthboundServer() error {
	s := northbound.NewServer(northbound.NewServerCfg(
		m.config.CAPath,
		m.config.KeyPath,
		m.config.CertPath,
		int16(m.config.GRPCPort),
		true,
		northbound.SecurityConfig{}))

	s.AddService(NewA1EIService())
	s.AddService(NewA1PService())

	doneCh := make(chan error)
	go func() {
		err := s.Serve(func(started string) {
			log.Info("Started NBI on ", started)
			close(doneCh)
		})
		if err != nil {
			doneCh <- err
		}
	}()
	return <-doneCh
}
