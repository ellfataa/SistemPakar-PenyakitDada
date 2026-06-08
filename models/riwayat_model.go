package models

type RiwayatKonsultasi struct {
	IDRiwayat       int     `json:"id_riwayat"`
	IDUser          int     `json:"id_user"`
	GejalaDipilih   string  `json:"gejala_dipilih"`
	HasilDiagnosa   string  `json:"hasil_diagnosa"`
	Probabilitas    float64 `json:"probabilitas"`
	KodePenyakit    string  `json:"kode_penyakit"`
	WaktuKonsultasi string  `json:"waktu_konsultasi"`
}