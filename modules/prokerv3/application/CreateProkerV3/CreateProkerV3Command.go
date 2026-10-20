package CreateProkerV3

type TimelineItemInput struct {
	UUID             string `json:"uuid" form:"uuid"`
	NoAktivitas      int    `json:"no_aktivitas" form:"no_aktivitas"`
	NamaAktivitas    string `json:"nama_aktivitas" form:"nama_aktivitas"`
	PIC              string `json:"pic" form:"pic"`
	TitikKontrol     string `json:"titik_kontrol" form:"titik_kontrol"`
	TanggalRkAwal    string `json:"tanggal_rk_awal" form:"tanggal_rk_awal"`
	TanggalRkAkhir   string `json:"tanggal_rk_akhir" form:"tanggal_rk_akhir"`
	TanggalRAwal     string `json:"tanggal_r_awal" form:"tanggal_r_awal"`
	TanggalRAkhir    string `json:"tanggal_r_akhir" form:"tanggal_r_akhir"`
	JadwalRk         string `json:"jadwal_rk" form:"jadwal_rk"`
	JadwalR          string `json:"jadwal_r" form:"jadwal_r"`
	Capaian          string `json:"capaian" form:"capaian"`
	Masalah          string `json:"masalah" form:"masalah"`
	PemecahanMasalah string `json:"pemecahan_masalah" form:"pemecahan_masalah"`
	Status           string `json:"status" form:"status"`
}

type CreateProkerV3Command struct {
	IdFakultasUnit    *int                `json:"id_fakultas_unit" form:"id_fakultas_unit"`
	FakultasUuid      string              `json:"fakultas_unit" form:"fakultas_unit"`
	IdMataProgram     *int                `json:"id_mata_program" form:"id_mata_program"`
	IdIndikator       *int                `json:"id_indikator" form:"id_indikator"`
	Tahun             string              `json:"tahun" form:"tahun"`
	NoProker          *int64              `json:"no_proker" form:"no_proker"`
	NamaProgram       string              `json:"nama_program" form:"nama_program"`
	StandarSubStandar string              `json:"standar_sub_standar" form:"standar_sub_standar"`
	IndikatorKinerja  string              `json:"indikator_kinerja" form:"indikator_kinerja"`
	TargetUnpak       string              `json:"target_unpak" form:"target_unpak"`
	CapaianUnitLalu   string              `json:"capaian_unit_lalu" form:"capaian_unit_lalu"`
	TargetUnit        string              `json:"target_unit" form:"target_unit"`
	Output            string              `json:"output" form:"output"`
	Outcome           string              `json:"outcome" form:"outcome"`
	Anggaran          float64             `json:"anggaran" form:"anggaran"`
	DokumenReferensi  string              `json:"dokumen_referensi" form:"dokumen_referensi"`
	LuaranProduk      string              `json:"luaran_produk" form:"luaran_produk"`
	PIC               string              `json:"pic" form:"pic"`
	Aktivitas         string              `json:"aktivitas" form:"aktivitas"`
	Status            string              `json:"status" form:"status"`
	Timelines         []TimelineItemInput `json:"timelines" form:"timelines"`
}
