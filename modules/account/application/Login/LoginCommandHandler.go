package application

import (
	"context"
	
	domainaccount "UnpakSiamida/modules/account/domain"
	helper "UnpakSiamida/common/helper"
	"time"
)

type LoginCommandHandler struct{
	Repo domainaccount.IAccountRepository
}

func (h *LoginCommandHandler) Handle(
	ctx context.Context,
	cmd LoginCommand,
) (*domainaccount.LoginResult, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	user, err := h.Repo.Auth(ctx, cmd.Username, cmd.Password)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, domainaccount.InvalidCredential()
	}

	sid := user.UUID
	accessToken, refreshToken, err := helper.GenerateToken(sid)
	if err != nil {
		return nil, err
	}

	return &domainaccount.LoginResult{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		UserID:       sid,
	}, nil
}
