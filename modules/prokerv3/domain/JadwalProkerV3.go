package domain

import (
	"time"

	"github.com/google/uuid"
)

type JadwalProkerV3 struct {
	ID           uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	UUID         uuid.UUID `gorm:"type:char(36);uniqueIndex" json:"uuid"`
	Tahun        string    `gorm:"type:varchar(10);not null" json:"tahun"`
	TanggalMulai time.Time `gorm:"type:date;not null" json:"tanggal_mulai"`
	TanggalAkhir time.Time `gorm:"type:date;not null" json:"tanggal_akhir"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func (JadwalProkerV3) TableName() string {
	return "jadwal_proker_v3"
}

type JadwalProkerV3DTO struct {
	ID           uint   `json:"id"`
	UUID         string `json:"uuid"`
	Tahun        string `json:"tahun"`
	TanggalMulai string `json:"tanggal_mulai"`
	TanggalAkhir string `json:"tanggal_akhir"`
	Status       string `json:"status"` // "berjalan", "belum_dibuka", "ditutup"
	StatusLabel  string `json:"status_label"`
	CreatedAt    string `json:"created_at"`
	UpdatedAt    string `json:"updated_at"`
}
