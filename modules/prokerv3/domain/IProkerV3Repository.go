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

	// Options & Target
	GetOptions(ctx context.Context) (*ProkerOptionsResult, error)
	GetTargetUnpak(ctx context.Context, tahun, indikatorParam, unitParam string) (*TargetUnpakResult, error)

	// Aktivitas V2
	GetAktivitasByMatrikId(ctx context.Context, matrikID uint) ([]AktivitasV2Default, error)
	GetAktivitasByUuid(ctx context.Context, uid uuid.UUID) (*AktivitasV2Default, error)
	GetAktivitasEntityByUuid(ctx context.Context, uid uuid.UUID) (*AktivitasV2, error)
	CreateAktivitas(ctx context.Context, data *AktivitasV2, timelines []ProkerTimelineV3) error
	UpdateAktivitas(ctx context.Context, data *AktivitasV2) error
	DeleteAktivitas(ctx context.Context, uid uuid.UUID) error

	// Status & History
	UpdateStatus(ctx context.Context, req UpdateStatusRequest) (*UpdateStatusResult, error)
	GetHistoryState(ctx context.Context, entityUuid string) ([]HistoryState, error)
	RecordHistoryState(ctx context.Context, h *HistoryState) error

	// Timeline
	GetTimelineByUuid(ctx context.Context, uid uuid.UUID) (*ProkerTimelineV3, error)
	CreateTimeline(ctx context.Context, tl *ProkerTimelineV3) error
	UpdateTimeline(ctx context.Context, tl *ProkerTimelineV3) error
	DeleteTimeline(ctx context.Context, uid uuid.UUID) error
	ReorderTimeline(ctx context.Context, items []ReorderTimelineItem) error
	AddTimelineFiles(ctx context.Context, files []ProkerTimelineFileV3) error
	DeleteTimelineFile(ctx context.Context, fileUid uuid.UUID) error
	GetTimelineFileByUuid(ctx context.Context, fileUid uuid.UUID) (*ProkerTimelineFileV3, error)

	// Dokumen
	GetDokumenByMatrikId(ctx context.Context, matrikID uint) ([]MatrixProkerDokumen, error)
	GetDokumenByUuid(ctx context.Context, uid uuid.UUID) (*MatrixProkerDokumen, error)
	CreateDokumen(ctx context.Context, doc *MatrixProkerDokumen) error
	DeleteDokumen(ctx context.Context, uid uuid.UUID) error
	UpdateVerifikasiDokumen(ctx context.Context, uid uuid.UUID, status, catatan, catatanPenolakan string) (*MatrixProkerDokumen, error)
}
