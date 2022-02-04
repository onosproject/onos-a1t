// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0

package sbclient

import "context"

type Client interface {
	Run(ctx context.Context) error
	Close()
}
