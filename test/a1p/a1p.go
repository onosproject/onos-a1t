// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-Only-1.0

package a1p

import (
	"context"
	"testing"

	"github.com/onosproject/onos-a1t/test/utils"
	"github.com/onosproject/onos-a1t/test/utils/nonrtric"
	"github.com/stretchr/testify/assert"
)

// TestA1TPMService is the function for Helmit-based integration test
func (s *TestSuite) TestA1TPMService(t *testing.T) {

	mgr, err := nonrtric.NewManager(utils.NonRTRicBaseURL, utils.NearRTRicBaseURL)
	assert.NoError(t, err)

	mgr.Run()

	control := mgr.GetController()

	ctx, cancel := context.WithTimeout(context.Background(), utils.TestTimeout)
	defer cancel()

	policyTypes, err := control.A1PMGetPolicytypes(ctx)
	assert.NoError(t, err)
	assert.EqualValues(t, policyTypes, utils.ExpectedA1PMTypeIDs)

	policyTypeObject, err := control.A1PMGetPolicytypesPolicyTypeId(ctx, utils.PolicyTypeId)
	assert.NoError(t, err)

	err = utils.ValidateSchema(utils.ExpectedPolicyObject, policyTypeObject.PolicySchema)
	assert.NoError(t, err)

	t.Log("A1T Policy Management suite test passed")
}
