package GetHistoryStateV3

import (
	domain "UnpakSiamida/modules/prokerv3/domain"
	"context"
)

type GetHistoryStateV3QueryHandler struct {
	Repo domain.IProkerV3Repository
}

func (h *GetHistoryStateV3QueryHandler) Handle(ctx context.Context, query GetHistoryStateV3Query) ([]domain.HistoryState, error) {
	return h.Repo.GetHistoryState(ctx, query.UUID)
}
