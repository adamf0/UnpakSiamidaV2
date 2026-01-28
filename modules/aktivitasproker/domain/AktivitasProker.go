package domain

import (
	"time"

	common "UnpakSiamida/common/domain"
	event "UnpakSiamida/modules/aktivitasproker/event"

	"github.com/google/uuid"
)

type AktivitasProker struct {
	common.Entity

	ID             uint      `gorm:"primaryKey;autoIncrement"`
	UUID           uuid.UUID `gorm:"type:char(36);uniqueIndex"`
	MataProgram    uint
	FakultasUnit   uint `gorm:"column:fakultas_unit"`
	Aktivitas      string
	PIC            string
	TanggalRKAwal  time.Time `gorm:"column:tanggal_rk_awal"`
	TanggalRKAkhir time.Time `gorm:"column:tanggal_rk_akhir"`
}

func (AktivitasProker) TableName() string {
	return "aktivitas"
}

// === CREATE ===
func NewAktivitasProker(mataprogram uint, fakultasunit uint, aktivitas string, pic string, tanggalrkawal string, tanggalrkakhir string) common.ResultValue[*AktivitasProker] {
	format := "2006-01-02"

	if fakultasunit <= 0 {
		return common.FailureValue[*AktivitasProker](NotFoundFakultas())
	}
	if mataprogram <= 0 {
		return common.FailureValue[*AktivitasProker](NotFoundMataProgram())
	}
	tanggalRkAwalTime, err := time.Parse(format, tanggalrkawal)
	if err != nil {
		return common.FailureValue[*AktivitasProker](InvalidDate("tanggal rk awal"))
	}
	tanggalRkAkhirTime, err := time.Parse(format, tanggalrkakhir)
	if err != nil {
		return common.FailureValue[*AktivitasProker](InvalidDate("tanggal rk akhir"))
	}

	if isOverlap(tanggalRkAwalTime, tanggalRkAkhirTime) {
		return common.FailureValue[*AktivitasProker](InvalidDateRange())
	}

	aktivitasproker := &AktivitasProker{
		UUID:           uuid.New(),
		FakultasUnit:   fakultasunit,
		MataProgram:    mataprogram,
		Aktivitas:      aktivitas,
		PIC:            pic,
		TanggalRKAwal:  tanggalRkAwalTime,
		TanggalRKAkhir: tanggalRkAkhirTime,
	}

	aktivitasproker.Raise(event.AktivitasProkerCreatedEvent{
		EventID:             uuid.New(),
		OccurredOn:          time.Now().UTC(),
		AktivitasProkerUUID: aktivitasproker.UUID,
	})

	return common.SuccessValue(aktivitasproker)
}

// === UPDATE ===
func UpdateAktivitasProker(
	prev *AktivitasProker,
	uid uuid.UUID,
	mataprogram uint,
	fakultasunit uint,
	aktivitas string,
	pic string,
	tanggalrkawal string,
	tanggalrkakhir string,
) common.ResultValue[*AktivitasProker] {

	if prev == nil {
		return common.FailureValue[*AktivitasProker](EmptyData())
	}

	if prev.UUID != uid {
		return common.FailureValue[*AktivitasProker](InvalidData())
	}

	format := "2006-01-02"

	if fakultasunit <= 0 {
		return common.FailureValue[*AktivitasProker](NotFoundFakultas())
	}
	if mataprogram <= 0 {
		return common.FailureValue[*AktivitasProker](NotFoundMataProgram())
	}
	tanggalRkAwalTime, err := time.Parse(format, tanggalrkawal)
	if err != nil {
		return common.FailureValue[*AktivitasProker](InvalidDate("tanggal rk awal"))
	}
	tanggalRkAkhirTime, err := time.Parse(format, tanggalrkakhir)
	if err != nil {
		return common.FailureValue[*AktivitasProker](InvalidDate("tanggal rk akhir"))
	}

	if isOverlap(tanggalRkAwalTime, tanggalRkAkhirTime) {
		return common.FailureValue[*AktivitasProker](InvalidDateRange())
	}

	prev.FakultasUnit = fakultasunit
	prev.MataProgram = mataprogram
	prev.Aktivitas = aktivitas
	prev.PIC = pic
	prev.TanggalRKAwal = tanggalRkAwalTime
	prev.TanggalRKAkhir = tanggalRkAkhirTime

	prev.Raise(event.AktivitasProkerUpdatedEvent{
		EventID:             uuid.New(),
		OccurredOn:          time.Now().UTC(),
		AktivitasProkerUUID: prev.UUID,
	})

	return common.SuccessValue(prev)
}

func isOverlap(start1, end1 time.Time) bool {
	return !end1.After(start1)
}
