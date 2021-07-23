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
	Long:  `Test runs the notebooks in the given path start to end. A test fails if an error is thrown in any of the cells.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			log.Fatal("Please specify a path to test")
		}
		keepAlive, err := cmd.Flags().GetBool("serve")
		if err != nil {
			log.Fatal(err)
		}
		nbtest.TestPath(args[0], keepAlive)
	},
	PreRun: bindTestCmdToViper,
}

func init() {
	rootCmd.AddCommand(nbtestCmd)

	nbtestCmd.Flags().Bool("serve", false, "Keep the NB test server alive so you can access the files on the `/nbtest/<filepath>` endpoint")

	nbtestCmd.Flags().String("static_folder", "", "Override where static assets are loaded from, it uses the embedded assets if not set")
	nbtestCmd.Flags().String("templates_folder", "", "Override where templates are loaded from, it uses the embedded assets if not set")
}


func bindTestCmdToViper(cmd *cobra.Command, args []string) {
	viper.BindPFlag("static_folder", cmd.Flags().Lookup("static_folder"))
	viper.BindPFlag("templates_folder", cmd.Flags().Lookup("templates_folder"))
}