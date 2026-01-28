package domain

import (
	"time"

	common "UnpakSiamida/common/domain"
	event "UnpakSiamida/modules/jadwalproker/event"

	"github.com/google/uuid"
)

type JadwalProker struct {
	common.Entity

	ID                  uint      `gorm:"primaryKey;autoIncrement"`
	UUID                uuid.UUID `gorm:"type:char(36);uniqueIndex"`
	FakultasUnit        uint      `gorm:"column:fakultas_unit"`
	TanggalTutupEntry   time.Time `gorm:"column:tanggal_tutup"`
	TanggalTutupDokumen time.Time `gorm:"column:tanggal_tutup_dokumen"`
}

func (JadwalProker) TableName() string {
	return "jadwal_proker"
}

// === CREATE ===
func NewJadwalProker(fakultasunit uint, tanggalTutupEntry string, tanggalTutupDokumen string) common.ResultValue[*JadwalProker] {
	format := "2006-01-02"

	if fakultasunit <= 0 {
		return common.FailureValue[*JadwalProker](NotFoundFakultas())
	}
	tanggalTutupEntryTime, err := time.Parse(format, tanggalTutupEntry)
	if err != nil {
		return common.FailureValue[*JadwalProker](InvalidDate("tanggal input"))
	}
	tanggalTutupDokumenTime, err := time.Parse(format, tanggalTutupDokumen)
	if err != nil {
		return common.FailureValue[*JadwalProker](InvalidDate("tanggal upload dokumen"))
	}

	if isOverlap(tanggalTutupEntryTime, tanggalTutupDokumenTime) {
		return common.FailureValue[*JadwalProker](InvalidDateRange())
	}

	jadwalproker := &JadwalProker{
		UUID:                uuid.New(),
		FakultasUnit:        fakultasunit,
		TanggalTutupEntry:   tanggalTutupEntryTime,
		TanggalTutupDokumen: tanggalTutupDokumenTime,
	}

	jadwalproker.Raise(event.JadwalProkerCreatedEvent{
		EventID:          uuid.New(),
		OccurredOn:       time.Now().UTC(),
		JadwalProkerUUID: jadwalproker.UUID,
	})

	return common.SuccessValue(jadwalproker)
}

// === UPDATE ===
func UpdateJadwalProker(
	prev *JadwalProker,
	uid uuid.UUID,
	fakultasunit uint,
	tanggalTutupEntry string,
	tanggalTutupDokumen string,
) common.ResultValue[*JadwalProker] {

	if prev == nil {
		return common.FailureValue[*JadwalProker](EmptyData())
	}

	if prev.UUID != uid {
		return common.FailureValue[*JadwalProker](InvalidData())
	}

	format := "2006-01-02"

	if fakultasunit <= 0 {
		return common.FailureValue[*JadwalProker](NotFoundFakultas())
	}
	tanggalTutupEntryTime, err := time.Parse(format, tanggalTutupEntry)
	if err != nil {
		return common.FailureValue[*JadwalProker](InvalidDate("tanggal input"))
	}
	tanggalTutupDokumenTime, err := time.Parse(format, tanggalTutupDokumen)
	if err != nil {
		return common.FailureValue[*JadwalProker](InvalidDate("tanggal upload dokumen"))
	}

	if isOverlap(tanggalTutupEntryTime, tanggalTutupDokumenTime) {
		return common.FailureValue[*JadwalProker](InvalidDateRange())
	}

	prev.FakultasUnit = fakultasunit
	prev.TanggalTutupEntry = tanggalTutupEntryTime
	prev.TanggalTutupDokumen = tanggalTutupDokumenTime

	prev.Raise(event.JadwalProkerUpdatedEvent{
		EventID:          uuid.New(),
		OccurredOn:       time.Now().UTC(),
		JadwalProkerUUID: prev.UUID,
	})

	return common.SuccessValue(prev)
}

func isOverlap(start1, end1 time.Time) bool {
	return !end1.Before(start1)
}
