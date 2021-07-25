/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package server

import (
	"fmt"
	"log"
	"path/filepath"

	"github.com/gzuidhof/starlit/starlit/internal/fs/assetfs"
	"github.com/spf13/viper"
)

func Start(servePath string) {
	portPrimary := viper.GetString("server.port")
	portSecondary := viper.GetString("server.port_secondary")
	
	serveFolder := servePath
	serveFolderAbs, err := filepath.Abs(serveFolder)
	if err != nil {
		log.Fatalf("Invalid serve folder, could not get absolute path: %v", err)
	}

	viper.Set("serve_filepath", serveFolderAbs)
	serveFS := assetfs.GetAssetFileSystems()

	done := make(chan bool)
	go func() {
		log.Fatal(startPrimaryServer(serveFolderAbs, serveFS, portPrimary))
	}()
	go func() {
		log.Fatal(startSecondaryServer(serveFolderAbs, serveFS, portSecondary))
	}()
	log.Printf("\nListening on :%v (and :%s for sandboxing)\nhttp://localhost:%v", portPrimary, portSecondary, portPrimary)

	<-done
}


func StartNBTestServer(testPath string) string {
	portPrimary := viper.GetString("server.port")

	viper.Set("serve_filepath", testPath)
	serveFS := assetfs.GetAssetFileSystems()

	go func() {
		log.Fatal(startNBTestServer(testPath, serveFS, portPrimary))
	}()

	return fmt.Sprintf("http://localhost:%s", portPrimary)
}
