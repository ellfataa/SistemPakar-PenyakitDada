package models

type Gejala struct {
	KodeGejala string `json:"kode_gejala"`
	NamaGejala string `json:"nama_gejala"`
}

type GejalaRequest struct {
	KodeGejala string `json:"kode_gejala"`
	NamaGejala string `json:"nama_gejala"`
}