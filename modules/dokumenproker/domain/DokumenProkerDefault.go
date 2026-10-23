package domain

import (
	"time"

	"github.com/google/uuid"
)

type DokumenProkerDefault struct {
	Id               uint
	UUID             uuid.UUID
	FakultasUnitId   uint       `json:"FakultasUnitId"`
	FakultasUnitUUID *uuid.UUID `json:"FakultasUnitUuid"`
	FakultasUnit     string     `json:"FakultasUnit"`
	Jenjang          *string    `json:"Jenjang"`
	Type             string     `json:"Type"`
	Fakultas         *string    `json:"Fakultas"`

	MataProgramId   uint       `json:"MataProgramId"`
	MataProgramUUID *uuid.UUID `json:"MataProgramUuid"`
	MataProgram     string     `json:"MataProgram"`

	TahunId   uint       `json:"TahunId"`
	TahunUUID *uuid.UUID `json:"TahunUuid"`
	Tahun     string     `json:"Tahun"`

	JenisDokumen string  `json:"JenisDokumen"`
	File         string  `json:"File"`
	Status       string  `json:"Status"`
	Catatan      *string `json:"Catatan"`

	NamaFile  *string    `json:"NamaFile,omitempty"`
	FileURL   *string    `json:"FileURL,omitempty"`
	FileSize  *int64     `json:"FileSize,omitempty"`
	FileType  *string    `json:"FileType,omitempty"`
	CreatedAt *time.Time `json:"CreatedAt,omitempty"`
}

