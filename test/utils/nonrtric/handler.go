package nonrtric

import (
	"fmt"
	"net/http"

	"github.com/deepmap/oapi-codegen/pkg/runtime"
	"github.com/labstack/echo/v4"

	a1ei "github.com/onosproject/onos-a1t/pkg/northbound/a1ap/enrichment_information"
	a1p "github.com/onosproject/onos-a1t/pkg/northbound/a1ap/policy_management"
)

type a1eiWraper struct {
	version    string
	controller Controller
}

func SetRESTA1EIWraper(e *echo.Echo, version string, controller Controller) {
	wraper := &a1eiWraper{
		version:    version,
		controller: controller,
	}

	a1ei.RegisterHandlers(e, wraper)
}

// EI job identifiers
// (GET /A1-EI/v1/eijobs)
func (a1eiw *a1eiWraper) GetEiJobIdsUsingGET(ctx echo.Context, params a1ei.GetEiJobIdsUsingGETParams) error {
	return nil
}

// Individual EI job
// (DELETE /A1-EI/v1/eijobs/{eiJobId})
func (a1eiw *a1eiWraper) DeleteIndividualEiJobUsingDELETE(ctx echo.Context, eiJobId string) error {
	return nil
}

// Individual EI job
// (GET /A1-EI/v1/eijobs/{eiJobId})
func (a1eiw *a1eiWraper) GetIndividualEiJobUsingGET(ctx echo.Context, eiJobId string) error {
	return nil
}

// Individual EI job
// (PUT /A1-EI/v1/eijobs/{eiJobId})
func (a1eiw *a1eiWraper) PutIndividualEiJobUsingPUT(ctx echo.Context, eiJobId string) error {
	return nil
}

// EI job status
// (GET /A1-EI/v1/eijobs/{eiJobId}/status)
func (a1eiw *a1eiWraper) GetEiJobStatusUsingGET(ctx echo.Context, eiJobId string) error { return nil }

// EI type identifiers
// (GET /A1-EI/v1/eitypes)
func (a1eiw *a1eiWraper) GetEiTypeIdentifiersUsingGET(ctx echo.Context) error { return nil }

// Individual EI type
// (GET /A1-EI/v1/eitypes/{eiTypeId})
func (a1eiw *a1eiWraper) GetEiTypeUsingGET(ctx echo.Context, eiTypeId string) error { return nil }

type a1pWraper struct {
	version    string
	controller Controller
}

func SetRESTA1PMWraper(e *echo.Echo, version string, controller Controller) {
	wraper := &a1pWraper{
		version:    version,
		controller: controller,
	}

	RegisterPolicyHandlers(version, e, wraper)
}

type PMServerInterface interface {
	// Individual Policy
	// (POST /A1-P/v1/policies/{policyId})
	PostIndividualPolicyUsingPOST(ctx echo.Context, policyId string) error
}

// ServerInterfaceWrapper converts echo contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler PMServerInterface
}

// PostIndividualPolicyUsingPOST converts echo context to params.
func (w *ServerInterfaceWrapper) PostIndividualPolicyUsingPOST(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "policyId" -------------
	var policyId string

	err = runtime.BindStyledParameterWithLocation("simple", false, "policyId", runtime.ParamLocationPath, ctx.Param("policyId"), &policyId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter policyId: %s", err))
	}

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.PostIndividualPolicyUsingPOST(ctx, policyId)
	return err
}

// RegisterHandlers adds each server route to the EchoRouter.
func RegisterPolicyHandlers(version string, router a1p.EchoRouter, si PMServerInterface) {
	registerHandlersWithBaseURL(version, router, si, "")
}

// Registers handlers, and prepends BaseURL to the paths, so that the paths
// can be served under a prefix.
func registerHandlersWithBaseURL(version string, router a1p.EchoRouter, si PMServerInterface, baseURL string) {

	wrapper := ServerInterfaceWrapper{
		Handler: si,
	}

	router.POST(baseURL+"/A1-P/"+version+"/policies/:policyId/notify", wrapper.PostIndividualPolicyUsingPOST)
}

func (a1pw *a1pWraper) PostIndividualPolicyUsingPOST(ctx echo.Context, policyId string) error {

	eiJobObjNot := make(map[string]interface{})

	if err := ctx.Bind(&eiJobObjNot); err != nil {
		return ctx.JSON(http.StatusOK, err)
	}

	/* 	err := a1pw.a1pController.HandlePolicyNotify(ctx.Request().Context(), policyId, eiJobObjNot)
	   	if err != nil {
	   		return err
	   	}
	*/
	return nil
}
