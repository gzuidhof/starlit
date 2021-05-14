/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */
package pages

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gzuidhof/starlit/starlit/internal/config"
	"github.com/gzuidhof/starlit/starlit/internal/templaterenderer"
	"github.com/spf13/viper"
)

type PagesHandler struct {
	templaterenderer.TemplateRenderer
	config *viper.Viper
}

func NewPagesHandler(spaceName string, templateRenderer templaterenderer.TemplateRenderer) *PagesHandler {
	vconfig := config.GetSpaceConfig(spaceName, viper.GetViper())

	return &PagesHandler{
		TemplateRenderer: templateRenderer,
		config: vconfig,
	}
}

func (h *PagesHandler) renderNotebookPage(c *fiber.Ctx, notebookBytes []byte) error {
	return h.Render(c, "views/pages/notebook", fiber.Map{
		"sandboxUrl": "/static/vendor/starboard-notebook@0.9.1/dist/index.html",
		"notebookContent": notebookBytes,
	})
}

func (h *PagesHandler) renderMarkdownPage(c *fiber.Ctx, markdownBytes []byte) error {
	return h.Render(c, "views/pages/markdown", fiber.Map{
		"markdownContent": markdownBytes,
	})
}

func (h *PagesHandler) Handle(c *fiber.Ctx) error {
	r := []byte("# Title")
	return h.renderNotebookPage(c, r)
}
