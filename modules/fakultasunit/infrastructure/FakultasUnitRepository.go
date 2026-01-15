package infrastructure

import (
	commondomainfakultasunit "UnpakSiamida/common/domain"
	domainfakultasunit "UnpakSiamida/modules/fakultasunit/domain"
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type FakultasUnitRepository struct {
	db *gorm.DB
}

func NewFakultasUnitRepository(db *gorm.DB) domainfakultasunit.IFakultasUnitRepository {
	return &FakultasUnitRepository{db: db}
}

// ------------------------
// GET BY UUID
// ------------------------
// [PR] ini bukan default
func (r *FakultasUnitRepository) GetDefaultByUuid(ctx context.Context, uid uuid.UUID) (*domainfakultasunit.FakultasUnit, error) {
	var fakultasunit domainfakultasunit.FakultasUnit

	err := r.db.WithContext(ctx).
		Where("uuid = ?", uid).
		First(&fakultasunit).Error

	// if errors.Is(err, gorm.ErrRecordNotFound) {
	// 	return nil, nil
	// }

	if err != nil {
		return nil, err
	}

	return &fakultasunit, nil
}

var allowedSearchColumns = map[string]string{
	// key:param -> db column
	"nama_fak_prod_unit": "nama_fak_prod_unit",
	"jenjang":            "jenjang",
	"fakultas":           "fakultas",
}

// ------------------------
// GET ALL
// ------------------------
func (r *FakultasUnitRepository) GetAll(
	ctx context.Context,
	search string,
	searchFilters []commondomainfakultasunit.SearchFilter,
	page, limit *int,
) ([]domainfakultasunit.FakultasUnit, int64, error) {

	var fakultasunits []domainfakultasunit.FakultasUnit
	var total int64

	db := r.db.WithContext(ctx).Model(&domainfakultasunit.FakultasUnit{})

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

		db = db.Where("("+strings.Join(orParts, " OR ")+")", params...)
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
	if err := db.Find(&fakultasunits).Error; err != nil {
		return nil, 0, err
	}

	return fakultasunits, total, nil
}

func (r *FakultasUnitRepository) SetupUuid(ctx context.Context) error {
	const chunkSize = 500

	var ids []uint
	if err := r.db.WithContext(ctx).
		Model(&domainfakultasunit.FakultasUnit{}).
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
			"UPDATE sijamu_fakultas_unit SET uuid = %s WHERE id IN (?)",
			caseSQL,
		)

		if err := r.db.WithContext(ctx).
			Exec(query, args...).Error; err != nil {
			return err
		}
	}

	return nil
}
