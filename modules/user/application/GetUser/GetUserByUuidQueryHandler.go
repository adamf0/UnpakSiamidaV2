package application

import (
    "context"

    domainuser "UnpakSiamida/modules/user/domain"
    "github.com/google/uuid"
    "errors"
    "gorm.io/gorm"
    "time"
)

type GetUserByUuidQueryHandler struct {
    Repo domainuser.IUserRepository
}

func (h *GetUserByUuidQueryHandler) Handle(
    ctx context.Context,
    q GetUserByUuidQuery,
) (*domainuser.User, error) {
    ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

    parsed, err := uuid.Parse(q.Uuid)
	if err != nil {
		return nil, domainuser.NotFound(q.Uuid)
	}

    user, err := h.Repo.GetByUuid(ctx, parsed)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domainuser.NotFound(q.Uuid)
		}
		return nil, err
	}

    return user, nil
}
