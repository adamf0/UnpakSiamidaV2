package application

import (
	domainberitaacara "UnpakSiamida/modules/beritaacara/domain"
	"context"
	"time"
)

type GetAllBeritaAcarasQueryHandler struct {
	Repo domainberitaacara.IBeritaAcaraRepository
}

func (h *GetAllBeritaAcarasQueryHandler) Handle(
	ctx context.Context,
	q GetAllBeritaAcarasQuery,
) (domainberitaacara.PagedBeritaAcaras, error) {
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	BeritaAcaras, total, err := h.Repo.GetAll(
		ctx,
		q.Search,
		q.SearchFilters,
		q.Page,
		q.Limit,
	)
	if err != nil {
		return domainberitaacara.PagedBeritaAcaras{}, err
	}

	currentPage := 1
	totalPages := 1

	if q.Page != nil {
		currentPage = *q.Page
	}
	if q.Limit != nil && *q.Limit > 0 {
		totalPages = int((total + int64(*q.Limit) - 1) / int64(*q.Limit))
	}

	return domainberitaacara.PagedBeritaAcaras{
		Data:        BeritaAcaras,
		Total:       total,
		CurrentPage: currentPage,
		TotalPages:  totalPages,
	}, nil
}
