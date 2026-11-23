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
		case "templaterenstrasasaran", "templaterenstra", "template_renstra", "indikator", "standar", "sasaran", "templaterenstraorunit":
			like := "%" + val + "%"
			q = q.Where(`EXISTS (
				SELECT 1 FROM mata_programv2_template_renstra mptr2
				JOIN template_renstra tr2 ON mptr2.id_template_renstra = tr2.id
				LEFT JOIN master_indikator_renstra mir2 ON tr2.indikator = mir2.id
				LEFT JOIN master_standar_renstra msr2 ON mir2.id_master_standar = msr2.id
				LEFT JOIN v_fakultas_unit vfu2 ON tr2.fakultas_unit = vfu2.id
				WHERE (mptr2.id_mata_programv2 = mp.id OR mptr2.id_mata_programv2 IN (
					SELECT mp_sub.id FROM mata_programv2 mp_sub 
					WHERE mp_sub.mata_program = mp.mata_program 
					  AND (mp_sub.id_master_tahun = mp.id_master_tahun OR (mp_sub.id_master_tahun IS NULL AND mp.id_master_tahun IS NULL))
				))
				  AND (
					LOWER(mir2.indikator) LIKE ? OR
					LOWER(msr2.nama) LIKE ? OR
					LOWER(vfu2.nama_fak_prod_unit) LIKE ? OR
					LOWER(COALESCE(vfu2.jenjang, '')) LIKE ? OR
					LOWER(CONCAT(vfu2.nama_fak_prod_unit, ' (', COALESCE(vfu2.jenjang, ''), ')')) LIKE ? OR
					LOWER(vfu2.fakultas) LIKE ?
				  )
			)`, like, like, like, like, like, like)
		case "info", "infounit", "info_unit", "prokerunit", "proker_unit", "matrixproker", "matrix_proker":
			like := "%" + val + "%"
			q = q.Where(`EXISTS (
				SELECT 1 FROM matrix_proker mpk
				JOIN mata_programv2 mp_src ON mpk.id_mata_program = mp_src.id
				JOIN v_fakultas_unit vfu_p ON mpk.id_fakultas_unit = vfu_p.id
				WHERE (mpk.id_mata_program = mp.id OR (
					mp_src.mata_program = mp.mata_program
					AND (mp_src.id_master_tahun = mp.id_master_tahun OR (mp_src.id_master_tahun IS NULL AND mp.id_master_tahun IS NULL))
				))
				  AND mpk.deleted_at IS NULL
				  AND (
					LOWER(vfu_p.nama_fak_prod_unit) LIKE ? OR
					LOWER(COALESCE(vfu_p.jenjang, '')) LIKE ? OR
					LOWER(CONCAT(vfu_p.nama_fak_prod_unit, ' (', COALESCE(vfu_p.jenjang, ''), ')')) LIKE ?
				  )
			)`, like, like, like)
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
		q = q.Where(`(
			LOWER(mp.mata_program) LIKE ? OR 
			LOWER(vfu.nama_fak_prod_unit) LIKE ? OR 
			LOWER(mt.tahun) LIKE ? OR 
			LOWER(vfu.fakultas) LIKE ? OR
			EXISTS (
				SELECT 1 FROM mata_programv2_template_renstra mptr2
				JOIN template_renstra tr2 ON mptr2.id_template_renstra = tr2.id
				LEFT JOIN master_indikator_renstra mir2 ON tr2.indikator = mir2.id
				LEFT JOIN master_standar_renstra msr2 ON mir2.id_master_standar = msr2.id
				LEFT JOIN v_fakultas_unit vfu2 ON tr2.fakultas_unit = vfu2.id
				WHERE (mptr2.id_mata_programv2 = mp.id OR mptr2.id_mata_programv2 IN (
					SELECT mp_sub.id FROM mata_programv2 mp_sub 
					WHERE mp_sub.mata_program = mp.mata_program 
					  AND (mp_sub.id_master_tahun = mp.id_master_tahun OR (mp_sub.id_master_tahun IS NULL AND mp.id_master_tahun IS NULL))
				))
				  AND (
					LOWER(mir2.indikator) LIKE ? OR
					LOWER(msr2.nama) LIKE ? OR
					LOWER(vfu2.nama_fak_prod_unit) LIKE ? OR
					LOWER(COALESCE(vfu2.jenjang, '')) LIKE ? OR
					LOWER(CONCAT(vfu2.nama_fak_prod_unit, ' (', COALESCE(vfu2.jenjang, ''), ')')) LIKE ? OR
					LOWER(vfu2.fakultas) LIKE ?
				  )
			) OR
			EXISTS (
				SELECT 1 FROM matrix_proker mpk
				JOIN mata_programv2 mp_src ON mpk.id_mata_program = mp_src.id
				JOIN v_fakultas_unit vfu_p ON mpk.id_fakultas_unit = vfu_p.id
				WHERE (mpk.id_mata_program = mp.id OR (
					mp_src.mata_program = mp.mata_program
					AND (mp_src.id_master_tahun = mp.id_master_tahun OR (mp_src.id_master_tahun IS NULL AND mp.id_master_tahun IS NULL))
				))
				  AND mpk.deleted_at IS NULL
				  AND (
					LOWER(vfu_p.nama_fak_prod_unit) LIKE ? OR
					LOWER(COALESCE(vfu_p.jenjang, '')) LIKE ? OR
					LOWER(CONCAT(vfu_p.nama_fak_prod_unit, ' (', COALESCE(vfu_p.jenjang, ''), ')')) LIKE ?
				  )
			)
		)`, like, like, like, like, like, like, like, like, like, like, like, like, like)
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

