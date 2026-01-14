package infrastructure

import (
	"context"
	commondomain "UnpakSiamida/common/domain"
	domaindokumentambahan "UnpakSiamida/modules/dokumentambahan/domain"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"strings"
	"fmt"
	"errors"
)

type DokumenTambahanRepository struct {
	db *gorm.DB
}

func NewDokumenTambahanRepository(db *gorm.DB) domaindokumentambahan.IDokumenTambahanRepository {
	return &DokumenTambahanRepository{db: db}
}

// ------------------------
// GET BY UUID
// ------------------------
func (r *DokumenTambahanRepository) GetByUuid(ctx context.Context, uid uuid.UUID) (*domaindokumentambahan.DokumenTambahan, error) {
	var dokumentambahan domaindokumentambahan.DokumenTambahan

	err := r.db.WithContext(ctx).
		Where("uuid = ?", uid).
		First(&dokumentambahan).Error

	// if errors.Is(err, gorm.ErrRecordNotFound) {
	// 	return nil, nil
	// }

	if err != nil {
		return nil, err
	}

	return &dokumentambahan, nil
}

// ------------------------
// GET DEFAULT BY UUID
// ------------------------
func (r *DokumenTambahanRepository) GetDefaultByUuid(
    ctx context.Context,
    id uuid.UUID,
) (*domaindokumentambahan.DokumenTambahanDefault, error) {

    query := `
        SELECT 
			dt.id as ID,
			dt.uuid as UUID,
			r.id as RenstraId,
			r.uuid as RenstraUUID,
			tdt.uuid as TemplateDokumenTambahanUUID,
			fu.nama_fak_prod_unit as TargetAudit,
            fu.jenjang as Jenjang,
            fu.fakultas as Fakultas,
            fu.type as Type,

			tdt.pertanyaan as Pertanyaan,
			jf.nama as Dokumen,
			tdt.klasifikasi as Klasifikasi,

			r.tahun as TahunRenstra,
			tdt.tahun as TahunDokumenTambahan,
			tdt.tugas as Tugas,
			file as Link,
			dt.capaian_auditor as CapaianAuditor,
			dt.catatan_auditor as CatatanAuditor 
		FROM dokumen_tambahan dt
		join template_dokumen_tambahan tdt on dt.id_template_dokumen_tambahan = tdt.id
		join jenis_file_renstra jf on tdt.jenis_file = jf.id 
		join renstra r on dt.id_renstra = r.id
		join v_fakultas_unit fu on r.fakultas_unit = fu.id
		where r.id = = ?
        LIMIT 1
    `

    var rowData domaindokumentambahan.DokumenTambahanDefault

    res := r.db.WithContext(ctx).Raw(query, id).Scan(&rowData)
    if res.Error != nil {
        if errors.Is(res.Error, gorm.ErrRecordNotFound) {
            return nil, gorm.ErrRecordNotFound
        }
        return nil, res.Error
    }

    if rowData.ID == 0 {
        return nil, gorm.ErrRecordNotFound
    }

    return &rowData, nil
}

var allowedSearchColumns = map[string]string{
    // key:param -> db column
	"uuidrenstra":      "r.uuid",
	"dokumen":        	"jf.nama",
	"pertanyaan":   	"tdt.pertanyaan",
	"tahun":           	"tdt.tahun",
	"targetaudit":      "fu.nama_fak_prod_unit",
	"jenjang":      	"fu.jenjang",
	"tipe":      		"fu.type",
}

func (r *DokumenTambahanRepository) GetAll(
	ctx context.Context,
	search string,
	searchFilters []commondomain.SearchFilter,
	page, limit *int,
) ([]domaindokumentambahan.DokumenTambahanDefault, int64, error) {

	var (
		result     []domaindokumentambahan.DokumenTambahanDefault
		total      int64
		conditions []string
		args       []interface{}
	)

	// =====================================================
	// BASE FROM + JOIN (WAJIB PAKAI INI)
	// =====================================================
	baseFrom := `
		FROM dokumen_tambahan dt
		join template_dokumen_tambahan tdt on dt.id_template_dokumen_tambahan = tdt.id
		join jenis_file_renstra jf on tdt.jenis_file = jf.id 
		join renstra r on dt.id_renstra = r.id
		join v_fakultas_unit fu on r.fakultas_unit = fu.id
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

	orderBy := " ORDER BY r.tahun DESC, dt.id DESC"

	// =====================================================
	// SELECT QUERY (INI YANG KAMU MINTA)
	// =====================================================
	selectQuery := `
		SELECT
			dt.id as ID,
			dt.uuid as UUID,
			r.id as RenstraId,
			r.uuid as RenstraUUID,
			tdt.uuid as TemplateDokumenTambahanUUID,
			fu.nama_fak_prod_unit as TargetAudit,
            fu.jenjang as Jenjang,
            fu.fakultas as Fakultas,
            fu.type as Type,

			tdt.pertanyaan as Pertanyaan,
			jf.nama as Dokumen,
			tdt.klasifikasi as Klasifikasi,

			r.tahun as TahunRenstra,
			tdt.tahun as TahunDokumenTambahan,
			tdt.tugas as Tugas,
			file as Link,
			dt.capaian_auditor as CapaianAuditor,
			dt.catatan_auditor as CatatanAuditor 
	` + baseFrom + whereClause + orderBy + pagination

	if err := r.db.WithContext(ctx).
		Raw(selectQuery, args...).
		Scan(&result).Error; err != nil {
		return nil, 0, err
	}

	return result, total, nil
}


// ------------------------
// UPDATE
// ------------------------
func (r *DokumenTambahanRepository) Update(ctx context.Context, dokumentambahan *domaindokumentambahan.DokumenTambahan) error {
	return r.db.WithContext(ctx).Save(dokumentambahan).Error
}

// ------------------------
// DELETE (by UUID)
// ------------------------
func (r *DokumenTambahanRepository) Delete(ctx context.Context, uid uuid.UUID) error {
	return r.db.WithContext(ctx).
		Where("uuid = ?", uid).
		Delete(&domaindokumentambahan.DokumenTambahan{}).Error
}

func (r *DokumenTambahanRepository) SetupUuid(ctx context.Context) error {
	const chunkSize = 500

	var ids []uint
	if err := r.db.WithContext(ctx).
		Model(&domaindokumentambahan.DokumenTambahan{}).
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
			"UPDATE dokumen_tambahan SET uuid = %s WHERE id IN (?)",
			caseSQL,
		)

		if err := r.db.WithContext(ctx).
			Exec(query, args...).Error; err != nil {
			return err
		}
	}

	return nil
}