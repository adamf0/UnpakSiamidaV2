package UpdateProkerV3

import "UnpakSiamida/modules/prokerv3/application/CreateProkerV3"

type UpdateProkerV3Command struct {
	UUID              string                             `json:"uuid"`
	IdFakultasUnit    *int                               `json:"id_fakultas_unit"`
	FakultasUuid      string                             `json:"fakultas_unit"`
	IdMataProgram     *int                               `json:"id_mata_program"`
	IdIndikator       *int                               `json:"id_indikator"`
	Tahun             string                             `json:"tahun"`
	NoProker          *int64                             `json:"no_proker"`
	NamaProgram       string                             `json:"nama_program"`
	StandarSubStandar string                             `json:"standar_sub_standar"`
	IndikatorKinerja  string                             `json:"indikator_kinerja"`
	TargetUnpak       string                             `json:"target_unpak"`
	CapaianUnitLalu   string                             `json:"capaian_unit_lalu"`
	TargetUnit        string                             `json:"target_unit"`
	Output            string                             `json:"output"`
	Outcome           string                             `json:"outcome"`
	Anggaran          float64                            `json:"anggaran"`
	DokumenReferensi  string                             `json:"dokumen_referensi"`
	LuaranProduk      string                             `json:"luaran_produk"`
	PIC               string                             `json:"pic"`
	Aktivitas         string                             `json:"aktivitas"`
	Status            string                             `json:"status"`
	Timelines         []CreateProkerV3.TimelineItemInput `json:"timelines"`
}
