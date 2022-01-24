// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-Only-1.0

package sbclient

import "context"

type Client interface {
	Run(ctx context.Context) error
	Close()
}
