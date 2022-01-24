// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-Only-1.0

package stream

import "fmt"

func GetEndpointIDWithTargetXAppID(targetXAppID string, a1Service A1Service) EndpointID {
	return EndpointID(fmt.Sprintf("%s-%s", targetXAppID, a1Service.String()))
}

func GetEndpointIDWithElement(elementName string) EndpointID {
	return EndpointID(elementName)
}
