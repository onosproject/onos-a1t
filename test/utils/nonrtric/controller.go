// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-Only-1.0

package nonrtric

import (
	"context"

	"github.com/labstack/echo/v4"
	a1pm "github.com/onosproject/onos-a1t/pkg/northbound/a1ap/policy_management"
	"github.com/onosproject/onos-a1t/pkg/store"
)

type Controller interface {
	/*
		A1P Client Requests
	*/

	// GetPolicytypes request
	A1PMGetPolicytypes(ctx context.Context) ([]string, error)

	// GetPolicytypesPolicyTypeId request
	A1PMGetPolicytypesPolicyTypeId(ctx context.Context, policyTypeId string) (string, error)

	// GetPolicytypesPolicyTypeIdPolicies request
	A1PMGetPolicytypesPolicyTypeIdPolicies(ctx context.Context, policyTypeId string) ([]string, error)

	// DeletePolicytypesPolicyTypeIdPoliciesPolicyId request
	A1PMDeletePolicytypesPolicyTypeIdPoliciesPolicyId(ctx context.Context, policyTypeId, policyId string) error

	// GetPolicytypesPolicyTypeIdPoliciesPolicyId request
	A1PMGetPolicytypesPolicyTypeIdPoliciesPolicyId(ctx context.Context, policyTypeId, policyId string) (string, error)

	// PutPolicytypesPolicyTypeIdPoliciesPolicyId request with any body
	A1PMPutPolicytypesPolicyTypeIdPoliciesPolicyIdWithBody(ctx context.Context, policyTypeId, policyId string, body string) error

	A1PMPutPolicytypesPolicyTypeIdPoliciesPolicyId(ctx context.Context, policyTypeId, policyId, param, body string) error

	// GetPolicytypesPolicyTypeIdPoliciesPolicyIdStatus request
	A1PMGetPolicytypesPolicyTypeIdPoliciesPolicyIdStatus(ctx context.Context, policyTypeId, policyId string) (string, error)

	/*
		A1EI Server Handlers
	*/

	// EI job identifiers
	// (GET /A1-EI/v1/eijobs)
	A1EIGetEiJobIdsUsingGET(ctx echo.Context, eiTypeId string) error

	// Individual EI job
	// (DELETE /A1-EI/v1/eijobs/{eiJobId})
	A1EIDeleteIndividualEiJobUsingDELETE(ctx echo.Context, eiJobId string) error

	// Individual EI job
	// (GET /A1-EI/v1/eijobs/{eiJobId})
	A1EIGetIndividualEiJobUsingGET(ctx echo.Context, eiJobId string) error

	// Individual EI job
	// (PUT /A1-EI/v1/eijobs/{eiJobId})
	A1EIPutIndividualEiJobUsingPUT(ctx echo.Context, eiJobId string) error

	// EI job status
	// (GET /A1-EI/v1/eijobs/{eiJobId}/status)
	A1EIGetEiJobStatusUsingGET(ctx echo.Context, eiJobId string) error

	// EI type identifiers
	// (GET /A1-EI/v1/eitypes)
	A1EIGetEiTypeIdentifiersUsingGET(ctx echo.Context) error

	// Individual EI type
	// (GET /A1-EI/v1/eitypes/{eiTypeId})
	A1EIGetEiTypeUsingGET(ctx echo.Context, eiTypeId string) error
}

type controller struct {
	nearRTRicBaseURL string
	policyStore      store.Store
	eijobsStore      store.Store
	a1pClient        a1pm.ClientWithResponsesInterface
}

func NewController(nearRTRicBaseURL string, policyStore, eijobsStore store.Store) Controller {
	a1pClient, err := a1pm.NewClientWithResponses(nearRTRicBaseURL)
	if err != nil {
		log.Fatal(err)
	}

	return &controller{
		nearRTRicBaseURL: nearRTRicBaseURL,
		policyStore:      policyStore,
		eijobsStore:      eijobsStore,
		a1pClient:        a1pClient,
	}
}
