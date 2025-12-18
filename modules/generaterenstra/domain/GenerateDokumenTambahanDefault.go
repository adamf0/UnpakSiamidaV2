package domain

import (
	"github.com/google/uuid"
	// "database/sql"
)

type GenerateDokumenTambahanDefault struct {
	ID              				uint
	UUID            				uuid.UUID
	RenstraId            			uint
	TemplateDokumenTambahanUuid     uuid.UUID
	TemplateDokumenTambahan         uint
	JenisFileId           			uint
	JenisFile           			string
	Pertanyaan           			string
	Tugas  							string
	FakultasUnit      				uint
	RenstraTahun      				string
	TemplateTahun      				string
}