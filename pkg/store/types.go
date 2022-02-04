// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0

package store

import topoapi "github.com/onosproject/onos-api/go/onos/topo"

type Entry struct {
	Key   interface{}
	Value interface{}
}

// For watcher

type Event struct {
	Key   interface{}
	Value interface{}
	Type  interface{}
}

type EventType int

const (
	// None none cell event
	None EventType = iota
	// Created created entity event
	Created
	// Updated updated entity event
	Updated
	// Deleted deleted entity event
	Deleted
)

func (e EventType) String() string {
	return [...]string{"None", "Created", "Update", "Deleted"}[e]
}

// For service definition

type A1Service int

const (
	PolicyManagement A1Service = iota
	EnrichmentInformation
)

func (s A1Service) String() string {
	return [...]string{"PolicyManagement", "EnrichmentInformation"}[s]
}

type A1ServiceType struct {
	A1Service A1Service
	TypeID    string
}

// For subscription manager

type SubscriptionKey struct {
	TargetXAppID topoapi.ID
}

type SubscriptionValue struct {
	A1EndpointIP          string
	A1EndpointPort        uint32
	A1ServiceCapabilities []*A1ServiceType
}

// For A1-PM/EI - A1Type/Obj - xApp mapping

type A1Key struct {
	TargetXAppID topoapi.ID
}

type A1PolicyObjectID string

type A1EIJobObjectID string

type A1PMValue struct {
	A1PolicyObjects map[A1PolicyObjectID]A1ServiceType
}

type A1EIValue struct {
	A1EIJobObjects map[A1EIJobObjectID]A1ServiceType
}
