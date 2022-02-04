// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0

package utils

var (
	NotifyURL             = "127.0.0.1:9640/"
	PolicyId              = "000"
	PolicyTypeId          = "111"
	ExpectedA1PMPolicyIDs = []string{PolicyId}
	ExpectedA1PMTypeIDs   = []string{PolicyTypeId}
	ExpectedPolicyObject  = `
{
	"scope": {
	"ueId": "855",
	"qosId": 67
	},
	"qosObjectives": {
	"priorityLevel": 50
	},
	"tspResources": [
	{
	"cellIdList": [
	39,
	40
	],
	"preference": "PREFER"
	},
	{
	"cellIdList": [
	81,
	82,
	83
	],
	"preference": "AVOID"
	}
	]
}`
	ExpectedPolicySchema = `
{
	"$schema": "http://json-schema.org/draft-07/schema#",
	"description": "O-RAN standard QoS Target policy",
	"type": "object",
	"properties": {
	"scope": {
	"anyOf": [
	{
	"type": "object",
	"properties": {
	"ueId": {"$ref": '#/definitions/UeId'},
	"groupId": {"$ref": '#/definitions/GroupId'},
	"qosId": {"$ref": '#/definitions/QosId'},
	"cellId": {"$ref": '#/definitions/CellId'}
	},
	"additionalProperties": false,
	"required": ["ueId", "qosId"]
	},
	{
	"type": "object",
	"properties": {
	"ueId": {"$ref": '#/definitions/UeId'},
	"sliceId": {"$ref": '#/definitions/SliceId'},
	"qosId": {"$ref": '#/definitions/QosId'},
	"cellId": {"$ref": '#/definitions/CellId'}
	},
	"additionalProperties": false,
	"required": ["ueId", "qosId"]
	},
	{
	"type": "object",
	"properties": {
	"groupId": {"$ref": '#/definitions/GroupId'},
	"qosId": {"$ref": '#/definitions/QosId'},
	"cellId": {"$ref": '#/definitions/CellId'}
	},
	"additionalProperties": false,
	"required": ["groupId", "qosId"]
	},
	{
	"type": "object",
	"properties": {
	"sliceId": {"$ref": '#/definitions/SliceId'},
	"qosId": {"$ref": '#/definitions/QosId'},
	"cellId": {"$ref": '#/definitions/CellId'}
	},
	"additionalProperties": false,
	"required": ["sliceId", "qosId"]
	},
	{
	"type": "object",
	"properties": {
	"qosId": {"$ref": '#/definitions/QosId'},
	"cellId": {"$ref": '#/definitions/CellId'}
	},
	"additionalProperties": false,
	"required": ["qosId"]
	}
	]
	},
	"qosObjectives": {
	"type": "object",
	"properties": {
	"gfbr": {"type": "number"},
	"mfbr": {"type": "number"},
	"priorityLevel": {"type": "number"},
	"pdb": {"type": "number"}
	},
	"minProperties": 1,
	"additionalProperties": false
	}
	},
	"additionalProperties": false,
	"required": ["scope", "qosObjectives"],
	"definitions": {
	"UeId": {"type": "string"},
	"GroupId": {"type": "number"},
	"SliceId": {"type": "number"},
	"QosId": {"type": "number"},
	"CellId": {"type": "number"}
	}
}`

	ExpectedPolicyStatusSchema = `
{
	"$schema": "http://json-schema.org/draft-07/schema#",
	"description": "O-RAN standard policy status",
	"type": "object",
	"properties": {
	"enforceStatus": {
	"type": "string",
	"enum": [
	"ENFORCED",
	"NOT_ENFORCED"
	]
	},
	"enforceReason": {
	"type": "string",
	"enum": [
	"SCOPE_NOT_APPLICABLE",
	"STATEMENT_NOT_APPLICABLE",
	"OTHER_REASON"
	]
	}
	},
	"additionalProperties": false,
	"required": ["enforceStatus"]
}`
)
