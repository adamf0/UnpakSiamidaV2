package domain

import (
	"time"
	"github.com/google/uuid"
	// "UnpakSiamida/common/domain"
)

type JenisFileUpdatedEvent struct {
	EventID    uuid.UUID
	JenisFileUUID    uuid.UUID
	OccurredOn time.Time
}

func (e JenisFileUpdatedEvent) ID() string {
	return e.EventID.String()
}

func (e JenisFileUpdatedEvent) OccurredOnUTC() time.Time {
	return e.OccurredOn
}