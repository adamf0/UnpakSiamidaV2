package application

import (
	"context"

	domainrenstranilai "UnpakSiamida/modules/renstranilai/domain"
	"github.com/google/uuid"
	"time"
)

type DeleteRenstraNilaiCommandHandler struct {
	Repo domainrenstranilai.IRenstraNilaiRepository
}

func (h *DeleteRenstraNilaiCommandHandler) Handle(
	ctx context.Context,
	cmd DeleteRenstraNilaiCommand,
) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// Validate UUID
	renstranilaiUUID, err := uuid.Parse(cmd.Uuid)
	if err != nil {
		return "", domainrenstranilai.InvalidUuid()
	}

	// Get existing renstranilai
	existingRenstraNilai, err := h.Repo.GetByUuid(ctx, renstranilaiUUID)
	if err != nil {
		return "", err
	}
	if existingRenstraNilai == nil {
		return "", domainrenstranilai.NotFound(cmd.Uuid)
	}

	// Delete by UUID
	if err := h.Repo.Delete(ctx, renstranilaiUUID); err != nil { // FIXED
		return "", err
	}

	return cmd.Uuid, nil
}
