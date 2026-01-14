package domain

import (
	"context"
	commonDomain "UnpakSiamida/common/domain"
	"github.com/google/uuid"
)

type IRenstraRepository interface {
	IsUnique(ctx context.Context, fakultas_unit uint, tahun string) (bool, error)
	GetByUuid(ctx context.Context, uid uuid.UUID) (*Renstra, error)
	GetDefaultByUuid(ctx context.Context, uid uuid.UUID) (*RenstraDefault, error)
	GetAll(
        ctx context.Context,
        search string,
        searchFilters []commonDomain.SearchFilter,
        page, limit *int,
		scope string,
    ) ([]RenstraDefault, int64, error)
	Create(ctx context.Context, renstra *Renstra) error
	Update(ctx context.Context, renstra *Renstra) error
	Delete(ctx context.Context, uid uuid.UUID) error
	SetupUuid(ctx context.Context) error
}
