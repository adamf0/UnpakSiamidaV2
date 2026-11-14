package UpdateTimelineV3

type UpdateTimelineV3Command struct {
	UUID             string `json:"uuid" form:"uuid"`
	TanggalRkAwal    string `json:"tanggal_rk_awal" form:"tanggal_rk_awal"`
	TanggalRkAkhir   string `json:"tanggal_rk_akhir" form:"tanggal_rk_akhir"`
	TanggalRAwal     string `json:"tanggal_r_awal" form:"tanggal_r_awal"`
	TanggalRAkhir    string `json:"tanggal_r_akhir" form:"tanggal_r_akhir"`
	NoAktivitas      int    `json:"no_aktivitas" form:"no_aktivitas"`
	NamaAktivitas    string `json:"nama_aktivitas" form:"nama_aktivitas"`
	PIC              string `json:"pic" form:"pic"`
	TitikKontrol     string   `json:"titik_kontrol" form:"titik_kontrol"`
	Anggaran         *float64 `json:"anggaran" form:"anggaran"`
	RAB              string   `json:"rab" form:"rab"`
	JadwalRk         string   `json:"jadwal_rk" form:"jadwal_rk"`
	JadwalR          string `json:"jadwal_r" form:"jadwal_r"`
	Capaian          string `json:"capaian" form:"capaian"`
	Masalah          string `json:"masalah" form:"masalah"`
	PemecahanMasalah string `json:"pemecahan_masalah" form:"pemecahan_masalah"`
	Status           string `json:"status" form:"status"`
}