func (r *MataProgramV2Repository) SaveMataProgramV2WithTemplates(ctx context.Context, tahunID uint, mataProgram string, templateRenstraIDs []uint) (string, uint, error) {
	var mp domain.MataProgramV2Record
	errFind := r.db.WithContext(ctx).Table("mata_programv2").
		Select("id, uuid, deleted_at").
		Where("mata_program = ? AND id_master_tahun = ?", mataProgram, tahunID).
		Order("id DESC").
		Take(&mp).Error

	var mpID uint
	var mpUUID string

	if errFind == nil && mp.ID > 0 {
		mpID = mp.ID
		mpUUID = mp.UUID
		if mp.DeletedAt != nil {
			_ = r.db.WithContext(ctx).Table("mata_programv2").
				Where("id = ?", mp.ID).
				Update("deleted_at", nil).Error
		}
	} else {
		newUUID := uuid.New().String()
		errInsert := r.db.WithContext(ctx).Exec(
			"INSERT INTO mata_programv2 (uuid, id_master_tahun, mata_program) VALUES (?, ?, ?)",
			newUUID, tahunID, mataProgram,
		).Error
		if errInsert != nil {
			return "", 0, errInsert
		}
		var created domain.MataProgramV2Record
		if err := r.db.WithContext(ctx).Table("mata_programv2").
			Select("id, uuid").
			Where("uuid = ?", newUUID).
			Take(&created).Error; err != nil {
			return "", 0, err
		}
		mpID = created.ID
		mpUUID = created.UUID
	}

	// Hapus relasi lama lalu insert yang baru
	_ = r.db.WithContext(ctx).Exec("DELETE FROM mata_programv2_template_renstra WHERE id_mata_programv2 = ?", mpID).Error

	var insertedCount uint = 0
	for _, trID := range templateRenstraIDs {
		if trID == 0 {
			continue
		}
		linkUUID := uuid.New().String()
		errLink := r.db.WithContext(ctx).Exec(
			"INSERT INTO mata_programv2_template_renstra (uuid, id_mata_programv2, id_template_renstra) VALUES (?, ?, ?) ON DUPLICATE KEY UPDATE updated_at = CURRENT_TIMESTAMP",
			linkUUID, mpID, trID,
		).Error
		if errLink == nil {
			insertedCount++
		}
	}

	return mpUUID, insertedCount, nil
}

