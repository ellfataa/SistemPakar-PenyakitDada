package controllers

import (
	"context"
	"encoding/json"
	"net/http"

	"sispak-dada/config"
	"sispak-dada/models"

	"github.com/gin-gonic/gin"
)

type riwayatTemp struct {
	Riwayat           models.RiwayatKonsultasi
	GejalaDipilihRaw string
}

func GetRiwayatByUser(c *gin.Context) {
	idUser := c.Param("id_user")

	rows, err := config.DB.Query(context.Background(), `
		SELECT 
			r.id_riwayat,
			r.id_user,
			u.nama,
			r.gejala_dipilih::text,
			r.hasil_diagnosa,
			r.probabilitas,
			r.kode_penyakit,
			p.nama_penyakit,
			TO_CHAR(r.waktu_konsultasi, 'YYYY-MM-DD HH24:MI:SS') AS waktu_konsultasi
		FROM riwayat_konsultasi r
		JOIN users u 
			ON r.id_user = u.id_user
		JOIN penyakit p 
			ON r.kode_penyakit = p.kode_penyakit
		WHERE r.id_user = $1
		ORDER BY r.waktu_konsultasi DESC
	`, idUser)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Gagal mengambil riwayat konsultasi",
			"error":   err.Error(),
		})
		return
	}

	var tempList []riwayatTemp

	for rows.Next() {
		var temp riwayatTemp

		err := rows.Scan(
			&temp.Riwayat.IDRiwayat,
			&temp.Riwayat.IDUser,
			&temp.Riwayat.NamaUser,
			&temp.GejalaDipilihRaw,
			&temp.Riwayat.HasilDiagnosa,
			&temp.Riwayat.Probabilitas,
			&temp.Riwayat.KodePenyakit,
			&temp.Riwayat.NamaPenyakit,
			&temp.Riwayat.WaktuKonsultasi,
		)

		if err != nil {
			rows.Close()
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Gagal membaca data riwayat",
				"error":   err.Error(),
			})
			return
		}

		tempList = append(tempList, temp)
	}

	rows.Close()

	if rows.Err() != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Terjadi kesalahan saat membaca data riwayat",
			"error":   rows.Err().Error(),
		})
		return
	}

	var riwayatList []models.RiwayatKonsultasi

	for _, temp := range tempList {
		gejalaList, err := getDetailGejalaFromJSON(temp.GejalaDipilihRaw)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Gagal membaca detail gejala riwayat",
				"error":   err.Error(),
			})
			return
		}

		temp.Riwayat.Gejala = gejalaList
		riwayatList = append(riwayatList, temp.Riwayat)
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Riwayat konsultasi berhasil diambil",
		"data":    riwayatList,
	})
}

func getDetailGejalaFromJSON(gejalaDipilihText string) ([]models.RiwayatGejala, error) {
	var kodeGejalaList []string

	err := json.Unmarshal([]byte(gejalaDipilihText), &kodeGejalaList)
	if err != nil {
		return nil, err
	}

	if len(kodeGejalaList) == 0 {
		return []models.RiwayatGejala{}, nil
	}

	rows, err := config.DB.Query(context.Background(), `
		SELECT kode_gejala, nama_gejala
		FROM gejala
		WHERE kode_gejala = ANY($1)
		ORDER BY kode_gejala ASC
	`, kodeGejalaList)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var gejalaList []models.RiwayatGejala

	for rows.Next() {
		var gejala models.RiwayatGejala

		err := rows.Scan(
			&gejala.KodeGejala,
			&gejala.NamaGejala,
		)

		if err != nil {
			return nil, err
		}

		gejalaList = append(gejalaList, gejala)
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return gejalaList, nil
}