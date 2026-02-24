package application

import (
	commondomain "UnpakSiamida/common/domain"
	domainuser "UnpakSiamida/modules/user/domain"
	"context"
	"time"
)

type GetAllUsersQueryHandler struct {
	Repo domainuser.IUserRepository
}

func (h *GetAllUsersQueryHandler) Handle(
	ctx context.Context,
	q GetAllUsersQuery,
) (commondomain.Paged[domainuser.User], error) {
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	users, total, err := h.Repo.GetAll(
		ctx,
		q.Search,
		q.SearchFilters,
		q.Page,
		q.Limit,
	)
	if err != nil {
		return commondomain.Paged[domainuser.User]{}, err
	}

	currentPage := 1
	totalPages := 1

	if q.Page != nil {
		currentPage = *q.Page
	}
	if q.Limit != nil && *q.Limit > 0 {
		totalPages = int((total + int64(*q.Limit) - 1) / int64(*q.Limit))
	}

	return commondomain.Paged[domainuser.User]{
		Data:        users,
		Total:       total,
		CurrentPage: currentPage,
		TotalPages:  totalPages,
	}, nil
}
