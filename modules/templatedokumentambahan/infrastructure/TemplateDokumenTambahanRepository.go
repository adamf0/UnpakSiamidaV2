package infrastructure

import (
	commondomaintemplatedokumentambahan "UnpakSiamida/common/domain"
	domaintemplatedokumentambahan "UnpakSiamida/modules/templatedokumentambahan/domain"
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
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
	"tahun":         "dt.tahun",
	"pertanyaan":    "dt.pertanyaan",
	"kategori":      "dt.fakultas_prodi_unit",
	"jenisfile":     "jf.nama",
	"jenisfileuuid": "jf.uuid",
}

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

func (r *TemplateDokumenTambahanRepository) SetupUuid(ctx context.Context) error {
	const chunkSize = 500

	var ids []uint
	if err := r.db.WithContext(ctx).
		Model(&domaintemplatedokumentambahan.TemplateDokumenTambahan{}).
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
			"UPDATE template_dokumen_tambahan SET uuid = %s WHERE id IN (?)",
			caseSQL,
		)

		if err := r.db.WithContext(ctx).
			Exec(query, args...).Error; err != nil {
			return err
		}
	}

	return nil
}
