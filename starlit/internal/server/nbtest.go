/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */
package server

import (
	"fmt"
	"log"
	"strings"

	"github.com/fatih/color"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"github.com/gofiber/template/jet"
	"github.com/gzuidhof/starlit/starlit/internal/fs/assetfs"
	"github.com/gzuidhof/starlit/starlit/internal/fs/stripprefix"
	"github.com/gzuidhof/starlit/starlit/internal/middleware/coopcoep"
	nbtesthandler "github.com/gzuidhof/starlit/starlit/internal/nbtest/handler"
	"github.com/gzuidhof/starlit/starlit/internal/templaterenderer"
	"github.com/spf13/afero"
	"github.com/spf13/viper"
)

func fixWasmContentTypeMiddleware(ctx *fiber.Ctx) error {
	err := ctx.Next()
	if strings.HasSuffix(ctx.Path(), ".wasm") {
		ctx.Set("content-type", "application/wasm")
	}
	return err
}

func CreateNBTestApp(serveFolderAbs string, serveFS assetfs.ServeFS, starboardArtifactsFolder string, pyodideArtifactsFolder string, crossOriginIsolated bool) (*fiber.App, error) {
	app := fiber.New(fiber.Config{CaseSensitive: true, DisableStartupMessage: true})

	if crossOriginIsolated {
		app.Use(coopcoep.AddCOOPCOEPHeadersMiddleware)
	}

	engine := jet.NewFileSystem(afero.NewHttpFs(serveFS.Templates), ".jet.html")
	renderer := templaterenderer.NewRenderer(engine)

	err := engine.Load()
	if err != nil {
		log.Fatalf("Could not parse templates: %v", err)
	}

	fs := afero.NewReadOnlyFs(afero.NewBasePathFs(afero.NewOsFs(), serveFolderAbs))
	nbTestHandler := nbtesthandler.NewNBTestHandler(fs, renderer, starboardArtifactsFolder, pyodideArtifactsFolder)

	if starboardArtifactsFolder != "" {
		app.Use("/static/starboardArtifacts/*",
			fixWasmContentTypeMiddleware,
			filesystem.New(filesystem.Config{
				Root: afero.NewHttpFs(
					stripprefix.New("/static/starboardArtifacts/",
						afero.NewReadOnlyFs(afero.NewBasePathFs(afero.NewOsFs(), starboardArtifactsFolder)))),
				Browse: true,
			}),
		)
	}

	if pyodideArtifactsFolder != "" {
		app.Use("/static/pyodideArtifacts/*",
			fixWasmContentTypeMiddleware,
			filesystem.New(filesystem.Config{
				Root: afero.NewHttpFs(
					stripprefix.New("/static/pyodideArtifacts/",
						afero.NewReadOnlyFs(afero.NewBasePathFs(afero.NewOsFs(), pyodideArtifactsFolder)))),
				Browse: true,
			}),
		)
	}

	staticServes := viper.GetStringSlice("nbtest.serve_static")
	for _, ss := range staticServes {
		parts := strings.Split(ss, "=")
		if len(parts) != 2 {
			log.Fatalf("Invalid serve_static flag \"%s\", it should contain a single `=`.", ss)
		}
		routePath := fmt.Sprintf("/static/%s/", strings.Trim(parts[0], "/"))
		app.Use(routePath + "*", fixWasmContentTypeMiddleware, filesystem.New(filesystem.Config{
			Root: afero.NewHttpFs(
				stripprefix.New(routePath,
					afero.NewReadOnlyFs(afero.NewBasePathFs(afero.NewOsFs(), parts[1])))),
			Browse: true,
		}))
		fmt.Fprintf(color.Output,
				"%s serving files in %s\n",
				color.CyanString(routePath),
				color.BlueString(parts[1]),
		)
	}


	app.Use("/static/*", fixWasmContentTypeMiddleware, filesystem.New(filesystem.Config{
		Root: afero.NewHttpFs(stripprefix.New("/static/", serveFS.Static)),
	}))

	app.Use("/nbtest/*", nbTestHandler.Handle)

	return app, nil
}

// TODO: refactor these many args into an options object, or just read from viper directly
func StartNBTestServer(serveFolderAbs string, serveFS assetfs.ServeFS, port string, starboardArtifactsFolder string, pyodideArtifactsFolder string, crossOriginIsolated bool) error {
	app, err := CreateNBTestApp(serveFolderAbs, serveFS, starboardArtifactsFolder, pyodideArtifactsFolder, crossOriginIsolated)

	if err != nil {
		log.Fatal(err)
	}

	return app.Listen(":" + port)
}
