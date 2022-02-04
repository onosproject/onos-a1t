// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"github.com/onosproject/helmit/pkg/registry"
	"github.com/onosproject/helmit/pkg/test"
	"github.com/onosproject/onos-a1t/test/a1ei"
	"github.com/onosproject/onos-a1t/test/a1p"
)

func main() {
	registry.RegisterTestSuite("a1pm", &a1p.TestSuite{})
	registry.RegisterTestSuite("a1ei", &a1ei.TestSuite{})
	test.Main()
}
