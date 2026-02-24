package application

import (
	commondomain "UnpakSiamida/common/domain"
	domainrenstranilai "UnpakSiamida/modules/renstranilai/domain"
	"context"
	"time"
)

type GetAllRenstraNilaisQueryHandler struct {
	Repo domainrenstranilai.IRenstraNilaiRepository
}

func (h *GetAllRenstraNilaisQueryHandler) Handle(
	ctx context.Context,
	q GetAllRenstraNilaisQuery,
) (commondomain.Paged[domainrenstranilai.RenstraNilaiDefault], error) {
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	renstranilais, total, err := h.Repo.GetAll(
		ctx,
		q.Search,
		q.SearchFilters,
		q.Page,
		q.Limit,
	)
	if err != nil {
		return commondomain.Paged[domainrenstranilai.RenstraNilaiDefault]{}, err
	}

	currentPage := 1
	totalPages := 1

	if q.Page != nil {
		currentPage = *q.Page
	}
	if q.Limit != nil && *q.Limit > 0 {
		totalPages = int((total + int64(*q.Limit) - 1) / int64(*q.Limit))
	}

	return commondomain.Paged[domainrenstranilai.RenstraNilaiDefault]{
		Data:        renstranilais,
		Total:       total,
		CurrentPage: currentPage,
		TotalPages:  totalPages,
	}, nil
}
