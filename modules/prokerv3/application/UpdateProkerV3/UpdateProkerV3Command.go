package UpdateProkerV3

import "UnpakSiamida/modules/prokerv3/application/CreateProkerV3"

type UpdateProkerV3Command struct {
	UUID              string                             `json:"uuid" form:"uuid"`
	IdFakultasUnit    *int                               `json:"id_fakultas_unit" form:"id_fakultas_unit"`
	FakultasUuid      string                             `json:"fakultas_unit" form:"fakultas_unit"`
	IdMataProgram     *int                               `json:"id_mata_program" form:"id_mata_program"`
	IdIndikator       *int                               `json:"id_indikator" form:"id_indikator"`
	Tahun             string                             `json:"tahun" form:"tahun"`
	NoProker          *int64                             `json:"no_proker" form:"no_proker"`
	NamaProgram       string                             `json:"nama_program" form:"nama_program"`
	StandarSubStandar string                             `json:"standar_sub_standar" form:"standar_sub_standar"`
	IndikatorKinerja  string                             `json:"indikator_kinerja" form:"indikator_kinerja"`
	TargetUnpak       string                             `json:"target_unpak" form:"target_unpak"`
	CapaianUnitLalu   string                             `json:"capaian_unit_lalu" form:"capaian_unit_lalu"`
	TargetUnit        string                             `json:"target_unit" form:"target_unit"`
	Output            string                             `json:"output" form:"output"`
	Outcome           string                             `json:"outcome" form:"outcome"`
	Anggaran          float64                            `json:"anggaran" form:"anggaran"`
	DokumenReferensi  string                             `json:"dokumen_referensi" form:"dokumen_referensi"`
	LuaranProduk      string                             `json:"luaran_produk" form:"luaran_produk"`
	PIC               string                             `json:"pic" form:"pic"`
	Aktivitas         string                             `json:"aktivitas" form:"aktivitas"`
	Status            string                             `json:"status" form:"status"`
	Timelines         []CreateProkerV3.TimelineItemInput `json:"timelines" form:"timelines"`
}
