package application

import (
    "context"
    domainindikatorrenstra "UnpakSiamida/modules/indikatorrenstra/domain"
)

type GetAllIndikatorRenstrasQueryHandler struct {
    Repo domainindikatorrenstra.IIndikatorRenstraRepository
}

func (h *GetAllIndikatorRenstrasQueryHandler) Handle(
    ctx context.Context,
    q GetAllIndikatorRenstrasQuery,
) (domainindikatorrenstra.PagedIndikatorRenstras, error) {

    indikatorrenstras, total, err := h.Repo.GetAll(
        ctx,
        q.Search,
        q.SearchFilters,
        q.Page,
        q.Limit,
    )
    if err != nil {
        return domainindikatorrenstra.PagedIndikatorRenstras{}, err
    }

    currentPage := 1
    totalPages := 1

    if q.Page != nil {
        currentPage = *q.Page
    }
    if q.Limit != nil && *q.Limit > 0 {
        totalPages = int((total + int64(*q.Limit) - 1) / int64(*q.Limit))
    }

    return domainindikatorrenstra.PagedIndikatorRenstras{
        Data:  indikatorrenstras,
        Total: total,
        CurrentPage: currentPage,
        TotalPages:  totalPages,
    }, nil
}