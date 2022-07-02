// Copyright 2020-2022 Dave Shanley / Quobix
// SPDX-License-Identifier: MIT

package cmd

import (
	"github.com/daveshanley/vacuum/cui"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
	"os"
)

// TODO: This is a temporary UI, it's to help figure out the best experience, and it is not intended as a final face
// of vacuum. It's going to change around a good bit, so don't get too comfy with it :)

var (
	rootCmd = &cobra.Command{
		Use:   "vacuum lint <your-openapi-file.yaml>",
		Short: "vacuum is a very, very fast OpenAPI linter",
		Long:  `vacuum is a very, very fast OpenAPI linter. It will suck all the lint off your spec in milliseconds`,
		RunE: func(cmd *cobra.Command, args []string) error {

			pterm.Println()

			pterm.DefaultBigText.WithLetters(
				pterm.NewLettersFromStringWithRGB("vacuum", pterm.NewRGB(153, 51, 255))).
				Render()

			pterm.Println("To see something useful, try 'vacuum lint <my-openapi-spec.yaml>'")

			pterm.Println()

			return nil
		},
	}
)

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().BoolP("time", "t", false, "Show how long vacuum took to run")
	rootCmd.PersistentFlags().StringP("ruleset", "r", "", "Path to a spectral ruleset configuration")

	lintCommand := cui.GetLintCommand()

	rootCmd.AddCommand(lintCommand)
	rootCmd.AddCommand(GetVacuumReportCommand())
	rootCmd.AddCommand(GetSpectralReportCommand())
	rootCmd.AddCommand(GetHTMLReportCommand())
	rootCmd.AddCommand(GetDashboardCommand())
	rootCmd.AddCommand(GetGenerateRulesetCommand())

}

func initConfig() {

	// do something with this later, we don't need any configuration files right now

	//if cfgFile != "" {
	//	// Use config file from the flag.
	//	viper.SetConfigFile(cfgFile)
	//} else {
	//	// Find home directory.
	//	home, err := os.UserHomeDir()
	//	cobra.CheckErr(err)
	//
	//	// Search config in home directory with name ".cobra" (without extension).
	//	viper.AddConfigPath(home)
	//	viper.SetConfigType("yaml")
	//	viper.SetConfigName(".cobra")
	//}
	//
	//viper.AutomaticEnv()
	//
	//if err := viper.ReadInConfig(); err == nil {
	//	fmt.Println("Using config file:", viper.ConfigFileUsed())
	//}
}
