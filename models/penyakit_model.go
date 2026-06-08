package models

type Penyakit struct {
	KodePenyakit string `json:"kode_penyakit"`
	NamaPenyakit string `json:"nama_penyakit"`
	Deskripsi   string `json:"deskripsi"`
	Solusi string `json:"solusi"`
}