// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-Only-1.0

package controller

import (
	"context"
	"fmt"
	"sort"

	"github.com/onosproject/onos-lib-go/pkg/logging"

	a1psbi "github.com/onosproject/onos-a1t/pkg/southbound/a1p"
	a1pstore "github.com/onosproject/onos-a1t/pkg/store/a1p"
	substore "github.com/onosproject/onos-a1t/pkg/store/subscription"
	// a1pnbi "github.com/onosproject/onos-a1t/pkg/northbound/a1ap/policy_management"
)

var log = logging.GetLogger("controller", "a1p")

type A1PController interface {
	HandlePolicyCreate(ctx context.Context, policyID, policyTypeID string, params map[string]string, policyObject map[string]string) error
	HandlePolicyDelete(ctx context.Context, policyID, policyTypeID string) error
	HandlePolicyUpdate() error
	HandleGetPolicyTypes(ctx context.Context) []string
	HandleGetPoliciesTypeID(ctx context.Context, policyTypeID string) ([]*a1pstore.Value, error)
	HandleGetPolicy(ctx context.Context, policyID, policyTypeID string) (*a1pstore.Value, error)
	HandleGetPolicyStatus(ctx context.Context, policyID, policyTypeID string) (bool, error)
}

type a1pController struct {
	policiesStore     a1pstore.Store
	subscriptionStore substore.Store
}

func NewA1PController(subscriptionStore substore.Store, policiesStore a1pstore.Store) A1PController {
	return &a1pController{
		policiesStore:     policiesStore,
		subscriptionStore: subscriptionStore,
	}
}

func (a1p *a1pController) HandleGetPolicyTypes(ctx context.Context) []string {
	policyTypes := []string{}

	tmpSubs := getSubscriptionPolicyTypes(ctx, a1p)

	for k := range tmpSubs {
		policyTypes = append(policyTypes, k)
	}
	sort.Strings(policyTypes)

	return policyTypes
}

func (a1p *a1pController) HandleGetPoliciesTypeID(ctx context.Context, policyTypeID string) ([]*a1pstore.Value, error) {
	policyEntries := []*a1pstore.Value{}
	policychEntries := make(chan *a1pstore.Entry)
	done := make(chan bool)

	go func(ch chan *a1pstore.Entry) {
		for entry := range policychEntries {
			value, ok := entry.Value.(a1pstore.Value)
			if ok {
				policyEntries = append(policyEntries, &value)
			}
		}
		done <- true
	}(policychEntries)

	err := a1pstore.GetPoliciesByTypeID(ctx, a1p.policiesStore, policyTypeID, policychEntries)
	if err != nil {
		close(policychEntries)
		return policyEntries, err
	}

	<-done
	return policyEntries, nil
}

func (a1p *a1pController) HandleGetPolicy(ctx context.Context, policyID, policyTypeID string) (*a1pstore.Value, error) {
	a1pEntry, err := a1pstore.GetPolicyByID(ctx, a1p.policiesStore, policyID, policyTypeID)
	if err != nil {
		return nil, err
	}

	a1pEntryValue := a1pEntry.Value.(a1pstore.Value)
	return &a1pEntryValue, nil
}

func (a1p *a1pController) HandleGetPolicyStatus(ctx context.Context, policyID, policyTypeID string) (bool, error) {
	a1pEntry, err := a1pstore.GetPolicyByID(ctx, a1p.policiesStore, policyID, policyTypeID)
	if err != nil {
		return false, err
	}

	a1pEntryValue := a1pEntry.Value.(a1pstore.Value)
	a1pEntryValueStatus := a1pEntryValue.PolicyStatus
	return a1pEntryValueStatus, nil
}

func (a1p *a1pController) HandlePolicyCreate(ctx context.Context, policyID, policyTypeID string, params map[string]string, policyObject map[string]string) error {
	policyTypes := getSubscriptionPolicyTypes(ctx, a1p)

	if _, ok := policyTypes[policyTypeID]; !ok {
		return fmt.Errorf("policyTypeID does not exist")
	}

	ch := make(chan *substore.Entry)
	err := substore.SubscriptionsByTypeID(ctx, a1p.subscriptionStore, substore.POLICY, policyTypeID, ch)
	if err != nil {
		return err
	}

	policyTargets := make(map[string]a1pstore.PolicyTarget)

	policyStatus := true

	for subEntry := range ch {
		subValue := subEntry.Value.(substore.Value)
		subAddress := subValue.Client.Address

		policyStatusTarget := a1psbi.CreatePolicy(ctx, subAddress, "", "", policyID, policyTypeID, policyObject)
		if policyStatusTarget != nil {
			policyStatus = false
		}

		policyTarget := a1pstore.PolicyTarget{
			Address:            subAddress,
			PolicyStatusObject: map[string]string{"status": policyStatusTarget.Error()},
		}
		policyTargets[subAddress] = policyTarget
	}

	a1pKey := a1pstore.Key{
		PolicyId:     policyID,
		PolicyTypeId: policyTypeID,
	}
	a1pValue := a1pstore.Value{
		NotificationDestination: params["notificationDestination"],
		PolicyObject:            policyObject,
		Targets:                 policyTargets,
		PolicyStatus:            policyStatus,
	}

	_, err = a1p.policiesStore.Put(ctx, a1pKey, a1pValue)
	if err != nil {
		return err
	}

	return nil
}

func (a1p *a1pController) HandlePolicyDelete(ctx context.Context, policyID, policyTypeID string) error {

	a1pEntry, err := a1pstore.GetPolicyByID(ctx, a1p.policiesStore, policyID, policyTypeID)
	if err != nil {
		return err
	}

	a1pValue := a1pEntry.Value.(a1pstore.Value)

	for _, targetValue := range a1pValue.Targets {
		err := a1psbi.DeletePolicy(ctx, targetValue.Address, "", "", policyID, policyTypeID)
		if err != nil {
			log.Warn(err)
		}
	}

	err = a1p.policiesStore.Delete(ctx, a1pEntry.Key)
	if err != nil {
		return err
	}

	return nil
}

func (a1p *a1pController) HandlePolicyUpdate() error {
	return nil
}

func getSubscriptionPolicyTypes(ctx context.Context, a1p *a1pController) map[string]struct{} {
	var exists = struct{}{}
	tmpSubs := make(map[string]struct{})
	ch := make(chan *substore.Entry)

	err := substore.SubscriptionsByType(ctx, a1p.subscriptionStore, substore.POLICY, ch)
	if err != nil {
		return tmpSubs
	}

	for subEntry := range ch {
		subValue := subEntry.Value.(substore.Value)
		for _, sub := range subValue.Subscriptions {
			if _, ok := tmpSubs[sub.TypeID]; !ok {
				tmpSubs[sub.TypeID] = exists
			}
		}
	}

	return tmpSubs
}
