package GetAktivitasByMatrix

import (
	domain "UnpakSiamida/modules/aktivitasv2/domain"
	"context"
)

type GetAktivitasByMatrixQueryHandler struct {
	Repo domain.IAktivitasV2Repository
}

func (h *GetAktivitasByMatrixQueryHandler) Handle(ctx context.Context, query GetAktivitasByMatrixQuery) (*domain.AktivitasByMatrixResult, error) {
	proker, err := h.Repo.GetMatrixProkerByUuid(ctx, query.UUID)
	if err != nil {
		return nil, err
	}
	if proker == nil {
		return nil, domain.ProkerNotFound()
	}

	activities, err := h.Repo.GetAktivitasByMatrikId(ctx, proker.ID)
	if err != nil {
		return nil, err
	}

	return &domain.AktivitasByMatrixResult{
		Proker: proker,
		Data:   activities,
		Total:  len(activities),
	}, nil
}
