package domain

import (
	"context"
	commonDomain "UnpakSiamida/common/domain"
	"github.com/google/uuid"
)

type IFakultasUnitRepository interface {
	GetDefaultByUuid(ctx context.Context, uid uuid.UUID) (*FakultasUnit, error)
	GetAll(
		ctx context.Context,
		search string,
		searchFilters []commonDomain.SearchFilter,
		page, limit *int,
		includeHidden bool,
	) ([]FakultasUnit, int64, error)
	UpdateTarget(ctx context.Context, uid uuid.UUID, nama *string, isActive *bool) error
	SetupUuid(ctx context.Context) error
	EnsureSchema(ctx context.Context) error
}
