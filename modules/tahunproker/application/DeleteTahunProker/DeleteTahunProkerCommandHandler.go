package application

import (
	"context"

	domaintahunproker "UnpakSiamida/modules/tahunproker/domain"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type DeleteTahunProkerCommandHandler struct {
	Repo domaintahunproker.ITahunProkerRepository
}

func (h *DeleteTahunProkerCommandHandler) Handle(
	ctx context.Context,
	cmd DeleteTahunProkerCommand,
) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// Validate UUID
	tahunprokerUUID, err := uuid.Parse(cmd.Uuid)
	if err != nil {
		return "", domaintahunproker.InvalidUuid()
	}

	// Get existing tahunproker
	_, err = h.Repo.GetByUuid(ctx, tahunprokerUUID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return "", domaintahunproker.NotFound(cmd.Uuid)
		}
		return "", err
	}

	// Delete by UUID
	if err := h.Repo.Delete(ctx, tahunprokerUUID); err != nil { // FIXED
		return "", err
	}

	return cmd.Uuid, nil
}
