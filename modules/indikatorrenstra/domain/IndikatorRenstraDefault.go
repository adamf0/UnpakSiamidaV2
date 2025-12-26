package domain

import (
	"github.com/google/uuid"
)

type IndikatorRenstraDefault struct {
		Id          uint
		Uuid        uuid.UUID
		Indikator   string
		StandarID   *uint
		Standar     *string
		UuidStandar *string
		Parent      *uint
		UuidParent  *string
		Tahun       string
		TipeTarget  string
		Operator    *string
}