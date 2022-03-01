// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0

package a1p

import (
	"context"
	"github.com/onosproject/onos-a1t/test/utils"
	"github.com/onosproject/onos-a1t/test/utils/nonrtric"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

// TestRimdeoTSxApp is the function for Helmit-based integration test
func (s *TestSuite) TestRimdeoTSxApp(t *testing.T) {

	t.Log("RIMEDO Lab TS xApp suite test started")
	mgr, err := nonrtric.NewManager(utils.NonRTRicBaseURL, utils.NearRTRicBaseURL)
	assert.NoError(t, err)
	mgr.Run()
	control := mgr.GetController()

	ctx, cancel := context.WithTimeout(context.Background(), utils.TestTimeout)
	defer cancel()

	sim := utils.CreateRanSimulator(t)
	assert.NotNil(t, sim)
	time.Sleep(waitPeriod * time.Second)

	e2t := utils.CreateE2T(t)
	assert.NotNil(t, e2t)
	time.Sleep(waitPeriod * time.Second)

	a1t := utils.CreateA1T(t)
	assert.NotNil(t, a1t)
	time.Sleep(waitPeriod * time.Second)

	time.Sleep(waitPeriod * time.Second)

	// Read policies and make sure there is no policies
	policyTypes, err := control.A1PMGetPolicytypes(ctx)
	assert.NoError(t, err)
	assert.ElementsMatch(t, []string{}, policyTypes)
	t.Log("A1T Policy Management: Expected empty policy type IDs passed")

	time.Sleep(waitPeriod * time.Second)

	// Run TS xApp
	ts := utils.CreateTSxApp(t)
	assert.NotNil(t, ts)
	time.Sleep(waitPeriod * time.Second)

	// ToDo - Download policy in xApp (should be with curl)

	// Wait for 60 seconds or more (max. 2 minutes) to this policy to apply (TS should happen)
	time.Sleep(time.Minute)

	// Read out policies from A1T and make sure there are some policies
	policyTypesNew, err := control.A1PMGetPolicytypes(ctx)
	assert.NoError(t, err)
	assert.NotNil(t, policyTypesNew)
	t.Logf("A1T Policy Management: Expected one stored policy\n%v", policyTypesNew)
}
