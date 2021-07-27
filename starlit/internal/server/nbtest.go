package server

import (
	"log"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"github.com/gofiber/template/jet"
	"github.com/gzuidhof/starlit/starlit/internal/fs/assetfs"
	"github.com/gzuidhof/starlit/starlit/internal/fs/stripprefix"
	nbtesthandler "github.com/gzuidhof/starlit/starlit/internal/nbtest/handler"
	"github.com/gzuidhof/starlit/starlit/internal/templaterenderer"
	"github.com/spf13/afero"
)

func fixWasmContentTypeMiddleware(ctx *fiber.Ctx) error { 
			err := ctx.Next()
			if strings.HasSuffix(ctx.Path(), ".wasm") {
				ctx.Set("content-type", "application/wasm")
			}
			return err
}

func CreateNBTestApp(serveFolderAbs string, serveFS assetfs.ServeFS, starboardArtifactsFolder string, pyodideArtifactsFolder string) (*fiber.App, error) {
	app := fiber.New(fiber.Config{CaseSensitive: true, DisableStartupMessage: true})

	engine := jet.NewFileSystem(afero.NewHttpFs(serveFS.Templates), ".jet.html")
	renderer := templaterenderer.NewRenderer(engine)

	err := engine.Load()
	if err != nil {
		log.Fatal(err)
	}

	// engine.Reload(true)

	fs := afero.NewReadOnlyFs(afero.NewBasePathFs(afero.NewOsFs(), serveFolderAbs))
	nbTestHandler := nbtesthandler.NewNBTestHandler(fs, renderer, starboardArtifactsFolder, pyodideArtifactsFolder)

	if starboardArtifactsFolder != "" {
		app.Use("/static/starboardArtifacts/*", 
			fixWasmContentTypeMiddleware, 
			filesystem.New(filesystem.Config{
				Root: afero.NewHttpFs(
					stripprefix.New("/static/starboardArtifacts/",
					afero.NewReadOnlyFs(afero.NewBasePathFs(afero.NewOsFs(), starboardArtifactsFolder)))),
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
			}),
		)
	}

	app.Use("/static/*", fixWasmContentTypeMiddleware, filesystem.New(filesystem.Config{
		Root: afero.NewHttpFs(stripprefix.New("/static/", serveFS.Static)),
	}))

	app.Use("/nbtest/*", nbTestHandler.Handle)

	return app, nil
}

func StartNBTestServer(serveFolderAbs string, serveFS assetfs.ServeFS, port string, starboardArtifactsFolder string, pyodideArtifactsFolder string) error {
	app, err := CreateNBTestApp(serveFolderAbs, serveFS, starboardArtifactsFolder, pyodideArtifactsFolder)

	if err != nil {
		log.Fatal(err)
	}

	return app.Listen(":" + port)
}
