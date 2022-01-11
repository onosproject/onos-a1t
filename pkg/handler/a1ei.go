// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-Only-1.0

package handler

import (
	"fmt"
	"net/http"

	"github.com/deepmap/oapi-codegen/pkg/runtime"
	"github.com/labstack/echo/v4"
	"github.com/onosproject/onos-a1t/pkg/controller"
	a1ei "github.com/onosproject/onos-a1t/pkg/northbound/a1ap/enrichment_information"
)

type a1eiWraper struct {
	version        string
	a1eiController controller.A1EIController
}

func SetRESTA1EIWraper(e *echo.Echo, version string, a1eiController controller.A1EIController) {
	wraper := &a1eiWraper{
		version:        version,
		a1eiController: a1eiController,
	}

	RegisterEIJobHandlers(version, e, wraper)
}

type EIServerInterface interface {
	// Individual EI job
	// (POST /A1-EI/v1/eijobs/{eiJobId})
	PostIndividualEiJobUsingPOST(ctx echo.Context, eiJobId string) error
}

// ServerInterfaceWrapper converts echo contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler EIServerInterface
}

// PostIndividualEiJobUsingPOST converts echo context to params.
func (w *ServerInterfaceWrapper) PostIndividualEiJobUsingPOST(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "eiJobId" -------------
	var eiJobId string

	err = runtime.BindStyledParameterWithLocation("simple", false, "eiJobId", runtime.ParamLocationPath, ctx.Param("eiJobId"), &eiJobId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter eiJobId: %s", err))
	}

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.PostIndividualEiJobUsingPOST(ctx, eiJobId)
	return err
}

// RegisterHandlers adds each server route to the EchoRouter.
func RegisterEIJobHandlers(version string, router a1ei.EchoRouter, si EIServerInterface) {
	registerHandlersWithBaseURL(version, router, si, "")
}

// Registers handlers, and prepends BaseURL to the paths, so that the paths
// can be served under a prefix.
func registerHandlersWithBaseURL(version string, router a1ei.EchoRouter, si EIServerInterface, baseURL string) {

	wrapper := ServerInterfaceWrapper{
		Handler: si,
	}

	router.POST(baseURL+"/A1-EI/"+version+"/eijobs/:eiJobId/notify", wrapper.PostIndividualEiJobUsingPOST)
}

func (a1eiw *a1eiWraper) PostIndividualEiJobUsingPOST(ctx echo.Context, eiJobId string) error {

	eiJobObjNot := make(map[string]interface{})

	if err := ctx.Bind(&eiJobObjNot); err != nil {
		return ctx.JSON(http.StatusOK, err)
	}

	err := a1eiw.a1eiController.HandleEIJobNotify(ctx.Request().Context(), eiJobId, eiJobObjNot)
	if err != nil {
		return err
	}

	return nil
}

