// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-Only-1.0

package xapp

import (
	"context"

	"github.com/onosproject/onos-api/go/onos/topo"
	a1connection "github.com/onosproject/onos-ric-sdk-go/pkg/a1/connection"
)

func NewNBIManager(caPath string, keyPath string, certPath string, grpcPort int, xAppName string, a1PolicyTypes []*topo.A1PolicyType) (*NBIManager, error) {
	a1ConnManager, err := a1connection.NewManager(caPath, keyPath, certPath, grpcPort, a1PolicyTypes)
	if err != nil {
		return nil, err
	}
	return &NBIManager{
		a1ConnManager: a1ConnManager,
	}, nil
}

type NBIManager struct {
	a1ConnManager *a1connection.Manager
}

func (m *NBIManager) Start() {
	m.a1ConnManager.Start(context.Background())
}

func (m *NBIManager) Close(ctx context.Context) {
	err := m.a1ConnManager.DeleteXAppElementOnTopo(ctx)
	if err != nil {
		log.Error(err)
	}
}
