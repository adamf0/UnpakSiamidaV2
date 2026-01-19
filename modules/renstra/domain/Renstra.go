package domain

import (
	"time"

	common "UnpakSiamida/common/domain"
	event "UnpakSiamida/modules/renstra/event"

	"github.com/google/uuid"
)

type Renstra struct {
	common.Entity

	ID                            uint      `gorm:"primaryKey;autoIncrement"`
	UUID                          uuid.UUID `gorm:"type:char(36);uniqueIndex"`
	Tahun                         string    `gorm:"size:255;not null"`
	FakultasUnit                  uint      `gorm:"column:fakultas_unit;"`
	PeriodeUploadMulai            string    `gorm:"column:periode_upload_mulai;"`
	PeriodeUploadAkhir            string    `gorm:"column:periode_upload_akhir;"`
	PeriodeAssesmentDokumenMulai  string    `gorm:"column:periode_assesment_dokumen_mulai;"`
	PeriodeAssesmentDokumenAkhir  string    `gorm:"column:periode_assesment_dokumen_akhir;"`
	PeriodeAssesmentLapanganMulai string    `gorm:"column:periode_assesment_lapangan_mulai;"`
	PeriodeAssesmentLapanganAkhir string    `gorm:"column:periode_assesment_lapangan_akhir;"`
	KodeAkses                     *string   `gorm:"column:kodeAkses;"`
	Auditee                       uint      `gorm:""`
	Auditor1                      uint      `gorm:""`
	Auditor2                      uint      `gorm:""`
	Catatan1                      *string   `gorm:"column:catatan;"`
	Catatan2                      *string   `gorm:"column:catatan2;"`
}

func (Renstra) TableName() string {
	return "renstra"
}

// ------------------------
// Helper: check overlap
// ------------------------
func isOverlap(start1, end1, start2, end2 time.Time) bool {
	return !end1.Before(start2) && !end2.Before(start1)
}

// ------------------------
// Create New Renstra
// ------------------------
func NewRenstra(
	tahun string,
	fakultasUnit uint,
	periodeUploadMulai, periodeUploadAkhir string,
	periodeDokumenMulai, periodeDokumenAkhir string,
	periodeLapanganMulai, periodeLapanganAkhir string,
	auditee, auditor1, auditor2 uint,
	isUniqueData bool,
) common.ResultValue[*Renstra] {

	if fakultasUnit <= 0 {
		return common.FailureValue[*Renstra](InvalidFakultasUnit())
	}

	if !isUniqueData {
		return common.FailureValue[*Renstra](DataExisting())
	}

	if auditee <= 0 {
		return common.FailureValue[*Renstra](MissingAuditee())
	}
	if auditor1 <= 0 {
		return common.FailureValue[*Renstra](MissingAuditor1())
	}
	if auditor2 <= 0 {
		return common.FailureValue[*Renstra](MissingAuditor2())
	}

	if auditee == auditor1 || auditee == auditor2 || auditor1 == auditor2 {
		return common.FailureValue[*Renstra](DuplicateAssigment())
	}

	format := "2006-01-02"

	// Periode Upload
	startUpload, err := time.Parse(format, periodeUploadMulai)
	if err != nil {
		return common.FailureValue[*Renstra](InvalidDate("upload"))
	}
	endUpload, err := time.Parse(format, periodeUploadAkhir)
	if err != nil {
		return common.FailureValue[*Renstra](InvalidDate("upload"))
	}

	// Periode Assessment Dokumen
	startDokumen, err := time.Parse(format, periodeDokumenMulai)
	if err != nil {
		return common.FailureValue[*Renstra](InvalidDate("assessment dokumen"))
	}
	endDokumen, err := time.Parse(format, periodeDokumenAkhir)
	if err != nil {
		return common.FailureValue[*Renstra](InvalidDate("assessment dokumen"))
	}

	// Periode Assessment Lapangan
	startLapangan, err := time.Parse(format, periodeLapanganMulai)
	if err != nil {
		return common.FailureValue[*Renstra](InvalidDate("assessment lapangan"))
	}
	endLapangan, err := time.Parse(format, periodeLapanganAkhir)
	if err != nil {
		return common.FailureValue[*Renstra](InvalidDate("assessment lapangan"))
	}

	// Cek overlap
	if isOverlap(startUpload, endUpload, startDokumen, endDokumen) {
		return common.FailureValue[*Renstra](PeriodOverlapUploadDokumen())
	}

	if isOverlap(startUpload, endUpload, startLapangan, endLapangan) {
		return common.FailureValue[*Renstra](PeriodOverlapUploadLapangan())
	}

	if isOverlap(startDokumen, endDokumen, startLapangan, endLapangan) {
		return common.FailureValue[*Renstra](PeriodOverlapDokumenLapangan())
	}

	renstra := &Renstra{
		UUID:                          uuid.New(),
		Tahun:                         tahun,
		FakultasUnit:                  fakultasUnit,
		PeriodeUploadMulai:            periodeUploadMulai,
		PeriodeUploadAkhir:            periodeUploadAkhir,
		PeriodeAssesmentDokumenMulai:  periodeDokumenMulai,
		PeriodeAssesmentDokumenAkhir:  periodeDokumenAkhir,
		PeriodeAssesmentLapanganMulai: periodeLapanganMulai,
		PeriodeAssesmentLapanganAkhir: periodeLapanganAkhir,
		Auditee:                       auditee,
		Auditor1:                      auditor1,
		Auditor2:                      auditor2,
	}

	renstra.Raise(event.RenstraCreatedEvent{
		EventID:     uuid.New(),
		OccurredOn:  time.Now().UTC(),
		RenstraUUID: renstra.UUID,
	})

	return common.SuccessValue(renstra)
}

