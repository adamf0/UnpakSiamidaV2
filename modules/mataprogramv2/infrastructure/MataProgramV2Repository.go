package infrastructure

import (
	domain "UnpakSiamida/modules/mataprogramv2/domain"
	"context"
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type MataProgramV2Repository struct {
	db *gorm.DB
}

func NewMataProgramV2Repository(db *gorm.DB) domain.IMataProgramV2Repository {
	return &MataProgramV2Repository{db: db}
}

const flatSelectSQL = `
	mp.id,
	mp.uuid,
	mp.id_master_tahun,
	COALESCE(mt.uuid, '') as tahun_uuid,
	COALESCE(mt.tahun, '') as tahun,
	mp.mata_program,
	mp.mata_program as nama,
	mp.id_fakultas_unit,
	COALESCE(vfu.uuid, '') as fakultas_unit_uuid,
	CASE 
		WHEN vfu.type = 'prodi' AND vfu.jenjang IS NOT NULL AND vfu.jenjang != '' 
		THEN CONCAT(vfu.nama_fak_prod_unit, ' (', UPPER(vfu.jenjang), ')') 
		ELSE COALESCE(vfu.nama_fak_prod_unit, '') 
	END as fakultas_unit_label,
	COALESCE(vfu.jenjang, '') as jenjang,
	COALESCE(vfu.type, '') as type,
	COALESCE(vfu.fakultas, '') as fakultas,
	COALESCE((SELECT COUNT(*) FROM matrix_proker mpk WHERE mpk.id_mata_program = mp.id AND mpk.deleted_at IS NULL), 0) AS proker_count,
	mp.created_at,
	mp.updated_at,
	mp.deleted_at
`

const unitBadgeSelectSQL = `
	mp.id,
	mp.uuid,
	mp.id_fakultas_unit,
	COALESCE(vfu.uuid, '') as fakultas_unit_uuid,
	CASE 
		WHEN vfu.type = 'prodi' AND vfu.jenjang IS NOT NULL AND vfu.jenjang != '' 
		THEN CONCAT(vfu.nama_fak_prod_unit, ' (', UPPER(vfu.jenjang), ')') 
		ELSE COALESCE(vfu.nama_fak_prod_unit, '') 
	END as fakultas_unit_label,
	COALESCE(vfu.jenjang, '') as jenjang,
	COALESCE(vfu.type, '') as type,
	COALESCE(vfu.fakultas, '') as fakultas,
	COALESCE((SELECT COUNT(*) FROM matrix_proker mpk WHERE mpk.id_mata_program = mp.id AND mpk.deleted_at IS NULL), 0) AS proker_count,
	mp.deleted_at
`

func (r *MataProgramV2Repository) applyBaseFilters(q *gorm.DB, filter domain.MataProgramV2Filter) *gorm.DB {
	if filter.ShowDeleted {
		q = q.Where("mp.deleted_at IS NOT NULL")
	} else {
		q = q.Where("mp.deleted_at IS NULL")
	}

	for _, f := range filter.Filters {
		val := strings.ToLower(f.Value)
		switch f.Field {
		case "mataprogram", "mata_program", "nama", "namaprogramkerja":
			col := "LOWER(mp.mata_program)"
			switch f.Op {
			case "eq":
				q = q.Where(col+" = ?", val)
			case "neq":
				q = q.Where(col+" <> ?", val)
			default:
				q = q.Where(col+" LIKE ?", "%"+val+"%")
			}
		case "tahun", "mastertahun", "tahunrenstra":
			col := "LOWER(mt.tahun)"
			switch f.Op {
			case "eq":
				q = q.Where(col+" = ?", val)
			case "neq":
				q = q.Where(col+" <> ?", val)
			default:
				q = q.Where(col+" LIKE ?", "%"+val+"%")
			}
		case "fakultasunit", "fakultas_unit", "unit", "stakeholder", "target":
			like := "%" + val + "%"
			q = q.Where("(LOWER(vfu.nama_fak_prod_unit) LIKE ? OR LOWER(vfu.jenjang) LIKE ? OR LOWER(vfu.fakultas) LIKE ?)", like, like, like)
		default:
			q = q.Where("LOWER(mp.mata_program) LIKE ?", "%"+val+"%")
		}
	}

	if filter.Tahun != "" && filter.Tahun != "all" {
		if _, err := uuid.Parse(filter.Tahun); err == nil {
			q = q.Where("mt.uuid = ?", filter.Tahun)
		} else if id, err := strconv.Atoi(filter.Tahun); err == nil && id > 0 {
			q = q.Where("(mp.id_master_tahun = ? OR mt.tahun = ?)", id, filter.Tahun)
		} else {
			q = q.Where("mt.tahun = ?", filter.Tahun)
		}
	}

	if filter.Unit != "" && filter.Unit != "all" {
		if _, err := uuid.Parse(filter.Unit); err == nil {
			q = q.Where("vfu.uuid = ?", filter.Unit)
		} else if id, err := strconv.Atoi(filter.Unit); err == nil && id > 0 {
			q = q.Where("mp.id_fakultas_unit = ?", id)
		}
	}

	if filter.Search != "" {
		like := "%" + strings.ToLower(filter.Search) + "%"
		q = q.Where("(LOWER(mp.mata_program) LIKE ? OR LOWER(vfu.nama_fak_prod_unit) LIKE ? OR LOWER(mt.tahun) LIKE ? OR LOWER(vfu.fakultas) LIKE ?)", like, like, like, like)
	}

	return q
}

func (r *MataProgramV2Repository) baseJoinedQuery(ctx context.Context, filter domain.MataProgramV2Filter) *gorm.DB {
	q := r.db.WithContext(ctx).
		Table("mata_programv2 mp").
		Joins("LEFT JOIN master_tahun mt ON mp.id_master_tahun = mt.id").
		Joins("LEFT JOIN v_fakultas_unit vfu ON mp.id_fakultas_unit = vfu.id")
	return r.applyBaseFilters(q, filter)
}

func (r *MataProgramV2Repository) GetAllFlat(ctx context.Context, filter domain.MataProgramV2Filter) ([]domain.MataProgramV2Row, error) {
	var list []domain.MataProgramV2Row
	err := r.baseJoinedQuery(ctx, filter).
		Select(flatSelectSQL).
		Order("mp.mata_program ASC, mp.id ASC").
		Scan(&list).Error
	if err != nil {
		return nil, err
	}
	return list, nil
}

func (r *MataProgramV2Repository) CountGroups(ctx context.Context, filter domain.MataProgramV2Filter) (int64, error) {
	var total int64
	sub := r.baseJoinedQuery(ctx, filter).
		Select("mp.id_master_tahun, mp.mata_program").
		Group("mp.id_master_tahun, mp.mata_program")
	err := r.db.WithContext(ctx).Table("(?) as grps", sub).Count(&total).Error
	return total, err
}

func (r *MataProgramV2Repository) GetGroupSummaries(ctx context.Context, filter domain.MataProgramV2Filter, offset, limit int) ([]domain.GroupSummary, error) {
	var summaries []domain.GroupSummary
	err := r.baseJoinedQuery(ctx, filter).
		Select(`
			MIN(mp.id) as id,
			MIN(mp.uuid) as uuid,
			mp.id_master_tahun,
			mp.mata_program,
			COALESCE(mt.uuid, '') as tahun_uuid,
			COALESCE(mt.tahun, '') as tahun,
			MAX(mp.created_at) as created_at,
			MAX(mp.updated_at) as updated_at,
			MAX(mp.deleted_at) as deleted_at
		`).
		Group("mp.id_master_tahun, mp.mata_program").
		Order("id DESC").
		Offset(offset).
		Limit(limit).
		Scan(&summaries).Error
	if err != nil {
		return nil, err
	}
	return summaries, nil
}

func (r *MataProgramV2Repository) GetGroupUnits(ctx context.Context, mataProgram string, tahunID *uint, showDeleted bool) ([]domain.MataProgramV2UnitBadge, error) {
	var units []domain.MataProgramV2UnitBadge
	q := r.db.WithContext(ctx).
		Table("mata_programv2 mp").
		Select(unitBadgeSelectSQL).
		Joins("LEFT JOIN v_fakultas_unit vfu ON mp.id_fakultas_unit = vfu.id").
		Where("mp.mata_program = ?", mataProgram)

	if showDeleted {
		q = q.Where("mp.deleted_at IS NOT NULL")
	} else {
		q = q.Where("mp.deleted_at IS NULL")
	}

	if tahunID != nil {
		q = q.Where("mp.id_master_tahun = ?", *tahunID)
	} else {
		q = q.Where("mp.id_master_tahun IS NULL")
	}

	err := q.Order("vfu.type ASC, fakultas_unit_label ASC").Scan(&units).Error
	return units, err
}

func (r *MataProgramV2Repository) GetRowByUuid(ctx context.Context, rawUuid string) (*domain.MataProgramV2Row, error) {
	var row domain.MataProgramV2Row
	err := r.db.WithContext(ctx).
		Table("mata_programv2 mp").
		Select(flatSelectSQL).
		Joins("LEFT JOIN master_tahun mt ON mp.id_master_tahun = mt.id").
		Joins("LEFT JOIN v_fakultas_unit vfu ON mp.id_fakultas_unit = vfu.id").
		Where("mp.uuid = ? AND mp.deleted_at IS NULL", rawUuid).
		Take(&row).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.NotFound()
		}
		return nil, err
	}
	return &row, nil
}

