// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-Only-1.0

package controller

import (
	"context"
	"fmt"
	a1eisbi "github.com/onosproject/onos-a1t/pkg/southbound/a1ei"
	a1eistore "github.com/onosproject/onos-a1t/pkg/store/a1ei"
	substore "github.com/onosproject/onos-a1t/pkg/store/subscription"
	"sort"
	// a1einbi "github.com/onosproject/onos-a1t/pkg/northbound/a1ei/enrichment_information"
)

type A1EIController interface {
	HandleEIJobCreate(ctx context.Context, eiJobID, eiJobTypeID string, params map[string]string, eiJobObject map[string]string) error
	HandleEIJobDelete(ctx context.Context, eiTypeID, eiJobTypeID string) error
	HandleEIJobUpdate(ctx context.Context, eiJobID, eiJobTypeID string, params map[string]string, eiJobObject map[string]string) error
	HandleGetEIJobStatus(ctx context.Context, eiJobID, eiJobTypeID string) (bool, error)
	HandleGetEIJobTypes(ctx context.Context) []string
	HandleGetEIJobs(ctx context.Context, eiJobs map[string]string) ([]*a1eistore.Value, error)
	HandleGetEIJob(ctx context.Context, eiJobID, eiJobTypeID string) (*a1eistore.Value, error)
}

type a1eiController struct {
	eijobsStore       a1eistore.Store
	subscriptionStore substore.Store
}

func NewA1EIController(subscriptionStore substore.Store, eijobsStore a1eistore.Store) A1EIController {
	return &a1eiController{
		eijobsStore:       eijobsStore,
		subscriptionStore: subscriptionStore,
	}
}

func (a1ei *a1eiController) HandleEIJobCreate(ctx context.Context, eiJobID, eiJobTypeID string, params map[string]string, eiJobObject map[string]string) error {
	eiJobTypes := getSubscriptionEIJobTypes(ctx, a1ei)

	if _, ok := eiJobTypes[eiJobTypeID]; !ok {
		return fmt.Errorf("eiJobTypeID does not exist")
	}

	ch := make(chan *substore.Entry)
	err := substore.SubscriptionsByTypeID(ctx, a1ei.subscriptionStore, substore.EIJOB, eiJobID, ch)
	if err != nil {
		return err
	}

	eiJobTargets := make(map[string]a1eistore.EIJobTarget)

	eiJobStatus := true

	for subEntry := range ch {
		subValue := subEntry.Value.(substore.Value)
		subAddress := subValue.Client.Address

		eiJobStatusTarget := a1eisbi.CreateEIjob(ctx, subAddress, "", "", eiJobID, eiJobTypeID, eiJobObject)
		if eiJobStatusTarget != nil {
			eiJobStatus = false
		}

		eiJobTarget := a1eistore.EIJobTarget{
			Address:           subAddress,
			EIJobStatusObject: map[string]string{"status": eiJobStatusTarget.Error()},
		}
		eiJobTargets[subAddress] = eiJobTarget
	}

	key := a1eistore.Key{
		EIJobID:   eiJobID,
		EIJobtype: eiJobTypeID,
	}

	value := a1eistore.Value{
		NotificationDestination: params["notificationDestination"],
		EIJobObject:             eiJobObject,
		//EIJobStatusObjects:      eiJobStatus, // Should be Enabled or Disabled per specification
		Targets:     eiJobTargets,
		EIJobStatus: eiJobStatus,
	}

	_, err = a1ei.eijobsStore.Put(ctx, key, value)
	if err != nil {
		return err
	}

	return nil
}

func (a1ei *a1eiController) HandleEIJobDelete(ctx context.Context, eiJobID, eiJobTypeID string) error {

	a1eiEntry, err := a1eistore.GetEIjobByID(ctx, a1ei.eijobsStore, eiJobID, eiJobTypeID)
	if err != nil {
		return err
	}

	a1eiValue := a1eiEntry.Value.(a1eistore.Value)

	for _, targetValue := range a1eiValue.Targets {
		err := a1eisbi.DeleteEIjob(ctx, targetValue.Address, "", "", eiJobID, eiJobTypeID)
		if err != nil {
			log.Warn(err)
		}
	}

	err = a1ei.eijobsStore.Delete(ctx, a1eiEntry.Key)
	if err != nil {
		return err
	}

	return nil
}

// HandleEIJobUpdate should be the same routine as for HandleEIJobCreate
func (a1ei *a1eiController) HandleEIJobUpdate(ctx context.Context, eiJobID, eiJobTypeID string, params map[string]string, eiJobObject map[string]string) error {
	return a1ei.HandleEIJobCreate(ctx, eiJobID, eiJobTypeID, params, eiJobObject)
}

func (a1ei *a1eiController) HandleGetEIJobStatus(ctx context.Context, eiJobID, eiJobTypeID string) (bool, error) {
	a1eiEntry, err := a1eistore.GetEIjobByID(ctx, a1ei.eijobsStore, eiJobID, eiJobTypeID)
	if err != nil {
		return false, err
	}

	a1eiEntryValue := a1eiEntry.Value.(a1eistore.Value)
	a1eiEntryValueStatus := a1eiEntryValue.EIJobStatus
	return a1eiEntryValueStatus, nil
}

func (a1ei *a1eiController) HandleGetEIJobTypes(ctx context.Context) []string {
	eiJobTypes := []string{}

	tmpSubs := getSubscriptionEIJobTypes(ctx, a1ei)

	for k := range tmpSubs {
		eiJobTypes = append(eiJobTypes, k)
	}
	sort.Strings(eiJobTypes)

	return eiJobTypes
}

// HandleGetEIJobs is expecting to get a map of strings, where Key is an EIJobID and a content is EIJobTypeID
func (a1ei *a1eiController) HandleGetEIJobs(ctx context.Context, eiJobs map[string]string) ([]*a1eistore.Value, error) {
	a1eiEntryValues := make([]*a1eistore.Value, 0)

	for eiJobID, eiJobTypeID := range eiJobs {
		a1eiEntry, err := a1eistore.GetEIjobByID(ctx, a1ei.eijobsStore, eiJobID, eiJobTypeID)
		if err != nil {
			return nil, err
		}
		a1eiEntryValue := a1eiEntry.Value.(a1eistore.Value)

		a1eiEntryValues = append(a1eiEntryValues, &a1eiEntryValue)
	}

	return a1eiEntryValues, nil
}

func (a1ei *a1eiController) HandleGetEIJob(ctx context.Context, eiJobID, eiJobTypeID string) (*a1eistore.Value, error) {
	a1eiEntry, err := a1eistore.GetEIjobByID(ctx, a1ei.eijobsStore, eiJobID, eiJobTypeID)
	if err != nil {
		return nil, err
	}

	a1eiEntryValue := a1eiEntry.Value.(a1eistore.Value)
	return &a1eiEntryValue, nil
}

func getSubscriptionEIJobTypes(ctx context.Context, a1ei *a1eiController) map[string]struct{} {
	var exists = struct{}{}
	tmpSubs := make(map[string]struct{})
	ch := make(chan *substore.Entry)

	err := substore.SubscriptionsByType(ctx, a1ei.subscriptionStore, substore.EIJOB, ch)
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
