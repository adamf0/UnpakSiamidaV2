package domain

import (
	"context"
	commonDomain "UnpakSiamida/common/domain"
	"github.com/google/uuid"
)

type IRenstraNilaiRepository interface {
	GetByUuid(ctx context.Context, uid uuid.UUID) (*RenstraNilai, error)
	GetDefaultByUuid(ctx context.Context, uid uuid.UUID) (*RenstraNilaiDefault, error)
	GetAll(
        ctx context.Context,
        search string,
        searchFilters []commonDomain.SearchFilter,
        page, limit *int,
    ) ([]RenstraNilaiDefault, int64, error)
	Update(ctx context.Context, renstranilai *RenstraNilai) error
	Delete(ctx context.Context, uid uuid.UUID) error
	SetupUuid(ctx context.Context) error
}
