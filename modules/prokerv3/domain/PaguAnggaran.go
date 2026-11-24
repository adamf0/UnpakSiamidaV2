package domain

import (
	"time"
)

type PaguAnggaran struct {
	ID             uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	IdFakultasUnit *int      `gorm:"column:id_fakultas_unit" json:"id_fakultas_unit"`
	FakultasUnit   string    `gorm:"column:fakultas_unit;size:100" json:"fakultas_unit"`
	Tahun          string    `gorm:"column:tahun;size:10;not null" json:"tahun"`
	Pagu           float64   `gorm:"column:pagu;default:0" json:"pagu"`
	CreatedAt      time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt      time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
}

func (PaguAnggaran) TableName() string {
	return "pagu_anggaran"
}

type SavePaguAnggaranRequest struct {
	ID             uint     `json:"id"`
	IdFakultasUnit *int     `json:"id_fakultas_unit"`
	FakultasUnit   string   `json:"fakultas_unit"`
	Tahun          string   `json:"tahun"`
	Pagu           *float64 `json:"pagu"`
}

type UnitPaguSummaryResult struct {
	HasPagu          bool    `json:"has_pagu"`
	Pagu             float64 `json:"pagu"`
	TotalAnggaran    float64 `json:"total_anggaran"`
	TotalRealisasi   float64 `json:"total_realisasi"`
	HasRealisasi     bool    `json:"has_realisasi"`
	HasRealDiff      bool    `json:"has_real_diff"`
	IsMelebihiPagu   bool    `json:"is_melebihi_pagu"`
	SisaPagu         float64 `json:"sisa_pagu"`
	NamaFakultasUnit string  `json:"nama_fakultas_unit"`
	IdFakultasUnit   int     `json:"id_fakultas_unit"`
	Tahun            string  `json:"tahun"`
}