func (r *MataProgramV2Repository) UpdateMataProgramV2WithTemplates(ctx context.Context, rawUuid string, tahunID uint, mataProgram string, templateRenstraIDs []uint) error {
	var mp domain.MataProgramV2Record
	if err := r.db.WithContext(ctx).Table("mata_programv2").
		Select("id, uuid").
		Where("uuid = ?", rawUuid).
		Take(&mp).Error; err != nil {
		return err
	}

	updates := map[string]interface{}{}
	if mataProgram != "" {
		updates["mata_program"] = mataProgram
	}
	if tahunID > 0 {
		updates["id_master_tahun"] = tahunID
	}
	if len(updates) > 0 {
		if err := r.db.WithContext(ctx).Table("mata_programv2").
			Where("id = ?", mp.ID).
			Updates(updates).Error; err != nil {
			return err
		}
	}

	// Update relasi template renstra
	_ = r.db.WithContext(ctx).Exec("DELETE FROM mata_programv2_template_renstra WHERE id_mata_programv2 = ?", mp.ID).Error

	for _, trID := range templateRenstraIDs {
		if trID == 0 {
			continue
		}
		linkUUID := uuid.New().String()
		_ = r.db.WithContext(ctx).Exec(
			"INSERT INTO mata_programv2_template_renstra (uuid, id_mata_programv2, id_template_renstra) VALUES (?, ?, ?) ON DUPLICATE KEY UPDATE updated_at = CURRENT_TIMESTAMP",
			linkUUID, mp.ID, trID,
		).Error
	}

	return nil
}

func (r *MataProgramV2Repository) GetTemplatesByMataProgramID(ctx context.Context, mpID uint) ([]domain.MataProgramV2TemplateItem, []uint, error) {
	var items []domain.MataProgramV2TemplateItem
	var ids []uint

	sql := `
		SELECT 
			tr.id as id,
			tr.uuid as uuid,
			tr.id as id_template_renstra,
			tr.indikator as indikator_id,
			COALESCE(mir.indikator, '') as indikator,
			COALESCE(msr.nama, '') as standar_renstra,
			tr.fakultas_unit as fakultas_unit_id,
			CASE
				WHEN vfu.type = 'prodi' AND vfu.jenjang IS NOT NULL AND vfu.jenjang != ''
				THEN CONCAT(vfu.nama_fak_prod_unit, ' (', UPPER(vfu.jenjang), ')')
				ELSE COALESCE(vfu.nama_fak_prod_unit, '')
			END as fakultas_unit,
			COALESCE(vfu.jenjang, '') as jenjang,
			COALESCE(vfu.type, '') as type,
			COALESCE(tr.kategori, '') as kategori,
			COALESCE(tr.target, '') as target,
			COALESCE(tr.satuan, '') as satuan
		FROM mata_programv2_template_renstra mptr
		JOIN template_renstra tr ON mptr.id_template_renstra = tr.id
		LEFT JOIN master_indikator_renstra mir ON tr.indikator = mir.id
		LEFT JOIN master_standar_renstra msr ON mir.id_master_standar = msr.id
		LEFT JOIN v_fakultas_unit vfu ON tr.fakultas_unit = vfu.id
		WHERE mptr.id_mata_programv2 = ?
		ORDER BY msr.nama ASC, mir.indikator ASC, fakultas_unit ASC
	`
	err := r.db.WithContext(ctx).Raw(sql, mpID).Scan(&items).Error
	if err != nil {
		return nil, nil, err
	}

	for _, it := range items {
		ids = append(ids, it.ID)
	}

	return items, ids, nil
}

