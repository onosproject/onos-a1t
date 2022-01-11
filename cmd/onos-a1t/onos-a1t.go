// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-Only-1.0

package main

import (
	"flag"

	"github.com/onosproject/onos-a1t/pkg/manager"
	"github.com/onosproject/onos-lib-go/pkg/certs"
	"github.com/onosproject/onos-lib-go/pkg/logging"
)

var log = logging.GetLogger("main")

func main() {
	caPath := flag.String("caPath", "", "path to CA certificate")
	keyPath := flag.String("keyPath", "", "path to client private key")
	certPath := flag.String("certPath", "", "path to client certificate")
	configPath := flag.String("configPath", "/etc/onos/config/config.json", "path to config.json file")
	grpcPort := flag.Int("grpcPort", 5150, "grpc Port number")
	baseURL := flag.String("baseURL", "0.0.0.0:9639", "base URL for NBI A1T restfull server")
	nonRTRICURL := flag.String("nonRTRICURL", "127.0.0.1:9640", "base URL of A1 in Non-RT RIC")

	ready := make(chan bool)

	flag.Parse()

	_, err := certs.HandleCertPaths(*caPath, *keyPath, *certPath, true)
	if err != nil {
		log.Fatal(err)
	}

	cfg := manager.Config{
		CAPath:      *caPath,
		KeyPath:     *keyPath,
		CertPath:    *certPath,
		GRPCPort:    *grpcPort,
		ConfigPath:  *configPath,
		BaseURL:     *baseURL,
		NonRTRICURL: *nonRTRICURL,
	}

	log.Info("Starting onos-a1t")
	mgr, err := manager.NewManager(cfg)
	if err != nil {
		log.Fatal(err)
	}

	mgr.Run()
	<-ready
}
