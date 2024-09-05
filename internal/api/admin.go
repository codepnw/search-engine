package api

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/codepnw/search-engine/internal/db"
	"github.com/codepnw/search-engine/internal/utils"
	"github.com/codepnw/search-engine/internal/views"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

type loginForm struct {
	Email    string `form:"email"`
	Password string `form:"password"`
}

type AdminClaims struct {
	User                 string `json:"user"`
	Id                   string `json:"id"`
	jwt.RegisteredClaims `json:"claims"`
}

type settingsForm struct {
	Amount   int    `form:"amount"`
	SearchOn string `form:"searchOn"`
	AddNew   string `form:"addNew"`
}

func LoginHandler(c *fiber.Ctx) error {
	return render(c, views.Login())
}

func LoginPostHandler(c *fiber.Ctx) error {
	time.Sleep(2 * time.Second)

	input := loginForm{}
	if err := c.BodyParser(&input); err != nil {
		c.Status(fiber.ErrBadRequest.Code)
		return c.SendString("<h2>Error: something went wrong</h2>")
	}

	user := &db.User{}
	user, err := user.LoginAsAdmin(input.Email, input.Password)
	if err != nil {
		c.Status(fiber.ErrUnauthorized.Code)
		return c.SendString("<h2>Error: unauthorized")
	}

	signedToken, err := utils.CreateNewAuthToken(user.ID, user.Email, user.IsAdmin)
	if err != nil {
		c.Status(fiber.ErrUnauthorized.Code)
		return c.SendString("<h2>Error: something went wrong logging in")
	}

	cookie := fiber.Cookie{
		Name:     "admin",
		Value:    signedToken,
		Expires:  time.Now().Add(24 * time.Hour),
		HTTPOnly: true,
	}

	c.Cookie(&cookie)
	c.Append("HX-Redirect", "/")
	return c.SendStatus(fiber.StatusOK)
}

func LogoutHandler(c *fiber.Ctx) error {
	c.ClearCookie("admin")
	c.Set("HX-Redirect", "/login")
	return c.SendStatus(fiber.StatusOK)
}

func AuthMiddleware(c *fiber.Ctx) error {
	cookie := c.Cookies("admin")
	if cookie == "" {
		return c.Redirect("/login", fiber.StatusFound)
	}

	token, err := jwt.ParseWithClaims(cookie, &AdminClaims{}, func(t *jwt.Token) (any, error) {
		return []byte(os.Getenv("SECRET_KEY")), nil
	})
	if err != nil {
		return c.Redirect("/login", fiber.StatusFound)
	}

	_, ok := token.Claims.(*AdminClaims)
	if ok && token.Valid {
		return c.Next()
	}
	return c.Redirect("/login", fiber.StatusFound)
}

func DashboardHandler(c *fiber.Ctx) error {
	settings := db.SearchSettings{}

	if err := settings.Get(); err != nil {
		c.Status(fiber.ErrInternalServerError.Code)
		return c.SendString("<h2>Error: cannot get settings</h2>")
	}

	amount := strconv.FormatUint(uint64(settings.Amount), 10)
	return render(c, views.Home(amount, settings.SearchOn, settings.AddNew))
}

func DashboardPostHandler(c *fiber.Ctx) error {
	input := settingsForm{}

	if err := c.BodyParser(&input); err != nil {
		c.Status(fiber.ErrBadRequest.Code)
		return c.SendString("<h2>Error: cannot get settings</h2>")
	}

	addNew := false
	if input.AddNew == "on" {
		addNew = true
	}

	searchOn := false
	if input.SearchOn == "on" {
		searchOn = true
	}

	settings := &db.SearchSettings{}
	settings.Amount = uint(input.Amount)
	settings.SearchOn = searchOn
	settings.AddNew = addNew

	if err := settings.Update(); err != nil {
		fmt.Println(err)
		return c.SendString("<h2>Error: cannot update settings</h2>")
	}

	c.Append("HX-Refresh", "true")
	return c.SendStatus(fiber.StatusOK)
}
