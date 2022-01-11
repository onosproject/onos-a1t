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

//ToDo - eiJobObject should be of type EIJobObject as in onos-a1t/pkg/northbound/a1ap/enrichment_information/a1ap_ei.go
func CreateEIjob(ctx context.Context, address, certPath, keyPath string, eiJobID, eiJobTypeID, eiJobObject string) error {
	conn, err := a1tsb.GetConnection(ctx, address, certPath, keyPath)
	if err != nil {
		return err
	}
	defer conn.Close()

	request := a1tapi.CreateRequest{
		Object: &a1tapi.Object{
			//ToDo - add in ID and Revision in the future
			Type: a1tapi.Object_EIJOB,
			Obj: &a1tapi.Object_Eijob{
				Eijob: &a1tapi.EIJob{
					Id:     eiJobID,
					Typeid: eiJobTypeID,
					Object: []byte(eiJobObject),
					//ToDo - add in status in the future
					//Status: &a1tapi.Status{}
				},
			},
		},
	}
	client := a1tapi.NewA1TClient(conn)

	respCreate, err := client.Create(context.Background(), &request)
	if err != nil {
		return err
	}

	if respCreate.GetObject().String() == "" {
		return fmt.Errorf("EI Job object create failed")
	}

	return nil
}

//ToDo - eiJobObject should be of type EIJobObject as in onos-a1t/pkg/northbound/a1ap/enrichment_information/a1ap_ei.go
func DeleteEIjob(ctx context.Context, address, certPath, keyPath string, eiJobID, eiJobTypeID, eiJobObject string) error {
	conn, err := a1tsb.GetConnection(ctx, address, certPath, keyPath)
	if err != nil {
		return err
	}
	defer conn.Close()

	request := a1tapi.DeleteRequest{
		Object: &a1tapi.Object{
			//ToDo - add in ID and Revision in the future
			Type: a1tapi.Object_EIJOB,
			Obj: &a1tapi.Object_Eijob{
				Eijob: &a1tapi.EIJob{
					Id:     eiJobID,
					Typeid: eiJobTypeID,
					Object: []byte(eiJobObject),
					//ToDo - add in status in the future
					//Status: &a1tapi.Status{}
				},
			},
		},
	}
	client := a1tapi.NewA1TClient(conn)

	respCreate, err := client.Delete(context.Background(), &request)
	if err != nil {
		return err
	}

	if respCreate.GetObject().String() == "" {
		return fmt.Errorf("EI Job object delete failed")
	}

	return nil
}
