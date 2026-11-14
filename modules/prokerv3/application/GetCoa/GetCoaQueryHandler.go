package GetCoa

import (
	domain "UnpakSiamida/modules/prokerv3/domain"
	"context"
)

type GetCoaQueryHandler struct {
	Repo domain.IProkerV3Repository
}

func (h *GetCoaQueryHandler) Handle(ctx context.Context, q GetCoaQuery) ([]domain.CoaItem, error) {
	return h.Repo.GetCoa(ctx, q.Search)
}
