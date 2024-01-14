// Package cmd implements the command line parsing.
//
// Copyright 2020-2022 Dave Shanley / Quobix
// SPDX-License-Identifier: MIT
package cmd

import (
	"errors"

	"github.com/daveshanley/vacuum/cui"
	"github.com/daveshanley/vacuum/shared"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

// GetDashboardCommand gets the cobra.Command instance for the dashboard command.
func GetDashboardCommand() *cobra.Command {

	cmd := &cobra.Command{
		Use:     "dashboard",
		Short:   "Show vacuum dashboard for linting report",
		Long:    "Interactive console dashboard to explore linting report in detail",
		Example: "vacuum dashboard my-awesome-spec.yaml",
		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			if len(args) != 0 {
				return nil, cobra.ShellCompDirectiveNoFileComp
			}
			return []string{"yaml", "yml", "json"}, cobra.ShellCompDirectiveFilterFileExt
		},
		RunE: func(cmd *cobra.Command, args []string) error {

			// check for file args
			if len(args) == 0 {
				errText := "please supply an OpenAPI specification to generate a report"
				pterm.Error.Println(errText)
				pterm.Println()
				return errors.New(errText)
			}
			baseFlag, _ := cmd.Flags().GetString("base")
			skipCheckFlag, _ := cmd.Flags().GetBool("skip-check")
			timeoutFlag, _ := cmd.Flags().GetInt("timeout")
			hardModeFlag, _ := cmd.Flags().GetBool("hard-mode")
			followFlag, _ := cmd.Flags().GetBool("follow")
			functionsFlag, _ := cmd.Flags().GetString("functions")
			rulesetFlag, _ := cmd.Flags().GetString("ruleset")
			customFunctions, _ := shared.LoadCustomFunctions(functionsFlag)

			openapiFile, err := cui.NewFile(
				args[0],
				baseFlag,
				skipCheckFlag,
				timeoutFlag,
				hardModeFlag,
				followFlag,
				customFunctions,
				rulesetFlag,
			)
			cobra.CheckErr(err)
			cobra.CheckErr(openapiFile.ReadFile())

			dash := cui.CreateDashboard(openapiFile)
			dash.Version = Version
			return dash.Render()
		},
	}

	cmd.Flags().BoolP("follow", "F", false, "Follow any changes to the OpenAPI spec.")

	return cmd
}
