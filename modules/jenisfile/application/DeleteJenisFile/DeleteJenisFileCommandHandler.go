package application

import (
	"context"

	domainjenisfile "UnpakSiamida/modules/jenisfile/domain"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type DeleteJenisFileCommandHandler struct {
	Repo domainjenisfile.IJenisFileRepository
}

func (h *DeleteJenisFileCommandHandler) Handle(
	ctx context.Context,
	cmd DeleteJenisFileCommand,
) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// Validate UUID
	jenisfileUUID, err := uuid.Parse(cmd.Uuid)
	if err != nil {
		return "", domainjenisfile.InvalidUuid()
	}

	// Get existing jenisfile
	_, err = h.Repo.GetByUuid(ctx, jenisfileUUID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return "", domainjenisfile.NotFound(cmd.Uuid)
		}
		return "", err
	}

	// Delete by UUID
	if err := h.Repo.Delete(ctx, jenisfileUUID); err != nil { // FIXED
		return "", err
	}

	return cmd.Uuid, nil
}
