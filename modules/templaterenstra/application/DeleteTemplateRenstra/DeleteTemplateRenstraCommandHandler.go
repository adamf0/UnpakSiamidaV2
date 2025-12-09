package application

import (
	"context"

	domaintemplaterenstra "UnpakSiamida/modules/templaterenstra/domain"
	"github.com/google/uuid"
)

type DeleteTemplateRenstraCommandHandler struct {
	Repo domaintemplaterenstra.ITemplateRenstraRepository
}

func (h *DeleteTemplateRenstraCommandHandler) Handle(
	ctx context.Context,
	cmd DeleteTemplateRenstraCommand,
) (string, error) {

	// Validate UUID
	templaterenstraUUID, err := uuid.Parse(cmd.Uuid)
	if err != nil {
		return "", domaintemplaterenstra.InvalidUuid()
	}

	// Get existing templaterenstra
	existingTemplateRenstra, err := h.Repo.GetByUuid(ctx, templaterenstraUUID)
	if err != nil {
		return "", err
	}
	if existingTemplateRenstra == nil {
		return "", domaintemplaterenstra.NotFound(cmd.Uuid)
	}

	// Delete by UUID
	if err := h.Repo.Delete(ctx, templaterenstraUUID); err != nil { // FIXED
		return "", err
	}

	return cmd.Uuid, nil
}
