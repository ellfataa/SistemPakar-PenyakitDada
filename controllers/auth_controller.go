package controllers

import (
	"context"
	"net/http"

	"sispak-dada/config"
	"sispak-dada/models"
	"sispak-dada/utils"

	"github.com/gin-gonic/gin"
)

func Register(c *gin.Context) {
	var request models.RegisterRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Format request tidak valid",
			"error":   err.Error(),
		})
		return
	}

	if request.Nama == "" || request.Username == "" || request.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Nama, username, dan password wajib diisi",
		})
		return
	}

	if request.Role == "" {
		request.Role = "user"
	}

	var count int
	err := config.DB.QueryRow(context.Background(), `
		SELECT COUNT(*) 
		FROM users 
		WHERE username = $1
	`, request.Username).Scan(&count)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Gagal mengecek username",
			"error":   err.Error(),
		})
		return
	}

	if count > 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Username sudah digunakan",
		})
		return
	}

	hashedPassword, err := utils.HashPassword(request.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Gagal mengenkripsi password",
			"error":   err.Error(),
		})
		return
	}

	var user models.User

	err = config.DB.QueryRow(context.Background(), `
		INSERT INTO users (nama, username, password, role)
		VALUES ($1, $2, $3, $4)
		RETURNING id_user, nama, username, role
	`,
		request.Nama,
		request.Username,
		hashedPassword,
		request.Role,
	).Scan(
		&user.IDUser,
		&user.Nama,
		&user.Username,
		&user.Role,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Gagal melakukan register",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Register berhasil",
		"data":    user,
	})
}

func Login(c *gin.Context) {
	var request models.LoginRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Format request tidak valid",
			"error":   err.Error(),
		})
		return
	}

	if request.Username == "" || request.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Username dan password wajib diisi",
		})
		return
	}

	var user models.User

	err := config.DB.QueryRow(context.Background(), `
		SELECT id_user, nama, username, password, role
		FROM users
		WHERE username = $1
	`,
		request.Username,
	).Scan(
		&user.IDUser,
		&user.Nama,
		&user.Username,
		&user.Password,
		&user.Role,
	)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "Username atau password salah",
		})
		return
	}

	if !utils.CheckPasswordHash(request.Password, user.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "Username atau password salah",
		})
		return
	}

	token, err := utils.GenerateToken(
		user.IDUser,
		user.Nama,
		user.Username,
		user.Role,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Gagal membuat token",
			"error":   err.Error(),
		})
		return
	}

	user.Password = ""

	c.JSON(http.StatusOK, gin.H{
		"message": "Login berhasil",
		"token":   token,
		"data":    user,
	})
}

func Profile(c *gin.Context) {
	idUser, _ := c.Get("id_user")
	nama, _ := c.Get("nama")
	username, _ := c.Get("username")
	role, _ := c.Get("role")

	c.JSON(http.StatusOK, gin.H{
		"message": "Data profile berhasil diambil",
		"data": gin.H{
			"id_user":  idUser,
			"nama":     nama,
			"username": username,
			"role":     role,
		},
	})
}