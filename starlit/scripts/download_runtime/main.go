/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */
package main

import (
	"log"
	"os"
	"path"

	"github.com/gzuidhof/starlit/starlit/internal/npm"
)

const defaultRuntimePackageName = "starboard-notebook"
const defaultWrapPackageName = "starboard-wrap"

// Deletes the src and test folder in the output, saves some KB in executable size.
func deleteUselessFiles(fromFolder string, toRemove []string) {
	for _, folder := range toRemove {
		err := os.RemoveAll(path.Join(fromFolder, folder))
		if err != nil {
			log.Fatalf("Failed to delete: %v", err)
		}
	}
}

func dirExists(dirpath string) bool {
	info, err := os.Stat(dirpath)
	if os.IsNotExist(err) {
		return false
	}
	return info.IsDir()
}

func main() {
	if len(os.Args) < 4 {
		log.Print("Not enough arguments, supply 3 arguments: the package name, version and output folder")
		os.Exit(1)
	}
	packageName := os.Args[1]
	version := os.Args[2]
	outFolder := os.Args[3]

	if version != "latest" && dirExists(path.Join(outFolder, packageName+"@"+version)) {
		log.Printf("Skipping NPM fetch of %s as version %v already seems to be vendored already", packageName, version)
		os.Exit(0)
	}

	// id has the form <packagename>@<version>
	id, err := npm.DownloadPackageIntoFolder(packageName, version, outFolder)
	packageFolder := path.Join(outFolder, id)

	if packageName == defaultRuntimePackageName {
		deleteUselessFiles(packageFolder, []string{"dist/src", "dist/test"})
	} else if packageName == defaultWrapPackageName {
		// Possible improvement: create a "keep-these-files" function.. we really only need one file
		deleteUselessFiles(packageFolder, []string{
			".github",
			"dist/index.iife.js", // Some large markdown files we really don't have to statically bundle
			"dist/index.js",
			"dist/index.cjs",
			"dist/index.d.ts",
			"dist/embed.d.ts",
			"examples",
			"src",
			"README.md",
			"rollup.config.js",
			"tsconfig.json",
		})
	}

	if err != nil {
		log.Fatalf("Failed to fetch %s: %v", packageName, err)
	}
	log.Printf("Downloaded %s into %s", id, outFolder)
}
