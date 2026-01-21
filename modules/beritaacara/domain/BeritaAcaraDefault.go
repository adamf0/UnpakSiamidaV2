package domain

import (
	"time"

	"github.com/google/uuid"
)

type BeritaAcaraDefault struct {
	Id   uint
	UUID uuid.UUID

	Tahun        string
	FakultasUnit int
	Tanggal      time.Time
	Auditee      *int
	NamaAuditee  *string
	Auditor1     *int
	NamaAuditor1 *string
	Auditor2     *int
	NamaAuditor2 *string
}
