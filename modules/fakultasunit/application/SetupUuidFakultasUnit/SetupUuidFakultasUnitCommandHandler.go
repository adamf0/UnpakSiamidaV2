package application

import (
	"context"

	domainfakultasunit "UnpakSiamida/modules/fakultasunit/domain"
)

type SetupUuidFakultasUnitCommandHandler struct {
	Repo domainfakultasunit.IFakultasUnitRepository
}

func (h *SetupUuidFakultasUnitCommandHandler) Handle(
	ctx context.Context,
	cmd SetupUuidFakultasUnitCommand,
) (string, error) {

	err := h.Repo.SetupUuid(ctx)
	if err != nil {
		return "", err
	}

	return "berhasil setup uuid pada data", nil
}
