package domain

import (
	"github.com/google/uuid"
	// "database/sql"
)

type GenerateRenstraDefault struct {
	ID              				uint
	UUID            				uuid.UUID
	RenstraId            			uint
	TemplateRenstraUuid           	uuid.UUID
	TemplateRenstra           		uint
	Indikator           			string
	Tugas  							string
	FakultasUnit      				uint
	RenstraTahun      				string
	TemplateTahun      				string
}