package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()
	app.Get("/api/diseases", func(c *fiber.Ctx) error {
		return c.JSON(&fiber.Map{
			"success": true,
			"message": "The start has been established!",
		})
	})

	log.Fatal(app.Listen(":3000"))
}
