package application

import (
	"context"
	
	domainjenisfile "UnpakSiamida/modules/jenisfile/domain"
	"time"
)

type CreateJenisFileCommandHandler struct{
	Repo domainjenisfile.IJenisFileRepository
}

func (h *CreateJenisFileCommandHandler) Handle(
	ctx context.Context,
	cmd CreateJenisFileCommand,
) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

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
