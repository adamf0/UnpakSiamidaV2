package GetUnitPaguSummary

import (
	domain "UnpakSiamida/modules/prokerv3/domain"
	"context"
)

type GetUnitPaguSummaryQueryHandler struct {
	Repo domain.IProkerV3Repository
}

func (h *GetUnitPaguSummaryQueryHandler) Handle(ctx context.Context, q GetUnitPaguSummaryQuery) (*domain.UnitPaguSummaryResult, error) {
	return h.Repo.GetUnitPaguSummary(ctx, q.Tahun, q.ReqFU, q.UserFU)
}
