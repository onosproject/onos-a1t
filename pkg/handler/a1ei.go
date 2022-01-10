// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-Only-1.0

package handler

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/onosproject/onos-a1t/pkg/controller"
	a1ei "github.com/onosproject/onos-a1t/pkg/northbound/a1ap/enrichment_information"
	"io"
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
		//ToDo - Do I have to define http responses this way?
		return &http.Response{
			Status:     "404 Enrichment Information type is not found",
			StatusCode: http.StatusNotFound,
		}, fmt.Errorf("EItypeID must be defined")
	}

	eiJobIDs, err := a1eiw.a1eiController.HandleGetEIJobs(ctx.Request().Context(), *params.EiTypeId)
	if err != nil {
		return &http.Response{
			Status:     "404 Enrichment Information type is not found",
			StatusCode: http.StatusNotFound,
			Body:       eiJobIDs,
		}, err
	}

	return &http.Response{
		Status:     "200 EI job identifiers Enrichment Information type is not found",
		StatusCode: http.StatusNotFound,
	}, nil
}

// DeleteIndividualEiJobUsingDELETE request
func (a1eiw *a1eiWraper) DeleteIndividualEiJobUsingDELETE(ctx echo.Context, eiJobId string, reqEditors ...a1ei.RequestEditorFn) (*http.Response, error) {

	// no idea for what reqEditor could be used for..
	//for _, item := range reqEditors {
	//	//item
	//}

	err := a1eiw.a1eiController.HandleEIJobDelete(ctx.Request().Context(), eiJobId)
	if err != nil {
		return ctx.JSON(http.StatusNotFound, err)
	}

	return ctx.JSON(http.StatusNoContent, err)

}

// GetIndividualEiJobUsingGET request
func (a1eiw *a1eiWraper) GetIndividualEiJobUsingGET(ctx echo.Context, eiJobId string, reqEditors ...a1ei.RequestEditorFn) (*http.Response, error) {

	// no idea for what reqEditor could be used for..
	//for _, item := range reqEditors {
	//	//item
	//}

	resp, err := a1eiw.a1eiController.HandleGetEIJob(ctx.Request().Context(), eiJobId)
	if err != nil {
		return ctx.JSON(http.StatusNotFound, err)
	}

	return ctx.JSON(http.StatusOK, resp)
}

// PutIndividualEiJobUsingPUT request with any body
func (a1eiw *a1eiWraper) PutIndividualEiJobUsingPUTWithBody(ctx echo.Context, eiJobId string, contentType string, body io.Reader, reqEditors ...a1ei.RequestEditorFn) (*http.Response, error) {

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
	return fmt.Errorf("PutIndividualEiJobUsingPUTWithBody() method is not (yet) implemented :/")
}

func (a1eiw *a1eiWraper) PutIndividualEiJobUsingPUT(ctx echo.Context, eiJobId string, body a1ei.PutIndividualEiJobUsingPUTJSONRequestBody, reqEditors ...a1ei.RequestEditorFn) (*http.Response, error) {

	// no idea for what reqEditor could be used for..
	//for _, item := range reqEditors {
	//	//item
	//}

	a1eiEntryValue, err := a1eiw.a1eiController.HandleEIJobCreate(ctx.Request().Context(), eiJobId, a1ei.EiJobObject(body))
	if err != nil {
		return ctx.JSON(http.StatusNotFound, err)
	}

	return ctx.JSON(http.StatusCreated, a1eiEntryValue)
}

// GetEiJobStatusUsingGET request
func (a1eiw *a1eiWraper) GetEiJobStatusUsingGET(ctx echo.Context, eiJobId string, reqEditors ...a1ei.RequestEditorFn) (*http.Response, error) {

	// no idea for what reqEditor could be used for..
	//for _, item := range reqEditors {
	//	//item
	//}

	status, err := a1eiw.a1eiController.HandleGetEIJobStatus(ctx.Request().Context(), eiJobId)
	if err != nil {
		return ctx.JSON(http.StatusNotFound, err)
	}

	return ctx.JSON(http.StatusCreated, status)
}

// GetEiTypeIdentifiersUsingGET request
func (a1eiw *a1eiWraper) GetEiTypeIdentifiersUsingGET(ctx echo.Context, reqEditors ...a1ei.RequestEditorFn) (*http.Response, error) {

	// no idea for what reqEditor could be used for..
	//for _, item := range reqEditors {
	//	//item
	//}

	eiTypeIDs, err := a1eiw.a1eiController.HandleGetEIJobTypes(ctx.Request().Context())
	if err != nil {
		return ctx.JSON(http.StatusNotFound, err)
	}

	return ctx.JSON(http.StatusCreated, eiTypeIDs)
}

// GetEiTypeUsingGET request
func (a1eiw *a1eiWraper) GetEiTypeUsingGET(ctx echo.Context, eiTypeId string, reqEditors ...a1ei.RequestEditorFn) (*http.Response, error) {

	// no idea for what reqEditor could be used for..
	//for _, item := range reqEditors {
	//	//item
	//}

	eiJobIDs, err := a1eiw.a1eiController.HandleGetEIJobs(ctx.Request().Context(), eiTypeId)
	if err != nil {
		return ctx.JSON(http.StatusNotFound, err)
	}

	return ctx.JSON(http.StatusCreated, eiJobIDs)
}
