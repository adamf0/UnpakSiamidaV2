package application

import (
	"context"

	domaindokumentambahan "UnpakSiamida/modules/dokumentambahan/domain"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
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
	_, err = h.Repo.GetByUuid(ctx, dokumentambahanUUID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return "", domaindokumentambahan.NotFound(cmd.Uuid)
		}
		return "", err
	}

	// Delete by UUID
	if err := h.Repo.Delete(ctx, dokumentambahanUUID); err != nil { // FIXED
		return "", err
	}

	return cmd.Uuid, nil
}
