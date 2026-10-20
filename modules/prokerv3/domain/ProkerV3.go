package domain

import (
	common "UnpakSiamida/common/domain"
	"time"

	"github.com/google/uuid"
)

type ProkerV3 struct {
	common.Entity    `json:"-"`
	ID               uint      `gorm:"primaryKey;autoIncrement" json:"ID"`
	UUID             uuid.UUID `gorm:"type:varchar(36);uniqueIndex" json:"UUID"`
	IdFakultasUnit   *int      `gorm:"column:id_fakultas_unit" json:"IdFakultasUnit"`
	IdMataProgram    *int      `gorm:"column:id_mata_program" json:"IdMataProgram"`
	IdIndikator      *int      `gorm:"column:id_indikator" json:"IdIndikator"`
	TargetUnpak      string    `gorm:"column:target_unpak;type:text" json:"TargetUnpak"`
	CapaianTerdahulu string    `gorm:"column:capaian_terdahulu;type:text" json:"CapaianTerdahulu"`
	TargetUnit       string    `gorm:"column:target_unit;type:text" json:"TargetUnit"`
	Output           string    `gorm:"column:output;type:text" json:"Output"`
	Outcome          string    `gorm:"column:outcome;type:text" json:"Outcome"`
	Status            string    `gorm:"column:status;type:varchar(50);default:'draft'" json:"Status"`
	CatatanPenolakan  string    `gorm:"column:catatan_penolakan;type:text" json:"CatatanPenolakan"`
	StatusVerifikator string    `gorm:"column:status_verifikator;type:varchar(50);default:'draft'" json:"StatusVerifikator"`
	CatatanVerifikator string   `gorm:"column:catatan_verifikator;type:text" json:"CatatanVerifikator"`
	StatusAnggaran    string    `gorm:"column:status_anggaran;type:varchar(50);default:'menunggu_verifikator'" json:"StatusAnggaran"`
	CatatanAnggaran   string    `gorm:"column:catatan_anggaran;type:text" json:"CatatanAnggaran"`
	CreatedAt         time.Time `gorm:"column:created_at" json:"CreatedAt"`
	UpdatedAt        time.Time `gorm:"column:updated_at" json:"UpdatedAt"`
}

func (ProkerV3) TableName() string {
	return "matrix_proker"
}
