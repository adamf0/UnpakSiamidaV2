package application

import (
	"context"

	domaintahunproker "UnpakSiamida/modules/tahunproker/domain"
	"time"
)

type CreateTahunProkerCommandHandler struct {
	Repo domaintahunproker.ITahunProkerRepository
}

func (h *CreateTahunProkerCommandHandler) Handle(
	ctx context.Context,
	cmd CreateTahunProkerCommand,
) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	result := domaintahunproker.NewTahunProker(
		cmd.Tahun,
		cmd.Status,
	)

	if !result.IsSuccess {
		return "", result.Error
	}

	createTahunProker := result.Value
	if err := h.Repo.Create(ctx, createTahunProker); err != nil {
		return "", err
	}

	return result.Value.UUID.String(), nil
}
