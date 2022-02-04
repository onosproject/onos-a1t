// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0

package utils

import (
	"encoding/json"
	policyschemas "github.com/onosproject/onos-a1-dm/go/policy_schemas"
	"github.com/onosproject/onos-lib-go/pkg/errors"
	"github.com/onosproject/onos-lib-go/pkg/logging"
	"github.com/xeipuuv/gojsonschema"
)

var jsonValLog = logging.GetLogger("utils", "json-validator")

const NotificationDestination = "notificationDestination"

func JsonValidateWithTypeID(policyTypeID string, jsonDoc string) bool {
	schemeDoc, ok := policyschemas.PolicySchemas[policyTypeID]
	if !ok {
		jsonValLog.Errorf("Policy ID %v not supports", policyTypeID)
	}

	return JsonValidate(schemeDoc, jsonDoc)
}

func JsonValidate(schemaDoc string, jsonDoc string) bool {
	schemaLoader := gojsonschema.NewStringLoader(schemaDoc)
	documentLoader := gojsonschema.NewStringLoader(jsonDoc)
	result, err := gojsonschema.Validate(schemaLoader, documentLoader)
	if err != nil {
		jsonValLog.Error(err)
		return false
	}

	if result.Valid() {
		return true
	}

	for _, desc := range result.Errors() {
		jsonValLog.Error(desc)
	}
	return false
}

func ConvertStringFormatJsonToMap(doc string) map[string]interface{} {
	var result map[string]interface{}
	err := json.Unmarshal([]byte(doc), &result)
	if err != nil {
		jsonValLog.Error(err)
		return nil
	}
	return result
}

func GetPolicyObject(query map[string]interface{}) map[string]interface{} {
	result := make(map[string]interface{})

	for k, v := range query {
		if k != "policyId" && k != "policyTypeId" {
			result[k] = v
		}
	}

	return result
}

func PolicyObjListValidate(objs interface{}) (bool, error) {
	switch objs := objs.(type) {
	case []map[string]interface{}:
		if len(objs) == 0 {
			return false, errors.NewNotFound("there is no policy object")
		}

		targetDoc, err := json.Marshal(objs)
		if err != nil {
			return false, err
		}

		for i := 1; i < len(objs); i++ {
			tmpDoc, err := json.Marshal(objs[0])
			if err != nil {
				return false, err
			}

			if string(targetDoc) != string(tmpDoc) {
				return false, errors.NewConflict("PolicyObject is inconsistent")
			}
		}

		return true, nil
	case [][]string:
		if len(objs) == 0 {
			return false, errors.NewNotFound("there is no policy object")
		}

		targetDoc, err := json.Marshal(objs[0])
		if err != nil {
			return false, err
		}

		for i := 1; i < len(objs); i++ {
			tmpDoc, err := json.Marshal(objs[0])
			if err != nil {
				return false, err
			}

			if string(targetDoc) != string(tmpDoc) {
				return false, errors.NewConflict("PolicyObject is inconsistent")
			}
		}

		return true, nil
	}

	return false, errors.NewNotSupported("object type should be map[string]interface{} or []string")
}
