package domain

import (
	"context"
	commonDomain "UnpakSiamida/common/domain"
	"github.com/google/uuid"
)

type IKtsRepository interface {
	GetByUuid(ctx context.Context, uid uuid.UUID) (*Kts, error)
	GetDefaultByUuid(ctx context.Context, uid uuid.UUID) (*KtsDefault, error)
	GetAll(
        ctx context.Context,
        search string,
        searchFilters []commonDomain.SearchFilter,
        page, limit *int,
    ) ([]KtsDefault, int64, error)
	Create(ctx context.Context, kts *Kts) error
	Update(ctx context.Context, kts *Kts) error
	Delete(ctx context.Context, uid uuid.UUID) error
	SetupUuid(ctx context.Context) error
	GetTemuanKts(ctx context.Context, renstraUuid string, tugas string, tahun string) ([]TemuanKts, error)
	GenerateKtsFromTemuan(ctx context.Context, renstraNilaiId *uint, renstraNilaiUuid *string, dokumenTambahanId *uint, dokumenTambahanUuid *string, auditorId string) (*uuid.UUID, error)
	WithTx(ctx context.Context, fn func(txRepo IKtsRepositoryTx) error) error
}

