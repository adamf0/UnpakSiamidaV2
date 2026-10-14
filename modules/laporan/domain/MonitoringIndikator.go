package domain

type MonitoringIndikator struct {
	Tahun            int    `json:"Tahun"`
	FakultasUnitId   int    `json:"FakultasUnitId"`
	FakultasUnitUuid string `json:"FakultasUnitUuid"`
	FakultasUnit     string `json:"FakultasUnit"`
	Jenjang          string `json:"Jenjang"`
	Type             string `json:"Type"`
	Fakultas         string `json:"Fakultas"`

	IndikatorId   int    `json:"IndikatorId"`
	IndikatorUuid string `json:"IndikatorUuid"`
	UndikatorUuid string `json:"UndikatorUuid,omitempty"`
	Indikator     string `json:"Indikator"`
	TipeTarget    string `json:"TipeTarget"`
	Operator      *string `json:"Operator"`

	Target      *string `json:"Target"`
	TargetMin   *string `json:"TargetMin"`
	TargetMax   *string `json:"TargetMax"`
	RangeTarget *string `json:"RangeTarget"`
	Satuan      *string `json:"Satuan"`

	Capaian        *string `json:"Capaian"`
	CapaianAuditor *string `json:"CapaianAuditor"`
}
