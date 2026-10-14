package domain

import (
	commonDomain "UnpakSiamida/common/domain"
	"context"

	"github.com/google/uuid"
)

type ISwotRepository interface {
	GetByUuid(ctx context.Context, uid uuid.UUID) (*Swot, error)
	GetDefaultByUuid(ctx context.Context, uid uuid.UUID) (*SwotDefault, error)
	GetAll(
		ctx context.Context,
		search string,
		searchFilters []commonDomain.SearchFilter,
		page, limit *int,
	) ([]SwotDefault, int64, error)
	Create(ctx context.Context, data *Swot) error
	Update(ctx context.Context, data *Swot) error
	Delete(ctx context.Context, uid uuid.UUID) error
}
