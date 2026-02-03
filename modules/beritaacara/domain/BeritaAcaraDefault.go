package domain

import (
	"time"

	"github.com/google/uuid"
)

type BeritaAcaraDefault struct {
	Id   uint
	UUID uuid.UUID

	Tahun            string
	FakultasUnitId   int
	FakultasUnitUuid uuid.UUID
	FakultasUnit     string
	Tanggal          time.Time

	AuditeeId   *int
	Auditee     *string
	AuditeeUuid *uuid.UUID

	Auditor1Id   *int
	Auditor1     *string
	Auditor1Uuid *uuid.UUID

	Auditor2Id   *int
	Auditor2     *string
	Auditor2Uuid *uuid.UUID
}
