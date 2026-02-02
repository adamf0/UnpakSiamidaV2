package domain

import "github.com/google/uuid"

type IndikatorTree struct {
	IndikatorId       int
	IndikatorUuid     uuid.UUID
	Indikator         string
	ParentIndikatorId *int
	Pointing          string
}
