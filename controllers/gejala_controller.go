package controllers

import (
	"context"
	"errors"
	"net/http"

	"sispak-dada/config"
	"sispak-dada/models"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

func GetGejala(c *gin.Context) {
	rows, err := config.DB.Query(context.Background(), `
		SELECT kode_gejala, nama_gejala
		FROM gejala
		ORDER BY kode_gejala ASC
	`)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Gagal mengambil data gejala",
			"error":   err.Error(),
		})
		return
	}

	defer rows.Close()

	var gejalaList []models.Gejala

	for rows.Next() {
		var gejala models.Gejala

		err := rows.Scan(
			&gejala.KodeGejala,
			&gejala.NamaGejala,
		)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Gagal membaca data gejala",
				"error":   err.Error(),
			})
			return
		}

		gejalaList = append(gejalaList, gejala)
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Data gejala berhasil diambil",
		"data":    gejalaList,
	})
}

func GetGejalaByKode(c *gin.Context) {
	kodeGejala := c.Param("kode_gejala")

	var gejala models.Gejala

	err := config.DB.QueryRow(context.Background(), `
		SELECT kode_gejala, nama_gejala
		FROM gejala
		WHERE kode_gejala = $1
	`, kodeGejala).Scan(
		&gejala.KodeGejala,
		&gejala.NamaGejala,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{
				"message": "Data gejala tidak ditemukan",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Gagal mengambil detail gejala",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Detail gejala berhasil diambil",
		"data":    gejala,
	})
}

func CreateGejala(c *gin.Context) {
	var request models.GejalaRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Format request tidak valid",
			"error":   err.Error(),
		})
		return
	}

	if request.KodeGejala == "" || request.NamaGejala == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "kode_gejala dan nama_gejala wajib diisi",
		})
		return
	}

	var count int
	err := config.DB.QueryRow(context.Background(), `
		SELECT COUNT(*)
		FROM gejala
		WHERE kode_gejala = $1
	`, request.KodeGejala).Scan(&count)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Gagal mengecek kode gejala",
			"error":   err.Error(),
		})
		return
	}

	if count > 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Kode gejala sudah digunakan",
		})
		return
	}

	var gejala models.Gejala

	err = config.DB.QueryRow(context.Background(), `
		INSERT INTO gejala (kode_gejala, nama_gejala)
		VALUES ($1, $2)
		RETURNING kode_gejala, nama_gejala
	`,
		request.KodeGejala,
		request.NamaGejala,
	).Scan(
		&gejala.KodeGejala,
		&gejala.NamaGejala,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Gagal menambahkan data gejala",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Data gejala berhasil ditambahkan",
		"data":    gejala,
	})
}

func UpdateGejala(c *gin.Context) {
	kodeGejala := c.Param("kode_gejala")

	var request models.GejalaRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Format request tidak valid",
			"error":   err.Error(),
		})
		return
	}

	if request.NamaGejala == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "nama_gejala wajib diisi",
		})
		return
	}

	var gejala models.Gejala

	err := config.DB.QueryRow(context.Background(), `
		UPDATE gejala
		SET nama_gejala = $1
		WHERE kode_gejala = $2
		RETURNING kode_gejala, nama_gejala
	`,
		request.NamaGejala,
		kodeGejala,
	).Scan(
		&gejala.KodeGejala,
		&gejala.NamaGejala,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{
				"message": "Data gejala tidak ditemukan",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Gagal memperbarui data gejala",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Data gejala berhasil diperbarui",
		"data":    gejala,
	})
}

func DeleteGejala(c *gin.Context) {
	kodeGejala := c.Param("kode_gejala")

	var countRelasi int
	err := config.DB.QueryRow(context.Background(), `
		SELECT COUNT(*)
		FROM penyakit_gejala
		WHERE kode_gejala = $1
	`, kodeGejala).Scan(&countRelasi)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Gagal mengecek relasi gejala",
			"error":   err.Error(),
		})
		return
	}

	if countRelasi > 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Gejala tidak bisa dihapus karena masih digunakan pada relasi penyakit_gejala",
		})
		return
	}

	result, err := config.DB.Exec(context.Background(), `
		DELETE FROM gejala
		WHERE kode_gejala = $1
	`, kodeGejala)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Gagal menghapus data gejala",
			"error":   err.Error(),
		})
		return
	}

	if result.RowsAffected() == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Data gejala tidak ditemukan",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Data gejala berhasil dihapus",
	})
}