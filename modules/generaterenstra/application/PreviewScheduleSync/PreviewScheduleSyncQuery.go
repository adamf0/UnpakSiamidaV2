package PreviewScheduleSync

import domain "UnpakSiamida/modules/generaterenstra/domain"

type PreviewScheduleSyncQuery struct {
	UUIDs []string
}

type PreviewScheduleSyncResult struct {
	Items                []domain.SchedulePreviewComparison `json:"items"`
	TotalSchedules       int                                `json:"totalSchedules"`
	TotalBeforeRows      int                                `json:"totalBeforeRows"`
	TotalAfterRows       int                                `json:"totalAfterRows"`
	TotalInserted        int                                `json:"totalInserted"`
	TotalDeleted         int                                `json:"totalDeleted"`
	TotalMerged          int                                `json:"totalMerged"`
	TotalUpdatedTugas    int                                `json:"totalUpdatedTugas"`
	RequiredConfirmation string                             `json:"requiredConfirmation"`
}
