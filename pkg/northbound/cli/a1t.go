// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0

package cli

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/onosproject/onos-a1t/pkg/controller"
	a1p "github.com/onosproject/onos-a1t/pkg/northbound/a1ap/policy_management"
	"github.com/onosproject/onos-a1t/pkg/rnib"
	"github.com/onosproject/onos-a1t/pkg/store"
	"github.com/onosproject/onos-lib-go/pkg/logging"
	"google.golang.org/grpc"
	"time"

	a1tadminapi "github.com/onosproject/onos-api/go/onos/a1t/admin"
	"github.com/onosproject/onos-lib-go/pkg/logging/service"
)

const TimeoutTimer = time.Second * 5

var cliLog = logging.GetLogger("northbound", "cli")

// NewService returns a new A1T interface service.
func NewService(subscriptionStore store.Store, policiesStore store.Store, eijobsStore store.Store, controllerBroker controller.Broker, rnibClient rnib.TopoClient) service.Service {
	return &Service{
		subscriptionStore: subscriptionStore,
		policiesStore:     policiesStore,
		eijobsStore:       eijobsStore,
		rnibClient:        rnibClient,
		ctrlBroker:        controllerBroker,
	}
}

// Service is a service implementation for administration.
type Service struct {
	service.Service
	subscriptionStore store.Store
	policiesStore     store.Store
	eijobsStore       store.Store
	rnibClient        rnib.TopoClient
	ctrlBroker        controller.Broker
}

func (s Service) Register(r *grpc.Server) {
	server := &Server{
		subscriptionStore: s.subscriptionStore,
		policiesStore:     s.policiesStore,
		eijobsStore:       s.eijobsStore,
		rnibClient:        s.rnibClient,
		ctrlBroker:        s.ctrlBroker,
	}
	a1tadminapi.RegisterA1TAdminServiceServer(r, server)
}

type Server struct {
	subscriptionStore store.Store
	policiesStore     store.Store
	eijobsStore       store.Store
	rnibClient        rnib.TopoClient
	ctrlBroker        controller.Broker
}

func (s *Server) GetXAppConnections(request *a1tadminapi.GetXAppConnectionsRequest, server a1tadminapi.A1TAdminService_GetXAppConnectionsServer) error {
	cliLog.Info("Get xApp Connection")
	ch := make(chan *store.Entry)
	ctx, cancel := context.WithTimeout(context.Background(), TimeoutTimer)
	defer cancel()
	go s.subscriptionStore.Entries(ctx, ch)

	for e := range ch {
		sKey := e.Key.(store.SubscriptionKey)
		sValue := e.Value.(*store.SubscriptionValue)
		if request.XappId != "" && request.XappId != string(sKey.TargetXAppID) {
			continue
		}

		endPoint := fmt.Sprintf("%s:%d", sValue.A1EndpointIP, sValue.A1EndpointPort)
		for _, c := range sValue.A1ServiceCapabilities {
			resp := &a1tadminapi.GetXAppConnectionResponse{
				XappId:                   string(sKey.TargetXAppID),
				SupportedA1Service:       c.A1Service.String(),
				SupportedA1ServiceTypeId: c.TypeID,
				XappA1Endpoint:           endPoint,
			}
			err := server.Send(resp)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (s *Server) GetPolicyTypeObject(request *a1tadminapi.GetPolicyTypeObjectRequest, server a1tadminapi.A1TAdminService_GetPolicyTypeObjectServer) error {
	cliLog.Info("Get policy type object")
	ctx, cancel := context.WithTimeout(context.Background(), TimeoutTimer)
	defer cancel()

	policyTypes := s.ctrlBroker.A1PController().HandleGetPolicyTypes(ctx)
	for _, pt := range policyTypes {
		policyTypeSchema, statusSchema, err := s.ctrlBroker.A1PController().HandleGetPolicytypesPolicyTypeId(ctx, pt)
		if err != nil {
			return err
		}

		if request.PolicyTypeId != "" && request.PolicyTypeId != pt {
			continue
		}

		policyTypeObject := a1p.PolicyTypeObject{
			PolicySchema: policyTypeSchema,
			StatusSchema: (*a1p.JsonSchema)(&statusSchema),
		}
		pto, err := json.Marshal(policyTypeObject)
		if err != nil {
			return err
		}

		oIDs, err := s.ctrlBroker.A1PController().HandleGetPolicytypesPolicyTypeIdPolicies(ctx, pt)
		if err != nil {
			return err
		}

		resp := &a1tadminapi.GetPolicyTypeObjectResponse{
			PolicyTypeId:     pt,
			PolicyIds:        oIDs,
			PolicyTypeObject: string(pto),
		}

		err = server.Send(resp)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *Server) GetPolicyObject(request *a1tadminapi.GetPolicyObjectRequest, server a1tadminapi.A1TAdminService_GetPolicyObjectServer) error {
	cliLog.Info("Get policy object")
	ctx, cancel := context.WithTimeout(context.Background(), TimeoutTimer)
	defer cancel()

	var err error
	tIDs := make([]string, 0)
	pIDs := make([]string, 0)

	if request.PolicyTypeId == "" {
		tIDs = s.ctrlBroker.A1PController().HandleGetPolicyTypes(ctx)
	} else {
		tIDs = append(tIDs, request.PolicyTypeId)
	}

	for _, t := range tIDs {
		if request.PolicyObjectId == "" {
			pIDs, err = s.ctrlBroker.A1PController().HandleGetPolicytypesPolicyTypeIdPolicies(ctx, t)
			if err != nil {
				return err
			}
		} else {
			pIDs = append(pIDs, request.PolicyObjectId)
		}

		for _, i := range pIDs {
			obj, err := s.ctrlBroker.A1PController().HandleGetPolicy(ctx, i, t)
			if err != nil {
				return err
			}
			objJson, err := json.Marshal(obj)
			if err != nil {
				return err
			}
			resp := &a1tadminapi.GetPolicyObjectResponse{
				PolicyTypeId:   t,
				PolicyObjectId: i,
				PolicyObject:   string(objJson),
			}
			err = server.Send(resp)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (s *Server) GetPolicyObjectStatus(request *a1tadminapi.GetPolicyObjectStatusRequest, server a1tadminapi.A1TAdminService_GetPolicyObjectStatusServer) error {
	cliLog.Info("Get policy type object status")
	ctx, cancel := context.WithTimeout(context.Background(), TimeoutTimer)
	defer cancel()

	var err error
	tIDs := make([]string, 0)
	pIDs := make([]string, 0)

	if request.PolicyTypeId == "" {
		tIDs = s.ctrlBroker.A1PController().HandleGetPolicyTypes(ctx)
	} else {
		tIDs = append(tIDs, request.PolicyTypeId)
	}

	for _, t := range tIDs {
		if request.PolicyObjectId == "" {
			pIDs, err = s.ctrlBroker.A1PController().HandleGetPolicytypesPolicyTypeIdPolicies(ctx, t)
			if err != nil {
				return err
			}
		} else {
			pIDs = append(pIDs, request.PolicyObjectId)
		}

		for _, i := range pIDs {
			obj, err := s.ctrlBroker.A1PController().HandleGetPolicyStatus(ctx, i, t)
			if err != nil {
				return err
			}
			objJson, err := json.Marshal(obj)
			if err != nil {
				return err
			}
			resp := &a1tadminapi.GetPolicyObjectStatusResponse{
				PolicyTypeId:       t,
				PolicyObjectId:     i,
				PolicyObjectStatus: string(objJson),
			}
			err = server.Send(resp)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
