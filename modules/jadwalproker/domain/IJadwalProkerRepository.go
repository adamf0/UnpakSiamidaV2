package domain

import (
	commonDomain "UnpakSiamida/common/domain"
	"context"

	"github.com/google/uuid"
)

type IJadwalProkerRepository interface {
	GetByUuid(ctx context.Context, uid uuid.UUID) (*JadwalProker, error)
	GetDefaultByUuid(ctx context.Context, uid uuid.UUID) (*JadwalProkerDefault, error)
	GetAll(
		ctx context.Context,
		search string,
		searchFilters []commonDomain.SearchFilter,
		page, limit *int,
	) ([]JadwalProkerDefault, int64, error)
	Create(ctx context.Context, jadwalproker *JadwalProker) error
	Update(ctx context.Context, jadwalproker *JadwalProker) error
	Delete(ctx context.Context, uid uuid.UUID) error
	SetupUuid(ctx context.Context) error
}
