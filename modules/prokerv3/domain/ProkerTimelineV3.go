package domain

import (
	common "UnpakSiamida/common/domain"
	"time"

	"github.com/google/uuid"
)

type ProkerTimelineV3 struct {
	common.Entity    `json:"-"`
	ID               uint       `gorm:"primaryKey;autoIncrement" json:"ID"`
	UUID             uuid.UUID  `gorm:"type:varchar(36)" json:"UUID"`
	IdAktivitasV2    uint       `gorm:"column:id_aktivitasv2" json:"IdAktivitasV2"`
	TanggalRkAwal    time.Time  `gorm:"column:tanggal_rk_awal;type:date" json:"TanggalRkAwal"`
	TanggalRkAkhir   time.Time  `gorm:"column:tanggal_rk_akhir;type:date" json:"TanggalRkAkhir"`
	TanggalRAwal     *time.Time `gorm:"column:tanggal_r_awal;type:date" json:"TanggalRAwal"`
	TanggalRAkhir    *time.Time `gorm:"column:tanggal_r_akhir;type:date" json:"TanggalRAkhir"`
	NoAktivitas      int        `gorm:"column:no_aktivitas;default:1" json:"NoAktivitas"`
	NamaAktivitas    string     `gorm:"column:nama_aktivitas;size:500" json:"NamaAktivitas"`
	PIC              string     `gorm:"column:pic;size:255" json:"PIC"`
	TitikKontrol     string     `gorm:"column:titik_kontrol;type:text" json:"TitikKontrol"`
	JadwalRk         string     `gorm:"column:jadwal_rk;type:text" json:"JadwalRk"`
	JadwalR          string     `gorm:"column:jadwal_r;type:text" json:"JadwalR"`
	Capaian          string     `gorm:"column:capaian;type:text" json:"Capaian"`
	Masalah          string     `gorm:"column:masalah;type:text" json:"Masalah"`
	PemecahanMasalah string     `gorm:"column:pemecahan_masalah;type:text" json:"PemecahanMasalah"`
	Status           string     `gorm:"column:status;size:50;default:'Belum Terlaksana'" json:"Status"`
	CreatedAt        time.Time  `gorm:"column:created_at;autoCreateTime" json:"CreatedAt"`
	UpdatedAt        time.Time  `gorm:"column:updated_at;autoUpdateTime" json:"UpdatedAt"`
}

func (ProkerTimelineV3) TableName() string {
	return "aktivitasv2_timeline"
}
