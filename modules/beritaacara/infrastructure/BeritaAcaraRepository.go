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
	"gorm.io/gorm/clause"
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

	var rowData domainberitaacara.BeritaAcaraDefault

	err := r.db.WithContext(ctx).
		Table("berita_acara ba").
		Select(`
		ba.id as Id,
		ba.uuid as Uuid,
		ba.tahun as Tahun,
		ba.fakultas_unit as FakultasUnitId,
		fu.nama_fak_prod_unit as FakultasUnit,
		ba.tanggal as Tanggal,
		
		ba.auditee as AuditeeId,
		u1.uuid as AuditeeUuid,
		u1.name as Auditee,

		ba.auditor1 as Auditor1Id,
		u2.uuid as Auditor1Uuid,
		u2.name as Auditor1,
		
		ba.auditor2 as Auditor2Id,
		u3.uuid as Auditor2Uuid,
		u3.name as Auditor2
	`).
		Joins("LEFT JOIN v_fakultas_unit fu ON ba.fakultas_unit = fu.id").
		Joins("LEFT JOIN users u1 ON ba.auditee = u1.id").
		Joins("LEFT JOIN users u2 ON ba.auditor1 = u2.id").
		Joins("LEFT JOIN users u3 ON ba.auditor2 = u3.id").
		Where("ba.uuid = ?", id).
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
	"nama_fak_prod_unit": "fu.nama_fak_prod_unit",
	"auditee":            "u1.name",
	"auditor1":           "u2.name",
	"auditor2":           "u3.name",
	"fakultas":           "fu.nama_fak_prod_unit",
	"target":             "fu.uuid",
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

	var rows = make([]domainberitaacara.BeritaAcaraDefault, 0)
	var total int64

	db := r.db.WithContext(ctx).
		Table("berita_acara ba").
		Select(`
			ba.id as Id,
			ba.uuid as Uuid,
			ba.tahun as Tahun,
			ba.fakultas_unit as FakultasUnitId,
			fu.uuid as FakultasUnitUuid,
			fu.nama_fak_prod_unit as FakultasUnit,
			ba.tanggal as Tanggal,

			ba.auditee as AuditeeId,
			u1.uuid as AuditeeUuid,
			u1.name as Auditee,

			ba.auditor1 as Auditor1Id,
			u2.uuid as Auditor1Uuid,
			u2.name as Auditor1,

			ba.auditor2 as Auditor2Id,
			u3.uuid as Auditor2Uuid,
			u3.name as Auditor2
		`).
		Joins(`LEFT JOIN v_fakultas_unit fu ON ba.fakultas_unit = fu.id`).
		Joins(`LEFT JOIN users u1 ON ba.auditee = u1.id`).
		Joins(`LEFT JOIN users u2 ON ba.auditor1 = u2.id`).
		Joins(`LEFT JOIN users u3 ON ba.auditor2 = u3.id`)

	// ==========================================
	// ADVANCED FILTERS
	// ==========================================
	for _, f := range searchFilters {

		targetColumn := strings.ToLower(strings.TrimSpace(f.Field))
		if targetColumn == "" {
			continue
		}

		val := ""
		if f.Value != nil {
			val = strings.TrimSpace(*f.Value)
		}
		if val == "" {
			continue
		}

		// Special case: pic (search across 3 users)
		if targetColumn == "pic" {
			db = db.Where("(u1.uuid = ? OR u2.uuid = ? OR u3.uuid = ?)", val, val, val)
			continue
		}

		col, ok := allowedSearchColumns[targetColumn]
		if !ok || strings.TrimSpace(col) == "" {
			continue
		}

		switch strings.ToLower(f.Operator) {
		case "eq":
			db = db.Where(clause.Eq{
				Column: col,
				Value:  val,
			})
		case "neq":
			db = db.Where(clause.Neq{
				Column: col,
				Value:  val,
			})
		case "like":
			db = db.Where(clause.Like{
				Column: col,
				Value:  "%" + val + "%",
			})
		case "in":
			rawVals := strings.Split(val, ",")
			vals := make([]interface{}, 0, len(rawVals))

			for _, v := range rawVals {
				v = strings.TrimSpace(v)
				if v != "" {
					vals = append(vals, v)
				}
			}

			if len(vals) > 0 {
				db = db.Where(clause.IN{
					Column: col,
					Values: vals,
				})
			}
		}
	}

	// ==========================================
	// GLOBAL SEARCH (FIXED)
	// ==========================================
	if strings.TrimSpace(search) != "" {
		like := "%" + strings.TrimSpace(search) + "%"
		var conditions []clause.Expression

		for _, col := range allowedSearchColumns {

			// skip empty column (anti bug OR  LIKE ...)
			if strings.TrimSpace(col) == "" {
				continue
			}

			conditions = append(conditions, clause.Like{
				Column: col,
				Value:  like,
			})
		}

		if len(conditions) > 0 {
			db = db.Where(clause.Or(conditions...))
		}
	}

	// ==========================================
	// COUNT (SAFE)
	// ==========================================
	countDB := db.Session(&gorm.Session{})
	if err := countDB.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// ==========================================
	// ORDER + PAGINATION
	// ==========================================
	db = db.Order("ba.id DESC")

	if page != nil && limit != nil && *limit > 0 {
		offset := (*page - 1) * (*limit)
		db = db.Offset(offset).Limit(*limit)
	}

	// ==========================================
	// EXECUTE
	// ==========================================
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
