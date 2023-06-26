package handler

import (
	"github.com/firmfoundation/dbquery/config"
	db "github.com/firmfoundation/dbquery/init"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

var DbInstanceMap map[string]*gorm.DB = make(map[string]*gorm.DB)

func CreateConnectionHandler(c *fiber.Ctx) error {
	config := new(config.DbConfig)

	err := c.BodyParser(config)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err})
	}

	if ok, dbInstance := db.ConnectDB(config); ok {
		//multiple instance support
		cacheKey := uuid.New()
		DbInstanceMap[cacheKey.String()] = dbInstance
		return c.Status(fiber.StatusAccepted).JSON(fiber.Map{"status": "database connected", "database_instance_id": cacheKey})
	}

	return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "database not connected"})
}
