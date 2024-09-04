package api

import (
	"github.com/a-h/templ"
	"github.com/codepnw/search-engine/internal/views"
	"github.com/gofiber/fiber/v2"
)

func render(c *fiber.Ctx, component templ.Component) error {
	c.Set("Content-Type", "text/html")
	return component.Render(c.Context(), c.Response().BodyWriter())
}

func NewRoutes(app *fiber.App) {
	app.Get("/", func(c *fiber.Ctx) error {
		return render(c, views.Home())
	})
}
