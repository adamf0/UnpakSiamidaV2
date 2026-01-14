package domain

import (
	"github.com/google/uuid"
)

type RenstraNilaiDefault struct {
	ID              	uint
	UUID            	uuid.UUID
	RenstraId         	*uint
	RenstraUUID         uuid.UUID
	TemplateRenstraUUID uuid.UUID
	TargetAudit			string
	Jenjang				*string
	Fakultas			*string
	Type				string
	
	NamaStandarRenstra  *string
	Indikator       	*string
	Kategori       		string
	Klasifikasi       	string
	Satuan       		*string
	Target       		*string
	TargetMin       	*string
	TargetMax       	*string
	Operator       		*string

	TahunRenstra		string
	TahunIndikator		string
	Tugas				string
	CapaianAuditee		*string
	CatatanAuditee		*string
	LinkBukti			*string
	CapaianAuditor		*string
	CatatanAuditor		*string
	// Version				*string
}