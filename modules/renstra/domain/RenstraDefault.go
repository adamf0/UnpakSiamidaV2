package domain

import (
	"github.com/google/uuid"
)

type RenstraDefault struct {
	ID     uint      `json:"ID"`
	UUID   uuid.UUID `json:"UUID"`
	Tahun  string    `json:"Tahun"`

	FakultasUnitId   uint        `json:"FakultasUnitId"`
	FakultasUnitUUID *uuid.UUID `json:"FakultasUnitUuid"`
	FakultasUnit     string     `json:"FakultasUnit"`

	PeriodeUploadMulai   string `json:"PeriodeUploadMulai"`
	PeriodeUploadAkhir   string `json:"PeriodeUploadAkhir"`
	PeriodeAssesmentDokumenMulai  string `json:"PeriodeAssesmentDokumenMulai"`
	PeriodeAssesmentDokumenAkhir  string `json:"PeriodeAssesmentDokumenAkhir"`
	PeriodeAssesmentLapanganMulai string `json:"PeriodeAssesmentLapanganMulai"`
	PeriodeAssesmentLapanganAkhir string `json:"PeriodeAssesmentLapanganAkhir"`

	KodeAkses *string `json:"KodeAkses"`

	AuditeeId  *uint `json:"AuditeeId"`
	Auditor1Id *uint `json:"Auditor1Id"`
	Auditor2Id *uint `json:"Auditor2Id"`

	AuditeeUuid  *uuid.UUID `json:"AuditeeUuid"`
	Auditor1Uuid *uuid.UUID `json:"Auditor1Uuid"`
	Auditor2Uuid *uuid.UUID `json:"Auditor2Uuid"`

	Auditee  *string `json:"Auditee"`
	Auditor1 *string `json:"Auditor1"`
	Auditor2 *string `json:"Auditor2"`

	Catatan1 *string `json:"Catatan1"`
	Catatan2 *string `json:"Catatan2"`

	Jenjang  *string `json:"Jenjang"`
	Type     string  `json:"Type"`
	Fakultas *string `json:"Fakultas"`
}