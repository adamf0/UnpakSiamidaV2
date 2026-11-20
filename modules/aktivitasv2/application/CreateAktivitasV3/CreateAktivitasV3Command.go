package CreateAktivitasV3

import CreateProker "UnpakSiamida/modules/prokerv3/application/CreateProkerV3"

type CreateAktivitasV3Command struct {
	MatrixUUID       string                           `json:"matrix_uuid" form:"matrix_uuid"`
	Aktivitas        string                           `json:"aktivitas" form:"aktivitas"`
	PIC              string                           `json:"pic" form:"pic"`
	Anggaran         float64                          `json:"anggaran" form:"anggaran"`
	RAB              string                           `json:"rab" form:"rab"`
	LuaranProduk     string                           `json:"luaran_produk" form:"luaran_produk"`
	DokumenReferensi string                           `json:"dokumen_referensi" form:"dokumen_referensi"`
	Status           string                           `json:"status" form:"status"`
	Capaian          string                           `json:"capaian" form:"capaian"`
	Masalah          string                           `json:"masalah" form:"masalah"`
	PemecahanMasalah string                           `json:"pemecahan_masalah" form:"pemecahan_masalah"`
	Output           string                           `json:"output" form:"output"`
	Outcome          string                           `json:"outcome" form:"outcome"`
	BebanYayasan     bool                             `json:"beban_yayasan" form:"beban_yayasan"`
	Timelines        []CreateProker.TimelineItemInput `json:"timelines" form:"timelines"`
}
