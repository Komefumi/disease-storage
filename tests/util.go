package tests

import (
	"testing"

	"github.com/gofiber/fiber/v2"
)

type APITest struct {
	description     string
	route           string
	method          string
	isErrorExpected bool
	expectedCode    int
	preRunner       func(app *fiber.App, t *testing.T) error
	postRunner      func(bodyString string, t *testing.T) error
}

func emptyPreRunner(app *fiber.App, t *testing.T) error {
	return nil
}

func emptyPostRunner(bodyString string, t *testing.T) error {
	return nil
}
