package GetAktivitasByMatrix

import (
	domain "UnpakSiamida/modules/prokerv3/domain"
	"context"
	"errors"
)

type GetAktivitasByMatrixQueryHandler struct {
	Repo domain.IProkerV3Repository
}

func (h *GetAktivitasByMatrixQueryHandler) Handle(ctx context.Context, query GetAktivitasByMatrixQuery) (*domain.AktivitasByMatrixResult, error) {
	proker, err := h.Repo.GetDefaultByUuid(ctx, query.UUID)
	if err != nil {
		return nil, err
	}
	if proker == nil {
		return nil, errors.New("Program kerja tidak ditemukan")
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
