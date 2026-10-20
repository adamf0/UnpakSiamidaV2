package GetProkerV3Options

import (
	domain "UnpakSiamida/modules/prokerv3/domain"
	"context"
)

type GetProkerV3OptionsQueryHandler struct {
	Repo domain.IProkerV3Repository
}

func (h *GetProkerV3OptionsQueryHandler) Handle(ctx context.Context, query GetProkerV3OptionsQuery) (*domain.ProkerOptionsResult, error) {
	return h.Repo.GetOptions(ctx)
}
