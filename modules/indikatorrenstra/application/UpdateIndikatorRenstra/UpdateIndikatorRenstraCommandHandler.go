package application

import (
	"context"
	"errors"

	// "strconv"
	domainindikatorrenstra "UnpakSiamida/modules/indikatorrenstra/domain"
	domainstandarrenstra "UnpakSiamida/modules/standarrenstra/domain"

	// helper "UnpakSiamida/common/helper"
	"github.com/google/uuid"
	"gorm.io/gorm"

	// "errors"
	// "gorm.io/gorm"
	"time"
)

type UpdateIndikatorRenstraCommandHandler struct {
	Repo               domainindikatorrenstra.IIndikatorRenstraRepository
	RepoStandarRenstra domainstandarrenstra.IStandarRenstraRepository
}

func (h *UpdateIndikatorRenstraCommandHandler) Handle(
	ctx context.Context,
	cmd UpdateIndikatorRenstraCommand,
) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	indikatorrenstraUUID, err := uuid.Parse(cmd.Uuid)
	if err != nil {
		return "", domainindikatorrenstra.InvalidUuid()
	}

	existingIndikatorRenstra, err := h.Repo.GetByUuid(ctx, indikatorrenstraUUID) // ← memastikan pakai nama interface yg benar
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", domainindikatorrenstra.NotFound(cmd.Uuid)
		}
		return "", err
	}

	standarrenstraUUID, err := uuid.Parse(cmd.StandarRenstra)
	if err != nil {
		return "", domainindikatorrenstra.InvalidParent()
	}

	var standar *uint
	standarRenstra, err := h.RepoStandarRenstra.GetByUuid(ctx, standarrenstraUUID)
	if err != nil {
		standar = nil
	} else {
		standar = &standarRenstra.ID
	}

	var parentUUID uuid.UUID
	if cmd.Parent != nil && *cmd.Parent != "" {
		parsed, err := uuid.Parse(*cmd.Parent)
		if err != nil {
			parentUUID = uuid.Nil
		} else {
			parentUUID = parsed
		}
	} else {
		parentUUID = uuid.Nil
	}

	var parent *uint
	parentIndikator, err := h.Repo.GetDefaultByUuid(ctx, parentUUID)
	if err != nil {
		parent = nil
	} else {
		parent = &parentIndikator.Id
	}

	// -------------------------
	// AGGREGATE ROOT LOGIC
	// -------------------------
	result := domainindikatorrenstra.UpdateIndikatorRenstra(
		existingIndikatorRenstra,
		indikatorrenstraUUID,
		cmd.Indikator,
		standar,
		parent,
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
