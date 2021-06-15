/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */
package server

import (
	"log"
	"sort"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gzuidhof/starlit/starlit/internal/fs/assetfs"
	"github.com/gzuidhof/starlit/starlit/internal/fs/stripprefix"
	"github.com/gzuidhof/starlit/starlit/internal/spaces/pages"
	"github.com/spf13/afero"
	"github.com/spf13/viper"
)

func CreatePrimaryApp(serveFolderAbs string, serveFS assetfs.ServeFS) (*fiber.App, error) {
	app := fiber.New(fiber.Config{CaseSensitive: true, DisableStartupMessage: true})

	app.Use(recover.New())
	app.Use(logger.New())

	app.Use("/static/*", filesystem.New(filesystem.Config{
		Root: afero.NewHttpFs(stripprefix.New("/static/", serveFS.Static)),
	}))

	spaces := viper.GetStringMap("spaces")
	spaceNames := make([]string, len(spaces))
	i := 0;
	for k := range spaces {
		spaceNames[i] = k
		i++;
	}
	sort.Strings(spaceNames)

	spaceNamesReverseOrder := append([]string(nil), spaceNames...)
	// TODO: this is such a hack. We should order spaces by their weight.
	sort.Sort(sort.Reverse(sort.StringSlice(spaceNamesReverseOrder)))
	viper.Set("space_names", spaceNamesReverseOrder)

	for _, appName := range spaceNames {
		mountPath := viper.GetString("spaces." + appName + ".path")
		if (mountPath == "") {
			mountPath = "/" + appName
			viper.Set("spaces." + appName + ".path", mountPath)
			log.Printf("No path set for space %s, defaulting to \"%s\"", appName, mountPath)
		}
		subApp := pages.CreateApp(appName, serveFS, true)

		log.Printf("Mounting app %s on %s\n", appName, mountPath)
		app.Mount(mountPath, subApp)
	}

	return app, nil
}

func startPrimaryServer(serveFolderAbs string, serveFS assetfs.ServeFS, port string) error {
	app, err := CreatePrimaryApp(serveFolderAbs, serveFS)

	if err != nil {
		log.Fatal(err)
	}

	return app.Listen(":" + port)
}
