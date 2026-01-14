package application

import (
    "context"
    domaintemplatedokumentambahan "UnpakSiamida/modules/templatedokumentambahan/domain"
    "time"
)

type GetAllTemplateDokumenTambahansQueryHandler struct {
    Repo domaintemplatedokumentambahan.ITemplateDokumenTambahanRepository
}

func (h *GetAllTemplateDokumenTambahansQueryHandler) Handle(
    ctx context.Context,
    q GetAllTemplateDokumenTambahansQuery,
) (domaintemplatedokumentambahan.PagedTemplateDokumenTambahans, error) {
    ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

    templatedokumentambahans, total, err := h.Repo.GetAll(
        ctx,
        q.Search,
        q.SearchFilters,
        q.Page,
        q.Limit,
    )
    if err != nil {
        return domaintemplatedokumentambahan.PagedTemplateDokumenTambahans{}, err
    }

    currentPage := 1
    totalPages := 1

    if q.Page != nil {
        currentPage = *q.Page
    }
    if q.Limit != nil && *q.Limit > 0 {
        totalPages = int((total + int64(*q.Limit) - 1) / int64(*q.Limit))
    }

    return domaintemplatedokumentambahan.PagedTemplateDokumenTambahans{
        Data:  templatedokumentambahans,
        Total: total,
        CurrentPage: currentPage,
        TotalPages:  totalPages,
    }, nil
}