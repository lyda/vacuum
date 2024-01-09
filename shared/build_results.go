package shared

import (
	"os"
	"time"

	"github.com/daveshanley/vacuum/model"
	"github.com/daveshanley/vacuum/motor"
	"github.com/daveshanley/vacuum/rulesets"
	"github.com/pterm/pterm"
)

// BuildResults get results with no skip checks.
func BuildResults(
	silent bool,
	hardMode bool,
	rulesetFlag string,
	specBytes []byte,
	customFunctions map[string]model.RuleFunction,
	base string,
	timeout time.Duration) (*model.RuleResultSet, *motor.RuleSetExecutionResult, error) {
	return BuildResultsWithDocCheckSkip(silent, hardMode, rulesetFlag, specBytes, customFunctions, base, false, timeout)
}

// BuildResultsWithDocCheckSkip get results with configurable skip checks.
func BuildResultsWithDocCheckSkip(
	silent bool,
	hardMode bool,
	rulesetFlag string,
	specBytes []byte,
	customFunctions map[string]model.RuleFunction,
	base string,
	skipCheck bool,
	timeout time.Duration) (*model.RuleResultSet, *motor.RuleSetExecutionResult, error) {

	// read spec and parse
	defaultRuleSets := rulesets.BuildDefaultRuleSets()

	// default is recommended rules, based on spectral (for now anyway)
	selectedRS := defaultRuleSets.GenerateOpenAPIRecommendedRuleSet()

	// HARD MODE
	if hardMode {
		selectedRS = defaultRuleSets.GenerateOpenAPIDefaultRuleSet()

		// extract all OWASP Rules
		owaspRules := rulesets.GetAllOWASPRules()
		allRules := selectedRS.Rules
		for k, v := range owaspRules {
			allRules[k] = v
		}
		if !silent {
			box := pterm.DefaultBox.WithLeftPadding(5).WithRightPadding(5)
			box.BoxStyle = pterm.NewStyle(pterm.FgLightRed)
			box.Println(pterm.LightRed("🚨 HARD MODE ENABLED 🚨"))
			pterm.Println()
		}
	}

	// if ruleset has been supplied, lets make sure it exists, then load it in
	// and see if it's valid. If so - let's go!
	if rulesetFlag != "" {

		rsBytes, rsErr := os.ReadFile(rulesetFlag)
		if rsErr != nil {
			return nil, nil, rsErr
		}
		selectedRS, rsErr = BuildRuleSetFromUserSuppliedSet(rsBytes, defaultRuleSets)
		if rsErr != nil {
			return nil, nil, rsErr
		}
	}

	pterm.Info.Printf("Linting against %d rules: %s\n", len(selectedRS.Rules), selectedRS.DocumentationURI)

	ruleset := motor.ApplyRulesToRuleSet(&motor.RuleSetExecution{
		RuleSet:           selectedRS,
		Spec:              specBytes,
		CustomFunctions:   customFunctions,
		Base:              base,
		SkipDocumentCheck: skipCheck,
		AllowLookup:       true,
		Timeout:           timeout,
	})

	resultSet := model.NewRuleResultSet(ruleset.Results)
	resultSet.SortResultsByLineNumber()
	return resultSet, ruleset, nil
}
