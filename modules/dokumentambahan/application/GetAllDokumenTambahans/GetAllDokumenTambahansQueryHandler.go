package application

import (
    "context"
    domaindokumentambahan "UnpakSiamida/modules/dokumentambahan/domain"
    "time"
)

type GetAllDokumenTambahansQueryHandler struct {
    Repo domaindokumentambahan.IDokumenTambahanRepository
}

func (h *GetAllDokumenTambahansQueryHandler) Handle(
    ctx context.Context,
    q GetAllDokumenTambahansQuery,
) (domaindokumentambahan.PagedDokumenTambahans, error) {
    ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

    dokumentambahans, total, err := h.Repo.GetAll(
        ctx,
        q.Search,
        q.SearchFilters,
        q.Page,
        q.Limit,
    )
    if err != nil {
        return domaindokumentambahan.PagedDokumenTambahans{}, err
    }

    currentPage := 1
    totalPages := 1

    if q.Page != nil {
        currentPage = *q.Page
    }
    if q.Limit != nil && *q.Limit > 0 {
        totalPages = int((total + int64(*q.Limit) - 1) / int64(*q.Limit))
    }

    return domaindokumentambahan.PagedDokumenTambahans{
        Data:  dokumentambahans,
        Total: total,
        CurrentPage: currentPage,
        TotalPages:  totalPages,
    }, nil
}