package domain

import (
	"context"
	commonDomain "UnpakSiamida/common/domain"
	"github.com/google/uuid"
)

type ITemplateRenstraRepository interface {
	GetByUuid(ctx context.Context, uid uuid.UUID) (*TemplateRenstra, error)
	GetAll(
        ctx context.Context,
        search string,
        searchFilters []commonDomain.SearchFilter,
        page, limit *int,
    ) ([]TemplateRenstraDefault, int64, error)
	GetAllByTahunFakUnit(ctx context.Context, tahun string, fakultasUnit uint) ([]TemplateRenstra, error)
	GetAllByTahunFakUnitDefault(ctx context.Context, tahun string, fakultasUnit uint) ([]TemplateRenstraDefault, error)
	Create(ctx context.Context, templaterenstra *TemplateRenstra) error
	Update(ctx context.Context, templaterenstra *TemplateRenstra) error
	Delete(ctx context.Context, uid uuid.UUID) error
	SetupUuid(ctx context.Context) error
}
