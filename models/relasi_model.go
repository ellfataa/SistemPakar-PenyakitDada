package models

type RelasiPenyakitGejala struct {
	KodePenyakit string `json:"kode_penyakit"`
	NamaPenyakit string `json:"nama_penyakit"`
	KodeGejala   string `json:"kode_gejala"`
	NamaGejala   string `json:"nama_gejala"`
}

type RelasiRequest struct {
	KodePenyakit string `json:"kode_penyakit"`
	KodeGejala   string `json:"kode_gejala"`
}