// Copyright 2022 Dave Shanley / Quobix
// SPDX-License-Identifier: MIT

package openapi

import (
	"fmt"
	"github.com/daveshanley/vacuum/model"
	"github.com/pb33f/doctor/model/high/base"
	"github.com/pb33f/libopenapi/utils"
	"gopkg.in/yaml.v3"
)

// DuplicatedEnum will check enum values match the types provided
type DuplicatedEnum struct {
}

// GetSchema returns a model.RuleFunctionSchema defining the schema of the DuplicatedEnum rule.
func (de DuplicatedEnum) GetSchema() model.RuleFunctionSchema {
	return model.RuleFunctionSchema{
		Name: "duplicated_enum",
	}
}

// RunRule will execute the DuplicatedEnum rule, based on supplied context and a supplied []*yaml.Node slice.
func (de DuplicatedEnum) RunRule(_ []*yaml.Node, context model.RuleFunctionContext) []model.RuleFunctionResult {

	if context.DrDocument == nil {
		return nil
	}

	var results []model.RuleFunctionResult

	schemas := context.DrDocument.Schemas

	for _, schema := range schemas {

		if schema.Value.Enum != nil {
			node := schema.Value.GoLow().Enum.KeyNode

			duplicates := utils.CheckEnumForDuplicates(schema.Value.Enum)

			// iterate through duplicate results and add results.
			for _, res := range duplicates {
				result := model.RuleFunctionResult{
					Message:   fmt.Sprintf("enum contains a duplicate: `%s`", res.Value),
					StartNode: node,
					EndNode:   node,
					Path:      fmt.Sprintf("%s.%s", schema.GenerateJSONPath(), "enum"),
					Rule:      context.Rule,
				}
				schema.AddRuleFunctionResult(base.ConvertRuleResult(&result))
				results = append(results, result)
			}
		}
	}

	return results
}
