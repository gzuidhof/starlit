/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */
package server

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gzuidhof/starlit/starlit/internal/fs/assetfs"
	"github.com/gzuidhof/starlit/starlit/internal/fs/stripprefix"
	"github.com/gzuidhof/starlit/starlit/internal/spaces"
	"github.com/spf13/afero"
)

func CreatePrimaryApp(serveFolderAbs string, serveFS assetfs.ServeFS) (*fiber.App, error) {
	app := fiber.New(fiber.Config{CaseSensitive: true, DisableStartupMessage: true})

	app.Use(recover.New())
	app.Use(logger.New())

	app.Use("/static/*", filesystem.New(filesystem.Config{
		Root: afero.NewHttpFs(stripprefix.New("/static/", serveFS.Static)),
	}))

	spaceApps := spaces.SetupSpaces(serveFS)
	spaceApps.MountSpacesOnApp(app)
	
	return app, nil
}

func startPrimaryServer(serveFolderAbs string, serveFS assetfs.ServeFS, port string) error {
	app, err := CreatePrimaryApp(serveFolderAbs, serveFS)

	if err != nil {
		log.Fatal(err)
	}

	return app.Listen(":" + port)
}
