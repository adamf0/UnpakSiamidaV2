package application

import (
	"context"

	domainjenisfile "UnpakSiamida/modules/jenisfile/domain"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UpdateJenisFileCommandHandler struct {
	Repo domainjenisfile.IJenisFileRepository
}

func (h *UpdateJenisFileCommandHandler) Handle(
	ctx context.Context,
	cmd UpdateJenisFileCommand,
) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// -------------------------
	// VALIDATE UUID
	// -------------------------
	jenisfileUUID, err := uuid.Parse(cmd.Uuid)
	if err != nil {
		return "", domainjenisfile.InvalidUuid()
	}

	// -------------------------
	// GET EXISTING jenisfile
	// -------------------------
	existingJenisFile, err := h.Repo.GetByUuid(ctx, jenisfileUUID) // ← memastikan pakai nama interface yg benar
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return "", domainjenisfile.NotFound(cmd.Uuid)
		}
		return "", err
	}

	// -------------------------
	// AGGREGATE ROOT LOGIC
	// -------------------------
	result := domainjenisfile.UpdateJenisFile(
		existingJenisFile,
		jenisfileUUID,
		cmd.Nama,
	)

	if !result.IsSuccess {
		return "", result.Error
	}

	updatedJenisFile := result.Value

	// -------------------------
	// SAVE TO REPOSITORY
	// -------------------------
	if err := h.Repo.Update(ctx, updatedJenisFile); err != nil {
		return "", err
	}

	return updatedJenisFile.UUID.String(), nil
}
