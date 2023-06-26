package init

import (
	"fmt"
	"log"

	"github.com/firmfoundation/dbquery/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB(config *config.DbConfig) (bool, *gorm.DB) {
	var err error
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=CET", config.DBHost, config.DBUserName, config.DBUserPassword, config.DBName, config.DBPort)

	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Println("Failed to connect to the Database")
		return false, nil
	}

	fmt.Println("? Connected Successfully to the Database")
	return true, DB
}
