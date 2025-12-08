package application

import (
    "context"

    domainTahunRenstra "UnpakSiamida/modules/tahunrenstra/domain"
    "errors"
    "gorm.io/gorm"
)

type GetActiveTahunRenstraQueryHandler struct {
    Repo domainTahunRenstra.ITahunRenstraRepository
}

func (h *GetActiveTahunRenstraQueryHandler) Handle(
    ctx context.Context,
    q GetActiveTahunRenstraQuery,
) (*domainTahunRenstra.TahunRenstra, error) {

    tahunrenstra, err := h.Repo.GetActive(ctx)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domainTahunRenstra.EmptyData()
		}
		return nil, err
	}

    return tahunrenstra, nil
}
