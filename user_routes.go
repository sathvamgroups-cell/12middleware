package routes

import (
	"12middleware/controllers"
	"12middleware/middleware"

	"github.com/gin-gonic/gin"
)

func UserRoutes(router *gin.Engine) {

	router.POST("/users", controllers.CreateUser)

	router.GET("/users", middleware.JWTMiddleware(), controllers.GetUsers)

	router.GET("/users/:id", controllers.GetUser)

	router.PUT("/users/:id", controllers.UpdateUser)

	router.DELETE(
		"/users/:id",
		middleware.JWTMiddleware(),
		middleware.AdminMiddleware(),
		controllers.DeleteUser,
	)

	router.POST("/login", controllers.LoginUser)
}
