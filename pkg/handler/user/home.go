package user

import (
	"github.com/gofiber/fiber/v2"
)

// HomeHandler godoc
// @Summary Ana sayfa
// @Description Hoş geldin mesajı döner
// @Tags Home
// @Success 200 {string} string "Hoş geldin"
// @Router / [get]
func HomeHandler(c *fiber.Ctx) error {
	return c.SendString("Hoş geldin!")
}