func (r *MataProgramV2Repository) GetRecordByUuid(ctx context.Context, rawUuid string, includeDeleted bool) (*domain.MataProgramV2Record, error) {
	var rec domain.MataProgramV2Record
	q := r.db.WithContext(ctx).
		Table("mata_programv2").
		Select("id, uuid, mata_program, id_master_tahun, deleted_at").
		Where("uuid = ?", rawUuid)
	if !includeDeleted {
		q = q.Where("deleted_at IS NULL")
	}
	err := q.Take(&rec).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			if !includeDeleted {
				return nil, domain.NotFoundOrDeleted()
			}
			return nil, domain.NotFound()
		}
		return nil, err
	}
	return &rec, nil
}

func (r *MataProgramV2Repository) ResolveTahunID(ctx context.Context, tahunStr string) (uint, error) {
	var mt struct {
		ID uint
	}
	var err error
	if _, parseErr := uuid.Parse(tahunStr); parseErr == nil {
		err = r.db.WithContext(ctx).Table("master_tahun").Select("id").Where("uuid = ?", tahunStr).Take(&mt).Error
	} else if id, convErr := strconv.Atoi(tahunStr); convErr == nil && id > 0 {
		err = r.db.WithContext(ctx).Table("master_tahun").Select("id").Where("id = ? OR tahun = ?", id, tahunStr).Take(&mt).Error
	} else {
		err = r.db.WithContext(ctx).Table("master_tahun").Select("id").Where("tahun = ?", tahunStr).Take(&mt).Error
	}
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return 0, err
	}
	return mt.ID, nil
}

