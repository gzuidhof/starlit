/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */
package config

import (
	"fmt"

	"github.com/gzuidhof/starlit/starlit/web"
	"github.com/spf13/viper"
)


var defaultSandboxUrl = fmt.Sprintf("/static/vendor/%s/dist/index.html", web.GetVendoredPackage("starboard-notebook"))

func setDefaults() {
	viper.SetDefault("title", "Starlit")

	viper.SetDefault("appbar.enabled", true)
	viper.SetDefault("sidebar.enabled", true)
	viper.SetDefault("sandbox_url", defaultSandboxUrl)

	viper.SetDefault("reload", true)

	viper.SetDefault("output_folder", "build")
	viper.SetDefault("server.port", 8585)
	viper.SetDefault("server.secondary_port", 5742) //star in 1337speak

}
