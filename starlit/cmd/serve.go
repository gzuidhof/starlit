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

var serverCmd = &cobra.Command{
	Use:   "serve",
	Short: "Serves files in the given folder or path.",
	Long:  `Serves files in the given folder or path.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			log.Fatal("Please specify a path to serve")
		}
		server.Start(args[0])
	},
	PreRun: bindServeCmdToViper,
}

func init() {
	rootCmd.AddCommand(serverCmd)

	serverCmd.Flags().StringP("port", "p", "8585", "Port to serve files on")
	serverCmd.Flags().String("port_secondary", "5742", "Port used as secondary origin (for sandboxing)")

	serverCmd.Flags().String("static_folder", "", "Override where static assets are served from, it uses the embedded assets if not set")
	serverCmd.Flags().String("templates_folder", "", "Override where templates are loaded from, it uses the embedded assets if not set")
}

func bindServeCmdToViper(cmd *cobra.Command, args []string) {
	viper.BindPFlag("server.port", cmd.Flags().Lookup("port"))
	viper.BindPFlag("server.port_secondary", cmd.Flags().Lookup("port_secondary"))

	viper.BindPFlag("static_folder", cmd.Flags().Lookup("static_folder"))
	viper.BindPFlag("templates_folder", cmd.Flags().Lookup("templates_folder"))
}