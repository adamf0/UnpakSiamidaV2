package event

import (
	"time"

	"github.com/google/uuid"
	// "UnpakSiamida/common/domain"
)

type KtsPdfRequestedEvent struct {
	EventID    uuid.UUID
	KtsUUID    uuid.UUID
	Token      string
	OccurredOn time.Time
}

func (e KtsPdfRequestedEvent) ID() string {
	return e.EventID.String()
}

func (e KtsPdfRequestedEvent) OccurredOnUTC() time.Time {
	return e.OccurredOn
}
