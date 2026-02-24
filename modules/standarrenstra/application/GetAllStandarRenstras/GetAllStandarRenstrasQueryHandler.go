package application

import (
	commondomain "UnpakSiamida/common/domain"
	domainstandarrenstra "UnpakSiamida/modules/standarrenstra/domain"
	"context"
	"time"
)

type GetAllStandarRenstrasQueryHandler struct {
	Repo domainstandarrenstra.IStandarRenstraRepository
}

func (h *GetAllStandarRenstrasQueryHandler) Handle(
	ctx context.Context,
	q GetAllStandarRenstrasQuery,
) (commondomain.Paged[domainstandarrenstra.StandarRenstra], error) {
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	standarrenstras, total, err := h.Repo.GetAll(
		ctx,
		q.Search,
		q.SearchFilters,
		q.Page,
		q.Limit,
	)
	if err != nil {
		return commondomain.Paged[domainstandarrenstra.StandarRenstra]{}, err
	}

	currentPage := 1
	totalPages := 1

	if q.Page != nil {
		currentPage = *q.Page
	}
	if q.Limit != nil && *q.Limit > 0 {
		totalPages = int((total + int64(*q.Limit) - 1) / int64(*q.Limit))
	}

	return commondomain.Paged[domainstandarrenstra.StandarRenstra]{
		Data:        standarrenstras,
		Total:       total,
		CurrentPage: currentPage,
		TotalPages:  totalPages,
	}, nil
}
