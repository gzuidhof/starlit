/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */
package pages

import (
	"fmt"
	"log"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gzuidhof/starlit/starlit/internal/config"
	"github.com/gzuidhof/starlit/starlit/internal/content"
	"github.com/gzuidhof/starlit/starlit/internal/templaterenderer"
	"github.com/spf13/afero"
	"github.com/spf13/viper"
)

type PagesHandler struct {
	spaceName string
	templaterenderer.TemplateRenderer
	config *viper.Viper
	fs afero.Fs
	menu MenuData
}

func NewPagesHandler(spaceName string, templateRenderer templaterenderer.TemplateRenderer) *PagesHandler {
	vconfig := config.GetSpaceConfig(spaceName, viper.GetViper())
	servePath := vconfig.GetString("serve_filepath")

	fs := afero.NewReadOnlyFs(afero.NewBasePathFs(afero.NewOsFs(), servePath))

	menu, err := buildContentTree(".", fs, vconfig.GetString("space.path"))
	if err != nil {
		log.Fatalf("Failed to construct initial menu for space %s: %v", spaceName, err)
	}

	return &PagesHandler{
		spaceName: spaceName,
		fs: fs,
		TemplateRenderer: templateRenderer,
		config: vconfig,
		menu: menu,
	}
}

func (h *PagesHandler) reloadContent() {
	menu, err := buildContentTree(".", h.fs, h.config.GetString("space.path"))
	if err != nil {
		log.Print(fmt.Errorf("Failed to construct initial menu for space %s: %v", h.spaceName, err))
	}

	h.menu = menu;
}

func (h *PagesHandler) renderNotebookPage(c *fiber.Ctx, templateData map[string]interface{}) error {
	return h.Render(c, "views/pages/notebook", templateData)
}

func (h *PagesHandler) renderMarkdownPage(c *fiber.Ctx, templateData map[string]interface{}) error {
	return h.Render(c, "views/pages/markdown", templateData)
}

func (h *PagesHandler) Handle(c *fiber.Ctx) error {
	if (h.config.GetBool("reload")) {
		h.reloadContent()
	}

	path := "/" + strings.TrimPrefix(c.Path(), h.config.GetString("space.path"))
	log.Printf("Space %s, path %s", h.spaceName, path)
	page := h.menu.ResolvePage(path)

	if (page == nil) {
		pageConfig := h.config
		pageConfig.Set("menu", h.menu)
		pageConfig.Set("path", path)
		return h.Render(c, "views/pages/404", pageConfig.AllSettings())
	}

	// TODO: load files live perhaps
	// file, err := h.fs.Open(path)
	// if err != nil {
	// 	// log.Printf("Failed to open %s: %v", path, err)
	// 	return c.SendStatus(404)
	// }
	// defer file.Close()
	
	// page, err  := content.ReadPageFile(path, file)
	// if err != nil {
	// 	return fmt.Errorf("failed to read file as page %s: %v", path, err)
	// }

	
	pageConfig := config.GetPageConfig(path, h.config, page.FrontMatter)
	pageConfig.Set("page", *page)
	pageConfig.Set("menu", h.menu)
	pageConfig.Set("path", path)
	templateData := pageConfig.AllSettings()

	if page.Type == content.Markdown {
		return h.renderMarkdownPage(c, templateData)
	} else if page.Type == content.Notebook {
		return h.renderNotebookPage(c, templateData)
	}

	return c.SendString("Unsupported filetype")
}
