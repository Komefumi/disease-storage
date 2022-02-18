package handler

import (
	"log"

	"github.com/Komefumi/disease-storage/validation"
	"github.com/gofiber/fiber/v2"
)

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

func handleInvalidBodyError(errors []*validation.ErrorResponse, c *fiber.Ctx) error {
	return c.Status(fiber.StatusBadRequest).JSON(errors)
}

func handleDBError(responseMessage string, err error, errorStatus int, c *fiber.Ctx) error {
	log.Println(err.Error())
	return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
		"message": responseMessage,
	})
}
