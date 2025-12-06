package application

import (
	"context"

	domainstandarrenstra "UnpakSiamida/modules/standarrenstra/domain"
	"github.com/google/uuid"
)

type DeleteStandarRenstraCommandHandler struct {
	Repo domainstandarrenstra.IStandarRenstraRepository
}

func (h *DeleteStandarRenstraCommandHandler) Handle(
	ctx context.Context,
	cmd DeleteStandarRenstraCommand,
) (string, error) {

	// Validate UUID
	standarrenstraUUID, err := uuid.Parse(cmd.Uuid)
	if err != nil {
		return "", domainstandarrenstra.InvalidUuid()
	}

	// Get existing standarrenstra
	existingStandarRenstra, err := h.Repo.GetByUuid(ctx, standarrenstraUUID)
	if err != nil {
		return "", err
	}
	if existingStandarRenstra == nil {
		return "", domainstandarrenstra.NotFound(cmd.Uuid)
	}

	// Delete by UUID
	if err := h.Repo.Delete(ctx, standarrenstraUUID); err != nil { // FIXED
		return "", err
	}

	return cmd.Uuid, nil
}
