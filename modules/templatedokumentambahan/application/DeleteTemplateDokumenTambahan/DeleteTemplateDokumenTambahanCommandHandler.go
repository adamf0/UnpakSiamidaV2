package application

import (
	"context"

	domaintemplatedokumentambahan "UnpakSiamida/modules/templatedokumentambahan/domain"
	"github.com/google/uuid"
)

type DeleteTemplateDokumenTambahanCommandHandler struct {
	Repo domaintemplatedokumentambahan.ITemplateDokumenTambahanRepository
}

func (h *DeleteTemplateDokumenTambahanCommandHandler) Handle(
	ctx context.Context,
	cmd DeleteTemplateDokumenTambahanCommand,
) (string, error) {

	// Validate UUID
	templatedokumentambahanUUID, err := uuid.Parse(cmd.Uuid)
	if err != nil {
		return "", domaintemplatedokumentambahan.InvalidUuid()
	}

	// Get existing templatedokumentambahan
	existingTemplateDokumenTambahan, err := h.Repo.GetByUuid(ctx, templatedokumentambahanUUID)
	if err != nil {
		return "", err
	}
	if existingTemplateDokumenTambahan == nil {
		return "", domaintemplatedokumentambahan.NotFound(cmd.Uuid)
	}

	// Delete by UUID
	if err := h.Repo.Delete(ctx, templatedokumentambahanUUID); err != nil { // FIXED
		return "", err
	}

	return cmd.Uuid, nil
}
