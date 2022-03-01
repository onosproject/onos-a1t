// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0

package utils

import (
	"context"
	"testing"

	"github.com/onosproject/helmit/pkg/helm"
	"github.com/onosproject/helmit/pkg/input"
	"github.com/onosproject/helmit/pkg/kubernetes"
	"github.com/onosproject/onos-test/pkg/onostest"
	"github.com/stretchr/testify/assert"
)

func getCredentials() (string, string, error) {
	kubClient, err := kubernetes.New()
	if err != nil {
		return "", "", err
	}
	secrets, err := kubClient.CoreV1().Secrets().Get(context.Background(), onostest.SecretsName)
	if err != nil {
		return "", "", err
	}
	username := string(secrets.Object.Data["sd-ran-username"])
	password := string(secrets.Object.Data["sd-ran-password"])

	return username, password, nil
}

// CreateSdranRelease creates a helm release for an sd-ran instance
func CreateSdranRelease(c *input.Context) (*helm.HelmRelease, error) {
	username, password, err := getCredentials()
	registry := c.GetArg("registry").String("")

	if err != nil {
		return nil, err
	}

	sdran := helm.Chart("sd-ran", onostest.SdranChartRepo).
		Release("sd-ran").
		SetUsername(username).
		SetPassword(password).
		Set("import.onos-a1t.enabled", true).
		Set("import.onos-topo.enabled", true).
		Set("import.ran-simulator.enabled", false).
		Set("import.onos-kpimon.enabled", false).
		Set("global.image.tag", "latest").
		Set("global.image.registry", registry)

	return sdran, nil
}

// CreateRanSimulatorWithName creates a ran simulator
func CreateA1TXapp(t *testing.T, name string) *helm.HelmRelease {

	a1txapp := helm.
		Chart("./a1txapp").
		Release(name).
		Set("image.tag", "latest").
		Set("fullnameOverride", "")
	err := a1txapp.Install(true)
	assert.NoError(t, err, "could not install a1txapp %v", err)

	return a1txapp
}

// CreateRanSimulator creates a ran simulator
func CreateRanSimulator(t *testing.T) *helm.HelmRelease {
	username, password, err := getCredentials()
	assert.NoError(t, err)

	simulator := helm.
		Chart("ran-simulator", onostest.SdranChartRepo).
		Release("ran-simulator").
		SetUsername(username).
		SetPassword(password).
		Set("image.tag", "latest").
		Set("fullnameOverride", "")
	err = simulator.Install(true)
	assert.NoError(t, err, "could not install device simulator %v", err)

	return simulator
}

// CreateTSxApp creates a ran simulator
func CreateTSxApp(t *testing.T) *helm.HelmRelease {
	username, password, err := getCredentials()
	assert.NoError(t, err)

	ts := helm.
		Chart("rimedo-ts", onostest.SdranChartRepo).
		Release("rimedo-ts").
		SetUsername(username).
		SetPassword(password).
		Set("image.tag", "latest").
		Set("fullnameOverride", "")
	err = ts.Install(true)
	assert.NoError(t, err, "could not install TS xApp %v", err)

	return ts
}

// CreateA1T creates a ran simulator
func CreateA1T(t *testing.T) *helm.HelmRelease {
	username, password, err := getCredentials()
	assert.NoError(t, err)

	a1t := helm.
		Chart("onos-a1t", onostest.SdranChartRepo).
		Release("onos-a1t").
		SetUsername(username).
		SetPassword(password).
		Set("image.tag", "latest").
		Set("fullnameOverride", "")
	err = a1t.Install(true)
	assert.NoError(t, err, "could not install A1T %v", err)

	return a1t
}

// CreateE2T creates a ran simulator
func CreateE2T(t *testing.T) *helm.HelmRelease {
	username, password, err := getCredentials()
	assert.NoError(t, err)

	e2t := helm.
		Chart("onos-e2t", onostest.SdranChartRepo).
		Release("onos-e2t").
		SetUsername(username).
		SetPassword(password).
		Set("image.tag", "latest").
		Set("fullnameOverride", "")
	err = e2t.Install(true)
	assert.NoError(t, err, "could not install E2T %v", err)

	return e2t
}
