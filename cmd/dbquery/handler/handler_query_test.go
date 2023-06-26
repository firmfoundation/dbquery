package handler

import (
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/firmfoundation/dbquery/config"
	db "github.com/firmfoundation/dbquery/init"
	"github.com/gofiber/fiber/v2"
)

func TestQueryStateHandler(t *testing.T) {
	db.ConnectRedis(&config.DbConfig{})
	app := fiber.New()

	app.Get("/querystate", QueryStateHandler)

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
