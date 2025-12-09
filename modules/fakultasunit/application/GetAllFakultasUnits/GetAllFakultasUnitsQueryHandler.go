package application

import (
    "context"
    domainfakultasunit "UnpakSiamida/modules/fakultasunit/domain"
)

type GetAllFakultasUnitsQueryHandler struct {
    Repo domainfakultasunit.IFakultasUnitRepository
}

func (h *GetAllFakultasUnitsQueryHandler) Handle(
    ctx context.Context,
    q GetAllFakultasUnitsQuery,
) (domainfakultasunit.PagedFakultasUnits, error) {

    fakultasunits, total, err := h.Repo.GetAll(
        ctx,
        q.Search,
        q.SearchFilters,
        q.Page,
        q.Limit,
    )
    if err != nil {
        return domainfakultasunit.PagedFakultasUnits{}, err
    }

    currentPage := 1
    totalPages := 1

    if q.Page != nil {
        currentPage = *q.Page
    }
    if q.Limit != nil && *q.Limit > 0 {
        totalPages = int((total + int64(*q.Limit) - 1) / int64(*q.Limit))
    }

    return domainfakultasunit.PagedFakultasUnits{
        Data:  fakultasunits,
        Total: total,
        CurrentPage: currentPage,
        TotalPages:  totalPages,
    }, nil
}