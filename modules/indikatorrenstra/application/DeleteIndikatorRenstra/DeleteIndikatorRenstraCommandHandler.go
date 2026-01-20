package application

import (
	"context"
	"errors"

	domainindikatorrenstra "UnpakSiamida/modules/indikatorrenstra/domain"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
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
	_, err = h.Repo.GetByUuid(ctx, indikatorrenstraUUID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", domainindikatorrenstra.NotFound(cmd.Uuid)
		}
		return "", err
	}

	// Delete by UUID
	if err := h.Repo.Delete(ctx, indikatorrenstraUUID); err != nil { // FIXED
		return "", err
	}

	return cmd.Uuid, nil
}
