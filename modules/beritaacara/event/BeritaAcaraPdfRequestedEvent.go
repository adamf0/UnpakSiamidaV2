package event

import (
	"time"

	"github.com/google/uuid"
	// "UnpakSiamida/common/domain"
)

type BeritaAcaraPdfRequestedEvent struct {
	EventID         uuid.UUID
	BeritaAcaraUUID uuid.UUID
	Token           string
	OccurredOn      time.Time
}

func (e BeritaAcaraPdfRequestedEvent) ID() string {
	return e.EventID.String()
}

func (e BeritaAcaraPdfRequestedEvent) OccurredOnUTC() time.Time {
	return e.OccurredOn
}
