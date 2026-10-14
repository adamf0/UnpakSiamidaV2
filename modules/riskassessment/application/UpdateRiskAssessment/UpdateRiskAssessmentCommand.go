package application

type UpdateRiskAssessmentCommand struct {
	UUID           string `json:"uuid"`
	FakultasUuid   string `json:"fakultas_unit"`
	IdFakultasUnit *int   `json:"id_fakultas_unit"`
	TahunUuid      string `json:"tahun_uuid"`
	IdMasterTahun  *int   `json:"id_master_tahun"`
	Tahun          string `json:"tahun"`

	// Excel Sheet: ANALISIS RISIKO Fields
	IdMataProgram        int     `json:"id_mata_program"`
	MataProgramUuid      string  `json:"mata_program_uuid"`
	IsuInternalEksternal string  `json:"isu_internal_eksternal"`
	DeskripsiRisiko      string  `json:"deskripsi_risiko"`
	KemungkinanAwal      int     `json:"kemungkinan_awal"`
	DampakAwal           int     `json:"dampak_awal"`
	TingkatRisikoAwal    string  `json:"tingkat_risiko_awal"`
	TindakanMitigasi     string  `json:"tindakan_mitigasi"`
	Peluang              string  `json:"peluang"`
	KenyataanResidu      string  `json:"kenyataan_residu"`
	KemungkinanResidu    *int    `json:"kemungkinan_residu"`
	DampakResidu         *int    `json:"dampak_residu"`
	TingkatRisikoResidu  *string `json:"tingkat_risiko_residu"`
	RencanaTindakLanjut  string  `json:"rencana_tindak_lanjut"`

	// Legacy / Compatibility fields
	KodeRisiko       string `json:"kode_risiko"`
	Kategori         string `json:"kategori"`
	PernyataanRisiko string `json:"pernyataan_risiko"`
	Penyebab         string `json:"penyebab"`
	Dampak           string `json:"dampak"`
	Kemungkinan      *int   `json:"kemungkinan"`
	Keparahan        *int   `json:"keparahan"`
	TingkatRisiko    *int   `json:"tingkat_risiko"`
	LevelRisiko      string `json:"level_risiko"`
	RencanaMitigasi  string `json:"rencana_mitigasi"`
	Pengendalian     string `json:"pengendalian"`
	PenanggungJawab  string `json:"penanggung_jawab"`
	TargetWaktu      string `json:"target_waktu"`
	Status           string `json:"status"`
}
