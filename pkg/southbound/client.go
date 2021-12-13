// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-Only-1.0

package southbound

import (
	"context"

	"google.golang.org/grpc"

	sbi "github.com/onosproject/onos-lib-go/pkg/southbound"
)

func GetConnection(ctx context.Context, address, certPath, keyPath string) (*grpc.ClientConn, error) {
	clientConn, err := sbi.Connect(ctx, address, certPath, keyPath)

	if err != nil {
		return nil, err
	}

	return clientConn, err
}
