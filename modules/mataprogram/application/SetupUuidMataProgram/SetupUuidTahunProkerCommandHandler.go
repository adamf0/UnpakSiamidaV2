package application

import (
	"context"

	domainmataprogram "UnpakSiamida/modules/mataprogram/domain"
)

type SetupUuidMataProgramCommandHandler struct {
	Repo domainmataprogram.IMataProgramRepository
}

func (h *SetupUuidMataProgramCommandHandler) Handle(
	ctx context.Context,
	cmd SetupUuidMataProgramCommand,
) (string, error) {

	err := h.Repo.SetupUuid(ctx)
	if err != nil {
		return "", err
	}

	return "berhasil setup uuid pada data", nil
}
