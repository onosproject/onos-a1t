// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-Only-1.0

package controller

import (
	"context"
	"fmt"
	A1apEnrichmentInformation "github.com/onosproject/onos-a1t/pkg/northbound/a1ap/enrichment_information"
	a1eisbi "github.com/onosproject/onos-a1t/pkg/southbound/a1ei"
	a1eistore "github.com/onosproject/onos-a1t/pkg/store/a1ei"
	substore "github.com/onosproject/onos-a1t/pkg/store/subscription"
	"sort"
	// a1einbi "github.com/onosproject/onos-a1t/pkg/northbound/a1ei/enrichment_information"
)

// ToDo - start return instead of text errors appropriate HTTP response codes..

type A1EIController interface {
	HandleEIJobCreate(ctx context.Context, eiJobID string, eiJobObject A1apEnrichmentInformation.EiJobObject) (*a1eistore.Entry, error)
	HandleEIJobDelete(ctx context.Context, eiTypeID string) error
	//HandleEIJobUpdate(ctx context.Context, eiJobID, eiJobTypeID string, params map[string]string, eiJobObject map[string]string) error
	HandleGetEIJobStatus(ctx context.Context, eiJobID string) (bool, error)
	HandleGetEIJobTypes(ctx context.Context) ([]string, error)
	HandleGetEIJobs(ctx context.Context, eiTypeID string) ([]*string, error)
	HandleGetEIJob(ctx context.Context, eiJobID string) (*a1eistore.Value, error)
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

func (a1ei *a1eiController) HandleEIJobCreate(ctx context.Context, eiJobID string, eiJobObject A1apEnrichmentInformation.EiJobObject) (*a1eistore.Entry, error) {
	eiJobTypes, err := getSubscriptionEIJobTypes(ctx, a1ei)
	if err != nil {
		return nil, err
	}

	if _, ok := eiJobTypes[eiJobObject.EiTypeId]; !ok {
		return nil, fmt.Errorf("eiJobTypeID does not exist")
	}

	ch := make(chan *substore.Entry)
	err = substore.SubscriptionsByTypeID(ctx, a1ei.subscriptionStore, substore.EIJOB, eiJobID, ch)
	if err != nil {
		return nil, err
	}

	eiJobTargets := make(map[string]a1eistore.EIJobTarget)
	eiJobStatus := true

	for subEntry := range ch {
		subValue := subEntry.Value.(substore.Value)
		subAddress := subValue.Client.Address

		//ToDo - eiJobObject should be of type EIJobObject as in onos-a1t/pkg/northbound/a1ap/enrichment_information/a1ap_ei.go
		eiJobStatusTarget := a1eisbi.CreateEIjob(ctx, subAddress, "", "", eiJobID, eiJobObject.EiTypeId, eiJobObject.JobResultUri)
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
		EIJobID: eiJobID,
	}

	value := a1eistore.Value{
		EIJobObject: eiJobObject,
		//EIJobStatusObjects:      eiJobStatus, // ToDo - Should be Enabled or Disabled per specification
		EIJobStatus: eiJobStatus,
		EIJobtype:   eiJobObject.EiTypeId,
		Targets:     eiJobTargets,
	}
	if eiJobObject.JobStatusNotificationUri != nil {
		value.NotificationDestination = *eiJobObject.JobStatusNotificationUri
	}

	entry, err := a1ei.eijobsStore.Put(ctx, key, value)
	if err != nil {
		return nil, err
	}

	return entry, nil
}

func (a1ei *a1eiController) HandleEIJobDelete(ctx context.Context, eiJobID string) error {

	a1eiEntry, err := a1eistore.GetEIjobByID(ctx, a1ei.eijobsStore, eiJobID)
	if err != nil {
		return err
	}

	a1eiValue := a1eiEntry.Value.(a1eistore.Value)

	for _, targetValue := range a1eiValue.Targets {
		//ToDo - eiJobObject should be of type EIJobObject as in onos-a1t/pkg/northbound/a1ap/enrichment_information/a1ap_ei.go
		err := a1eisbi.DeleteEIjob(ctx, targetValue.Address, "", "", eiJobID, a1eiValue.EIJobtype, a1eiValue.EIJobObject.JobResultUri)
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

//// HandleEIJobUpdate should be the same routine as for HandleEIJobCreate
//func (a1ei *a1eiController) HandleEIJobUpdate(ctx context.Context, eiJobID, eiJobTypeID string, params map[string]string, eiJobObject map[string]string) error {
//	return a1ei.HandleEIJobCreate(ctx, eiJobID, eiJobTypeID, params, eiJobObject)
//}

func (a1ei *a1eiController) HandleGetEIJobStatus(ctx context.Context, eiJobID string) (bool, error) {
	a1eiEntry, err := a1eistore.GetEIjobByID(ctx, a1ei.eijobsStore, eiJobID)
	if err != nil {
		return false, err
	}

	a1eiEntryValue := a1eiEntry.Value.(a1eistore.Value)
	a1eiEntryValueStatus := a1eiEntryValue.EIJobStatus
	return a1eiEntryValueStatus, nil
}

func (a1ei *a1eiController) HandleGetEIJobTypes(ctx context.Context) ([]string, error) {
	eiJobTypes := []string{}

	tmpSubs, err := getSubscriptionEIJobTypes(ctx, a1ei)
	if err != nil {
		return nil, err
	}

	for k := range tmpSubs {
		eiJobTypes = append(eiJobTypes, k)
	}
	sort.Strings(eiJobTypes)

	return eiJobTypes, nil
}

// ToDo - not confident in this piece of code, test it out..
// HandleGetEIJobs returning an array of IDs which correspond to EiTypeID
func (a1ei *a1eiController) HandleGetEIJobs(ctx context.Context, eiTypeID string) ([]*string, error) {
	eiJobIDs := make([]*string, 0)

	eiJobchEntries := make(chan *a1eistore.Entry)
	done := make(chan bool)

	go func(ch chan *a1eistore.Entry) {
		for entry := range eiJobchEntries {
			value, ok := entry.Value.(a1eistore.Value)
			if ok {
				if value.EIJobtype == eiTypeID {
					eiJobIDs = append(eiJobIDs, &entry.Key.EIJobID)
				}
			}
		}
		done <- true
	}(eiJobchEntries)

	return eiJobIDs, nil
}

func (a1ei *a1eiController) HandleGetEIJob(ctx context.Context, eiJobID string) (*a1eistore.Value, error) {
	a1eiEntry, err := a1eistore.GetEIjobByID(ctx, a1ei.eijobsStore, eiJobID)
	if err != nil {
		return nil, err
	}

	a1eiEntryValue := a1eiEntry.Value.(a1eistore.Value)
	return &a1eiEntryValue, nil
}

func getSubscriptionEIJobTypes(ctx context.Context, a1ei *a1eiController) (map[string]struct{}, error) {
	var exists = struct{}{}
	tmpSubs := make(map[string]struct{})
	ch := make(chan *substore.Entry)

	err := substore.SubscriptionsByType(ctx, a1ei.subscriptionStore, substore.EIJOB, ch)
	if err != nil {
		return nil, err
	}

	for subEntry := range ch {
		subValue := subEntry.Value.(substore.Value)
		for _, sub := range subValue.Subscriptions {
			if _, ok := tmpSubs[sub.TypeID]; !ok {
				tmpSubs[sub.TypeID] = exists
			}
		}
	}

	return tmpSubs, nil
}
