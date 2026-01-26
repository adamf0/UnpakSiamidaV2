package application

import (
	"context"

	domainmataprogram "UnpakSiamida/modules/mataprogram/domain"
	domaintahunproker "UnpakSiamida/modules/tahunproker/domain"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UpdateMataProgramCommandHandler struct {
	Repo            domainmataprogram.IMataProgramRepository
	RepoTahunProker domaintahunproker.ITahunProkerRepository
}

func (h *UpdateMataProgramCommandHandler) Handle(
	ctx context.Context,
	cmd UpdateMataProgramCommand,
) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// -------------------------
	// VALIDATE UUID
	// -------------------------
	mataprogramUUID, err := uuid.Parse(cmd.Uuid)
	if err != nil {
		return "", domainmataprogram.InvalidUuid()
	}
	tahunprokerUUID, err := uuid.Parse(cmd.TahunUuid)
	if err != nil {
		return "", domainmataprogram.InvalidParseTahun()
	}

	// -------------------------
	// GET EXISTING mataprogram
	// -------------------------
	existingMataProgram, err := h.Repo.GetByUuid(ctx, mataprogramUUID) // ← memastikan pakai nama interface yg benar
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return "", domainmataprogram.NotFound(cmd.Uuid)
		}
		return "", err
	}

	var tahunprokerId uint
	tahunproker, err := h.RepoTahunProker.GetByUuid(ctx, tahunprokerUUID)
	if err != nil {
		tahunprokerId = tahunproker.ID
	} else {
		tahunprokerId = 0
	}

	// -------------------------
	// AGGREGATE ROOT LOGIC
	// -------------------------
	result := domainmataprogram.UpdateMataProgram(
		existingMataProgram,
		mataprogramUUID,
		tahunprokerId,
		cmd.MataProgram,
	)

	if !result.IsSuccess {
		return "", result.Error
	}

	updatedMataProgram := result.Value

	// -------------------------
	// SAVE TO REPOSITORY
	// -------------------------
	if err := h.Repo.Update(ctx, updatedMataProgram); err != nil {
		return "", err
	}

	return updatedMataProgram.UUID.String(), nil
}
