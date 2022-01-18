// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-Only-1.0

package nonrtric

type A1EIValue struct {
	NotificationDestination string
	EIJobObject             map[string]interface{}
	EIJobStatusObjects      map[string]interface{}
}

type A1EIKey struct {
	EIJobID   string
	EIJobtype string
}

type A1EIEntry struct {
	Key   A1EIKey
	Value interface{}
}

type A1PMKey struct {
	PolicyId     string
	PolicyTypeId string
}

type A1PMValue struct {
	NotificationDestination string
	PolicyObject            map[string]interface{}
	PolicyStatusObjects     map[string]interface{}
}

type A1PMEntry struct {
	Key   A1PMKey
	Value interface{}
}
