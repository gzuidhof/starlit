/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */
package book

import (
	"github.com/gofiber/fiber/v2"
)

type BookHandler struct {
}

func (h *BookHandler) Handle(c *fiber.Ctx) error {
	return c.Render("views/book/book", fiber.Map{})
}
