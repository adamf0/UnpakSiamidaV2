package application

import (
	commondomain "UnpakSiamida/common/domain"
	domainrenstra "UnpakSiamida/modules/renstra/domain"
	"context"
	"time"
)

type GetAllRenstrasQueryHandler struct {
	Repo domainrenstra.IRenstraRepository
}

func (h *GetAllRenstrasQueryHandler) Handle(
	ctx context.Context,
	q GetAllRenstrasQuery,
) (commondomain.Paged[domainrenstra.RenstraDefault], error) {
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	renstras, total, err := h.Repo.GetAll(
		ctx,
		q.Search,
		q.SearchFilters,
		q.Page,
		q.Limit,
		q.Scope,
	)
	if err != nil {
		return commondomain.Paged[domainrenstra.RenstraDefault]{}, err
	}

	currentPage := 1
	totalPages := 1

	if q.Page != nil {
		currentPage = *q.Page
	}
	if q.Limit != nil && *q.Limit > 0 {
		totalPages = int((total + int64(*q.Limit) - 1) / int64(*q.Limit))
	}

	return commondomain.Paged[domainrenstra.RenstraDefault]{
		Data:        renstras,
		Total:       total,
		CurrentPage: currentPage,
		TotalPages:  totalPages,
	}, nil
}
