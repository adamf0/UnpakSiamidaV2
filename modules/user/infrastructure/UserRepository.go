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

	err := r.db.WithContext(ctx).Table("usersv3").
		Select(`
			usersv3.id,
			usersv3.uuid,
			usersv3.nidn_username,
			usersv3.password,
			usersv3.name,
			usersv3.email,
			usersv3.level,
			usersv3.id_user_merge,
			mu.name as nama_user_merge
		`).
		Joins("LEFT JOIN usersv3 mu ON usersv3.id_user_merge = mu.id").
		Where("usersv3.uuid = ?", uid).
		First(&user).Error

	if err != nil {
		return nil, err
	}

	var abacSummaries []domainuser.UserAbacSummary
	_ = r.db.WithContext(ctx).Table("userv3_abac").
		Select(`
			userv3_abac.id,
			userv3_abac.uuid,
			userv3_abac.id_user,
			userv3_abac.tahun,
			userv3_abac.id_fakultas_unit,
			userv3_abac.role,
			userv3_abac.status,
			CASE 
				WHEN userv3_abac.id_fakultas_unit IS NULL OR userv3_abac.id_fakultas_unit = 0 THEN 'Tingkat Universitas'
				WHEN fu.type = 'prodi' AND fu.jenjang IS NOT NULL AND fu.jenjang != '' 
				THEN CONCAT(fu.nama_fak_prod_unit, ' (', UPPER(fu.jenjang), ')') 
				ELSE fu.nama_fak_prod_unit 
			END as nama_fakultas_unit,
			CASE
				WHEN userv3_abac.role IS NOT NULL AND userv3_abac.role != '' THEN LOWER(userv3_abac.role)
				WHEN fu.type = 'unit' THEN 'upml'
				WHEN LOWER(fu.nama_fak_prod_unit) LIKE '%pascasarjana%' OR LOWER(fu.nama_fak_prod_unit) LIKE '%vokasi%' OR LOWER(fu.nama_fak_prod_unit) LIKE '%sekolah%' THEN 'upms'
				WHEN userv3_abac.id_fakultas_unit IS NULL OR userv3_abac.id_fakultas_unit = 0 THEN 'upms'
				ELSE 'upmf'
			END AS computed_role
		`).
		Joins("LEFT JOIN v_fakultas_unit fu ON userv3_abac.id_fakultas_unit = fu.id").
		Where("userv3_abac.id_user = ?", user.ID).
		Order("userv3_abac.tahun DESC, userv3_abac.id DESC").
		Find(&abacSummaries).Error
	user.UserAbac = abacSummaries
	if user.Level == "" {
		user.Level = "user"
	}
	if len(abacSummaries) > 0 {
		if user.NamaFakultasUnit == nil || *user.NamaFakultasUnit == "" {
			user.NamaFakultasUnit = &abacSummaries[0].NamaFakultasUnit
		}
		if abacSummaries[0].IdFakultasUnit != nil {
			idFuInt := int(*abacSummaries[0].IdFakultasUnit)
			user.FakultasUnit = &idFuInt
		}
	}

	return &user, nil
}

var allowedSearchColumns = map[string]string{
	// key:param -> db column
	"name":          "usersv3.name",
	"username":      "usersv3.nidn_username",
	"email":         "usersv3.email",
	"level":         "usersv3.level",
	"id_user_merge": "usersv3.id_user_merge",
}

