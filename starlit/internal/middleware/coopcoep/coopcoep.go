/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */
package coopcoep

import "github.com/gofiber/fiber/v2"

func AddCOOPCOEPHeadersMiddleware(ctx *fiber.Ctx) error {
	err := ctx.Next()

	ctx.Set("Cross-Origin-Embedder-Policy", "require-corp")
	ctx.Set("Cross-Origin-Opener-Policy", "same-origin")
	return err
}
