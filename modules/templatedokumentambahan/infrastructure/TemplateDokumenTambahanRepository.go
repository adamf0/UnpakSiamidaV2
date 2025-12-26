package infrastructure

import (
	"context"
	commondomaintemplatedokumentambahan "UnpakSiamida/common/domain"
	domaintemplatedokumentambahan "UnpakSiamida/modules/templatedokumentambahan/domain"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"strings"
	"fmt"
)

type TemplateDokumenTambahanRepository struct {
	db *gorm.DB
}

func NewTemplateDokumenTambahanRepository(db *gorm.DB) domaintemplatedokumentambahan.ITemplateDokumenTambahanRepository {
	return &TemplateDokumenTambahanRepository{db: db}
}

// ------------------------
// GET BY UUID
// ------------------------
func (r *TemplateDokumenTambahanRepository) GetByUuid(ctx context.Context, uid uuid.UUID) (*domaintemplatedokumentambahan.TemplateDokumenTambahan, error) {
	var templatedokumentambahan domaintemplatedokumentambahan.TemplateDokumenTambahan

	err := r.db.WithContext(ctx).
		Where("uuid = ?", uid).
		First(&templatedokumentambahan).Error

	// if errors.Is(err, gorm.ErrRecordNotFound) {
	// 	return nil, nil
	// }

	if err != nil {
		return nil, err
	}

	return &templatedokumentambahan, nil
}

var allowedSearchColumns = map[string]string{
    // key:param -> db column
    "tahun":          	"dt.tahun",
	"pertanyaan":     	"dt.pertanyaan",
	"kategori":         "dt.fakultas_prodi_unit",
	"jenisfileuuid":    "jf.uuid",
}

// ------------------------
// GET ALL
// ------------------------
// func (r *TemplateDokumenTambahanRepository) GetAll(
// 	ctx context.Context,
// 	search string,
// 	searchFilters []commondomaintemplatedokumentambahan.SearchFilter,
// 	page, limit *int,
// ) ([]domaintemplatedokumentambahan.TemplateDokumenTambahan, int64, error) {

// 	var templatedokumentambahans []domaintemplatedokumentambahan.TemplateDokumenTambahan
// 	var total int64

// 	db := r.db.WithContext(ctx).Model(&domaintemplatedokumentambahan.TemplateDokumenTambahan{})

// 	// -------------------------------
// 	// SEARCH FILTERS (ADVANCED)
// 	// -------------------------------
// 	if len(searchFilters) > 0 {
// 		for _, f := range searchFilters {
// 			field := strings.TrimSpace(strings.ToLower(f.Field))
// 			operator := strings.TrimSpace(strings.ToLower(f.Operator))
			
// 			var value string
// 			if f.Value != nil {
// 				value = strings.TrimSpace(*f.Value)
// 			} else {
// 				value = "" // nil dianggap kosong
// 			}

// 			// if value == "" {
// 			// 	continue
// 			// }

// 			// Validate allowed column
// 			col, ok := allowedSearchColumns[field]
// 			if !ok {
// 				continue // skip unknown field
// 			}

// 			switch operator {
// 			case "eq":
// 				db = db.Where(fmt.Sprintf("%s = ?", col), value)
// 			case "neq":
// 				db = db.Where(fmt.Sprintf("%s <> ?", col), value)
// 			case "like":
// 				db = db.Where(fmt.Sprintf("%s LIKE ?", col), "%"+value+"%")
// 			case "gt":
// 				db = db.Where(fmt.Sprintf("%s > ?", col), value)
// 			case "gte":
// 				db = db.Where(fmt.Sprintf("%s >= ?", col), value)
// 			case "lt":
// 				db = db.Where(fmt.Sprintf("%s < ?", col), value)
// 			case "lte":
// 				db = db.Where(fmt.Sprintf("%s <= ?", col), value)
// 			case "in":
// 				db = db.Where(fmt.Sprintf("%s IN (?)", col), strings.Split(value, ","))
// 			default:
// 				// default fallback → LIKE
// 				db = db.Where(fmt.Sprintf("%s LIKE ?", col), "%"+value+"%")
// 			}
// 		}

// 	}
// 	if strings.TrimSpace(search) != "" {

// 		// -------------------------------
// 		// GLOBAL SEARCH
// 		// -------------------------------
// 		like := "%" + search + "%"
// 		var orParts []string
// 		var params []interface{}

// 		for _, col := range allowedSearchColumns {
// 			orParts = append(orParts, fmt.Sprintf("%s LIKE ?", col))
// 			params = append(params, like)
// 		}

