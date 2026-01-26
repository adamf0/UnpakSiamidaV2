package domain

import (
	"github.com/google/uuid"
)

type MataProgramDefault struct {
	Id     uint
	UUID   uuid.UUID `gorm:"type:char(36);uniqueIndex"`
	Tahun  string    `gorm:"type:longtext;not null"`
	Status string    `gorm:"type:longtext;not null"`
}
