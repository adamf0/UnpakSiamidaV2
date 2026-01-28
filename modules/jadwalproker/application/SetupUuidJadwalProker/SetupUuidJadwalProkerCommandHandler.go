package application

import (
	"context"

	domainjadwalproker "UnpakSiamida/modules/jadwalproker/domain"
)

type SetupUuidJadwalProkerCommandHandler struct {
	Repo domainjadwalproker.IJadwalProkerRepository
}

func (h *SetupUuidJadwalProkerCommandHandler) Handle(
	ctx context.Context,
	cmd SetupUuidJadwalProkerCommand,
) (string, error) {

	err := h.Repo.SetupUuid(ctx)
	if err != nil {
		return "", err
	}

	return "berhasil setup uuid pada data", nil
}
