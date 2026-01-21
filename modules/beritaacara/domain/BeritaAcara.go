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
	FakultasUnit int       `gorm:"column:fakultas_unit"`
	Tanggal      time.Time `gorm:"column:tanggal;type:date"`
	Auditee      *int      `gorm:"column:auditee"`
	Auditor1     *int      `gorm:"column:auditor1"`
	Auditor2     *int      `gorm:"column:auditor2"`
}

func (BeritaAcara) TableName() string {
	return "berita_acara"
}

// === CREATE ===
func NewBeritaAcara(
	tahun string,
	fakultasUnit int,
	tanggal time.Time,
	auditee *int,
	auditor1 *int,
	auditor2 *int,
) common.ResultValue[*BeritaAcara] {

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
	fakultasUnit int,
	tanggal time.Time,
	auditee *int,
	auditor1 *int,
	auditor2 *int,
) common.ResultValue[*BeritaAcara] {

	if prev == nil {
		return common.FailureValue[*BeritaAcara](EmptyData())
	}

	if prev.UUID != uid {
		return common.FailureValue[*BeritaAcara](InvalidData())
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
