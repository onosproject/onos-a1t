// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-Only-1.0

package southbound

import (
	"context"
	"fmt"

	a1tsb "github.com/onosproject/onos-a1t/pkg/southbound"
	a1tapi "github.com/onosproject/onos-a1t/pkg/southbound/a1t"
)

func CreatePolicy(ctx context.Context, clientcfg a1tsb.Client) error {
	conn, err := a1tsb.GetConnection(ctx, clientcfg)
	if err != nil {
		return err
	}
	defer conn.Close()

	request := a1tapi.CreateRequest{}
	client := a1tapi.NewA1TClient(conn)

	respCreate, err := client.Create(context.Background(), &request)
	if err != nil {
		return err
	}

	if respCreate.GetObject().Id != "" {
		return fmt.Errorf("policy object create failed")
	}

	return nil
}
