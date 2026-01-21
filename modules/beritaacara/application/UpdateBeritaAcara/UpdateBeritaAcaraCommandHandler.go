package application

import (
	"context"

	domainberitaacara "UnpakSiamida/modules/beritaacara/domain"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UpdateBeritaAcaraCommandHandler struct {
	Repo domainberitaacara.IBeritaAcaraRepository
}

func (h *UpdateBeritaAcaraCommandHandler) Handle(
	ctx context.Context,
	cmd UpdateBeritaAcaraCommand,
) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// -------------------------
	// VALIDATE UUID
	// -------------------------
	BeritaAcaraUUID, err := uuid.Parse(cmd.Uuid)
	if err != nil {
		return "", domainberitaacara.InvalidUuid()
	}

	// -------------------------
	// GET EXISTING BeritaAcara
	// -------------------------
	existingBeritaAcara, err := h.Repo.GetByUuid(ctx, BeritaAcaraUUID) // ← memastikan pakai nama interface yg benar
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return "", domainberitaacara.NotFound(cmd.Uuid)
		}
		return "", err
	}

	// -------------------------
	// AGGREGATE ROOT LOGIC
	// -------------------------

	tanggal, err := time.Parse("2006-01-02", cmd.Tanggal)
	if err != nil {
		return "", domainberitaacara.InvalidTanggal()
	}

	result := domainberitaacara.UpdateBeritaAcara(
		existingBeritaAcara,
		BeritaAcaraUUID,
		cmd.Tahun,
		cmd.FakultasUnit,
		tanggal,
		cmd.Auditee,
		cmd.Auditor1,
		cmd.Auditor2,
	)

	if !result.IsSuccess {
		return "", result.Error
	}

	updatedBeritaAcara := result.Value

	// -------------------------
	// SAVE TO REPOSITORY
	// -------------------------
	if err := h.Repo.Update(ctx, updatedBeritaAcara); err != nil {
		return "", err
	}

	return updatedBeritaAcara.UUID.String(), nil
}
