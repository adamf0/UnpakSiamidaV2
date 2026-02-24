package application

import (
	commondomain "UnpakSiamida/common/domain"
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
) (commondomain.Paged[domainJadwalProker.JadwalProkerDefault], error) {
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
		return commondomain.Paged[domainJadwalProker.JadwalProkerDefault]{}, err
	}

	currentPage := 1
	totalPages := 1

	if q.Page != nil {
		currentPage = *q.Page
	}
	if q.Limit != nil && *q.Limit > 0 {
		totalPages = int((total + int64(*q.Limit) - 1) / int64(*q.Limit))
	}

	return commondomain.Paged[domainJadwalProker.JadwalProkerDefault]{
		Data:        JadwalProkers,
		Total:       total,
		CurrentPage: currentPage,
		TotalPages:  totalPages,
	}, nil
}
