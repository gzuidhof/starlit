/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */
package assetfs

import (
	"github.com/gzuidhof/starlit/starlit/internal/fs/stripprefix"
	"github.com/gzuidhof/starlit/starlit/web/static"
	"github.com/gzuidhof/starlit/starlit/web/templates"
	"github.com/spf13/afero"
	"github.com/spf13/viper"
)

type ServeFS struct {
	Static    afero.Fs
	Templates afero.Fs
}

func GetAssetFileSystems() ServeFS {
	var staticFS afero.Fs
	var templatesFS afero.Fs

	if viper.GetString("static_folder") != "" {
		staticFS = afero.NewReadOnlyFs(afero.NewBasePathFs(afero.NewOsFs(), viper.GetString("static_folder")))
	} else {
		staticFS = afero.NewReadOnlyFs(afero.FromIOFS{static.FS})
	}

	if viper.GetString("templates_folder") != "" {
		templatesFS = afero.NewReadOnlyFs(afero.NewBasePathFs(afero.NewOsFs(), viper.GetString("templates_folder")))
	} else {
		templatesFS = afero.NewReadOnlyFs(stripprefix.New("/", afero.FromIOFS{templates.FS}))
	}

	return ServeFS{
		Static:    staticFS,
		Templates: templatesFS,
	}
}
