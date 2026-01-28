package domain

import (
	commonDomain "UnpakSiamida/common/domain"
	"context"

	"github.com/google/uuid"
)

type IDokumenProkerRepository interface {
	GetByUuid(ctx context.Context, uid uuid.UUID) (*DokumenProker, error)
	GetDefaultByUuid(ctx context.Context, uid uuid.UUID) (*DokumenProkerDefault, error)
	GetAll(
		ctx context.Context,
		search string,
		searchFilters []commonDomain.SearchFilter,
		page, limit *int,
	) ([]DokumenProkerDefault, int64, error)
	Create(ctx context.Context, dokumenproker *DokumenProker) error
	Update(ctx context.Context, dokumenproker *DokumenProker) error
	Delete(ctx context.Context, uid uuid.UUID) error
	SetupUuid(ctx context.Context) error
}
