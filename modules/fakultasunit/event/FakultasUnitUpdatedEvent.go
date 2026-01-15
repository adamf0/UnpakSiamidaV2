package event

import (
	"time"

	"github.com/google/uuid"
	// "UnpakSiamida/common/domain"
)

type FakultasUnitUpdatedEvent struct {
	EventID          uuid.UUID
	FakultasUnitUUID uuid.UUID
	OccurredOn       time.Time
}

func (e FakultasUnitUpdatedEvent) ID() string {
	return e.EventID.String()
}

func (e FakultasUnitUpdatedEvent) OccurredOnUTC() time.Time {
	return e.OccurredOn
}
