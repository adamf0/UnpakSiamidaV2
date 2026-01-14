package application

import (
	"context"

	domaintemplatedokumentambahan "UnpakSiamida/modules/templatedokumentambahan/domain"
)

type SetupUuidTemplateDokumenTambahanCommandHandler struct {
	Repo domaintemplatedokumentambahan.ITemplateDokumenTambahanRepository
}

func (h *SetupUuidTemplateDokumenTambahanCommandHandler) Handle(
	ctx context.Context,
	cmd SetupUuidTemplateDokumenTambahanCommand,
) (string, error) {

	err := h.Repo.SetupUuid(ctx)
	if err != nil {
		return "", err
	}

	return "berhasil setup uuid pada data", nil
}
