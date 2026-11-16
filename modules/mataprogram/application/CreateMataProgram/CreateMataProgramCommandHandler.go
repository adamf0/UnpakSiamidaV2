package application

import (
	"context"
	"strconv"
	"strings"
	"time"

	domainmataprogram "UnpakSiamida/modules/mataprogram/domain"
	domaintahunproker "UnpakSiamida/modules/tahunproker/domain"

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
	if err != nil || tahunproker == nil {
		tahunprokerId = 0
	} else {
		tahunprokerId = tahunproker.ID
	}

	var fuID uint = cmd.IdFakultasUnit
	if fuID == 0 && strings.TrimSpace(cmd.FakultasUnitUuid) != "" {
		trimmed := strings.TrimSpace(cmd.FakultasUnitUuid)
		if parsed, err := uuid.Parse(trimmed); err == nil {
			if id, err := h.Repo.GetFakultasUnitIDByUuid(ctx, parsed); err == nil {
				fuID = id
			}
		} else if id, err := strconv.Atoi(trimmed); err == nil && id > 0 {
			fuID = uint(id)
		}
	}

	result := domainmataprogram.NewMataProgram(
		tahunprokerId,
		cmd.MataProgram,
		fuID,
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
