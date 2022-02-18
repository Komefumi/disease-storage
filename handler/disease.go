package handler

import (
	"fmt"
	"strconv"

	"github.com/Komefumi/disease-storage/database"
	"github.com/Komefumi/disease-storage/model"
	"github.com/Komefumi/disease-storage/validation"
	"github.com/gofiber/fiber/v2"
)

var db = database.DB

func GetAllDiseases(c *fiber.Ctx) error {
	var diseases []model.Disease
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
}

func GetDiseaseByID(c *fiber.Ctx) error {
	var disease model.Disease
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
}

func CreateDiseaseRecord(c *fiber.Ctx) error {
	disease := new(model.Disease)
	if err := c.BodyParser(disease); err != nil {
		return handleBodyParseError(err, c)
		fmt.Println(disease)
	}
	errors := validation.ValidateStruct(disease)
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
}

func UpdateDiseaseByID(c *fiber.Ctx) error {
	var disease model.Disease
	diseaseUpdates := new(model.Disease)
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
	errors := validation.ValidateStruct(diseaseUpdates)
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
}

func DeleteDiseaseByID(c *fiber.Ctx) error {
	var disease model.Disease
	id, errId := strconv.Atoi(c.Params("id"))
	if errId != nil {
		return handleInvalidIDError(c)
	}
	result := db.First(&disease).Where("id = ?", id)
	if result.Error != nil {
		fmt.Println(id)
		return handleDBError(fmt.Sprintf("Failed to find disease of id %v", id), result.Error, c)
	}
	result = db.Delete(&model.Disease{}, id)
	if result.Error != nil {
		return handleDBError(fmt.Sprintf("Failed to delete disease of id %v", id), result.Error, c)
	}
	return c.JSON(fiber.Map{
		"success": true,
		"message": "Successfully deleted disease",
		"disease": disease,
	})
}
