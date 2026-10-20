package application

type UpdateRiskAssessmentCommand struct {
	UUID           string `json:"uuid" form:"uuid"`
	FakultasUuid   string `json:"fakultas_unit" form:"fakultas_unit"`
	IdFakultasUnit *int   `json:"id_fakultas_unit" form:"id_fakultas_unit"`
	TahunUuid      string `json:"tahun_uuid" form:"tahun_uuid"`
	IdMasterTahun  *int   `json:"id_master_tahun" form:"id_master_tahun"`
	Tahun          string `json:"tahun" form:"tahun"`

	// Excel Sheet: ANALISIS RISIKO Fields
	IdMataProgram        int     `json:"id_mata_program" form:"id_mata_program"`
	MataProgramUuid      string  `json:"mata_program_uuid" form:"mata_program_uuid"`
	IsuInternalEksternal string  `json:"isu_internal_eksternal" form:"isu_internal_eksternal"`
	DeskripsiRisiko      string  `json:"deskripsi_risiko" form:"deskripsi_risiko"`
	KemungkinanAwal      int     `json:"kemungkinan_awal" form:"kemungkinan_awal"`
	DampakAwal           int     `json:"dampak_awal" form:"dampak_awal"`
	TingkatRisikoAwal    string  `json:"tingkat_risiko_awal" form:"tingkat_risiko_awal"`
	TindakanMitigasi     string  `json:"tindakan_mitigasi" form:"tindakan_mitigasi"`
	Peluang              string  `json:"peluang" form:"peluang"`
	KenyataanResidu      string  `json:"kenyataan_residu" form:"kenyataan_residu"`
	KemungkinanResidu    *int    `json:"kemungkinan_residu" form:"kemungkinan_residu"`
	DampakResidu         *int    `json:"dampak_residu" form:"dampak_residu"`
	TingkatRisikoResidu  *string `json:"tingkat_risiko_residu" form:"tingkat_risiko_residu"`
	RencanaTindakLanjut  string  `json:"rencana_tindak_lanjut" form:"rencana_tindak_lanjut"`

	// Legacy / Compatibility fields
	KodeRisiko       string `json:"kode_risiko" form:"kode_risiko"`
	Kategori         string `json:"kategori" form:"kategori"`
	PernyataanRisiko string `json:"pernyataan_risiko" form:"pernyataan_risiko"`
	Penyebab         string `json:"penyebab" form:"penyebab"`
	Dampak           string `json:"dampak" form:"dampak"`
	Kemungkinan      *int   `json:"kemungkinan" form:"kemungkinan"`
	Keparahan        *int   `json:"keparahan" form:"keparahan"`
	TingkatRisiko    *int   `json:"tingkat_risiko" form:"tingkat_risiko"`
	LevelRisiko      string `json:"level_risiko" form:"level_risiko"`
	RencanaMitigasi  string `json:"rencana_mitigasi" form:"rencana_mitigasi"`
	Pengendalian     string `json:"pengendalian" form:"pengendalian"`
	PenanggungJawab  string `json:"penanggung_jawab" form:"penanggung_jawab"`
	TargetWaktu      string `json:"target_waktu" form:"target_waktu"`
	Status           string `json:"status" form:"status"`
}
