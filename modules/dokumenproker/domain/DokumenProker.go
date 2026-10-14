package domain

import (
	"slices"
	"strings"
	"time"

	common "UnpakSiamida/common/domain"
	event "UnpakSiamida/modules/dokumenproker/event"

	"github.com/google/uuid"
)

type DokumenProker struct {
	common.Entity

	ID           uint      `gorm:"primaryKey;autoIncrement"`
	UUID         uuid.UUID `gorm:"type:char(36);uniqueIndex"`
	MataProgram  uint      `gorm:"column:id_mata_program"`
	FakultasUnit uint      `gorm:"column:id_fakultas_unit"`
	JenisDokumen string    `gorm:"column:jenis_dokumen"`
	File         string    `gorm:"column:file"`
	Status       string    `gorm:"column:status_verifikasi"`
	Catatan      *string   `gorm:"column:catatan"`
}

func (DokumenProker) TableName() string {
	return "dokumen_realisasi_proker"
}

// === CREATE ===
func NewDokumenProker(mataprogram uint, fakultasunit uint, jenis_dokumen string, file string, status string, catatan *string) common.ResultValue[*DokumenProker] {
	if fakultasunit <= 0 {
		return common.FailureValue[*DokumenProker](NotFoundFakultas())
	}
	if mataprogram <= 0 {
		return common.FailureValue[*DokumenProker](NotFoundMataProgram())
	}
	if !isValidJenisDokumen(jenis_dokumen) {
		return common.FailureValue[*DokumenProker](InvalidJenisDokumen())
	}
	if status == "" {
		status = "belum_terverifikasi"
	}
	if !isValidStatus(status) {
		return common.FailureValue[*DokumenProker](InvalidStatus())
	}

	dokumenproker := &DokumenProker{
		UUID:         uuid.New(),
		FakultasUnit: fakultasunit,
		MataProgram:  mataprogram,
		JenisDokumen: strings.ToUpper(strings.TrimSpace(jenis_dokumen)),
		File:         file,
		Status:       strings.ToLower(strings.TrimSpace(status)),
		Catatan:      catatan,
	}

	dokumenproker.Raise(event.DokumenProkerCreatedEvent{
		EventID:           uuid.New(),
		OccurredOn:        time.Now().UTC(),
		DokumenProkerUUID: dokumenproker.UUID,
	})

	return common.SuccessValue(dokumenproker)
}

// === UPDATE ===
func UpdateDokumenProker(
	prev *DokumenProker,
	uid uuid.UUID,
	mataprogram uint,
	fakultasunit uint,
	jenis_dokumen string,
	file string,
	status string,
	catatan *string,
) common.ResultValue[*DokumenProker] {

	if prev == nil {
		return common.FailureValue[*DokumenProker](EmptyData())
	}

	if prev.UUID != uid {
		return common.FailureValue[*DokumenProker](InvalidData())
	}

	if fakultasunit <= 0 {
		return common.FailureValue[*DokumenProker](NotFoundFakultas())
	}
	if mataprogram <= 0 {
		return common.FailureValue[*DokumenProker](NotFoundMataProgram())
	}
	if !isValidJenisDokumen(jenis_dokumen) {
		return common.FailureValue[*DokumenProker](InvalidJenisDokumen())
	}
	if status != "" && !isValidStatus(status) {
		return common.FailureValue[*DokumenProker](InvalidStatus())
	}

	prev.FakultasUnit = fakultasunit
	prev.MataProgram = mataprogram
	prev.JenisDokumen = strings.ToUpper(strings.TrimSpace(jenis_dokumen))
	if file != "" {
		prev.File = file
	}
	if status != "" {
		prev.Status = strings.ToLower(strings.TrimSpace(status))
	}
	prev.Catatan = catatan

	prev.Raise(event.DokumenProkerUpdatedEvent{
		EventID:           uuid.New(),
		OccurredOn:        time.Now().UTC(),
		DokumenProkerUUID: prev.UUID,
	})

	return common.SuccessValue(prev)
}

func isValidJenisDokumen(jenis string) bool {
	j := strings.ToUpper(strings.TrimSpace(jenis))
	return slices.Contains([]string{"PROPOSAL/TOR", "SK", "LAPORAN", "SOP", "DOCUMENT PENDUKUNG", "DOKUMEN PENDUKUNG"}, j)
}

func isValidStatus(status string) bool {
	s := strings.ToLower(strings.TrimSpace(status))
	if s == "" {
		return true
	}
	return slices.Contains([]string{"belum_terverifikasi", "gagal_terverifikasi", "terverifikasi"}, s)
}
