package PreviewScheduleSync

import (
	"context"

	domain "UnpakSiamida/modules/generaterenstra/domain"
)

type PreviewScheduleSyncQueryHandler struct {
	Repo domain.ISyncRenstraRepository
}

func (h *PreviewScheduleSyncQueryHandler) Handle(ctx context.Context, query PreviewScheduleSyncQuery) (*PreviewScheduleSyncResult, error) {
	if h.Repo == nil {
		return nil, domain.RepoNotInitialized()
	}

	if len(query.UUIDs) == 0 {
		return nil, domain.PilihMinimalSatuJadwal()
	}

	var items []domain.SchedulePreviewComparison
	var totalBeforeRows, totalAfterRows, totalInserted, totalDeleted, totalMerged, totalUpdatedTugas int

	for _, u := range query.UUIDs {
		comp, err := h.Repo.BuildSchedulePreview(ctx, u)
		if err != nil {
			return nil, err
		}
		items = append(items, *comp)
		totalBeforeRows += comp.Summary.BeforeTotalRows
		totalAfterRows += comp.Summary.AfterTotalRows
		totalInserted += comp.Summary.InsertedCount
		totalDeleted += comp.Summary.DeletedCount
		totalMerged += comp.Summary.MergedGroupCount
		totalUpdatedTugas += comp.Summary.UpdatedTugasCount
	}

	return &PreviewScheduleSyncResult{
		Items:                items,
		TotalSchedules:       len(items),
		TotalBeforeRows:      totalBeforeRows,
		TotalAfterRows:       totalAfterRows,
		TotalInserted:        totalInserted,
		TotalDeleted:         totalDeleted,
		TotalMerged:          totalMerged,
		TotalUpdatedTugas:    totalUpdatedTugas,
		RequiredConfirmation: domain.RequiredGenerateConfirmationPhrase,
	}, nil
}
