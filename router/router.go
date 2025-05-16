package router

import (
	"PrGoRestApi/controllers"

	"github.com/gin-gonic/gin"

	"gorm.io/gorm"

	_ "PrGoRestApi/docs"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Person API
// @version 1.0
// @description API для работы с людьми
// @host localhost:8081
// @BasePath /

func SetupRouter(db *gorm.DB) *gin.Engine {
	r := gin.Default()

	api := r.Group("/api")
	{
		api.POST("/persons", controllers.CreatePerson(db))
		api.GET("/persons", controllers.GetPersons(db))
		api.GET("/persons/:id", controllers.GetPersonByID(db))
		api.PUT("/persons/:id", controllers.UpdatePerson(db))
		api.DELETE("/persons/:id", controllers.DeletePerson(db))
	}

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return r
}
