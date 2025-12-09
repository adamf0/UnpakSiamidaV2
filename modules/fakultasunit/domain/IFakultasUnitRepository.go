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
    ) ([]FakultasUnit, int64, error)
}
