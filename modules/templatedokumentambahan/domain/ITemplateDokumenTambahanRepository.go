package domain

import (
	"context"
	commonDomain "UnpakSiamida/common/domain"
	"github.com/google/uuid"
)

type ITemplateDokumenTambahanRepository interface {
	GetByUuid(ctx context.Context, uid uuid.UUID) (*TemplateDokumenTambahan, error)
	GetAll(
        ctx context.Context,
        search string,
        searchFilters []commonDomain.SearchFilter,
        page, limit *int,
    ) ([]TemplateDokumenTambahan, int64, error)
	Create(ctx context.Context, templatedokumentambahan *TemplateDokumenTambahan) error
	Update(ctx context.Context, templatedokumentambahan *TemplateDokumenTambahan) error
	Delete(ctx context.Context, uid uuid.UUID) error
}
