package controllers

import (
	"context"
	"net/http"

	"sispak-dada/config"
	"sispak-dada/models"

	"github.com/gin-gonic/gin"
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