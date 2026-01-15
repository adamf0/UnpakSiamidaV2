package event

import (
	"time"

	"github.com/google/uuid"
	// "UnpakSiamida/common/domain"
)

type StandarRenstraCreatedEvent struct {
	EventID            uuid.UUID
	StandarRenstraUUID uuid.UUID
	OccurredOn         time.Time
}

func (e StandarRenstraCreatedEvent) ID() string {
	return e.EventID.String()
}

func (e StandarRenstraCreatedEvent) OccurredOnUTC() time.Time {
	return e.OccurredOn
}
