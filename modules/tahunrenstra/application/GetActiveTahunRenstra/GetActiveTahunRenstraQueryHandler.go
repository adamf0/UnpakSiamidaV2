package application

import (
    "context"

    domainTahunRenstra "UnpakSiamida/modules/tahunrenstra/domain"
)

type GetActiveTahunRenstraQueryHandler struct {
    Repo domainTahunRenstra.ITahunRenstraRepository
}

func (h *GetActiveTahunRenstraQueryHandler) Handle(
    ctx context.Context,
    q GetActiveTahunRenstraQuery,
) (*domainTahunRenstra.TahunRenstra, error) {

    return h.Repo.GetActive(ctx)
}
