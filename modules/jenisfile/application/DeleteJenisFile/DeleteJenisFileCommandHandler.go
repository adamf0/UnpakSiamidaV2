package application

import (
	"context"

	domainjenisfile "UnpakSiamida/modules/jenisfile/domain"
	"github.com/google/uuid"
)

type DeleteJenisFileCommandHandler struct {
	Repo domainjenisfile.IJenisFileRepository
}

func (h *DeleteJenisFileCommandHandler) Handle(
	ctx context.Context,
	cmd DeleteJenisFileCommand,
) (string, error) {

	// Validate UUID
	jenisfileUUID, err := uuid.Parse(cmd.Uuid)
	if err != nil {
		return "", domainjenisfile.InvalidUuid()
	}

	// Get existing jenisfile
	existingJenisFile, err := h.Repo.GetByUuid(ctx, jenisfileUUID)
	if err != nil {
		return "", err
	}
	if existingJenisFile == nil {
		return "", domainjenisfile.NotFound(cmd.Uuid)
	}

	// Delete by UUID
	if err := h.Repo.Delete(ctx, jenisfileUUID); err != nil { // FIXED
		return "", err
	}

	return cmd.Uuid, nil
}
