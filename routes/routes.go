package routes

import (
	"sispak-dada/controllers"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	r.GET("/", controllers.Home)
	r.GET("/gejala", controllers.GetGejala)
}