func (r *MataProgramV2Repository) GetTemplatesByMataProgramIDs(ctx context.Context, mpIDs []uint) (map[uint][]domain.MataProgramV2TemplateItem, error) {
	result := make(map[uint][]domain.MataProgramV2TemplateItem)
	if len(mpIDs) == 0 {
		return result, nil
	}

	type QueryRow struct {
		IdMataProgramV2    uint   `gorm:"column:id_mata_programv2"`
		ID                 uint   `gorm:"column:id"`
		UUID               string `gorm:"column:uuid"`
		IdTemplateRenstra  uint   `gorm:"column:id_template_renstra"`
		IndikatorRenstraID uint   `gorm:"column:indikator_id"`
		Indikator          string `gorm:"column:indikator"`
		StandarRenstra     string `gorm:"column:standar_renstra"`
		FakultasUnitID     uint   `gorm:"column:fakultas_unit_id"`
		FakultasUnit       string `gorm:"column:fakultas_unit"`
		Jenjang            string `gorm:"column:jenjang"`
		Type               string `gorm:"column:type"`
		Kategori           string `gorm:"column:kategori"`
		Target             string `gorm:"column:target"`
		Satuan             string `gorm:"column:satuan"`
	}

	var rows []QueryRow
	sql := `
		SELECT 
			COALESCE(mp_target.id, mptr.id_mata_programv2) as id_mata_programv2,
			tr.id as id,
			tr.uuid as uuid,
			tr.id as id_template_renstra,
			tr.indikator as indikator_id,
			COALESCE(mir.indikator, '') as indikator,
			COALESCE(msr.nama, '') as standar_renstra,
			tr.fakultas_unit as fakultas_unit_id,
			CASE
				WHEN vfu.type = 'prodi' AND vfu.jenjang IS NOT NULL AND vfu.jenjang != ''
				THEN CONCAT(vfu.nama_fak_prod_unit, ' (', UPPER(vfu.jenjang), ')')
				ELSE COALESCE(vfu.nama_fak_prod_unit, '')
			END as fakultas_unit,
			COALESCE(vfu.jenjang, '') as jenjang,
			COALESCE(vfu.type, '') as type,
			COALESCE(tr.kategori, '') as kategori,
			COALESCE(tr.target, '') as target,
			COALESCE(tr.satuan, '') as satuan
		FROM mata_programv2_template_renstra mptr
		JOIN template_renstra tr ON mptr.id_template_renstra = tr.id
		JOIN mata_programv2 mp_src ON mptr.id_mata_programv2 = mp_src.id
		JOIN mata_programv2 mp_target ON mp_src.mata_program = mp_target.mata_program 
		  AND (mp_src.id_master_tahun = mp_target.id_master_tahun OR (mp_src.id_master_tahun IS NULL AND mp_target.id_master_tahun IS NULL))
		LEFT JOIN master_indikator_renstra mir ON tr.indikator = mir.id
		LEFT JOIN master_standar_renstra msr ON mir.id_master_standar = msr.id
		LEFT JOIN v_fakultas_unit vfu ON tr.fakultas_unit = vfu.id
		WHERE mp_target.id IN (?)
		GROUP BY mp_target.id, tr.id
		ORDER BY msr.nama ASC, mir.indikator ASC, fakultas_unit ASC
	`
	err := r.db.WithContext(ctx).Raw(sql, mpIDs).Scan(&rows).Error
	if err != nil {
		return result, err
	}

	for _, rw := range rows {
		item := domain.MataProgramV2TemplateItem{
			ID:                 rw.ID,
			UUID:               rw.UUID,
			IdTemplateRenstra:  rw.IdTemplateRenstra,
			IndikatorRenstraID: rw.IndikatorRenstraID,
			Indikator:          rw.Indikator,
			StandarRenstra:     rw.StandarRenstra,
			FakultasUnitID:     rw.FakultasUnitID,
			FakultasUnit:       rw.FakultasUnit,
			Jenjang:            rw.Jenjang,
			Type:               rw.Type,
			Kategori:           rw.Kategori,
			Target:             rw.Target,
			Satuan:             rw.Satuan,
		}
		result[rw.IdMataProgramV2] = append(result[rw.IdMataProgramV2], item)
	}

	return result, nil
}

