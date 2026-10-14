package CreateProkerV3

type TimelineItemInput struct {
	UUID             string `json:"uuid"`
	NoAktivitas      int    `json:"no_aktivitas"`
	NamaAktivitas    string `json:"nama_aktivitas"`
	PIC              string `json:"pic"`
	TitikKontrol     string `json:"titik_kontrol"`
	TanggalRkAwal    string `json:"tanggal_rk_awal"`
	TanggalRkAkhir   string `json:"tanggal_rk_akhir"`
	TanggalRAwal     string `json:"tanggal_r_awal"`
	TanggalRAkhir    string `json:"tanggal_r_akhir"`
	JadwalRk         string `json:"jadwal_rk"`
	JadwalR          string `json:"jadwal_r"`
	Capaian          string `json:"capaian"`
	Masalah          string `json:"masalah"`
	PemecahanMasalah string `json:"pemecahan_masalah"`
	Status           string `json:"status"`
}

type CreateProkerV3Command struct {
	IdFakultasUnit    *int                `json:"id_fakultas_unit"`
	FakultasUuid      string              `json:"fakultas_unit"`
	IdMataProgram     *int                `json:"id_mata_program"`
	IdIndikator       *int                `json:"id_indikator"`
	Tahun             string              `json:"tahun"`
	NoProker          *int64              `json:"no_proker"`
	NamaProgram       string              `json:"nama_program"`
	StandarSubStandar string              `json:"standar_sub_standar"`
	IndikatorKinerja  string              `json:"indikator_kinerja"`
	TargetUnpak       string              `json:"target_unpak"`
	CapaianUnitLalu   string              `json:"capaian_unit_lalu"`
	TargetUnit        string              `json:"target_unit"`
	Output            string              `json:"output"`
	Outcome           string              `json:"outcome"`
	Anggaran          float64             `json:"anggaran"`
	DokumenReferensi  string              `json:"dokumen_referensi"`
	LuaranProduk      string              `json:"luaran_produk"`
	PIC               string              `json:"pic"`
	Aktivitas         string              `json:"aktivitas"`
	Status            string              `json:"status"`
	Timelines         []TimelineItemInput `json:"timelines"`
}
