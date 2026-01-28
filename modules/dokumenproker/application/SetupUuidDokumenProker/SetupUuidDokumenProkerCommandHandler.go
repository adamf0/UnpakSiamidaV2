package application

import (
	"context"

	domaindokumenproker "UnpakSiamida/modules/dokumenproker/domain"
)

type SetupUuidDokumenProkerCommandHandler struct {
	Repo domaindokumenproker.IDokumenProkerRepository
}

func (h *SetupUuidDokumenProkerCommandHandler) Handle(
	ctx context.Context,
	cmd SetupUuidDokumenProkerCommand,
) (string, error) {

	err := h.Repo.SetupUuid(ctx)
	if err != nil {
		return "", err
	}

	return "berhasil setup uuid pada data", nil
}
