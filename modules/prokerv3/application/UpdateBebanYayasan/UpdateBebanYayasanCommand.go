package UpdateBebanYayasan

import (
	"github.com/google/uuid"
)

type UpdateBebanYayasanCommand struct {
	UUID         uuid.UUID
	BebanYayasan *bool
}

type UpdateBebanYayasanResult struct {
	Message      string `json:"message"`
	UUID         string `json:"uuid"`
	BebanYayasan bool   `json:"beban_yayasan"`
}
