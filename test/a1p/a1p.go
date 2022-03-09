// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0

package a1p

import (
	"context"
	"testing"
	"time"

	"github.com/onosproject/helmit/pkg/helm"
	"github.com/onosproject/onos-a1t/test/utils"
	"github.com/onosproject/onos-a1t/test/utils/nonrtric"
	"github.com/stretchr/testify/assert"
)

var (
	waitPeriod = time.Duration(0)
)

func startA1TXapp(t *testing.T, xappName, policyTypeIDs string) *helm.HelmRelease {
	a1txapp := utils.CreateA1TXapp(t, xappName, policyTypeIDs)
	t.Log("A1T xApp started")
	return a1txapp
}

// // TestA1TPMService is the function for Helmit-based integration test
// func (s *TestSuite) TestA1TPMService(t *testing.T) {

// 	xAppName := "onos-a1txapp"
// 	xAppPolicyTypeIDs := "ORAN_TrafficSteeringPreference_2.0.0"

// 	t.Log("A1T Policy Management suite test started")

// 	mgr, err := nonrtric.NewManager(utils.NonRTRicBaseURL, utils.NearRTRicBaseURL)
// 	assert.NoError(t, err)
// 	mgr.Run()
// 	control := mgr.GetController()

// 	ctx, cancel := context.WithTimeout(context.Background(), utils.TestTimeout)
// 	defer cancel()

// 	// 1. Before starting xApp with subscription to policy typeID
// 	policyTypes, err := control.A1PMGetPolicytypes(ctx)
// 	assert.NoError(t, err)
// 	assert.ElementsMatch(t, []string{}, policyTypes)
// 	t.Log("A1T Policy Management: Expected empty policy type IDs passed")

// 	time.Sleep(waitPeriod * time.Second)

// 	// 2. Start xApp with a1t NBI enabled for policy management
// 	a1txapp := startA1TXapp(t, xAppName, xAppPolicyTypeIDs)
// 	time.Sleep(waitPeriod * time.Second)

// 	// 3. After starting xApp with subscription to policy typeID
// 	policyTypes, err = control.A1PMGetPolicytypes(ctx)
// 	assert.NoError(t, err)
// 	assert.ElementsMatch(t, utils.ExpectedA1PMTypeIDs, policyTypes)
// 	t.Log("A1T Policy Management: Expected PM type IDs passed")

// 	time.Sleep(waitPeriod * time.Second)

// 	// 4. Retrieves policy typeID schema
// 	policySchema, err := control.A1PMGetPolicytypesPolicyTypeId(ctx, utils.PolicyTypeId)
// 	// t.Log("Received policySchema ", policySchema)
// 	assert.NoError(t, err)

// 	// 5. Validates policy typeID schema
// 	err = utils.ValidateSchema(utils.ExpectedPolicyObject, policySchema)
// 	assert.NoError(t, err)
// 	t.Log("A1T Policy Management: Expected policy type schema passed")

// 	// 6. Query policy IDs with a particular policy typeID
// 	existingPolicyIDs, err := control.A1PMGetPolicytypesPolicyTypeIdPolicies(ctx, utils.PolicyTypeId)
// 	assert.NoError(t, err)
// 	assert.ElementsMatch(t, []string{}, existingPolicyIDs)
// 	t.Log("A1T Policy Management: Expected policy IDs for typeID passed")

// 	time.Sleep(waitPeriod * time.Second)

// 	// 7. Creates policy
// 	err = control.A1PMPutPolicytypesPolicyTypeIdPoliciesPolicyId(ctx, utils.PolicyTypeId, utils.PolicyId, utils.NotifyURL, utils.ExpectedPolicyObject)
// 	assert.NoError(t, err)
// 	t.Log("A1T Policy Management: Expected policy create passed")

// 	time.Sleep(waitPeriod * time.Second)

// 	// 8. Query policy IDs with a particular policy typeID (must be list with expected policy ID)
// 	PolicyIDs, err := control.A1PMGetPolicytypesPolicyTypeIdPolicies(ctx, utils.PolicyTypeId)
// 	assert.NoError(t, err)
// 	assert.ElementsMatch(t, utils.ExpectedA1PMPolicyIDs, PolicyIDs)
// 	t.Log("A1T Policy Management: Expected policy IDs for typeID passed")

