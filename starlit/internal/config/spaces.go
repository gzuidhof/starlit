/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */
package config

import (
	"path"

	"github.com/spf13/viper"
)


func GetSpaceConfig(spaceName string, parentConfig *viper.Viper) *viper.Viper {
	vconfig := viper.New()
	// Merge the global config
	vconfig.MergeConfigMap(parentConfig.AllSettings())

	// This is not the root app
	if spaceName != "" {
		spaceConfig := viper.Sub("spaces." + spaceName)
		if spaceConfig != nil {
			vconfig.MergeConfigMap(spaceConfig.AllSettings())
		}
		vconfig.Set("serve_filepath", path.Join(parentConfig.GetString("serve_filepath"), spaceConfig.GetString("path")))
	}

	vconfig.Set("space", vconfig.AllSettings())
	return vconfig
}

func GetPageConfig(pageName string, spaceConfig *viper.Viper, frontmatterConfig *viper.Viper) *viper.Viper {
	vconfig := viper.New()
	// Merge the space config
	vconfig.MergeConfigMap(spaceConfig.AllSettings())

	pageConfig := viper.Sub("pages." + pageName)
	if pageConfig != nil {
		vconfig.MergeConfigMap(pageConfig.AllSettings())
	}

	vconfig.MergeConfigMap(frontmatterConfig.AllSettings())
	return vconfig
}