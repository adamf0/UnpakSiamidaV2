package infrastructure

import (
	commondomainfakultasunit "UnpakSiamida/common/domain"
	domainfakultasunit "UnpakSiamida/modules/fakultasunit/domain"
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type FakultasUnitRepository struct {
	db             *gorm.DB
	isActiveMu     sync.RWMutex
	isActiveCached *bool
}

func NewFakultasUnitRepository(db *gorm.DB) domainfakultasunit.IFakultasUnitRepository {
	return &FakultasUnitRepository{db: db}
}

func (r *FakultasUnitRepository) hasViewIsActiveColumn(ctx context.Context) bool {
	r.isActiveMu.RLock()
	if r.isActiveCached != nil {
		val := *r.isActiveCached
		r.isActiveMu.RUnlock()
		return val
	}
	r.isActiveMu.RUnlock()

	var colCount int64
	err := r.db.WithContext(ctx).Raw(`
		SELECT COUNT(*)
		FROM INFORMATION_SCHEMA.COLUMNS
		WHERE TABLE_SCHEMA = DATABASE()
		  AND TABLE_NAME = 'v_fakultas_unit'
		  AND COLUMN_NAME = 'is_active'
	`).Scan(&colCount).Error

	hasCol := err == nil && colCount > 0
	r.isActiveMu.Lock()
	r.isActiveCached = &hasCol
	r.isActiveMu.Unlock()
	return hasCol
}

func (r *FakultasUnitRepository) setHasViewIsActiveColumn(val bool) {
	r.isActiveMu.Lock()
	r.isActiveCached = &val
	r.isActiveMu.Unlock()
}

func (r *FakultasUnitRepository) selectColumns(hasIsActive bool) string {
	isActiveExpr := "1 AS is_active"
	if hasIsActive {
		isActiveExpr = "COALESCE(is_active, 1) AS is_active"
	}

	return fmt.Sprintf(`
		id,
		uuid,
		CASE
			WHEN type = 'prodi' AND jenjang IS NOT NULL AND jenjang != ''
			THEN CONCAT(nama_fak_prod_unit, ' (', UPPER(jenjang), ')')
			ELSE nama_fak_prod_unit
		END AS nama_fak_prod_unit,
		nama_fak_prod_unit AS nama_raw,
		kode_jenjang,
		jenjang,
		type,
		fakultas,
		kode_fakultas,
		%s
	`, isActiveExpr)
}

// ------------------------
// GET BY UUID
// ------------------------
func (r *FakultasUnitRepository) GetDefaultByUuid(ctx context.Context, uid uuid.UUID) (*domainfakultasunit.FakultasUnit, error) {
	var fakultasunit domainfakultasunit.FakultasUnit
	hasIsActive := r.hasViewIsActiveColumn(ctx)

	err := r.db.WithContext(ctx).
		Select(r.selectColumns(hasIsActive)).
		Where("uuid = ?", uid).
		First(&fakultasunit).Error

	if err != nil {
		return nil, err
	}

	return &fakultasunit, nil
}

var allowedSearchColumns = map[string]string{
	// key:param -> db column
	"nama_fak_prod_unit": "nama_fak_prod_unit",
	"nama":               "nama_fak_prod_unit",
	"jenjang":            "jenjang",
	"fakultas":           "fakultas",
	"type":               "type",
	"kode_fakultas":      "kode_fakultas",
	"is_active":          "is_active",
}

var globalSearchColumns = []string{
	"nama_fak_prod_unit",
	"jenjang",
	"fakultas",
	"type",
	"kode_fakultas",
}

// ------------------------
// GET ALL
// ------------------------
func (r *FakultasUnitRepository) GetAll(
	ctx context.Context,
	search string,
	searchFilters []commondomainfakultasunit.SearchFilter,
	page, limit *int,
	includeHidden bool,
) ([]domainfakultasunit.FakultasUnit, int64, error) {

	var fakultasunits = make([]domainfakultasunit.FakultasUnit, 0)
	var total int64
	hasIsActive := r.hasViewIsActiveColumn(ctx)

	db := r.db.WithContext(ctx).Model(&domainfakultasunit.FakultasUnit{})

	if !includeHidden && hasIsActive {
		db = db.Where("COALESCE(is_active, 1) = 1")
	}

	// -------------------------------
	// SEARCH FILTERS (ADVANCED)
	// -------------------------------
	if len(searchFilters) > 0 {
		for _, f := range searchFilters {
			field := strings.TrimSpace(strings.ToLower(f.Field))
			operator := strings.TrimSpace(strings.ToLower(f.Operator))

			var value string
			if f.Value != nil {
				value = strings.TrimSpace(*f.Value)
			} else {
				value = ""
			}

			col, ok := allowedSearchColumns[field]
			if !ok {
				continue
			}
			if col == "is_active" && !hasIsActive {
				continue
			}

			switch operator {
			case "eq":
				db = db.Where(fmt.Sprintf("%s = ?", col), value)
			case "neq":
				db = db.Where(fmt.Sprintf("%s <> ?", col), value)
			case "like":
				db = db.Where(fmt.Sprintf("%s LIKE ?", col), "%"+value+"%")
			case "gt":
				db = db.Where(fmt.Sprintf("%s > ?", col), value)
			case "gte":
				db = db.Where(fmt.Sprintf("%s >= ?", col), value)
			case "lt":
				db = db.Where(fmt.Sprintf("%s < ?", col), value)
			case "lte":
				db = db.Where(fmt.Sprintf("%s <= ?", col), value)
			case "in":
				db = db.Where(fmt.Sprintf("%s IN (?)", col), strings.Split(value, ","))
			default:
				db = db.Where(fmt.Sprintf("%s LIKE ?", col), "%"+value+"%")
			}
		}
	}

	if strings.TrimSpace(search) != "" {
		like := "%" + search + "%"
		var orParts []string
		var params []interface{}

		for _, col := range globalSearchColumns {
			orParts = append(orParts, fmt.Sprintf("%s LIKE ?", col))
			params = append(params, like)
		}

		db = db.Where("("+strings.Join(orParts, " OR ")+")", params...)
	}

	// -------------------------------
	// COUNT
	// -------------------------------
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// -------------------------------
	// PAGINATION
	// -------------------------------
	if page != nil && limit != nil && *limit > 0 {
		p := *page
		l := *limit

		if p < 1 {
			p = 1
		}

		offset := (p - 1) * l
		db = db.Offset(offset).Limit(l)
	}

	// -------------------------------
	// EXECUTE QUERY
	// -------------------------------
	if err := db.Select(r.selectColumns(hasIsActive)).Order("id ASC").Find(&fakultasunits).Error; err != nil {
		return nil, 0, err
	}

	return fakultasunits, total, nil
}

// ------------------------
// UPDATE TARGET (NO DELETE)
// ------------------------
func (r *FakultasUnitRepository) UpdateTarget(
	ctx context.Context,
	uid uuid.UUID,
	nama *string,
	isActive *bool,
) error {
	updates := map[string]interface{}{}
	if nama != nil {
		updates["nama"] = *nama
	}
	if isActive != nil {
		if *isActive {
			updates["is_active"] = 1
		} else {
			updates["is_active"] = 0
		}
	}
	if len(updates) == 0 {
		return nil
	}

	return r.db.WithContext(ctx).
		Table("sijamu_fakultas_unit").
		Where("uuid = ?", uid.String()).
		Updates(updates).Error
}

func (r *FakultasUnitRepository) ensureUuidColumn(ctx context.Context, tableName string) error {
	var colCount int64
	checkSQL := `
		SELECT COUNT(*)
		FROM INFORMATION_SCHEMA.COLUMNS
		WHERE TABLE_SCHEMA = DATABASE()
		  AND TABLE_NAME = ?
		  AND COLUMN_NAME = 'uuid'
	`
	if err := r.db.WithContext(ctx).Raw(checkSQL, tableName).Scan(&colCount).Error; err != nil {
		return err
	}

	if colCount == 0 {
		alterSQL := fmt.Sprintf("ALTER TABLE `%s` ADD COLUMN `uuid` CHAR(36) NULL AFTER `id`", tableName)
		if err := r.db.WithContext(ctx).Exec(alterSQL).Error; err != nil {
			return err
		}
	}

	return nil
}

func (r *FakultasUnitRepository) ensureIsActiveColumn(ctx context.Context) error {
	var colCount int64
	checkSQL := `
		SELECT COUNT(*)
		FROM INFORMATION_SCHEMA.COLUMNS
		WHERE TABLE_SCHEMA = DATABASE()
		  AND TABLE_NAME = 'sijamu_fakultas_unit'
		  AND COLUMN_NAME = 'is_active'
	`
	if err := r.db.WithContext(ctx).Raw(checkSQL).Scan(&colCount).Error; err != nil {
		return err
	}

	if colCount == 0 {
		alterSQL := "ALTER TABLE `sijamu_fakultas_unit` ADD COLUMN `is_active` TINYINT(1) NOT NULL DEFAULT 1 AFTER `standalone`"
		if err := r.db.WithContext(ctx).Exec(alterSQL).Error; err != nil {
			return err
		}
	}

	return nil
}

func (r *FakultasUnitRepository) ensureViewFakultasUnit(ctx context.Context) error {
	var tableCount int64
	if err := r.db.WithContext(ctx).Raw(`
		SELECT COUNT(*)
		FROM INFORMATION_SCHEMA.TABLES
		WHERE TABLE_SCHEMA = DATABASE()
		  AND TABLE_NAME IN ('sijamu_fakultas_unit', 'm_fakultas', 'm_program_studi')
	`).Scan(&tableCount).Error; err != nil {
		return err
	}

	if tableCount < 3 {
		return nil
	}

	var viewDef string
	_ = r.db.WithContext(ctx).Raw(`
		SELECT VIEW_DEFINITION
		FROM INFORMATION_SCHEMA.VIEWS
		WHERE TABLE_SCHEMA = DATABASE()
		  AND TABLE_NAME = 'v_fakultas_unit'
	`).Scan(&viewDef).Error

	hasActive := strings.Contains(viewDef, "is_active")
	hasStandaloneRename := strings.Contains(viewDef, "standalone") && strings.Contains(viewDef, "trim")

	if hasActive && hasStandaloneRename {
		r.setHasViewIsActiveColumn(true)
		return nil
	}

	createViewSQL := `
		CREATE OR REPLACE VIEW v_fakultas_unit AS
		SELECT
			sfu.id AS id,
			sfu.uuid AS uuid,
			sfu.kode_fakultas AS kode_fakultas,
			CASE
				WHEN sfu.standalone AND sfu.nama IS NOT NULL AND TRIM(sfu.nama) <> '' THEN sfu.nama
				WHEN sfu.standalone THEN fak.nama_fakultas
				WHEN sfu.kode_fakultas IS NOT NULL AND sfu.kode_prodi IS NOT NULL AND (sfu.nama IS NULL OR TRIM(sfu.nama) = '') THEN prod.nama_prodi
				WHEN sfu.kode_fakultas IS NOT NULL AND sfu.kode_prodi IS NOT NULL AND sfu.nama IS NOT NULL THEN sfu.nama
				WHEN sfu.kode_fakultas IS NOT NULL AND sfu.kode_prodi IS NULL AND sfu.nama IS NOT NULL THEN sfu.nama
				ELSE 'tidak diketahui'
			END AS nama_fak_prod_unit,
			prod.kode_jenjang AS kode_jenjang,
			CASE
				WHEN prod.kode_jenjang = 'C' THEN 's1'
				WHEN prod.kode_jenjang = 'B' THEN 's2'
				WHEN prod.kode_jenjang = 'A' THEN 's3'
				WHEN prod.kode_jenjang = 'E' THEN 'd3'
				WHEN prod.kode_jenjang = 'D' THEN 'd4'
				WHEN prod.kode_jenjang = 'J' THEN 'profesi'
				ELSE NULL
			END AS jenjang,
			CASE
				WHEN prod.kode_jenjang = 'C' THEN '1'
				WHEN prod.kode_jenjang = 'B' THEN '2'
				WHEN prod.kode_jenjang = 'A' THEN '3'
				WHEN prod.kode_jenjang = 'E' THEN '4'
				WHEN prod.kode_jenjang = 'D' THEN '5'
				WHEN prod.kode_jenjang = 'J' THEN '6'
				ELSE '7'
			END AS jenjang_int,
			CASE
				WHEN sfu.standalone THEN 'fakultas'
				WHEN sfu.kode_fakultas IS NOT NULL AND sfu.kode_prodi IS NOT NULL THEN 'prodi'
				WHEN sfu.kode_fakultas IS NOT NULL AND sfu.kode_prodi IS NULL AND sfu.nama IS NOT NULL THEN 'unit'
				ELSE NULL
			END AS type,
			CASE
				WHEN sfu.standalone AND sfu.nama IS NOT NULL AND TRIM(sfu.nama) <> '' THEN sfu.nama
				WHEN sfu.standalone THEN fak.nama_fakultas
				WHEN sfu.kode_fakultas IS NOT NULL AND sfu.kode_prodi IS NOT NULL THEN fak.nama_fakultas
				ELSE NULL
			END AS fakultas,
			sfu.is_active AS is_active
		FROM sijamu_fakultas_unit sfu
		LEFT JOIN m_fakultas fak ON sfu.kode_fakultas = fak.kode_fakultas
		LEFT JOIN m_program_studi prod ON sfu.kode_prodi = prod.kode_prodi
	`

	if err := r.db.WithContext(ctx).Exec(createViewSQL).Error; err != nil {
		return err
	}

	r.setHasViewIsActiveColumn(true)
	return nil
}

func (r *FakultasUnitRepository) EnsureSchema(ctx context.Context) error {
	if err := r.SetupUuid(ctx); err != nil {
		return err
	}
	if err := r.ensureIsActiveColumn(ctx); err != nil {
		return err
	}
	if err := r.ensureViewFakultasUnit(ctx); err != nil {
		return err
	}
	return nil
}

func (r *FakultasUnitRepository) SetupUuid(ctx context.Context) error {
	if err := r.ensureUuidColumn(ctx, "sijamu_fakultas_unit"); err != nil {
		return err
	}

	const chunkSize = 500

	var ids []uint
	if err := r.db.WithContext(ctx).
		Table("sijamu_fakultas_unit").
		Where("uuid IS NULL OR uuid = ''").
		Pluck("id", &ids).Error; err != nil {
		return err
	}

	if len(ids) == 0 {
		return nil
	}

	for i := 0; i < len(ids); i += chunkSize {
		end := i + chunkSize
		if end > len(ids) {
			end = len(ids)
		}

		chunk := ids[i:end]

		caseSQL := "CASE id "
		args := make([]any, 0, len(chunk)*2+1)

		for _, id := range chunk {
			u := uuid.NewString()
			caseSQL += "WHEN ? THEN ? "
			args = append(args, id, u)
		}

		caseSQL += "END"
		args = append(args, chunk)

		query := fmt.Sprintf(
			"UPDATE sijamu_fakultas_unit SET uuid = %s WHERE id IN (?)",
			caseSQL,
		)

		if err := r.db.WithContext(ctx).
			Exec(query, args...).Error; err != nil {
			return err
		}
	}

	return nil
}
