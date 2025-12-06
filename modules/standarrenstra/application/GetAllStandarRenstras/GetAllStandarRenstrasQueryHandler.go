package application

import (
    "context"
    domainstandarrenstra "UnpakSiamida/modules/standarrenstra/domain"
)

type GetAllStandarRenstrasQueryHandler struct {
    Repo domainstandarrenstra.IStandarRenstraRepository
}

func (h *GetAllStandarRenstrasQueryHandler) Handle(
    ctx context.Context,
    q GetAllStandarRenstrasQuery,
) (domainstandarrenstra.PagedStandarRenstras, error) {

    standarrenstras, total, err := h.Repo.GetAll(
        ctx,
        q.Search,
        q.SearchFilters,
        q.Page,
        q.Limit,
    )
    if err != nil {
        return domainstandarrenstra.PagedStandarRenstras{}, err
    }

    currentPage := 1
    totalPages := 1

    if q.Page != nil {
        currentPage = *q.Page
    }
    if q.Limit != nil && *q.Limit > 0 {
        totalPages = int((total + int64(*q.Limit) - 1) / int64(*q.Limit))
    }

    return domainstandarrenstra.PagedStandarRenstras{
        Data:  standarrenstras,
        Total: total,
        CurrentPage: currentPage,
        TotalPages:  totalPages,
    }, nil
}