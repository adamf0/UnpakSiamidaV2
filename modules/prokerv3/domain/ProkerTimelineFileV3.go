package domain

import (
	common "UnpakSiamida/common/domain"
	"time"

	"github.com/google/uuid"
)

type ProkerTimelineFileV3 struct {
	common.Entity         `json:"-"`
	ID                    uint      `gorm:"primaryKey;autoIncrement" json:"ID"`
	UUID                  uuid.UUID `gorm:"type:varchar(36)" json:"UUID"`
	IdAktivitasV2Timeline *uint     `gorm:"column:id_aktivitasv2_timeline" json:"IdAktivitasV2Timeline"`
	IdAktivitasV2         *uint     `gorm:"column:id_aktivitasv2" json:"IdAktivitasV2"`
	File                  string    `gorm:"column:file;type:text;not null" json:"File"`
	NamaFile              string    `gorm:"column:nama_file;size:255" json:"NamaFile"`
	FileURL               string    `gorm:"column:file_url;size:500" json:"FileURL"`
	FileSize              int64     `gorm:"column:file_size;default:0" json:"FileSize"`
	FileType              string    `gorm:"column:file_type;size:100" json:"FileType"`
	Keterangan            string    `gorm:"column:keterangan;type:text" json:"Keterangan"`
	CreatedAt             time.Time `gorm:"column:created_at;autoCreateTime" json:"CreatedAt"`
	UpdatedAt             time.Time `gorm:"column:updated_at;autoUpdateTime" json:"UpdatedAt"`
}

func (ProkerTimelineFileV3) TableName() string {
	return "aktivitasv2_timeline_file"
}
