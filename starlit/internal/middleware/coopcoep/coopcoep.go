package coopcoep

import "github.com/gofiber/fiber/v2"

func AddCOOPCOEPHeadersMiddleware(ctx *fiber.Ctx) error {
	err := ctx.Next()

	ctx.Set("Cross-Origin-Embedder-Policy", "require-corp")
	ctx.Set("Cross-Origin-Opener-Policy", "same-origin")
	return err
}
