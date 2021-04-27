/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */
package templaterenderer

import (
	"github.com/gofiber/fiber/v2"
	"github.com/valyala/bytebufferpool"
)

type TemplateRenderer struct {
	views fiber.Views
}

func NewRenderer(v fiber.Views) TemplateRenderer {
	return TemplateRenderer{
		views: v,
	}
}

func (r *TemplateRenderer) Render(ctx *fiber.Ctx, name string, bind interface{}, layouts ...string) error {
	buf := bytebufferpool.Get()
	defer bytebufferpool.Put(buf)

	err := r.views.Render(buf, name, bind, layouts...)
	if err != nil {
		return err
	}

	resp := ctx.Response()
	// Set Content-Type to text/html
	resp.Header.SetContentType(fiber.MIMETextHTMLCharsetUTF8)
	// Set rendered template to body
	resp.SetBody(buf.Bytes())

	return nil
}
