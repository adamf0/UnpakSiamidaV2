package application

import (
	"context"
	"time"

	domainindikatorrenstra "UnpakSiamida/modules/indikatorrenstra/domain"
)

type GetAllMasterIkuQueryHandler struct {
	Repo domainindikatorrenstra.IIndikatorRenstraRepository
}

func (h *GetAllMasterIkuQueryHandler) Handle(
	ctx context.Context,
	q GetAllMasterIkuQuery,
) ([]domainindikatorrenstra.MasterIku, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	return h.Repo.GetAllMasterIku(ctx)
}
