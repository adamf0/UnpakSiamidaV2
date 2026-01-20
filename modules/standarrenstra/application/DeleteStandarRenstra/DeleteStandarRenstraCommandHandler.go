package application

import (
	"context"
	"errors"

	domainstandarrenstra "UnpakSiamida/modules/standarrenstra/domain"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type DeleteStandarRenstraCommandHandler struct {
	Repo domainstandarrenstra.IStandarRenstraRepository
}

func (h *DeleteStandarRenstraCommandHandler) Handle(
	ctx context.Context,
	cmd DeleteStandarRenstraCommand,
) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// Validate UUID
	standarrenstraUUID, err := uuid.Parse(cmd.Uuid)
	if err != nil {
		return "", domainstandarrenstra.InvalidUuid()
	}

	// Get existing standarrenstra
	_, err = h.Repo.GetByUuid(ctx, standarrenstraUUID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", domainstandarrenstra.NotFound(cmd.Uuid)
		}
		return "", err
	}

	// Delete by UUID
	if err := h.Repo.Delete(ctx, standarrenstraUUID); err != nil { // FIXED
		return "", err
	}

	return cmd.Uuid, nil
}
