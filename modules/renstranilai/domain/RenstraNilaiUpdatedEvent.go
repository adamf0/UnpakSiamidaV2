package domain

import (
	"time"
	"github.com/google/uuid"
	// "UnpakSiamida/common/domain"
)

type RenstraNilaiUpdatedEvent struct {
	EventID    uuid.UUID
	RenstraNilaiUUID    uuid.UUID
	OccurredOn time.Time
}

func (e RenstraNilaiUpdatedEvent) ID() string {
	return e.EventID.String()
}

func (e RenstraNilaiUpdatedEvent) OccurredOnUTC() time.Time {
	return e.OccurredOn
}