// 	time.Sleep(waitPeriod * time.Second)

// 	// 7. Query policy status
// 	policyStatus, err := control.A1PMGetPolicytypesPolicyTypeIdPoliciesPolicyIdStatus(ctx, utils.PolicyTypeId, utils.PolicyId)
// 	assert.NoError(t, err)
// 	err = utils.ValidateSchema(policyStatus, utils.ExpectedPolicyStatusSchema)
// 	assert.NoError(t, err)
// 	t.Log("A1T Policy Management: Expected policy status passed")

// 	time.Sleep(waitPeriod * time.Second)

// 	// 8. Query policy
// 	policyObject, err := control.A1PMGetPolicytypesPolicyTypeIdPoliciesPolicyId(ctx, utils.PolicyTypeId, utils.PolicyId)
// 	assert.NoError(t, err)
// 	err = utils.ValidateSchema(policyObject, utils.ExpectedPolicySchema)
// 	assert.NoError(t, err)
// 	t.Log("A1T Policy Management: Expected policy object query passed")

// 	time.Sleep(waitPeriod * time.Second)

// 	// 9. Delete policy
// 	err = control.A1PMDeletePolicytypesPolicyTypeIdPoliciesPolicyId(ctx, utils.PolicyTypeId, utils.PolicyId)
// 	assert.NoError(t, err)
// 	t.Log("A1T Policy Management: Expected policy delete passed")

// 	time.Sleep(waitPeriod * time.Second)

// 	// 10. Query policy IDs with a particular policy typeID
// 	emptyPolicyIDs, err := control.A1PMGetPolicytypesPolicyTypeIdPolicies(ctx, utils.PolicyTypeId)
// 	assert.NoError(t, err)
// 	assert.ElementsMatch(t, []string{}, emptyPolicyIDs)
// 	t.Log("A1T Policy Management: Expected policy IDs for typeID passed")

// 	time.Sleep(waitPeriod * time.Second)

// 	t.Log("A1T Policy Management suite test finished")

// 	err = a1txapp.Uninstall()
// 	assert.NoError(t, err, "could not uninstall a1txapp %v", err)
// 	t.Log("A1T xApp stopped")

// 	mgr.End()
// 	t.Log("Ended A1T non-RT RIC manager")

// }

