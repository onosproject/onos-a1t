// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-Only-1.0

package manager

import (
	"context"
	"github.com/onosproject/onos-a1t/pkg/southbound"
	"github.com/onosproject/onos-a1t/pkg/store"
	"github.com/onosproject/onos-a1t/pkg/stream"
	"strconv"
	"strings"

	"github.com/onosproject/onos-a1t/pkg/controller"
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
	restServer        *nbirest.Server
	subManager        *subs.Manager
	sbManager         southbound.Manager
	broker            controller.Broker
	streamBroker      stream.Broker
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

	streamBroker := stream.NewBroker()

	subManager, err := subs.NewSubscriptionManager(subscriptionStore)
	if err != nil {
		return nil, err
	}

	sbManager := southbound.NewSouthboundManager(streamBroker, subscriptionStore)

	restServer, err := nbirest.NewRestServer(config.BaseURL, broker)
	if err != nil {
		return nil, err
	}

	rnibClient, err := rnib.NewClient()
	if err != nil {
		return nil, err
	}

	return &Manager{
		restServer:        restServer,
		subManager:        subManager,
		sbManager:         sbManager,
		broker:            broker,
		streamBroker:      streamBroker,
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

	//s.AddService(nbi.NewService(m.subscriptionStore, m.policyStore, m.eijobsStore))

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

	err = m.subManager.Start()
	if err != nil {
		log.Warn(err)
		return err
	}

	err = m.sbManager.Run(context.Background())
	if err != nil {
		log.Warn(err)
		return err
	}

	m.restServer.Start()

	return nil
}

func (m *Manager) Run() {
	err := m.start()
	if err != nil {
		log.Errorf("Error when starting A1T: %v", err)
	}
}
