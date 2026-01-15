package domain

import (
	"time"

	common "UnpakSiamida/common/domain"
	event "UnpakSiamida/modules/indikatorrenstra/event"

	"github.com/google/uuid"
)

type IndikatorRenstra struct {
	common.Entity

	ID             uint      `gorm:"primaryKey;autoIncrement"`
	UUID           uuid.UUID `gorm:"type:char(36);uniqueIndex"`
	StandarRenstra *uint     `gorm:"column:id_master_standar;"`
	Indikator      string    `gorm:"type:longtext;not null"`
	Parent         *uint     `gorm:""`
	Tahun          string    `gorm:"size:255;not null"`
	TipeTarget     string    `gorm:"size:255;not null"`
	Operator       *string   `gorm:""`
}

func (IndikatorRenstra) TableName() string {
	return "master_indikator_renstra"
}

func NewIndikatorRenstra(
	indikator string,
	standar *uint,
	parent *uint,
	tahun string,
	tipeTarget string,
	operator *string,
	isUniqueIndikator bool,
) common.ResultValue[*IndikatorRenstra] {

	if standar == nil {
		return common.FailureValue[*IndikatorRenstra](InvalidStandar())
	}

	if !isUniqueIndikator {
		return common.FailureValue[*IndikatorRenstra](NotUniqueIndikator())
	}

	ir := &IndikatorRenstra{
		UUID:           uuid.New(),
		Indikator:      indikator,
		StandarRenstra: standar,
		Parent:         parent,
		Tahun:          tahun,
		TipeTarget:     tipeTarget,
		Operator:       operator,
	}

	ir.Raise(event.IndikatorRenstraCreatedEvent{
		EventID:              uuid.New(),
		OccurredOn:           time.Now().UTC(),
		IndikatorRenstraUUID: ir.UUID,
	})

	return common.SuccessValue(ir)
}

func UpdateIndikatorRenstra(
	prev *IndikatorRenstra,
	uid uuid.UUID,
	indikator string,
	standar *uint,
	parent *uint,
	tahun string,
	tipeTarget string,
	operator *string,
	// isUniqueIndikator bool,
) common.ResultValue[*IndikatorRenstra] {

	if standar == nil {
		return common.FailureValue[*IndikatorRenstra](InvalidStandar())
	}

	// if !isUniqueIndikator {
	// 	return common.FailureValue[*IndikatorRenstra](NotUniqueIndikator())
	// }

	if prev == nil {
		return common.FailureValue[*IndikatorRenstra](EmptyData())
	}

	if prev.UUID != uid {
		return common.FailureValue[*IndikatorRenstra](InvalidData())
	}

	// always overwrite (required field)
	prev.Indikator = indikator
	prev.Tahun = tahun
	prev.TipeTarget = tipeTarget

	// optional pointer fields updated only if not nil
	prev.StandarRenstra = standar

	if parent != nil {
		prev.Parent = parent
	}

	if operator != nil {
		prev.Operator = operator
	}

	prev.Raise(event.IndikatorRenstraUpdatedEvent{
		EventID:              uuid.New(),
		OccurredOn:           time.Now().UTC(),
		IndikatorRenstraUUID: prev.UUID,
	})

	return common.SuccessValue(prev)
}
