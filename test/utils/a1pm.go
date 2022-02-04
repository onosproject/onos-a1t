// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0

package utils

import (
	"strings"

	"github.com/santhosh-tekuri/jsonschema/v5"
)

func ValidateSchema(object, schemaObject string) error {
	compiler := jsonschema.NewCompiler()
	if err := compiler.AddResource("schema.json", strings.NewReader(schemaObject)); err != nil {
		return err
	}
	schema, err := compiler.Compile("schema.json")
	if err != nil {
		return err
	}
	if err := schema.Validate(object); err != nil {
		return err
	}

	return nil
}
