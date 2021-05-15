/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */
package pages

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/gzuidhof/starlit/starlit/internal/config"
	"github.com/gzuidhof/starlit/starlit/internal/content"
	"github.com/gzuidhof/starlit/starlit/internal/templaterenderer"
	"github.com/spf13/afero"
	"github.com/spf13/viper"
)

type PagesHandler struct {
	templaterenderer.TemplateRenderer
	config *viper.Viper
	fs afero.Fs
}

func NewPagesHandler(spaceName string, templateRenderer templaterenderer.TemplateRenderer) *PagesHandler {
	vconfig := config.GetSpaceConfig(spaceName, viper.GetViper())
	fs := afero.NewReadOnlyFs(afero.NewBasePathFs(afero.NewOsFs(), vconfig.GetString("serve_filepath")))

	return &PagesHandler{
		fs: fs,
		TemplateRenderer: templateRenderer,
		config: vconfig,
	}
}

func (h *PagesHandler) renderNotebookPage(c *fiber.Ctx, page content.Page) error {
	return h.Render(c, "views/pages/notebook", fiber.Map{
		"sandboxUrl": "/static/vendor/starboard-notebook@0.9.4/dist/index.html",
		"page": page,
	})
}

func (h *PagesHandler) renderMarkdownPage(c *fiber.Ctx, page content.Page) error {
	return h.Render(c, "views/pages/markdown", fiber.Map{
		"page": page,
	})
}

func (h *PagesHandler) Handle(c *fiber.Ctx) error {
	path := c.Path()
	file, err := h.fs.Open(path)
	defer file.Close()

	if err != nil {
		return fmt.Errorf("failed to open %s: %v", path, err)
	}
	page, err  := content.ReadPageFile(path, file)
	if err != nil {
		return fmt.Errorf("failed to read file as page %s: %v", path, err)
	}

	if page.Type == content.Markdown {
		return h.renderMarkdownPage(c, page)
	} else if page.Type == content.Notebook {
		return h.renderNotebookPage(c, page)
	}

	return c.SendString("Unsupported filetype")
}
