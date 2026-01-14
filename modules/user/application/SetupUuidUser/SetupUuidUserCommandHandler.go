package application

import (
	"context"

	domainuser "UnpakSiamida/modules/user/domain"
)

type SetupUuidUserCommandHandler struct {
	Repo domainuser.IUserRepository
}

func (h *SetupUuidUserCommandHandler) Handle(
	ctx context.Context,
	cmd SetupUuidUserCommand,
) (string, error) {

	err := h.Repo.SetupUuid(ctx)
	if err != nil {
		return "", err
	}

	return "berhasil setup uuid pada data", nil
}
