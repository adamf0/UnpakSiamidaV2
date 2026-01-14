package event

import (
	"bytes"
	"html/template"
	"time"

	"github.com/google/uuid"
)

type KtsCreatedView struct {
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
}

func RenderKtsCreatedTemplate(e KtsCreatedEvent) string {
	view := KtsCreatedView{
		KtsUUID:    e.KtsUUID,
		OccurredOn: e.OccurredOn,

		TemplateRenstra: e.TemplateRenstra,
		TemplateDokumen: e.TemplateDokumen,

		Standar:    e.Standar,
		Indikator:  e.Indikator,
		Pertanyaan: e.Pertanyaan,
		JenisFile:  e.JenisFile,

		Tahun:  &e.Tahun,
		Target: &e.Target,
		Status: &e.Status,
	}

	tpl := template.Must(template.New("kts").Parse(ktsCreatedTemplate))

	var buf bytes.Buffer
	_ = tpl.Execute(&buf, view)

	return buf.String()
}
