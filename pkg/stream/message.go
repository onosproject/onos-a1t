// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0

package stream

type A1SBIMessageType int

const (
	PolicyRequestMessage A1SBIMessageType = iota
	PolicyResultMessage
	PolicyStatusMessage
	PolicyAckMessage
	EIRequestMessage
	EIResultMessage
	EIStatusMessage
	EIAckMessage
)

func (a A1SBIMessageType) String() string {
	return [...]string{"PolicyRequestMessage", "PolicyResultMessage", "PolicyStatusMessage", "PolicyAckMessage",
		"EIRequestMessage", "EIResultMessage", "EIStatusMessage", "EIAckMessage"}[a]
}

type A1SBIRPCType int

const (
	PolicySetup A1SBIRPCType = iota
	PolicyUpdate
	PolicyDelete
	PolicyQuery
	PolicyStatus
	EIQuery
	EIJobSetup
	EIJobUpdate
	EIJobDelete
	EIJobStatusQuery
	EIJobStatusNotify
	EIJobResultDelivery
)

func (r A1SBIRPCType) String() string {
	return [...]string{"PolicySetup", "PolicyUpdate", "PolicyDelete", "PolicyQuery", "PolicyStatus",
		"EIQuery", "EIJobSetup", "EIJobUpdate", "EIJobDelete", "EIJobStatusQuery", "EIJobStatusNotify", "EIJobResultDelivery"}[r]
}

type A1Service int

const (
	PolicyManagement A1Service = iota
	EnrichmentInformation
)

func (s A1Service) String() string {
	return [...]string{"PolicyManagement", "EnrichmentInformation"}[s]
}

type SBStreamMessage struct {
	TargetXAppID     string
	A1SBIMessageType A1SBIMessageType
	A1SBIRPCType     A1SBIRPCType
	A1Service        A1Service
	Payload          interface{}
}

const (
	A1PController  = "a1p-controller"
	A1EIController = "a1ei-controller"
)
