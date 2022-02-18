package main

import (
	"fmt"
	"log"

	"os"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Disease struct {
	gorm.Model
	Name        string `validate:"required,min=1" json:"name"`
	Description string `validate:"required,min=10" json:"description"`
}

type ErrorResponse struct {
	FailedField string
	Tag         string
	Value       string
}

var validate = validator.New()

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

	app.Post("/api/diseases", func(c *fiber.Ctx) error {
		disease := new(Disease)
		if err := c.BodyParser(disease); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": err.Error(),
			})
		}
		errors := ValidateStruct(*disease)

		if errors != nil {
			return c.Status(fiber.StatusBadRequest).JSON(errors)
		}

		result := db.Create(disease)
		if result.Error != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Failed to enter record",
			})
			log.Println(result.Error.Error())
		}

		return c.JSON(fiber.Map{
			"success": true,
			"message": "Successfully entered record",
			"disease": disease,
		})
	})

	log.Fatal(app.Listen(":3000"))
}

func ValidateStruct(payload interface{}) []*ErrorResponse {
	var errors []*ErrorResponse
	err := validate.Struct(payload)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			var element ErrorResponse
			element.FailedField = err.StructNamespace()
			element.Tag = err.Tag()
			element.Value = err.Param()
			errors = append(errors, &element)
		}
	}
	return errors
}
