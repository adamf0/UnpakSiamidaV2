package application

import (
	commondomain "UnpakSiamida/common/domain"
	domainfakultasunit "UnpakSiamida/modules/fakultasunit/domain"
	"context"
	"time"
)

type GetAllFakultasUnitsQueryHandler struct {
	Repo domainfakultasunit.IFakultasUnitRepository
}

func (h *GetAllFakultasUnitsQueryHandler) Handle(
	ctx context.Context,
	q GetAllFakultasUnitsQuery,
) (commondomain.Paged[domainfakultasunit.FakultasUnit], error) {
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	fakultasunits, total, err := h.Repo.GetAll(
		ctx,
		q.Search,
		q.SearchFilters,
		q.Page,
		q.Limit,
	)
	if err != nil {
		return commondomain.Paged[domainfakultasunit.FakultasUnit]{}, err
	}

	currentPage := 1
	totalPages := 1

	if q.Page != nil {
		currentPage = *q.Page
	}
	if q.Limit != nil && *q.Limit > 0 {
		totalPages = int((total + int64(*q.Limit) - 1) / int64(*q.Limit))
	}

	return commondomain.Paged[domainfakultasunit.FakultasUnit]{
		Data:        fakultasunits,
		Total:       total,
		CurrentPage: currentPage,
		TotalPages:  totalPages,
	}, nil
}
