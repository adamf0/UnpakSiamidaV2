package infrastructure

import (
	"context"
	commondomainstandarrenstra "UnpakSiamida/common/domain"
	domainstandarrenstra "UnpakSiamida/modules/standarrenstra/domain"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"strings"
	"fmt"
)

type StandarRenstraRepository struct {
	db *gorm.DB
}

func NewStandarRenstraRepository(db *gorm.DB) domainstandarrenstra.IStandarRenstraRepository {
	return &StandarRenstraRepository{db: db}
}

// ------------------------
// GET BY UUID
// ------------------------
func (r *StandarRenstraRepository) GetByUuid(ctx context.Context, uid uuid.UUID) (*domainstandarrenstra.StandarRenstra, error) {
	var standarrenstra domainstandarrenstra.StandarRenstra

	err := r.db.WithContext(ctx).
		Where("uuid = ?", uid).
		First(&standarrenstra).Error

	// if errors.Is(err, gorm.ErrRecordNotFound) {
	// 	return nil, nil
	// }

	if err != nil {
		return nil, err
	}

	return &standarrenstra, nil
}

var allowedSearchColumns = map[string]string{
    // key:param -> db column
    "nama":          "nama",
}

// ------------------------
// GET ALL
// ------------------------
func (r *StandarRenstraRepository) GetAll(
	ctx context.Context,
	search string,
	searchFilters []commondomainstandarrenstra.SearchFilter,
	page, limit *int,
) ([]domainstandarrenstra.StandarRenstra, int64, error) {

	var standarrenstras []domainstandarrenstra.StandarRenstra
	var total int64

	db := r.db.WithContext(ctx).Model(&domainstandarrenstra.StandarRenstra{})

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
	if err := db.Order("created_at DESC").Find(&standarrenstras).Error; err != nil {
		return nil, 0, err
	}

	return standarrenstras, total, nil
}


// ------------------------
// CREATE
// ------------------------
func (r *StandarRenstraRepository) Create(ctx context.Context, standarrenstra *domainstandarrenstra.StandarRenstra) error {
	return r.db.WithContext(ctx).Create(standarrenstra).Error
}

// ------------------------
// UPDATE
// ------------------------
func (r *StandarRenstraRepository) Update(ctx context.Context, standarrenstra *domainstandarrenstra.StandarRenstra) error {
	return r.db.WithContext(ctx).Save(standarrenstra).Error
}

// ------------------------
// DELETE (by UUID)
// ------------------------
func (r *StandarRenstraRepository) Delete(ctx context.Context, uid uuid.UUID) error {
	return r.db.WithContext(ctx).
		Where("uuid = ?", uid).
		Delete(&domainstandarrenstra.StandarRenstra{}).Error
}
