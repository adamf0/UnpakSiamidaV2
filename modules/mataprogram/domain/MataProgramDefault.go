package domain

import (
	"github.com/google/uuid"
)

type MataProgramDefault struct {
	Id               uint       `json:"Id"`
	UUID             uuid.UUID  `gorm:"type:char(36);uniqueIndex" json:"UUID"`
	TahunId          uint       `json:"TahunId"`
	TahunUUID        *uuid.UUID `json:"TahunUuid"`
	Tahun            string     `json:"Tahun"`
	Status           string     `json:"Status"`
	MataProgram      string     `json:"MataProgram"`
	FakultasUnitId   *uint      `json:"FakultasUnitId"`
	FakultasUnitUUID *uuid.UUID `json:"FakultasUnitUuid"`
	FakultasUnit     string     `json:"FakultasUnit"`
	Jenjang          *string    `json:"Jenjang"`
	Type             string     `json:"Type"`
	Fakultas         *string    `json:"Fakultas"`
}
