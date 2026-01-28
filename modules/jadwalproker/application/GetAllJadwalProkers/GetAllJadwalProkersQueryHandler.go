package application

import (
	domainJadwalProker "UnpakSiamida/modules/jadwalproker/domain"
	"context"
	"time"
)

type GetAllJadwalProkersQueryHandler struct {
	Repo domainJadwalProker.IJadwalProkerRepository
}

func (h *GetAllJadwalProkersQueryHandler) Handle(
	ctx context.Context,
	q GetAllJadwalProkersQuery,
) (domainJadwalProker.PagedJadwalProkers, error) {
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	JadwalProkers, total, err := h.Repo.GetAll(
		ctx,
		q.Search,
		q.SearchFilters,
		q.Page,
		q.Limit,
	)
	if err != nil {
		return domainJadwalProker.PagedJadwalProkers{}, err
	}

	currentPage := 1
	totalPages := 1

	if q.Page != nil {
		currentPage = *q.Page
	}
	if q.Limit != nil && *q.Limit > 0 {
		totalPages = int((total + int64(*q.Limit) - 1) / int64(*q.Limit))
	}

	return domainJadwalProker.PagedJadwalProkers{
		Data:        JadwalProkers,
		Total:       total,
		CurrentPage: currentPage,
		TotalPages:  totalPages,
	}, nil
}
