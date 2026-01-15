package event

import (
	"time"

	"github.com/google/uuid"
	// "UnpakSiamida/common/domain"
)

type GenerateRenstraCreatedEvent struct {
	EventID             uuid.UUID
	GenerateRenstraUUID uuid.UUID
	OccurredOn          time.Time
}

func (e GenerateRenstraCreatedEvent) ID() string {
	return e.EventID.String()
}

func (e GenerateRenstraCreatedEvent) OccurredOnUTC() time.Time {
	return e.OccurredOn
}
