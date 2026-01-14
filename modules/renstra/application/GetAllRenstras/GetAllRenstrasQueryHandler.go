package application

import (
    "context"
    domainrenstra "UnpakSiamida/modules/renstra/domain"
    "time"
)

type GetAllRenstrasQueryHandler struct {
    Repo domainrenstra.IRenstraRepository
}

func (h *GetAllRenstrasQueryHandler) Handle(
    ctx context.Context,
    q GetAllRenstrasQuery,
) (domainrenstra.PagedRenstras, error) {
    ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

    renstras, total, err := h.Repo.GetAll(
        ctx,
        q.Search,
        q.SearchFilters,
        q.Page,
        q.Limit,
        q.Scope,
    )
    if err != nil {
        return domainrenstra.PagedRenstras{}, err
    }

    currentPage := 1
    totalPages := 1

    if q.Page != nil {
        currentPage = *q.Page
    }
    if q.Limit != nil && *q.Limit > 0 {
        totalPages = int((total + int64(*q.Limit) - 1) / int64(*q.Limit))
    }

    return domainrenstra.PagedRenstras{
        Data:  renstras,
        Total: total,
        CurrentPage: currentPage,
        TotalPages:  totalPages,
    }, nil
}