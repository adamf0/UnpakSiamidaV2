package application

import (
	"context"

	domainjenisfile "UnpakSiamida/modules/jenisfile/domain"
)

type SetupUuidJenisFileCommandHandler struct {
	Repo domainjenisfile.IJenisFileRepository
}

func (h *SetupUuidJenisFileCommandHandler) Handle(
	ctx context.Context,
	cmd SetupUuidJenisFileCommand,
) (string, error) {

	err := h.Repo.SetupUuid(ctx)
	if err != nil {
		return "", err
	}

	return "berhasil setup uuid pada data", nil
}
