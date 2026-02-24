package application

import (
	commondomain "UnpakSiamida/common/domain"
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
) (commondomain.Paged[domaintahunproker.TahunProker], error) {
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
		return commondomain.Paged[domaintahunproker.TahunProker]{}, err
	}

	currentPage := 1
	totalPages := 1

	if q.Page != nil {
		currentPage = *q.Page
	}
	if q.Limit != nil && *q.Limit > 0 {
		totalPages = int((total + int64(*q.Limit) - 1) / int64(*q.Limit))
	}

	return commondomain.Paged[domaintahunproker.TahunProker]{
		Data:        TahunProkers,
		Total:       total,
		CurrentPage: currentPage,
		TotalPages:  totalPages,
	}, nil
}
