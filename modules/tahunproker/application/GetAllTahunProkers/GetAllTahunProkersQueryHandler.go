package application

import (
	domaintahunproker "UnpakSiamida/modules/tahunproker/domain"
	"context"
	"time"
)

type GetAllTahunProkersQueryHandler struct {
	Repo domaintahunproker.ITahunProkerRepository
}

func (h *GetAllTahunProkersQueryHandler) Handle(
	ctx context.Context,
	q GetAllTahunProkersQuery,
) (domaintahunproker.PagedTahunProkers, error) {
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	TahunProkers, total, err := h.Repo.GetAll(
		ctx,
		q.Search,
		q.SearchFilters,
		q.Page,
		q.Limit,
	)
	if err != nil {
		return domaintahunproker.PagedTahunProkers{}, err
	}

	currentPage := 1
	totalPages := 1

	if q.Page != nil {
		currentPage = *q.Page
	}
	if q.Limit != nil && *q.Limit > 0 {
		totalPages = int((total + int64(*q.Limit) - 1) / int64(*q.Limit))
	}

	return domaintahunproker.PagedTahunProkers{
		Data:        TahunProkers,
		Total:       total,
		CurrentPage: currentPage,
		TotalPages:  totalPages,
	}, nil
}
