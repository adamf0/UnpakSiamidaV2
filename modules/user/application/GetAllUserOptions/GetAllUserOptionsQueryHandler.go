package application

import (
	domainuser "UnpakSiamida/modules/user/domain"
	"context"
	"time"
)

type GetAllUserOptionsQueryHandler struct {
	Repo domainuser.IUserRepository
}

func (h *GetAllUserOptionsQueryHandler) Handle(
	ctx context.Context,
	q GetAllUserOptionsQuery,
) ([]domainuser.UserOptions, error) {
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	users, err := h.Repo.GetAllStrict(ctx)
	if err != nil {
		return []domainuser.UserOptions{}, err
	}

	return users, nil
}
