/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */
package handler

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gzuidhof/starlit/starlit/internal/content"
	"github.com/gzuidhof/starlit/starlit/internal/templaterenderer"
	"github.com/gzuidhof/starlit/starlit/web"
	"github.com/spf13/afero"
	"github.com/spf13/viper"
)

type NBTestHandler struct {
	templaterenderer.TemplateRenderer
	fs afero.Fs
	
	starboardArtifactsURL string
	pyodideArtifactsURL string
}

func NewNBTestHandler(fs afero.Fs, templateRenderer templaterenderer.TemplateRenderer, starboardArtifactsFolder string, pyodideArtifactsFolder string) *NBTestHandler {
	sbArtifacts := viper.GetString("nbtest.starboard_artifacts")
	if starboardArtifactsFolder != "" {
		sbArtifacts = "/static/starboardArtifacts"
	} else if sbArtifacts == "" {
		sbArtifacts = fmt.Sprintf("/static/vendor/%s/dist", web.GetVendoredPackage("starboard-notebook"))
	}

	pyArtifacts := viper.GetString("nbtest.pyodide_artifacts")
	if pyodideArtifactsFolder != "" {
		pyArtifacts = "/static/pyodideArtifacts"
	}


	return &NBTestHandler{
		TemplateRenderer: templateRenderer,
		fs: fs,
		starboardArtifactsURL: strings.TrimSuffix(sbArtifacts, "/"),
		pyodideArtifactsURL: pyArtifacts,
	}
}

func (h *NBTestHandler) Handle(c *fiber.Ctx) error {
	path, err := url.QueryUnescape(c.Params("*"))
	if err != nil {
		log.Fatal(err)
	}

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

	if (page.FrontMatter.GetBool("nbtest.skip")) {
		return h.Render(c, "views/nbtest/notebook-test-skip", nil)
	}

	return h.Render(c, "views/nbtest/notebook-test", map[string]interface{} {
		"notebook_content": page.Content,
		"starboard_artifacts_url": h.starboardArtifactsURL,
		"pyodide_artifacts_url":  h.pyodideArtifactsURL,
	})
}
