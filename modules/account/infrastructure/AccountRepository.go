package infrastructure

import (
	"context"
	"errors"
	"fmt"

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
		Table("users").
		Select(`
			users.id,
			users.uuid,
			users.nidn_username,
			users.password,
			users.level,
			users.name,
			users.email,
			users.fakultas_unit,
			CASE 
				WHEN fu.type = 'prodi' AND fu.jenjang IS NOT NULL AND fu.jenjang != '' 
				THEN CONCAT(fu.nama_fak_prod_unit, ' (', UPPER(fu.jenjang), ')') 
				ELSE fu.nama_fak_prod_unit 
			END AS nama_fakultas_unit
		`).
		Joins("LEFT JOIN v_fakultas_unit fu ON users.fakultas_unit = fu.id").
		Where("users.nidn_username = ? AND users.`password` = ?", username, password).
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
		Table("users").
		Select(`
			users.id,
			users.uuid,
			users.nidn_username,
			users.password,
			users.level,
			users.name,
			users.email,
			users.fakultas_unit,
			CASE 
				WHEN fu.type = 'prodi' AND fu.jenjang IS NOT NULL AND fu.jenjang != '' 
				THEN CONCAT(fu.nama_fak_prod_unit, ' (', UPPER(fu.jenjang), ')') 
				ELSE fu.nama_fak_prod_unit 
			END AS nama_fakultas_unit
		`).
		Joins("LEFT JOIN v_fakultas_unit fu ON users.fakultas_unit = fu.id").
		Where("users.uuid = ?", userUUID).
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
	_ = r.db.WithContext(ctx).Table("users").Where("id_user_merge = ?", idStr).Pluck("id", &secondaryIDs).Error
	for _, rid := range secondaryIDs {
		if rid != "" && rid != "0" {
			userIDsMap[rid] = true
		}
	}

	var parentMerge *string
	_ = r.db.WithContext(ctx).Table("users").Where("id = ? AND id_user_merge IS NOT NULL AND id_user_merge > 0", idStr).Pluck("id_user_merge", &parentMerge).Error
	if parentMerge != nil && *parentMerge != "" && *parentMerge != "0" {
		userIDsMap[*parentMerge] = true
		var siblingIDs []string
		_ = r.db.WithContext(ctx).Table("users").Where("id_user_merge = ?", *parentMerge).Pluck("id", &siblingIDs).Error
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

	// 2. Ambil role dari penugasan user_abac (UPML, UPMF, UPMS)
	var abacRoles []domain.ExtraRole
	_ = r.db.WithContext(ctx).Table("user_abac").
		Select(`
		user_abac.tahun,
		CASE
			WHEN user_abac.role IS NOT NULL AND user_abac.role != '' THEN LOWER(user_abac.role)
			WHEN fu.type = 'unit' THEN 'upml'
			WHEN LOWER(fu.nama_fak_prod_unit) LIKE '%pascasarjana%' OR LOWER(fu.nama_fak_prod_unit) LIKE '%vokasi%' OR LOWER(fu.nama_fak_prod_unit) LIKE '%sekolah%' THEN 'upms'
			WHEN user_abac.id_fakultas_unit IS NULL OR user_abac.id_fakultas_unit = 0 THEN 'upms'
			ELSE 'upmf'
		END AS role,
		COALESCE(user_abac.id_fakultas_unit, 0) as id_fakultas_unit,
		CASE
			WHEN user_abac.id_fakultas_unit IS NULL OR user_abac.id_fakultas_unit = 0 THEN 'Seluruh Unit (Tingkat Universitas)'
			WHEN fu.type = 'prodi' AND fu.jenjang IS NOT NULL AND fu.jenjang != ''
			THEN CONCAT(fu.nama_fak_prod_unit, ' (', UPPER(fu.jenjang), ')')
			ELSE fu.nama_fak_prod_unit
		END AS nama_fakultas_unit
	`).
		Joins("LEFT JOIN v_fakultas_unit fu ON user_abac.id_fakultas_unit = fu.id").
		Where("user_abac.id_user IN (?) AND (LOWER(user_abac.status) IN ('active', 'aktif') OR user_abac.status IS NULL OR user_abac.status = '') AND LOWER(user_abac.role) IN ('upms', 'upmf', 'upml')", userIDs).
		Group("user_abac.tahun, role, user_abac.id_fakultas_unit, fu.nama_fak_prod_unit, fu.type, fu.jenjang").
		Scan(&abacRoles).Error

	extraRoles = append(extraRoles, abacRoles...)

	return extraRoles, nil
}
