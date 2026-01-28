package application

import (
	"context"

	domainAktivitasproker "UnpakSiamida/modules/aktivitasproker/domain"
)

type SetupUuidAktivitasProkerCommandHandler struct {
	Repo domainAktivitasproker.IAktivitasProkerRepository
}

func (h *SetupUuidAktivitasProkerCommandHandler) Handle(
	ctx context.Context,
	cmd SetupUuidAktivitasProkerCommand,
) (string, error) {

	err := h.Repo.SetupUuid(ctx)
	if err != nil {
		return "", err
	}

	return "berhasil setup uuid pada data", nil
}
