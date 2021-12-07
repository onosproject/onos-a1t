// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-Only-1.0

package rest

import (
	"github.com/labstack/echo/v4"
	echomiddleware "github.com/labstack/echo/v4/middleware"

	"github.com/onosproject/onos-a1t/pkg/controller"
	handler "github.com/onosproject/onos-a1t/pkg/handler"
	api "github.com/onosproject/onos-a1t/pkg/northbound/a1ap/policy_management"
)

type RestServer struct {
	echo    *echo.Echo
	baseURL string
}

func NewRestServer(baseURL string, broker controller.Broker) (*RestServer, error) {

	// swagger, err := api.GetSwagger()
	// if err != nil {
	// 	fmt.Fprintf(os.Stderr, "Error loading swagger spec\n: %s", err)
	// 	return nil, err
	// }

	a1pWraper := handler.NewA1pWraper("v1", broker.A1PController())

	e := echo.New()
	// Log all requests
	e.Use(echomiddleware.Logger())
	// e.Use(middleware.OapiRequestValidator(swagger))

	api.RegisterHandlers(e, a1pWraper)

	rest := &RestServer{
		baseURL: baseURL,
		echo:    e,
	}

	return rest, nil
}

func (r *RestServer) Start() {
	r.echo.Logger.Fatal(r.echo.Start(r.baseURL))

}
