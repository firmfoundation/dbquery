package handler

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
)

type DbConfig struct {
	DBHost         string `mapstructure:"POSTGRES_HOST"`
	DBUserName     string `mapstructure:"POSTGRES_USER"`
	DBUserPassword string `mapstructure:"POSTGRES_PASSWORD"`
	DBName         string `mapstructure:"POSTGRES_DB"`
	DBPort         string `mapstructure:"POSTGRES_PORT"`
}

func TestCreateConnectionHandler(t *testing.T) {
	// Create a new Fiber app
	app := fiber.New()

	// Define the test case request body
	requestBody := map[string]string{
		"DBHost":         "127.0.0.1",
		"DBUserName":     "postgres",
		"DBUserPassword": "password123",
		"DBName":         "testdb",
		"DBPort":         "5432",
	}
	requestBodyBytes, _ := json.Marshal(requestBody)

	// Create a new test request
	req := httptest.NewRequest(http.MethodPost, "/connect", bytes.NewReader(requestBodyBytes))
	req.Header.Set("Content-Type", "application/json") // Set the JSON header
	// Create a new test response recorder
	//res := httptest.NewRecorder()

	// Call the handler function with the test request and response
	app.Post("/connect", CreateConnectionHandler)
	// Simulate the request and capture the response
	res, err := app.Test(req, -1)
	if err != nil {
		t.Errorf("Error testing request: %v", err)
	}

	// Check the response status code
	if res.StatusCode != http.StatusAccepted {
		t.Errorf("Expected status code %d but got %d", http.StatusAccepted, res.StatusCode)
	}

	// Check the response body status
	var responseBody struct {
		Status string `json:"status"`
	}
	bodyBytes, err := ioutil.ReadAll(res.Body)
	err = json.Unmarshal(bodyBytes, &responseBody)
	if err != nil {
		t.Errorf("Error decoding response body: %v", err)
	}
	if responseBody.Status != "database connected" {
		t.Errorf("Expected status 'database connected' but got '%s'", responseBody.Status)
	}
}
