package tests

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/Komefumi/disease-storage/database"
	"github.com/Komefumi/disease-storage/model"
	"github.com/Komefumi/disease-storage/server"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func TestDiseasesResource(t *testing.T) {
	t.Log("uh")
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
	app := server.Setup()

	for _, test := range tests {
		preErr := test.preRunner(app, t)
		assert.Nil(t, preErr, test.description)
		req, _ := http.NewRequest(test.method, test.route, nil)
		res, errInResponse := app.Test(req, -1)

		assert.Equal(t, test.isErrorExpected, errInResponse != nil, test.description)

		if test.isErrorExpected {
			continue
		}

		defer res.Body.Close()
		body, bodyErr := ioutil.ReadAll(res.Body)

		assert.Nil(t, bodyErr, test.description)

		t.Log(string(body))
		postErr := test.postRunner(string(body), t)
		assert.Nil(t, postErr, test.description)
	}
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
	/*
		if !diseasesCastIsOk {
			t.Log(body["diseases"])
			return errors.New("diseases must be a collection of disease records")
		}
	*/
	if len(diseaseIDs) == 0 {
		return errors.New("Must return non-empty collection of diseases")
	}
	if diseaseIDs[0] != model.PrototypeDisease.ID {
		return errors.New("Must return prototype disease record that was created")
	}

	return nil
}
