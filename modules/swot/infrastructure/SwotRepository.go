package infrastructure

import (
	commondomain "UnpakSiamida/common/domain"
	domainSwot "UnpakSiamida/modules/swot/domain"
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SwotRepository struct {
	db *gorm.DB
}

func NewSwotRepository(db *gorm.DB) domainSwot.ISwotRepository {
	return &SwotRepository{db: db}
}

// ------------------------
// GET BY UUID
// ------------------------
func (r *SwotRepository) GetByUuid(ctx context.Context, uid uuid.UUID) (*domainSwot.Swot, error) {
	var item domainSwot.Swot
	err := r.db.WithContext(ctx).Table("swot").Where("uuid = ?", uid).First(&item).Error
	if err != nil {
		return nil, err
	}
	return &item, nil
}

// ------------------------
// GET DEFAULT BY UUID
// ------------------------
func (r *SwotRepository) GetDefaultByUuid(ctx context.Context, uid uuid.UUID) (*domainSwot.SwotDefault, error) {
	var row domainSwot.SwotDefault
	err := r.db.WithContext(ctx).
		Table("swot s").
		Select(`
			s.id as Id,
			s.uuid as UUID,
			s.id_fakultas_unit as IdFakultasUnit,
			fu.nama_fak_prod_unit as NamaFakultasUnit,
			fu.uuid as UuidFakultasUnit,
			s.id_master_tahun as IdMasterTahun,
			COALESCE(mt.tahun, s.tahun) as Tahun,
			mt.uuid as TahunUuid,
			s.kategori as Kategori,
			s.faktor as Faktor,
			s.deskripsi as Deskripsi,
			s.strategi as Strategi,
			s.bobot as Bobot,
			s.rating as Rating,
			s.skor as Skor,
			s.created_at as CreatedAt,
			s.updated_at as UpdatedAt
		`).
		Joins("LEFT JOIN v_fakultas_unit fu ON s.id_fakultas_unit = fu.id").
		Joins("LEFT JOIN master_tahun mt ON s.id_master_tahun = mt.id").
		Where("s.uuid = ?", uid).
		Take(&row).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, err
	}

	return &row, nil
}

var allowedSearchColumns = map[string]string{
	"tahun":            "COALESCE(mt.tahun, s.tahun)",
	"tahunuuid":        "mt.uuid",
	"kategori":         "s.kategori",
	"faktor":           "s.faktor",
	"deskripsi":        "s.deskripsi",
	"strategi":         "s.strategi",
	"id_fakultas_unit": "s.id_fakultas_unit",
	"fakultas_unit":    "fu.nama_fak_prod_unit",
	"fakultasunituuid": "fu.uuid",
}

// ------------------------
// GET ALL
// ------------------------
func (r *SwotRepository) GetAll(
	ctx context.Context,
	search string,
	searchFilters []commondomain.SearchFilter,
	page, limit *int,
) ([]domainSwot.SwotDefault, int64, error) {

	var swotList = make([]domainSwot.SwotDefault, 0)
	var total int64

	db := r.db.WithContext(ctx).
		Table("swot s").
		Select(`
			s.id as Id,
			s.uuid as UUID,
			s.id_fakultas_unit as IdFakultasUnit,
			fu.nama_fak_prod_unit as NamaFakultasUnit,
			fu.uuid as UuidFakultasUnit,
			s.id_master_tahun as IdMasterTahun,
			COALESCE(mt.tahun, s.tahun) as Tahun,
			mt.uuid as TahunUuid,
			s.kategori as Kategori,
			s.faktor as Faktor,
			s.deskripsi as Deskripsi,
			s.strategi as Strategi,
			s.bobot as Bobot,
			s.rating as Rating,
			s.skor as Skor,
			s.created_at as CreatedAt,
			s.updated_at as UpdatedAt
		`).
		Joins("LEFT JOIN v_fakultas_unit fu ON s.id_fakultas_unit = fu.id").
		Joins("LEFT JOIN master_tahun mt ON s.id_master_tahun = mt.id")

	// -------------------------------
	// SEARCH FILTERS
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

	// -------------------------------
	// GLOBAL SEARCH
	// -------------------------------
	if strings.TrimSpace(search) != "" {
		like := "%" + search + "%"
		var orParts []string
		var params []interface{}

		for _, col := range []string{"s.kategori", "s.faktor", "s.deskripsi", "s.strategi", "COALESCE(mt.tahun, s.tahun)", "fu.nama_fak_prod_unit"} {
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

	db = db.Order("s.id DESC")

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

	if err := db.Find(&swotList).Error; err != nil {
		return nil, 0, err
	}

	return swotList, total, nil
}

// ------------------------
// CREATE
// ------------------------
func (r *SwotRepository) Create(ctx context.Context, data *domainSwot.Swot) error {
	return r.db.WithContext(ctx).Table("swot").Create(data).Error
}

// ------------------------
// UPDATE
// ------------------------
func (r *SwotRepository) Update(ctx context.Context, data *domainSwot.Swot) error {
	return r.db.WithContext(ctx).Table("swot").Save(data).Error
}

// ------------------------
// DELETE
// ------------------------
func (r *SwotRepository) Delete(ctx context.Context, uid uuid.UUID) error {
	return r.db.WithContext(ctx).Table("swot").Where("uuid = ?", uid).Delete(&domainSwot.Swot{}).Error
}
