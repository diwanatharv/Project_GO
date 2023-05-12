package routes

import (
	"github.com/gin-gonic/gin"

	"awesomeProject2/Authentication/controllers"
	"awesomeProject2/Authentication/middleware"
)

func UserRoutes(incoming *gin.Engine) {
	incoming.Use(middleware.Authenticate())
	incoming.GET("/users", controllers.GetUsers())
	incoming.GET("/users/:UserId", controllers.GetUser())
}
