package domain

import (
	"time"
	"github.com/google/uuid"
	// "UnpakSiamida/common/domain"
)

type TemplateRenstraUpdatedEvent struct {
	EventID    uuid.UUID
	TemplateRenstraUUID    uuid.UUID
	OccurredOn time.Time
}

func (e TemplateRenstraUpdatedEvent) ID() string {
	return e.EventID.String()
}

func (e TemplateRenstraUpdatedEvent) OccurredOnUTC() time.Time {
	return e.OccurredOn
}