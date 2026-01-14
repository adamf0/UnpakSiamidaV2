package domain

import (
	"context"
	commonDomain "UnpakSiamida/common/domain"
	"github.com/google/uuid"
)

type IDokumenTambahanRepository interface {
	GetByUuid(ctx context.Context, uid uuid.UUID) (*DokumenTambahan, error)
	GetDefaultByUuid(ctx context.Context, uid uuid.UUID) (*DokumenTambahanDefault, error)
	GetAll(
        ctx context.Context,
        search string,
        searchFilters []commonDomain.SearchFilter,
        page, limit *int,
    ) ([]DokumenTambahanDefault, int64, error)
	Update(ctx context.Context, dokumenTambahan *DokumenTambahan) error
	Delete(ctx context.Context, uid uuid.UUID) error
	SetupUuid(ctx context.Context) error
}
