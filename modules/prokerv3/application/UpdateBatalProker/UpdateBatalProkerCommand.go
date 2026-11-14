package UpdateBatalProker

import (
	"github.com/google/uuid"
)

type UpdateBatalProkerCommand struct {
	UUID         uuid.UUID
	IsBatal      *bool
	CatatanBatal string
}

type UpdateBatalProkerResult struct {
	Message      string `json:"message"`
	UUID         string `json:"uuid"`
	IsBatal      bool   `json:"is_batal"`
	CatatanBatal string `json:"catatan_batal"`
}
