package application

import (
	"context"

	domainindikatorrenstra "UnpakSiamida/modules/indikatorrenstra/domain"
)

type SetupUuidIndikatorRenstraCommandHandler struct {
	Repo domainindikatorrenstra.IIndikatorRenstraRepository
}

func (h *SetupUuidIndikatorRenstraCommandHandler) Handle(
	ctx context.Context,
	cmd SetupUuidIndikatorRenstraCommand,
) (string, error) {

	err := h.Repo.SetupUuid(ctx)
	if err != nil {
		return "", err
	}

	return "berhasil setup uuid pada data", nil
}
