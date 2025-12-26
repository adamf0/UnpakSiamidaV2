package infrastructure

import (
	"context"
	commondomaintemplaterenstra "UnpakSiamida/common/domain"
	domaintemplaterenstra "UnpakSiamida/modules/templaterenstra/domain"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"strings"
	"fmt"
)

type TemplateRenstraRepository struct {
	db *gorm.DB
}

func NewTemplateRenstraRepository(db *gorm.DB) domaintemplaterenstra.ITemplateRenstraRepository {
	return &TemplateRenstraRepository{db: db}
}

// ------------------------
// GET BY UUID
// ------------------------
func (r *TemplateRenstraRepository) GetByUuid(ctx context.Context, uid uuid.UUID) (*domaintemplaterenstra.TemplateRenstra, error) {
	var templaterenstra domaintemplaterenstra.TemplateRenstra

	err := r.db.WithContext(ctx).
		Where("uuid = ?", uid).
		First(&templaterenstra).Error

	// if errors.Is(err, gorm.ErrRecordNotFound) {
	// 	return nil, nil
	// }

	if err != nil {
		return nil, err
	}

	return &templaterenstra, nil
}

var allowedSearchColumns = map[string]string{
    // key:param -> db column
    "tahun":"tr.tahun",
	"indikator":"i.indikator",
	"indikatorrenstrauuid":"i.uuid",
	"fakultasunit":"fu.nama_fak_prod_unit",
	"kategori":"tr.kategori",
	"klasifikasi":"tr.klasifikasi",
	"tugas":"tr.tugas",
}

// ------------------------
// GET ALL
// ------------------------
// func (r *TemplateRenstraRepository) GetAll(
// 	ctx context.Context,
// 	search string,
// 	searchFilters []commondomaintemplaterenstra.SearchFilter,
// 	page, limit *int,
// ) ([]domaintemplaterenstra.TemplateRenstra, int64, error) {

// 	var templaterenstras []domaintemplaterenstra.TemplateRenstra
// 	var total int64

// 	db := r.db.WithContext(ctx).Model(&domaintemplaterenstra.TemplateRenstra{})

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
// 	if err := db.Find(&templaterenstras).Error; err != nil {
// 		return nil, 0, err
// 	}

