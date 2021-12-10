// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-Only-1.0

package manager

import (
	"github.com/onosproject/onos-a1t/pkg/controller"
	nbi "github.com/onosproject/onos-a1t/pkg/northbound/cli"
	nbirest "github.com/onosproject/onos-a1t/pkg/northbound/rest"
	subs "github.com/onosproject/onos-a1t/pkg/southbound/subscription"

	a1eistore "github.com/onosproject/onos-a1t/pkg/store/a1ei"
	a1pstore "github.com/onosproject/onos-a1t/pkg/store/a1p"
	substore "github.com/onosproject/onos-a1t/pkg/store/subscription"

	"github.com/onosproject/onos-lib-go/pkg/logging"
	"github.com/onosproject/onos-lib-go/pkg/northbound"
)

var log = logging.GetLogger("manager")

type Config struct {
	CAPath     string
	KeyPath    string
	CertPath   string
	GRPCPort   int
	ConfigPath string
	BaseURL    string
}

type Manager struct {
	restserver        *nbirest.Server
	submanager        *subs.SubscriptionManager
	broker            controller.Broker
	config            Config
	subscriptionStore substore.Store
	policyStore       a1pstore.Store
	eijobsStore       a1eistore.Store
}

func NewManager(config Config) (*Manager, error) {

	subscriptionStore := substore.NewStore()
	policyStore := a1pstore.NewStore()
	eijobsStore := a1eistore.NewStore()

	broker := controller.NewBroker(subscriptionStore, policyStore, eijobsStore)

	subManager, err := subs.NewSubscriptionManager(broker, subscriptionStore)
	if err != nil {
		return nil, err
	}

	restServer, err := nbirest.NewRestServer(config.BaseURL, broker)
	if err != nil {
		return nil, err
	}

	return &Manager{
		restserver:        restServer,
		submanager:        subManager,
		broker:            broker,
		subscriptionStore: subscriptionStore,
		policyStore:       policyStore,
		eijobsStore:       eijobsStore,
		config:            config,
	}, nil
}

func (m *Manager) startNorthboundServer() error {
	s := northbound.NewServer(northbound.NewServerCfg(
		m.config.CAPath,
		m.config.KeyPath,
		m.config.CertPath,
		int16(m.config.GRPCPort),
		true,
		northbound.SecurityConfig{}))

	s.AddService(nbi.NewService(m.subscriptionStore, m.policyStore, m.eijobsStore))

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

func (m *Manager) start() error {
	err := m.startNorthboundServer()
	if err != nil {
		log.Warn(err)
		return err
	}

	m.restserver.Start()

	err = m.submanager.Start()
	if err != nil {
		log.Warn(err)
		return err
	}

	return nil
}

func (m *Manager) Run() {
	err := m.start()
	if err != nil {
		log.Errorf("Error when starting KPIMON: %v", err)
	}
}
