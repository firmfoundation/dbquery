package handler

import (
	"fmt"

	"github.com/firmfoundation/dbquery/config"
	db "github.com/firmfoundation/dbquery/init"
	"github.com/gofiber/fiber/v2"
)

func CreateConnectionHandler(c *fiber.Ctx) error {
	config := new(config.DbConfig)

	err := c.BodyParser(config)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "data": err})
	}

	fmt.Println(db.ConnectDB(config))
	//fmt.Println(config)
	return c.Status(200).JSON(fiber.Map{"status": "Happy"})
}
