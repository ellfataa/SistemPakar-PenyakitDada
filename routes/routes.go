package routes

import (
	"net/http"
	"sispak-dada/controllers"
	"sispak-dada/middlewares"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "landing.html", gin.H{
			"title": "Sistem Pakar Penyakit Dada",
		})
	})

	// Frontend
	r.GET("/login-page", controllers.ShowLoginPage)
	r.GET("/register-page", controllers.ShowRegisterPage)
	// ADMIN
	r.GET("/admin/dashboard", controllers.ShowAdminDashboard)
	// menu Gejala
	r.GET("/admin/gejala", controllers.ShowAdminGejalaPage)
	r.GET("/admin/gejala/create", controllers.ShowAdminCreateGejalaPage)
	r.GET("/admin/gejala/edit/:kode_gejala", controllers.ShowAdminEditGejalaPage)
	// menu Penyakit
	r.GET("/admin/penyakit", controllers.ShowAdminPenyakitPage)
	r.GET("/admin/penyakit/create", controllers.ShowAdminCreatePenyakitPage)
	r.GET("/admin/penyakit/edit/:kode_penyakit", controllers.ShowAdminEditPenyakitPage)
	// menu Relasi
	r.GET("/admin/relasi", controllers.ShowAdminRelasiPage)
	// menu Profile Admin
	r.GET("/admin/profile", controllers.ShowAdminProfilePage)
	r.GET("/admin/profile/edit", controllers.ShowAdminEditProfilePage)
	// menu Manajemen Users
	r.GET("/admin/users", controllers.ShowAdminUsersPage)
	r.GET("/admin/users/create", controllers.ShowAdminCreateUserPage)
	r.GET("/admin/users/edit/:id_user", controllers.ShowAdminEditUserPage)

	// USER
	r.GET("/user/dashboard", controllers.ShowUserDashboard)
	r.GET("/user/konsultasi", controllers.ShowUserKonsultasiPage)
	r.GET("/user/hasil/:id_riwayat", controllers.ShowUserHasilPage)
	r.GET("/user/riwayat", controllers.ShowUserRiwayatPage)
	r.GET("/user/profile", controllers.ShowUserProfilePage)
	r.GET("/user/profile/edit", controllers.ShowUserEditProfilePage)


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
		userRoutes.PUT("/profile", controllers.UpdateProfile)

		userRoutes.POST("/konsultasi", controllers.Konsultasi)
		userRoutes.GET("/hasil-konsultasi/:id_riwayat", controllers.GetHasilKonsultasi)
		userRoutes.GET("/riwayat/:id_user", controllers.GetRiwayatByUser)
	}

	// Admin routes
	adminRoutes := r.Group("/")
	adminRoutes.Use(middlewares.AuthMiddleware())
	adminRoutes.Use(middlewares.RoleMiddleware("admin"))
	{
		adminRoutes.GET("/users", controllers.GetUsers)
		adminRoutes.GET("/users/:id_user", controllers.GetUserByID)
		adminRoutes.POST("/users", controllers.CreateUser)
		adminRoutes.PUT("/users/:id_user", controllers.UpdateUser)
		adminRoutes.DELETE("/users/:id_user", controllers.DeleteUser)

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