package application

import (
	"context"

	domainaktivitasproker "UnpakSiamida/modules/aktivitasproker/domain"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type DeleteAktivitasProkerCommandHandler struct {
	Repo domainaktivitasproker.IAktivitasProkerRepository
}

func (h *DeleteAktivitasProkerCommandHandler) Handle(
	ctx context.Context,
	cmd DeleteAktivitasProkerCommand,
) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// Validate UUID
	uuidAktivitasProker, err := uuid.Parse(cmd.Uuid)
	if err != nil {
		return "", domainaktivitasproker.InvalidUuid()
	}

	// Get existing Aktivitasproker
	_, err = h.Repo.GetByUuid(ctx, uuidAktivitasProker)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return "", domainaktivitasproker.NotFound(cmd.Uuid)
		}
		return "", err
	}

	// Delete by UUID
	if err := h.Repo.Delete(ctx, uuidAktivitasProker); err != nil { // FIXED
		return "", err
	}

	return cmd.Uuid, nil
}
