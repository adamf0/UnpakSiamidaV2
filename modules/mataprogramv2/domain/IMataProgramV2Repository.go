package domain

import (
	"context"
	"time"
)

type IMataProgramV2Repository interface {
	GetAllFlat(ctx context.Context, filter MataProgramV2Filter) ([]MataProgramV2Row, error)
	CountGroups(ctx context.Context, filter MataProgramV2Filter) (int64, error)
	GetGroupSummaries(ctx context.Context, filter MataProgramV2Filter, offset, limit int) ([]GroupSummary, error)
	GetGroupUnits(ctx context.Context, mataProgram string, tahunID *uint, showDeleted bool) ([]MataProgramV2UnitBadge, error)

	GetRowByUuid(ctx context.Context, rawUuid string) (*MataProgramV2Row, error)
	GetRecordByUuid(ctx context.Context, rawUuid string, includeDeleted bool) (*MataProgramV2Record, error)

	ResolveOrCreateTahunID(ctx context.Context, tahunStr string) (uint, error)
	ResolveTahunID(ctx context.Context, tahunStr string) (uint, error)
	ResolveFakultasUnitID(ctx context.Context, rawTarget string) (uint, error)

	UpsertOrRestoreAssignment(ctx context.Context, tahunID uint, unitID uint, mataProgram string) (string, bool, error)
	SaveMataProgramV2WithTemplates(ctx context.Context, tahunID uint, mataProgram string, templateRenstraIDs []uint) (string, uint, error)
	UpdateMataProgramV2WithTemplates(ctx context.Context, rawUuid string, tahunID uint, mataProgram string, templateRenstraIDs []uint) error
	GetTemplatesByMataProgramID(ctx context.Context, mpID uint) ([]MataProgramV2TemplateItem, []uint, error)
	GetTemplatesByMataProgramIDs(ctx context.Context, mpIDs []uint) (map[uint][]MataProgramV2TemplateItem, error)
	GetProkerUnitsByMataProgramIDs(ctx context.Context, mpIDs []uint) (map[uint][]MataProgramV2ProkerUnit, error)
	GetProkerUnitsByMataProgramID(ctx context.Context, mpID uint) ([]MataProgramV2ProkerUnit, error)

	UpdateGroupFields(ctx context.Context, originalMataProgram string, originalTahunID *uint, updates map[string]interface{}) error
	SoftDeleteGroupUnitsNotIn(ctx context.Context, mataProgram string, tahunID uint, keepUnitIDs []uint, deletedAt time.Time) error
	UpdateSingleByUuid(ctx context.Context, rawUuid string, updates map[string]interface{}) error

	SoftDeleteSingleByUuid(ctx context.Context, rawUuid string, deletedAt time.Time) (int64, error)
	SoftDeleteGroup(ctx context.Context, mataProgram string, tahunID *uint, deletedAt time.Time) (int64, error)
	RestoreSingleByUuid(ctx context.Context, rawUuid string) (int64, error)
	RestoreGroup(ctx context.Context, mataProgram string, tahunID *uint) (int64, error)
}
