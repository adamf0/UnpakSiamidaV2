package infrastructure

import (
	commondomain "UnpakSiamida/common/domain"
	domaindokumentambahan "UnpakSiamida/modules/dokumentambahan/domain"
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"gorm.io/gorm"
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

	var rowData domaindokumentambahan.DokumenTambahanDefault

	err := r.db.WithContext(ctx).
		Table("dokumen_tambahan dt").
		Select(`
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
		dt.file as Link,
		dt.capaian_auditor as CapaianAuditor,
		dt.catatan_auditor as CatatanAuditor
	`).
		Joins("JOIN template_dokumen_tambahan tdt ON dt.id_template_dokumen_tambahan = tdt.id").
		Joins("JOIN jenis_file_renstra jf ON tdt.jenis_file = jf.id").
		Joins("JOIN renstra r ON dt.id_renstra = r.id").
		Joins("JOIN v_fakultas_unit fu ON r.fakultas_unit = fu.id").
		Where("dt.uuid = ?", id).
		Take(&rowData).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, err
	}

	return &rowData, nil
}

var allowedSearchColumns = map[string]string{
	// key:param -> db column
	"uuidrenstra": "r.uuid",
	"dokumen":     "jf.nama",
	"pertanyaan":  "tdt.pertanyaan",
	"tahun":       "tdt.tahun",
	"targetaudit": "fu.nama_fak_prod_unit",
	"jenjang":     "fu.jenjang",
	"tipe":        "fu.type",
}

func (r *DokumenTambahanRepository) GetAll(
	ctx context.Context,
	search string,
	searchFilters []commondomain.SearchFilter,
	page, limit *int,
) ([]domaindokumentambahan.DokumenTambahanDefault, int64, error) {

	var (
		result = make([]domaindokumentambahan.DokumenTambahanDefault, 0)
		total  int64
		args   []interface{}
		db     = r.db.WithContext(ctx).Table("dokumen_tambahan dt").
			Joins("JOIN template_dokumen_tambahan tdt ON dt.id_template_dokumen_tambahan = tdt.id").
			Joins("JOIN jenis_file_renstra jf ON tdt.jenis_file = jf.id").
			Joins("JOIN renstra r ON dt.id_renstra = r.id").
			Joins("JOIN v_fakultas_unit fu ON r.fakultas_unit = fu.id")
	)

	// =====================================================
	// ADVANCED FILTERS
	// =====================================================
	for _, f := range searchFilters {
		field := strings.TrimSpace(strings.ToLower(f.Field))
		operator := strings.TrimSpace(strings.ToLower(f.Operator))
		col, ok := allowedSearchColumns[field]
		if !ok || f.Value == nil {
			continue
		}

		value := strings.TrimSpace(*f.Value)

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
			values := strings.Split(value, ",")
			db = db.Where(fmt.Sprintf("%s IN ?", col), values)
		default:
			db = db.Where(fmt.Sprintf("%s LIKE ?", col), "%"+value+"%")
		}
	}

	// =====================================================
	// GLOBAL SEARCH (opsional)
	// =====================================================
	if strings.TrimSpace(search) != "" {
		like := "%" + search + "%"
		var ors []string
		for _, col := range allowedSearchColumns {
			ors = append(ors, fmt.Sprintf("%s LIKE ?", col))
			args = append(args, like)
		}
		db = db.Where("("+strings.Join(ors, " OR ")+")", args...)
	}

	// =====================================================
	// COUNT TOTAL
	// =====================================================
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// =====================================================
	// PAGINATION
	// =====================================================
	if page != nil && limit != nil && *limit > 0 {
		p := *page
		l := *limit
		if p < 1 {
			p = 1
		}
		offset := (p - 1) * l
		db = db.Limit(l).Offset(offset)
	}

	// =====================================================
	// SELECT
	// =====================================================
	db = db.Select(`
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
		dt.file as Link,
		dt.capaian_auditor as CapaianAuditor,
		dt.catatan_auditor as CatatanAuditor
	`).Order("r.tahun DESC, dt.id DESC")

	if err := db.Scan(&result).Error; err != nil {
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
