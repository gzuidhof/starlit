/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */
package pages

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gzuidhof/starlit/starlit/internal/jetrender"
)

type BookHandler struct {
	jetrender.TemplateRenderer
}

func NewBookHandler(templateRenderer jetrender.TemplateRenderer) *BookHandler {
	return &BookHandler{
		templateRenderer,
	}

}

func (h *BookHandler) Handle(c *fiber.Ctx) error {
	return h.Render(c, "views/book/book", fiber.Map{})
}
