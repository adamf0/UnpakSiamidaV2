package event

import (
	"time"
	"github.com/google/uuid"
	// "UnpakSiamida/common/domain"
)

type KtsCreatedEvent struct {
	EventID    uuid.UUID
	KtsUUID    uuid.UUID
	OccurredOn time.Time

	// === Context ===
	TemplateRenstra  *uint
	TemplateDokumen  *uint

	// Renstra
	Standar    *string
	Indikator  *string

	// Dokumen
	Pertanyaan *string
	JenisFile *string
	Kts  *string

	// Shared
	Tahun    string
	IdTarget uint
	Target   string
	Status   string
}

func (e KtsCreatedEvent) ID() string {
	return e.EventID.String()
}

func (e KtsCreatedEvent) GetTemplateRenstra() *uint {
	return e.TemplateRenstra
}

func (e KtsCreatedEvent) GetTemplateDokumen() *uint {
	return e.TemplateDokumen
}

func (e KtsCreatedEvent) GetTahun() string {
	return e.Tahun
}

func (e KtsCreatedEvent) GetIdTarget() uint {
	return e.IdTarget
}

func (e KtsCreatedEvent) OccurredOnUTC() time.Time {
	return e.OccurredOn
}