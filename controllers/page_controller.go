package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func ShowLoginPage(c *gin.Context) {
	c.HTML(http.StatusOK, "login.html", gin.H{
		"title": "Login - Sistem Pakar Dada",
	})
}

func ShowRegisterPage(c *gin.Context) {
	c.HTML(http.StatusOK, "register.html", gin.H{
		"title": "Register - Sistem Pakar Dada",
	})
}

// ADMIN
func ShowAdminDashboard(c *gin.Context) {
	c.HTML(http.StatusOK, "admin_dashboard.html", gin.H{
		"title": "Dashboard Admin - Sistem Pakar Dada",
	})
}
// Admin menu Gejala
func ShowAdminGejalaPage(c *gin.Context) {
	c.HTML(http.StatusOK, "admin_gejala_index.html", gin.H{
		"title": "Data Gejala - Sistem Pakar Dada",
	})
}
func ShowAdminCreateGejalaPage(c *gin.Context) {
	c.HTML(http.StatusOK, "admin_gejala_create.html", gin.H{
		"title": "Tambah Gejala - Sistem Pakar Dada",
	})
}
func ShowAdminEditGejalaPage(c *gin.Context) {
	kodeGejala := c.Param("kode_gejala")
	
	c.HTML(http.StatusOK, "admin_gejala_edit.html", gin.H{
		"title":       "Edit Gejala - Sistem Pakar Dada",
		"kode_gejala": kodeGejala,
	})
}
// Admin menu Penyakit
func ShowAdminPenyakitPage(c *gin.Context) {
	c.HTML(http.StatusOK, "admin_penyakit_index.html", gin.H{
		"title": "Data Penyakit - Sistem Pakar Dada",
	})
}
func ShowAdminCreatePenyakitPage(c *gin.Context) {
	c.HTML(http.StatusOK, "admin_penyakit_create.html", gin.H{
		"title": "Tambah Penyakit - Sistem Pakar Dada",
	})
}
func ShowAdminEditPenyakitPage(c *gin.Context) {
	kodePenyakit := c.Param("kode_penyakit")

	c.HTML(http.StatusOK, "admin_penyakit_edit.html", gin.H{
		"title":         "Edit Penyakit - Sistem Pakar Dada",
		"kode_penyakit": kodePenyakit,
	})
}
// Admin menu Relasi
func ShowAdminRelasiPage(c *gin.Context) {
	c.HTML(http.StatusOK, "admin_relasi_index.html", gin.H{
		"title": "Relasi Penyakit dan Gejala - Sistem Pakar Dada",
	})
}

// USER
func ShowUserDashboard(c *gin.Context) {
	c.HTML(http.StatusOK, "user_dashboard.html", gin.H{
		"title": "Dashboard User - Sistem Pakar Dada",
	})
}
// User menu Konsultasi
func ShowUserKonsultasiPage(c *gin.Context) {
	c.HTML(http.StatusOK, "user_konsultasi_index.html", gin.H{
		"title": "Konsultasi - Sistem Pakar Dada",
	})
}
// User menu Hasil Konsultasi
func ShowUserHasilPage(c *gin.Context) {
	idRiwayat := c.Param("id_riwayat")

	c.HTML(http.StatusOK, "user_hasil_detail.html", gin.H{
		"title":      "Hasil Konsultasi - Sistem Pakar Dada",
		"id_riwayat": idRiwayat,
	})
}
// User menu Riwayat Konsultasi
func ShowUserRiwayatPage(c *gin.Context) {
	c.HTML(http.StatusOK, "user_riwayat_index.html", gin.H{
		"title": "Riwayat Konsultasi - Sistem Pakar Dada",
	})
}
// User menu Profile
func ShowUserProfilePage(c *gin.Context) {
	c.HTML(http.StatusOK, "user_profile_index.html", gin.H{
		"title": "Profil User - Sistem Pakar Dada",
	})
}