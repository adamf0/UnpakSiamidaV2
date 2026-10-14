package GetProkerV3

import (
	domain "UnpakSiamida/modules/prokerv3/domain"
	"context"
)

type GetProkerV3ByUuidQueryHandler struct {
	Repo domain.IProkerV3Repository
}

func (h *GetProkerV3ByUuidQueryHandler) Handle(
	ctx context.Context,
	query GetProkerV3ByUuidQuery,
) (*domain.ProkerV3Default, error) {
	return h.Repo.GetDefaultByUuid(ctx, query.UUID)
}
