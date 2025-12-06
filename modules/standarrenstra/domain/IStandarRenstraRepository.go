package domain

import (
	"context"
	commonDomain "UnpakSiamida/common/domain"
	"github.com/google/uuid"
)

type IStandarRenstraRepository interface {
	GetByUuid(ctx context.Context, uid uuid.UUID) (*StandarRenstra, error)
	GetAll(
        ctx context.Context,
        search string,
        searchFilters []commonDomain.SearchFilter,
        page, limit *int,
    ) ([]StandarRenstra, int64, error)
	Create(ctx context.Context, standarrenstra *StandarRenstra) error
	Update(ctx context.Context, standarrenstra *StandarRenstra) error
	Delete(ctx context.Context, uid uuid.UUID) error
}
