package application

import (
	"context"

	domainmataprogram "UnpakSiamida/modules/mataprogram/domain"
	domaintahunproker "UnpakSiamida/modules/tahunproker/domain"
	"time"

	"github.com/google/uuid"
)

type CreateMataProgramCommandHandler struct {
	Repo            domainmataprogram.IMataProgramRepository
	RepoTahunProker domaintahunproker.ITahunProkerRepository
}

func (h *CreateMataProgramCommandHandler) Handle(
	ctx context.Context,
	cmd CreateMataProgramCommand,
) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	tahunprokerUUID, err := uuid.Parse(cmd.TahunUuid)
	if err != nil {
		return "", domainmataprogram.InvalidParseTahun()
	}

	var tahunprokerId uint
	tahunproker, err := h.RepoTahunProker.GetByUuid(ctx, tahunprokerUUID)
	if err != nil {
		tahunprokerId = tahunproker.ID
	} else {
		tahunprokerId = 0
	}

	result := domainmataprogram.NewMataProgram(
		tahunprokerId,
		cmd.MataProgram,
	)

	if !result.IsSuccess {
		return "", result.Error
	}

	createMataProgram := result.Value
	if err := h.Repo.Create(ctx, createMataProgram); err != nil {
		return "", err
	}

	return result.Value.UUID.String(), nil
}
