package application

import (
	"context"

	domaintahunproker "UnpakSiamida/modules/tahunproker/domain"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UpdateTahunProkerCommandHandler struct {
	Repo domaintahunproker.ITahunProkerRepository
}

func (h *UpdateTahunProkerCommandHandler) Handle(
	ctx context.Context,
	cmd UpdateTahunProkerCommand,
) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// -------------------------
	// VALIDATE UUID
	// -------------------------
	tahunprokerUUID, err := uuid.Parse(cmd.Uuid)
	if err != nil {
		return "", domaintahunproker.InvalidUuid()
	}

	// -------------------------
	// GET EXISTING tahunproker
	// -------------------------
	existingTahunProker, err := h.Repo.GetByUuid(ctx, tahunprokerUUID) // ← memastikan pakai nama interface yg benar
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return "", domaintahunproker.NotFound(cmd.Uuid)
		}
		return "", err
	}

	// -------------------------
	// AGGREGATE ROOT LOGIC
	// -------------------------
	result := domaintahunproker.UpdateTahunProker(
		existingTahunProker,
		tahunprokerUUID,
		cmd.Tahun,
		cmd.Status,
	)

	if !result.IsSuccess {
		return "", result.Error
	}

	updatedTahunProker := result.Value

	// -------------------------
	// SAVE TO REPOSITORY
	// -------------------------
	if err := h.Repo.Update(ctx, updatedTahunProker); err != nil {
		return "", err
	}

	return updatedTahunProker.UUID.String(), nil
}
