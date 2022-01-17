package utils

import (
	"encoding/json"
	"strings"

	"github.com/santhosh-tekuri/jsonschema/v5"
)

func ValidateSchema(object, schemaObject map[string]interface{}) error {

	ps, err := json.Marshal(schemaObject)
	if err != nil {
		return err
	}

	compiler := jsonschema.NewCompiler()
	if err := compiler.AddResource("schema.json", strings.NewReader(string(ps))); err != nil {
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
