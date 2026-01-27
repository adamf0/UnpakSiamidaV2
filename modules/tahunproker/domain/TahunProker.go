package domain

import (
	"slices"
	"time"

	common "UnpakSiamida/common/domain"
	"UnpakSiamida/common/helper"
	event "UnpakSiamida/modules/tahunproker/event"

	"github.com/google/uuid"
)

type TahunProker struct {
	common.Entity

	ID     uint      `gorm:"primaryKey;autoIncrement"`
	UUID   uuid.UUID `gorm:"type:char(36);uniqueIndex"`
	Tahun  string    `gorm:"type:longtext;not null"`
	Status string    `gorm:"type:longtext;not null"`
}

func (TahunProker) TableName() string {
	return "master_tahun"
}

// === CREATE ===
func NewTahunProker(tahun string, status string) common.ResultValue[*TahunProker] {
	tahunint, err := helper.ParseInt64(tahun)
	if err != nil && err.Error() == "Number out of range" {
		return common.FailureValue[*TahunProker](TahunOOR())
	}
	if (err != nil && (err.Error() == "Must be a number" || err.Error() == "Invalid number")) || tahunint <= 2000 {
		return common.FailureValue[*TahunProker](InvalidTahun())
	}
	if !hasValidStatus(status) {
		return common.FailureValue[*TahunProker](InvalidStatus())
	}
	if !hasValidStatus(status) {
		return common.FailureValue[*TahunProker](InvalidStatus())
	}

	tahunproker := &TahunProker{
		UUID:   uuid.New(),
		Tahun:  tahun,
		Status: status,
	}

	tahunproker.Raise(event.TahunProkerCreatedEvent{
		EventID:         uuid.New(),
		OccurredOn:      time.Now().UTC(),
		TahunProkerUUID: tahunproker.UUID,
	})

	return common.SuccessValue(tahunproker)
}

// === UPDATE ===
func UpdateTahunProker(
	prev *TahunProker,
	uid uuid.UUID,
	tahun string,
	status string,
) common.ResultValue[*TahunProker] {

	if prev == nil {
		return common.FailureValue[*TahunProker](EmptyData())
	}

	if prev.UUID != uid {
		return common.FailureValue[*TahunProker](InvalidData())
	}

	tahunint, err := helper.ParseInt64(tahun)
	if err != nil && err.Error() == "Number out of range" {
		return common.FailureValue[*TahunProker](TahunOOR())
	}
	if (err != nil && (err.Error() == "Must be a number" || err.Error() == "Invalid number")) || tahunint <= 2000 {
		return common.FailureValue[*TahunProker](InvalidTahun())
	}
	if !hasValidStatus(status) {
		return common.FailureValue[*TahunProker](InvalidStatus())
	}

	prev.Tahun = tahun
	prev.Status = status

	prev.Raise(event.TahunProkerUpdatedEvent{
		EventID:         uuid.New(),
		OccurredOn:      time.Now().UTC(),
		TahunProkerUUID: prev.UUID,
	})

	return common.SuccessValue(prev)
}

func hasValidStatus(status string) bool {
	rule := []string{"aktif", "non-aktif"}
	return slices.Contains(rule, status)
}
