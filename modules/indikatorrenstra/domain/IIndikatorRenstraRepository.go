package domain

import (
	"context"
	commonDomain "UnpakSiamida/common/domain"
	"github.com/google/uuid"
)

type IIndikatorRenstraRepository interface {
	IsUniqueIndikator(ctx context.Context, idnikator string, tahun string) (bool, error)
	GetByUuid(ctx context.Context, uid uuid.UUID) (*IndikatorRenstra, error)
	GetDefaultByUuid(ctx context.Context, uid uuid.UUID) (*IndikatorRenstraDefault, error)
	GetAll(
        ctx context.Context,
        search string,
        searchFilters []commonDomain.SearchFilter,
        page, limit *int,
    ) ([]IndikatorRenstra, int64, error)
	Create(ctx context.Context, indikatorrenstra *IndikatorRenstra) error
	Update(ctx context.Context, indikatorrenstra *IndikatorRenstra) error
	Delete(ctx context.Context, uid uuid.UUID) error
}
