package GetTargetUnpak

import (
	domain "UnpakSiamida/modules/prokerv3/domain"
	"context"
)

type GetTargetUnpakQueryHandler struct {
	Repo domain.IProkerV3Repository
}

func (h *GetTargetUnpakQueryHandler) Handle(ctx context.Context, query GetTargetUnpakQuery) (*domain.TargetUnpakResult, error) {
	return h.Repo.GetTargetUnpak(ctx, query.Tahun, query.IndikatorParam, query.UnitParam)
}
