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

func GetPenyakit(c *gin.Context) {
	rows, err := config.DB.Query(context.Background(), `
		SELECT kode_penyakit, nama_penyakit, deskripsi, solusi
		FROM penyakit
		ORDER BY kode_penyakit ASC
	`)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Gagal mengambil data penyakit",
			"error":   err.Error(),
		})
		return
	}

	defer rows.Close()

	var penyakitList []models.Penyakit

	for rows.Next() {
		var penyakit models.Penyakit

		err := rows.Scan(
			&penyakit.KodePenyakit,
			&penyakit.NamaPenyakit,
			&penyakit.Deskripsi,
			&penyakit.Solusi,
		)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Gagal membaca data penyakit",
				"error":   err.Error(),
			})
			return
		}

		penyakitList = append(penyakitList, penyakit)
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Data penyakit berhasil diambil",
		"data":    penyakitList,
	})
}

func GetPenyakitByKode(c *gin.Context) {
	kodePenyakit := c.Param("kode_penyakit")

	var penyakit models.Penyakit

	err := config.DB.QueryRow(context.Background(), `
		SELECT kode_penyakit, nama_penyakit, deskripsi, solusi
		FROM penyakit
		WHERE kode_penyakit = $1
	`, kodePenyakit).Scan(
		&penyakit.KodePenyakit,
		&penyakit.NamaPenyakit,
		&penyakit.Deskripsi,
		&penyakit.Solusi,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{
				"message": "Data penyakit tidak ditemukan",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Gagal mengambil detail penyakit",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Detail penyakit berhasil diambil",
		"data":    penyakit,
	})
}

func CreatePenyakit(c *gin.Context) {
	var request models.PenyakitRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Format request tidak valid",
			"error":   err.Error(),
		})
		return
	}

	if request.KodePenyakit == "" || request.NamaPenyakit == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "kode_penyakit dan nama_penyakit wajib diisi",
		})
		return
	}

	var count int
	err := config.DB.QueryRow(context.Background(), `
		SELECT COUNT(*)
		FROM penyakit
		WHERE kode_penyakit = $1
	`, request.KodePenyakit).Scan(&count)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Gagal mengecek kode penyakit",
			"error":   err.Error(),
		})
		return
	}

	if count > 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Kode penyakit sudah digunakan",
		})
		return
	}

	var penyakit models.Penyakit

	err = config.DB.QueryRow(context.Background(), `
		INSERT INTO penyakit 
			(kode_penyakit, nama_penyakit, deskripsi, solusi)
		VALUES 
			($1, $2, $3, $4)
		RETURNING kode_penyakit, nama_penyakit, deskripsi, solusi
	`,
		request.KodePenyakit,
		request.NamaPenyakit,
		request.Deskripsi,
		request.Solusi,
	).Scan(
		&penyakit.KodePenyakit,
		&penyakit.NamaPenyakit,
		&penyakit.Deskripsi,
		&penyakit.Solusi,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Gagal menambahkan data penyakit",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Data penyakit berhasil ditambahkan",
		"data":    penyakit,
	})
}

func UpdatePenyakit(c *gin.Context) {
	kodePenyakit := c.Param("kode_penyakit")

	var request models.PenyakitRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Format request tidak valid",
			"error":   err.Error(),
		})
		return
	}

	if request.NamaPenyakit == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "nama_penyakit wajib diisi",
		})
		return
	}

	var penyakit models.Penyakit

	err := config.DB.QueryRow(context.Background(), `
		UPDATE penyakit
		SET 
			nama_penyakit = $1,
			deskripsi = $2,
			solusi = $3
		WHERE kode_penyakit = $4
		RETURNING kode_penyakit, nama_penyakit, deskripsi, solusi
	`,
		request.NamaPenyakit,
		request.Deskripsi,
		request.Solusi,
		kodePenyakit,
	).Scan(
		&penyakit.KodePenyakit,
		&penyakit.NamaPenyakit,
		&penyakit.Deskripsi,
		&penyakit.Solusi,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{
				"message": "Data penyakit tidak ditemukan",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Gagal memperbarui data penyakit",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Data penyakit berhasil diperbarui",
		"data":    penyakit,
	})
}

func DeletePenyakit(c *gin.Context) {
	kodePenyakit := c.Param("kode_penyakit")

	var countRelasi int
	err := config.DB.QueryRow(context.Background(), `
		SELECT COUNT(*)
		FROM penyakit_gejala
		WHERE kode_penyakit = $1
	`, kodePenyakit).Scan(&countRelasi)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Gagal mengecek relasi penyakit",
			"error":   err.Error(),
		})
		return
	}

	if countRelasi > 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Penyakit tidak bisa dihapus karena masih digunakan pada relasi penyakit_gejala",
		})
		return
	}

	var countRiwayat int
	err = config.DB.QueryRow(context.Background(), `
		SELECT COUNT(*)
		FROM riwayat_konsultasi
		WHERE kode_penyakit = $1
	`, kodePenyakit).Scan(&countRiwayat)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Gagal mengecek riwayat konsultasi",
			"error":   err.Error(),
		})
		return
	}

	if countRiwayat > 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Penyakit tidak bisa dihapus karena masih digunakan pada riwayat konsultasi",
		})
		return
	}

	result, err := config.DB.Exec(context.Background(), `
		DELETE FROM penyakit
		WHERE kode_penyakit = $1
	`, kodePenyakit)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Gagal menghapus data penyakit",
			"error":   err.Error(),
		})
		return
	}

	if result.RowsAffected() == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Data penyakit tidak ditemukan",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Data penyakit berhasil dihapus",
	})
}