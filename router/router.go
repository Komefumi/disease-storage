package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"

	"github.com/Komefumi/disease-storage/handler"
)

func SetupRoutes(app *fiber.App) {
	api := app.Group("/api", logger.New())
	diseases := api.Group("/diseases")
	diseases.Get("/", handler.GetAllDiseases)
	diseases.Post("/", handler.CreateDiseaseRecord)
	diseases.Get("/:id", handler.GetDiseaseByID)
	diseases.Put("/:id", handler.UpdateDiseaseByID)
	diseases.Delete("/:id", handler.DeleteDiseaseByID)
}
