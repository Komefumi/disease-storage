package main

import (
	"fmt"
	"log"
	"strconv"

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

var db *gorm.DB

func init() {
	err := os.Remove("test.db")
	if err != nil {
		fmt.Println(err)
	}
	dbOpened, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	dbOpened.AutoMigrate(&Disease{})

	// dbOpened.Create(&Disease{Name: "ProtoType Disease", Description: "Non real disease, made as a model to perform operations with"})

	db = dbOpened
}

func main() {
	app := fiber.New()
	app.Get("/api/diseases", func(c *fiber.Ctx) error {
		var diseases []Disease
		result := db.Find(&diseases)
		if result.RowsAffected == 0 {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"success": false,
				"message": "No disease data found",
			})
		}
		if result.Error != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"success": false,
				"message": "Error: Failed to retrieve disease data",
				"error":   result.Error.Error,
			})
		}
		return c.JSON(fiber.Map{
			"success":  false,
			"message":  "Successfully retrieved disease data",
			"diseases": diseases,
		})
	})

	app.Post("/api/diseases", func(c *fiber.Ctx) error {
		disease := new(Disease)
		if err := c.BodyParser(disease); err != nil {
			return handleBodyParseError(err, c)
			fmt.Println(disease)
		}
		errors := ValidateStruct(disease)
		if errors != nil {
			return handleInvalidBodyError(errors, c)
		}
		result := db.Create(disease)
		if result.Error != nil {
			return handleDBError("Failed to enter record", result.Error, c)
		}

		return c.JSON(fiber.Map{
			"success": true,
			"message": "Successfully entered record",
			"disease": disease,
		})
	})

	app.Get("/api/diseases/:id", func(c *fiber.Ctx) error {
		var disease Disease
		id, errId := strconv.Atoi(c.Params("id"))
		if errId != nil {
			return handleInvalidIDError(c)
		}
		result := db.First(&disease).Where("id = ?", id)
		if result.Error != nil {
			fmt.Println(id)
			return handleDBError(fmt.Sprintf("Failed to find disease of id %v", id), result.Error, c)
		}
		return c.JSON(fiber.Map{
			"success": true,
			"message": "Successfully found disease",
			"disease": disease,
		})
	})

	app.Put("/api/diseases/:id", func(c *fiber.Ctx) error {
		var disease Disease
		diseaseUpdates := new(Disease)
		id, errId := strconv.Atoi(c.Params("id"))
		if errId != nil {
			return handleInvalidIDError(c)
		}
		result := db.First(&disease).Where("id = ?", id)
		if result.Error != nil {
			fmt.Println(id)
			return handleDBError(fmt.Sprintf("Failed to find disease of id %v", id), result.Error, c)
		}
		if err := c.BodyParser(diseaseUpdates); err != nil {
			return handleBodyParseError(err, c)
			fmt.Println(disease)
		}
		errors := ValidateStruct(diseaseUpdates)
		if errors != nil {
			return handleInvalidBodyError(errors, c)
		}
		result = db.Model(&disease).Updates(diseaseUpdates)
		if result.Error != nil {
			return handleDBError(fmt.Sprintf("Failed to update disease of id %v", id), result.Error, c)
		}
		return c.JSON(fiber.Map{
			"success": true,
			"message": "Successfully updated disease",
			"disease": disease,
		})
	})

	app.Delete("/api/diseases/:id", func(c *fiber.Ctx) error {
		var disease Disease
		id, errId := strconv.Atoi(c.Params("id"))
		if errId != nil {
			return handleInvalidIDError(c)
		}
		result := db.First(&disease).Where("id = ?", id)
		if result.Error != nil {
			fmt.Println(id)
			return handleDBError(fmt.Sprintf("Failed to find disease of id %v", id), result.Error, c)
		}
		result = db.Delete(&Disease{}, id)
		if result.Error != nil {
			return handleDBError(fmt.Sprintf("Failed to delete disease of id %v", id), result.Error, c)
		}
		return c.JSON(fiber.Map{
			"success": true,
			"message": "Successfully deleted disease",
			"disease": disease,
		})
	})

	log.Fatal(app.Listen(":3000"))
}

func handleInvalidIDError(c *fiber.Ctx) error {
	return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
		"message": "Must provide valid numeric id",
	})
}

func handleBodyParseError(err error, c *fiber.Ctx) error {
	return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
		"message": err.Error(),
	})
}

func handleInvalidBodyError(errors []*ErrorResponse, c *fiber.Ctx) error {
	return c.Status(fiber.StatusBadRequest).JSON(errors)
}

func handleDBError(responseMessage string, err error, c *fiber.Ctx) error {
	log.Println(err.Error())
	return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
		"message": responseMessage,
	})
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
