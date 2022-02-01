// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-Only-1.0

package handler

import (
	"encoding/json"
	"github.com/onosproject/onos-a1t/pkg/utils"
	"github.com/onosproject/onos-lib-go/pkg/errors"
	"github.com/onosproject/onos-lib-go/pkg/logging"
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/onosproject/onos-a1t/pkg/controller"
	a1p "github.com/onosproject/onos-a1t/pkg/northbound/a1ap/policy_management"
)

var logA1P = logging.GetLogger("handler", "a1p")

type a1pWraper struct {
	version       string
	a1pController controller.A1PController
}

func SetRESTA1PWraper(e *echo.Echo, version string, a1pController controller.A1PController) {
	wraper := &a1pWraper{
		version:       version,
		a1pController: a1pController,
	}
	a1p.RegisterHandlers(e, wraper)
}

// (GET /policytypes)
func (a1pw *a1pWraper) GetPolicytypes(ctx echo.Context) error {
	policyTypes := a1pw.a1pController.HandleGetPolicyTypes(ctx.Request().Context())
	return ctx.JSONPretty(http.StatusOK, policyTypes, "  ")
}

// (GET /policytypes/{policyTypeId})
func (a1pw *a1pWraper) GetPolicytypesPolicyTypeId(ctx echo.Context, policyTypeId a1p.PolicyTypeId) error {
	policyType, err := a1pw.a1pController.HandleGetPolicytypesPolicyTypeId(ctx.Request().Context(), string(policyTypeId))
	if err != nil {
		return ctx.JSONPretty(http.StatusBadRequest, err, "  ")
	}
	return ctx.JSONPretty(http.StatusOK, policyType, "  ")
}

// (GET /policytypes/{policyTypeId}/policies)
func (a1pw *a1pWraper) GetPolicytypesPolicyTypeIdPolicies(ctx echo.Context, policyTypeId a1p.PolicyTypeId) error {
	a1pEntriesValues, err := a1pw.a1pController.HandleGetPolicytypesPolicyTypeIdPolicies(ctx.Request().Context(), string(policyTypeId))
	if err != nil {
		return ctx.JSONPretty(http.StatusBadRequest, err, "  ")
	}

	return ctx.JSONPretty(http.StatusOK, a1pEntriesValues, "  ")
}

// (DELETE /policytypes/{policyTypeId}/policies/{policyId})
func (a1pw *a1pWraper) DeletePolicytypesPolicyTypeIdPoliciesPolicyId(ctx echo.Context, policyTypeId a1p.PolicyTypeId, policyId a1p.PolicyId) error {
	err := a1pw.a1pController.HandlePolicyDelete(ctx.Request().Context(), string(policyId), string(policyTypeId))

	if err != nil {
		return ctx.JSONPretty(http.StatusServiceUnavailable, err, "  ")
	}

	return ctx.JSONPretty(http.StatusNoContent, err, "  ")
}

// (GET /policytypes/{policyTypeId}/policies/{policyId})
func (a1pw *a1pWraper) GetPolicytypesPolicyTypeIdPoliciesPolicyId(ctx echo.Context, policyTypeId a1p.PolicyTypeId, policyId a1p.PolicyId) error {
	// query - done
	a1pEntryValue, err := a1pw.a1pController.HandleGetPolicy(ctx.Request().Context(), string(policyId), string(policyTypeId))
	if err != nil {
		return ctx.JSONPretty(http.StatusServiceUnavailable, err, "  ")
	}

	return ctx.JSONPretty(http.StatusOK, a1pEntryValue, "  ")
}

// (PUT /policytypes/{policyTypeId}/policies/{policyId})
func (a1pw *a1pWraper) PutPolicytypesPolicyTypeIdPoliciesPolicyId(ctx echo.Context, policyTypeId a1p.PolicyTypeId, policyId a1p.PolicyId, params a1p.PutPolicytypesPolicyTypeIdPoliciesPolicyIdParams) error {
	// create (done) or update (todo)
	policyObject := make(map[string]interface{})
	paramsMap := make(map[string]string)

	if params.NotificationDestination == nil {
		return ctx.JSONPretty(http.StatusServiceUnavailable, "notificationDestination is missing in HTTPRequest", "  ")
	}
	paramsMap["notificationDestination"] = string(*params.NotificationDestination)

	if err := ctx.Bind(&policyObject); err != nil {
		return ctx.JSONPretty(http.StatusOK, err, "  ")
	}

	policyObject = utils.GetPolicyObject(policyObject)

	obj, err := json.Marshal(policyObject)
	if err != nil {
		return ctx.JSONPretty(http.StatusServiceUnavailable, err, "  ")
	}

	if !utils.JsonValidateWithTypeID(string(policyTypeId), string(obj)) {
		return ctx.JSONPretty(http.StatusServiceUnavailable, errors.NewInvalid("PolicyObject validation failed: policyObject %v", policyObject), "  ")
	}

	a1pEntriesValues, err := a1pw.a1pController.HandleGetPolicytypesPolicyTypeIdPolicies(ctx.Request().Context(), string(policyTypeId))
	if err != nil {
		return ctx.JSONPretty(http.StatusServiceUnavailable, err, "  ")
	}

	hasPolicyID := false

	for _, v := range a1pEntriesValues {
		if v == string(policyId) {
			hasPolicyID = true
		}
	}

	if hasPolicyID {
		err = a1pw.a1pController.HandlePolicyUpdate(ctx.Request().Context(), string(policyId), string(policyTypeId), paramsMap, policyObject)
		if err != nil {
			return ctx.JSONPretty(http.StatusServiceUnavailable, err, "  ")
		}
		return ctx.JSONPretty(http.StatusOK, policyObject, "  ")
	}

	err = a1pw.a1pController.HandlePolicyCreate(ctx.Request().Context(), string(policyId), string(policyTypeId), paramsMap, policyObject)
	if err != nil {
		return ctx.JSONPretty(http.StatusServiceUnavailable, err, "  ")
	}
	return ctx.JSONPretty(http.StatusCreated, policyObject, "  ")
}

// (GET /policytypes/{policyTypeId}/policies/{policyId}/status)
func (a1pw *a1pWraper) GetPolicytypesPolicyTypeIdPoliciesPolicyIdStatus(ctx echo.Context, policyTypeId a1p.PolicyTypeId, policyId a1p.PolicyId) error {
	a1pPolicyStatus, err := a1pw.a1pController.HandleGetPolicyStatus(ctx.Request().Context(), string(policyId), string(policyTypeId))
	if err != nil {
		return ctx.JSONPretty(http.StatusBadRequest, err, "  ")
	}

	return ctx.JSONPretty(http.StatusOK, a1pPolicyStatus, "  J")
}
