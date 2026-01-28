package domain

import (
	"time"

	"github.com/google/uuid"
)

type JadwalProkerDefault struct {
	Id               uint
	UUID             uuid.UUID
	FakultasUnitId   uint       `json:"FakultasUnitId"`
	FakultasUnitUUID *uuid.UUID `json:"FakultasUnitUuid"`
	FakultasUnit     string     `json:"FakultasUnit"`
	Jenjang          *string    `json:"Jenjang"`
	Type             string     `json:"Type"`
	Fakultas         *string    `json:"Fakultas"`

	TanggalTutupEntry   time.Time `json:"TanggalTutup"`
	TanggalTutupDokumen time.Time `json:"TanggalTutupDokumen"`
}
