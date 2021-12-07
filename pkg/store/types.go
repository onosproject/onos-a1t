// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-Only-1.0

package store

// EntityID is used for creating and entity definition
type EntityID struct {
	EntityID int64
}

type Key struct {
	EntityID EntityID
}

type Entry struct {
	Key   Key
	Value interface{}
}

// EntityEvent a entity event
type EntityEvent int

const (
	// None none cell event
	None EntityEvent = iota
	// Created created entity event
	Created
	// Updated updated entity event
	Updated
	// Deleted deleted entity event
	Deleted
)

// Event store event data structure
type Event struct {
	Key   interface{}
	Value interface{}
	Type  interface{}
}
