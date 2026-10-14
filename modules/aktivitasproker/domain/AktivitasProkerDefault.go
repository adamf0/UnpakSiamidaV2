package domain

import (
	"time"

	"github.com/google/uuid"
)

type AktivitasProkerDefault struct {
	Id               uint       `json:"Id"`
	UUID             uuid.UUID  `json:"UUID"`
	FakultasUnitId   uint       `json:"FakultasUnitId"`
	FakultasUnitUUID *uuid.UUID `json:"FakultasUnitUuid"`
	FakultasUnit     string     `json:"FakultasUnit"`
	Jenjang          *string    `json:"Jenjang"`
	Type             string     `json:"Type"`
	Fakultas         *string    `json:"Fakultas"`

	MataProgramId   uint       `json:"MataProgramId"`
	MataProgramUUID *uuid.UUID `json:"MataProgramUuid"`
	MataProgram     string     `json:"MataProgram"`

	TahunId   *uint      `json:"TahunId"`
	TahunUUID *uuid.UUID `json:"TahunUuid"`
	Tahun     *string    `json:"Tahun"`

	Aktivitas       string     `json:"Aktivitas"`
	PIC             string     `json:"PIC"`
	TanggalRKAwal   time.Time  `json:"TanggalRKAwal"`
	TanggalRKAkhir  time.Time  `json:"TanggalRKAkhir"`
	StatusRealisasi *string    `json:"StatusRealisasi"`
	TanggalRAwal    *time.Time `json:"TanggalRAwal"`
	TanggalRAkhir   *time.Time `json:"TanggalRAkhir"`
	Detail          *string    `json:"Detail"`
}
