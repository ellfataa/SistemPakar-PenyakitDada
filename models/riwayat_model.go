package models

type RiwayatGejala struct {
	KodeGejala string `json:"kode_gejala"`
	NamaGejala string `json:"nama_gejala"`
}

type RiwayatKonsultasi struct {
	IDRiwayat       int             `json:"id_riwayat"`
	IDUser          int             `json:"id_user"`
	NamaUser        string          `json:"nama_user"`
	HasilDiagnosa   string          `json:"hasil_diagnosa"`
	Probabilitas    float64         `json:"probabilitas"`
	KodePenyakit    string          `json:"kode_penyakit"`
	NamaPenyakit    string          `json:"nama_penyakit"`
	WaktuKonsultasi string          `json:"waktu_konsultasi"`
	Gejala          []RiwayatGejala `json:"gejala"`
}