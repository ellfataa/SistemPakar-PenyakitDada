package routes

import (
	"sispak-dada/controllers"
	"sispak-dada/middlewares"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	r.GET("/", controllers.Home)

	r.POST("/register", controllers.Register)
	r.POST("/login", controllers.Login)

	// Public routes
	r.GET("/gejala", controllers.GetGejala)
	r.GET("/gejala/:kode_gejala", controllers.GetGejalaByKode)

	r.GET("/penyakit", controllers.GetPenyakit)
	r.GET("/penyakit/:kode_penyakit", controllers.GetPenyakitByKode)

	r.GET("/relasi", controllers.GetRelasi)
	r.GET("/relasi/penyakit/:kode_penyakit", controllers.GetRelasiByPenyakit)

	// User routes wajib login
	userRoutes := r.Group("/")
	userRoutes.Use(middlewares.AuthMiddleware())
	{
		userRoutes.GET("/profile", controllers.Profile)

		userRoutes.POST("/konsultasi", controllers.Konsultasi)
		userRoutes.GET("/hasil-konsultasi/:id_riwayat", controllers.GetHasilKonsultasi)
		userRoutes.GET("/riwayat/:id_user", controllers.GetRiwayatByUser)
	}

	// Admin dan Pakar routes
	adminRoutes := r.Group("/")
	adminRoutes.Use(middlewares.AuthMiddleware())
	adminRoutes.Use(middlewares.RoleMiddleware("admin", "pakar"))
	{
		adminRoutes.POST("/gejala", controllers.CreateGejala)
		adminRoutes.PUT("/gejala/:kode_gejala", controllers.UpdateGejala)
		adminRoutes.DELETE("/gejala/:kode_gejala", controllers.DeleteGejala)

		adminRoutes.POST("/penyakit", controllers.CreatePenyakit)
		adminRoutes.PUT("/penyakit/:kode_penyakit", controllers.UpdatePenyakit)
		adminRoutes.DELETE("/penyakit/:kode_penyakit", controllers.DeletePenyakit)

		adminRoutes.POST("/relasi", controllers.CreateRelasi)
		adminRoutes.DELETE("/relasi/:kode_penyakit/:kode_gejala", controllers.DeleteRelasi)
	}
}