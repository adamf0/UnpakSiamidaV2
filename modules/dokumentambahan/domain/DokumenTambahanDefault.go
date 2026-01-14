package domain

import (
	"github.com/google/uuid"
)

type DokumenTambahanDefault struct {
	ID              	uint
	UUID            	uuid.UUID
	RenstraId         	*uint
	RenstraUUID         uuid.UUID
	TemplateDokumenTambahanUUID uuid.UUID
	TargetAudit			string
	Jenjang				*string
	Fakultas			*string
	Type				string
	
	Pertanyaan  		string
	Dokumen       		string
	Kategori       		*string
	Klasifikasi       	string

	TahunRenstra				string
	TahunDokumenTambahan		string
	Tugas						string
	Link						*string
	CapaianAuditor				*string
	CatatanAuditor				*string
	// Version				*string
}