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
	policyTypes := a1pw.a1pController.HandleGetPolicyTypes(ctx.Request().Context())
	return ctx.JSON(http.StatusOK, policyTypes)
}

// (GET /policytypes/{policyTypeId})
func (a1pw *a1pWraper) GetPolicytypesPolicyTypeId(ctx echo.Context, policyTypeId a1p.PolicyTypeId) error {
	return nil
}

// (GET /policytypes/{policyTypeId}/policies)
func (a1pw *a1pWraper) GetPolicytypesPolicyTypeIdPolicies(ctx echo.Context, policyTypeId a1p.PolicyTypeId) error {
	a1pEntriesValues, err := a1pw.a1pController.HandleGetPoliciesTypeID(ctx.Request().Context(), string(policyTypeId))
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, err)
	}

	// for a1pValue := range a1pEntriesValues {
	// 	a1pValue.PolicyObject
	// }
	return ctx.JSON(http.StatusOK, a1pEntriesValues)
}

// (DELETE /policytypes/{policyTypeId}/policies/{policyId})
func (a1pw *a1pWraper) DeletePolicytypesPolicyTypeIdPoliciesPolicyId(ctx echo.Context, policyTypeId a1p.PolicyTypeId, policyId a1p.PolicyId) error {
	err := a1pw.a1pController.HandlePolicyDelete(ctx.Request().Context(), string(policyTypeId), string(policyId))
	return ctx.JSON(http.StatusOK, err)
}

// (GET /policytypes/{policyTypeId}/policies/{policyId})
func (a1pw *a1pWraper) GetPolicytypesPolicyTypeIdPoliciesPolicyId(ctx echo.Context, policyTypeId a1p.PolicyTypeId, policyId a1p.PolicyId) error {
	a1pEntryValue, err := a1pw.a1pController.HandleGetPolicy(ctx.Request().Context(), string(policyId), string(policyTypeId))
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, err)
	}

	return ctx.JSON(http.StatusOK, a1pEntryValue)
}

// (PUT /policytypes/{policyTypeId}/policies/{policyId})
func (a1pw *a1pWraper) PutPolicytypesPolicyTypeIdPoliciesPolicyId(ctx echo.Context, policyTypeId a1p.PolicyTypeId, policyId a1p.PolicyId, params a1p.PutPolicytypesPolicyTypeIdPoliciesPolicyIdParams) error {
	policyObject := make(map[string]string)
	paramsMap := make(map[string]string)
	paramsMap["notificationDestination"] = string(*params.NotificationDestination)

	if err := ctx.Bind(&policyObject); err != nil {
		return ctx.JSON(http.StatusOK, err)
	}

	err := a1pw.a1pController.HandlePolicyCreate(ctx.Request().Context(), string(policyTypeId), string(policyId), paramsMap, policyObject)
	return ctx.JSON(http.StatusOK, err)
}

// (GET /policytypes/{policyTypeId}/policies/{policyId}/status)
func (a1pw *a1pWraper) GetPolicytypesPolicyTypeIdPoliciesPolicyIdStatus(ctx echo.Context, policyTypeId a1p.PolicyTypeId, policyId a1p.PolicyId) error {
	a1pPolicyStatus, err := a1pw.a1pController.HandleGetPolicyStatus(ctx.Request().Context(), string(policyId), string(policyTypeId))
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, err)
	}

	return ctx.JSON(http.StatusOK, a1pPolicyStatus)
}
