package application

import (
	"context"

	domainrenstra "UnpakSiamida/modules/renstra/domain"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type DeleteRenstraCommandHandler struct {
	Repo domainrenstra.IRenstraRepository
}

func (h *DeleteRenstraCommandHandler) Handle(
	ctx context.Context,
	cmd DeleteRenstraCommand,
) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// Validate UUID
	renstraUUID, err := uuid.Parse(cmd.Uuid)
	if err != nil {
		return "", domainrenstra.InvalidUuid()
	}

	// Get existing renstra
	_, err = h.Repo.GetByUuid(ctx, renstraUUID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return "", domainrenstra.NotFound(cmd.Uuid)
		}
		return "", err
	}

	// Delete by UUID
	if err := h.Repo.Delete(ctx, renstraUUID); err != nil { // FIXED
		return "", err
	}

	return cmd.Uuid, nil
}
