package domain

import (
	commondomain "UnpakSiamida/common/domain"
	"context"

	"github.com/google/uuid"
)

type IProkerV3Repository interface {
	GetByUuid(ctx context.Context, uid uuid.UUID) (*ProkerV3, error)
	GetDefaultByUuid(ctx context.Context, uid uuid.UUID) (*ProkerV3Default, error)
	GetAll(ctx context.Context, search string, searchFilters []commondomain.SearchFilter, page, limit *int) ([]ProkerV3Default, int64, error)
	Create(ctx context.Context, data *ProkerV3) error
	Update(ctx context.Context, data *ProkerV3) error
	Delete(ctx context.Context, uid uuid.UUID) error

	// Aktivitas V2
	GetAktivitasByMatrikId(ctx context.Context, matrikID uint) ([]AktivitasV2Default, error)
	GetAktivitasByUuid(ctx context.Context, uid uuid.UUID) (*AktivitasV2Default, error)
	CreateAktivitas(ctx context.Context, data *AktivitasV2, timelines []ProkerTimelineV3) error
	UpdateAktivitas(ctx context.Context, data *AktivitasV2) error
	DeleteAktivitas(ctx context.Context, uid uuid.UUID) error

	// Timeline
	GetTimelineByUuid(ctx context.Context, uid uuid.UUID) (*ProkerTimelineV3, error)
	AddTimelineFiles(ctx context.Context, files []ProkerTimelineFileV3) error
	DeleteTimelineFile(ctx context.Context, fileUid uuid.UUID) error
	GetTimelineFileByUuid(ctx context.Context, fileUid uuid.UUID) (*ProkerTimelineFileV3, error)
}
