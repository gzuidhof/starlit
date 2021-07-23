/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package cmd

import (
	"log"

	"github.com/gzuidhof/starlit/starlit/internal/server"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var renderCmd = &cobra.Command{
	Use:   "render",
	Short: "Renders files in the given folder or path.",
	Long:  `Renders files in the given folder or path.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			log.Fatal("Please specify a path to render")
		}
		server.Start(args[0])
	},
	PreRun: bindRenderCmdToViper,
}

func init() {
	rootCmd.AddCommand(renderCmd)

	renderCmd.Flags().String("static_folder", "", "Override where static assets are loaded from, it uses the embedded assets if not set")
	renderCmd.Flags().String("templates_folder", "", "Override where templates are loaded from, it uses the embedded assets if not set")
}


func bindRenderCmdToViper(cmd *cobra.Command, args []string) {
	viper.BindPFlag("static_folder", cmd.Flags().Lookup("static_folder"))
	viper.BindPFlag("templates_folder", cmd.Flags().Lookup("templates_folder"))
}
