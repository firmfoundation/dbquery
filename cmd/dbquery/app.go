package main

import (
	"github.com/firmfoundation/dbquery/cmd/dbquery/router"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	app := fiber.New()
	app.Use(logger.New())
	app.Use(cors.New())

	router.Routes(app)

	app.Listen(":3000")
}
