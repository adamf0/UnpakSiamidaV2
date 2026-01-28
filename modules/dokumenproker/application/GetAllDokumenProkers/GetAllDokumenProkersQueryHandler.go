package application

import (
	domainDokumenProker "UnpakSiamida/modules/dokumenproker/domain"
	"context"
	"time"
)

type GetAllDokumenProkersQueryHandler struct {
	Repo domainDokumenProker.IDokumenProkerRepository
}

func (h *GetAllDokumenProkersQueryHandler) Handle(
	ctx context.Context,
	q GetAllDokumenProkersQuery,
) (domainDokumenProker.PagedDokumenProkers, error) {
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	DokumenProkers, total, err := h.Repo.GetAll(
		ctx,
		q.Search,
		q.SearchFilters,
		q.Page,
		q.Limit,
	)
	if err != nil {
		return domainDokumenProker.PagedDokumenProkers{}, err
	}

	currentPage := 1
	totalPages := 1

	if q.Page != nil {
		currentPage = *q.Page
	}
	if q.Limit != nil && *q.Limit > 0 {
		totalPages = int((total + int64(*q.Limit) - 1) / int64(*q.Limit))
	}

	return domainDokumenProker.PagedDokumenProkers{
		Data:        DokumenProkers,
		Total:       total,
		CurrentPage: currentPage,
		TotalPages:  totalPages,
	}, nil
}
