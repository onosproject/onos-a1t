// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-Only-1.0

package handler

import (
	"github.com/onosproject/onos-a1t/pkg/controller"
	a1ei "github.com/onosproject/onos-a1t/pkg/northbound/a1ap/enrichment_information"
)

type a1eiWraper struct {
	version       string
	a1eiController controller.A1EIController
}

func NewA1eiWraper(version string, a1eiController controller.A1EIController) a1ei.ServerInterface {
	return &a1eiWraper{
		version:       version,
		a1eiController: a1eiController,
	}
}

// GetEiJobIdsUsingGET request
GetEiJobIdsUsingGET(ctx context.Context, params *GetEiJobIdsUsingGETParams, reqEditors ...RequestEditorFn) (*http.Response, error)

// DeleteIndividualEiJobUsingDELETE request
DeleteIndividualEiJobUsingDELETE(ctx context.Context, eiJobId string, reqEditors ...RequestEditorFn) (*http.Response, error)

// GetIndividualEiJobUsingGET request
GetIndividualEiJobUsingGET(ctx context.Context, eiJobId string, reqEditors ...RequestEditorFn) (*http.Response, error)

// PutIndividualEiJobUsingPUT request with any body
PutIndividualEiJobUsingPUTWithBody(ctx context.Context, eiJobId string, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error)

PutIndividualEiJobUsingPUT(ctx context.Context, eiJobId string, body PutIndividualEiJobUsingPUTJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error)

// GetEiJobStatusUsingGET request
GetEiJobStatusUsingGET(ctx context.Context, eiJobId string, reqEditors ...RequestEditorFn) (*http.Response, error)

// GetEiTypeIdentifiersUsingGET request
GetEiTypeIdentifiersUsingGET(ctx context.Context, reqEditors ...RequestEditorFn) (*http.Response, error)

// GetEiTypeUsingGET request
GetEiTypeUsingGET(ctx context.Context, eiTypeId string, reqEditors ...RequestEditorFn) (*http.Response, error)
