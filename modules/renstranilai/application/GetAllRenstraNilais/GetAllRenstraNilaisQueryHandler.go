package application

import (
    "context"
    domainrenstranilai "UnpakSiamida/modules/renstranilai/domain"
    "time"
)

type GetAllRenstraNilaisQueryHandler struct {
    Repo domainrenstranilai.IRenstraNilaiRepository
}

func (h *GetAllRenstraNilaisQueryHandler) Handle(
    ctx context.Context,
    q GetAllRenstraNilaisQuery,
) (domainrenstranilai.PagedRenstraNilais, error) {
    ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

    renstranilais, total, err := h.Repo.GetAll(
        ctx,
        q.Search,
        q.SearchFilters,
        q.Page,
        q.Limit,
    )
    if err != nil {
        return domainrenstranilai.PagedRenstraNilais{}, err
    }

    currentPage := 1
    totalPages := 1

    if q.Page != nil {
        currentPage = *q.Page
    }
    if q.Limit != nil && *q.Limit > 0 {
        totalPages = int((total + int64(*q.Limit) - 1) / int64(*q.Limit))
    }

    return domainrenstranilai.PagedRenstraNilais{
        Data:  renstranilais,
        Total: total,
        CurrentPage: currentPage,
        TotalPages:  totalPages,
    }, nil
}