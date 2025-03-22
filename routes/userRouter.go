package routes

import (
	"github.com/gin-gonic/gin"
	controller "restaurant.com/m/controllers"
)

func UserRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.GET("/users", controller.GetUsers())
	incomingRoutes.GET("/user", controller.GetUser())
	incomingRoutes.POST("/users/login", controller.Login())
	incomingRoutes.POST("/users/:signup", controller.SignUp())

}
