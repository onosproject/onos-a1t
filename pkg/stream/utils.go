// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0

package stream

import "fmt"

func GetEndpointIDWithTargetXAppID(targetXAppID string, a1Service A1Service) EndpointID {
	return EndpointID(fmt.Sprintf("%s-%s", targetXAppID, a1Service.String()))
}

func GetEndpointIDWithElement(elementName string) EndpointID {
	return EndpointID(elementName)
}

func NewSBStreamMessage(targetXAppID string, messageType A1SBIMessageType, sbirpcType A1SBIRPCType, service A1Service, payload interface{}) *SBStreamMessage {
	return &SBStreamMessage{
		TargetXAppID:     targetXAppID,
		A1SBIMessageType: messageType,
		A1SBIRPCType:     sbirpcType,
		A1Service:        service,
		Payload:          payload,
	}
}

// GetStreamID returns 1. southbound stream ID (sbID) and 2. northbound stream ID (nbID)
func GetStreamID(nbControllerID EndpointID, sbControllerID EndpointID) (ID, ID) {
	return ID{
			SrcEndpointID:  nbControllerID,
			DestEndpointID: sbControllerID,
		}, ID{
			SrcEndpointID:  sbControllerID,
			DestEndpointID: nbControllerID,
		}
}
