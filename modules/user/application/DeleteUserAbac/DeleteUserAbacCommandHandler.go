package DeleteUserAbac

import (
	"context"
	"strings"

	domainuser "UnpakSiamida/modules/user/domain"

	"github.com/google/uuid"
)

type DeleteUserAbacCommandHandler struct {
	Repo domainuser.IUserRepository
}

func (h *DeleteUserAbacCommandHandler) Handle(ctx context.Context, cmd DeleteUserAbacCommand) (string, error) {
	rawUUID := strings.TrimSpace(cmd.UUID)
	if rawUUID == "" {
		return "", domainuser.AbacUuidRequired()
	}

	parsedUUID, err := uuid.Parse(rawUUID)
	if err != nil {
		return "", domainuser.AbacInvalidUuid()
	}

	if err := h.Repo.DeleteUserAbac(ctx, parsedUUID); err != nil {
		return "", err
	}

	return "Penugasan role berhasil dihapus", nil
}
