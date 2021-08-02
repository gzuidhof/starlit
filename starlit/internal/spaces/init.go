/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package spaces

import (
	"log"
	"sort"

	"github.com/gofiber/fiber/v2"
	"github.com/gzuidhof/starlit/starlit/internal/fs/assetfs"
	"github.com/gzuidhof/starlit/starlit/internal/spaces/pages"
	"github.com/spf13/viper"
)

type Space struct {
	Name string
	MountPath string
	App *fiber.App
}

type Spaces []Space

func SetupSpaces(serveFS assetfs.ServeFS) Spaces {
	spacesConfigMap := viper.GetStringMap("spaces")
	spaceNames := make([]string, len(spacesConfigMap))
	i := 0;
	for k := range spacesConfigMap {
		spaceNames[i] = k
		i++;
	}
	sort.Strings(spaceNames)

	spaceNamesReverseOrder := append([]string(nil), spaceNames...)
	// TODO: this is such a hack. We should order spaces by their weight.
	sort.Sort(sort.Reverse(sort.StringSlice(spaceNamesReverseOrder)))
	viper.Set("space_names", spaceNamesReverseOrder)

	spaces := make(Spaces, len(spaceNames))

	for i, name := range spaceNames {
		mountPath := viper.GetString("spaces." + name + ".path")
		if (mountPath == "") {
			mountPath = "/" + name
			viper.Set("spaces." + name + ".path", mountPath)
			log.Printf("No path set for space %s, defaulting to \"%s\"", name, mountPath)
		}
		app := pages.CreateApp(name, serveFS, true)

		spaces[i] = Space{
			Name: name,
			App: app,
			MountPath: mountPath,
		}
	}

	return spaces
}

func (s Spaces) MountSpacesOnApp(app *fiber.App) {
	for _, space := range s {
		log.Printf("Mounting space %s on %s\n", space.Name, space.MountPath)
		app.Mount(space.MountPath, space.App)
	}
}
