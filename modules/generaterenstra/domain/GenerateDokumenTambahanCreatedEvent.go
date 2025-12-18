package domain

import (
	"time"
	"github.com/google/uuid"
	// "UnpakSiamida/common/domain"
)

type GenerateDokumenTambahanCreatedEvent struct {
	EventID    uuid.UUID
	GenerateDokumenTambahanUUID    uuid.UUID
	OccurredOn time.Time
}

func (e GenerateDokumenTambahanCreatedEvent) ID() string {
	return e.EventID.String()
}

func (e GenerateDokumenTambahanCreatedEvent) OccurredOnUTC() time.Time {
	return e.OccurredOn
}