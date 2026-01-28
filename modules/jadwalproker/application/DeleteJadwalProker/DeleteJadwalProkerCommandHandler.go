package application

import (
	"context"

	domainjadwalproker "UnpakSiamida/modules/jadwalproker/domain"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type DeleteJadwalProkerCommandHandler struct {
	Repo domainjadwalproker.IJadwalProkerRepository
}

func (h *DeleteJadwalProkerCommandHandler) Handle(
	ctx context.Context,
	cmd DeleteJadwalProkerCommand,
) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// Validate UUID
	jadwalprokerUUID, err := uuid.Parse(cmd.Uuid)
	if err != nil {
		return "", domainjadwalproker.InvalidUuid()
	}

	// Get existing jadwalproker
	_, err = h.Repo.GetByUuid(ctx, jadwalprokerUUID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return "", domainjadwalproker.NotFound(cmd.Uuid)
		}
		return "", err
	}

	// Delete by UUID
	if err := h.Repo.Delete(ctx, jadwalprokerUUID); err != nil { // FIXED
		return "", err
	}

	return cmd.Uuid, nil
}
