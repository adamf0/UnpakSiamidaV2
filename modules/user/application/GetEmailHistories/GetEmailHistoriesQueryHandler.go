package application

import (
	domainuser "UnpakSiamida/modules/user/domain"
	"context"
	"time"
)

type GetEmailHistoriesQueryHandler struct {
	Repo domainuser.IUserRepository
}

func (h *GetEmailHistoriesQueryHandler) Handle(
	ctx context.Context,
	q GetEmailHistoriesQuery,
) (*GetEmailHistoriesResult, error) {
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	page := q.Page
	if page <= 0 {
		page = 1
	}
	limit := q.Limit
	if limit <= 0 || limit > 100 {
		limit = 20
	}

	list, total, err := h.Repo.GetEmailHistories(ctx, q.Search, q.Status, page, limit)
	if err != nil {
		return nil, err
	}

	return &GetEmailHistoriesResult{
		Data:  list,
		Total: total,
		Page:  page,
		Limit: limit,
	}, nil
}
