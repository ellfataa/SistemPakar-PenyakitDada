package controllers

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"sispak-dada/config"
	"sispak-dada/models"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

func Konsultasi(c *gin.Context) {
	var request models.KonsultasiRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Format request tidak valid",
			"error":   err.Error(),
		})
		return
	}

	if request.IDUser == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "id_user wajib diisi",
		})
		return
	}

	if len(request.Gejala) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Minimal pilih satu gejala",
		})
		return
	}

	result, err := prosesDiagnosa(request.Gejala)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{
				"message": "Tidak ditemukan penyakit yang sesuai dengan gejala yang dipilih",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Gagal melakukan konsultasi",
			"error":   err.Error(),
		})
		return
	}

	idRiwayat, err := simpanRiwayatKonsultasi(request, result)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Konsultasi berhasil dihitung, tetapi gagal menyimpan riwayat",
			"error":   err.Error(),
			"data":    result,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":    "Konsultasi berhasil",
		"id_riwayat": idRiwayat,
		"data":       result,
	})
}

func prosesDiagnosa(gejala []string) (models.KonsultasiResult, error) {
	query := `
		SELECT 
			p.kode_penyakit,
			p.nama_penyakit,
			p.deskripsi,
			p.solusi,
			COUNT(pg.kode_gejala) FILTER (WHERE pg.kode_gejala = ANY($1)) AS jumlah_cocok,
			COUNT(pg.kode_gejala) AS total_gejala,
			ROUND(
				(
					COUNT(pg.kode_gejala) FILTER (WHERE pg.kode_gejala = ANY($1))::numeric 
					/ COUNT(pg.kode_gejala)::numeric
				) * 100, 
				2
			) AS probabilitas
		FROM penyakit p
		JOIN penyakit_gejala pg 
			ON p.kode_penyakit = pg.kode_penyakit
		GROUP BY 
			p.kode_penyakit,
			p.nama_penyakit,
			p.deskripsi,
			p.solusi
		HAVING COUNT(pg.kode_gejala) FILTER (WHERE pg.kode_gejala = ANY($1)) > 0
		ORDER BY 
			jumlah_cocok DESC,
			probabilitas DESC
		LIMIT 1
	`

	var result models.KonsultasiResult

	err := config.DB.QueryRow(
		context.Background(),
		query,
		gejala,
	).Scan(
		&result.KodePenyakit,
		&result.NamaPenyakit,
		&result.Deskripsi,
		&result.Solusi,
		&result.JumlahCocok,
		&result.TotalGejala,
		&result.Probabilitas,
	)

	return result, err
}

func simpanRiwayatKonsultasi(request models.KonsultasiRequest, result models.KonsultasiResult) (int, error) {
	gejalaJSON, err := json.Marshal(request.Gejala)
	if err != nil {
		return 0, err
	}

	var idRiwayat int

	err = config.DB.QueryRow(context.Background(), `
		INSERT INTO riwayat_konsultasi 
			(id_user, gejala_dipilih, hasil_diagnosa, probabilitas, kode_penyakit, waktu_konsultasi)
		VALUES 
			($1, $2, $3, $4, $5, NOW())
		RETURNING id_riwayat
	`,
		request.IDUser,
		string(gejalaJSON),
		result.NamaPenyakit,
		result.Probabilitas,
		result.KodePenyakit,
	).Scan(&idRiwayat)

	if err != nil {
		return 0, err
	}

	return idRiwayat, nil
}

func GetHasilKonsultasi(c *gin.Context) {
	idRiwayat := c.Param("id_riwayat")

	var hasil models.HasilKonsultasi

	err := config.DB.QueryRow(context.Background(), `
		SELECT 
			r.id_riwayat,
			r.id_user,
			u.nama,
			r.gejala_dipilih::text,
			r.hasil_diagnosa,
			r.probabilitas,
			r.kode_penyakit,
			p.nama_penyakit,
			p.deskripsi,
			p.solusi,
			TO_CHAR(r.waktu_konsultasi, 'YYYY-MM-DD HH24:MI:SS') AS waktu_konsultasi
		FROM riwayat_konsultasi r
		JOIN users u 
			ON r.id_user = u.id_user
		JOIN penyakit p 
			ON r.kode_penyakit = p.kode_penyakit
		WHERE r.id_riwayat = $1
	`, idRiwayat).Scan(
		&hasil.IDRiwayat,
		&hasil.IDUser,
		&hasil.NamaUser,
		&hasil.GejalaDipilih,
		&hasil.HasilDiagnosa,
		&hasil.Probabilitas,
		&hasil.KodePenyakit,
		&hasil.NamaPenyakit,
		&hasil.Deskripsi,
		&hasil.Solusi,
		&hasil.WaktuKonsultasi,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{
				"message": "Data hasil konsultasi tidak ditemukan",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Gagal mengambil hasil konsultasi",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Hasil konsultasi berhasil diambil",
		"data":    hasil,
	})
}