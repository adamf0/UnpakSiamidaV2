package application

import (
	"context"
	"errors"

	domainstandarrenstra "UnpakSiamida/modules/standarrenstra/domain"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UpdateStandarRenstraCommandHandler struct {
	Repo domainstandarrenstra.IStandarRenstraRepository
}

func (h *UpdateStandarRenstraCommandHandler) Handle(
	ctx context.Context,
	cmd UpdateStandarRenstraCommand,
) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// -------------------------
	// VALIDATE UUID
	// -------------------------
	standarrenstraUUID, err := uuid.Parse(cmd.Uuid)
	if err != nil {
		return "", domainstandarrenstra.InvalidUuid()
	}

	// -------------------------
	// GET EXISTING standarrenstra
	// -------------------------
	existingStandarRenstra, err := h.Repo.GetByUuid(ctx, standarrenstraUUID) // ← memastikan pakai nama interface yg benar
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", domainstandarrenstra.NotFound(cmd.Uuid)
		}
		return "", err
	}

	// -------------------------
	// AGGREGATE ROOT LOGIC
	// -------------------------
	result := domainstandarrenstra.UpdateStandarRenstra(
		existingStandarRenstra,
		standarrenstraUUID,
		cmd.Nama,
	)

	if !result.IsSuccess {
		return "", result.Error
	}

	updatedStandarRenstra := result.Value

	// -------------------------
	// SAVE TO REPOSITORY
	// -------------------------
	if err := h.Repo.Update(ctx, updatedStandarRenstra); err != nil {
		return "", err
	}

	return updatedStandarRenstra.UUID.String(), nil
}
