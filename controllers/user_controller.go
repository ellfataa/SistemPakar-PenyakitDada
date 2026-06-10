package controllers

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"sispak-dada/config"
	"sispak-dada/models"
	"sispak-dada/utils"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

func GetUsers(c *gin.Context) {
	rows, err := config.DB.Query(context.Background(), `
		SELECT id_user, nama, username, role
		FROM users
		ORDER BY id_user ASC
	`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Gagal mengambil data user", "error": err.Error()})
		return
	}
	defer rows.Close()

	var users []models.User

	for rows.Next() {
		var user models.User
		err := rows.Scan(&user.IDUser, &user.Nama, &user.Username, &user.Role)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Gagal membaca data user", "error": err.Error()})
			return
		}
		users = append(users, user)
	}

	c.JSON(http.StatusOK, gin.H{"message": "Data user berhasil diambil", "data": users})
}

func GetUserByID(c *gin.Context) {
	idUser := c.Param("id_user")

	var user models.User

	err := config.DB.QueryRow(context.Background(), `
		SELECT id_user, nama, username, role
		FROM users
		WHERE id_user = $1
	`, idUser).Scan(&user.IDUser, &user.Nama, &user.Username, &user.Role)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"message": "User tidak ditemukan"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Gagal mengambil detail user", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Detail user berhasil diambil", "data": user})
}

func CreateUser(c *gin.Context) {
	var request models.UserRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Format request tidak valid", "error": err.Error()})
		return
	}

	if request.Nama == "" || request.Username == "" || request.Password == "" || request.Role == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Nama, username, password, dan role wajib diisi"})
		return
	}

	if request.Role != "admin" && request.Role != "user" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Role harus admin atau user"})
		return
	}

	var count int
	err := config.DB.QueryRow(context.Background(), `
		SELECT COUNT(*) FROM users WHERE username = $1
	`, request.Username).Scan(&count)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Gagal mengecek username", "error": err.Error()})
		return
	}

	if count > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Username sudah digunakan"})
		return
	}

	hashedPassword, err := utils.HashPassword(request.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Gagal mengenkripsi password", "error": err.Error()})
		return
	}

	var user models.User

	err = config.DB.QueryRow(context.Background(), `
		INSERT INTO users (nama, username, password, role)
		VALUES ($1, $2, $3, $4)
		RETURNING id_user, nama, username, role
	`, request.Nama, request.Username, hashedPassword, request.Role).Scan(
		&user.IDUser,
		&user.Nama,
		&user.Username,
		&user.Role,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Gagal menambahkan user", "error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User berhasil ditambahkan", "data": user})
}

func UpdateUser(c *gin.Context) {
	idUser := c.Param("id_user")

	var request models.UpdateUserRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Format request tidak valid", "error": err.Error()})
		return
	}

	if request.Nama == "" || request.Username == "" || request.Role == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Nama, username, dan role wajib diisi"})
		return
	}

	if request.Role != "admin" && request.Role != "user" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Role harus admin atau user"})
		return
	}

	var count int
	err := config.DB.QueryRow(context.Background(), `
		SELECT COUNT(*) FROM users WHERE username = $1 AND id_user != $2
	`, request.Username, idUser).Scan(&count)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Gagal mengecek username", "error": err.Error()})
		return
	}

	if count > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Username sudah digunakan oleh user lain"})
		return
	}

	var user models.User

	if request.Password != "" {
		hashedPassword, err := utils.HashPassword(request.Password)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Gagal mengenkripsi password", "error": err.Error()})
			return
		}

		err = config.DB.QueryRow(context.Background(), `
			UPDATE users
			SET nama = $1, username = $2, password = $3, role = $4
			WHERE id_user = $5
			RETURNING id_user, nama, username, role
		`, request.Nama, request.Username, hashedPassword, request.Role, idUser).Scan(
			&user.IDUser, &user.Nama, &user.Username, &user.Role,
		)
	} else {
		err = config.DB.QueryRow(context.Background(), `
			UPDATE users
			SET nama = $1, username = $2, role = $3
			WHERE id_user = $4
			RETURNING id_user, nama, username, role
		`, request.Nama, request.Username, request.Role, idUser).Scan(
			&user.IDUser, &user.Nama, &user.Username, &user.Role,
		)
	}

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"message": "User tidak ditemukan"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Gagal memperbarui user", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User berhasil diperbarui", "data": user})
}

func DeleteUser(c *gin.Context) {
	idUser := c.Param("id_user")

	currentUserIDValue, exists := c.Get("id_user")
	if exists {
		currentUserID, ok := currentUserIDValue.(int)
		targetUserID, _ := strconv.Atoi(idUser)

		if ok && currentUserID == targetUserID {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Admin tidak bisa menghapus akun yang sedang digunakan"})
			return
		}
	}

	var countRiwayat int
	err := config.DB.QueryRow(context.Background(), `
		SELECT COUNT(*) FROM riwayat_konsultasi WHERE id_user = $1
	`, idUser).Scan(&countRiwayat)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Gagal mengecek riwayat user", "error": err.Error()})
		return
	}

	if countRiwayat > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"message": "User tidak bisa dihapus karena masih memiliki riwayat konsultasi"})
		return
	}

	result, err := config.DB.Exec(context.Background(), `
		DELETE FROM users WHERE id_user = $1
	`, idUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Gagal menghapus user", "error": err.Error()})
		return
	}

	if result.RowsAffected() == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "User tidak ditemukan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User berhasil dihapus"})
}