package application

import (
	"context"

	domainrenstranilai "UnpakSiamida/modules/renstranilai/domain"
)

type SetupUuidRenstraNilaiCommandHandler struct {
	Repo domainrenstranilai.IRenstraNilaiRepository
}

func (h *SetupUuidRenstraNilaiCommandHandler) Handle(
	ctx context.Context,
	cmd SetupUuidRenstraNilaiCommand,
) (string, error) {

	err := h.Repo.SetupUuid(ctx)
	if err != nil {
		return "", err
	}

	return "berhasil setup uuid pada data", nil
}
