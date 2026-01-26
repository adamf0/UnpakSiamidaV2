package domain

import (
	commonDomain "UnpakSiamida/common/domain"
	"context"

	"github.com/google/uuid"
)

type IMataProgramRepository interface {
	GetByUuid(ctx context.Context, uid uuid.UUID) (*MataProgram, error)
	GetDefaultByUuid(ctx context.Context, uid uuid.UUID) (*MataProgramDefault, error)
	GetAll(
		ctx context.Context,
		search string,
		searchFilters []commonDomain.SearchFilter,
		page, limit *int,
	) ([]MataProgram, int64, error)
	Create(ctx context.Context, tahunproker *MataProgram) error
	Update(ctx context.Context, tahunproker *MataProgram) error
	Delete(ctx context.Context, uid uuid.UUID) error
	SetupUuid(ctx context.Context) error
}
