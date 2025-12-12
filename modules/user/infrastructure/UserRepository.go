package infrastructure

import (
	"context"
	commondomainuser "UnpakSiamida/common/domain"
	domainuser "UnpakSiamida/modules/user/domain"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"strings"
	"fmt"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) domainuser.IUserRepository {
	return &UserRepository{db: db}
}

// ------------------------
// GET BY UUID
// ------------------------
func (r *UserRepository) GetByUuid(ctx context.Context, uid uuid.UUID) (*domainuser.User, error) {
	var user domainuser.User

	err := r.db.WithContext(ctx).
		Where("uuid = ?", uid).
		First(&user).Error

	// if errors.Is(err, gorm.ErrRecordNotFound) {
	// 	return nil, domainuser.NotFound(uid.String())
	// }

	if err != nil {
		return nil, err
	}

	return &user, nil
}

var allowedSearchColumns = map[string]string{
    // key:param -> db column
    "name":          "name",
    "username":      "nidn_username",
    "email":         "email",
	"level":         "level",
    "fakultas_unit":  "fakultas_unit",
}

// ------------------------
// GET ALL
// ------------------------
func (r *UserRepository) GetAll(
	ctx context.Context,
	search string,
	searchFilters []commondomainuser.SearchFilter,
	page, limit *int,
) ([]domainuser.User, int64, error) {

	var users []domainuser.User
	var total int64

	db := r.db.WithContext(ctx).Model(&domainuser.User{})

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
	if err := db.Order("created_at DESC").Find(&users).Error; err != nil {
		return nil, 0, err
	}

	return users, total, nil
}


// ------------------------
// CREATE
// ------------------------
func (r *UserRepository) Create(ctx context.Context, user *domainuser.User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

// ------------------------
// UPDATE
// ------------------------
func (r *UserRepository) Update(ctx context.Context, user *domainuser.User) error {
	return r.db.WithContext(ctx).Save(user).Error
}

// ------------------------
// DELETE (by UUID)
// ------------------------
func (r *UserRepository) Delete(ctx context.Context, uid uuid.UUID) error {
	return r.db.WithContext(ctx).
		Where("uuid = ?", uid).
		Delete(&domainuser.User{}).Error
}
