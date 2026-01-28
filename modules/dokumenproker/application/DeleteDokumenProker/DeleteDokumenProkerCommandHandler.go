package application

import (
	"context"

	domaindokumenproker "UnpakSiamida/modules/dokumenproker/domain"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type DeleteDokumenProkerCommandHandler struct {
	Repo domaindokumenproker.IDokumenProkerRepository
}

func (h *DeleteDokumenProkerCommandHandler) Handle(
	ctx context.Context,
	cmd DeleteDokumenProkerCommand,
) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// Validate UUID
	uuidDokumenProker, err := uuid.Parse(cmd.Uuid)
	if err != nil {
		return "", domaindokumenproker.InvalidUuid()
	}

	_, err = h.Repo.GetByUuid(ctx, uuidDokumenProker)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return "", domaindokumenproker.NotFound(cmd.Uuid)
		}
		return "", err
	}

	// Delete by UUID
	if err := h.Repo.Delete(ctx, uuidDokumenProker); err != nil { // FIXED
		return "", err
	}

	return cmd.Uuid, nil
}
