// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-Only-1.0

package a1ei

import (
	"context"
)

func GetEIjobByID(ctx context.Context, s Store, eiJobID string) (*Entry, error) {

	a1eiKey := Key{
		EIJobID: eiJobID,
	}
	a1eiEntry, err := s.Get(ctx, a1eiKey)
	if err != nil {
		return nil, err
	}

	return a1eiEntry, nil
}
