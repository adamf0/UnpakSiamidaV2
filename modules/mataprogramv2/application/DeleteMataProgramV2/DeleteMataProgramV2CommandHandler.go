package DeleteMataProgramV2

import (
	domain "UnpakSiamida/modules/mataprogramv2/domain"
	"context"
	"strings"
	"time"
)

type DeleteMataProgramV2CommandHandler struct {
	Repo domain.IMataProgramV2Repository
}

func (h *DeleteMataProgramV2CommandHandler) Handle(ctx context.Context, cmd DeleteMataProgramV2Command) (*DeleteMataProgramV2Result, error) {
	if h.Repo == nil {
		return nil, domain.RepoNotInitialized()
	}

	rawUuid := strings.TrimSpace(cmd.UUID)
	if rawUuid == "" {
		return nil, domain.UuidRequired()
	}

	row, err := h.Repo.GetRecordByUuid(ctx, rawUuid, false)
	if err != nil {
		return nil, err
	}
	if row == nil {
		return nil, domain.NotFoundOrDeleted()
	}

	now := time.Now()
	var rowsAffected int64
	if strings.TrimSpace(cmd.Scope) == "single" {
		rowsAffected, err = h.Repo.SoftDeleteSingleByUuid(ctx, rawUuid, now)
	} else {
		rowsAffected, err = h.Repo.SoftDeleteGroup(ctx, row.MataProgram, row.IdMasterTahun, now)
	}
	if err != nil {
		return nil, err
	}

	return &DeleteMataProgramV2Result{
		Status:       "success",
		Message:      "Berhasil menghapus data Mata Program V2",
		RowsAffected: rowsAffected,
	}, nil
}
