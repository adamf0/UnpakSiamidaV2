package domain

import (
	"time"

	"github.com/google/uuid"
)

type JadwalAktivitasV3 struct {
	ID             uint       `gorm:"primaryKey;autoIncrement" json:"id"`
	UUID           uuid.UUID  `gorm:"type:char(36);uniqueIndex" json:"uuid"`
	Tahun          string     `gorm:"type:varchar(10);not null" json:"tahun"`
	TanggalMulai   time.Time  `gorm:"type:date;not null" json:"tanggal_mulai"`
	TanggalAkhir   time.Time  `gorm:"type:date;not null" json:"tanggal_akhir"`
	IdFakultasUnit *uint64    `gorm:"column:id_fakultas_unit" json:"id_fakultas_unit"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
}

func (JadwalAktivitasV3) TableName() string {
	return "jadwal_aktivitas_v3"
}

type JadwalAktivitasV3DTO struct {
	ID               uint    `json:"id"`
	UUID             string  `json:"uuid"`
	Tahun            string  `json:"tahun"`
	TanggalMulai     string  `json:"tanggal_mulai"`
	TanggalAkhir     string  `json:"tanggal_akhir"`
	IdFakultasUnit   *uint64 `json:"id_fakultas_unit"`
	NamaFakultasUnit string  `json:"nama_fakultas_unit"`
	TypeFakultasUnit string  `json:"type_fakultas_unit"`
	Status           string  `json:"status"` // "berjalan", "belum_dibuka", "ditutup"
	StatusLabel      string  `json:"status_label"`
	CreatedAt        string  `json:"created_at"`
	UpdatedAt        string  `json:"updated_at"`
}
