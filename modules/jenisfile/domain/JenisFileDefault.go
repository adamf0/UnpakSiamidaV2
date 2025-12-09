package domain

import (
	"github.com/google/uuid"
)

type JenisFileDefault struct {
	Id          uint
	UUID            uuid.UUID  `gorm:"type:char(36);uniqueIndex"`
	Nama       		string     `gorm:"type:longtext;not null"`
}