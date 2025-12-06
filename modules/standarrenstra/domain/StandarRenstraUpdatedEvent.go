package domain

import (
	"time"
	"github.com/google/uuid"
	// "UnpakSiamida/common/domain"
)

type StandarRenstraUpdatedEvent struct {
	EventID    uuid.UUID
	StandarRenstraUUID    uuid.UUID
	OccurredOn time.Time
}

func (e StandarRenstraUpdatedEvent) ID() string {
	return e.EventID.String()
}

func (e StandarRenstraUpdatedEvent) OccurredOnUTC() time.Time {
	return e.OccurredOn
}