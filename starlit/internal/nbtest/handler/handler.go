/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */
package handler

import (
	"log"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/gzuidhof/starlit/starlit/internal/content"
	"github.com/gzuidhof/starlit/starlit/internal/templaterenderer"
	"github.com/spf13/afero"
)

type NBTestHandler struct {
	templaterenderer.TemplateRenderer
	fs afero.Fs
	starboardArtifactsURL string
	pyodideArtifactsURL string
}

func NewNBTestHandler(fs afero.Fs, starboardArtifactsURL string, pyodideArtifactsURL string, templateRenderer templaterenderer.TemplateRenderer) *NBTestHandler {
	return &NBTestHandler{
		TemplateRenderer: templateRenderer,
		fs: fs,
		starboardArtifactsURL: starboardArtifactsURL,
		pyodideArtifactsURL: pyodideArtifactsURL,
	}
}

func (h *NBTestHandler) Handle(c *fiber.Ctx) error {
	path := c.Params("*")

	f, err := h.fs.Open(path)
	// TODO handle 404 case
	if err != nil {
		log.Printf("Error in nbtest handler: %v", err)
		return c.Status(http.StatusInternalServerError).SendString(err.Error())
	}
	defer f.Close()

	page, err := content.ReadPageFile(path, f)
	if err != nil {
		log.Printf("Error in nbtest handler reading page: %v", err)
		return c.Status(http.StatusInternalServerError).SendString(err.Error())
	}

	return h.Render(c, "views/nbtest/notebook-test", map[string]interface{} {
		"notebook_content": page.Content,
		"starboard_cdn": h.starboardArtifactsURL,
		"pyodide_cdn":  h.pyodideArtifactsURL,
	})
}
