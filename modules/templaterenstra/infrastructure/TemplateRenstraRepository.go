package infrastructure

import (
	"context"
	commondomaintemplaterenstra "UnpakSiamida/common/domain"
	domaintemplaterenstra "UnpakSiamida/modules/templaterenstra/domain"
	"github.com/google/uuid"
	"gorm.io/gorm"
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
    "nama":          "nama",
}

// ------------------------
// GET ALL
// ------------------------
func (r *TemplateRenstraRepository) GetAll(
	ctx context.Context,
	search string,
	searchFilters []commondomaintemplaterenstra.SearchFilter,
	page, limit *int,
) ([]domaintemplaterenstra.TemplateRenstra, int64, error) {

	var templaterenstras []domaintemplaterenstra.TemplateRenstra
	var total int64

	db := r.db.WithContext(ctx).Model(&domaintemplaterenstra.TemplateRenstra{})

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
				value = "" // nil dianggap kosong
			}

			// if value == "" {
			// 	continue
			// }

			// Validate allowed column
			col, ok := allowedSearchColumns[field]
			if !ok {
				continue // skip unknown field
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
				// default fallback → LIKE
				db = db.Where(fmt.Sprintf("%s LIKE ?", col), "%"+value+"%")
			}
		}

	}
	if strings.TrimSpace(search) != "" {

		// -------------------------------
		// GLOBAL SEARCH
		// -------------------------------
		like := "%" + search + "%"
		var orParts []string
		var params []interface{}

		for _, col := range allowedSearchColumns {
			orParts = append(orParts, fmt.Sprintf("%s LIKE ?", col))
			params = append(params, like)
		}

		db = db.Where("(" + strings.Join(orParts, " OR ") + ")", params...)
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
	if err := db.Find(&templaterenstras).Error; err != nil {
		return nil, 0, err
	}

	return templaterenstras, total, nil
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
            i.uuid              	  AS IndikatorRenstraUuid,
			tr.indikator              AS IndikatorRenstraID,
			i.indikator               AS Indikator,
			tr.pertanyaan          	  AS IsPertanyaan,
			tr.fakultas_unit          AS FakultasUnit,
			tr.kategori               AS Kategori,
			tr.klasifikasi            AS Klasifikasi,
			tr.satuan                 AS Satuan,
			tr.target                 AS Target,
			tr.target_min             AS TargetMin,
			tr.target_max             AS TargetMax,
			tr.tugas                  AS Tugas
		FROM template_renstra tr
		INNER JOIN master_indikator_renstra i ON tr.indikator = i.id 
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
	return r.db.WithContext(ctx).Create(templaterenstra).Error
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
