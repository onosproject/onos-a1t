// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-Only-1.0

package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/onosproject/onos-a1t/pkg/controller"
	a1p "github.com/onosproject/onos-a1t/pkg/northbound/a1ap/policy_management"
)

type a1pWraper struct {
	version       string
	a1pController controller.A1PController
}

func NewA1pWraper(version string, a1pController controller.A1PController) a1p.ServerInterface {
	return &a1pWraper{
		version:       version,
		a1pController: a1pController,
	}
}

// (GET /policytypes)
func (a1pw *a1pWraper) GetPolicytypes(ctx echo.Context) error {
	policyTypes := a1pw.a1pController.HandleGetPolicyTypes()
	return ctx.JSON(http.StatusOK, policyTypes)
}

// (GET /policytypes/{policyTypeId})
func (a1pw *a1pWraper) GetPolicytypesPolicyTypeId(ctx echo.Context, policyTypeId a1p.PolicyTypeId) error {
	return nil
}

// (GET /policytypes/{policyTypeId}/policies)
func (a1pw *a1pWraper) GetPolicytypesPolicyTypeIdPolicies(ctx echo.Context, policyTypeId a1p.PolicyTypeId) error {
	return nil
}

// (DELETE /policytypes/{policyTypeId}/policies/{policyId})
func (a1pw *a1pWraper) DeletePolicytypesPolicyTypeIdPoliciesPolicyId(ctx echo.Context, policyTypeId a1p.PolicyTypeId, policyId a1p.PolicyId) error {
	return nil
}

// (GET /policytypes/{policyTypeId}/policies/{policyId})
func (a1pw *a1pWraper) GetPolicytypesPolicyTypeIdPoliciesPolicyId(ctx echo.Context, policyTypeId a1p.PolicyTypeId, policyId a1p.PolicyId) error {
	return nil
}

// (PUT /policytypes/{policyTypeId}/policies/{policyId})
func (a1pw *a1pWraper) PutPolicytypesPolicyTypeIdPoliciesPolicyId(ctx echo.Context, policyTypeId a1p.PolicyTypeId, policyId a1p.PolicyId, params a1p.PutPolicytypesPolicyTypeIdPoliciesPolicyIdParams) error {
	return nil
}

// (GET /policytypes/{policyTypeId}/policies/{policyId}/status)
func (a1pw *a1pWraper) GetPolicytypesPolicyTypeIdPoliciesPolicyIdStatus(ctx echo.Context, policyTypeId a1p.PolicyTypeId, policyId a1p.PolicyId) error {
	return nil
}
