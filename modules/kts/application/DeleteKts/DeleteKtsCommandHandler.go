package application

import (
	"context"
	"errors"

	domainkts "UnpakSiamida/modules/kts/domain"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type DeleteKtsCommandHandler struct {
	Repo domainkts.IKtsRepository
}

func (h *DeleteKtsCommandHandler) Handle(
	ctx context.Context,
	cmd DeleteKtsCommand,
) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// Validate UUID
	ktsUUID, err := uuid.Parse(cmd.Uuid)
	if err != nil {
		return "", domainkts.InvalidUuid()
	}

	// Get existing kts
	_, err = h.Repo.GetByUuid(ctx, ktsUUID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", domainkts.NotFound(cmd.Uuid)
		}
		return "", err
	}

	// Delete by UUID
	if err := h.Repo.Delete(ctx, ktsUUID); err != nil { // FIXED
		return "", err
	}

	return cmd.Uuid, nil
}
