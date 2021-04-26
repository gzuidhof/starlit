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
	"github.com/gofiber/template/jet"
	"github.com/gzuidhof/starlit/starlit/internal/fs/assetfs"
	"github.com/gzuidhof/starlit/starlit/internal/fs/stripprefix"
	"github.com/gzuidhof/starlit/starlit/internal/handler/book"
	"github.com/spf13/afero"
	"github.com/spf13/viper"
)

func startPrimaryServer(serveFolderAbs string, serveFS assetfs.ServeFS, port string) error {
	engine := jet.NewFileSystem(afero.NewHttpFs(serveFS.Templates), ".jet.html")

	engine.Debug(true)
	engine.Reload(true)
	err := engine.Load()
	if err != nil {
		log.Fatalf("Could not parse templates")
	}

	engine.AddFunc("appbar", fiber.Map{
		"enabled": viper.GetBool("appbar.enabled"),
	})

	// fmt.Printf("%s", serveFS.Templates)
	app := fiber.New(fiber.Config{CaseSensitive: true, DisableStartupMessage: true, Views: engine})

	app.Use(recover.New())
	app.Use(logger.New())

	b := &book.BookHandler{}

	app.Use("/static/*", filesystem.New(filesystem.Config{
		Root: afero.NewHttpFs(stripprefix.New("/static/", serveFS.Static)),
	}))

	app.Get("/*", b.Handle)

	return app.Listen(":" + port)
}
