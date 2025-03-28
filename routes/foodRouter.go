package routes

import (
	"github.com/gin-gonic/gin"
	"restaurant.com/m/controllers"
)

func FoodRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.GET("/foods", controllers.GetFoods())
	incomingRoutes.GET("foods/:food_id", controllers.GetFood())
	incomingRoutes.POST("foods", controllers.CreateFood())
	incomingRoutes.PATCH("foods/:food_id", controllers.UpdateFood())
}
