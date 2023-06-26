package handler

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/firmfoundation/dbquery/config"
	db "github.com/firmfoundation/dbquery/init"
	"github.com/gofiber/fiber/v2"
)

func TestQueryStateHandler(t *testing.T) {
	requestBody := map[string]string{
		"DBHost":         "127.0.0.1",
		"DBUserName":     "postgres",
		"DBUserPassword": "password123",
		"DBName":         "qrmenu",
		"DBPort":         "5432",
	}
	db.ConnectRedis(&config.DbConfig{})
	app := fiber.New()

	app.Get("/querystate", QueryStateHandler)

	t.Run("success", func(t *testing.T) {

		//lets get valid database instance id
		requestBodyBytes, _ := json.Marshal(requestBody)

		// Create a new test request
		req := httptest.NewRequest(http.MethodPost, "/connect", bytes.NewReader(requestBodyBytes))
		req.Header.Set("Content-Type", "application/json") // Set the JSON header

		// Call the handler function with the test request and response
		app.Post("/connect", CreateConnectionHandler)
		// Simulate the request and capture the response
		res, err := app.Test(req, -1)
		if err != nil {
			t.Errorf("Error testing request: %v", err)
		}

		// Check the connection api response status code
		if res.StatusCode != http.StatusAccepted {
			t.Errorf("Expected status code %d but got %d", http.StatusAccepted, res.StatusCode)
		}

		// Check the response body status
		var responseBody struct {
			Status       string `json:"status"`
			DBInstanceId string `json:"database_instance_id"`
		}
		bodyBytes, err := ioutil.ReadAll(res.Body)
		err = json.Unmarshal(bodyBytes, &responseBody)
		if err != nil {
			t.Errorf("Error decoding response body: %v", err)
		}
		if responseBody.Status != "database connected" {
			t.Errorf("Expected status 'database connected' but got '%s'", responseBody.Status)
		}

		//test api is fetching from postgrSQL
		req = httptest.NewRequest("GET", "/querystate?page=1&page_size=50&sort=fastest&filter_query=select&database_instance_id="+responseBody.DBInstanceId, nil)
		resp, err := app.Test(req, -1)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != fiber.StatusAccepted {
			t.Errorf("expected status code %d but got %d", fiber.StatusOK, resp.StatusCode)
		}

		//test api is fetching from cache redis
		req = httptest.NewRequest("GET", "/querystate?page=1&page_size=50&sort=fastest&filter_query=select&database_instance_id="+responseBody.DBInstanceId, nil)
		resp, err = app.Test(req, -1)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != fiber.StatusAccepted {
			t.Errorf("expected status code %d but got %d", fiber.StatusOK, resp.StatusCode)
		}
	})

	t.Run("instance not connect", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/querystate?page=1&page_size=50&sort=fastest&filter_query=select&database_instance_id=invalid", nil)
		resp, err := app.Test(req, -1)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != fiber.StatusInternalServerError {
			t.Errorf("expected status code %d but got %d", fiber.StatusInternalServerError, resp.StatusCode)
		}
	})

	t.Run("invalid query params", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/querystate?page=invalid&page_size=50&sort=fastest&filter_query=invalid&database_instance_id=invalid", nil)
		resp, err := app.Test(req, -1)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != fiber.StatusBadRequest {
			t.Errorf("expected status code %d but got %d", fiber.StatusBadRequest, resp.StatusCode)
		}

		var data map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if data["error"] != "invalid page number" {
			t.Errorf("expected error message 'invalid page number' but got %v", data["error"])
		}
	})
}
