package application

import (
	"context"

	domaintemplaterenstra "UnpakSiamida/modules/templaterenstra/domain"
)

type SetupUuidTemplateRenstraCommandHandler struct {
	Repo domaintemplaterenstra.ITemplateRenstraRepository
}

func (h *SetupUuidTemplateRenstraCommandHandler) Handle(
	ctx context.Context,
	cmd SetupUuidTemplateRenstraCommand,
) (string, error) {

	err := h.Repo.SetupUuid(ctx)
	if err != nil {
		return "", err
	}

	return "berhasil setup uuid pada data", nil
}
