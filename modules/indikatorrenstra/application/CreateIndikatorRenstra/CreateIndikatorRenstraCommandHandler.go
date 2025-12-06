package application

import (
	"context"
	"strconv"
	
	domainindikatorrenstra "UnpakSiamida/modules/indikatorrenstra/domain"
)

type CreateIndikatorRenstraCommandHandler struct{
	Repo domainindikatorrenstra.IIndikatorRenstraRepository
}

func (h *CreateIndikatorRenstraCommandHandler) Handle(
	ctx context.Context,
	cmd CreateIndikatorRenstraCommand,
) (string, error) {

	var standar *uint
	if cmd.StandarRenstra != "" {
		val, err := strconv.ParseUint(cmd.StandarRenstra, 10, 64)
		if err != nil {
			return "", err
		}
		v := uint(val)
		standar = &v
	}

	var parent *uint //[PR] untuk sekarang masih int seharusnya uuid yg diubaj jadi int. dimasukkan ke domain
	if cmd.Parent != nil && *cmd.Parent != "" {
		val, err := strconv.ParseUint(*cmd.Parent, 10, 64)
		if err != nil {
			return "", err
		}
		v := uint(val)
		parent = &v
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
