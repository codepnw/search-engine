package api

import (
	"time"

	"github.com/a-h/templ"
	"github.com/codepnw/search-engine/internal/views"
	"github.com/gofiber/fiber/v2"
)

func render(c *fiber.Ctx, component templ.Component) error {
	c.Set("Content-Type", "text/html")
	return component.Render(c.Context(), c.Response().BodyWriter())
}

type loginForm struct {
	Email    string `form:"email"`
	Password string `form:"password"`
}

type settingsForm struct {
	Amount   int  `form:"amount"`
	SearchOn bool `form:"searchOn"`
	AddNew   bool `form:"addNew"`
}

func NewRoutes(app *fiber.App) {
	// Home
	app.Get("/", func(c *fiber.Ctx) error {
		return render(c, views.Home())
	})
	app.Post("/", func(c *fiber.Ctx) error {
		time.Sleep(2 * time.Second)

		input := settingsForm{}
		if err := c.BodyParser(&input); err != nil {
			return c.SendString("<h2>Error: something went wrong</h2>")
		}
		return c.SendStatus(200)
	})

	// Login
	app.Get("/login", func(c *fiber.Ctx) error {
		return render(c, views.Login())
	})
	app.Post("/login", func(c *fiber.Ctx) error {
		time.Sleep(2 * time.Second)

		input := loginForm{}
		if err := c.BodyParser(&input); err != nil {
			return c.SendString("<h2>Error: something went wrong</h2>")
		}
		return c.SendStatus(200)
	})
}
