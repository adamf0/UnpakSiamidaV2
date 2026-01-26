package application

import (
	domainmataprogram "UnpakSiamida/modules/mataprogram/domain"
	"context"
	"time"
)

type GetAllMataProgramsQueryHandler struct {
	Repo domainmataprogram.IMataProgramRepository
}

func (h *GetAllMataProgramsQueryHandler) Handle(
	ctx context.Context,
	q GetAllMataProgramsQuery,
) (domainmataprogram.PagedMataPrograms, error) {
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	MataPrograms, total, err := h.Repo.GetAll(
		ctx,
		q.Search,
		q.SearchFilters,
		q.Page,
		q.Limit,
	)
	if err != nil {
		return domainmataprogram.PagedMataPrograms{}, err
	}

	currentPage := 1
	totalPages := 1

	if q.Page != nil {
		currentPage = *q.Page
	}
	if q.Limit != nil && *q.Limit > 0 {
		totalPages = int((total + int64(*q.Limit) - 1) / int64(*q.Limit))
	}

	return domainmataprogram.PagedMataPrograms{
		Data:        MataPrograms,
		Total:       total,
		CurrentPage: currentPage,
		TotalPages:  totalPages,
	}, nil
}
