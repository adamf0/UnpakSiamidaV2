package event

import (
	"time"

	"github.com/google/uuid"
	// "UnpakSiamida/common/domain"
)

type FakultasUnitCreatedEvent struct {
	EventID          uuid.UUID
	FakultasUnitUUID uuid.UUID
	OccurredOn       time.Time
}

func (e FakultasUnitCreatedEvent) ID() string {
	return e.EventID.String()
}

func (e FakultasUnitCreatedEvent) OccurredOnUTC() time.Time {
	return e.OccurredOn
}