// 		db = db.Where("(" + strings.Join(orParts, " OR ") + ")", params...)
// 	}

// 	// -------------------------------
// 	// COUNT
// 	// -------------------------------
// 	if err := db.Count(&total).Error; err != nil {
// 		return nil, 0, err
// 	}

// 	// -------------------------------
// 	// PAGINATION
// 	// -------------------------------
// 	if page != nil && limit != nil && *limit > 0 {
// 		p := *page
// 		l := *limit

// 		if p < 1 {
// 			p = 1
// 		}

// 		offset := (p - 1) * l
// 		db = db.Offset(offset).Limit(l)
// 	}

// 	// -------------------------------
// 	// EXECUTE QUERY
// 	// -------------------------------
// 	if err := db.Find(&templatedokumentambahans).Error; err != nil {
// 		return nil, 0, err
// 	}

// 	return templatedokumentambahans, total, nil
// }
func (r *TemplateDokumenTambahanRepository) GetAll(
	ctx context.Context,
	search string,
	searchFilters []commondomaintemplatedokumentambahan.SearchFilter,
	page, limit *int,
) ([]domaintemplatedokumentambahan.TemplateDokumenTambahanDefault, int64, error) {

	var (
		result     []domaintemplatedokumentambahan.TemplateDokumenTambahanDefault
		total      int64
		conditions []string
		args       []interface{}
	)

	// =====================================================
	// BASE FROM + JOIN (WAJIB PAKAI INI)
	// =====================================================
	baseFrom := `
		FROM template_dokumen_tambahan dt
		INNER JOIN jenis_file_renstra jf ON jf.id = dt.jenis_file
	`

	// =====================================================
	// ADVANCED FILTERS
	// =====================================================
	for _, f := range searchFilters {
		field := strings.TrimSpace(strings.ToLower(f.Field))
		operator := strings.TrimSpace(strings.ToLower(f.Operator))

		col, ok := allowedSearchColumns[field] // harus sudah pakai alias tr./i./fu.
		if !ok {
			continue
		}

		value := ""
		if f.Value != nil {
			value = strings.TrimSpace(*f.Value)
		}

		switch operator {
		case "eq":
			conditions = append(conditions, fmt.Sprintf("%s = ?", col))
			args = append(args, value)

		case "neq":
			conditions = append(conditions, fmt.Sprintf("%s <> ?", col))
			args = append(args, value)

		case "like":
			conditions = append(conditions, fmt.Sprintf("%s LIKE ?", col))
			args = append(args, "%"+value+"%")

		case "gt":
			conditions = append(conditions, fmt.Sprintf("%s > ?", col))
			args = append(args, value)

		case "gte":
			conditions = append(conditions, fmt.Sprintf("%s >= ?", col))
			args = append(args, value)

		case "lt":
			conditions = append(conditions, fmt.Sprintf("%s < ?", col))
			args = append(args, value)

		case "lte":
			conditions = append(conditions, fmt.Sprintf("%s <= ?", col))
			args = append(args, value)

		case "in":
			values := strings.Split(value, ",")
			if len(values) > 0 {
				placeholders := strings.TrimRight(strings.Repeat("?,", len(values)), ",")
				conditions = append(
					conditions,
					fmt.Sprintf("%s IN (%s)", col, placeholders),
				)
				for _, v := range values {
					args = append(args, strings.TrimSpace(v))
				}
			}

		default:
			conditions = append(conditions, fmt.Sprintf("%s LIKE ?", col))
			args = append(args, "%"+value+"%")
		}
	}

	// =====================================================
	// GLOBAL SEARCH
	// =====================================================
	if strings.TrimSpace(search) != "" {
		var orParts []string
		like := "%" + search + "%"

		for _, col := range allowedSearchColumns {
			orParts = append(orParts, fmt.Sprintf("%s LIKE ?", col))
			args = append(args, like)
		}

		conditions = append(
			conditions,
			"("+strings.Join(orParts, " OR ")+")",
		)
	}

	// =====================================================
	// WHERE CLAUSE
	// =====================================================
	whereClause := ""
	if len(conditions) > 0 {
		whereClause = " WHERE " + strings.Join(conditions, " AND ")
	}

	// =====================================================
	// COUNT QUERY (HARUS PAKAI JOIN JUGA)
	// =====================================================
	countQuery := `
		SELECT COUNT(DISTINCT dt.id)
	` + baseFrom + whereClause

	if err := r.db.WithContext(ctx).
		Raw(countQuery, args...).
		Scan(&total).Error; err != nil {
		return nil, 0, err
	}

	// =====================================================
	// PAGINATION
	// =====================================================
	pagination := ""
	if page != nil && limit != nil && *limit > 0 {
		p := *page
		l := *limit
		if p < 1 {
			p = 1
		}

		offset := (p - 1) * l
		pagination = " LIMIT ? OFFSET ?"
		args = append(args, l, offset)
	}

	orderBy := " ORDER BY dt.tahun DESC, dt.id DESC"

	// =====================================================
	// SELECT QUERY (INI YANG KAMU MINTA)
	// =====================================================
	selectQuery := `
		SELECT
			dt.id             		AS ID,
			dt.uuid           		AS UUID,
			dt.tahun          		AS Tahun,
			jf.id             		AS JenisFileId,
			jf.uuid           		AS JenisFileUuid,
			jf.nama           		AS JenisFile,
			dt.fakultas_prodi_unit 	AS FakultasProdiUnit,
			dt.pertanyaan     		AS Pertanyaan,
			dt.klasifikasi    		AS Klasifikasi,
			dt.tugas    			AS Tugas
	` + baseFrom + whereClause + orderBy + pagination

	if err := r.db.WithContext(ctx).
		Raw(selectQuery, args...).
		Scan(&result).Error; err != nil {
		return nil, 0, err
	}

	return result, total, nil
}

