package domain

import (
	"time"

	common "UnpakSiamida/common/domain"
	event "UnpakSiamida/modules/mataprogram/event"

	"github.com/google/uuid"
)

type MataProgram struct {
	common.Entity

	ID          uint      `gorm:"primaryKey;autoIncrement"`
	UUID        uuid.UUID `gorm:"type:char(36);uniqueIndex"`
	TahunID     uint      `gorm:"column:id_master_tahun"`
	MataProgram string    `gorm:"column:mata_program"`
}

func (MataProgram) TableName() string {
	return "mata_program"
}

// === CREATE ===
func NewMataProgram(tahun uint, mataprogram string) common.ResultValue[*MataProgram] {
	if tahun <= 0 {
		return common.FailureValue[*MataProgram](InvalidTahun())
	}

	data := &MataProgram{
		UUID:        uuid.New(),
		TahunID:     tahun,
		MataProgram: mataprogram,
	}

	data.Raise(event.MataProgramCreatedEvent{
		EventID:         uuid.New(),
		OccurredOn:      time.Now().UTC(),
		MataProgramUUID: data.UUID,
	})

	return common.SuccessValue(data)
}

// === UPDATE ===
func UpdateMataProgram(
	prev *MataProgram,
	uid uuid.UUID,
	tahun uint,
	mataprogram string,
) common.ResultValue[*MataProgram] {

	if prev == nil {
		return common.FailureValue[*MataProgram](EmptyData())
	}

	if prev.UUID != uid {
		return common.FailureValue[*MataProgram](InvalidData())
	}

	if tahun <= 0 {
		return common.FailureValue[*MataProgram](InvalidTahun())
	}

	prev.TahunID = tahun
	prev.MataProgram = mataprogram

	prev.Raise(event.MataProgramUpdatedEvent{
		EventID:         uuid.New(),
		OccurredOn:      time.Now().UTC(),
		MataProgramUUID: prev.UUID,
	})

	return common.SuccessValue(prev)
}
