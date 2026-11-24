package RestoreMataProgramV2

import (
	domain "UnpakSiamida/modules/mataprogramv2/domain"
	"context"
	"strings"
)

type RestoreMataProgramV2CommandHandler struct {
	Repo domain.IMataProgramV2Repository
}

func (h *RestoreMataProgramV2CommandHandler) Handle(ctx context.Context, cmd RestoreMataProgramV2Command) (*RestoreMataProgramV2Result, error) {
	if h.Repo == nil {
		return nil, domain.RepoNotInitialized()
	}

	rawUuid := strings.TrimSpace(cmd.UUID)
	if rawUuid == "" {
		return nil, domain.UuidRequired()
	}

	row, err := h.Repo.GetRecordByUuid(ctx, rawUuid, true)
	if err != nil {
		return nil, err
	}
	if row == nil {
		return nil, domain.NotFound()
	}

	rowsAffected, err := h.Repo.RestoreByUuid(ctx, rawUuid)
	if err != nil {
		return nil, err
	}

	return &RestoreMataProgramV2Result{
		Status:       "success",
		Message:      "Berhasil memulihkan data Mata Program V2",
		RowsAffected: rowsAffected,
	}, nil
}
