package CreateAktivitasV3

import CreateProker "UnpakSiamida/modules/prokerv3/application/CreateProkerV3"

type CreateAktivitasV3Command struct {
	MatrixUUID       string                           `json:"matrix_uuid" form:"matrix_uuid"`
	Aktivitas        string                           `json:"aktivitas" form:"aktivitas"`
	PIC              string                           `json:"pic" form:"pic"`
	Anggaran         float64                          `json:"anggaran" form:"anggaran"`
	LuaranProduk     string                           `json:"luaran_produk" form:"luaran_produk"`
	DokumenReferensi string                           `json:"dokumen_referensi" form:"dokumen_referensi"`
	Status           string                           `json:"status" form:"status"`
	Capaian          string                           `json:"capaian" form:"capaian"`
	Masalah          string                           `json:"masalah" form:"masalah"`
	PemecahanMasalah string                           `json:"pemecahan_masalah" form:"pemecahan_masalah"`
	Timelines        []CreateProker.TimelineItemInput `json:"timelines" form:"timelines"`
}
