package models

type KonsultasiRequest struct {
	IDUser int      `json:"id_user"`
	Gejala []string `json:"gejala"`
}

type KonsultasiResult struct {
	KodePenyakit string  `json:"kode_penyakit"`
	NamaPenyakit string  `json:"nama_penyakit"`
	Deskripsi    string  `json:"deskripsi"`
	Solusi        string  `json:"solusi"`
	JumlahCocok  int     `json:"jumlah_cocok"`
	TotalGejala  int     `json:"total_gejala"`
	Probabilitas float64 `json:"probabilitas"`
}

type HasilKonsultasi struct {
	IDRiwayat       int     `json:"id_riwayat"`
	IDUser          int     `json:"id_user"`
	NamaUser        string  `json:"nama_user"`
	GejalaDipilih   string  `json:"gejala_dipilih"`
	HasilDiagnosa   string  `json:"hasil_diagnosa"`
	Probabilitas    float64 `json:"probabilitas"`
	KodePenyakit    string  `json:"kode_penyakit"`
	NamaPenyakit    string  `json:"nama_penyakit"`
	Deskripsi       string  `json:"deskripsi"`
	Solusi           string  `json:"solusi"`
	WaktuKonsultasi string  `json:"waktu_konsultasi"`
}