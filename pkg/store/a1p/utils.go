// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-Only-1.0

package a1p

import (
	"context"
	"fmt"
)

func GetPoliciesByTypeID(ctx context.Context, s Store, policyTypeID string, ch chan<- *Entry) error {
	policychEntries := make(chan *Entry)
	done := make(chan bool)

	go func(policyTypeID string, ch chan<- *Entry) {

		for entry := range policychEntries {
			policyTypeId := entry.Key.PolicyTypeId
			if policyTypeId == policyTypeID {
				ch <- entry
			}
		}
		done <- true
	}(policyTypeID, ch)

	err := s.Entries(ctx, policychEntries)
	if err != nil {
		close(ch)
		return fmt.Errorf("no policy entries stored for Type %s", policyTypeID)
	}

	<-done
	close(ch)
	return nil
}

func GetPolicyByID(ctx context.Context, s Store, policyID, policyTypeID string) (*Entry, error) {
	a1pKey := Key{
		PolicyId:     policyID,
		PolicyTypeId: policyTypeID,
	}
	a1pEntry, err := s.Get(ctx, a1pKey)
	if err != nil {
		return nil, err
	}

	return a1pEntry, nil
}
