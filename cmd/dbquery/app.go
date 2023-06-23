package main

import (
	"github.com/firmfoundation/dbquery/cmd/dbquery/router"
	"github.com/firmfoundation/dbquery/config"
	db "github.com/firmfoundation/dbquery/init"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	db.ConnectRedis(&config.DbConfig{})

	app := fiber.New()
	app.Use(logger.New())
	app.Use(cors.New())

	router.Routes(app)

	app.Listen(":3000")
}
