package application

import (
	"context"
	
	domainstandarrenstra "UnpakSiamida/modules/standarrenstra/domain"
)

type CreateStandarRenstraCommandHandler struct{
	Repo domainstandarrenstra.IStandarRenstraRepository
}

func (h *CreateStandarRenstraCommandHandler) Handle(
	ctx context.Context,
	cmd CreateStandarRenstraCommand,
) (string, error) {

	result := domainstandarrenstra.NewStandarRenstra(
		cmd.Nama,
	)

	if !result.IsSuccess {
		return "", result.Error
	}

	createStandarRenstra := result.Value
	if err := h.Repo.Create(ctx, createStandarRenstra); err != nil {
		return "", err
	}

	return result.Value.UUID.String(), nil
}
