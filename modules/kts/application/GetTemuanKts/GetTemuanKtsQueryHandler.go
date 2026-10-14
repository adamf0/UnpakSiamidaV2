package application

import (
	domainkts "UnpakSiamida/modules/kts/domain"
	"context"
)

type GetTemuanKtsQueryHandler struct {
	Repo domainkts.IKtsRepository
}

func (h *GetTemuanKtsQueryHandler) Handle(
	ctx context.Context,
	query GetTemuanKtsQuery,
) ([]domainkts.TemuanKts, error) {
	return h.Repo.GetTemuanKts(ctx, query.RenstraUUID, query.Tugas, query.Tahun)
}
