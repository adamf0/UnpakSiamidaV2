package domain

import (
	"time"

	common "UnpakSiamida/common/domain"
	helper "UnpakSiamida/common/helper"
	event "UnpakSiamida/modules/generaterenstra/event"

	"github.com/google/uuid"
)

type GenerateRenstra struct {
	common.Entity

	ID              uint      `gorm:"primaryKey;autoIncrement"`
	UUID            uuid.UUID `gorm:"type:char(36);uniqueIndex"`
	RenstraId       uint      `gorm:"column:id_renstra;"`
	TemplateRenstra uint      `gorm:"column:template_renstra;"`
	Tugas           string    `gorm:""`
}

func (GenerateRenstra) TableName() string {
	return "renstra_nilai"
}

func NewGenerateRenstra(
	tahun string, //tahun template
	renstraTahun string, //tahun renstra

	fakultasUnit uint, //fakunit template
	renstraFakultasUnit uint, //fakunit renstra

	renstraId uint,

	template uint,
	templateUuid string,
	indikator string,
	tugas string,
	operation string,
) common.ResultValue[*GenerateRenstra] {

	if renstraTahun != tahun {
		return common.FailureValue[*GenerateRenstra](InvalidTahunRenstra(templateUuid, indikator, tahun, renstraTahun, operation))
	}
	if renstraFakultasUnit != fakultasUnit {
		return common.FailureValue[*GenerateRenstra](InvalidFakultasUnit())
	}
	if template <= 0 {
		return common.FailureValue[*GenerateRenstra](InvalidTemplate())
	}
	if renstraId <= 0 {
		return common.FailureValue[*GenerateRenstra](InvalidRenstra())
	}
	if !helper.IsValidTugas(tugas) {
		return common.FailureValue[*GenerateRenstra](InvalidTugas())
	}

	renstra := &GenerateRenstra{
		UUID:            uuid.New(),
		RenstraId:       renstraId,
		TemplateRenstra: template,
		Tugas:           tugas,
	}

	renstra.Raise(event.GenerateRenstraCreatedEvent{
		EventID:             uuid.New(),
		OccurredOn:          time.Now().UTC(),
		GenerateRenstraUUID: renstra.UUID,
	})

	return common.SuccessValue(renstra)
}
