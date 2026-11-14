package ExecuteScheduleSync

import (
	"context"
	"fmt"

	domain "UnpakSiamida/modules/generaterenstra/domain"
)

type ExecuteScheduleSyncCommandHandler struct {
	Repo domain.ISyncRenstraRepository
}

func (h *ExecuteScheduleSyncCommandHandler) Handle(ctx context.Context, cmd ExecuteScheduleSyncCommand) (*ExecuteScheduleSyncResult, error) {
	if h.Repo == nil {
		return nil, domain.RepoNotInitialized()
	}

	if len(cmd.UUIDs) == 0 {
		return nil, domain.PilihMinimalSatuJadwal()
	}

	if cmd.Confirmation != domain.RequiredGenerateConfirmationPhrase {
		return nil, domain.InvalidConfirmationPhrase(domain.RequiredGenerateConfirmationPhrase)
	}

	var processedTargets []string
	var processedUUIDs []string

	for _, u := range cmd.UUIDs {
		rec, err := h.Repo.ExecuteLegacySyncRenstra(ctx, u)
		if err != nil {
			return nil, err
		}
		processedUUIDs = append(processedUUIDs, rec.UUID)
		processedTargets = append(processedTargets, fmt.Sprintf("%s (%s)", rec.NamaFakProdUnit, rec.Tahun))
	}

	return &ExecuteScheduleSyncResult{
		Message:          fmt.Sprintf("Berhasil generate data pada %d jadwal audit.", len(processedUUIDs)),
		ProcessedCount:   len(processedUUIDs),
		ProcessedUUIDs:   processedUUIDs,
		ProcessedTargets: processedTargets,
	}, nil
}
