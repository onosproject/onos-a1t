// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-Only-1.0

package rest

import (
	"context"
	"os"
	"os/signal"
	"time"

	"github.com/labstack/echo/v4"

	"github.com/onosproject/onos-a1t/pkg/controller"
	"github.com/onosproject/onos-a1t/pkg/handler"
)

type Server struct {
	echo    *echo.Echo
	baseURL string
}

func NewRestServer(baseURL string, broker controller.Broker) (*Server, error) {
	e := echo.New()
	// Log all requests
	// e.Use(echomiddleware.Logger())

	handler.SetRESTA1PWraper(e, "v1", broker.A1PController())
	handler.SetRESTA1EIWraper(e, "v1", broker.A1EIController())

	rest := &Server{
		baseURL: baseURL,
		echo:    e,
	}
	return rest, nil
}

func (r *Server) Start() {
	r.echo.Logger.Fatal(r.echo.Start(r.baseURL))

}

func (r *Server) Stop() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	if err := r.echo.Shutdown(ctx); err != nil {
		r.echo.Logger.Fatal(err)
	}
}
