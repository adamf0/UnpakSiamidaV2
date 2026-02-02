package infrastructure

import (
	commondomaintemplaterenstra "UnpakSiamida/common/domain"
	domaintemplaterenstra "UnpakSiamida/modules/templaterenstra/domain"
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
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
	"tahun":                "tr.tahun",
	"indikator":            "i.indikator",
	"indikatorrenstrauuid": "i.uuid",
	"fakultasunit":         "fu.nama_fak_prod_unit",
	"kategori":             "tr.kategori",
	"klasifikasi":          "tr.klasifikasi",
	"tugas":                "tr.tugas",
}

// ------------------------
// GET ALL
// ------------------------
func (r *TemplateRenstraRepository) GetAll(
	ctx context.Context,
	search string,
	searchFilters []commondomaintemplaterenstra.SearchFilter,
	page, limit *int,
) ([]domaintemplaterenstra.TemplateRenstraDefault, int64, error) {

	var (
		result     = make([]domaintemplaterenstra.TemplateRenstraDefault, 0)
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

	var templaterenstras = make([]domaintemplaterenstra.TemplateRenstra, 0)

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

	var results = make([]domaintemplaterenstra.TemplateRenstraDefault, 0)

	err := r.db.WithContext(ctx).
		Table("template_renstra tr").
		Select(`
		tr.id AS ID,
		tr.uuid AS UUID,
		tr.tahun AS Tahun,

		s.uuid AS StandarRenstraUuid,
		s.id AS StandarRenstraID,
		s.nama AS StandarRenstra,

		i.uuid AS IndikatorRenstraUuid,
		tr.indikator AS IndikatorRenstraID,
		i.indikator AS Indikator,

		tr.pertanyaan AS IsPertanyaan,
		tr.fakultas_unit AS FakultasUnitID,
		fu.nama_fak_prod_unit AS FakultasUnit,
		tr.kategori AS Kategori,
		tr.klasifikasi AS Klasifikasi,
		tr.satuan AS Satuan,
		tr.target AS Target,
		tr.target_min AS TargetMin,
		tr.target_max AS TargetMax,
		tr.tugas AS Tugas
	`).
		Joins("INNER JOIN v_fakultas_unit fu ON tr.fakultas_unit = fu.id").
		Joins("INNER JOIN master_indikator_renstra i ON tr.indikator = i.id").
		Joins("INNER JOIN master_standar_renstra s ON i.id_master_standar = s.id").
		Where("tr.tahun = ? AND tr.fakultas_unit = ?", tahun, fakultasUnit).
		Find(&results).Error

	if err != nil {
		return nil, err
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

func (r *TemplateRenstraRepository) SetupUuid(ctx context.Context) error {
	const chunkSize = 500

	var ids []uint
	if err := r.db.WithContext(ctx).
		Model(&domaintemplaterenstra.TemplateRenstra{}).
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
			"UPDATE template_renstra SET uuid = %s WHERE id IN (?)",
			caseSQL,
		)

		if err := r.db.WithContext(ctx).
			Exec(query, args...).Error; err != nil {
			return err
		}
	}

	return nil
}
