package application

import (
	"context"
	"errors"

	domainrenstra "UnpakSiamida/modules/renstra/domain"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type GiveCodeAccessRenstraCommandHandler struct {
	Repo domainrenstra.IRenstraRepository
}

func (h *GiveCodeAccessRenstraCommandHandler) Handle(
	ctx context.Context,
	cmd GiveCodeAccessRenstraCommand,
) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	renstraUUID, err := uuid.Parse(cmd.Uuid)
	if err != nil {
		return "", domainrenstra.InvalidUuid()
	}

	existingRenstra, err := h.Repo.GetByUuid(ctx, renstraUUID) // ← memastikan pakai nama interface yg benar
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", domainrenstra.NotFound(cmd.Uuid)
		}
		return "", err
	}

	// -------------------------
	// AGGREGATE ROOT LOGIC
	// -------------------------
	result := domainrenstra.GiveCodeAccessRenstra(
		existingRenstra,
		renstraUUID,
		cmd.KodeAkses,
	)

	if !result.IsSuccess {
		return "", result.Error
	}

	updatedRenstra := result.Value

	// -------------------------
	// SAVE TO REPOSITORY
	// -------------------------
	if err := h.Repo.Update(ctx, updatedRenstra); err != nil {
		return "", err
	}

	return updatedRenstra.UUID.String(), nil
}
