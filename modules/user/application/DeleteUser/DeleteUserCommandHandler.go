package application

import (
	"context"

	domain "UnpakSiamida/modules/user/domain"
	"github.com/google/uuid"
)

type DeleteUserCommandHandler struct {
	Repo domain.IUserRepository
}

func (h *DeleteUserCommandHandler) Handle(
	ctx context.Context,
	cmd DeleteUserCommand,
) (string, error) {

	// Validate UUID
	userUUID, err := uuid.Parse(cmd.Uuid)
	if err != nil {
		return "", domain.InvalidUuid()
	}

	// Get existing user
	existingUser, err := h.Repo.GetByUuid(ctx, userUUID)
	if err != nil {
		return "", err
	}
	if existingUser == nil {
		return "", domain.NotFound(cmd.Uuid)
	}

	// Delete by UUID
	if err := h.Repo.Delete(ctx, userUUID); err != nil { // FIXED
		return "", err
	}

	return cmd.Uuid, nil
}
