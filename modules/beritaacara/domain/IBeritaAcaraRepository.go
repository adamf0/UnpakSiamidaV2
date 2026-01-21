package domain

import (
	commonDomain "UnpakSiamida/common/domain"
	"context"

	"github.com/google/uuid"
)

type IBeritaAcaraRepository interface {
	GetByUuid(ctx context.Context, uid uuid.UUID) (*BeritaAcara, error)
	GetDefaultByUuid(ctx context.Context, uid uuid.UUID) (*BeritaAcaraDefault, error)
	GetAll(
		ctx context.Context,
		search string,
		searchFilters []commonDomain.SearchFilter,
		page, limit *int,
	) ([]BeritaAcara, int64, error)
	Create(ctx context.Context, BeritaAcara *BeritaAcara) error
	Update(ctx context.Context, BeritaAcara *BeritaAcara) error
	Delete(ctx context.Context, uid uuid.UUID) error
	SetupUuid(ctx context.Context) error
}
