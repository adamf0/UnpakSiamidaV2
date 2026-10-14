package domain

import (
	common "UnpakSiamida/common/domain"
	"time"

	"github.com/google/uuid"
)

type Swot struct {
	common.Entity  `json:"-"`
	ID             uint      `gorm:"primaryKey;autoIncrement" json:"ID"`
	UUID           uuid.UUID `gorm:"type:char(36);uniqueIndex" json:"UUID"`
	IdFakultasUnit *int      `gorm:"column:id_fakultas_unit" json:"IdFakultasUnit"`
	IdMasterTahun  *int      `gorm:"column:id_master_tahun" json:"IdMasterTahun"`
	Tahun          string    `gorm:"column:tahun;size:10" json:"Tahun"`
	Kategori       string    `gorm:"column:kategori;size:50;not null" json:"Kategori"` // Kekuatan, Kelemahan, Peluang, Ancaman
	Faktor         string    `gorm:"column:faktor;size:255" json:"Faktor"`
	Deskripsi      string    `gorm:"column:deskripsi;type:text" json:"Deskripsi"`
	Strategi       string    `gorm:"column:strategi;type:text" json:"Strategi"`
	Bobot          *float64  `gorm:"column:bobot;type:decimal(5,2)" json:"Bobot"`
	Rating         *int      `gorm:"column:rating" json:"Rating"`
	Skor           *float64  `gorm:"column:skor;type:decimal(5,2)" json:"Skor"`
	CreatedAt      time.Time `gorm:"column:created_at;autoCreateTime" json:"CreatedAt"`
	UpdatedAt      time.Time `gorm:"column:updated_at;autoUpdateTime" json:"UpdatedAt"`
}

func (Swot) TableName() string {
	return "swot"
}

// === CREATE ===
func NewSwot(
	idFakultasUnit *int,
	idMasterTahun *int,
	tahun string,
	kategori string,
	faktor string,
	deskripsi string,
	strategi string,
	bobot *float64,
	rating *int,
	skor *float64,
) common.ResultValue[*Swot] {
	if kategori == "" {
		return common.FailureValue[*Swot](InvalidKategori())
	}

	data := &Swot{
		UUID:           uuid.New(),
		IdFakultasUnit: idFakultasUnit,
		IdMasterTahun:  idMasterTahun,
		Tahun:          tahun,
		Kategori:       kategori,
		Faktor:         faktor,
		Deskripsi:      deskripsi,
		Strategi:       strategi,
		Bobot:          bobot,
		Rating:         rating,
		Skor:           skor,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	return common.SuccessValue(data)
}

// === UPDATE ===
func UpdateSwot(
	prev *Swot,
	uid uuid.UUID,
	idFakultasUnit *int,
	idMasterTahun *int,
	tahun string,
	kategori string,
	faktor string,
	deskripsi string,
	strategi string,
	bobot *float64,
	rating *int,
	skor *float64,
) common.ResultValue[*Swot] {
	if prev == nil {
		return common.FailureValue[*Swot](EmptyData())
	}
	if prev.UUID != uid {
		return common.FailureValue[*Swot](InvalidData())
	}
	if kategori == "" {
		return common.FailureValue[*Swot](InvalidKategori())
	}

	prev.IdFakultasUnit = idFakultasUnit
	prev.IdMasterTahun = idMasterTahun
	prev.Tahun = tahun
	prev.Kategori = kategori
	prev.Faktor = faktor
	prev.Deskripsi = deskripsi
	prev.Strategi = strategi
	prev.Bobot = bobot
	prev.Rating = rating
	prev.Skor = skor
	prev.UpdatedAt = time.Now()

	return common.SuccessValue(prev)
}
