package api

import (
	"github.com/codepnw/search-engine/internal/db"
	"github.com/gofiber/fiber/v2"
)

type searchInput struct {
	Term string `json:"term"`
}

func HandleSearch(c *fiber.Ctx) error {
	input := searchInput{}
	if err := c.BodyParser(&input); err != nil {
		c.Status(fiber.ErrInternalServerError.Code)
		c.Append("content-type", "application/json")
		return c.JSON(fiber.Map{
			"success": false,
			"message": "Invalid input",
			"data":    nil,
		})
	}
	
	if input.Term == "" {
		c.Status(500)
		c.Append("content-type", "application/json")
		return c.JSON(fiber.Map{
			"success": false,
			"message": "Invalid input",
			"data":    nil,
		})
	}

	idx := &db.SearchIndex{}
	data, err := idx.FullTextSearch(input.Term)
	if err != nil {
		c.Status(fiber.ErrInternalServerError.Code)
		c.Append("content-type", "application/json")
		return c.JSON(fiber.Map{
			"success": false,
			"message": "Invalid input",
			"data":    nil,
		})
	}

	c.Status(fiber.StatusOK)
	c.Append("content-type", "application/json")
	return c.JSON(fiber.Map{
		"success": true,
		"message": "Search results",
		"data":    data,
	})
}
