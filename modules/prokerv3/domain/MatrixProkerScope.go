package domain

import (
	"time"

	"github.com/google/uuid"
)

type MatrixProkerScope struct {
	ID               uint      `gorm:"primaryKey;autoIncrement" json:"ID"`
	UUID             uuid.UUID `gorm:"type:char(36);uniqueIndex" json:"UUID"`
	IdMatrixProker   uint      `gorm:"column:id_matrix_proker" json:"IdMatrixProker"`
	IdMataProgram    *int      `gorm:"column:id_mata_program" json:"IdMataProgram"`
	IdIndikator      *int      `gorm:"column:id_indikator" json:"IdIndikator"`
	TargetUnpak      string    `gorm:"column:target_unpak;type:text" json:"TargetUnpak"`
	CapaianTerdahulu string    `gorm:"column:capaian_terdahulu;type:text" json:"CapaianTerdahulu"`
	TargetUnit       string    `gorm:"column:target_unit;type:text" json:"TargetUnit"`
	Output           string    `gorm:"column:output;type:text" json:"Output"`
	Outcome          string    `gorm:"column:outcome;type:text" json:"Outcome"`
	CreatedAt        time.Time `gorm:"column:created_at" json:"CreatedAt"`
	UpdatedAt        time.Time `gorm:"column:updated_at" json:"UpdatedAt"`
}

func (MatrixProkerScope) TableName() string {
	return "matrix_proker_scope"
}

type ProkerScopeDefault struct {
	ID                uint      `gorm:"column:id" json:"Id"`
	UUID              uuid.UUID `gorm:"column:uuid" json:"UUID"`
	IdMatrixProker    uint      `gorm:"column:id_matrix_proker" json:"IdMatrixProker"`
	IdMataProgram     *int      `gorm:"column:id_mata_program" json:"IdMataProgram"`
	IdIndikator       *int      `gorm:"column:id_indikator" json:"IdIndikator"`
	IndikatorKinerja  string    `gorm:"column:indikator_kinerja" json:"IndikatorKinerja"`
	IndikatorUUID     string    `gorm:"column:indikator_uuid" json:"IndikatorUUID"`
	IdMasterStandar   *int      `gorm:"column:id_master_standar" json:"IdMasterStandar"`
	StandarSubStandar string    `gorm:"column:standar_sub_standar" json:"StandarSubStandar"`
	StandarUUID       string    `gorm:"column:standar_uuid" json:"StandarUUID"`
	TargetUnpak       string    `gorm:"column:target_unpak" json:"TargetUnpak"`
	CapaianUnitLalu   string    `gorm:"column:capaian_unit_lalu" json:"CapaianUnitLalu"`
	TargetUnit        string    `gorm:"column:target_unit" json:"TargetUnit"`
	Output            string    `gorm:"column:output" json:"Output"`
	Outcome           string    `gorm:"column:outcome" json:"Outcome"`
	IkuList           string    `gorm:"column:iku_list" json:"IkuList"`
}
