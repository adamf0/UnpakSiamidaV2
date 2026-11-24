package infrastructure

import (
	domain "UnpakSiamida/modules/mataprogramv2/domain"
	"context"
	"errors"
	"fmt"
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
	COALESCE((SELECT COUNT(*) FROM matrix_proker mpk WHERE mpk.id_mata_program = mp.id AND mpk.deleted_at IS NULL), 0) AS proker_count,
	mp.created_at,
	mp.updated_at,
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
			q = q.Where(`EXISTS (
				SELECT 1 FROM mata_programv2_template_renstra mptr_u
				JOIN template_renstra tr_u ON mptr_u.id_template_renstra = tr_u.id
				JOIN v_fakultas_unit vfu_u ON tr_u.fakultas_unit = vfu_u.id
				WHERE mptr_u.id_mata_programv2 = mp.id
				  AND (mt.tahun IS NULL OR mt.tahun = '' OR CAST(tr_u.tahun AS CHAR) = CAST(mt.tahun AS CHAR))
				  AND (
					LOWER(vfu_u.nama_fak_prod_unit) LIKE ? OR
					LOWER(COALESCE(vfu_u.jenjang, '')) LIKE ? OR
					LOWER(CONCAT(vfu_u.nama_fak_prod_unit, ' (', COALESCE(vfu_u.jenjang, ''), ')')) LIKE ? OR
					LOWER(COALESCE(vfu_u.fakultas, '')) LIKE ?
				  )
			)`, like, like, like, like)
		case "templaterenstrasasaran", "templaterenstra", "template_renstra", "indikator", "standar", "sasaran", "templaterenstraorunit":
			like := "%" + val + "%"
			q = q.Where(`EXISTS (
				SELECT 1 FROM mata_programv2_template_renstra mptr2
				JOIN template_renstra tr2 ON mptr2.id_template_renstra = tr2.id
				LEFT JOIN master_indikator_renstra mir2 ON tr2.indikator = mir2.id
				LEFT JOIN master_standar_renstra msr2 ON mir2.id_master_standar = msr2.id
				LEFT JOIN v_fakultas_unit vfu2 ON tr2.fakultas_unit = vfu2.id
				WHERE mptr2.id_mata_programv2 = mp.id
				  AND (mt.tahun IS NULL OR mt.tahun = '' OR CAST(tr2.tahun AS CHAR) = CAST(mt.tahun AS CHAR))
				  AND (
					LOWER(mir2.indikator) LIKE ? OR
					LOWER(msr2.nama) LIKE ? OR
					LOWER(vfu2.nama_fak_prod_unit) LIKE ? OR
					LOWER(COALESCE(vfu2.jenjang, '')) LIKE ? OR
					LOWER(CONCAT(vfu2.nama_fak_prod_unit, ' (', COALESCE(vfu2.jenjang, ''), ')')) LIKE ? OR
					LOWER(COALESCE(vfu2.fakultas, '')) LIKE ?
				  )
			)`, like, like, like, like, like, like)
		case "informasi", "info", "infounit", "info_unit", "prokerunit", "proker_unit", "matrixproker", "matrix_proker":
			like := "%" + val + "%"
			q = q.Where(`EXISTS (
				SELECT 1 FROM matrix_proker mpk
				JOIN v_fakultas_unit vfu_p ON mpk.id_fakultas_unit = vfu_p.id
				WHERE mpk.id_mata_program = mp.id
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
			q = q.Where(`EXISTS (
				SELECT 1 FROM mata_programv2_template_renstra mptr_u
				JOIN template_renstra tr_u ON mptr_u.id_template_renstra = tr_u.id
				JOIN v_fakultas_unit vfu_u ON tr_u.fakultas_unit = vfu_u.id
				WHERE mptr_u.id_mata_programv2 = mp.id
				  AND (mt.tahun IS NULL OR mt.tahun = '' OR CAST(tr_u.tahun AS CHAR) = CAST(mt.tahun AS CHAR))
				  AND vfu_u.uuid = ?
			)`, filter.Unit)
		} else if id, err := strconv.Atoi(filter.Unit); err == nil && id > 0 {
			q = q.Where(`EXISTS (
				SELECT 1 FROM mata_programv2_template_renstra mptr_u
				JOIN template_renstra tr_u ON mptr_u.id_template_renstra = tr_u.id
				WHERE mptr_u.id_mata_programv2 = mp.id
				  AND (mt.tahun IS NULL OR mt.tahun = '' OR CAST(tr_u.tahun AS CHAR) = CAST(mt.tahun AS CHAR))
				  AND tr_u.fakultas_unit = ?
			)`, id)
		}
	}

	if filter.Search != "" {
		like := "%" + strings.ToLower(filter.Search) + "%"
		q = q.Where(`(
			LOWER(mp.mata_program) LIKE ? OR
			LOWER(COALESCE(mt.tahun, '')) LIKE ? OR
			EXISTS (
				SELECT 1 FROM mata_programv2_template_renstra mptr2
				JOIN template_renstra tr2 ON mptr2.id_template_renstra = tr2.id
				LEFT JOIN master_indikator_renstra mir2 ON tr2.indikator = mir2.id
				LEFT JOIN master_standar_renstra msr2 ON mir2.id_master_standar = msr2.id
				LEFT JOIN v_fakultas_unit vfu2 ON tr2.fakultas_unit = vfu2.id
				WHERE mptr2.id_mata_programv2 = mp.id
				  AND (mt.tahun IS NULL OR mt.tahun = '' OR CAST(tr2.tahun AS CHAR) = CAST(mt.tahun AS CHAR))
				  AND (
					LOWER(mir2.indikator) LIKE ? OR
					LOWER(msr2.nama) LIKE ? OR
					LOWER(vfu2.nama_fak_prod_unit) LIKE ? OR
					LOWER(COALESCE(vfu2.jenjang, '')) LIKE ? OR
					LOWER(CONCAT(vfu2.nama_fak_prod_unit, ' (', COALESCE(vfu2.jenjang, ''), ')')) LIKE ? OR
					LOWER(COALESCE(vfu2.fakultas, '')) LIKE ?
				  )
			) OR
			EXISTS (
				SELECT 1 FROM matrix_proker mpk
				JOIN v_fakultas_unit vfu_p ON mpk.id_fakultas_unit = vfu_p.id
				WHERE mpk.id_mata_program = mp.id
				  AND mpk.deleted_at IS NULL
				  AND (
					LOWER(vfu_p.nama_fak_prod_unit) LIKE ? OR
					LOWER(COALESCE(vfu_p.jenjang, '')) LIKE ? OR
					LOWER(CONCAT(vfu_p.nama_fak_prod_unit, ' (', COALESCE(vfu_p.jenjang, ''), ')')) LIKE ?
				  )
			)
		)`, like, like, like, like, like, like, like, like, like, like, like)
	}

	return q
}

func (r *MataProgramV2Repository) baseJoinedQuery(ctx context.Context, filter domain.MataProgramV2Filter) *gorm.DB {
	q := r.db.WithContext(ctx).
		Table("mata_programv2 mp").
		Joins("LEFT JOIN master_tahun mt ON mp.id_master_tahun = mt.id")
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

	// Jika difilter untuk satu unit tertentu, isi metadata unit dari v_fakultas_unit
	// agar dropdown konsumen (seperti Risk Assessment / Proker V3) mendapat info unitnya.
	if len(list) > 0 && filter.Unit != "" && filter.Unit != "all" {
		var fu struct {
			ID       uint   `gorm:"column:id"`
			UUID     string `gorm:"column:uuid"`
			Label    string `gorm:"column:fakultas_unit_label"`
			Jenjang  string `gorm:"column:jenjang"`
			Type     string `gorm:"column:type"`
			Fakultas string `gorm:"column:fakultas"`
		}
		fuQ := r.db.WithContext(ctx).Table("v_fakultas_unit vfu").
			Select(`
				vfu.id,
				COALESCE(vfu.uuid, '') as uuid,
				CASE
					WHEN vfu.type = 'prodi' AND vfu.jenjang IS NOT NULL AND vfu.jenjang != ''
					THEN CONCAT(vfu.nama_fak_prod_unit, ' (', UPPER(vfu.jenjang), ')')
					ELSE COALESCE(vfu.nama_fak_prod_unit, '')
				END as fakultas_unit_label,
				COALESCE(vfu.jenjang, '') as jenjang,
				COALESCE(vfu.type, '') as type,
				COALESCE(vfu.fakultas, '') as fakultas
			`)
		if _, parseErr := uuid.Parse(filter.Unit); parseErr == nil {
			fuQ = fuQ.Where("vfu.uuid = ?", filter.Unit)
		} else if id, convErr := strconv.Atoi(filter.Unit); convErr == nil && id > 0 {
			fuQ = fuQ.Where("vfu.id = ?", id)
		}
		if errFu := fuQ.Take(&fu).Error; errFu == nil && fu.ID > 0 {
			for i := range list {
				uidVal := fu.ID
				list[i].FakultasUnitId = &uidVal
				list[i].FakultasUnitUuid = fu.UUID
				list[i].FakultasUnit = fu.Label
				list[i].Jenjang = fu.Jenjang
				list[i].Type = fu.Type
				list[i].Fakultas = fu.Fakultas
			}
		}
	}

	return list, nil
}

func (r *MataProgramV2Repository) Count(ctx context.Context, filter domain.MataProgramV2Filter) (int64, error) {
	var total int64
	err := r.baseJoinedQuery(ctx, filter).Count(&total).Error
	return total, err
}

func (r *MataProgramV2Repository) GetPaged(ctx context.Context, filter domain.MataProgramV2Filter, offset, limit int) ([]domain.GroupSummary, error) {
	var summaries []domain.GroupSummary
	err := r.baseJoinedQuery(ctx, filter).
		Select(`
			mp.id,
			mp.uuid,
			mp.id_master_tahun,
			mp.mata_program,
			COALESCE(mt.uuid, '') as tahun_uuid,
			COALESCE(mt.tahun, '') as tahun,
			mp.created_at,
			mp.updated_at,
			mp.deleted_at
		`).
		Order("mp.id DESC").
		Offset(offset).
		Limit(limit).
		Scan(&summaries).Error
	if err != nil {
		return nil, err
	}
	return summaries, nil
}

func (r *MataProgramV2Repository) GetRowByUuid(ctx context.Context, rawUuid string) (*domain.MataProgramV2Row, error) {
	var row domain.MataProgramV2Row
	err := r.db.WithContext(ctx).
		Table("mata_programv2 mp").
		Select(flatSelectSQL).
		Joins("LEFT JOIN master_tahun mt ON mp.id_master_tahun = mt.id").
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
	if err := r.db.WithContext(ctx).Exec("INSERT INTO master_tahun (uuid, tahun, status) VALUES (?, ?, 'aktif')", newUUID, tahunStr).Error; err != nil {
		return 0, err
	}

	var mt struct {
		ID uint
	}
	if err := r.db.WithContext(ctx).Table("master_tahun").Select("id").Where("tahun = ?", tahunStr).Take(&mt).Error; err != nil {
		return 0, err
	}
	return mt.ID, nil
}

func (r *MataProgramV2Repository) filterValidTemplateIDsForYear(tx *gorm.DB, tahunID uint, templateRenstraIDs []uint) ([]uint, error) {
	seen := make(map[uint]bool, len(templateRenstraIDs))
	candidateIDs := make([]uint, 0, len(templateRenstraIDs))
	for _, id := range templateRenstraIDs {
		if id > 0 && !seen[id] {
			seen[id] = true
			candidateIDs = append(candidateIDs, id)
		}
	}
	if len(candidateIDs) == 0 {
		return nil, nil
	}

	var tahunStr string
	if tahunID > 0 {
		_ = tx.Table("master_tahun").Select("COALESCE(tahun, '')").Where("id = ?", tahunID).Take(&tahunStr).Error
		tahunStr = strings.TrimSpace(tahunStr)
	}

	type templateCandidateRow struct {
		ID           uint `gorm:"column:id"`
		Indikator    uint `gorm:"column:indikator"`
		FakultasUnit uint `gorm:"column:fakultas_unit"`
	}
	type comboKey struct {
		Indikator    uint
		FakultasUnit uint
	}

	// Deduplikasi berdasarkan (indikator, fakultas_unit) dengan mengambil ID template_renstra terbaru (MAX id)
	bestByCombo := make(map[comboKey]uint, len(candidateIDs))
	const checkChunk = 1000
	for i := 0; i < len(candidateIDs); i += checkChunk {
		end := i + checkChunk
		if end > len(candidateIDs) {
			end = len(candidateIDs)
		}
		var rows []templateCandidateRow
		q := tx.Table("template_renstra").
			Select("id, indikator, fakultas_unit").
			Where("id IN ?", candidateIDs[i:end])
		if tahunStr != "" {
			q = q.Where("CAST(tahun AS CHAR) = ?", tahunStr)
		}
		if err := q.Scan(&rows).Error; err != nil {
			return nil, err
		}
		for _, rw := range rows {
			if rw.ID == 0 || rw.FakultasUnit == 0 {
				continue
			}
			key := comboKey{Indikator: rw.Indikator, FakultasUnit: rw.FakultasUnit}
			if prev, exists := bestByCombo[key]; !exists || rw.ID > prev {
				bestByCombo[key] = rw.ID
			}
		}
	}

	validIDs := make([]uint, 0, len(bestByCombo))
	for _, id := range bestByCombo {
		validIDs = append(validIDs, id)
	}
	return validIDs, nil
}

func insertTemplateRenstraMappings(tx *gorm.DB, mpID uint, validIDs []uint) (uint, error) {
	if err := tx.Exec("DELETE FROM mata_programv2_template_renstra WHERE id_mata_programv2 = ?", mpID).Error; err != nil {
		return 0, err
	}

	var insertedCount uint
	const chunkSize = 250
	for i := 0; i < len(validIDs); i += chunkSize {
		end := i + chunkSize
		if end > len(validIDs) {
			end = len(validIDs)
		}
		chunk := validIDs[i:end]

		placeholders := make([]string, 0, len(chunk))
		args := make([]interface{}, 0, len(chunk)*3)
		for _, trID := range chunk {
			placeholders = append(placeholders, "(?, ?, ?)")
			args = append(args, uuid.New().String(), mpID, trID)
		}

		query := fmt.Sprintf(
			"INSERT INTO mata_programv2_template_renstra (uuid, id_mata_programv2, id_template_renstra) VALUES %s ON DUPLICATE KEY UPDATE updated_at = CURRENT_TIMESTAMP",
			strings.Join(placeholders, ", "),
		)
		if err := tx.Exec(query, args...).Error; err != nil {
			return insertedCount, err
		}
		insertedCount += uint(len(chunk))
	}
	return insertedCount, nil
}

func (r *MataProgramV2Repository) SaveMataProgramV2WithTemplates(ctx context.Context, tahunID uint, mataProgram string, templateRenstraIDs []uint) (string, uint, error) {
	var mpUUID string
	var insertedCount uint

	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		validIDs, err := r.filterValidTemplateIDsForYear(tx, tahunID, templateRenstraIDs)
		if err != nil {
			return err
		}
		if len(validIDs) == 0 {
			return domain.InvalidTemplates()
		}

		// Tolak duplikasi aktif untuk kombinasi (mata_program, id_master_tahun) yang sama
		var activeCount int64
		if err := tx.Table("mata_programv2").
			Where("LOWER(TRIM(mata_program)) = LOWER(TRIM(?)) AND id_master_tahun = ? AND deleted_at IS NULL", mataProgram, tahunID).
			Count(&activeCount).Error; err != nil {
			return err
		}
		if activeCount > 0 {
			return domain.DuplicateMataProgram()
		}

		var mpID uint
		var deletedRec domain.MataProgramV2Record
		errFindDeleted := tx.Table("mata_programv2").
			Select("id, uuid").
			Where("LOWER(TRIM(mata_program)) = LOWER(TRIM(?)) AND id_master_tahun = ? AND deleted_at IS NOT NULL", mataProgram, tahunID).
			Order("id DESC").
			Take(&deletedRec).Error

		if errFindDeleted == nil && deletedRec.ID > 0 {
			mpID = deletedRec.ID
			mpUUID = deletedRec.UUID
			if err := tx.Table("mata_programv2").
				Where("id = ?", mpID).
				Updates(map[string]interface{}{
					"mata_program": mataProgram,
					"deleted_at":   nil,
				}).Error; err != nil {
				return err
			}
		} else {
			newUUID := uuid.New().String()
			if err := tx.Exec(
				"INSERT INTO mata_programv2 (uuid, id_master_tahun, mata_program) VALUES (?, ?, ?)",
				newUUID, tahunID, mataProgram,
			).Error; err != nil {
				return err
			}
			var created domain.MataProgramV2Record
			if err := tx.Table("mata_programv2").
				Select("id, uuid").
				Where("uuid = ?", newUUID).
				Take(&created).Error; err != nil {
				return err
			}
			mpID = created.ID
			mpUUID = created.UUID
		}

		cnt, err := insertTemplateRenstraMappings(tx, mpID, validIDs)
		if err != nil {
			return err
		}
		insertedCount = cnt
		return nil
	})
	if err != nil {
		return "", 0, err
	}

	return mpUUID, insertedCount, nil
}

func (r *MataProgramV2Repository) UpdateMataProgramV2WithTemplates(ctx context.Context, rawUuid string, tahunID uint, mataProgram string, templateRenstraIDs []uint, replaceTemplates bool) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var mp domain.MataProgramV2Record
		if err := tx.Table("mata_programv2").
			Select("id, uuid, mata_program, id_master_tahun").
			Where("uuid = ? AND deleted_at IS NULL", rawUuid).
			Take(&mp).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return domain.NotFoundOrDeleted()
			}
			return err
		}

		targetNama := mp.MataProgram
		if mataProgram != "" {
			targetNama = mataProgram
		}
		var targetTahunID uint
		if mp.IdMasterTahun != nil {
			targetTahunID = *mp.IdMasterTahun
		}
		if tahunID > 0 {
			targetTahunID = tahunID
		}

		// Pastikan tidak bentrok dengan mata program aktif lain di tahun yang sama
		var dupCount int64
		dupQ := tx.Table("mata_programv2").
			Where("id != ? AND LOWER(TRIM(mata_program)) = LOWER(TRIM(?)) AND deleted_at IS NULL", mp.ID, targetNama)
		if targetTahunID > 0 {
			dupQ = dupQ.Where("id_master_tahun = ?", targetTahunID)
		} else {
			dupQ = dupQ.Where("id_master_tahun IS NULL")
		}
		if err := dupQ.Count(&dupCount).Error; err != nil {
			return err
		}
		if dupCount > 0 {
			return domain.DuplicateMataProgram()
		}

		var validIDs []uint
		if replaceTemplates {
			var err error
			validIDs, err = r.filterValidTemplateIDsForYear(tx, targetTahunID, templateRenstraIDs)
			if err != nil {
				return err
			}
			if len(validIDs) == 0 {
				return domain.InvalidTemplates()
			}
		}

		updates := map[string]interface{}{
			"mata_program": targetNama,
		}
		if targetTahunID > 0 {
			updates["id_master_tahun"] = targetTahunID
		}
		if err := tx.Table("mata_programv2").
			Where("id = ?", mp.ID).
			Updates(updates).Error; err != nil {
			return err
		}

		if replaceTemplates {
			if _, err := insertTemplateRenstraMappings(tx, mp.ID, validIDs); err != nil {
				return err
			}
		}

		return nil
	})
}

func (r *MataProgramV2Repository) SoftDeleteByUuid(ctx context.Context, rawUuid string, deletedAt time.Time) (int64, error) {
	res := r.db.WithContext(ctx).Table("mata_programv2").
		Where("uuid = ? AND deleted_at IS NULL", rawUuid).
		Update("deleted_at", deletedAt)
	return res.RowsAffected, res.Error
}

func (r *MataProgramV2Repository) RestoreByUuid(ctx context.Context, rawUuid string) (int64, error) {
	var rec domain.MataProgramV2Record
	if err := r.db.WithContext(ctx).Table("mata_programv2").
		Select("id, mata_program, id_master_tahun").
		Where("uuid = ?", rawUuid).
		Take(&rec).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, domain.NotFound()
		}
		return 0, err
	}

	var dupCount int64
	dupQ := r.db.WithContext(ctx).Table("mata_programv2").
		Where("id != ? AND LOWER(TRIM(mata_program)) = LOWER(TRIM(?)) AND deleted_at IS NULL", rec.ID, rec.MataProgram)
	if rec.IdMasterTahun != nil {
		dupQ = dupQ.Where("id_master_tahun = ?", *rec.IdMasterTahun)
	} else {
		dupQ = dupQ.Where("id_master_tahun IS NULL")
	}
	if err := dupQ.Count(&dupCount).Error; err != nil {
		return 0, err
	}
	if dupCount > 0 {
		return 0, domain.DuplicateMataProgram()
	}

	res := r.db.WithContext(ctx).Table("mata_programv2").
		Where("id = ?", rec.ID).
		Update("deleted_at", nil)
	return res.RowsAffected, res.Error
}

func (r *MataProgramV2Repository) GetTargetUnitsByMataProgramIDs(ctx context.Context, mpIDs []uint) (map[uint][]domain.MataProgramV2UnitBadge, error) {
	result := make(map[uint][]domain.MataProgramV2UnitBadge)
	if len(mpIDs) == 0 {
		return result, nil
	}

	type QueryRow struct {
		IdMataProgramV2  uint   `gorm:"column:id_mata_programv2"`
		FakultasUnitId   *uint  `gorm:"column:id_fakultas_unit"`
		FakultasUnitUuid string `gorm:"column:fakultas_unit_uuid"`
		FakultasUnit     string `gorm:"column:fakultas_unit_label"`
		Jenjang          string `gorm:"column:jenjang"`
		Type             string `gorm:"column:type"`
		Fakultas         string `gorm:"column:fakultas"`
		ProkerCount      int64  `gorm:"column:proker_count"`
	}

	var rows []QueryRow
	sql := `
		SELECT
			mptr.id_mata_programv2 AS id_mata_programv2,
			tr.fakultas_unit AS id_fakultas_unit,
			COALESCE(vfu.uuid, '') AS fakultas_unit_uuid,
			CASE
				WHEN vfu.type = 'prodi' AND vfu.jenjang IS NOT NULL AND vfu.jenjang != ''
				THEN CONCAT(vfu.nama_fak_prod_unit, ' (', UPPER(vfu.jenjang), ')')
				ELSE COALESCE(vfu.nama_fak_prod_unit, '')
			END AS fakultas_unit_label,
			COALESCE(vfu.jenjang, '') AS jenjang,
			COALESCE(vfu.type, '') AS type,
			COALESCE(vfu.fakultas, '') AS fakultas,
			COALESCE((
				SELECT COUNT(*) FROM matrix_proker mpk
				WHERE mpk.id_mata_program = mptr.id_mata_programv2
				  AND mpk.id_fakultas_unit = tr.fakultas_unit
				  AND mpk.deleted_at IS NULL
			), 0) AS proker_count
		FROM mata_programv2_template_renstra mptr
		JOIN mata_programv2 mp ON mp.id = mptr.id_mata_programv2
		LEFT JOIN master_tahun mt ON mt.id = mp.id_master_tahun
		JOIN template_renstra tr ON tr.id = mptr.id_template_renstra
		LEFT JOIN v_fakultas_unit vfu ON tr.fakultas_unit = vfu.id
		WHERE mptr.id_mata_programv2 IN (?)
		  AND tr.fakultas_unit IS NOT NULL AND tr.fakultas_unit > 0
		  AND (mt.tahun IS NULL OR mt.tahun = '' OR CAST(tr.tahun AS CHAR) = CAST(mt.tahun AS CHAR))
		GROUP BY mptr.id_mata_programv2, tr.fakultas_unit, vfu.uuid, vfu.nama_fak_prod_unit, vfu.jenjang, vfu.type, vfu.fakultas
		ORDER BY vfu.type ASC, fakultas_unit_label ASC
	`
	if err := r.db.WithContext(ctx).Raw(sql, mpIDs).Scan(&rows).Error; err != nil {
		return result, err
	}

	for _, rw := range rows {
		var idVal uint
		if rw.FakultasUnitId != nil {
			idVal = *rw.FakultasUnitId
		}
		result[rw.IdMataProgramV2] = append(result[rw.IdMataProgramV2], domain.MataProgramV2UnitBadge{
			Id:               idVal,
			UUID:             rw.FakultasUnitUuid,
			FakultasUnitId:   rw.FakultasUnitId,
			FakultasUnitUuid: rw.FakultasUnitUuid,
			FakultasUnit:     rw.FakultasUnit,
			Jenjang:          rw.Jenjang,
			Type:             rw.Type,
			Fakultas:         rw.Fakultas,
			ProkerCount:      rw.ProkerCount,
		})
	}
	return result, nil
}

func (r *MataProgramV2Repository) GetTargetUnitsByMataProgramID(ctx context.Context, mpID uint) ([]domain.MataProgramV2UnitBadge, error) {
	unitsMap, err := r.GetTargetUnitsByMataProgramIDs(ctx, []uint{mpID})
	if err != nil {
		return nil, err
	}
	return unitsMap[mpID], nil
}

func (r *MataProgramV2Repository) GetTemplatesByMataProgramID(ctx context.Context, mpID uint) ([]domain.MataProgramV2TemplateItem, []uint, error) {
	tplMap, err := r.GetTemplatesByMataProgramIDs(ctx, []uint{mpID})
	if err != nil {
		return nil, nil, err
	}
	items := tplMap[mpID]
	ids := make([]uint, 0, len(items))
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
			mptr.id_mata_programv2 as id_mata_programv2,
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
		JOIN mata_programv2 mp ON mp.id = mptr.id_mata_programv2
		LEFT JOIN master_tahun mt ON mt.id = mp.id_master_tahun
		JOIN template_renstra tr ON mptr.id_template_renstra = tr.id
		LEFT JOIN master_indikator_renstra mir ON tr.indikator = mir.id
		LEFT JOIN master_standar_renstra msr ON mir.id_master_standar = msr.id
		LEFT JOIN v_fakultas_unit vfu ON tr.fakultas_unit = vfu.id
		WHERE mptr.id_mata_programv2 IN (?)
		  AND (mt.tahun IS NULL OR mt.tahun = '' OR CAST(tr.tahun AS CHAR) = CAST(mt.tahun AS CHAR))
		GROUP BY mptr.id_mata_programv2, tr.id
		ORDER BY msr.nama ASC, mir.indikator ASC, fakultas_unit ASC, tr.id DESC
	`
	err := r.db.WithContext(ctx).Raw(sql, mpIDs).Scan(&rows).Error
	if err != nil {
		return result, err
	}

	type dedupKey struct {
		mpID      uint
		indikator uint
		unit      uint
	}
	seen := make(map[dedupKey]bool, len(rows))
	for _, rw := range rows {
		dk := dedupKey{mpID: rw.IdMataProgramV2, indikator: rw.IndikatorRenstraID, unit: rw.FakultasUnitID}
		if seen[dk] {
			continue
		}
		seen[dk] = true
		result[rw.IdMataProgramV2] = append(result[rw.IdMataProgramV2], domain.MataProgramV2TemplateItem{
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
		})
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
			mpk.id_mata_program as id_mata_program,
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
		JOIN v_fakultas_unit vfu ON mpk.id_fakultas_unit = vfu.id
		WHERE mpk.id_mata_program IN (?)
		  AND mpk.deleted_at IS NULL
		GROUP BY mpk.id_mata_program, mpk.id_fakultas_unit, vfu.uuid, vfu.nama_fak_prod_unit, vfu.jenjang, vfu.type
		ORDER BY label ASC
	`
	err := r.db.WithContext(ctx).Raw(sql, mpIDs).Scan(&rows).Error
	if err != nil {
		return result, err
	}

	for _, rw := range rows {
		result[rw.IdMataProgram] = append(result[rw.IdMataProgram], domain.MataProgramV2ProkerUnit{
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
		})
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
