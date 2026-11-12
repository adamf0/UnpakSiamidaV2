package infrastructure

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"gorm.io/gorm"

	domain "UnpakSiamida/modules/account/domain"
)

type AccountRepository struct {
	db *gorm.DB
}

func NewAccountRepository(db *gorm.DB) domain.IAccountRepository {
	return &AccountRepository{db: db}
}

func (r *AccountRepository) Auth(ctx context.Context, username string, password string) (*domain.Account, error) {
	var user domain.Account

	err := r.db.WithContext(ctx).
		Table("usersv3").
		Select(`
			usersv3.id,
			usersv3.uuid,
			usersv3.nidn_username,
			usersv3.password,
			usersv3.level,
			usersv3.name,
			usersv3.email,
			usersv3.fakultas_unit,
			CASE 
				WHEN fu.type = 'prodi' AND fu.jenjang IS NOT NULL AND fu.jenjang != '' 
				THEN CONCAT(fu.nama_fak_prod_unit, ' (', UPPER(fu.jenjang), ')') 
				ELSE fu.nama_fak_prod_unit 
			END AS nama_fakultas_unit
		`).
		Joins("LEFT JOIN v_fakultas_unit fu ON usersv3.fakultas_unit = fu.id").
		Where("usersv3.nidn_username = ? AND usersv3.`password` = ?", username, password).
		Take(&user).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, err
	}

	// Ambil extrarole dari renstra
	extraRoles, err := r.getExtraRole(ctx, user.ID)
	if err != nil {
		return nil, err
	}
	user.ExtraRole = extraRoles

	return &user, nil
}

func (r *AccountRepository) Get(ctx context.Context, userUUID string) (*domain.Account, error) {
	var user domain.Account

	err := r.db.WithContext(ctx).
		Table("usersv3").
		Select(`
			usersv3.id,
			usersv3.uuid,
			usersv3.nidn_username,
			usersv3.password,
			usersv3.level,
			usersv3.name,
			usersv3.email,
			usersv3.fakultas_unit,
			CASE 
				WHEN fu.type = 'prodi' AND fu.jenjang IS NOT NULL AND fu.jenjang != '' 
				THEN CONCAT(fu.nama_fak_prod_unit, ' (', UPPER(fu.jenjang), ')') 
				ELSE fu.nama_fak_prod_unit 
			END AS nama_fakultas_unit
		`).
		Joins("LEFT JOIN v_fakultas_unit fu ON usersv3.fakultas_unit = fu.id").
		Where("usersv3.uuid = ?", userUUID).
		Take(&user).Error
	if err != nil {
		return nil, err
	}

	extraRoles, err := r.getExtraRole(ctx, user.ID)
	if err != nil {
		return nil, err
	}
	user.ExtraRole = extraRoles

	return &user, nil
}