/* // ToDo - handle jobIDs by owner as well
// GetEiJobIdsUsingGET request
func (a1eiw *a1eiWraper) GetEiJobIdsUsingGET(ctx context.Context, params *a1ei.GetEiJobIdsUsingGETParams, reqEditors ...a1ei.RequestEditorFn) (*http.Response, error) {

	// no idea for what reqEditor could be used for..
	//for _, item := range reqEditors {
	//	//item
	//}

	if params.EiTypeId != nil {
		//ToDo - ProblemDetails should be included
		return &http.Response{
			Status:     "404 Enrichment Information type is not found",
			StatusCode: http.StatusNotFound,
		}, fmt.Errorf("GetEiJobIdsUsingGET() EItypeID must be defined")
	}

	eiJobIDs, err := a1eiw.a1eiController.HandleGetEIJobs(ctx, *params.EiTypeId)
	if err != nil {
		//ToDo - ProblemDetails should be included
		return &http.Response{
			Status:     "404 Enrichment Information type is not found",
			StatusCode: http.StatusNotFound,
		}, fmt.Errorf("GetEiJobIdsUsingGET() error searching for JobIDs with defined EiTypeID (%v): %v", *params.EiTypeId, err)
	}
	responseBody, err := json.Marshal(eiJobIDs)
	if err != nil {
		return nil, fmt.Errorf("GetEiJobIdsUsingGET() error marshalling response body to bytes: %v", err)
	}

	return &http.Response{
		Status:     "200 OK EI job identifiers",
		StatusCode: http.StatusOK,
		Body:       ioutil.NopCloser(bytes.NewBuffer(responseBody)),
	}, nil
}

// DeleteIndividualEiJobUsingDELETE request
func (a1eiw *a1eiWraper) DeleteIndividualEiJobUsingDELETE(ctx context.Context, eiJobID string, reqEditors ...a1ei.RequestEditorFn) (*http.Response, error) {

	// no idea for what reqEditor could be used for..
	//for _, item := range reqEditors {
	//	//item
	//}

	err := a1eiw.a1eiController.HandleEIJobDelete(ctx, eiJobID)
	if err != nil {
		return &http.Response{
			Status:     "404 Enrichment Information job is not found",
			StatusCode: http.StatusNotFound,
			//ToDo - ProblemDetails should be included
		}, fmt.Errorf("DeleteIndividualEiJobUsingDELETE() error deleting EiJob with JobID %v: %v", eiJobID, err)
	}

	responseBody, err := json.Marshal("Job deleted")
	if err != nil {
		return nil, fmt.Errorf("DeleteIndividualEiJobUsingDELETE() error marshalling response body to bytes: %v", err)
	}

	return &http.Response{
		Status:     "204 OK",
		StatusCode: http.StatusNoContent,
		Body:       ioutil.NopCloser(bytes.NewBuffer(responseBody)),
	}, nil
}

// GetIndividualEiJobUsingGET request
func (a1eiw *a1eiWraper) GetIndividualEiJobUsingGET(ctx context.Context, eiJobID string, reqEditors ...a1ei.RequestEditorFn) (*http.Response, error) {

	// no idea for what reqEditor could be used for..
	//for _, item := range reqEditors {
	//	//item
	//}

	resp, err := a1eiw.a1eiController.HandleGetEIJob(ctx, eiJobID)
	if err != nil {
		return &http.Response{
			Status:     "404 Enrichment Information job is not found",
			StatusCode: http.StatusNotFound,
			//ToDo - ProblemDetails should be included
		}, fmt.Errorf("GetIndividualEiJobUsingGET() error retrieving EiJob with JobID %v: %v", eiJobID, err)
	}

	responseBody, err := json.Marshal(resp)
	if err != nil {
		return nil, fmt.Errorf("GetIndividualEiJobUsingGET() error marshalling response body to bytes: %v", err)
	}

	return &http.Response{
		Status:     "200 OK",
		StatusCode: http.StatusOK,
		Body:       ioutil.NopCloser(bytes.NewBuffer(responseBody)),
	}, nil
}

// PutIndividualEiJobUsingPUT request with any body
func (a1eiw *a1eiWraper) PutIndividualEiJobUsingPUTWithBody(ctx context.Context, eiJobID string, contentType string, body io.Reader, reqEditors ...a1ei.RequestEditorFn) (*http.Response, error) {

	// no idea for what reqEditor could be used for..
	//for _, item := range reqEditors {
	//	//item
	//}
	//body.Read()

	// ToDo - how to deal with body defined through io.Reader??
	//a1eiEntryValue, err := a1eiw.a1eiController.HandleEIJobCreate(ctx, eiJobId, a1ei.EiJobObject(body))
	//if err != nil {
	//	return ctx.JSON(http.StatusNotFound, err)
	//}
	//
	//return ctx.JSON(http.StatusCreated, a1eiEntryValue)
	return nil, fmt.Errorf("PutIndividualEiJobUsingPUTWithBody() method is not (yet) implemented :/")
}

// ToDo - handle EI Job update
func (a1eiw *a1eiWraper) PutIndividualEiJobUsingPUT(ctx context.Context, eiJobID string, body a1ei.PutIndividualEiJobUsingPUTJSONRequestBody, reqEditors ...a1ei.RequestEditorFn) (*http.Response, error) {

	// no idea for what reqEditor could be used for..
	//for _, item := range reqEditors {
	//	//item
	//}

	a1eiEntryValue, err := a1eiw.a1eiController.HandleEIJobCreate(ctx, eiJobID, a1ei.EiJobObject(body))
	if err != nil {
		return &http.Response{
			Status:     "404 Enrichment Information type is not found",
			StatusCode: http.StatusNotFound,
			//ToDo - ProblemDetails should be included
		}, fmt.Errorf("PutIndividualEiJobUsingPUT() error creating EiJob with JobID %v: %v", eiJobID, err)
	}

	responseBody, err := json.Marshal(a1eiEntryValue)
	if err != nil {
		return nil, fmt.Errorf("PutIndividualEiJobUsingPUT() error marshalling response body to bytes: %v", err)
	}

	return &http.Response{
		Status:     "201 OK",
		StatusCode: http.StatusCreated,
		Body:       ioutil.NopCloser(bytes.NewBuffer(responseBody)),
	}, nil
}

// GetEiJobStatusUsingGET request
func (a1eiw *a1eiWraper) GetEiJobStatusUsingGET(ctx context.Context, eiJobID string, reqEditors ...a1ei.RequestEditorFn) (*http.Response, error) {

	// no idea for what reqEditor could be used for..
	//for _, item := range reqEditors {
	//	//item
	//}

	status, err := a1eiw.a1eiController.HandleGetEIJobStatus(ctx, eiJobID)
	if err != nil {
		return &http.Response{
			Status:     "404 Enrichment Information job is not found",
			StatusCode: http.StatusNotFound,
			//ToDo - ProblemDetails should be included
		}, fmt.Errorf("GetEiJobStatusUsingGET() error retrieving EiJob with JobID %v: %v", eiJobID, err)
	}

	responseBody, err := json.Marshal(status)
	if err != nil {
		return nil, fmt.Errorf("GetEiJobStatusUsingGET() error marshalling response body to bytes: %v", err)
	}

	return &http.Response{
		Status:     "200 OK",
		StatusCode: http.StatusOK,
		Body:       ioutil.NopCloser(bytes.NewBuffer(responseBody)),
	}, nil
}

// GetEiTypeIdentifiersUsingGET request
func (a1eiw *a1eiWraper) GetEiTypeIdentifiersUsingGET(ctx context.Context, reqEditors ...a1ei.RequestEditorFn) (*http.Response, error) {

	// no idea for what reqEditor could be used for..
	//for _, item := range reqEditors {
	//	//item
	//}

	eiTypeIDs, err := a1eiw.a1eiController.HandleGetEIJobTypes(ctx)
	if err != nil {
		return &http.Response{
			Status:     "404 Not Found",
			StatusCode: http.StatusNotFound,
			//ToDo - ProblemDetails should be included
		}, fmt.Errorf("GetEiTypeIdentifiersUsingGET() error retrieving EiJob types: %v", err)
	}

	responseBody, err := json.Marshal(eiTypeIDs)
	if err != nil {
		return nil, fmt.Errorf("GetEiTypeIdentifiersUsingGET() error marshalling response body to bytes: %v", err)
	}

	return &http.Response{
		Status:     "200 OK",
		StatusCode: http.StatusOK,
		Body:       ioutil.NopCloser(bytes.NewBuffer(responseBody)),
	}, nil
}

// GetEiTypeUsingGET request
func (a1eiw *a1eiWraper) GetEiTypeUsingGET(ctx context.Context, eiTypeId string, reqEditors ...a1ei.RequestEditorFn) (*http.Response, error) {

	// no idea for what reqEditor could be used for..
	//for _, item := range reqEditors {
	//	//item
	//}

	eiJobIDs, err := a1eiw.a1eiController.HandleGetEIJobs(ctx, eiTypeId)
	if err != nil {
		return &http.Response{
			Status:     "404 Not Found",
			StatusCode: http.StatusNotFound,
			//ToDo - ProblemDetails should be included
		}, fmt.Errorf("GetEiTypeUsingGET() error retrieving EiJobs by EiTypeID (%v): %v", eiTypeId, err)
	}

	responseBody, err := json.Marshal(eiJobIDs)
	if err != nil {
		return nil, fmt.Errorf("GetEiTypeUsingGET() error marshalling response body to bytes: %v", err)
	}

	return &http.Response{
		Status:     "200 OK",
		StatusCode: http.StatusOK,
		Body:       ioutil.NopCloser(bytes.NewBuffer(responseBody)),
	}, nil
}
*/
