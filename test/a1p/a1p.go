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

	t.Log("A1T Policy Management suite test started")

	mgr, err := nonrtric.NewManager(utils.NonRTRicBaseURL, utils.NearRTRicBaseURL)
	assert.NoError(t, err)

	mgr.Run()

	control := mgr.GetController()

	ctx, cancel := context.WithTimeout(context.Background(), utils.TestTimeout)
	defer cancel()

	// 1. Before starting xApp with subscription to policy typeID
	policyTypes, err := control.A1PMGetPolicytypes(ctx)
	assert.NoError(t, err)
	assert.EqualValues(t, policyTypes, []string{})
	t.Log("A1T Policy Management: Expected empty policy type IDs passed")

	// 2. Start xApp with a1t NBI enabled for policy management
	// For instance: s.sdran.Set("import.onos-kpimon.enabled", true)

	// 3. After starting xApp with subscription to policy typeID
	policyTypes, err = control.A1PMGetPolicytypes(ctx)
	assert.NoError(t, err)
	assert.EqualValues(t, policyTypes, utils.ExpectedA1PMTypeIDs)
	t.Log("A1T Policy Management: Expected PM type IDs passed")

	// 4. Retrieves policy typeID schema
	policySchema, err := control.A1PMGetPolicytypesPolicyTypeId(ctx, utils.PolicyTypeId)
	assert.NoError(t, err)

	// 5. Validates policy typeID schema
	err = utils.ValidateSchema(utils.ExpectedPolicyObject, policySchema)
	assert.NoError(t, err)
	t.Log("A1T Policy Management: Expected policy type schema passed")

	// 6. Query policy IDs with a particular policy typeID (must be empty)
	emptyPolicyIDs, err := control.A1PMGetPolicytypesPolicyTypeIdPolicies(ctx, utils.PolicyTypeId)
	assert.NoError(t, err)
	assert.EqualValues(t, emptyPolicyIDs, []string{})
	t.Log("A1T Policy Management: Expected empty policy IDs for typeID passed")

	// 7. Creates policy
	err = control.A1PMPutPolicytypesPolicyTypeIdPoliciesPolicyId(ctx, utils.PolicyTypeId, utils.PolicyId, utils.NotifyURL, utils.ExpectedPolicyObject)
	assert.NoError(t, err)
	t.Log("A1T Policy Management: Expected policy create passed")

	// 8. Query policy IDs with a particular policy typeID (must be list with expected policy ID)
	PolicyIDs, err := control.A1PMGetPolicytypesPolicyTypeIdPolicies(ctx, utils.PolicyTypeId)
	assert.NoError(t, err)
	assert.EqualValues(t, PolicyIDs, utils.ExpectedA1PMPolicyIDs)
	t.Log("A1T Policy Management: Expected policy IDs for typeID passed")

	// 7. Query policy status
	policyStatus, err := control.A1PMGetPolicytypesPolicyTypeIdPoliciesPolicyIdStatus(ctx, utils.PolicyTypeId, utils.PolicyId)
	assert.NoError(t, err)
	err = utils.ValidateSchema(policyStatus, utils.ExpectedPolicyStatusSchema)
	assert.NoError(t, err)
	t.Log("A1T Policy Management: Expected policy status passed")

	// 8. Query policy
	policyObject, err := control.A1PMGetPolicytypesPolicyTypeIdPoliciesPolicyId(ctx, utils.PolicyTypeId, utils.PolicyId)
	assert.NoError(t, err)
	err = utils.ValidateSchema(policyObject, utils.ExpectedPolicySchema)
	assert.NoError(t, err)
	t.Log("A1T Policy Management: Expected policy object query passed")

	// 9. Delete policy
	err = control.A1PMDeletePolicytypesPolicyTypeIdPoliciesPolicyId(ctx, utils.PolicyTypeId, utils.PolicyId)
	assert.NoError(t, err)
	t.Log("A1T Policy Management: Expected policy delete passed")

	// 10. Query policy IDs with a particular policy typeID (must be empty)
	emptyPolicyIDs, err = control.A1PMGetPolicytypesPolicyTypeIdPolicies(ctx, utils.PolicyTypeId)
	assert.NoError(t, err)
	assert.EqualValues(t, emptyPolicyIDs, []string{})
	t.Log("A1T Policy Management: Expected empty policy IDs for typeID passed")

	t.Log("A1T Policy Management suite test finished")
}