// 	return templaterenstras, total, nil
// }
func (r *TemplateRenstraRepository) GetAll(
	ctx context.Context,
	search string,
	searchFilters []commondomaintemplaterenstra.SearchFilter,
	page, limit *int,
) ([]domaintemplaterenstra.TemplateRenstraDefault, int64, error) {

	var (
		result     []domaintemplaterenstra.TemplateRenstraDefault
		total      int64
		conditions []string
		args       []interface{}
	)

	// =====================================================
	// BASE FROM + JOIN (WAJIB PAKAI INI)
	// =====================================================
	baseFrom := `
		FROM template_renstra tr
		INNER JOIN v_fakultas_unit fu ON tr.fakultas_unit = fu.id
		INNER JOIN master_indikator_renstra i ON tr.indikator = i.id
		INNER JOIN master_standar_renstra s ON i.id_master_standar = s.id
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
		SELECT COUNT(DISTINCT tr.id)
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

	orderBy := " ORDER BY tr.tahun DESC, tr.id DESC"

	// =====================================================
	// SELECT QUERY (INI YANG KAMU MINTA)
	// =====================================================
	selectQuery := `
		SELECT
			tr.id                     AS ID,
			tr.uuid                   AS UUID,
			tr.tahun                  AS Tahun,
			
			s.uuid                    AS StandarRenstraUuid,
			s.id              		  AS StandarRenstraID,
			s.nama               	  AS StandarRenstra,

			i.uuid                    AS IndikatorRenstraUuid,
			tr.indikator              AS IndikatorRenstraID,
			i.indikator               AS Indikator,

			tr.pertanyaan             AS IsPertanyaan,
			tr.fakultas_unit          AS FakultasUnitID,
			fu.nama_fak_prod_unit     AS FakultasUnit,
			tr.kategori               AS Kategori,
			tr.klasifikasi            AS Klasifikasi,
			tr.satuan                 AS Satuan,
			tr.target                 AS Target,
			tr.target_min             AS TargetMin,
			tr.target_max             AS TargetMax,
			tr.tugas                  AS Tugas
	` + baseFrom + whereClause + orderBy + pagination

	if err := r.db.WithContext(ctx).
		Raw(selectQuery, args...).
		Scan(&result).Error; err != nil {
		return nil, 0, err
	}

	return result, total, nil
}

// ------------------------
// GET ALL BY Tahun & FakultasUnit
// ------------------------
func (r *TemplateRenstraRepository) GetAllByTahunFakUnit(
	ctx context.Context,
	tahun string,
	fakultasUnit uint,
) ([]domaintemplaterenstra.TemplateRenstra, error) {

	var templaterenstras []domaintemplaterenstra.TemplateRenstra

	err := r.db.WithContext(ctx).
		Where("tahun = ? AND fakultas_unit = ?", tahun, fakultasUnit).
		Find(&templaterenstras).Error

	if err != nil {
		return nil, err
	}

	return templaterenstras, nil
}

// ------------------------
// GET ALL BY Tahun & FakultasUnit (default)
// ------------------------
func (r *TemplateRenstraRepository) GetAllByTahunFakUnitDefault(
	ctx context.Context,
	tahun string,
	fakultasUnit uint,
) ([]domaintemplaterenstra.TemplateRenstraDefault, error) {

	query := `
		SELECT
			tr.id                     AS ID,
			tr.uuid                   AS UUID,
			tr.tahun                  AS Tahun,
            
			s.uuid                    AS StandarRenstraUuid,
			s.id              		  AS StandarRenstraID,
			s.nama               	  AS StandarRenstra,

			i.uuid              	  AS IndikatorRenstraUuid,
			tr.indikator              AS IndikatorRenstraID,
			i.indikator               AS Indikator,

			tr.pertanyaan          	  AS IsPertanyaan,
			tr.fakultas_unit          AS FakultasUnitID,
			fu.nama_fak_prod_unit     AS FakultasUnit,
			tr.kategori               AS Kategori,
			tr.klasifikasi            AS Klasifikasi,
			tr.satuan                 AS Satuan,
			tr.target                 AS Target,
			tr.target_min             AS TargetMin,
			tr.target_max             AS TargetMax,
			tr.tugas                  AS Tugas
		FROM template_renstra tr
		INNER JOIN v_fakultas_unit fu ON tr.fakultas_unit = fu.id 
		INNER JOIN master_indikator_renstra i ON tr.indikator = i.id 
		INNER JOIN master_standar_renstra s ON i.id_master_standar = s.id
		WHERE tr.tahun = ?
		  AND tr.fakultas_unit = ?
	`

	rows, err := r.db.WithContext(ctx).Raw(query, tahun, fakultasUnit).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	results := make([]domaintemplaterenstra.TemplateRenstraDefault, 0)

	for rows.Next() {
		var item domaintemplaterenstra.TemplateRenstraDefault

		err := rows.Scan(
			&item.ID,
			&item.UUID,
			&item.Tahun,
			
			&item.StandarRenstraUuid,
			&item.StandarRenstraID,
			&item.StandarRenstra,

			&item.IndikatorRenstraUuid,
			&item.IndikatorRenstraID,
			&item.Indikator,
			
			&item.IsPertanyaan,
			&item.FakultasUnit,
			&item.Kategori,
			&item.Klasifikasi,
			&item.Satuan,
			&item.Target,
			&item.TargetMin,
			&item.TargetMax,
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
func (r *TemplateRenstraRepository) Create(ctx context.Context, templaterenstra *domaintemplaterenstra.TemplateRenstra) error {
	// return r.db.WithContext(ctx).Create(templaterenstra).Error
	return r.db.WithContext(ctx).
		Clauses(clause.OnConflict{
			Columns: []clause.Column{
				{Name: "tahun"},
				{Name: "indikator"},
				{Name: "fakultas_unit"},
			},
			DoUpdates: clause.AssignmentColumns([]string{
				"pertanyaan",
				"kategori",
				"klasifikasi",
				"satuan",
				"target",
				"target_min",
				"target_max",
				"tugas",
			}),
		}).
		Create(templaterenstra).
		Error
}

// ------------------------
// UPDATE
// ------------------------
func (r *TemplateRenstraRepository) Update(ctx context.Context, templaterenstra *domaintemplaterenstra.TemplateRenstra) error {
	return r.db.WithContext(ctx).Save(templaterenstra).Error
}

// ------------------------
// DELETE (by UUID)
// ------------------------
func (r *TemplateRenstraRepository) Delete(ctx context.Context, uid uuid.UUID) error {
	return r.db.WithContext(ctx).
		Where("uuid = ?", uid).
		Delete(&domaintemplaterenstra.TemplateRenstra{}).Error
}
