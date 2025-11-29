package domain

import (
	"time"
	"github.com/google/uuid"
	// "UnpakSiamida/common/domain"
)

type AkurasiPenelitianCreatedEvent struct {
	EventID    uuid.UUID
	OccurredOn time.Time
}

func (e AkurasiPenelitianCreatedEvent) ID() string {
	return e.EventID.String()
}

func (e AkurasiPenelitianCreatedEvent) OccurredOnUTC() time.Time {
	return e.OccurredOn
}
