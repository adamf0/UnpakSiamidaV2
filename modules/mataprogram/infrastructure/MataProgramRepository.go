package infrastructure

import (
	commondomainMataProgram "UnpakSiamida/common/domain"
	domainMataProgram "UnpakSiamida/modules/mataprogram/domain"
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type MataProgramRepository struct {
	db *gorm.DB
}

func NewMataProgramRepository(db *gorm.DB) domainMataProgram.IMataProgramRepository {
	return &MataProgramRepository{db: db}
}

// ------------------------
// GET BY UUID
// ------------------------
func (r *MataProgramRepository) GetByUuid(ctx context.Context, uid uuid.UUID) (*domainMataProgram.MataProgram, error) {
	var MataProgram domainMataProgram.MataProgram

	err := r.db.WithContext(ctx).
		Where("uuid = ?", uid).
		First(&MataProgram).Error

	// if errors.Is(err, gorm.ErrRecordNotFound) {
	// 	return nil, nil
	// }

	if err != nil {
		return nil, err
	}

	return &MataProgram, nil
}

// ------------------------
// GET DEFAULT BY UUID
// ------------------------
func (r *MataProgramRepository) GetDefaultByUuid(
	ctx context.Context,
	id uuid.UUID,
) (*domainMataProgram.MataProgramDefault, error) {

	// Ambil hanya kolom yang benar-benar ada di struct MataProgramDefault
	query := `
		SELECT id, uuid, nama
		FROM mata_program
		WHERE uuid = ?
		LIMIT 1
	`

	var rowData domainMataProgram.MataProgramDefault

	err := r.db.WithContext(ctx).Raw(query, id).Scan(&rowData).Error
	if err != nil {
		return nil, err
	}

	// Jika tidak ada row → struct kosong → anggap record not found
	if rowData.Id == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	return &rowData, nil
}

var allowedSearchColumns = map[string]string{
	// key:param -> db column
	"nama": "nama",
}

// ------------------------
// GET ALL
// ------------------------
func (r *MataProgramRepository) GetAll(
	ctx context.Context,
	search string,
	searchFilters []commondomainMataProgram.SearchFilter,
	page, limit *int,
) ([]domainMataProgram.MataProgram, int64, error) {

	var MataPrograms []domainMataProgram.MataProgram
	var total int64

	db := r.db.WithContext(ctx).Model(&domainMataProgram.MataProgram{})

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

	db = db.Order("id DESC")

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
	if err := db.Find(&MataPrograms).Error; err != nil {
		return nil, 0, err
	}

	return MataPrograms, total, nil
}

// ------------------------
// CREATE
// ------------------------
func (r *MataProgramRepository) Create(ctx context.Context, mataprogram *domainMataProgram.MataProgram) error {
	return r.db.WithContext(ctx).Create(mataprogram).Error
}

// ------------------------
// UPDATE
// ------------------------
func (r *MataProgramRepository) Update(ctx context.Context, mataprogram *domainMataProgram.MataProgram) error {
	return r.db.WithContext(ctx).Save(mataprogram).Error
}

// ------------------------
// DELETE (by UUID)
// ------------------------
func (r *MataProgramRepository) Delete(ctx context.Context, uid uuid.UUID) error {
	return r.db.WithContext(ctx).
		Where("uuid = ?", uid).
		Delete(&domainMataProgram.MataProgram{}).Error
}

func (r *MataProgramRepository) SetupUuid(ctx context.Context) error {
	const chunkSize = 500

	var ids []uint
	if err := r.db.WithContext(ctx).
		Model(&domainMataProgram.MataProgram{}).
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
			"UPDATE mata_program SET uuid = %s WHERE id IN (?)",
			caseSQL,
		)

		if err := r.db.WithContext(ctx).
			Exec(query, args...).Error; err != nil {
			return err
		}
	}

	return nil
}
