package application

import (
    "context"
    domaintemplaterenstra "UnpakSiamida/modules/templaterenstra/domain"
)

type GetAllTemplateRenstrasQueryHandler struct {
    Repo domaintemplaterenstra.ITemplateRenstraRepository
}

func (h *GetAllTemplateRenstrasQueryHandler) Handle(
    ctx context.Context,
    q GetAllTemplateRenstrasQuery,
) (domaintemplaterenstra.PagedTemplateRenstras, error) {

    templaterenstras, total, err := h.Repo.GetAll(
        ctx,
        q.Search,
        q.SearchFilters,
        q.Page,
        q.Limit,
    )
    if err != nil {
        return domaintemplaterenstra.PagedTemplateRenstras{}, err
    }

    currentPage := 1
    totalPages := 1

    if q.Page != nil {
        currentPage = *q.Page
    }
    if q.Limit != nil && *q.Limit > 0 {
        totalPages = int((total + int64(*q.Limit) - 1) / int64(*q.Limit))
    }

    return domaintemplaterenstra.PagedTemplateRenstras{
        Data:  templaterenstras,
        Total: total,
        CurrentPage: currentPage,
        TotalPages:  totalPages,
    }, nil
}