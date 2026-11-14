package domain

import "context"

type ISyncRenstraRepository interface {
	LoadRenstraRecordByUUID(ctx context.Context, renstraUUID string) (*RenstraRecord, error)
	BuildSchedulePreview(ctx context.Context, renstraUUID string) (*SchedulePreviewComparison, error)
	ExecuteLegacySyncRenstra(ctx context.Context, renstraUUID string) (*RenstraRecord, error)
}
