package event

import (
	"time"

	"github.com/google/uuid"
	// "UnpakSiamida/common/domain"
)

type KtsUpdatedEvent struct {
	EventID    uuid.UUID
	KtsUUID    uuid.UUID
	OccurredOn time.Time

	// === Context ===
	TemplateRenstra *uint
	TemplateDokumen *uint

	// Renstra
	StandarId   *uint
	Standar     *string
	IndikatorId *uint
	Indikator   *string

	// Dokumen
	Pertanyaan  *string
	JenisFileId *uint
	JenisFile   *string
	Kts         *string

	// Shared
	Tahun    *string
	IdTarget *uint
	Target   *string
	Status   *string

	//auditor
	NomorLaporan     *string
	TanggalLaporan   *string
	KetidaksesuaianP *string
	KetidaksesuaianL *string
	KetidaksesuaianO *string
	KetidaksesuaianR *string
	AkarMasalah      *string
	TindakanKoreksi  *string

	//auditee
	StatusAccAuditee  *uint
	KeteranganTolak   *string
	TindakanPerbaikan *string

	//auditor
	TanggalPenyelesaian *string

	//auditee
	TinjauanTindakanPerbaikan *string
	TanggalClosing            *string

	//auditor
	TanggalClosingFinal *string
	WmmUpmfUpmps        *string
}

func (e KtsUpdatedEvent) ID() string {
	return e.EventID.String()
}

func (e KtsUpdatedEvent) GetTemplateRenstra() *uint {
	return e.TemplateRenstra
}

func (e KtsUpdatedEvent) GetTemplateDokumen() *uint {
	return e.TemplateDokumen
}

func (e KtsUpdatedEvent) GetTahun() *string {
	return e.Tahun
}

func (e KtsUpdatedEvent) GetIdTarget() *uint {
	return e.IdTarget
}

func (e KtsUpdatedEvent) OccurredOnUTC() time.Time {
	return e.OccurredOn
}
