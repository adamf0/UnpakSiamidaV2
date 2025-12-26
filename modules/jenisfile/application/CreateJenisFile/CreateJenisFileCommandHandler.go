package application

import (
	"context"
	
	domainjenisfile "UnpakSiamida/modules/jenisfile/domain"
)

type CreateJenisFileCommandHandler struct{
	Repo domainjenisfile.IJenisFileRepository
}

func (h *CreateJenisFileCommandHandler) Handle(
	ctx context.Context,
	cmd CreateJenisFileCommand,
) (string, error) {

	result := domainjenisfile.NewJenisFile(
		cmd.Nama,
	)

	if !result.IsSuccess {
		return "", result.Error
	}

	createJenisFile := result.Value
	if err := h.Repo.Create(ctx, createJenisFile); err != nil {
		return "", err
	}

	return result.Value.UUID.String(), nil
}
