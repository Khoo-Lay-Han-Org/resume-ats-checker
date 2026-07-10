package database

import (
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	systemconfig "resuming/system-config"
)

var DB *gorm.DB

var err error

func DatabaseConnect() error {
	DB, err = gorm.Open(postgres.Open(systemconfig.DatabaseDSN), &gorm.Config{})
	if err != nil {
		log.Println("Failed to Connect to Database")
		log.Println(err)
		return err
	}
	log.Println("Successfully Connected to Database")
	return nil
}
