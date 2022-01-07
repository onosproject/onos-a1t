// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-Only-1.0

package a1ei

import (
	"context"
	"fmt"
)

func GetEIjobByTypeID(ctx context.Context, s Store, eiTypeID string, ch chan<- *Entry) error {
	eiJobChEntries := make(chan *Entry)
	done := make(chan bool)

	go func(policyTypeID string, ch chan<- *Entry) {

		for entry := range eiJobChEntries {
			eiJobTypeID := entry.Key.EIJobtype
			if eiJobTypeID == eiTypeID {
				ch <- entry
			}
		}
		done <- true
	}(eiTypeID, ch)

	err := s.Entries(ctx, eiJobChEntries)
	if err != nil {
		close(ch)
		return fmt.Errorf("no EI Job entries stored for Type %s", eiTypeID)
	}

	<-done
	close(ch)
	return nil
}

func GetEIjobByID(ctx context.Context, s Store, eiJobID, eiTypeID string) (*Entry, error) {
	a1eiKey := Key{
		EIJobID:   eiJobID,
		EIJobtype: eiTypeID,
	}
	a1eiEntry, err := s.Get(ctx, a1eiKey)
	if err != nil {
		return nil, err
	}

	return a1eiEntry, nil
}
