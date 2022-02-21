package tests

import (
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/Komefumi/disease-storage/server"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
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

func testAPIResource(tests []APITest, t *testing.T) {
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
