package CreateJadwalProkerV3

import (
	domain "UnpakSiamida/modules/prokerv3/domain"
	"time"
)

type CreateJadwalProkerV3Command struct {
	Tahun          string
	TanggalMulai   time.Time
	TanggalAkhir   time.Time
	IdFakultasUnit *uint64
}

type CreateJadwalProkerV3Result struct {
	Message string                `json:"message"`
	UUID    string                `json:"uuid"`
	Data    domain.JadwalProkerV3 `json:"data"`
}
