package routes

import (
	"github.com/gin-gonic/gin"
	controller "restaurant.com/m/restaurant_management/controllers"
)

func UserRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.GET("/users", controller.GetUsers())
	incomingRoutes.GET("/users", controller.GetUser())
	incomingRoutes.POST("/users/login", controller.Login())
	incomingRoutes.POST("/users/:signup", controller.SignUp())

}
