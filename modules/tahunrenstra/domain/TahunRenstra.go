package domain

import (
	common "UnpakSiamida/common/domain"
)

type TahunRenstra struct {
	common.Entity
	Tahun        string     `gorm:""`
	Status       string     `gorm:""`
}
func (TahunRenstra) TableName() string {
	return "v_tahun_renstra"
}