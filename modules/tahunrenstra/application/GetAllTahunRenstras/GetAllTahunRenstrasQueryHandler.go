package application

import (
    "context"
    domainTahunRenstra "UnpakSiamida/modules/tahunrenstra/domain"
    "time"
)

type GetAllTahunRenstrasQueryHandler struct {
    Repo domainTahunRenstra.ITahunRenstraRepository
}

func (h *GetAllTahunRenstrasQueryHandler) Handle(
    ctx context.Context,
    q GetAllTahunRenstrasQuery,
) (domainTahunRenstra.PagedTahunRenstras, error) {
    ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

    TahunRenstras, total, err := h.Repo.GetAll(
        ctx,
        q.Search,
        q.SearchFilters,
        q.Page,
        q.Limit,
    )
    if err != nil {
        return domainTahunRenstra.PagedTahunRenstras{}, err
    }

    currentPage := 1
    totalPages := 1

    if q.Page != nil {
        currentPage = *q.Page
    }
    if q.Limit != nil && *q.Limit > 0 {
        totalPages = int((total + int64(*q.Limit) - 1) / int64(*q.Limit))
    }

    return domainTahunRenstra.PagedTahunRenstras{
        Data:  TahunRenstras,
        Total: total,
        CurrentPage: currentPage,
        TotalPages:  totalPages,
    }, nil
}