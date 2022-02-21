package tests

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/Komefumi/disease-storage/database"
	"github.com/Komefumi/disease-storage/model"
	"github.com/gofiber/fiber/v2"
)

func TestDiseasesResource(t *testing.T) {
	tests := []APITest{
		{
			description:     "Get all the disease records",
			route:           "/api/diseases",
			method:          fiber.MethodGet,
			isErrorExpected: false,
			expectedCode:    200,
			preRunner:       GetAllDiseasesPreRunner,
			postRunner:      GetAllDiseasesPostRunner,
		},
	}
	testAPIResource(tests, t)
}

func GetAllDiseasesPreRunner(app *fiber.App, t *testing.T) error {
	err := database.InsertPrototypes()
	return err
}

func GetAllDiseasesPostRunner(bodyString string, t *testing.T) error {
	var body map[string]interface{}
	err := json.Unmarshal([]byte(bodyString), &body)
	if err != nil {
		return err
	}
	var diseaseIDs []uint
	for _, currentDiseaseData := range body["diseases"].([]interface{}) {
		currentID := uint(currentDiseaseData.(map[string]interface{})["ID"].(float64))
		diseaseIDs = append(diseaseIDs, currentID)
	}
	t.Log(diseaseIDs)
	if len(diseaseIDs) == 0 {
		return errors.New("Must return non-empty collection of diseases")
	}
	if diseaseIDs[0] != model.PrototypeDisease.ID {
		return errors.New("Must return prototype disease record that was created")
	}

	return nil
}
