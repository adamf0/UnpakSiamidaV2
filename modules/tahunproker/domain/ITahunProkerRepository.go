package domain

import (
	commonDomain "UnpakSiamida/common/domain"
	"context"

	"github.com/google/uuid"
)

type ITahunProkerRepository interface {
	GetByUuid(ctx context.Context, uid uuid.UUID) (*TahunProker, error)
	GetDefaultByUuid(ctx context.Context, uid uuid.UUID) (*TahunProkerDefault, error)
	GetAll(
		ctx context.Context,
		search string,
		searchFilters []commonDomain.SearchFilter,
		page, limit *int,
	) ([]TahunProker, int64, error)
	Create(ctx context.Context, tahunproker *TahunProker) error
	Update(ctx context.Context, tahunproker *TahunProker) error
	Delete(ctx context.Context, uid uuid.UUID) error
	SetupUuid(ctx context.Context) error
}
