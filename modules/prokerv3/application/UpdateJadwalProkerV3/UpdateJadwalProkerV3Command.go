package UpdateJadwalProkerV3

import (
	domain "UnpakSiamida/modules/prokerv3/domain"
	"time"

	"github.com/google/uuid"
)

type UpdateJadwalProkerV3Command struct {
	UUID           uuid.UUID
	Tahun          string
	TanggalMulai   *time.Time
	TanggalAkhir   *time.Time
	HasFuKey       bool
	IdFakultasUnit *uint64
}

type UpdateJadwalProkerV3Result struct {
	Message string                `json:"message"`
	UUID    string                `json:"uuid"`
	Data    domain.JadwalProkerV3 `json:"data"`
}
