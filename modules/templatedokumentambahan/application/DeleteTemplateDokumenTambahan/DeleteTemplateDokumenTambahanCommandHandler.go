package application

import (
	"context"
	"errors"

	domaintemplatedokumentambahan "UnpakSiamida/modules/templatedokumentambahan/domain"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type DeleteTemplateDokumenTambahanCommandHandler struct {
	Repo domaintemplatedokumentambahan.ITemplateDokumenTambahanRepository
}

func (h *DeleteTemplateDokumenTambahanCommandHandler) Handle(
	ctx context.Context,
	cmd DeleteTemplateDokumenTambahanCommand,
) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// Validate UUID
	templatedokumentambahanUUID, err := uuid.Parse(cmd.Uuid)
	if err != nil {
		return "", domaintemplatedokumentambahan.InvalidUuid()
	}

	// Get existing templatedokumentambahan
	_, err = h.Repo.GetByUuid(ctx, templatedokumentambahanUUID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", domaintemplatedokumentambahan.NotFound(cmd.Uuid)
		}
		return "", err
	}

	// Delete by UUID
	if err := h.Repo.Delete(ctx, templatedokumentambahanUUID); err != nil { // FIXED
		return "", err
	}

	return cmd.Uuid, nil
}
