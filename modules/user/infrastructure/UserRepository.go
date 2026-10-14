package infrastructure

import (
	commondomainuser "UnpakSiamida/common/domain"
	commoninfra "UnpakSiamida/common/infrastructure"
	domainuser "UnpakSiamida/modules/user/domain"
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRepository struct {
	db  *gorm.DB
	uow *commoninfra.UnitOfWork
}

func NewUserRepository(db *gorm.DB) domainuser.IUserRepository {
	return &UserRepository{db: db, uow: commoninfra.NewUnitOfWork(db)}
}

// ------------------------
// ------------------------
// GET BY UUID
// ------------------------
func (r *UserRepository) GetByUuid(ctx context.Context, uid uuid.UUID) (*domainuser.User, error) {
	var user domainuser.User

	err := r.db.WithContext(ctx).Table("users").
		Select(`
			users.id,
			users.uuid,
			users.nidn_username,
			users.password,
			users.name,
			users.email,
			users.level,
			users.fakultas_unit,
			users.id_user_merge,
			mu.name as nama_user_merge,
			fu.nama_fak_prod_unit,
			fu.uuid as uuid_fakultas_unit
		`).
		Joins("LEFT JOIN v_fakultas_unit fu ON users.fakultas_unit = fu.id").
		Joins("LEFT JOIN users mu ON users.id_user_merge = mu.id").
		Where("users.uuid = ?", uid).
		First(&user).Error

	if err != nil {
		return nil, err
	}

	return &user, nil
}

var allowedSearchColumns = map[string]string{
	// key:param -> db column
	"name":          "users.name",
	"username":      "users.nidn_username",
	"email":         "users.email",
	"level":         "users.level",
	"fakultas_unit": "fu.nama_fak_prod_unit",
	"id_user_merge": "users.id_user_merge",
}

func (r *UserRepository) GetAllStrict(ctx context.Context) ([]domainuser.UserOptions, error) {
	var users = make([]domainuser.UserOptions, 0)

	// Ambil semua data dari table "users"
	if err := r.db.WithContext(ctx).Table("users").Where("level != ?", "admin").Find(&users).Error; err != nil {
		return nil, err
	}

	return users, nil
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

	var users = make([]domainuser.User, 0)
	var total int64

	db := r.db.WithContext(ctx).Table("users").
		Select(`
			users.id,
			users.uuid,
			users.nidn_username,
			users.password,
			users.name,
			users.email,
			users.level,
			users.fakultas_unit,
			users.id_user_merge,
			mu.name as nama_user_merge,
			fu.nama_fak_prod_unit,
			fu.uuid as uuid_fakultas_unit
		`).
		Joins("LEFT JOIN v_fakultas_unit fu ON users.fakultas_unit = fu.id").
		Joins("LEFT JOIN users mu ON users.id_user_merge = mu.id")

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
	if err := db.Order("users.id DESC").Find(&users).Error; err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

// ------------------------
// CREATE
// ------------------------
func (r *UserRepository) Create(ctx context.Context, user *domainuser.User) error {
	if err := r.db.WithContext(ctx).Table("users").Omit("NamaFakultasUnit", "UuidFakultasUnit", "nama_fak_prod_unit", "uuid_fakultas_unit").Save(user).Error; err != nil {
		return err
	}

	return r.uow.Save(&user.Entity)
}

// ------------------------
// UPDATE
// ------------------------
func (r *UserRepository) Update(ctx context.Context, user *domainuser.User) error {
	if err := r.db.WithContext(ctx).Table("users").Omit("NamaFakultasUnit", "UuidFakultasUnit", "nama_fak_prod_unit", "uuid_fakultas_unit").Save(user).Error; err != nil {
		return err
	}

	return r.uow.Save(&user.Entity)
}

// ------------------------
// DELETE (by UUID)
// ------------------------
func (r *UserRepository) Delete(ctx context.Context, uid uuid.UUID) error {
	return r.db.WithContext(ctx).Table("users").
		Where("uuid = ?", uid).
		Delete(&domainuser.User{}).Error
}

// ------------------------
// MERGE USER
// ------------------------
func (r *UserRepository) MergeUser(ctx context.Context, primaryID uint, secondaryIDs []uint) error {
	_ = r.ensureIdUserMergeColumn(ctx)
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 1. Akun primary tidak boleh merujuk ke merge lain
		if err := tx.Exec("UPDATE `users` SET `id_user_merge` = NULL WHERE `id` = ?", primaryID).Error; err != nil {
			return err
		}

		// 2. Akun secondary di-update id_user_merge = primaryID
		if len(secondaryIDs) > 0 {
			if err := tx.Table("users").Where("id IN (?)", secondaryIDs).Update("id_user_merge", primaryID).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

// ------------------------
// UNMERGE USER
// ------------------------
func (r *UserRepository) UnmergeUser(ctx context.Context, userID uint) error {
	return r.db.WithContext(ctx).Exec("UPDATE `users` SET `id_user_merge` = NULL WHERE `id` = ?", userID).Error
}

func (r *UserRepository) ensureIdUserMergeColumn(ctx context.Context) error {
	var colCount int64
	checkSQL := `
		SELECT COUNT(*)
		FROM INFORMATION_SCHEMA.COLUMNS
		WHERE TABLE_SCHEMA = DATABASE()
		  AND TABLE_NAME = 'users'
		  AND COLUMN_NAME = 'id_user_merge'
	`
	if err := r.db.WithContext(ctx).Raw(checkSQL).Scan(&colCount).Error; err != nil {
		return err
	}

	if colCount == 0 {
		alterSQL := "ALTER TABLE `users` ADD COLUMN `id_user_merge` INT UNSIGNED NULL AFTER `fakultas_unit`"
		if err := r.db.WithContext(ctx).Exec(alterSQL).Error; err != nil {
			return err
		}
	}

	return nil
}

func (r *UserRepository) ensureUuidColumn(ctx context.Context, tableName string) error {
	var colCount int64
	checkSQL := `
		SELECT COUNT(*)
		FROM INFORMATION_SCHEMA.COLUMNS
		WHERE TABLE_SCHEMA = DATABASE()
		  AND TABLE_NAME = ?
		  AND COLUMN_NAME = 'uuid'
	`
	if err := r.db.WithContext(ctx).Raw(checkSQL, tableName).Scan(&colCount).Error; err != nil {
		return err
	}

	if colCount == 0 {
		alterSQL := fmt.Sprintf("ALTER TABLE `%s` ADD COLUMN `uuid` CHAR(36) NULL AFTER `id`", tableName)
		if err := r.db.WithContext(ctx).Exec(alterSQL).Error; err != nil {
			return err
		}
	}

	return nil
}

func (r *UserRepository) SetupUuid(ctx context.Context) error {
	_ = r.ensureIdUserMergeColumn(ctx)
	if err := r.ensureUuidColumn(ctx, "users"); err != nil {
		return err
	}

	const chunkSize = 500

	var ids []uint
	if err := r.db.WithContext(ctx).
		Model(&domainuser.User{}).
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
			"UPDATE users SET uuid = %s WHERE id IN (?)",
			caseSQL,
		)

		if err := r.db.WithContext(ctx).
			Exec(query, args...).Error; err != nil {
			return err
		}
	}

	return nil
}

func (r *UserRepository) WithTx(
	ctx context.Context,
	fn func(txRepo domainuser.IUserRepositoryTx) error,
) error {

	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		txRepo := &UserRepository{
			db:  tx,
			uow: commoninfra.NewUnitOfWork(tx),
		}
		return fn(txRepo)
	})
}

func (r *UserRepository) InsertOutbox(
	ctx context.Context,
	msg *commoninfra.OutboxMessage,
) error {
	return r.db.WithContext(ctx).Create(msg).Error
}
