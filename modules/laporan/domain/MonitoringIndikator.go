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
	IndikatorUuid string `json:"UndikatorUuid"`
	Indikator     string `json:"Indikator"`
	TipeTarget    string `json:"TipeTarget"`

	Capaian        float64 `json:"Capaian"`
	CapaianAuditor float64 `json:"CapaianAuditor"`
}
