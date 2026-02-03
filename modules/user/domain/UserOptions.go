package domain

import (
	"github.com/google/uuid"
)

type UserOptions struct {
	ID       uint
	UUID     uuid.UUID
	Name     string
	Level    string
	Fakultas string
}
