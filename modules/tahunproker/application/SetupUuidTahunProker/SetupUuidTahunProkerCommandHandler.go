package application

import (
	"context"

	domaintahunproker "UnpakSiamida/modules/tahunproker/domain"
)

type SetupUuidTahunProkerCommandHandler struct {
	Repo domaintahunproker.ITahunProkerRepository
}

func (h *SetupUuidTahunProkerCommandHandler) Handle(
	ctx context.Context,
	cmd SetupUuidTahunProkerCommand,
) (string, error) {

	err := h.Repo.SetupUuid(ctx)
	if err != nil {
		return "", err
	}

	return "berhasil setup uuid pada data", nil
}
