package domain

import (
	"time"
	"github.com/google/uuid"
	// "UnpakSiamida/common/domain"
)

type RenstraGiveCodeAccessEvent struct {
	EventID    uuid.UUID
	RenstraUUID    uuid.UUID
	OccurredOn time.Time
}

func (e RenstraGiveCodeAccessEvent) ID() string {
	return e.EventID.String()
}

func (e RenstraGiveCodeAccessEvent) OccurredOnUTC() time.Time {
	return e.OccurredOn
}