package application

import (
	"context"
	"strconv"

	domainindikatorrenstra "UnpakSiamida/modules/indikatorrenstra/domain"
	"github.com/google/uuid"
)

type UpdateIndikatorRenstraCommandHandler struct {
	Repo domainindikatorrenstra.IIndikatorRenstraRepository
}

func (h *UpdateIndikatorRenstraCommandHandler) Handle(
	ctx context.Context,
	cmd UpdateIndikatorRenstraCommand,
) (string, error) {

	// -------------------------
	// VALIDATE UUID
	// -------------------------
	indikatorrenstraUUID, err := uuid.Parse(cmd.Uuid)
	if err != nil {
		return "", domainindikatorrenstra.InvalidUuid()
	}

	// -------------------------
	// GET EXISTING indikatorrenstra
	// -------------------------
	existingIndikatorRenstra, err := h.Repo.GetByUuid(ctx, indikatorrenstraUUID) // ← memastikan pakai nama interface yg benar
	if err != nil {
		return "", err
	}
	if existingIndikatorRenstra == nil {
		return "", domainindikatorrenstra.NotFound(cmd.Uuid)
	}

	// -------------------------
	// CHECK UNIQUE (jika indikator berubah)
	// -------------------------
	// isUnique, err := h.Repo.IsUniqueIndikator(ctx, cmd.Indikator, cmd.Tahun)
	// if err != nil {
	// 	return "", err
	// }

	var standarPtr *uint
	if cmd.StandarRenstra != "" {
		v, err := strconv.ParseUint(cmd.StandarRenstra, 10, 64)
		if err != nil {
			return "", domainindikatorrenstra.InvalidStandar()
		}
		tmp := uint(v)
		standarPtr = &tmp
	}

	var parentPtr *uint //[PR] untuk sekarang masih int seharusnya uuid yg diubaj jadi int. dimasukkan ke domain
	if cmd.Parent != nil && *cmd.Parent != "" {
		v, err := strconv.ParseUint(*cmd.Parent, 10, 64)
		if err != nil {
			return "", domainindikatorrenstra.InvalidParent()
		}
		tmp := uint(v)
		parentPtr = &tmp
	}

	// -------------------------
	// AGGREGATE ROOT LOGIC
	// -------------------------
	result := domainindikatorrenstra.UpdateIndikatorRenstra(
		existingIndikatorRenstra,
		indikatorrenstraUUID,
		cmd.Indikator,
		standarPtr,
		parentPtr,
		cmd.Tahun,
		cmd.TipeTarget,
		cmd.Operator,
		// isUnique,
	)

	if !result.IsSuccess {
		return "", result.Error
	}

	updatedIndikatorRenstra := result.Value

	// -------------------------
	// SAVE TO REPOSITORY
	// -------------------------
	if err := h.Repo.Update(ctx, updatedIndikatorRenstra); err != nil {
		return "", err
	}

	return updatedIndikatorRenstra.UUID.String(), nil
}
