// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0

package utils

import (
	"encoding/json"

	"github.com/santhosh-tekuri/jsonschema/v5"
)

func ValidateSchema(instance, schema string) error {

	sch, err := jsonschema.CompileString("schema.json", schema)
	if err != nil {
		return err
	}

	var v interface{}
	if err := json.Unmarshal([]byte(instance), &v); err != nil {
		return err
	}

	if err = sch.Validate(v); err != nil {
		return err
	}

	return nil
}
