package main

import (
	"bytes"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func TestHealthCheck(t *testing.T) {

	tests := []struct {
		Name           string
		Path           string
		method         string
		expectedStatus int
		ExpectedBody   string
	}{
		{
			Name:           "Health check",
			Path:           "/health",
			method:         "GET",
			expectedStatus: 200,
			ExpectedBody:   `{"message":"Healthy"}`,
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			app := fiber.New()
			app.Get(test.Path, HealthCheck)
			req := httptest.NewRequest(test.method, test.Path, nil)
			res, _ := app.Test(req)

			assert.Equal(t, test.expectedStatus, res.StatusCode, test.Name)
			buff := new(bytes.Buffer)
			_, err := buff.ReadFrom(res.Body)
			assert.NoError(t, err)
			assert.Equal(t, test.ExpectedBody, buff.String())
		})
	}
}

func TestAnotherHealthCheck(t *testing.T) {

	tests := []struct {
		Name           string
		Path           string
		method         string
		expectedStatus int
		ExpectedBody   string
	}{
		{
			Name:           "Health check",
			Path:           "/health",
			method:         "GET",
			expectedStatus: 200,
			ExpectedBody:   `{"message":"Healthy"}`,
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			app := fiber.New()
			app.Get(test.Path, HealthCheck)
			req := httptest.NewRequest(test.method, test.Path, nil)
			res, _ := app.Test(req)

			assert.Equal(t, test.expectedStatus, res.StatusCode, test.Name)
			buff := new(bytes.Buffer)
			_, err := buff.ReadFrom(res.Body)
			assert.NoError(t, err)
			assert.Equal(t, test.ExpectedBody, buff.String())
		})
	}
}
