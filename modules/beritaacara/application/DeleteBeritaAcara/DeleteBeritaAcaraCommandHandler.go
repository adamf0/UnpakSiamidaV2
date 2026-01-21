package application

import (
	"context"

	domainberitaacara "UnpakSiamida/modules/beritaacara/domain"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type DeleteBeritaAcaraCommandHandler struct {
	Repo domainberitaacara.IBeritaAcaraRepository
}

func (h *DeleteBeritaAcaraCommandHandler) Handle(
	ctx context.Context,
	cmd DeleteBeritaAcaraCommand,
) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// Validate UUID
	BeritaAcaraUUID, err := uuid.Parse(cmd.Uuid)
	if err != nil {
		return "", domainberitaacara.InvalidUuid()
	}

	// Get existing BeritaAcara
	_, err = h.Repo.GetByUuid(ctx, BeritaAcaraUUID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return "", domainberitaacara.NotFound(cmd.Uuid)
		}
		return "", err
	}

	// Delete by UUID
	if err := h.Repo.Delete(ctx, BeritaAcaraUUID); err != nil { // FIXED
		return "", err
	}

	return cmd.Uuid, nil
}
