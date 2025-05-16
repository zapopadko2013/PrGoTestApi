package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"

	_ "PrGoRestApi/docs"

	"PrGoRestApi/internal/database"
	"PrGoRestApi/internal/migration"
	"PrGoRestApi/router"
)

// @title           API по людям с обогащением
// @version         1.0
// @description     API для обогащения информации о человеке по имени (возраст, пол, национальность)
// @host            localhost:8081
// @BasePath        /
func main() {
	// Загрузка .env
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	// Инициализация базы
	db := database.Init()

	// Миграции
	migration.RunMigrations()

	// Запуск API
	r := router.SetupRouter(db)
	if err := r.Run(":" + os.Getenv("PORT")); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
