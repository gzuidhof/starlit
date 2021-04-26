/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package cmd

import (
	"fmt"
	"os"

	"github.com/gzuidhof/starlit/starlit/internal/config"
	"github.com/spf13/cobra"
)

var cfgFile string

// Version (injected by goreleaser)
var version = "<unknown>"
var date = "date unknown"
var commit = ""
var target = ""

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "starlit",
	Short: "Starlit is a tool for collaborating on notebooks and other content.",
	Long:  `Starlit is a tool that takes Starboard Notebooks and transforms them into a website.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//	Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(func() { config.InitConfig(cfgFile) })

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.starlit.yaml)")

	rootCmd.Version = version + " " + target + " (" + date + ") " + commit

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	// rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
