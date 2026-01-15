package event

import (
	"time"

	"github.com/google/uuid"
	// "UnpakSiamida/common/domain"
)

type TemplateRenstraCreatedEvent struct {
	EventID             uuid.UUID
	TemplateRenstraUUID uuid.UUID
	OccurredOn          time.Time
}

func (e TemplateRenstraCreatedEvent) ID() string {
	return e.EventID.String()
}

func (e TemplateRenstraCreatedEvent) OccurredOnUTC() time.Time {
	return e.OccurredOn
}
