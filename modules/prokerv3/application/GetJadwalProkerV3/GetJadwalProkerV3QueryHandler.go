package GetJadwalProkerV3

import (
	domain "UnpakSiamida/modules/prokerv3/domain"
	"context"
)

type GetJadwalProkerV3QueryHandler struct {
	Repo domain.IProkerV3Repository
}

func (h *GetJadwalProkerV3QueryHandler) Handle(ctx context.Context, q GetJadwalProkerV3Query) (*domain.JadwalProkerV3DTO, error) {
	return h.Repo.GetJadwalProkerV3(ctx, q.UUID)
}
