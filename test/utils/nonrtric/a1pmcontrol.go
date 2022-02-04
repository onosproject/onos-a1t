// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0

package nonrtric

import (
	"context"
	"encoding/json"
	"fmt"

	a1pm "github.com/onosproject/onos-a1t/pkg/northbound/a1ap/policy_management"
)

/*
	PM Implementations of Controller interface
*/

// GetPolicytypes request
func (c *controller) A1PMGetPolicytypes(ctx context.Context) ([]string, error) {
	response, err := c.a1pClient.GetPolicytypesWithResponse(ctx)
	if err != nil {
		return []string{}, err
	}

	policyTypes := []string{}
	for _, typeID := range *response.JSON200 {
		policyTypes = append(policyTypes, string(typeID))
	}

	return policyTypes, nil
}

// GetPolicytypesPolicyTypeId request
func (c *controller) A1PMGetPolicytypesPolicyTypeId(ctx context.Context, policyTypeId string) (string, error) {
	response, err := c.a1pClient.GetPolicytypesPolicyTypeIdWithResponse(ctx, a1pm.PolicyTypeId(policyTypeId))
	if err != nil {
		return "", err
	}

	policySchema, err := json.Marshal(response.JSON200.PolicySchema)
	if err != nil {
		return "", err
	}

	return string(policySchema), nil
}

// GetPolicytypesPolicyTypeIdPolicies request
func (c *controller) A1PMGetPolicytypesPolicyTypeIdPolicies(ctx context.Context, policyTypeId string) ([]string, error) {
	response, err := c.a1pClient.GetPolicytypesPolicyTypeIdPoliciesWithResponse(ctx, a1pm.PolicyTypeId(policyTypeId))
	if err != nil {
		return []string{}, err
	}

	policyIDs := []string{}
	for _, polID := range *response.JSON200 {
		policyIDs = append(policyIDs, string(polID))
	}

	return policyIDs, nil

}

// DeletePolicytypesPolicyTypeIdPoliciesPolicyId request
func (c *controller) A1PMDeletePolicytypesPolicyTypeIdPoliciesPolicyId(ctx context.Context, policyTypeId, policyId string) error {
	response, err := c.a1pClient.DeletePolicytypesPolicyTypeIdPoliciesPolicyIdWithResponse(ctx, a1pm.PolicyTypeId(policyTypeId), a1pm.PolicyId(policyId))
	if err != nil {
		return err
	}

	if response.StatusCode() != 204 {
		return fmt.Errorf("policy delete status %d - body %s", response.StatusCode(), response.Body)
	}

	return nil
}

// GetPolicytypesPolicyTypeIdPoliciesPolicyId request
func (c *controller) A1PMGetPolicytypesPolicyTypeIdPoliciesPolicyId(ctx context.Context, policyTypeId, policyId string) (string, error) {

	response, err := c.a1pClient.GetPolicytypesPolicyTypeIdPoliciesPolicyIdWithResponse(ctx, a1pm.PolicyTypeId(policyTypeId), a1pm.PolicyId(policyId))
	if err != nil {
		return "", err
	}

	polObj, err := json.Marshal(response.JSON200)
	if err != nil {
		return "", err
	}

	return string(polObj), nil
}

// PutPolicytypesPolicyTypeIdPoliciesPolicyId request with any body
func (c *controller) A1PMPutPolicytypesPolicyTypeIdPoliciesPolicyIdWithBody(ctx context.Context, policyTypeId, policyId string, body string) error {
	return nil
}

func (c *controller) A1PMPutPolicytypesPolicyTypeIdPoliciesPolicyId(ctx context.Context, policyTypeId, policyId, paramNotifyURL, body string) error {

	notDest := a1pm.NotificationDestination(paramNotifyURL)
	params := &a1pm.PutPolicytypesPolicyTypeIdPoliciesPolicyIdParams{
		NotificationDestination: &notDest,
	}

	bodyInterface := make(map[string]interface{})
	err := json.Unmarshal([]byte(body), &bodyInterface)
	if err != nil {
		return err
	}

	policyBody := a1pm.PutPolicytypesPolicyTypeIdPoliciesPolicyIdJSONRequestBody(a1pm.PutPolicytypesPolicyTypeIdPoliciesPolicyIdJSONBody(a1pm.PolicyObject(bodyInterface)))

	response, err := c.a1pClient.PutPolicytypesPolicyTypeIdPoliciesPolicyIdWithResponse(ctx, a1pm.PolicyTypeId(policyTypeId), a1pm.PolicyId(policyId), params, policyBody)
	if err != nil {
		return err
	}

	if response.StatusCode() != 201 {
		return fmt.Errorf("policy create status %d - body %s", response.StatusCode(), response.Body)
	}

	return nil
}

// GetPolicytypesPolicyTypeIdPoliciesPolicyIdStatus request
func (c *controller) A1PMGetPolicytypesPolicyTypeIdPoliciesPolicyIdStatus(ctx context.Context, policyTypeId, policyId string) (string, error) {

	response, err := c.a1pClient.GetPolicytypesPolicyTypeIdPoliciesPolicyIdStatusWithResponse(ctx, a1pm.PolicyTypeId(policyTypeId), a1pm.PolicyId(policyId))
	if err != nil {
		return "", err
	}

	polStatus, err := json.Marshal(response.JSON200)
	if err != nil {
		return "", err
	}

	return string(polStatus), nil
}
