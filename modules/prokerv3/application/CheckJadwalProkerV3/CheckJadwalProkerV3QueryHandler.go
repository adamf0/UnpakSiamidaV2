package CheckJadwalProkerV3

import (
	domain "UnpakSiamida/modules/prokerv3/domain"
	"context"
)

type CheckJadwalProkerV3QueryHandler struct {
	Repo domain.IProkerV3Repository
}

func (h *CheckJadwalProkerV3QueryHandler) Handle(ctx context.Context, q CheckJadwalProkerV3Query) (*domain.CheckJadwalResult, error) {
	return h.Repo.CheckJadwalProkerV3(ctx, q.Request)
}
