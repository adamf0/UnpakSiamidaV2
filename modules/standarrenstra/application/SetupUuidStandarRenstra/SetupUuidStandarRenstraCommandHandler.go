package application

import (
	"context"

	domainstandarrenstra "UnpakSiamida/modules/standarrenstra/domain"
)

type SetupUuidStandarRenstraCommandHandler struct {
	Repo domainstandarrenstra.IStandarRenstraRepository
}

func (h *SetupUuidStandarRenstraCommandHandler) Handle(
	ctx context.Context,
	cmd SetupUuidStandarRenstraCommand,
) (string, error) {

	err := h.Repo.SetupUuid(ctx)
	if err != nil {
		return "", err
	}

	return "berhasil setup uuid pada data", nil
}
