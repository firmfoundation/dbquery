package router

import (
	"github.com/firmfoundation/dbquery/cmd/dbquery/handler"
	"github.com/gofiber/fiber/v2"
)

func Routes(app *fiber.App) {

	//group
	queryApi := app.Group("/connect")
	queryState := app.Group("/querystate")

	//routes
	queryApi.Post("/", handler.CreateConnectionHandler)
	queryState.Get("/", handler.QueryStateHandler)
}
