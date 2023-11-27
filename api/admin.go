package api

import (
	"fmt"

	"github.com/DebasisOnDev/E-Commerce-Backend/types"
	"github.com/gofiber/fiber/v2"
)

func AdminAuth(c *fiber.Ctx) error {
	user, ok := c.Context().UserValue("user").(*types.User)
	if !ok {
		fmt.Println("received user is not ok")
		return c.JSON(fiber.Map{"error": "Unauthorized user user is not ok"})
	}
	if !user.IsAdmin {
		fmt.Println("user is not admin")
		return c.JSON(fiber.Map{"error": "Unauthorized user user is not admin"})
	}
	return c.Next()
}
