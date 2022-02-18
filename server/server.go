package server

import (
	"github.com/Komefumi/disease-storage/router"
	"github.com/gofiber/fiber/v2"
)

func Setup() *fiber.App {
	// Initialize a new app
	app := fiber.New()
	router.SetupRoutes(app)
	// Return the configured app
	return app
}
