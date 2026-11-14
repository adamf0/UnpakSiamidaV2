package domain

import (
	"context"

	"github.com/google/uuid"
)

type IAktivitasV2Repository interface {
	GetMatrixProkerByUuid(ctx context.Context, uid uuid.UUID) (*MatrixProkerSummary, error)
	GetMatrixProkerById(ctx context.Context, id uint) (*MatrixProkerSummary, error)
	GetAktivitasByMatrikId(ctx context.Context, matrikID uint) ([]AktivitasV2Default, error)
	GetAktivitasByUuid(ctx context.Context, uid uuid.UUID) (*AktivitasV2Default, error)
	GetAktivitasEntityByUuid(ctx context.Context, uid uuid.UUID) (*AktivitasV2, error)
	GetAktivitasEntityById(ctx context.Context, id uint) (*AktivitasV2, error)
	CreateAktivitas(ctx context.Context, data *AktivitasV2, timelines []ProkerTimelineV3) error
	UpdateAktivitas(ctx context.Context, data *AktivitasV2) error
	DeleteAktivitas(ctx context.Context, uid uuid.UUID) error

	// Status
	UpdateStatus(ctx context.Context, req UpdateStatusRequest) (*UpdateStatusResult, error)

	// Timeline
	GetTimelineByUuid(ctx context.Context, uid uuid.UUID) (*ProkerTimelineV3, error)
	CreateTimeline(ctx context.Context, tl *ProkerTimelineV3) error
	UpdateTimeline(ctx context.Context, tl *ProkerTimelineV3) error
	DeleteTimeline(ctx context.Context, uid uuid.UUID) error
	ReorderTimeline(ctx context.Context, items []ReorderTimelineItem) error
	CountTimelineByAktivitasId(ctx context.Context, aktivitasId uint) (int64, error)
	AddTimelineFiles(ctx context.Context, files []ProkerTimelineFileV3) error
	DeleteTimelineFile(ctx context.Context, fileUid uuid.UUID) error
	GetTimelineFileByUuid(ctx context.Context, fileUid uuid.UUID) (*ProkerTimelineFileV3, error)
}
