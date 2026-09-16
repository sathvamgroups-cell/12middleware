package main

import (
	"log"

	"12middleware/config"
	"12middleware/middleware"
	"12middleware/models"
	"12middleware/routes"

	"github.com/gin-gonic/gin"
)

func main() {

	models.DB = config.ConnectDatabase()

	models.DB.AutoMigrate(&models.User{}, &models.Post{})

	router := gin.Default()

	router.Use(middleware.LoggerMiddleware())

	routes.UserRoutes(router)

	log.Println("Server running on port 8080")

	router.Run(":8080")
}
