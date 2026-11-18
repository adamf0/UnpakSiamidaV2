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

	scope := strings.TrimSpace(cmd.Scope)
	if scope == "" {
		scope = "single"
	}

	row, err := h.Repo.GetRecordByUuid(ctx, rawUuid, true)
	if err != nil {
		return nil, err
	}
	if row == nil {
		return nil, domain.NotFound()
	}

	var rowsAffected int64
	if scope == "single" {
		rowsAffected, err = h.Repo.RestoreSingleByUuid(ctx, rawUuid)
	} else {
		rowsAffected, err = h.Repo.RestoreGroup(ctx, row.MataProgram, row.IdMasterTahun)
	}
	if err != nil {
		return nil, err
	}

	return &RestoreMataProgramV2Result{
		Status:       "success",
		Message:      "Berhasil memulihkan data Mata Program V2",
		RowsAffected: rowsAffected,
	}, nil
}
