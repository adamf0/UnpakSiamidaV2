package application

import (
	commondomain "UnpakSiamida/common/domain"
	domainKts "UnpakSiamida/modules/kts/domain"
	"context"

	"time"
)

type GetAllKtssQueryHandler struct {
	Repo domainKts.IKtsRepository
}

func (h *GetAllKtssQueryHandler) Handle(
	ctx context.Context,
	q GetAllKtssQuery,
) (commondomain.Paged[domainKts.KtsDefault], error) {
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	Ktss, total, err := h.Repo.GetAll(
		ctx,
		q.Search,
		q.SearchFilters,
		q.Page,
		q.Limit,
	)
	if err != nil {
		return commondomain.Paged[domainKts.KtsDefault]{}, err
	}

	currentPage := 1
	totalPages := 1

	if q.Page != nil {
		currentPage = *q.Page
	}
	if q.Limit != nil && *q.Limit > 0 {
		totalPages = int((total + int64(*q.Limit) - 1) / int64(*q.Limit))
	}

	return commondomain.Paged[domainKts.KtsDefault]{
		Data:        Ktss,
		Total:       total,
		CurrentPage: currentPage,
		TotalPages:  totalPages,
	}, nil
}
