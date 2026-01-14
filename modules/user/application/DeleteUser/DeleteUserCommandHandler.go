package application

import (
	"context"

	domainuser "UnpakSiamida/modules/user/domain"
	"github.com/google/uuid"
	"time"
)

type DeleteUserCommandHandler struct {
	Repo domainuser.IUserRepository
}

func (h *DeleteUserCommandHandler) Handle(
	ctx context.Context,
	cmd DeleteUserCommand,
) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// Validate UUID
	userUUID, err := uuid.Parse(cmd.Uuid)
	if err != nil {
		return "", domainuser.InvalidUuid()
	}

	// Get existing user
	existingUser, err := h.Repo.GetByUuid(ctx, userUUID)
	if err != nil {
		return "", err
	}
	if existingUser == nil {
		return "", domainuser.NotFound(cmd.Uuid)
	}

	// Delete by UUID
	if err := h.Repo.Delete(ctx, userUUID); err != nil { // FIXED
		return "", err
	}

	return cmd.Uuid, nil
}
