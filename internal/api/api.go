package api

import (
	"github.com/a-h/templ"
	"github.com/gofiber/fiber/v2"
)

func render(c *fiber.Ctx, component templ.Component) error {
	c.Set("Content-Type", "text/html")
	return component.Render(c.Context(), c.Response().BodyWriter())
}

func NewRoutes(app *fiber.App) {
	// Home
	app.Get("/", AuthMiddleware, DashboardHandler)
	app.Post("/", AuthMiddleware, DashboardPostHandler)

	// Auth
	app.Get("/login", LoginHandler)
	app.Post("/login", LoginPostHandler)
	app.Get("/logout", LogoutHandler)

	// Create Admin
	// app.Get("/create", func(c *fiber.Ctx) error {
	// 	u := &db.User{}
	// 	u.CreateAdmin()
	// 	return c.SendString("created")
	// })

	app.Post("/search", HandleSearch)
}
