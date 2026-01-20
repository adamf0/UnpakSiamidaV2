package application

import (
	"context"
	"errors"

	domaintemplaterenstra "UnpakSiamida/modules/templaterenstra/domain"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type DeleteTemplateRenstraCommandHandler struct {
	Repo domaintemplaterenstra.ITemplateRenstraRepository
}

func (h *DeleteTemplateRenstraCommandHandler) Handle(
	ctx context.Context,
	cmd DeleteTemplateRenstraCommand,
) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// Validate UUID
	templaterenstraUUID, err := uuid.Parse(cmd.Uuid)
	if err != nil {
		return "", domaintemplaterenstra.InvalidUuid()
	}

	// Get existing templaterenstra
	_, err = h.Repo.GetByUuid(ctx, templaterenstraUUID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", domaintemplaterenstra.NotFound(cmd.Uuid)
		}
		return "", err
	}

	// Delete by UUID
	if err := h.Repo.Delete(ctx, templaterenstraUUID); err != nil { // FIXED
		return "", err
	}

	return cmd.Uuid, nil
}
