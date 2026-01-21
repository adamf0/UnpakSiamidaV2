package infrastructure

import (
	commondomainberitaacara "UnpakSiamida/common/domain"
	domainberitaacara "UnpakSiamida/modules/beritaacara/domain"
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BeritaAcaraRepository struct {
	db *gorm.DB
}

func NewBeritaAcaraRepository(db *gorm.DB) domainberitaacara.IBeritaAcaraRepository {
	return &BeritaAcaraRepository{db: db}
}

// ------------------------
// GET BY UUID
// ------------------------
func (r *BeritaAcaraRepository) GetByUuid(ctx context.Context, uid uuid.UUID) (*domainberitaacara.BeritaAcara, error) {
	var BeritaAcara domainberitaacara.BeritaAcara

	err := r.db.WithContext(ctx).
		Where("uuid = ?", uid).
		First(&BeritaAcara).Error

	// if errors.Is(err, gorm.ErrRecordNotFound) {
	// 	return nil, nil
	// }

	if err != nil {
		return nil, err
	}

	return &BeritaAcara, nil
}

// ------------------------
// GET DEFAULT BY UUID
// ------------------------
func (r *BeritaAcaraRepository) GetDefaultByUuid(
	ctx context.Context,
	id uuid.UUID,
) (*domainberitaacara.BeritaAcaraDefault, error) {

	query := `
        SELECT 
            ba.id as Id,
            ba.uuid as Uuid,
            
            ba.tahun as Tahun,
            ba.fakultas_unit as FakultasUnitId,
            fu.nama_fak_prod_unit as FakultasUnit,
            ba.tanggal as Tanggal,
            
            ba.auditee as AuditeeId,
            u1.name as NamaAuditee,
            
            ba.auditor1 as Auditor1,
            u2.name as NamaAuditor1,
            
            ba.auditor2 as Auditor2,
            u3.name as NamaAuditor2
        FROM berita_acara ba
        LEFT JOIN v_fakultas_unit fu ON ba.fakultas_unit = fu.id
        LEFT JOIN users u1 ON ba.auditee = u1.id
        LEFT JOIN users u2 ON ba.auditor1 = u2.id 
        LEFT JOIN users u3 ON ba.auditor2 = u3.id
        WHERE ba.uuid = ?
        LIMIT 1
    `

	var rowData domainberitaacara.BeritaAcaraDefault

	res := r.db.WithContext(ctx).Raw(query, id).Scan(&rowData)
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, res.Error
	}

	if rowData.Id == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	return &rowData, nil
}

var allowedSearchColumns = map[string]string{
	// key:param -> db column
	"nama_fak_prod_unit": "fu.nama_fak_prod_unit",
}

// ------------------------
// GET ALL
// ------------------------
func (r *BeritaAcaraRepository) GetAll(
	ctx context.Context,
	search string,
	searchFilters []commondomainberitaacara.SearchFilter,
	page, limit *int,
) ([]domainberitaacara.BeritaAcaraDefault, int64, error) {

	var rows []domainberitaacara.BeritaAcaraDefault
	var total int64

	db := r.db.WithContext(ctx).
		Table("berita_acara ba").
		Select(`
			ba.id AS id,
			ba.uuid AS uuid,

			ba.tahun AS tahun,
			ba.fakultas_unit AS fakultas_unit_id,
			fu.nama_fak_prod_unit AS fakultas_unit,
			ba.tanggal AS tanggal,

			ba.auditee AS auditee_id,
			u1.name AS nama_auditee,

			ba.auditor1 AS auditor1,
			u2.name AS nama_auditor1,

			ba.auditor2 AS auditor2,
			u3.name AS nama_auditor2
		`).
		Joins(`LEFT JOIN v_fakultas_unit fu ON ba.fakultas_unit = fu.id`).
		Joins(`LEFT JOIN users u1 ON ba.auditee = u1.id`).
		Joins(`LEFT JOIN users u2 ON ba.auditor1 = u2.id`).
		Joins(`LEFT JOIN users u3 ON ba.auditor2 = u3.id`)

	// -----------------------------------
	// ADVANCED FILTERS
	// -----------------------------------
	for _, f := range searchFilters {
		col, ok := allowedSearchColumns[strings.ToLower(f.Field)]
		if !ok {
			continue
		}

		val := ""
		if f.Value != nil {
			val = strings.TrimSpace(*f.Value)
		}

		switch strings.ToLower(f.Operator) {
		case "eq":
			db = db.Where(col+" = ?", val)
		case "neq":
			db = db.Where(col+" <> ?", val)
		case "like":
			db = db.Where(col+" LIKE ?", "%"+val+"%")
		case "in":
			db = db.Where(col+" IN ?", strings.Split(val, ","))
		}
	}

	// -----------------------------------
	// GLOBAL SEARCH
	// -----------------------------------
	if strings.TrimSpace(search) != "" {
		like := "%" + search + "%"
		var or []string
		var args []interface{}

		for _, col := range allowedSearchColumns {
			or = append(or, col+" LIKE ?")
			args = append(args, like)
		}

		db = db.Where("("+strings.Join(or, " OR ")+")", args...)
	}

	// -----------------------------------
	// COUNT (AMAN)
	// -----------------------------------
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// -----------------------------------
	// ORDER + PAGINATION
	// -----------------------------------
	db = db.Order("ba.id DESC")

	if page != nil && limit != nil && *limit > 0 {
		offset := (*page - 1) * (*limit)
		db = db.Offset(offset).Limit(*limit)
	}

	// -----------------------------------
	// EXECUTE
	// -----------------------------------
	if err := db.Find(&rows).Error; err != nil {
		return nil, 0, err
	}

	return rows, total, nil
}

// ------------------------
// CREATE
// ------------------------
func (r *BeritaAcaraRepository) Create(ctx context.Context, jenisfile *domainberitaacara.BeritaAcara) error {
	return r.db.WithContext(ctx).Create(jenisfile).Error
}

// ------------------------
// UPDATE
// ------------------------
func (r *BeritaAcaraRepository) Update(ctx context.Context, jenisfile *domainberitaacara.BeritaAcara) error {
	return r.db.WithContext(ctx).Save(jenisfile).Error
}

// ------------------------
// DELETE (by UUID)
// ------------------------
func (r *BeritaAcaraRepository) Delete(ctx context.Context, uid uuid.UUID) error {
	return r.db.WithContext(ctx).
		Where("uuid = ?", uid).
		Delete(&domainberitaacara.BeritaAcara{}).Error
}

func (r *BeritaAcaraRepository) SetupUuid(ctx context.Context) error {
	const chunkSize = 500

	var ids []uint
	if err := r.db.WithContext(ctx).
		Model(&domainberitaacara.BeritaAcara{}).
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
			"UPDATE berita_acara SET uuid = %s WHERE id IN (?)",
			caseSQL,
		)

		if err := r.db.WithContext(ctx).
			Exec(query, args...).Error; err != nil {
			return err
		}
	}

	return nil
}
