package nonrtric

import (
	"context"

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
func (c *controller) A1PMGetPolicytypesPolicyTypeId(ctx context.Context, policyTypeId string) (*a1pm.PolicyTypeObject, error) {
	response, err := c.a1pClient.GetPolicytypesPolicyTypeIdWithResponse(ctx, a1pm.PolicyTypeId(policyTypeId))
	if err != nil {
		return &a1pm.PolicyTypeObject{}, err
	}

	return response.JSON200, nil
}

// GetPolicytypesPolicyTypeIdPolicies request
func (c *controller) A1PMGetPolicytypesPolicyTypeIdPolicies(ctx context.Context, policyTypeId string) ([]string, error) {
	return []string{}, nil
}

// DeletePolicytypesPolicyTypeIdPoliciesPolicyId request
func (c *controller) A1PMDeletePolicytypesPolicyTypeIdPoliciesPolicyId(ctx context.Context, policyTypeId, policyId string) error {
	return nil
}

// GetPolicytypesPolicyTypeIdPoliciesPolicyId request
func (c *controller) A1PMGetPolicytypesPolicyTypeIdPoliciesPolicyId(ctx context.Context, policyTypeId, policyId string) ([]string, error) {
	return []string{}, nil
}

// PutPolicytypesPolicyTypeIdPoliciesPolicyId request with any body
func (c *controller) A1PMPutPolicytypesPolicyTypeIdPoliciesPolicyIdWithBody(ctx context.Context, policyTypeId, policyId string, body string) error {
	return nil
}

func (c *controller) A1PMPutPolicytypesPolicyTypeIdPoliciesPolicyId(ctx context.Context, policyTypeId, policyId string, body string) error {
	return nil
}

// GetPolicytypesPolicyTypeIdPoliciesPolicyIdStatus request
func (c *controller) A1PMGetPolicytypesPolicyTypeIdPoliciesPolicyIdStatus(ctx context.Context, policyTypeId, policyId string) ([]string, error) {
	return []string{}, nil
}
