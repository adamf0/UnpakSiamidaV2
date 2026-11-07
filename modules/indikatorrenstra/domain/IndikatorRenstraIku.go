package domain

import (
	"time"
)

type IndikatorRenstraIku struct {
	ID                 uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	IdIndikatorRenstra uint64    `gorm:"column:id_indikator_renstra;not null" json:"id_indikator_renstra"`
	IdIku              uint64    `gorm:"column:id_iku;not null" json:"id_iku"`
	CreatedAt          time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt          time.Time `gorm:"column:updated_at" json:"updated_at"`
}

func (IndikatorRenstraIku) TableName() string {
	return "indikator_renstra_iku"
}

type MasterIku struct {
	ID           uint64  `gorm:"primaryKey;column:id" json:"id"`
	IdKonteks    uint64  `gorm:"column:id_konteks" json:"id_konteks"`
	Iku          string  `gorm:"column:iku" json:"iku"`
	Kategori     string  `gorm:"column:kategori" json:"kategori"`
	IdSub        *uint64 `gorm:"column:id_sub" json:"id_sub"`
	FullKategori string  `gorm:"column:full_kategori" json:"full_kategori"`
	Satuan       string  `gorm:"column:satuan" json:"satuan"`
	JenisIku     string  `gorm:"column:jenis_iku" json:"jenis_iku"`
	Target       string  `gorm:"column:target" json:"target"`
	FormulaText  string  `gorm:"column:formula_text" json:"formula_text"`
	SumberData   string  `gorm:"column:sumber_data" json:"sumber_data"`
}

func (MasterIku) TableName() string {
	return "master_indikator"
}

