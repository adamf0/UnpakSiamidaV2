package application

import (
	"context"
	"errors"

	domainuser "UnpakSiamida/modules/user/domain"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
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
	_, err = h.Repo.GetByUuid(ctx, userUUID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", domainuser.NotFound(cmd.Uuid)
		}
		return "", err
	}

	// Delete by UUID
	if err := h.Repo.Delete(ctx, userUUID); err != nil { // FIXED
		return "", err
	}

	return cmd.Uuid, nil
}
