// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-Only-1.0

package southbound

import (
	"context"
	"encoding/json"
	"fmt"

	a1tsb "github.com/onosproject/onos-a1t/pkg/southbound"
	a1tapi "github.com/onosproject/onos-api/go/onos/a1t/a1"
)

func CreatePolicy(ctx context.Context, address, certPath, keyPath string, policyID, policyTypeID string, policyObject map[string]interface{}) error {
	conn, err := a1tsb.GetConnection(ctx, address, certPath, keyPath)
	if err != nil {
		return err
	}
	defer conn.Close()

	policyObjectValue, err := json.Marshal(policyObject)
	if err != nil {
		return err
	}

	request := a1tapi.PolicyRequestMessage{
		PolicyType: &a1tapi.PolicyType{
			Id: policyTypeID,
		},
		Message: &a1tapi.RequestMessage{
			Payload: policyObjectValue,
		},
	}
	client := a1tapi.NewPolicyServiceClient(conn)

	respCreate, err := client.PolicySetup(context.Background(), &request)
	if err != nil {
		return err
	}

	if respCreate.String() == "" {
		return fmt.Errorf("policy object create failed")
	}

	return nil
}

func DeletePolicy(ctx context.Context, address, certPath, keyPath string, policyID, policyTypeID string) error {
	conn, err := a1tsb.GetConnection(ctx, address, certPath, keyPath)
	if err != nil {
		return err
	}
	defer conn.Close()

	request := a1tapi.PolicyRequestMessage{
		PolicyType: &a1tapi.PolicyType{
			Id: policyTypeID,
		},
	}

	client := a1tapi.NewPolicyServiceClient(conn)

	respCreate, err := client.PolicyDelete(context.Background(), &request)
	if err != nil {
		return err
	}

	if respCreate.String() == "" {
		return fmt.Errorf("policy object delete failed")
	}

	return nil
}
