package routes

import (
	"sispak-dada/controllers"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	r.GET("/", controllers.Home)

	r.POST("/register", controllers.Register)
	r.POST("/login", controllers.Login)

	r.GET("/gejala", controllers.GetGejala)
	r.GET("/penyakit", controllers.GetPenyakit)

	r.POST("/konsultasi", controllers.Konsultasi)
	r.GET("/hasil-konsultasi/:id_riwayat", controllers.GetHasilKonsultasi)

	r.GET("/riwayat/:id_user", controllers.GetRiwayatByUser)
}