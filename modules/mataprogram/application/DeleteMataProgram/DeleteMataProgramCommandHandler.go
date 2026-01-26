package application

import (
	"context"

	domainmataprogram "UnpakSiamida/modules/mataprogram/domain"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type DeleteMataProgramCommandHandler struct {
	Repo domainmataprogram.IMataProgramRepository
}

func (h *DeleteMataProgramCommandHandler) Handle(
	ctx context.Context,
	cmd DeleteMataProgramCommand,
) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// Validate UUID
	mataprogramUUID, err := uuid.Parse(cmd.Uuid)
	if err != nil {
		return "", domainmataprogram.InvalidUuid()
	}

	// Get existing mataprogram
	_, err = h.Repo.GetByUuid(ctx, mataprogramUUID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return "", domainmataprogram.NotFound(cmd.Uuid)
		}
		return "", err
	}

	// Delete by UUID
	if err := h.Repo.Delete(ctx, mataprogramUUID); err != nil { // FIXED
		return "", err
	}

	return cmd.Uuid, nil
}
