// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-Only-1.0

package nonrtric

import (
	"context"
	"os"
	"os/signal"
	"time"

	"github.com/labstack/echo/v4"
)

type RestServer interface {
	Start()
	Stop()
}

type server struct {
	echo    *echo.Echo
	baseURL string
}

func NewRestServer(baseURL string, controller Controller) (RestServer, error) {
	e := echo.New()
	// Log all requests
	// e.Use(echomiddleware.Logger())

	SetRESTA1PMWraper(e, "v1", controller)
	SetRESTA1EIWraper(e, "v1", controller)

	rest := &server{
		baseURL: baseURL,
		echo:    e,
	}
	return rest, nil
}

func (r *server) Start() {
	r.echo.Logger.Fatal(r.echo.Start(r.baseURL))

}

func (r *server) Stop() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	if err := r.echo.Shutdown(ctx); err != nil {
		r.echo.Logger.Fatal(err)
	}
}
