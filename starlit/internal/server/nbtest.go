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
	"github.com/spf13/viper"
)

func CreateNBTestApp(serveFolderAbs string, serveFS assetfs.ServeFS) (*fiber.App, error) {
	app := fiber.New(fiber.Config{CaseSensitive: true, DisableStartupMessage: true})

	engine := jet.NewFileSystem(afero.NewHttpFs(serveFS.Templates), ".jet.html")
	renderer := templaterenderer.NewRenderer(engine)

	engine.Reload(true)
	// app.Use(logger.New())
	

	starboardArtifactsURL := strings.TrimSuffix(viper.GetString("sandbox_url"), "/index.html")
	pyodideArtifactsURL := ""

	fs := afero.NewReadOnlyFs(afero.NewBasePathFs(afero.NewOsFs(), serveFolderAbs))
	nbTestHandler := nbtesthandler.NewNBTestHandler(fs, starboardArtifactsURL, pyodideArtifactsURL, renderer)

	app.Use("/static/*", filesystem.New(filesystem.Config{
		Root: afero.NewHttpFs(stripprefix.New("/static/", serveFS.Static)),
	}))

	app.Use("/nbtest/*", nbTestHandler.Handle)

	return app, nil
}

func startNBTestServer(serveFolderAbs string, serveFS assetfs.ServeFS, port string) error {
	app, err := CreateNBTestApp(serveFolderAbs, serveFS)

	if err != nil {
		log.Fatal(err)
	}

	return app.Listen(":" + port)
}