// TestA1TPMTwoxAppService is the function for Helmit-based integration test
// for 2 xApps interacting with onos-a1t
// TODO: Test case that supports xApp uninstall/install TBD
func (s *TestSuite) TestA1TPMTwoxAppService(t *testing.T) {

	xApp1Name := "onos-a1txapp2"
	xApp1PolicyTypeIDs := "ORAN_TrafficSteeringPreference_2.0.0"

	xApp2Name := "onos-a1txapp3"
	xApp2PolicyTypeIDs := "ORAN_TrafficSteeringPreference_2.0.0"

	t.Log("A1T Policy Management 2 xApps suite test started")

	mgr, err := nonrtric.NewManager(utils.NonRTRicBaseURL, utils.NearRTRicBaseURL)
	assert.NoError(t, err)
	mgr.Run()
	control := mgr.GetController()

	ctx, cancel := context.WithTimeout(context.Background(), utils.TestTimeout)
	defer cancel()

	// 1. Before starting xApp with subscription to policy typeID
	policyTypes, err := control.A1PMGetPolicytypes(ctx)
	assert.NoError(t, err)
	assert.ElementsMatch(t, []string{}, policyTypes)
	t.Log("A1T Policy Management 2xApps: Expected empty policy type IDs passed")

	time.Sleep(waitPeriod * time.Second)

	// 2. Start xApp with a1t NBI enabled for policy management
	a1txapp1 := startA1TXapp(t, xApp1Name, xApp1PolicyTypeIDs)
	time.Sleep(waitPeriod * time.Second)

	a1txapp2 := startA1TXapp(t, xApp2Name, xApp2PolicyTypeIDs)
	time.Sleep(waitPeriod * time.Second)

	// 3. After starting xApp with subscription to policy typeID
	policyTypes, err = control.A1PMGetPolicytypes(ctx)
	assert.NoError(t, err)
	assert.ElementsMatch(t, utils.ExpectedA1PMTypeIDs, policyTypes)
	t.Log("A1T Policy Management: Expected PM type IDs passed")

	// 4. Retrieves policy typeID schema
	policySchema, err := control.A1PMGetPolicytypesPolicyTypeId(ctx, utils.PolicyTypeId)
	// t.Log("Received policySchema ", policySchema)
	assert.NoError(t, err)

	// 5. Validates policy typeID schema
	err = utils.ValidateSchema(utils.ExpectedPolicyObject, policySchema)
	assert.NoError(t, err)
	t.Log("A1T Policy Management: Expected policy type schema passed")

	time.Sleep(waitPeriod * time.Second)

	// 6. Query policy IDs with a particular policy typeID
	existingPolicyIDs, err := control.A1PMGetPolicytypesPolicyTypeIdPolicies(ctx, utils.PolicyTypeId)
	assert.NoError(t, err)
	assert.ElementsMatch(t, []string{}, existingPolicyIDs)
	t.Log("A1T Policy Management: Expected policy IDs for typeID passed")

	time.Sleep(waitPeriod * time.Second)

	// 7. Creates policy
	err = control.A1PMPutPolicytypesPolicyTypeIdPoliciesPolicyId(ctx, utils.PolicyTypeId, utils.PolicyId, utils.NotifyURL, utils.ExpectedPolicyObject)
	assert.NoError(t, err)
	t.Log("A1T Policy Management: Expected policy create passed")

	time.Sleep(waitPeriod * time.Second)

	// 8. Query policy IDs with a particular policy typeID (must be list with expected policy ID)
	PolicyIDs, err := control.A1PMGetPolicytypesPolicyTypeIdPolicies(ctx, utils.PolicyTypeId)
	assert.NoError(t, err)
	assert.ElementsMatch(t, utils.ExpectedA1PMPolicyIDs, PolicyIDs)
	t.Log("A1T Policy Management: Expected policy IDs for typeID passed")

	time.Sleep(waitPeriod * time.Second)

	// 7. Query policy status
	policyStatus, err := control.A1PMGetPolicytypesPolicyTypeIdPoliciesPolicyIdStatus(ctx, utils.PolicyTypeId, utils.PolicyId)
	assert.NoError(t, err)
	err = utils.ValidateSchema(policyStatus, utils.ExpectedPolicyStatusSchema)
	assert.NoError(t, err)
	t.Log("A1T Policy Management: Expected policy status passed")

	time.Sleep(waitPeriod * time.Second)

	// 8. Query policy
	policyObject, err := control.A1PMGetPolicytypesPolicyTypeIdPoliciesPolicyId(ctx, utils.PolicyTypeId, utils.PolicyId)
	assert.NoError(t, err)
	err = utils.ValidateSchema(policyObject, utils.ExpectedPolicySchema)
	assert.NoError(t, err)
	t.Log("A1T Policy Management: Expected policy object query passed")

	time.Sleep(waitPeriod * time.Second)

	// 9. Delete policy
	err = control.A1PMDeletePolicytypesPolicyTypeIdPoliciesPolicyId(ctx, utils.PolicyTypeId, utils.PolicyId)
	assert.NoError(t, err)
	t.Log("A1T Policy Management: Expected policy delete passed")

	time.Sleep(waitPeriod * time.Second)

	// 10. Query policy IDs with a particular policy typeID
	emptyPolicyIDs, err := control.A1PMGetPolicytypesPolicyTypeIdPolicies(ctx, utils.PolicyTypeId)
	assert.NoError(t, err)
	assert.ElementsMatch(t, []string{}, emptyPolicyIDs)
	t.Log("A1T Policy Management: Expected policy IDs for typeID passed")

	time.Sleep(waitPeriod * time.Second)

	t.Log("A1T Policy Management suite test finished")

	err = a1txapp1.Uninstall()
	assert.NoError(t, err, "could not uninstall a1txapp1 %v", err)
	err = a1txapp2.Uninstall()
	assert.NoError(t, err, "could not uninstall a1txapp2 %v", err)
	t.Log("A1T xApps stopped")

	mgr.End()
	t.Log("Ended A1T non-RT RIC manager")
}
