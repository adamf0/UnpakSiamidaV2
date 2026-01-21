package domain

import (
	"time"

	"github.com/google/uuid"
)

type BeritaAcaraDefault struct {
	Id   uint
	UUID uuid.UUID

	Tahun          string
	FakultasUnitId int
	FakultasUnit   string
	Tanggal        time.Time
	AuditeeId      *int
	Auditee        *string
	Auditor1Id     *int
	Auditor1       *string
	Auditor2ID     *int
	Auditor2       *string
}
