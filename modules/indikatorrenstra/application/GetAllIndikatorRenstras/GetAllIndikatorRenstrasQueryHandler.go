package application

import (
	commondomain "UnpakSiamida/common/domain"
	domainindikatorrenstra "UnpakSiamida/modules/indikatorrenstra/domain"
	"context"
	"time"
)

type GetAllIndikatorRenstrasQueryHandler struct {
	Repo domainindikatorrenstra.IIndikatorRenstraRepository
}

func (h *GetAllIndikatorRenstrasQueryHandler) Handle(
	ctx context.Context,
	q GetAllIndikatorRenstrasQuery,
) (commondomain.Paged[domainindikatorrenstra.IndikatorRenstraDefault], error) {
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	indikatorrenstras, total, err := h.Repo.GetAll(
		ctx,
		q.Search,
		q.SearchFilters,
		q.Page,
		q.Limit,
	)
	if err != nil {
		return commondomain.Paged[domainindikatorrenstra.IndikatorRenstraDefault]{}, err
	}

	currentPage := 1
	totalPages := 1

	if q.Page != nil {
		currentPage = *q.Page
	}
	if q.Limit != nil && *q.Limit > 0 {
		totalPages = int((total + int64(*q.Limit) - 1) / int64(*q.Limit))
	}

	return commondomain.Paged[domainindikatorrenstra.IndikatorRenstraDefault]{
		Data:        indikatorrenstras,
		Total:       total,
		CurrentPage: currentPage,
		TotalPages:  totalPages,
	}, nil
}
