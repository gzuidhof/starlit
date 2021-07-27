/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package cmd

import (
	"log"

	"github.com/gzuidhof/starlit/starlit/internal/nbtest"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var nbtestCmd = &cobra.Command{
	Use:   "nbtest",
	Short: "Tests the notebook files in the given path.",
	Long:  `nbest runs the notebooks in the given path start to end in a headless browser over the chrome debugging protocol. A test fails if an error is thrown in any of the cells.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			log.Fatal("Please specify a path to test")
		}
		nbtest.TestPath(args[0])
	},
	PreRun: bindTestCmdToViper,
}

func init() {
	rootCmd.AddCommand(nbtestCmd)

	nbtestCmd.Flags().Bool("serve", false, "Keep the NB test server alive so you can access the files on the `/nbtest/<filepath>` endpoint")
	nbtestCmd.Flags().Bool("headless", true, "Run in headless mode (turn off using `headless=false`)")

	nbtestCmd.Flags().Float64("timeout", 30, "Timeout for individual tests in seconds (0 = no timeout)")
	nbtestCmd.Flags().Int("concurrency", 0, "Number of tests to run in parallel (0 = num CPUs)")
	nbtestCmd.Flags().String("browser_exec_path", "", "(Advanced) browser binary path to use for testing, should be Chromium-based (Chrome, Edge, Safari, Opera, ...)")

	nbtestCmd.Flags().String("starboard_artifacts", "", "URL prefix of starboard artifacts (perhaps from a dev server on your localhost during local development) or a path to a folder. Defaults to the static assets embedded in the binary or served from the `static_folder` flag")
	nbtestCmd.Flags().String("pyodide_artifacts", "", "URL prefix of pyodide artifacts (e.g. \"https://cdn.jsdelivr.net/pyodide/v0.17.0/full/\")")

	nbtestCmd.Flags().String("static_folder", "", "Override where static assets are loaded from, it uses the embedded assets if not set")
	nbtestCmd.Flags().String("templates_folder", "", "Override where templates are loaded from, it uses the embedded assets if not set")
}

func bindTestCmdToViper(cmd *cobra.Command, args []string) {	
	viper.BindPFlag("static_folder", cmd.Flags().Lookup("static_folder"))
	viper.BindPFlag("templates_folder", cmd.Flags().Lookup("templates_folder"))

	viper.BindPFlag("nbtest.serve", cmd.Flags().Lookup("serve"))
	viper.BindPFlag("nbtest.headless", cmd.Flags().Lookup("headless"))

	viper.BindPFlag("nbtest.timeout", cmd.Flags().Lookup("timeout"))
	viper.BindPFlag("nbtest.concurrency", cmd.Flags().Lookup("concurrency"))
	viper.BindPFlag("nbtest.browser_exec_path", cmd.Flags().Lookup("browser_exec_path"))

	viper.BindPFlag("nbtest.starboard_artifacts", cmd.Flags().Lookup("starboard_artifacts"))
	viper.BindPFlag("nbtest.pyodide_artifacts", cmd.Flags().Lookup("pyodide_artifacts"))
}