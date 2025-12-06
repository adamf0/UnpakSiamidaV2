package application

import (
    "context"
    domainuser "UnpakSiamida/modules/user/domain"
)

type GetAllUsersQueryHandler struct {
    Repo domainuser.IUserRepository
}

func (h *GetAllUsersQueryHandler) Handle(
    ctx context.Context,
    q GetAllUsersQuery,
) (domainuser.PagedUsers, error) {

    users, total, err := h.Repo.GetAll(
        ctx,
        q.Search,
        q.SearchFilters,
        q.Page,
        q.Limit,
    )
    if err != nil {
        return domainuser.PagedUsers{}, err
    }

    currentPage := 1
    totalPages := 1

    if q.Page != nil {
        currentPage = *q.Page
    }
    if q.Limit != nil && *q.Limit > 0 {
        totalPages = int((total + int64(*q.Limit) - 1) / int64(*q.Limit))
    }

    return domainuser.PagedUsers{
        Data:  users,
        Total: total,
        CurrentPage: currentPage,
        TotalPages:  totalPages,
    }, nil
}