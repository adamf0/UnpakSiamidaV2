package domain

import (
	"time"

	"github.com/google/uuid"
)

type JadwalProkerV3 struct {
	ID             uint       `gorm:"primaryKey;autoIncrement" json:"id"`
	UUID           uuid.UUID  `gorm:"type:char(36);uniqueIndex" json:"uuid"`
	Tahun          string     `gorm:"type:varchar(10);not null" json:"tahun"`
	TanggalMulai   time.Time  `gorm:"type:date;not null" json:"tanggal_mulai"`
	TanggalAkhir   time.Time  `gorm:"type:date;not null" json:"tanggal_akhir"`
	IdFakultasUnit *uint64    `gorm:"column:id_fakultas_unit" json:"id_fakultas_unit"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
}

func (JadwalProkerV3) TableName() string {
	return "jadwal_proker_v3"
}

type JadwalProkerV3Row struct {
	JadwalProkerV3
	NamaFakultasUnit string `json:"nama_fakultas_unit"`
	TypeFakultasUnit string `json:"type_fakultas_unit"`
}

type JadwalProkerV3DTO struct {
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

type JadwalFilter struct {
	Search     string
	Tahun      string
	FilterUnit string
}

type CheckJadwalRequest struct {
	Tahun     string
	UnitParam string
}

type CheckJadwalResult struct {
	HasSchedule      bool    `json:"has_schedule"`
	IsOpen           bool    `json:"is_open"`
	Status           string  `json:"status"`
	StatusLabel      string  `json:"status_label"`
	Tahun            string  `json:"tahun"`
	TanggalMulai     string  `json:"tanggal_mulai,omitempty"`
	TanggalAkhir     string  `json:"tanggal_akhir,omitempty"`
	IdFakultasUnit   *uint64 `json:"id_fakultas_unit,omitempty"`
	IsUnitSpecific   bool    `json:"is_unit_specific"`
	NamaFakultasUnit string  `json:"nama_fakultas_unit,omitempty"`
	Message          string  `json:"message"`
}
