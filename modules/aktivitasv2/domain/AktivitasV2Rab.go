package domain

import (
	common "UnpakSiamida/common/domain"
	"time"

	"github.com/google/uuid"
)

type AktivitasV2Rab struct {
	common.Entity `json:"-"`
	ID            uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	UUID          uuid.UUID `gorm:"type:char(36);uniqueIndex" json:"uuid"`
	IdAktivitasV2 uint      `gorm:"column:id_aktivitasv2" json:"id_aktivitasv2"`
	NoAkun        string    `gorm:"column:no_akun;size:50" json:"no_akun"`
	NamaAkun      string    `gorm:"column:nama_akun;type:text" json:"nama_akun"`
	Kategori      string    `gorm:"column:kategori;size:255" json:"kategori"`
	Rincian       string    `gorm:"column:rincian;type:text" json:"rincian"`
	Keterangan    string    `gorm:"column:keterangan;type:text" json:"keterangan"`
	QtySatuan     float64   `gorm:"column:qty_satuan;default:1" json:"qty_satuan"`
	QtyPaket      float64   `gorm:"column:qty_paket;default:1" json:"qty_paket"`
	SBU           float64   `gorm:"column:sbu;default:0" json:"sbu"`
	Total         float64   `gorm:"column:total;default:0" json:"total"`
	CreatedAt     time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt     time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
}

func (AktivitasV2Rab) TableName() string {
	return "aktivitasv2_rab"
}
