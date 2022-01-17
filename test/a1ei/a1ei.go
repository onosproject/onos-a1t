// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-Only-1.0

package a1ei

import (
	"testing"

	"github.com/onosproject/onos-a1t/test/utils"
	"github.com/onosproject/onos-a1t/test/utils/nonrtric"
	"github.com/stretchr/testify/assert"
)

// TestA1TEIService is the function for Helmit-based integration test
func (s *TestSuite) TestA1TEIService(t *testing.T) {

	mgr, err := nonrtric.NewManager(utils.NonRTRicBaseURL, utils.NearRTRicBaseURL)
	assert.NoError(t, err)

	mgr.Run()

	/*
		ctx, cancel := context.WithTimeout(context.Background(), utils.TestTimeout)
		defer cancel()

		assert.NoError(t, err)
	*/
	t.Log("A1T Enrichment Information suite test passed")
}
