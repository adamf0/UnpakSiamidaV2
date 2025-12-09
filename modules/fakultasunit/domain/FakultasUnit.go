package domain

import (
	common "UnpakSiamida/common/domain"

	"github.com/google/uuid"
)

type FakultasUnit struct {
	common.Entity
	ID           uint       `gorm:""`
	UUID         uuid.UUID  `gorm:"type:char(36);"`
	Nama         string     `gorm:"column:nama_fak_prod_unit;"`
	KodeJenjang  string     `gorm:"column:kode_jenjang;"`
	Jenjang      string     `gorm:"column:jenjang;"`
	Type         string     `gorm:"column:type;"`
	Fakultas     string     `gorm:"column:fakultas;"`
}
func (FakultasUnit) TableName() string {
	return "v_fakultas_unit"
}