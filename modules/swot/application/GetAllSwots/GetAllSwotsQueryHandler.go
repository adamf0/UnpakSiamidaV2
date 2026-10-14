package application

import (
	commondomain "UnpakSiamida/common/domain"
	domainSwot "UnpakSiamida/modules/swot/domain"
	"context"
	"time"
)

type GetAllSwotsQueryHandler struct {
	Repo domainSwot.ISwotRepository
}

func (h *GetAllSwotsQueryHandler) Handle(
	ctx context.Context,
	q GetAllSwotsQuery,
) (commondomain.Paged[domainSwot.SwotDefault], error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	data, total, err := h.Repo.GetAll(ctx, q.Search, q.SearchFilters, q.Page, q.Limit)
	if err != nil {
		return commondomain.Paged[domainSwot.SwotDefault]{}, err
	}

	return commondomain.Paged[domainSwot.SwotDefault]{
		Data:  data,
		Total: total,
	}, nil
}
