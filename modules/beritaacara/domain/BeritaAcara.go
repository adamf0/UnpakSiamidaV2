package domain

import (
	"time"

	common "UnpakSiamida/common/domain"
	event "UnpakSiamida/modules/beritaacara/event"

	"github.com/google/uuid"
)

type BeritaAcara struct {
	common.Entity

	ID           uint      `gorm:"primaryKey;autoIncrement"`
	UUID         uuid.UUID `gorm:"type:char(36);uniqueIndex"`
	Tahun        string    `gorm:"column:tahun"`
	FakultasUnit uint      `gorm:"column:fakultas_unit"`
	Tanggal      time.Time `gorm:"column:tanggal;type:date"`
	Auditee      uint      `gorm:"column:auditee"`
	Auditor1     *uint     `gorm:"column:auditor1"`
	Auditor2     *uint     `gorm:"column:auditor2"`
}

func (BeritaAcara) TableName() string {
	return "berita_acara"
}

// === CREATE ===
func NewBeritaAcara(
	tahun string,
	fakultasUnit uint,
	tanggal time.Time,
	auditee uint,
	auditor1 *uint,
	auditor2 *uint,
) common.ResultValue[*BeritaAcara] {
	if fakultasUnit <= 0 {
		return common.FailureValue[*BeritaAcara](NotFoundFakultas())
	}
	if auditee <= 0 {
		return common.FailureValue[*BeritaAcara](NotFoundAuditee())
	}
	if (auditor1 == nil || auditor2 == nil) || (*auditor1 <= 0 || *auditor2 <= 0) {
		return common.FailureValue[*BeritaAcara](NotFoundAuditor())
	}
	if auditee == *auditor1 || auditee == *auditor2 {
		return common.FailureValue[*BeritaAcara](DuplicateAssigment())
	}

	BeritaAcara := &BeritaAcara{
		UUID:         uuid.New(),
		Tahun:        tahun,
		FakultasUnit: fakultasUnit,
		Tanggal:      tanggal,
		Auditee:      auditee,
		Auditor1:     auditor1,
		Auditor2:     auditor2,
	}

	BeritaAcara.Raise(event.BeritaAcaraCreatedEvent{
		EventID:         uuid.New(),
		OccurredOn:      time.Now().UTC(),
		BeritaAcaraUUID: BeritaAcara.UUID,
	})

	return common.SuccessValue(BeritaAcara)
}

// === UPDATE ===
func UpdateBeritaAcara(
	prev *BeritaAcara,
	uid uuid.UUID,
	tahun string,
	fakultasUnit uint,
	tanggal time.Time,
	auditee uint,
	auditor1 *uint,
	auditor2 *uint,
) common.ResultValue[*BeritaAcara] {

	if prev == nil {
		return common.FailureValue[*BeritaAcara](EmptyData())
	}

	if prev.UUID != uid {
		return common.FailureValue[*BeritaAcara](InvalidData())
	}

	if fakultasUnit <= 0 {
		return common.FailureValue[*BeritaAcara](NotFoundFakultas())
	}
	if auditee <= 0 {
		return common.FailureValue[*BeritaAcara](NotFoundAuditee())
	}
	if (auditor1 == nil || auditor2 == nil) || (*auditor1 <= 0 || *auditor2 <= 0) {
		return common.FailureValue[*BeritaAcara](NotFoundAuditor())
	}
	if auditee == *auditor1 || auditee == *auditor2 {
		return common.FailureValue[*BeritaAcara](DuplicateAssigment())
	}

	prev.Tahun = tahun
	prev.FakultasUnit = fakultasUnit
	prev.Tanggal = tanggal
	prev.Auditee = auditee
	prev.Auditor1 = auditor1
	prev.Auditor2 = auditor2

	prev.Raise(event.BeritaAcaraUpdatedEvent{
		EventID:         uuid.New(),
		OccurredOn:      time.Now().UTC(),
		BeritaAcaraUUID: prev.UUID,
	})

	return common.SuccessValue(prev)
}

// func (b *BeritaAcara) RequestPdfGeneration(token string) {
// 	b.Raise(event.BeritaAcaraPdfRequestedEvent{
// 		EventID:         uuid.New(),
// 		OccurredOn:      time.Now().UTC(),
// 		BeritaAcaraUUID: b.UUID,
// 		Token:           token,
// 	})
// }
