package handler

import (
	"github.com/firmfoundation/dbquery/config"
	db "github.com/firmfoundation/dbquery/init"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func CreateConnectionHandler(c *fiber.Ctx) error {
	config := new(config.DbConfig)

	err := c.BodyParser(config)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "data": err})
	}

	if ok := db.ConnectDB(config); ok {
		/*

			to do
			the query statistics api need to have database instance id
			to be used for clients to access multiple databases

		*/
		uid := uuid.New()
		return c.Status(200).JSON(fiber.Map{"status": "database connected", "database_instance_id": uid})
	}

	return c.Status(500).JSON(fiber.Map{"status": "database not connected"})
}
