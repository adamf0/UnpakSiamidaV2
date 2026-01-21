package application

import (
	"context"

	domainberitaacara "UnpakSiamida/modules/beritaacara/domain"
)

type SetupUuidBeritaAcaraCommandHandler struct {
	Repo domainberitaacara.IBeritaAcaraRepository
}

func (h *SetupUuidBeritaAcaraCommandHandler) Handle(
	ctx context.Context,
	cmd SetupUuidBeritaAcaraCommand,
) (string, error) {

	err := h.Repo.SetupUuid(ctx)
	if err != nil {
		return "", err
	}

	return "berhasil setup uuid pada data", nil
}
