package main

import (
	"fmt"
	"log"

	"os"

	"github.com/gofiber/fiber/v2"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Disease struct {
	gorm.Model
	Name        string `json:"name"`
	Description string `json:"description"`
}

func main() {
	err := os.Remove("test.db")
	if err != nil {
		fmt.Println(err)
	}
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	db.AutoMigrate(&Disease{})

	db.Create(&Disease{Name: "ProtoType Disease", Description: "Non real disease, made as a model to perform operations with"})

	app := fiber.New()
	app.Get("/api/diseases", func(c *fiber.Ctx) error {
		var diseases []Disease
		var finalMessage *fiber.Map
		result := db.Find(&diseases)
		if result.RowsAffected > 0 && result.Error == nil {
			finalMessage = &fiber.Map{
				"success":  false,
				"message":  "Successfully retrieved disease data",
				"diseases": diseases,
			}
		}
		if result.RowsAffected == 0 {
			finalMessage = &fiber.Map{
				"success": false,
				"message": "No disease data found",
			}
		}
		if result.Error != nil {
			finalMessage = &fiber.Map{
				"success": false,
				"message": "Error: Failed to retrieve disease data",
				"error":   result.Error.Error,
			}
		}
		return c.JSON(finalMessage)
	})

	log.Fatal(app.Listen(":3000"))
}
