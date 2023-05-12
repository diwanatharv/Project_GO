package routes

import (
	"github.com/gin-gonic/gin"

	"awesomeProject2/Authentication/controllers"
)

func AuthRoutes(incoming *gin.Engine) {
	incoming.POST("user/signup", controllers.SingnUp())
	incoming.POST("user/login", controllers.Login())
}