func (r *MataProgramV2Repository) GetProkerUnitsByMataProgramIDs(ctx context.Context, mpIDs []uint) (map[uint][]domain.MataProgramV2ProkerUnit, error) {
	result := make(map[uint][]domain.MataProgramV2ProkerUnit)
	if len(mpIDs) == 0 {
		return result, nil
	}

	type QueryRow struct {
		IdMataProgram    uint   `gorm:"column:id_mata_program"`
		Id               uint   `gorm:"column:id"`
		UUID             string `gorm:"column:uuid"`
		FakultasUnitID   uint   `gorm:"column:id_fakultas_unit"`
		FakultasUnitUuid string `gorm:"column:fakultas_unit_uuid"`
		NamaFakultasUnit string `gorm:"column:nama_fakultas_unit"`
		Label            string `gorm:"column:label"`
		Jenjang          string `gorm:"column:jenjang"`
		Type             string `gorm:"column:type"`
		Status           string `gorm:"column:status"`
	}

	var rows []QueryRow
	sql := `
		SELECT 
			COALESCE(mp_target.id, mpk.id_mata_program) as id_mata_program,
			MIN(mpk.id) as id,
			MIN(mpk.uuid) as uuid,
			mpk.id_fakultas_unit as id_fakultas_unit,
			COALESCE(vfu.uuid, '') as fakultas_unit_uuid,
			COALESCE(vfu.nama_fak_prod_unit, '') as nama_fakultas_unit,
			CASE 
				WHEN vfu.type = 'prodi' AND vfu.jenjang IS NOT NULL AND vfu.jenjang != '' 
				THEN CONCAT(vfu.nama_fak_prod_unit, ' (', UPPER(vfu.jenjang), ')') 
				ELSE COALESCE(vfu.nama_fak_prod_unit, '') 
			END as label,
			COALESCE(vfu.jenjang, '') as jenjang,
			COALESCE(vfu.type, '') as type,
			COALESCE(MIN(mpk.status), 'draft') as status
		FROM matrix_proker mpk
		JOIN mata_programv2 mp_src ON mpk.id_mata_program = mp_src.id
		JOIN mata_programv2 mp_target ON mp_src.mata_program = mp_target.mata_program 
		  AND (mp_src.id_master_tahun = mp_target.id_master_tahun OR (mp_src.id_master_tahun IS NULL AND mp_target.id_master_tahun IS NULL))
		JOIN v_fakultas_unit vfu ON mpk.id_fakultas_unit = vfu.id
		WHERE mp_target.id IN (?)
		  AND mpk.deleted_at IS NULL
		GROUP BY mp_target.id, mpk.id_fakultas_unit, vfu.uuid, vfu.nama_fak_prod_unit, vfu.jenjang, vfu.type
		ORDER BY label ASC
	`
	err := r.db.WithContext(ctx).Raw(sql, mpIDs).Scan(&rows).Error
	if err != nil {
		return result, err
	}

	for _, rw := range rows {
		item := domain.MataProgramV2ProkerUnit{
			Id:               rw.Id,
			UUID:             rw.UUID,
			IdMataProgram:    rw.IdMataProgram,
			FakultasUnitID:   rw.FakultasUnitID,
			FakultasUnitUuid: rw.FakultasUnitUuid,
			NamaFakultasUnit: rw.NamaFakultasUnit,
			Label:            rw.Label,
			Jenjang:          rw.Jenjang,
			Type:             rw.Type,
			Status:           rw.Status,
		}
		result[rw.IdMataProgram] = append(result[rw.IdMataProgram], item)
	}

	return result, nil
}

func (r *MataProgramV2Repository) GetProkerUnitsByMataProgramID(ctx context.Context, mpID uint) ([]domain.MataProgramV2ProkerUnit, error) {
	mpMap, err := r.GetProkerUnitsByMataProgramIDs(ctx, []uint{mpID})
	if err != nil {
		return nil, err
	}
	return mpMap[mpID], nil
}