func GiveCodeAccessRenstra(
	prev *Renstra,
	uid uuid.UUID,
	kodeAkses string,
) common.ResultValue[*Renstra] {

	if prev == nil {
		return common.FailureValue[*Renstra](EmptyData())
	}

	if prev.UUID != uid {
		return common.FailureValue[*Renstra](InvalidData())
	}

	prev.KodeAkses = &kodeAkses

	prev.Raise(event.RenstraGiveCodeAccessEvent{
		EventID:     uuid.New(),
		OccurredOn:  time.Now().UTC(),
		RenstraUUID: prev.UUID,
	})

	return common.SuccessValue(prev)
}

func UpdateRenstra(
	prev *Renstra,
	uid uuid.UUID,
	tahun string,
	fakultasUnit uint,
	periodeUploadMulai, periodeUploadAkhir string,
	periodeDokumenMulai, periodeDokumenAkhir string,
	periodeLapanganMulai, periodeLapanganAkhir string,
	auditee, auditor1, auditor2 uint,
) common.ResultValue[*Renstra] {

	if prev == nil {
		return common.FailureValue[*Renstra](EmptyData())
	}

	if prev.UUID != uid {
		return common.FailureValue[*Renstra](InvalidData())
	}

	if fakultasUnit <= 0 {
		return common.FailureValue[*Renstra](InvalidFakultasUnit())
	}

	if auditee <= 0 {
		return common.FailureValue[*Renstra](MissingAuditee())
	}
	if auditor1 <= 0 {
		return common.FailureValue[*Renstra](MissingAuditor1())
	}
	if auditor2 <= 0 {
		return common.FailureValue[*Renstra](MissingAuditor2())
	}

	if auditee == auditor1 || auditee == auditor2 || auditor1 == auditor2 {
		return common.FailureValue[*Renstra](DuplicateAssigment())
	}

	format := "2006-01-02"

	// Periode Upload
	startUpload, err := time.Parse(format, periodeUploadMulai)
	if err != nil {
		return common.FailureValue[*Renstra](InvalidDate("upload"))
	}
	endUpload, err := time.Parse(format, periodeUploadAkhir)
	if err != nil {
		return common.FailureValue[*Renstra](InvalidDate("upload"))
	}

	// Periode Assessment Dokumen
	startDokumen, err := time.Parse(format, periodeDokumenMulai)
	if err != nil {
		return common.FailureValue[*Renstra](InvalidDate("assessment dokumen"))
	}
	endDokumen, err := time.Parse(format, periodeDokumenAkhir)
	if err != nil {
		return common.FailureValue[*Renstra](InvalidDate("assessment dokumen"))
	}

	// Periode Assessment Lapangan
	startLapangan, err := time.Parse(format, periodeLapanganMulai)
	if err != nil {
		return common.FailureValue[*Renstra](InvalidDate("assessment lapangan"))
	}
	endLapangan, err := time.Parse(format, periodeLapanganAkhir)
	if err != nil {
		return common.FailureValue[*Renstra](InvalidDate("assessment lapangan"))
	}

	// Cek overlap
	if isOverlap(startUpload, endUpload, startDokumen, endDokumen) {
		return common.FailureValue[*Renstra](PeriodOverlapUploadDokumen())
	}

	if isOverlap(startUpload, endUpload, startLapangan, endLapangan) {
		return common.FailureValue[*Renstra](PeriodOverlapUploadLapangan())
	}

	if isOverlap(startDokumen, endDokumen, startLapangan, endLapangan) {
		return common.FailureValue[*Renstra](PeriodOverlapDokumenLapangan())
	}

	// always overwrite (required field)
	prev.Tahun = tahun
	prev.FakultasUnit = fakultasUnit
	prev.PeriodeUploadMulai = periodeUploadMulai
	prev.PeriodeUploadAkhir = periodeUploadAkhir
	prev.PeriodeAssesmentDokumenMulai = periodeDokumenMulai
	prev.PeriodeAssesmentDokumenAkhir = periodeDokumenAkhir
	prev.PeriodeAssesmentLapanganMulai = periodeLapanganMulai
	prev.PeriodeAssesmentLapanganAkhir = periodeLapanganAkhir
	prev.Auditee = auditee
	prev.Auditor1 = auditor1
	prev.Auditor2 = auditor2

	prev.Raise(event.RenstraUpdatedEvent{
		EventID:     uuid.New(),
		OccurredOn:  time.Now().UTC(),
		RenstraUUID: prev.UUID,
	})

	return common.SuccessValue(prev)
}
