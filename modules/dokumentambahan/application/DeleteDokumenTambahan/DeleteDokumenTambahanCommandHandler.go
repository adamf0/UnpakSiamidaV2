package application

import (
	"context"

	domaindokumentambahan "UnpakSiamida/modules/dokumentambahan/domain"
	"github.com/google/uuid"
	"time"
)

type DeleteDokumenTambahanCommandHandler struct {
	Repo domaindokumentambahan.IDokumenTambahanRepository
}

func (h *DeleteDokumenTambahanCommandHandler) Handle(
	ctx context.Context,
	cmd DeleteDokumenTambahanCommand,
) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// Validate UUID
	dokumentambahanUUID, err := uuid.Parse(cmd.Uuid)
	if err != nil {
		return "", domaindokumentambahan.InvalidUuid()
	}

	// Get existing dokumentambahan
	existingDokumenTambahan, err := h.Repo.GetByUuid(ctx, dokumentambahanUUID)
	if err != nil {
		return "", err
	}
	if existingDokumenTambahan == nil {
		return "", domaindokumentambahan.NotFound(cmd.Uuid)
	}

	// Delete by UUID
	if err := h.Repo.Delete(ctx, dokumentambahanUUID); err != nil { // FIXED
		return "", err
	}

	return cmd.Uuid, nil
}
