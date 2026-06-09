package controllers

import (
	"context"
	"net/http"

	"sispak-dada/config"
	"sispak-dada/models"

	"github.com/gin-gonic/gin"
)

func GetRelasi(c *gin.Context) {
	rows, err := config.DB.Query(context.Background(), `
		SELECT 
			p.kode_penyakit,
			p.nama_penyakit,
			g.kode_gejala,
			g.nama_gejala
		FROM penyakit_gejala pg
		JOIN penyakit p 
			ON pg.kode_penyakit = p.kode_penyakit
		JOIN gejala g 
			ON pg.kode_gejala = g.kode_gejala
		ORDER BY p.kode_penyakit ASC, g.kode_gejala ASC
	`)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Gagal mengambil data relasi penyakit dan gejala",
			"error":   err.Error(),
		})
		return
	}

	defer rows.Close()

	var relasiList []models.RelasiPenyakitGejala

	for rows.Next() {
		var relasi models.RelasiPenyakitGejala

		err := rows.Scan(
			&relasi.KodePenyakit,
			&relasi.NamaPenyakit,
			&relasi.KodeGejala,
			&relasi.NamaGejala,
		)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Gagal membaca data relasi",
				"error":   err.Error(),
			})
			return
		}

		relasiList = append(relasiList, relasi)
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Data relasi berhasil diambil",
		"data":    relasiList,
	})
}

func GetRelasiByPenyakit(c *gin.Context) {
	kodePenyakit := c.Param("kode_penyakit")

	rows, err := config.DB.Query(context.Background(), `
		SELECT 
			p.kode_penyakit,
			p.nama_penyakit,
			g.kode_gejala,
			g.nama_gejala
		FROM penyakit_gejala pg
		JOIN penyakit p 
			ON pg.kode_penyakit = p.kode_penyakit
		JOIN gejala g 
			ON pg.kode_gejala = g.kode_gejala
		WHERE p.kode_penyakit = $1
		ORDER BY g.kode_gejala ASC
	`, kodePenyakit)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Gagal mengambil relasi berdasarkan penyakit",
			"error":   err.Error(),
		})
		return
	}

	defer rows.Close()

	var relasiList []models.RelasiPenyakitGejala

	for rows.Next() {
		var relasi models.RelasiPenyakitGejala

		err := rows.Scan(
			&relasi.KodePenyakit,
			&relasi.NamaPenyakit,
			&relasi.KodeGejala,
			&relasi.NamaGejala,
		)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Gagal membaca data relasi",
				"error":   err.Error(),
			})
			return
		}

		relasiList = append(relasiList, relasi)
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Data relasi berdasarkan penyakit berhasil diambil",
		"data":    relasiList,
	})
}

func CreateRelasi(c *gin.Context) {
	var request models.RelasiRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Format request tidak valid",
			"error":   err.Error(),
		})
		return
	}

	if request.KodePenyakit == "" || request.KodeGejala == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "kode_penyakit dan kode_gejala wajib diisi",
		})
		return
	}

	var countPenyakit int
	err := config.DB.QueryRow(context.Background(), `
		SELECT COUNT(*) 
		FROM penyakit 
		WHERE kode_penyakit = $1
	`, request.KodePenyakit).Scan(&countPenyakit)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Gagal mengecek data penyakit",
			"error":   err.Error(),
		})
		return
	}

	if countPenyakit == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Kode penyakit tidak ditemukan",
		})
		return
	}

	var countGejala int
	err = config.DB.QueryRow(context.Background(), `
		SELECT COUNT(*) 
		FROM gejala 
		WHERE kode_gejala = $1
	`, request.KodeGejala).Scan(&countGejala)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Gagal mengecek data gejala",
			"error":   err.Error(),
		})
		return
	}

	if countGejala == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Kode gejala tidak ditemukan",
		})
		return
	}

	var countRelasi int
	err = config.DB.QueryRow(context.Background(), `
		SELECT COUNT(*) 
		FROM penyakit_gejala 
		WHERE kode_penyakit = $1 AND kode_gejala = $2
	`, request.KodePenyakit, request.KodeGejala).Scan(&countRelasi)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Gagal mengecek data relasi",
			"error":   err.Error(),
		})
		return
	}

	if countRelasi > 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Relasi penyakit dan gejala sudah ada",
		})
		return
	}

	_, err = config.DB.Exec(context.Background(), `
		INSERT INTO penyakit_gejala (kode_penyakit, kode_gejala)
		VALUES ($1, $2)
	`, request.KodePenyakit, request.KodeGejala)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Gagal menambahkan relasi penyakit dan gejala",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Relasi penyakit dan gejala berhasil ditambahkan",
		"data": gin.H{
			"kode_penyakit": request.KodePenyakit,
			"kode_gejala":   request.KodeGejala,
		},
	})
}

func DeleteRelasi(c *gin.Context) {
	kodePenyakit := c.Param("kode_penyakit")
	kodeGejala := c.Param("kode_gejala")

	result, err := config.DB.Exec(context.Background(), `
		DELETE FROM penyakit_gejala
		WHERE kode_penyakit = $1 AND kode_gejala = $2
	`, kodePenyakit, kodeGejala)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Gagal menghapus relasi penyakit dan gejala",
			"error":   err.Error(),
		})
		return
	}

	if result.RowsAffected() == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Relasi penyakit dan gejala tidak ditemukan",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Relasi penyakit dan gejala berhasil dihapus",
	})
}