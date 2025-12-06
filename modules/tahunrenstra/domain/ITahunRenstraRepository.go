package domain

import (
	"context"
	commonDomain "UnpakSiamida/common/domain"
)

type ITahunRenstraRepository interface {
	GetActive(ctx context.Context) (*TahunRenstra, error)
	GetAll(
        ctx context.Context,
        search string,
        searchFilters []commonDomain.SearchFilter,
        page, limit *int,
    ) ([]TahunRenstra, int64, error)
}