func (r *MataProgramV2Repository) ResolveOrCreateTahunID(ctx context.Context, tahunStr string) (uint, error) {
	id, err := r.ResolveTahunID(ctx, tahunStr)
	if err != nil {
		return 0, err
	}
	if id > 0 {
		return id, nil
	}

	newUUID := uuid.New().String()
	_ = r.db.WithContext(ctx).Exec("INSERT INTO master_tahun (uuid, tahun, status) VALUES (?, ?, 'aktif')", newUUID, tahunStr).Error

	var mt struct {
		ID uint
	}
	_ = r.db.WithContext(ctx).Table("master_tahun").Select("id").Where("tahun = ?", tahunStr).Take(&mt).Error
	return mt.ID, nil
}

func (r *MataProgramV2Repository) ResolveFakultasUnitID(ctx context.Context, rawTarget string) (uint, error) {
	if id, err := strconv.Atoi(rawTarget); err == nil && id > 0 {
		return uint(id), nil
	}
	if parsed, err := uuid.Parse(rawTarget); err == nil {
		var fu struct {
			ID uint
		}
		if err := r.db.WithContext(ctx).Table("v_fakultas_unit").Select("id").Where("uuid = ?", parsed).Take(&fu).Error; err == nil {
			return fu.ID, nil
		}
	}
	return 0, nil
}

