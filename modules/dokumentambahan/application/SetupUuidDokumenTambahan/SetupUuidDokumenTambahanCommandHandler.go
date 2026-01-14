package application

import (
	"context"

	domaindokumentambahan "UnpakSiamida/modules/dokumentambahan/domain"
)

type SetupUuidDokumenTambahanCommandHandler struct {
	Repo domaindokumentambahan.IDokumenTambahanRepository
}

func (h *SetupUuidDokumenTambahanCommandHandler) Handle(
	ctx context.Context,
	cmd SetupUuidDokumenTambahanCommand,
) (string, error) {

	err := h.Repo.SetupUuid(ctx)
	if err != nil {
		return "", err
	}

	return "berhasil setup uuid pada data", nil
}
