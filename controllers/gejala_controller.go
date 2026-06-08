package controllers

import (
	"context"
	"net/http"

	"sispak-dada/config"
	"sispak-dada/models"

	"github.com/gin-gonic/gin"
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