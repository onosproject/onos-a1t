// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-Only-1.0

package manager

import (
	"context"
	"github.com/onosproject/onos-a1t/pkg/store"
	"strconv"
	"strings"

	"github.com/onosproject/onos-a1t/pkg/controller"
	nbi "github.com/onosproject/onos-a1t/pkg/northbound/cli"
	nbirest "github.com/onosproject/onos-a1t/pkg/northbound/rest"
	"github.com/onosproject/onos-a1t/pkg/rnib"
	subs "github.com/onosproject/onos-a1t/pkg/subscription"

	"github.com/onosproject/onos-lib-go/pkg/logging"
	"github.com/onosproject/onos-lib-go/pkg/northbound"
)

var log = logging.GetLogger("manager")

type Config struct {
	CAPath      string
	KeyPath     string
	CertPath    string
	GRPCPort    int
	ConfigPath  string
	BaseURL     string
	NonRTRICURL string
}

type Manager struct {
	restserver        *nbirest.Server
	submanager        *subs.SubscriptionManager
	broker            controller.Broker
	config            Config
	subscriptionStore store.Store
	policyStore       store.Store
	eijobsStore       store.Store
	rnibClient        rnib.Client
}

func NewManager(config Config) (*Manager, error) {

	subscriptionStore := store.NewStore()
	policyStore := store.NewStore()
	eijobsStore := store.NewStore()

	broker := controller.NewBroker(config.NonRTRICURL, subscriptionStore, policyStore, eijobsStore)

	subManager, err := subs.NewSubscriptionManager(broker, subscriptionStore)
	if err != nil {
		return nil, err
	}

	restServer, err := nbirest.NewRestServer(config.BaseURL, broker)
	if err != nil {
		return nil, err
	}

	rnibClient, err := rnib.NewClient()
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
		rnibClient:        rnibClient,
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

func (m *Manager) registerA1TtoRnib() error {
	nbPort, err := strconv.Atoi(strings.Split(m.config.BaseURL, ":")[1])
	if err != nil {
		return err
	}
	return m.rnibClient.AddA1TEntity(context.Background(), uint32(nbPort))
}

func (m *Manager) start() error {
	err := m.registerA1TtoRnib()
	if err != nil {
		return err
	}

	err = m.startNorthboundServer()
	if err != nil {
		log.Warn(err)
		return err
	}

	err = m.submanager.Start()
	if err != nil {
		log.Warn(err)
		return err
	}

	m.restserver.Start()

	return nil
}

func (m *Manager) Run() {
	err := m.start()
	if err != nil {
		log.Errorf("Error when starting KPIMON: %v", err)
	}
}
