package domain

import (
	"time"

	"github.com/google/uuid"
)

type AktivitasProkerDefault struct {
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

	Aktivitas      string
	PIC            string
	TanggalRKAwal  time.Time `gorm:"tanggal_rk_awal"`
	TanggalRKAkhir time.Time `gorm:"tanggal_rk_akhir"`
}
