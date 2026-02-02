package domain

import (
	"github.com/google/uuid"
)

type MonitoringProker struct {
	FakultasUnitId   uint       `json:"FakultasUnitId"`
	FakultasUnitUUID *uuid.UUID `json:"FakultasUnitUuid"`
	FakultasUnit     string     `json:"FakultasUnit"`
	Jenjang          *string    `json:"Jenjang"`
	Type             string     `json:"Type"`
	Fakultas         *string    `json:"Fakultas"`

	TahunId   uint       `json:"TahunId"`
	TahunUUID *uuid.UUID `json:"TahunUuid"`
	Tahun     string     `json:"Tahun"`

	MataProgramId   uint       `json:"MataProgramId"`
	MataProgramUUID *uuid.UUID `json:"MataProgramUuid"`
	MataProgram     string     `json:"MataProgram"`

	Sk   int64 `json:"SK"`
	SkR0 int64 `json:"SKR0"`
	SkR1 int64 `json:"SKR1"`
	SkR2 int64 `json:"SKR2"`
	SkR3 int64 `json:"SKR3"`

	Sop   int64 `json:"SOP"`
	SopR0 int64 `json:"SOPR0"`
	SopR1 int64 `json:"SOPR1"`
	SopR2 int64 `json:"SOPR2"`
	SopR3 int64 `json:"SOPR3"`

	ProposalTor   int64 `json:"ProposalTOR"`
	ProposalTorR0 int64 `json:"ProposalTORR0"`
	ProposalTorR1 int64 `json:"ProposalTORR1"`
	ProposalTorR2 int64 `json:"ProposalTORR2"`
	ProposalTorR3 int64 `json:"ProposalTORR3"`

	Laporan   int64 `json:"Laporan"`
	LaporanR0 int64 `json:"LaporanR0"`
	LaporanR1 int64 `json:"LaporanR1"`
	LaporanR2 int64 `json:"LaporanR2"`
	LaporanR3 int64 `json:"LaporanR3"`

	DokumenPendukung   int64 `json:"DokumenPendukung"`
	DokumenPendukungR0 int64 `json:"DokumenPendukungR0"`
	DokumenPendukungR1 int64 `json:"DokumenPendukungR1"`
	DokumenPendukungR2 int64 `json:"DokumenPendukungR2"`
	DokumenPendukungR3 int64 `json:"DokumenPendukungR3"`
}
