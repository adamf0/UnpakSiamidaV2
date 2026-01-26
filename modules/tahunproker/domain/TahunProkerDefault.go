package domain

import (
	"github.com/google/uuid"
)

type TahunProkerDefault struct {
	Id     uint
	UUID   uuid.UUID `gorm:"type:char(36);uniqueIndex"`
	Tahun  string    `gorm:"type:longtext;not null"`
	Status string    `gorm:"type:longtext;not null"`
}
