package GetSankeyData

import (
	domain "UnpakSiamida/modules/prokerv3/domain"
	"context"
)

type GetSankeyDataQueryHandler struct {
	Repo domain.IProkerV3Repository
}

func (h *GetSankeyDataQueryHandler) Handle(ctx context.Context, q GetSankeyDataQuery) (*domain.SankeyDataResult, error) {
	return h.Repo.GetSankeyData(ctx, q.Tahun)
}
