package application

import (
	"context"
	"errors"

	domainrenstranilai "UnpakSiamida/modules/renstranilai/domain"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
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
	_, err = h.Repo.GetByUuid(ctx, renstranilaiUUID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", domainrenstranilai.NotFound(cmd.Uuid)
		}
		return "", err
	}

	// Delete by UUID
	if err := h.Repo.Delete(ctx, renstranilaiUUID); err != nil { // FIXED
		return "", err
	}

	return cmd.Uuid, nil
}
