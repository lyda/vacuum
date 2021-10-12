package model

import (
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/daveshanley/vaccum/utils"
	"github.com/xeipuuv/gojsonschema"
	"gopkg.in/yaml.v3"
	"strings"
)

//go:embed schemas/ruleset.schema.json
var rulesetSchema string

type RuleFunctionContext struct {
	RuleAction *RuleAction
	Options    interface{}
}

type RuleFunctionResult struct {
	Message string
	Path    string
}

type RuleFunction interface {
	RunRule(nodes []*yaml.Node, context RuleFunctionContext) []RuleFunctionResult
}

type RuleAction struct {
	Field           string      `json:"field"`
	FunctionName    string      `json:"function"`
	FunctionOptions interface{} `json:"functionOptions"`
}

type Rule struct {
	Description string      `json:"description"`
	Given       string      `json:"given"`
	Formats     []string    `json:"formats"`
	Resolved    bool        `json:"resolved"`
	Recommended bool        `json:"recommended"`
	Severity    string      `json:"severity"`
	Then        *RuleAction `json:"then"`
}

func (r Rule) ToJSON() string {
	d, _ := json.Marshal(r)
	return string(d)
}

type RuleSet struct {
	DocumentationURI string           `json:"documentationUrl"`
	Formats          []string         `json:"formats"`
	Rules            map[string]*Rule `json:"rules"`
	schemaLoader     gojsonschema.JSONLoader
}

func CreateRuleSetUsingJSON(jsonData []byte) (*RuleSet, error) {
	jsonString := string(jsonData)
	if !utils.IsJSON(jsonString) {
		return nil, errors.New("data is not JSON")
	}

	jsonLoader := gojsonschema.NewStringLoader(jsonString)
	schemaLoader := LoadRulesetSchema()

	// check blob is a valid contract, before creating ruleset.
	res, err := gojsonschema.Validate(schemaLoader, jsonLoader)
	if err != nil {
		return nil, err
	}

	if !res.Valid() {
		var buf strings.Builder
		for _, e := range res.Errors() {
			buf.WriteString(fmt.Sprintf("%s (line),", e.Description()))
		}

		return nil, errors.New(fmt.Sprintf("rules not valid: %s", buf.String()))
	}

	// unmarshal JSON into new RuleSet
	rs := &RuleSet{}
	err = json.Unmarshal(jsonData, rs)
	if err != nil {
		return nil, err
	}

	// save our loaded schema for later.
	rs.schemaLoader = schemaLoader
	return rs, nil
}

func LoadRulesetSchema() gojsonschema.JSONLoader {
	return gojsonschema.NewStringLoader(rulesetSchema)
}
