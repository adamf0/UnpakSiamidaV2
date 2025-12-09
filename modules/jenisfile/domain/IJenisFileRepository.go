package domain

import (
	"context"
	commonDomain "UnpakSiamida/common/domain"
	"github.com/google/uuid"
)

type IJenisFileRepository interface {
	GetByUuid(ctx context.Context, uid uuid.UUID) (*JenisFile, error)
	GetDefaultByUuid(ctx context.Context, uid uuid.UUID) (*JenisFileDefault, error)
	GetAll(
        ctx context.Context,
        search string,
        searchFilters []commonDomain.SearchFilter,
        page, limit *int,
    ) ([]JenisFile, int64, error)
}
