package domain

import (
	"github.com/google/uuid"
)

type RenstraNilaiDefault struct {
	// Renstra
	RenstraID    int
	RenstraUUID  uuid.UUID
	TahunRenstra string

	// Renstra Nilai
	ID   int
	UUID string

	// Fakultas / Unit
	TargetAudit string
	Jenjang     *string
	Fakultas    *string
	Type        string

	// Standar
	StandarID   int
	StandarUUID uuid.UUID
	NamaStandar string

	// Indikator
	IndikatorID    int
	IndikatorUUID  uuid.UUID
	NamaIndikator  string
	TahunIndikator string
	TipeTarget     string
	Operator       string

	// Template Renstra
	TemplateRenstraId   int
	TemplateRenstraUUID uuid.UUID
	Satuan              string
	Target              *string
	TargetMin           *string
	TargetMax           *string
	TugasTemplate       string
	TahunTemplate       string
	IsPertanyaan        bool

	// Nilai
	Tugas          string
	CapaianAuditee *string
	CatatanAuditee *string
	LinkBukti      *string
	CapaianAuditor *string
	CatatanAuditor *string
}
