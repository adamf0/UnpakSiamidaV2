package domain

import (
	"time"

	common "UnpakSiamida/common/domain"
	event "UnpakSiamida/modules/templaterenstra/event"

	"github.com/google/uuid"
)

type TemplateRenstra struct {
	common.Entity
	ID                 uint      `gorm:"primaryKey;autoIncrement"`
	UUID               uuid.UUID `gorm:"type:char(36);uniqueIndex"`
	Tahun              string    `gorm:""`
	IndikatorRenstraID uint      `gorm:"column:indikator;"`
	IsPertanyaan       bool      `gorm:"column:pertanyaan;"`
	FakultasUnit       uint      `gorm:"column:fakultas_unit;"`
	Kategori           string    `gorm:""`
	Klasifikasi        string    `gorm:""`
	Satuan             *string   `gorm:""`
	Target             *string   `gorm:""`
	TargetMin          *string   `gorm:"column:target_min;"`
	TargetMax          *string   `gorm:"column:target_max;"`
	Tugas              string    `gorm:""`
}

func (TemplateRenstra) TableName() string {
	return "template_renstra"
}

// === CREATE ===
func NewTemplateRenstra(
	tahun string,
	indikatorRenstraID uint,
	isPertanyaan bool,
	fakultasUnit uint,
	kategori string,
	klasifikasi string,
	satuan *string,
	target *string,
	targetMin *string,
	targetMax *string,
	tugas string,
) common.ResultValue[*TemplateRenstra] {
	hasTarget := target != nil && *target != ""
	hasTargetMin := targetMin != nil && *targetMin != ""
	hasTargetMax := targetMax != nil && *targetMax != ""

	// Valid hanya jika:
	// - Mode A: hasTarget == true AND hasTargetMin == false AND hasTargetMax == false
	// - Mode B: hasTarget == false AND hasTargetMin == true AND hasTargetMax == true
	if !((hasTarget && !hasTargetMin && !hasTargetMax) ||
		(!hasTarget && hasTargetMin && hasTargetMax)) {
		return common.FailureValue[*TemplateRenstra](InvalidValueTarget())
	}
	if indikatorRenstraID <= 0 {
		return common.FailureValue[*TemplateRenstra](IndikatorNotFound())
	}
	if fakultasUnit <= 0 {
		return common.FailureValue[*TemplateRenstra](FakultasUnitNotFound())
	}

	templaterenstra := &TemplateRenstra{
		UUID:               uuid.New(),
		Tahun:              tahun,
		IndikatorRenstraID: indikatorRenstraID,
		IsPertanyaan:       isPertanyaan,
		FakultasUnit:       fakultasUnit,
		Kategori:           kategori,
		Klasifikasi:        klasifikasi,
		Satuan:             satuan,
		Target:             target,
		TargetMin:          targetMin,
		TargetMax:          targetMax,
		Tugas:              tugas,
	}

	templaterenstra.Raise(event.TemplateRenstraCreatedEvent{
		EventID:             uuid.New(),
		OccurredOn:          time.Now().UTC(),
		TemplateRenstraUUID: templaterenstra.UUID,
	})

	return common.SuccessValue(templaterenstra)
}

// === UPDATE ===
func UpdateTemplateRenstra(
	prev *TemplateRenstra,
	uid uuid.UUID,
	tahun string,
	indikatorRenstraID uint,
	isPertanyaan bool,
	fakultasUnit uint,
	kategori string,
	klasifikasi string,
	satuan *string,
	target *string,
	targetMin *string,
	targetMax *string,
	tugas string,
) common.ResultValue[*TemplateRenstra] {

	hasTarget := target != nil && *target != ""
	hasTargetMin := targetMin != nil && *targetMin != ""
	hasTargetMax := targetMax != nil && *targetMax != ""

	if prev == nil {
		return common.FailureValue[*TemplateRenstra](EmptyData())
	}
	if prev.UUID != uid {
		return common.FailureValue[*TemplateRenstra](InvalidData())
	}
	// Valid hanya jika:
	// - Mode A: hasTarget == true AND hasTargetMin == false AND hasTargetMax == false
	// - Mode B: hasTarget == false AND hasTargetMin == true AND hasTargetMax == true
	if !((hasTarget && !hasTargetMin && !hasTargetMax) ||
		(!hasTarget && hasTargetMin && hasTargetMax)) {
		return common.FailureValue[*TemplateRenstra](InvalidValueTarget())
	}
	if indikatorRenstraID <= 0 {
		return common.FailureValue[*TemplateRenstra](IndikatorNotFound())
	}
	if fakultasUnit <= 0 {
		return common.FailureValue[*TemplateRenstra](FakultasUnitNotFound())
	}

	prev.Tahun = tahun
	prev.IndikatorRenstraID = indikatorRenstraID
	prev.IsPertanyaan = isPertanyaan
	prev.FakultasUnit = fakultasUnit
	prev.Kategori = kategori
	prev.Klasifikasi = klasifikasi
	prev.Satuan = satuan
	prev.Target = target
	prev.TargetMin = targetMin
	prev.TargetMax = targetMax
	prev.Tugas = tugas

	prev.Raise(event.TemplateRenstraUpdatedEvent{
		EventID:             uuid.New(),
		OccurredOn:          time.Now().UTC(),
		TemplateRenstraUUID: prev.UUID,
	})

	return common.SuccessValue(prev)
}
