package database

import (
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"PrGoRestApi/config"
)

var DB *gorm.DB

func Init() *gorm.DB {
	dsn := config.GetDSN()
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to DB:", err)
	}

	log.Println("Database connection established")
	DB = db
	return DB
}