func (r *UserRepository) GetAllStrict(ctx context.Context) ([]domainuser.UserOptions, error) {
	var users = make([]domainuser.UserOptions, 0)

	// Ambil semua data dari table "usersv3" kecuali yang memiliki level admin
	if err := r.db.WithContext(ctx).Table("usersv3").Where("LOWER(usersv3.level) <> 'admin'").Find(&users).Error; err != nil {
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

	db := r.db.WithContext(ctx).Table("usersv3").
		Select(`
			usersv3.id,
			usersv3.uuid,
			usersv3.nidn_username,
			usersv3.password,
			usersv3.name,
			usersv3.email,
			usersv3.level,
			usersv3.id_user_merge,
			mu.name as nama_user_merge
		`).
		Joins("LEFT JOIN usersv3 mu ON usersv3.id_user_merge = mu.id")

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
	if err := db.Order("usersv3.id DESC").Find(&users).Error; err != nil {
		return nil, 0, err
	}

	if len(users) > 0 {
		userIds := make([]uint, len(users))
		for i, u := range users {
			userIds[i] = u.ID
		}

		type AbacRow struct {
			ID               uint   `gorm:"column:id"`
			UUID             string `gorm:"column:uuid"`
			IdUser           uint   `gorm:"column:id_user"`
			Tahun            string `gorm:"column:tahun"`
			IdFakultasUnit   *uint  `gorm:"column:id_fakultas_unit"`
			Role             string `gorm:"column:role"`
			Status           string `gorm:"column:status"`
			NamaFakultasUnit string `gorm:"column:nama_fak_prod_unit"`
			ComputedRole     string `gorm:"column:computed_role"`
		}

		var abacRows []AbacRow
		err := r.db.WithContext(ctx).Table("userv3_abac").
			Select(`
				userv3_abac.id,
				userv3_abac.uuid,
				userv3_abac.id_user,
				userv3_abac.tahun,
				userv3_abac.id_fakultas_unit,
				userv3_abac.role,
				userv3_abac.status,
				CASE 
					WHEN userv3_abac.id_fakultas_unit IS NULL OR userv3_abac.id_fakultas_unit = 0 THEN 'Seluruh Unit (Tingkat Universitas)'
					WHEN fu.type = 'prodi' AND fu.jenjang IS NOT NULL AND fu.jenjang != '' 
					THEN CONCAT(fu.nama_fak_prod_unit, ' (', UPPER(fu.jenjang), ')') 
					ELSE fu.nama_fak_prod_unit 
				END as nama_fak_prod_unit,
				CASE
					WHEN userv3_abac.role IS NOT NULL AND userv3_abac.role != '' THEN LOWER(userv3_abac.role)
					WHEN fu.type = 'unit' THEN 'upml'
					WHEN LOWER(fu.nama_fak_prod_unit) LIKE '%pascasarjana%' OR LOWER(fu.nama_fak_prod_unit) LIKE '%vokasi%' OR LOWER(fu.nama_fak_prod_unit) LIKE '%sekolah%' THEN 'upms'
					WHEN userv3_abac.id_fakultas_unit IS NULL OR userv3_abac.id_fakultas_unit = 0 THEN 'upms'
					ELSE 'upmf'
				END AS computed_role
			`).
			Joins("LEFT JOIN v_fakultas_unit fu ON userv3_abac.id_fakultas_unit = fu.id").
			Where("userv3_abac.id_user IN ?", userIds).
			Order("userv3_abac.tahun DESC, userv3_abac.id DESC").
			Find(&abacRows).Error

		if err == nil {
			abacMap := make(map[uint][]domainuser.UserAbacSummary)
			for _, row := range abacRows {
				abacMap[row.IdUser] = append(abacMap[row.IdUser], domainuser.UserAbacSummary{
					ID:               row.ID,
					UUID:             row.UUID,
					IdUser:           row.IdUser,
					Tahun:            row.Tahun,
					IdFakultasUnit:   row.IdFakultasUnit,
					NamaFakultasUnit: row.NamaFakultasUnit,
					Role:             row.Role,
					ComputedRole:     row.ComputedRole,
					Status:           row.Status,
				})
			}

			for i := range users {
				if items, ok := abacMap[users[i].ID]; ok && len(items) > 0 {
					users[i].UserAbac = items
					if users[i].NamaFakultasUnit == nil || *users[i].NamaFakultasUnit == "" {
						users[i].NamaFakultasUnit = &items[0].NamaFakultasUnit
					}
					if items[0].IdFakultasUnit != nil {
						idFuInt := int(*items[0].IdFakultasUnit)
						users[i].FakultasUnit = &idFuInt
					}
				} else {
					users[i].UserAbac = []domainuser.UserAbacSummary{}
				}
				if users[i].Level == "" {
					users[i].Level = "user"
				}
			}
		}
	}

	return users, total, nil
}

// ------------------------
// CREATE
// ------------------------
func (r *UserRepository) Create(ctx context.Context, user *domainuser.User) error {
	if err := r.db.WithContext(ctx).Table("usersv3").Omit("NamaFakultasUnit", "UuidFakultasUnit", "nama_fak_prod_unit", "uuid_fakultas_unit").Save(user).Error; err != nil {
		return err
	}

	return r.uow.Save(&user.Entity)
}

// ------------------------
// UPDATE
// ------------------------
func (r *UserRepository) Update(ctx context.Context, user *domainuser.User) error {
	if err := r.db.WithContext(ctx).Table("usersv3").Omit("NamaFakultasUnit", "UuidFakultasUnit", "nama_fak_prod_unit", "uuid_fakultas_unit").Save(user).Error; err != nil {
		return err
	}

	return r.uow.Save(&user.Entity)
}

// ------------------------
// DELETE (by UUID)
// ------------------------
func (r *UserRepository) Delete(ctx context.Context, uid uuid.UUID) error {
	return r.db.WithContext(ctx).Table("usersv3").
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
		if err := tx.Exec("UPDATE `usersv3` SET `id_user_merge` = NULL WHERE `id` = ?", primaryID).Error; err != nil {
			return err
		}

		// 2. Akun secondary di-update id_user_merge = primaryID
		if len(secondaryIDs) > 0 {
			if err := tx.Table("usersv3").Where("id IN (?)", secondaryIDs).Update("id_user_merge", primaryID).Error; err != nil {
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
	return r.db.WithContext(ctx).Exec("UPDATE `usersv3` SET `id_user_merge` = NULL WHERE `id` = ?", userID).Error
}

func (r *UserRepository) ensureIdUserMergeColumn(ctx context.Context) error {
	var colCount int64
	checkSQL := `
		SELECT COUNT(*)
		FROM INFORMATION_SCHEMA.COLUMNS
		WHERE TABLE_SCHEMA = DATABASE()
		  AND TABLE_NAME = 'usersv3'
		  AND COLUMN_NAME = 'id_user_merge'
	`
	if err := r.db.WithContext(ctx).Raw(checkSQL).Scan(&colCount).Error; err != nil {
		return err
	}

	if colCount == 0 {
		alterSQL := "ALTER TABLE `usersv3` ADD COLUMN `id_user_merge` INT UNSIGNED NULL AFTER `fakultas_unit`"
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
	if err := r.ensureUuidColumn(ctx, "usersv3"); err != nil {
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
			"UPDATE usersv3 SET uuid = %s WHERE id IN (?)",
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

// ------------------------
// USER ABAC METHODS
// ------------------------
func (r *UserRepository) GetAllUserAbac(ctx context.Context, tahun, status, role, search, userUuid, fuUuid string) ([]domainuser.UserAbac, error) {
	query := r.db.WithContext(ctx).Table("userv3_abac").
		Select(`
			userv3_abac.id,
			userv3_abac.uuid,
			userv3_abac.id_user,
			userv3_abac.tahun,
			userv3_abac.id_fakultas_unit,
			userv3_abac.role,
			userv3_abac.status,
			userv3_abac.created_at,
			userv3_abac.updated_at,
			u.name as user_name,
			u.nidn_username as user_username,
			u.email as user_email,
			u.uuid as user_uuid,
			CASE 
				WHEN userv3_abac.id_fakultas_unit IS NULL OR userv3_abac.id_fakultas_unit = 0 THEN 'Seluruh Unit (Tingkat Universitas)'
				WHEN fu.type = 'prodi' AND fu.jenjang IS NOT NULL AND fu.jenjang != '' 
				THEN CONCAT(fu.nama_fak_prod_unit, ' (', UPPER(fu.jenjang), ')') 
				ELSE fu.nama_fak_prod_unit 
			END as nama_fakultas_unit,
			fu.jenjang as jenjang_fakultas_unit,
			fu.uuid as uuid_fakultas_unit,
			fu.type as type_fakultas_unit,
			fu.fakultas as fakultas_unit_parent,
			CASE
				WHEN userv3_abac.role IS NOT NULL AND userv3_abac.role != '' THEN LOWER(userv3_abac.role)
				WHEN fu.type = 'unit' THEN 'upml'
				WHEN LOWER(fu.nama_fak_prod_unit) LIKE '%pascasarjana%' OR LOWER(fu.nama_fak_prod_unit) LIKE '%vokasi%' OR LOWER(fu.nama_fak_prod_unit) LIKE '%sekolah%' THEN 'upms'
				WHEN userv3_abac.id_fakultas_unit IS NULL OR userv3_abac.id_fakultas_unit = 0 THEN 'upms'
				ELSE 'upmf'
			END AS computed_role
		`).
		Joins("LEFT JOIN usersv3 u ON userv3_abac.id_user = u.id").
		Joins("LEFT JOIN v_fakultas_unit fu ON userv3_abac.id_fakultas_unit = fu.id")

	if tahun != "" && tahun != "all" {
		query = query.Where("userv3_abac.tahun = ?", tahun)
	}

	if status != "" && status != "all" {
		query = query.Where("userv3_abac.status = ?", status)
	}

	if userUuid != "" {
		query = query.Where("u.uuid = ?", userUuid)
	}

	if fuUuid != "" {
		query = query.Where("fu.uuid = ?", fuUuid)
	}

	if search != "" {
		like := "%" + strings.ToLower(search) + "%"
		query = query.Where("LOWER(u.name) LIKE ? OR LOWER(u.nidn_username) LIKE ? OR LOWER(fu.nama_fak_prod_unit) LIKE ?", like, like, like)
	}

	var results []domainuser.UserAbac
	if err := query.Order("userv3_abac.id DESC").Find(&results).Error; err != nil {
		return nil, err
	}

	if role != "" && role != "all" {
		filtered := make([]domainuser.UserAbac, 0, len(results))
		for _, item := range results {
			rRole := strings.ToLower(item.ComputedRole)
			if rRole == "" {
				rRole = strings.ToLower(item.Role)
			}
			if rRole == role {
				filtered = append(filtered, item)
			}
		}
		results = filtered
	}

	return results, nil
}

func (r *UserRepository) GetUserAbacByUuid(ctx context.Context, uid uuid.UUID) (*domainuser.UserAbac, error) {
	var abac domainuser.UserAbac
	if err := r.db.WithContext(ctx).Table("userv3_abac").Where("uuid = ?", uid).First(&abac).Error; err != nil {
		return nil, err
	}
	return &abac, nil
}

func (r *UserRepository) CreateUserAbac(ctx context.Context, abac *domainuser.UserAbac) error {
	return r.db.WithContext(ctx).Table("userv3_abac").Create(abac).Error
}

func (r *UserRepository) UpdateUserAbac(ctx context.Context, id uint, updates map[string]interface{}) error {
	return r.db.WithContext(ctx).Table("userv3_abac").Where("id = ?", id).Updates(updates).Error
}

func (r *UserRepository) DeleteUserAbac(ctx context.Context, uid uuid.UUID) error {
	res := r.db.WithContext(ctx).Table("userv3_abac").Where("uuid = ?", uid).Delete(&domainuser.UserAbac{})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *UserRepository) CheckUserAbacExists(ctx context.Context, userID uint, tahun string, fuID *uint) (bool, error) {
	var count int64
	q := r.db.WithContext(ctx).Table("userv3_abac").
		Where("id_user = ? AND tahun = ?", userID, tahun)
	if fuID != nil && *fuID > 0 {
		q = q.Where("id_fakultas_unit = ?", *fuID)
	} else {
		q = q.Where("id_fakultas_unit IS NULL OR id_fakultas_unit = 0")
	}
	err := q.Count(&count).Error
	return count > 0, err
}

func (r *UserRepository) GetFakultasUnitInfo(ctx context.Context, id *uint, uidStr string) (idRes uint, typeRes string, namaRes string, err error) {
	type FuRow struct {
		ID   uint   `gorm:"column:id"`
		Type string `gorm:"column:type"`
		Nama string `gorm:"column:nama_fak_prod_unit"`
	}
	var row FuRow
	q := r.db.WithContext(ctx).Table("v_fakultas_unit").Select("id, type, nama_fak_prod_unit")
	if id != nil && *id > 0 {
		q = q.Where("id = ?", *id)
	} else if strings.TrimSpace(uidStr) != "" {
		q = q.Where("uuid = ?", strings.TrimSpace(uidStr))
	} else {
		return 0, "", "", gorm.ErrRecordNotFound
	}
	if err := q.First(&row).Error; err != nil {
		return 0, "", "", err
	}
	return row.ID, row.Type, row.Nama, nil
}
