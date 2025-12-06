package application

import (
    "context"

    domainuser "UnpakSiamida/modules/user/domain"
    "github.com/google/uuid"
)

type GetUserByUuidQueryHandler struct {
    Repo domainuser.IUserRepository
}

func (h *GetUserByUuidQueryHandler) Handle(
    ctx context.Context,
    q GetUserByUuidQuery,
) (*domainuser.User, error) {

    parsed, err := uuid.Parse(q.Uuid)
    if err != nil {
        return nil, err
    }

    return h.Repo.GetByUuid(ctx, parsed)
}