func (r *MataProgramV2Repository) UpsertOrRestoreAssignment(ctx context.Context, tahunID uint, unitID uint, mataProgram string) (string, bool, error) {
	var existingAny struct {
		ID        uint
		UUID      string
		DeletedAt *time.Time
	}
	errFind := r.db.WithContext(ctx).Table("mata_programv2").
		Select("id, uuid, deleted_at").
		Where("mata_program = ? AND id_fakultas_unit = ? AND id_master_tahun = ?", mataProgram, unitID, tahunID).
		Order("id DESC").
		Take(&existingAny).Error

	if errFind == nil && existingAny.ID > 0 {
		if existingAny.DeletedAt != nil {
			_ = r.db.WithContext(ctx).Table("mata_programv2").
				Where("id = ?", existingAny.ID).
				Update("deleted_at", nil).Error
		}
		return existingAny.UUID, true, nil
	}

	newUUID := uuid.New().String()
	errInsert := r.db.WithContext(ctx).Exec(
		"INSERT INTO mata_programv2 (uuid, id_master_tahun, id_fakultas_unit, mata_program) VALUES (?, ?, ?, ?)",
		newUUID, tahunID, unitID, mataProgram,
	).Error
	if errInsert != nil {
		return "", false, errInsert
	}
	return newUUID, true, nil
}

func (r *MataProgramV2Repository) UpdateGroupFields(ctx context.Context, originalMataProgram string, originalTahunID *uint, updates map[string]interface{}) error {
	q := r.db.WithContext(ctx).Table("mata_programv2").
		Where("mata_program = ? AND deleted_at IS NULL", originalMataProgram)
	if originalTahunID != nil {
		q = q.Where("id_master_tahun = ?", *originalTahunID)
	}
	return q.Updates(updates).Error
}

func (r *MataProgramV2Repository) SoftDeleteGroupUnitsNotIn(ctx context.Context, mataProgram string, tahunID uint, keepUnitIDs []uint, deletedAt time.Time) error {
	return r.db.WithContext(ctx).Table("mata_programv2").
		Where("mata_program = ? AND id_master_tahun = ? AND id_fakultas_unit NOT IN (?) AND deleted_at IS NULL", mataProgram, tahunID, keepUnitIDs).
		Update("deleted_at", deletedAt).Error
}

func (r *MataProgramV2Repository) UpdateSingleByUuid(ctx context.Context, rawUuid string, updates map[string]interface{}) error {
	return r.db.WithContext(ctx).Table("mata_programv2").
		Where("uuid = ? AND deleted_at IS NULL", rawUuid).
		Updates(updates).Error
}

func (r *MataProgramV2Repository) SoftDeleteSingleByUuid(ctx context.Context, rawUuid string, deletedAt time.Time) (int64, error) {
	res := r.db.WithContext(ctx).Table("mata_programv2").
		Where("uuid = ? AND deleted_at IS NULL", rawUuid).
		Update("deleted_at", deletedAt)
	return res.RowsAffected, res.Error
}

func (r *MataProgramV2Repository) SoftDeleteGroup(ctx context.Context, mataProgram string, tahunID *uint, deletedAt time.Time) (int64, error) {
	q := r.db.WithContext(ctx).Table("mata_programv2").
		Where("mata_program = ? AND deleted_at IS NULL", mataProgram)
	if tahunID != nil {
		q = q.Where("id_master_tahun = ?", *tahunID)
	} else {
		q = q.Where("id_master_tahun IS NULL")
	}
	res := q.Update("deleted_at", deletedAt)
	return res.RowsAffected, res.Error
}

func (r *MataProgramV2Repository) RestoreSingleByUuid(ctx context.Context, rawUuid string) (int64, error) {
	res := r.db.WithContext(ctx).Table("mata_programv2").
		Where("uuid = ?", rawUuid).
		Update("deleted_at", nil)
	return res.RowsAffected, res.Error
}

func (r *MataProgramV2Repository) RestoreGroup(ctx context.Context, mataProgram string, tahunID *uint) (int64, error) {
	q := r.db.WithContext(ctx).Table("mata_programv2").
		Where("mata_program = ?", mataProgram)
	if tahunID != nil {
		q = q.Where("id_master_tahun = ?", *tahunID)
	} else {
		q = q.Where("id_master_tahun IS NULL")
	}
	res := q.Update("deleted_at", nil)
	return res.RowsAffected, res.Error
}
