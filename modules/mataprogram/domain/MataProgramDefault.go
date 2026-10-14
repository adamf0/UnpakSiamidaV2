package domain

import (
	"github.com/google/uuid"
)

type MataProgramDefault struct {
	Id          uint
	UUID        uuid.UUID  `gorm:"type:char(36);uniqueIndex"`
	TahunId     uint       `json:"TahunId"`
	TahunUUID   *uuid.UUID `json:"TahunUuid"`
	Tahun       string     `json:"Tahun"`
	Status      string     `json:"Status"`
	MataProgram string     `json:"MataProgram"`
}
