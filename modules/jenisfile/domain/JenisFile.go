package domain

import (
	common "UnpakSiamida/common/domain"

	"github.com/google/uuid"
)

type JenisFile struct {
	common.Entity

	ID              uint       `gorm:"primaryKey;autoIncrement"`
	UUID            uuid.UUID  `gorm:"type:char(36);uniqueIndex"`
	Nama       		string     `gorm:"type:longtext;not null"`
}

func (JenisFile) TableName() string {
	return "jenis_file_renstra"
}