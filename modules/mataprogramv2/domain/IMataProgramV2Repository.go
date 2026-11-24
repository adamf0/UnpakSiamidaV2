package domain

import (
	"context"
	"time"
)

type IMataProgramV2Repository interface {
	GetAllFlat(ctx context.Context, filter MataProgramV2Filter) ([]MataProgramV2Row, error)
	Count(ctx context.Context, filter MataProgramV2Filter) (int64, error)
	GetPaged(ctx context.Context, filter MataProgramV2Filter, offset, limit int) ([]GroupSummary, error)

	GetRowByUuid(ctx context.Context, rawUuid string) (*MataProgramV2Row, error)
	GetRecordByUuid(ctx context.Context, rawUuid string, includeDeleted bool) (*MataProgramV2Record, error)

	ResolveOrCreateTahunID(ctx context.Context, tahunStr string) (uint, error)
	ResolveTahunID(ctx context.Context, tahunStr string) (uint, error)

	SaveMataProgramV2WithTemplates(ctx context.Context, tahunID uint, mataProgram string, templateRenstraIDs []uint) (string, uint, error)
	UpdateMataProgramV2WithTemplates(ctx context.Context, rawUuid string, tahunID uint, mataProgram string, templateRenstraIDs []uint, replaceTemplates bool) error

	GetTargetUnitsByMataProgramIDs(ctx context.Context, mpIDs []uint) (map[uint][]MataProgramV2UnitBadge, error)
	GetTargetUnitsByMataProgramID(ctx context.Context, mpID uint) ([]MataProgramV2UnitBadge, error)
	GetTemplatesByMataProgramID(ctx context.Context, mpID uint) ([]MataProgramV2TemplateItem, []uint, error)
	GetTemplatesByMataProgramIDs(ctx context.Context, mpIDs []uint) (map[uint][]MataProgramV2TemplateItem, error)
	GetProkerUnitsByMataProgramIDs(ctx context.Context, mpIDs []uint) (map[uint][]MataProgramV2ProkerUnit, error)
	GetProkerUnitsByMataProgramID(ctx context.Context, mpID uint) ([]MataProgramV2ProkerUnit, error)

	SoftDeleteByUuid(ctx context.Context, rawUuid string, deletedAt time.Time) (int64, error)
	RestoreByUuid(ctx context.Context, rawUuid string) (int64, error)
}
