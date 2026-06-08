package controllers

import (
	"context"
	"net/http"

	"sispak-dada/config"
	"sispak-dada/models"

	"github.com/gin-gonic/gin"
)

func GetRiwayatByUser(c *gin.Context) {
	idUser := c.Param("id_user")

	rows, err := config.DB.Query(context.Background(), `
		SELECT 
			id_riwayat,
			id_user,
			gejala_dipilih::text,
			hasil_diagnosa,
			probabilitas,
			kode_penyakit,
			TO_CHAR(waktu_konsultasi, 'YYYY-MM-DD HH24:MI:SS') AS waktu_konsultasi
		FROM riwayat_konsultasi
		WHERE id_user = $1
		ORDER BY waktu_konsultasi DESC
	`, idUser)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Gagal mengambil riwayat konsultasi",
			"error":   err.Error(),
		})
		return
	}

	defer rows.Close()

	var riwayatList []models.RiwayatKonsultasi

	for rows.Next() {
		var riwayat models.RiwayatKonsultasi

		err := rows.Scan(
			&riwayat.IDRiwayat,
			&riwayat.IDUser,
			&riwayat.GejalaDipilih,
			&riwayat.HasilDiagnosa,
			&riwayat.Probabilitas,
			&riwayat.KodePenyakit,
			&riwayat.WaktuKonsultasi,
		)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Gagal membaca data riwayat",
				"error":   err.Error(),
			})
			return
		}

		riwayatList = append(riwayatList, riwayat)
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Riwayat konsultasi berhasil diambil",
		"data":    riwayatList,
	})
}