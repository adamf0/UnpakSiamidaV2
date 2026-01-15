package application

import (
	"context"
	"errors"

	domainaccount "UnpakSiamida/modules/account/domain"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type WhoamiCommandHandler struct {
	Repo domainaccount.IAccountRepository
}

func (h *WhoamiCommandHandler) Handle(
	ctx context.Context,
	cmd WhoamiCommand,
) (*domainaccount.Account, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	_, err := uuid.Parse(cmd.SID)
	if err != nil {
		return nil, domainaccount.NotFound(cmd.SID)
	}

	user, err := h.Repo.Get(ctx, cmd.SID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domainaccount.InvalidCredential()
		}
		return nil, err
	}
	if user.ExtraRole == nil {
		user.ExtraRole = []domainaccount.ExtraRole{}
	}

	return user, nil
}
