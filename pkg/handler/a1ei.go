// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-Only-1.0

package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/onosproject/onos-a1t/pkg/controller"
	a1ei "github.com/onosproject/onos-a1t/pkg/northbound/a1ap/enrichment_information"
	"io"
	"io/ioutil"
	"net/http"
)

type a1eiWraper struct {
	version        string
	a1eiController controller.A1EIController
}

func NewA1eiWraper(version string, a1eiController controller.A1EIController) a1ei.ClientInterface {
	return &a1eiWraper{
		version:        version,
		a1eiController: a1eiController,
	}
}

// ToDo - handle jobIDs by owner as well
// GetEiJobIdsUsingGET request
func (a1eiw *a1eiWraper) GetEiJobIdsUsingGET(ctx echo.Context, params a1ei.GetEiJobIdsUsingGETParams, reqEditors ...a1ei.RequestEditorFn) (*http.Response, error) {

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

	eiJobIDs, err := a1eiw.a1eiController.HandleGetEIJobs(ctx.Request().Context(), *params.EiTypeId)
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
func (a1eiw *a1eiWraper) DeleteIndividualEiJobUsingDELETE(ctx echo.Context, eiJobID string, reqEditors ...a1ei.RequestEditorFn) (*http.Response, error) {

	// no idea for what reqEditor could be used for..
	//for _, item := range reqEditors {
	//	//item
	//}

	err := a1eiw.a1eiController.HandleEIJobDelete(ctx.Request().Context(), eiJobID)
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
func (a1eiw *a1eiWraper) GetIndividualEiJobUsingGET(ctx echo.Context, eiJobID string, reqEditors ...a1ei.RequestEditorFn) (*http.Response, error) {

	// no idea for what reqEditor could be used for..
	//for _, item := range reqEditors {
	//	//item
	//}

	resp, err := a1eiw.a1eiController.HandleGetEIJob(ctx.Request().Context(), eiJobID)
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
func (a1eiw *a1eiWraper) PutIndividualEiJobUsingPUTWithBody(ctx echo.Context, eiJobID string, contentType string, body io.Reader, reqEditors ...a1ei.RequestEditorFn) (*http.Response, error) {

	// no idea for what reqEditor could be used for..
	//for _, item := range reqEditors {
	//	//item
	//}
	//body.Read()

	// ToDo - how to deal with body defined through io.Reader??
	//a1eiEntryValue, err := a1eiw.a1eiController.HandleEIJobCreate(ctx.Request().Context(), eiJobId, a1ei.EiJobObject(body))
	//if err != nil {
	//	return ctx.JSON(http.StatusNotFound, err)
	//}
	//
	//return ctx.JSON(http.StatusCreated, a1eiEntryValue)
	return nil, fmt.Errorf("PutIndividualEiJobUsingPUTWithBody() method is not (yet) implemented :/")
}

// ToDo - handle EI Job update
func (a1eiw *a1eiWraper) PutIndividualEiJobUsingPUT(ctx echo.Context, eiJobID string, body a1ei.PutIndividualEiJobUsingPUTJSONRequestBody, reqEditors ...a1ei.RequestEditorFn) (*http.Response, error) {

	// no idea for what reqEditor could be used for..
	//for _, item := range reqEditors {
	//	//item
	//}

	a1eiEntryValue, err := a1eiw.a1eiController.HandleEIJobCreate(ctx.Request().Context(), eiJobID, a1ei.EiJobObject(body))
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
func (a1eiw *a1eiWraper) GetEiJobStatusUsingGET(ctx echo.Context, eiJobID string, reqEditors ...a1ei.RequestEditorFn) (*http.Response, error) {

	// no idea for what reqEditor could be used for..
	//for _, item := range reqEditors {
	//	//item
	//}

	status, err := a1eiw.a1eiController.HandleGetEIJobStatus(ctx.Request().Context(), eiJobID)
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
func (a1eiw *a1eiWraper) GetEiTypeIdentifiersUsingGET(ctx echo.Context, reqEditors ...a1ei.RequestEditorFn) (*http.Response, error) {

	// no idea for what reqEditor could be used for..
	//for _, item := range reqEditors {
	//	//item
	//}

	eiTypeIDs, err := a1eiw.a1eiController.HandleGetEIJobTypes(ctx.Request().Context())
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
func (a1eiw *a1eiWraper) GetEiTypeUsingGET(ctx echo.Context, eiTypeId string, reqEditors ...a1ei.RequestEditorFn) (*http.Response, error) {

	// no idea for what reqEditor could be used for..
	//for _, item := range reqEditors {
	//	//item
	//}

	eiJobIDs, err := a1eiw.a1eiController.HandleGetEIJobs(ctx.Request().Context(), eiTypeId)
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
