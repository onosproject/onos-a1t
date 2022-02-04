// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0

package utils

var (
	NotifyURL             = "http://127.0.0.1:9640/"
	PolicyId              = "000"
	PolicyTypeId          = "ORAN_TrafficSteeringPreference_2.0.0"
	ExpectedA1PMPolicyIDs = []string{"1", "2", PolicyId}
	ExpectedA1PMTypeIDs   = []string{PolicyTypeId}
	ExpectedPolicyObject  = `
	{
	  "scope": {
		"ueId": "0000000000000002"
	  },
	  "tspResources": [
		{
		  "cellIdList": [
			{"plmnId": {"mcc": "248","mnc": "35"},
			  "cId": {"ncI": 39}},
			{"plmnId": {"mcc": "248","mnc": "35"},
			 "cId": {"ncI": 40}}
		  ], 
		  "preference": "PREFER"
		},
		{
		  "cellIdList": [
			{"plmnId": {"mcc": "248","mnc": "35"},
			  "cId": {"ncI": 81}},
			{"plmnId": {"mcc": "248","mnc": "35"},
			  "cId": {"ncI": 82}},
			{"plmnId": {"mcc": "248","mnc": "35"},
			 "cId": {"ncI": 83}}
		  ],
		  "preference": "FORBID"
		}
	  ]
	}
	`
	ExpectedPolicySchema = `
	{
		"$schema": "http://json-schema.org/draft-07/schema#",
		"description": "O-RAN standard Traffic Steering Preference policy",
		"type": "object",
		"properties": {
		  "scope": {
			"anyOf": [
			  {
				"type": "object",
				"properties": {
				  "ueId": {
					"$ref": "#/definitions/UeId"
				  },
				  "sliceId": {
					"$ref": "#/definitions/SliceId"
				  },
				  "qosId": {
					"$ref": "#/definitions/QosId"
				  },
				  "cellId": {
					"$ref": "#/definitions/CellId"
				  }
				},
				"additionalProperties": false,
				"required": [
				  "ueId"
				]
			  },
			  {
				"type": "object",
				"properties": {
				  "sliceId": {
					"$ref": "#/definitions/SliceId"
				  },
				  "qosId": {
					"$ref": "#/definitions/QosId"
				  },
				  "cellId": {
					"$ref": "#/definitions/CellId"
				  }
				},
				"additionalProperties": false,
				"required": [
				  "sliceId"
				]
			  }
			]
		  },
		  "tspResources": {
			"type": "array",
			"items": {
			  "$ref": "#/definitions/TspResource"
			},
			"minItems": 1
		  }
		},
		"additionalProperties": false,
		"required": [
		  "scope",
		  "tspResources"
		],
		"definitions": {
		  "UeId": {
			"type": "string",
			"pattern": "^[A-Fa-f0-9]{16}$"
		  },
		  "SliceId": {
			"type": "object",
			"properties": {
			  "sst": {
				"type": "integer",
				"minimum": 0,
				"maximum": 255
			  },
			  "sd": {
				"type": "string",
				"pattern": "^[A-Fa-f0-9]{6}$"
			  },
			  "plmnId": {
				"$ref": "#/definitions/PlmnId"
			  }
			},
			"additionalProperties": false,
			"required": [
			  "sst",
			  "plmnId"
			]
		  },
		  "QosId": {
			"oneOf": [
			  {
				"type": "object",
				"properties": {
				  "5qI": {
					"type": "integer",
					"minimum": 1,
					"maximum": 256
				  }
				},
				"additionalProperties": false,
				"required": [
				  "5qI"
				]
			  },
			  {
				"type": "object",
				"properties": {
				  "qcI": {
					"type": "integer",
					"minimum": 1,
					"maximum": 256
				  }
				},
				"additionalProperties": false,
				"required": [
				  "qcI"
				]
			  }
			]
		  },
		  "CellId": {
			"type": "object",
			"properties": {
			  "plmnId": {
				"$ref": "#/definitions/PlmnId"
			  },
			  "cId": {
				"$ref": "#/definitions/CId"
			  }
			},
			"additionalProperties": false,
			"required": [
			  "plmnId",
			  "cId"
			]
		  },
		  "CId": {
			"oneOf": [
			  {
				"type": "object",
				"properties": {
				  "ncI": {
					"$ref": "#/definitions/NcI"
				  }
				},
				"additionalProperties": false,
				"required": [
				  "ncI"
				]
			  },
			  {
				"type": "object",
				"properties": {
				  "ecI": {
					"$ref": "#/definitions/EcI"
				  }
				},
				"additionalProperties": false,
				"required": [
				  "ecI"
				]
			  }
			]
		  },
		  "NcI": {
			"type": "integer",
			"minimum": 0,
			"maximum": 68719476735
		  },
		  "EcI": {
			"type": "integer",
			"minimum": 0,
			"maximum": 268435455
		  },
		  "PlmnId": {
			"type": "object",
			"properties": {
			  "mcc": {
				"type": "string",
				"pattern": "^[0-9]{3}$"
			  },
			  "mnc": {
				"type": "string",
				"pattern": "^[0-9]{2,3}$"
			  }
			},
			"additionalProperties": false,
			"required": [
			  "mcc",
			  "mnc"
			]
		  },
		  "PreferenceType": {
			"type": "string",
			"enum": [
			  "SHALL",
			  "PREFER",
			  "AVOID",
			  "FORBID"
			]
		  },
		  "CellIdList": {
			"type": "array",
			"items": {
			  "$ref": "#/definitions/CellId"
			}
		  },
		  "TspResource": {
			"type": "object",
			"properties": {
			  "cellIdList": {
				"$ref": "#/definitions/CellIdList"
			  },
			  "preference": {
				"$ref": "#/definitions/PreferenceType"
			  },
			  "primary": {
				"type": "boolean"
			  }
			},
			"required": [
			  "cellIdList",
			  "preference"
			],
			"additionalProperties": false
		  }
		}
	  }
	`

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
		"required": [
		  "enforceStatus"
		]
	  }
	`
)
