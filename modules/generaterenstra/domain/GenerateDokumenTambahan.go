package domain

import (
	"time"

	common "UnpakSiamida/common/domain"
	helper "UnpakSiamida/common/helper"

	"github.com/google/uuid"
)

type GenerateDokumenTambahan struct {
	common.Entity

	ID              				uint       `gorm:"primaryKey;autoIncrement"`
	UUID            				uuid.UUID  `gorm:"type:char(36);uniqueIndex"`
	RenstraId            			uint  	   `gorm:"column:id_renstra;"`
	TemplateDokumenTambahan         uint       `gorm:"column:id_template_dokumen_tambahan;"`
	Tugas  							string     `gorm:""`
}

func (GenerateDokumenTambahan) TableName() string {
	return "dokumen_tambahan"
}

func NewGenerateDokumenTambahan(
	tahun string, //tahun template
	renstraTahun string, //tahun renstra
	
	fakultasUnit uint, //fakunit template
	renstraFakultasUnit uint, //fakunit renstra
	
	renstraId uint,
	
	template uint,
	templateUuid string,
	jenisFile string,
	tugas string,
	operation string,
) common.ResultValue[*GenerateDokumenTambahan] {

	if renstraTahun != tahun {
		return common.FailureValue[*GenerateDokumenTambahan](InvalidTahunDokumenTambahan(templateUuid, jenisFile, tahun, renstraTahun, operation))
	}
	if renstraFakultasUnit != fakultasUnit {
		return common.FailureValue[*GenerateDokumenTambahan](InvalidFakultasUnit())
	}
	if template <= 0 {
		return common.FailureValue[*GenerateDokumenTambahan](InvalidTemplate())
	}
	if renstraId <= 0 {
		return common.FailureValue[*GenerateDokumenTambahan](InvalidRenstra())
	}
	if !helper.IsValidTugas(tugas) {
		return common.FailureValue[*GenerateDokumenTambahan](InvalidTugas())
	}

	renstra := &GenerateDokumenTambahan{
		UUID:                   uuid.New(),
		RenstraId:              renstraId,
		TemplateDokumenTambahan:template,
		Tugas:         			tugas,
	}

	renstra.Raise(GenerateDokumenTambahanCreatedEvent{
		EventID:    uuid.New(),
		OccurredOn: time.Now().UTC(),
		GenerateDokumenTambahanUUID: renstra.UUID,
	})

	return common.SuccessValue(renstra)
}
