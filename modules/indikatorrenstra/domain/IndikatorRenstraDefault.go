package domain

import (
	"github.com/google/uuid"
	"database/sql"
)

type IndikatorRenstraDefault struct {
		Id          uint
		Uuid        uuid.UUID
		Indikator   string
		Standar     sql.NullInt64
		UuidStandar sql.NullString
		Parent      sql.NullInt64
		UuidParent  sql.NullString
		Tahun       string
		TipeTarget  string
		Operator    sql.NullString
}