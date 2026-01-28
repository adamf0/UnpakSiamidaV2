package domain

import (
	commonDomain "UnpakSiamida/common/domain"
	"context"

	"github.com/google/uuid"
)

type IAktivitasProkerRepository interface {
	GetByUuid(ctx context.Context, uid uuid.UUID) (*AktivitasProker, error)
	GetDefaultByUuid(ctx context.Context, uid uuid.UUID) (*AktivitasProkerDefault, error)
	GetAll(
		ctx context.Context,
		search string,
		searchFilters []commonDomain.SearchFilter,
		page, limit *int,
	) ([]AktivitasProkerDefault, int64, error)
	Create(ctx context.Context, aktivitasproker *AktivitasProker) error
	Update(ctx context.Context, aktivitasproker *AktivitasProker) error
	Delete(ctx context.Context, uid uuid.UUID) error
	SetupUuid(ctx context.Context) error
}
