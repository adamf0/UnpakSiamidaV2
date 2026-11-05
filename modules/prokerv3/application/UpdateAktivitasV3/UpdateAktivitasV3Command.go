package UpdateAktivitasV3

type UpdateAktivitasV3Command struct {
	UUID             string   `json:"uuid" form:"uuid"`
	Aktivitas        string   `json:"aktivitas" form:"aktivitas"`
	PIC              string   `json:"pic" form:"pic"`
	Anggaran         float64  `json:"anggaran" form:"anggaran"`
	RAB              string   `json:"rab" form:"rab"`
	Realisasi        *float64 `json:"realisasi" form:"realisasi"`
	Pencairan        *float64 `json:"pencairan" form:"pencairan"`
	LuaranProduk     string   `json:"luaran_produk" form:"luaran_produk"`
	DokumenReferensi string   `json:"dokumen_referensi" form:"dokumen_referensi"`
	Status           string   `json:"status" form:"status"`
	Capaian          string   `json:"capaian" form:"capaian"`
	Masalah          string   `json:"masalah" form:"masalah"`
	PemecahanMasalah string   `json:"pemecahan_masalah" form:"pemecahan_masalah"`
}
