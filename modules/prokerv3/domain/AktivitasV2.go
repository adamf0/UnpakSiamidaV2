package domain

import (
	common "UnpakSiamida/common/domain"
	"time"

	"github.com/google/uuid"
)

type AktivitasV2 struct {
	common.Entity    `json:"-"`
	ID               uint      `gorm:"primaryKey;autoIncrement" json:"ID"`
	UUID             uuid.UUID `gorm:"type:char(36);uniqueIndex" json:"UUID"`
	IdMatrikProker   *uint64   `gorm:"column:id_matrik_proker" json:"IdMatrikProker"`
	IdFakultasUnit   *uint64   `gorm:"column:id_fakultas_unit" json:"IdFakultasUnit"`
	Aktivitas        string    `gorm:"column:aktivitas;type:text" json:"Aktivitas"`
	PIC              string    `gorm:"column:PIC;type:text" json:"PIC"`
	Capaian          string    `gorm:"column:capaian;type:text" json:"Capaian"`
	Masalah          string    `gorm:"column:masalah;type:text" json:"Masalah"`
	PemecahanMasalah string    `gorm:"column:pemecahan_masalah;type:text" json:"PemecahanMasalah"`
	Anggaran         float64   `gorm:"column:anggaran;default:0" json:"Anggaran"`
	Realisasi        *float64  `gorm:"column:realisasi" json:"Realisasi"`
	Status           string    `gorm:"column:status;size:50;default:'draft'" json:"Status"`
	CatatanPenolakan string    `gorm:"column:catatan_penolakan;type:text" json:"CatatanPenolakan"`
	StatusVerifikator string   `gorm:"column:status_verifikator;size:50;default:'draft'" json:"StatusVerifikator"`
	CatatanVerifikator string  `gorm:"column:catatan_verifikator;type:text" json:"CatatanVerifikator"`
	StatusAnggaran   string    `gorm:"column:status_anggaran;size:50;default:'menunggu_verifikator'" json:"StatusAnggaran"`
	CatatanAnggaran  string    `gorm:"column:catatan_anggaran;type:text" json:"CatatanAnggaran"`
	StatusUpml       string    `gorm:"column:status_upml;size:50;default:'menunggu'" json:"StatusUpml"`
	CatatanTolakUpml string    `gorm:"column:catatan_tolak_upml;type:text" json:"CatatanTolakUpml"`
	StatusUpmf       string    `gorm:"column:status_upmf;size:50;default:'menunggu'" json:"StatusUpmf"`
	CatatanTolakUpmf string    `gorm:"column:catatan_tolak_upmf;type:text" json:"CatatanTolakUpmf"`
	StatusUpms       string    `gorm:"column:status_upms;size:50;default:'menunggu'" json:"StatusUpms"`
	CatatanTolakUpms string    `gorm:"column:catatan_tolak_upms;type:text" json:"CatatanTolakUpms"`
	StatusLpm        string    `gorm:"column:status_lpm;size:50;default:'menunggu'" json:"StatusLpm"`
	CatatanTolakLpm  string    `gorm:"column:catatan_tolak_lpm;type:text" json:"CatatanTolakLpm"`
	CreatedAt        time.Time `gorm:"column:created_at;autoCreateTime" json:"CreatedAt"`
	UpdatedAt        time.Time `gorm:"column:updated_at;autoUpdateTime" json:"UpdatedAt"`
}

func (AktivitasV2) TableName() string {
	return "aktivitasv2"
}
