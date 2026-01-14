package application

import (
	"context"

	domainkts "UnpakSiamida/modules/kts/domain"
)

type SetupUuidKtsCommandHandler struct {
	Repo domainkts.IKtsRepository
}

func (h *SetupUuidKtsCommandHandler) Handle(
	ctx context.Context,
	cmd SetupUuidKtsCommand,
) (string, error) {

	err := h.Repo.SetupUuid(ctx)
	if err != nil {
		return "", err
	}

	return "berhasil setup uuid pada data", nil
}
