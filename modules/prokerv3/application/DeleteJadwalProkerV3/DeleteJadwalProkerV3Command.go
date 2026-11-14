package DeleteJadwalProkerV3

import (
	"github.com/google/uuid"
)

type DeleteJadwalProkerV3Command struct {
	UUID uuid.UUID
}

type DeleteJadwalProkerV3Result struct {
	Message string `json:"message"`
	UUID    string `json:"uuid"`
}
