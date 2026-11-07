package event

import (
	"time"

	"github.com/google/uuid"
)

type RenstraEmailRecipient struct {
	Role  string `json:"role"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type RenstraSendEmailEvent struct {
	EventID     uuid.UUID               `json:"event_id"`
	RenstraUUID uuid.UUID               `json:"renstra_uuid"`
	Tahun       string                  `json:"tahun"`
	TargetAudit string                  `json:"target_audit"`
	KodeAkses   string                  `json:"kode_akses"`
	Recipients  []RenstraEmailRecipient `json:"recipients"`
	OccurredOn  time.Time               `json:"occurred_on"`
}

func (e RenstraSendEmailEvent) ID() string {
	return e.EventID.String()
}

func (e RenstraSendEmailEvent) OccurredOnUTC() time.Time {
	return e.OccurredOn
}
