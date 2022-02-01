// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-Only-1.0

package utils

import (
	"encoding/json"
	policyschemas "github.com/onosproject/onos-a1-dm/go/policy_schemas"
	"github.com/onosproject/onos-lib-go/pkg/errors"
	"github.com/onosproject/onos-lib-go/pkg/logging"
	"github.com/xeipuuv/gojsonschema"
)

var jsonValLog = logging.GetLogger("utils", "json-validator")

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
	json.Unmarshal([]byte(doc), &result)
	return result
}

func GetPolicyObject(query map[string]interface{}) map[string]interface{} {
	if _, ok := query["policyId"]; ok {
		delete(query, "policyId")
	}

	if _, ok := query["policyTypeId"]; ok {
		delete(query, "policyTypeId")
	}

	return query
}

func PolicyObjListValidate(objs interface{}) (bool, error) {
	switch objs.(type) {
	case []map[string]interface{}:
		if len(objs.([]map[string]interface{})) == 0 {
			return false, errors.NewNotFound("there is no policy object")
		}

		targetDoc, err := json.Marshal(objs.([]map[string]interface{})[0])
		if err != nil {
			return false, err
		}

		for i := 1; i < len(objs.([]map[string]interface{})); i++ {
			tmpDoc, err := json.Marshal(objs.([]map[string]interface{})[0])
			if err != nil {
				return false, err
			}

			if string(targetDoc) != string(tmpDoc) {
				return false, errors.NewConflict("PolicyObject is inconsistent")
			}
		}

		return true, nil
	case [][]string:
		if len(objs.([][]string)) == 0 {
			return false, errors.NewNotFound("there is no policy object")
		}

		targetDoc, err := json.Marshal(objs.([][]string)[0])
		if err != nil {
			return false, err
		}

		for i := 1; i < len(objs.([][]string)); i++ {
			tmpDoc, err := json.Marshal(objs.([][]string)[0])
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
