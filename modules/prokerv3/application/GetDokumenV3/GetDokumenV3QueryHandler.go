package GetDokumenV3

import (
	domain "UnpakSiamida/modules/prokerv3/domain"
	"context"
)

type GetDokumenV3QueryHandler struct {
	Repo domain.IProkerV3Repository
}

func (h *GetDokumenV3QueryHandler) Handle(ctx context.Context, query GetDokumenV3Query) (*domain.ProkerDokumenResult, error) {
	proker, err := h.Repo.GetDefaultByUuid(ctx, query.ProkerUUID)
	if err != nil {
		return nil, err
	}
	if proker == nil {
		return nil, domain.ProkerNotFound()
	}

	docs, err := h.Repo.GetDokumenByMatrikId(ctx, proker.ID)
	if err != nil {
		return nil, err
	}

	activities, _ := h.Repo.GetAktivitasOptionsByMatrikId(ctx, proker.ID)

	return &domain.ProkerDokumenResult{
		Proker:     proker,
		Activities: activities,
		Data:       docs,
	}, nil
}
