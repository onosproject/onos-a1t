// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-Only-1.0

package southbound

import (
	"context"
	"fmt"

	prototypes "github.com/gogo/protobuf/types"

	a1tsb "github.com/onosproject/onos-a1t/pkg/southbound"
	a1tapi "github.com/onosproject/onos-a1t/pkg/southbound/a1t"
)

func CreateEIjob(ctx context.Context, address, certPath, keyPath string, eiJobID, eiJobTypeID string, eiJobObject map[string]string) error {
	conn, err := a1tsb.GetConnection(ctx, address, certPath, keyPath)
	if err != nil {
		return err
	}
	defer conn.Close()

	var eiJobObjectValue *prototypes.Any
	objValue := &a1tapi.ObjectValue{Value: eiJobObject}
	eiJobObjectValue, err = prototypes.MarshalAny(objValue)
	if err != nil {
		return err
	}

	request := a1tapi.CreateRequest{
		Object: &a1tapi.Object{
			Type: a1tapi.Object_EIJOB,
			Obj: &a1tapi.Object_Policy{
				Policy: &a1tapi.Policy{
					Id:     eiJobID,
					Typeid: eiJobTypeID,
					Object: eiJobObjectValue,
				},
			},
		},
	}
	client := a1tapi.NewA1TClient(conn)

	respCreate, err := client.Create(context.Background(), &request)
	if err != nil {
		return err
	}

	if respCreate.GetObject().Id != "" {
		return fmt.Errorf("EI Job object create failed")
	}

	return nil
}

func DeleteEIjob(ctx context.Context, address, certPath, keyPath string, eiJobID, eiJobTypeID string) error {
	conn, err := a1tsb.GetConnection(ctx, address, certPath, keyPath)
	if err != nil {
		return err
	}
	defer conn.Close()

	request := a1tapi.DeleteRequest{
		Object: &a1tapi.Object{
			Type: a1tapi.Object_EIJOB,
			Obj: &a1tapi.Object_Policy{
				Policy: &a1tapi.Policy{
					Id:     eiJobID,
					Typeid: eiJobTypeID,
				},
			},
		},
	}
	client := a1tapi.NewA1TClient(conn)

	respCreate, err := client.Delete(context.Background(), &request)
	if err != nil {
		return err
	}

	if respCreate.GetObject().Id != "" {
		return fmt.Errorf("EI Job object delete failed")
	}

	return nil
}
