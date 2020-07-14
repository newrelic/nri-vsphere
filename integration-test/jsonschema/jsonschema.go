/*
* Copyright 2020 New Relic Corporation. All rights reserved.
* SPDX-License-Identifier: Apache-2.0
 */

package jsonschema

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/xeipuuv/gojsonschema"
)

// Validate validates the input argument against JSON schema. If the
// input is not valid the error is returned. The first argument is the file name
// of the JSON schema. It is used to build file URI required to load the JSON schema.
// The second argument is the input string that is validated.
func Validate(fileName string, input string) error {
	pwd, err := os.Getwd()
	if err != nil {
		return err
	}
	schemaURI := fmt.Sprintf("file://%s", filepath.ToSlash(filepath.Join(pwd, fileName)))
	schemaLoader := gojsonschema.NewReferenceLoader(schemaURI)
	documentLoader := gojsonschema.NewStringLoader(input)

	schema, err := gojsonschema.NewSchema(schemaLoader)
	if err != nil {
		return fmt.Errorf("error loading JSON schema, error: %v, schemaURI: %s", err, schemaURI)
	}

	result, err := schema.Validate(documentLoader)
	if err != nil {
		return fmt.Errorf("error validating, error: %v", err)
	}

	if result.Valid() {
		return nil
	}
	fmt.Printf("Errors for JSON schema: '%s'\n", schemaURI)
	for _, desc := range result.Errors() {
		fmt.Printf("\t- %s\n", desc)
	}
	fmt.Printf("\n")
	return fmt.Errorf("the output of the integration doesn't have expected JSON format")
}

// ValidationField is a struct used in JSON schema
type ValidationField struct {
	Keyword      string
	KeywordValue interface{}
}
