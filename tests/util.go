package tests

import (
	"errors"
	"fmt"
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

func testAPISuite(passedBaseURL interface{}, tests []APITest, t *testing.T) {
	app := server.Setup()
	var isBaseURLNil bool
	var baseURL string

	switch passedBaseURL.(type) {
	case string:
		isBaseURLNil = false
	case nil:
		isBaseURLNil = true
	default:
		panic(errors.New("passedBaseURL has to be a string if provided"))
	}

	if !isBaseURLNil {
		baseURL = passedBaseURL.(string)
	}

	for _, test := range tests {
		preErr := test.preRunner(app, t)
		assert.Nil(t, preErr, test.description)
		route := fmt.Sprintf("%s%s", baseURL, test.route)
		req, _ := http.NewRequest(test.method, route, nil)
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
