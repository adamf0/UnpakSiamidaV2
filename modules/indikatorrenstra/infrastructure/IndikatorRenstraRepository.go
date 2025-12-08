package infrastructure

import (
	"context"
	commondomainindikatorrenstra "UnpakSiamida/common/domain"
	domainindikatorrenstra "UnpakSiamida/modules/indikatorrenstra/domain"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"strings"
	"fmt"
)

type IndikatorRenstraRepository struct {
	db *gorm.DB
}

func NewIndikatorRenstraRepository(db *gorm.DB) domainindikatorrenstra.IIndikatorRenstraRepository {
	return &IndikatorRenstraRepository{db: db}
}

func (r *IndikatorRenstraRepository) IsUniqueIndikator(
	ctx context.Context,
	indikator string,
	tahun string,
) (bool, error) {

	var count int64

	err := r.db.WithContext(ctx).
		Model(&domainindikatorrenstra.IndikatorRenstra{}).
		Where("indikator = ?", indikator).
		Where("tahun = ?", tahun).
		Count(&count).Error

	if err != nil {
		return false, err
	}

	return count == 0, nil
}

// ------------------------
// GET BY UUID
// ------------------------
func (r *IndikatorRenstraRepository) GetByUuid(ctx context.Context, uid uuid.UUID) (*domainindikatorrenstra.IndikatorRenstra, error) {
	var indikatorrenstra domainindikatorrenstra.IndikatorRenstra

	err := r.db.WithContext(ctx).
		Where("uuid = ?", uid).
		First(&indikatorrenstra).Error

	// if errors.Is(err, gorm.ErrRecordNotFound) {
	// 	return nil, nil
	// }

	if err != nil {
		return nil, err
	}

	return &indikatorrenstra, nil
}

var allowedSearchColumns = map[string]string{
    // key:param -> db column
    "nama":           	"nama",
	"standar_renstra": 	"id_master_standar",
	"indikator": 	  	"indikator",
	"parent":		  	"parent",
	"tahun": 		  	"tahun",
	"tipe_target": 	  	"tipe_target",
	"operator": 	  	"operator",
}

// ------------------------
// GET ALL
// ------------------------
func (r *IndikatorRenstraRepository) GetAll(
	ctx context.Context,
	search string,
	searchFilters []commondomainindikatorrenstra.SearchFilter,
	page, limit *int,
) ([]domainindikatorrenstra.IndikatorRenstra, int64, error) {

	var indikatorrenstras []domainindikatorrenstra.IndikatorRenstra
	var total int64

	db := r.db.WithContext(ctx).Model(&domainindikatorrenstra.IndikatorRenstra{})

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
	if err := db.Order("created_at DESC").Find(&indikatorrenstras).Error; err != nil {
		return nil, 0, err
	}

	return indikatorrenstras, total, nil
}


// ------------------------
// CREATE
// ------------------------
func (r *IndikatorRenstraRepository) Create(ctx context.Context, indikatorrenstra *domainindikatorrenstra.IndikatorRenstra) error {
	return r.db.WithContext(ctx).Create(indikatorrenstra).Error
}

// ------------------------
// UPDATE
// ------------------------
func (r *IndikatorRenstraRepository) Update(ctx context.Context, indikatorrenstra *domainindikatorrenstra.IndikatorRenstra) error {
	return r.db.WithContext(ctx).Save(indikatorrenstra).Error
}

// ------------------------
// DELETE (by UUID)
// ------------------------
func (r *IndikatorRenstraRepository) Delete(ctx context.Context, uid uuid.UUID) error {
	return r.db.WithContext(ctx).
		Where("uuid = ?", uid).
		Delete(&domainindikatorrenstra.IndikatorRenstra{}).Error
}
