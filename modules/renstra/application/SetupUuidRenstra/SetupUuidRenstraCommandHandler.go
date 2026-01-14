package application

import (
	"context"

	domainrenstra "UnpakSiamida/modules/renstra/domain"
)

type SetupUuidRenstraCommandHandler struct {
	Repo domainrenstra.IRenstraRepository
}

func (h *SetupUuidRenstraCommandHandler) Handle(
	ctx context.Context,
	cmd SetupUuidRenstraCommand,
) (string, error) {

	err := h.Repo.SetupUuid(ctx)
	if err != nil {
		return "", err
	}

	return "berhasil setup uuid pada data", nil
}
