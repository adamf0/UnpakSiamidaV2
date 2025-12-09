package domain

import (
	"time"
	"github.com/google/uuid"
	// "UnpakSiamida/common/domain"
)

type TemplateDokumenTambahanCreatedEvent struct {
	EventID    uuid.UUID
	TemplateDokumenTambahanUUID    uuid.UUID
	OccurredOn time.Time
}

func (e TemplateDokumenTambahanCreatedEvent) ID() string {
	return e.EventID.String()
}

func (e TemplateDokumenTambahanCreatedEvent) OccurredOnUTC() time.Time {
	return e.OccurredOn
}