package domain

import (
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
	MataProgram     string

	JenisDokumen string  `gorm:"JenisDokumen"`
	File         string  `gorm:"File"`
	Status       string  `gorm:"Status"`
	Catatan      *string `gorm:"Catatan"`
}
