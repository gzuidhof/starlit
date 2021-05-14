/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */
package pages

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/jet"
	"github.com/gzuidhof/starlit/starlit/internal/format"
	"github.com/gzuidhof/starlit/starlit/internal/fs/assetfs"
	"github.com/gzuidhof/starlit/starlit/internal/templaterenderer"
	"github.com/spf13/afero"
	"github.com/spf13/viper"
)

func CreateApp(name string, serveFS assetfs.ServeFS, reloadTemplates bool) *fiber.App {
	engine := jet.NewFileSystem(afero.NewHttpFs(serveFS.Templates), ".jet.html")

	engine.Reload(reloadTemplates)
	err := engine.Load()
	if err != nil {
		log.Fatalf("Could not parse templates: %v", err)
	}

	engine.AddFunc("appbar", fiber.Map{
		"enabled": viper.GetBool("appbar.enabled"),
	})

	engine.AddFunc("renderMarkdown", format.MarkdownToHTML)

	renderer := templaterenderer.NewRenderer(engine)
	app := fiber.New(fiber.Config{CaseSensitive: true, DisableStartupMessage: true})
	b := NewPagesHandler(name, renderer)
	app.Get("/:path?+", b.Handle)

	return app
}
