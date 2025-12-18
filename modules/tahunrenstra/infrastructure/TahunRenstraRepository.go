package infrastructure

import (
	"context"
	commondomainTahunRenstra "UnpakSiamida/common/domain"
	domainTahunRenstra "UnpakSiamida/modules/tahunrenstra/domain"
	"gorm.io/gorm"
	"strings"
	"fmt"
	"strconv"
	"time"
)

type TahunRenstraRepository struct {
	db *gorm.DB
}

func NewTahunRenstraRepository(db *gorm.DB) domainTahunRenstra.ITahunRenstraRepository {
	return &TahunRenstraRepository{db: db}
}

// ------------------------
// GET ACTIVE
// ------------------------
func (r *TahunRenstraRepository) GetActive(ctx context.Context) (*domainTahunRenstra.TahunRenstra, error) {
	var TahunRenstra domainTahunRenstra.TahunRenstra

	tahunNow := strconv.Itoa(time.Now().Year())

	err := r.db.WithContext(ctx).
		Where("tahun = ?", tahunNow).
		Order("tahun DESC").
		First(&TahunRenstra).Error

	// if errors.Is(err, gorm.ErrRecordNotFound) {
	// 	return nil, nil
	// }

	if err != nil {
		return nil, err
	}

	return &TahunRenstra, nil
}

var allowedSearchColumns = map[string]string{
    // key:param -> db column
    "tahun":          "tahun",
	"status":         "status",
}

// ------------------------
// GET ALL
// ------------------------
func (r *TahunRenstraRepository) GetAll(
	ctx context.Context,
	search string,
	searchFilters []commondomainTahunRenstra.SearchFilter,
	page, limit *int,
) ([]domainTahunRenstra.TahunRenstra, int64, error) {

	var TahunRenstras []domainTahunRenstra.TahunRenstra
	var total int64

	db := r.db.WithContext(ctx).Model(&domainTahunRenstra.TahunRenstra{})
	db = db.Where("status IN ('active', 'no-active')")
    db = db.Where("tahun IS NOT NULL AND tahun != '0000'")

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
	if err := db.Order("tahun DESC").Find(&TahunRenstras).Error; err != nil {
		return nil, 0, err
	}

	return TahunRenstras, total, nil
}