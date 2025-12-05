package application

import (
    "context"

    "UnpakSiamida/modules/user/domain"
    "github.com/google/uuid"
)

type GetUserByUuidQueryHandler struct {
    Repo domain.IUserRepository
}

func (h *GetUserByUuidQueryHandler) Handle(
    ctx context.Context,
    q GetUserByUuidQuery,
) (*domain.User, error) {

    parsed, err := uuid.Parse(q.Uuid)
    if err != nil {
        return nil, err
    }

    return h.Repo.GetByUuid(ctx, parsed)
}