// Ambil extrarole tetap pakai GORM
func (r *AccountRepository) getExtraRole(ctx context.Context, idUser any) ([]domain.ExtraRole, error) {
	var extraRoles = make([]domain.ExtraRole, 0)
	idStr := fmt.Sprint(idUser)

	// Kumpulkan ID pengguna ini beserta ID yang dimerge (jika ada)
	userIDsMap := make(map[string]bool)
	if idStr != "" && idStr != "0" {
		userIDsMap[idStr] = true
	}

	var secondaryIDs []string
	_ = r.db.WithContext(ctx).Table("usersv3").Where("id_user_merge = ?", idStr).Pluck("id", &secondaryIDs).Error
	for _, rid := range secondaryIDs {
		if rid != "" && rid != "0" {
			userIDsMap[rid] = true
		}
	}

	var parentMerge *string
	_ = r.db.WithContext(ctx).Table("usersv3").Where("id = ? AND id_user_merge IS NOT NULL AND id_user_merge > 0", idStr).Pluck("id_user_merge", &parentMerge).Error
	if parentMerge != nil && *parentMerge != "" && *parentMerge != "0" {
		userIDsMap[*parentMerge] = true
		var siblingIDs []string
		_ = r.db.WithContext(ctx).Table("usersv3").Where("id_user_merge = ?", *parentMerge).Pluck("id", &siblingIDs).Error
		for _, sid := range siblingIDs {
			if sid != "" && sid != "0" {
				userIDsMap[sid] = true
			}
		}
	}

	var userIDs []string
	for uid := range userIDsMap {
		userIDs = append(userIDs, uid)
	}

	if len(userIDs) == 0 {
		return extraRoles, nil
	}

	// 1. Cek tabel renstra untuk role: auditee, auditor1, dan auditor2
	unionRenstraSQL := `
		SELECT DISTINCT
			r.tahun,
			'auditee' AS role,
			r.fakultas_unit AS id_fakultas_unit,
			CASE 
				WHEN fu.type = 'prodi' AND fu.jenjang IS NOT NULL AND fu.jenjang != '' 
				THEN CONCAT(fu.nama_fak_prod_unit, ' (', UPPER(fu.jenjang), ')') 
				ELSE fu.nama_fak_prod_unit 
			END AS nama_fakultas_unit
		FROM renstra r
		LEFT JOIN v_fakultas_unit fu ON r.fakultas_unit = fu.id
		WHERE r.auditee IN (?)
		UNION
		SELECT DISTINCT
			r.tahun,
			'auditor1' AS role,
			r.fakultas_unit AS id_fakultas_unit,
			CASE 
				WHEN fu.type = 'prodi' AND fu.jenjang IS NOT NULL AND fu.jenjang != '' 
				THEN CONCAT(fu.nama_fak_prod_unit, ' (', UPPER(fu.jenjang), ')') 
				ELSE fu.nama_fak_prod_unit 
			END AS nama_fakultas_unit
		FROM renstra r
		LEFT JOIN v_fakultas_unit fu ON r.fakultas_unit = fu.id
		WHERE r.auditor1 IN (?)
		UNION
		SELECT DISTINCT
			r.tahun,
			'auditor2' AS role,
			r.fakultas_unit AS id_fakultas_unit,
			CASE 
				WHEN fu.type = 'prodi' AND fu.jenjang IS NOT NULL AND fu.jenjang != '' 
				THEN CONCAT(fu.nama_fak_prod_unit, ' (', UPPER(fu.jenjang), ')') 
				ELSE fu.nama_fak_prod_unit 
			END AS nama_fakultas_unit
		FROM renstra r
		LEFT JOIN v_fakultas_unit fu ON r.fakultas_unit = fu.id
		WHERE r.auditor2 IN (?)
	`

	var renstraRoles []domain.ExtraRole
	if err := r.db.WithContext(ctx).Raw(unionRenstraSQL, userIDs, userIDs, userIDs).Scan(&renstraRoles).Error; err != nil {
		return nil, err
	}
	extraRoles = append(extraRoles, renstraRoles...)

	// 2. Ambil role dari penugasan userv3_abac (UPML, UPMF, UPMS)
	var abacRoles []domain.ExtraRole
	_ = r.db.WithContext(ctx).Table("userv3_abac").
		Select(`
		userv3_abac.tahun,
		CASE
			WHEN userv3_abac.role IS NOT NULL AND userv3_abac.role != '' THEN LOWER(userv3_abac.role)
			WHEN fu.type = 'unit' THEN 'upml'
			WHEN LOWER(fu.nama_fak_prod_unit) LIKE '%pascasarjana%' OR LOWER(fu.nama_fak_prod_unit) LIKE '%vokasi%' OR LOWER(fu.nama_fak_prod_unit) LIKE '%sekolah%' THEN 'upms'
			WHEN userv3_abac.id_fakultas_unit IS NULL OR userv3_abac.id_fakultas_unit = 0 THEN 'upms'
			ELSE 'upmf'
		END AS role,
		COALESCE(userv3_abac.id_fakultas_unit, 0) as id_fakultas_unit,
		CASE
			WHEN userv3_abac.id_fakultas_unit IS NULL OR userv3_abac.id_fakultas_unit = 0 THEN 'Seluruh Unit (Tingkat Universitas)'
			WHEN fu.type = 'prodi' AND fu.jenjang IS NOT NULL AND fu.jenjang != ''
			THEN CONCAT(fu.nama_fak_prod_unit, ' (', UPPER(fu.jenjang), ')')
			ELSE fu.nama_fak_prod_unit
		END AS nama_fakultas_unit
	`).
		Joins("LEFT JOIN v_fakultas_unit fu ON userv3_abac.id_fakultas_unit = fu.id").
		Where("userv3_abac.id_user IN (?) AND (LOWER(userv3_abac.status) IN ('active', 'aktif') OR userv3_abac.status IS NULL OR userv3_abac.status = '') AND LOWER(userv3_abac.role) IN ('upms', 'upmf', 'upml')", userIDs).
		Group("userv3_abac.tahun, role, userv3_abac.id_fakultas_unit, fu.nama_fak_prod_unit, fu.type, fu.jenjang").
		Scan(&abacRoles).Error

	extraRoles = append(extraRoles, abacRoles...)

	return extraRoles, nil
}

func (r *AccountRepository) IsUsernameTakenByOther(ctx context.Context, username string, excludeUserUUID string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Table("usersv3").
		Where("LOWER(TRIM(nidn_username)) = ? AND uuid <> ?", strings.ToLower(strings.TrimSpace(username)), excludeUserUUID).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *AccountRepository) UpdateProfile(ctx context.Context, userUUID string, username string, name string, email *string, password *string) error {
	updates := map[string]interface{}{
		"nidn_username": username,
		"name":          name,
		"email":         nil,
		"updated_at":    time.Now(),
	}
	if email != nil && strings.TrimSpace(*email) != "" {
		updates["email"] = strings.TrimSpace(*email)
	}
	if password != nil && strings.TrimSpace(*password) != "" {
		updates["password"] = strings.TrimSpace(*password)
	}

	res := r.db.WithContext(ctx).
		Table("usersv3").
		Where("uuid = ?", userUUID).
		Updates(updates)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
