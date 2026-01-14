package application

import (
	"context"
	// "strconv"
	// "errors"
    // "gorm.io/gorm"
	"github.com/google/uuid"
	domainindikatorrenstra "UnpakSiamida/modules/indikatorrenstra/domain"
	domainstandarrenstra "UnpakSiamida/modules/standarrenstra/domain"
	// helper "UnpakSiamida/common/helper"
	"time"
)

type CreateIndikatorRenstraCommandHandler struct{
	Repo domainindikatorrenstra.IIndikatorRenstraRepository
	RepoStandarRenstra domainstandarrenstra.IStandarRenstraRepository
}

func (h *CreateIndikatorRenstraCommandHandler) Handle(
	ctx context.Context,
	cmd CreateIndikatorRenstraCommand,
) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

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
		} else{
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

	operator := cmd.Operator

	isUnique, err := h.Repo.IsUniqueIndikator(ctx, cmd.Indikator, cmd.Tahun)
	if err != nil {
		return "", err
	}

	result := domainindikatorrenstra.NewIndikatorRenstra(
		cmd.Indikator,
		standar,
		parent,
		cmd.Tahun,
		cmd.TipeTarget,
		operator,
		isUnique,
	)

	if !result.IsSuccess {
		return "", result.Error
	}

	indikatorRenstra := result.Value

	if err := h.Repo.Create(ctx, indikatorRenstra); err != nil {
		return "", err
	}

	return indikatorRenstra.UUID.String(), nil
}