// ------------------------
// GET ALL BY Tahun & FakulktasProdiUnit
// ------------------------
func (r *TemplateDokumenTambahanRepository) GetAllByTahunFakUnitDefault(
	ctx context.Context,
	tahun string,
	fakultasProdiUnit string,
) ([]domaintemplatedokumentambahan.TemplateDokumenTambahanDefault, error) {

	query := `
		SELECT
			dt.id             		AS ID,
			dt.uuid           		AS UUID,
			dt.tahun          		AS Tahun,
			jf.id             		AS JenisFileId,
			jf.uuid           		AS JenisFileUuid,
			jf.nama           		AS JenisFile,
			dt.fakultas_prodi_unit 	AS FakultasProdiUnit,
			dt.pertanyaan     		AS Pertanyaan,
			dt.klasifikasi    		AS Klasifikasi,
			dt.tugas    			AS Tugas
		FROM template_dokumen_tambahan dt
		INNER JOIN jenis_file_renstra jf ON jf.id = dt.jenis_file
		WHERE dt.tahun = ?
		  AND dt.fakultas_prodi_unit = ?
	`

	rows, err := r.db.WithContext(ctx).Raw(query, tahun, fakultasProdiUnit).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	results := make([]domaintemplatedokumentambahan.TemplateDokumenTambahanDefault, 0)

	for rows.Next() {
		var item domaintemplatedokumentambahan.TemplateDokumenTambahanDefault

		err := rows.Scan(
			&item.ID,
			&item.UUID,
			&item.Tahun,
			&item.JenisFileID,
			&item.JenisFileUuid,
			&item.JenisFile,
			&item.FakultasProdiUnit,
			&item.Pertanyaan,
			&item.Klasifikasi,
			&item.Tugas,
		)
		if err != nil {
			return nil, err
		}

		results = append(results, item)
	}

	return results, nil
}

// ------------------------
// CREATE
// ------------------------
func (r *TemplateDokumenTambahanRepository) Create(ctx context.Context, templatedokumentambahan *domaintemplatedokumentambahan.TemplateDokumenTambahan) error {
	// return r.db.WithContext(ctx).Create(templatedokumentambahan).Error
	return r.db.WithContext(ctx).
		Clauses(clause.OnConflict{
			Columns: []clause.Column{
				{Name: "tahun"},
				{Name: "jenis_file"},
				{Name: "fakultas_prodi_unit"},
			},
			DoUpdates: clause.AssignmentColumns([]string{
				"pertanyaan",
				"klasifikasi",
				"tugas",
			}),
		}).
		Create(templatedokumentambahan).
		Error
}

// ------------------------
// UPDATE
// ------------------------
func (r *TemplateDokumenTambahanRepository) Update(ctx context.Context, templatedokumentambahan *domaintemplatedokumentambahan.TemplateDokumenTambahan) error {
	return r.db.WithContext(ctx).Save(templatedokumentambahan).Error
}

// ------------------------
// DELETE (by UUID)
// ------------------------
func (r *TemplateDokumenTambahanRepository) Delete(ctx context.Context, uid uuid.UUID) error {
	return r.db.WithContext(ctx).
		Where("uuid = ?", uid).
		Delete(&domaintemplatedokumentambahan.TemplateDokumenTambahan{}).Error
}
