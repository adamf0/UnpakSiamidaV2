package GetAktivitasDetail

import (
	domain "UnpakSiamida/modules/aktivitasv2/domain"
	"context"
)

type GetAktivitasDetailQueryHandler struct {
	Repo domain.IAktivitasV2Repository
}

func (h *GetAktivitasDetailQueryHandler) Handle(ctx context.Context, query GetAktivitasDetailQuery) (*domain.AktivitasDetailResult, error) {
	act, err := h.Repo.GetAktivitasByUuid(ctx, query.UUID)
	if err != nil {
		return nil, err
	}
	if act == nil {
		return nil, domain.AktivitasNotFound()
	}

	var proker *domain.MatrixProkerSummary
	if act.IdMatrikProker != nil && *act.IdMatrikProker > 0 {
		proker, _ = h.Repo.GetMatrixProkerById(ctx, uint(*act.IdMatrikProker))
	}

	return &domain.AktivitasDetailResult{
		Aktivitas: act,
		Proker:    proker,
		Timelines: act.Timelines,
		Data:      act.Timelines,
	}, nil
}
