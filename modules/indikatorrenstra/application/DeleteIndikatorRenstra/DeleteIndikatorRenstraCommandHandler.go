package application

import (
	"context"

	domainindikatorrenstra "UnpakSiamida/modules/indikatorrenstra/domain"
	"github.com/google/uuid"
	"time"
)

type DeleteIndikatorRenstraCommandHandler struct {
	Repo domainindikatorrenstra.IIndikatorRenstraRepository
}

func (h *DeleteIndikatorRenstraCommandHandler) Handle(
	ctx context.Context,
	cmd DeleteIndikatorRenstraCommand,
) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// Validate UUID
	indikatorrenstraUUID, err := uuid.Parse(cmd.Uuid)
	if err != nil {
		return "", domainindikatorrenstra.InvalidUuid()
	}

	// Get existing indikatorrenstra
	existingIndikatorRenstra, err := h.Repo.GetByUuid(ctx, indikatorrenstraUUID)
	if err != nil {
		return "", err
	}
	if existingIndikatorRenstra == nil {
		return "", domainindikatorrenstra.NotFound(cmd.Uuid)
	}

	// Delete by UUID
	if err := h.Repo.Delete(ctx, indikatorrenstraUUID); err != nil { // FIXED
		return "", err
	}

	return cmd.Uuid, nil
}
