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

	// GORM query
	err := r.db.WithContext(ctx).
		Where("nidn_username = ? AND `password` = ?", username, password).
		First(&user).Error

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
		Where("uuid = ?", userUUID).
		First(&user).Error
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

	db := r.db.WithContext(ctx).Table("renstra").
		Select(`
		tahun,
		CASE
			WHEN auditee = ? THEN 'auditee'
			WHEN auditor1 = ? THEN 'auditor1'
			WHEN auditor2 = ? THEN 'auditor2'
		END AS role,
		fakultas_unit AS id_fakultas_unit
	`, idStr, idStr, idStr).
		Where("auditee = ? OR auditor1 = ? OR auditor2 = ?", idStr, idStr, idStr).
		Group("tahun, role, fakultas_unit")

	err := db.Scan(&extraRoles).Error
	if err != nil {
		return nil, err
	}

	// Ambil role dari penugasan user_abac (UPML, UPMF, UPMS)
	var abacRoles []domain.ExtraRole
	_ = r.db.WithContext(ctx).Table("user_abac").
		Select(`
		user_abac.tahun,
		CASE
			WHEN user_abac.role IS NOT NULL AND user_abac.role != '' THEN LOWER(user_abac.role)
			WHEN fu.type = 'unit' THEN 'upml'
			WHEN LOWER(fu.nama_fak_prod_unit) LIKE '%pascasarjana%' OR LOWER(fu.nama_fak_prod_unit) LIKE '%vokasi%' OR LOWER(fu.nama_fak_prod_unit) LIKE '%sekolah%' THEN 'upms'
			ELSE 'upmf'
		END AS role,
		user_abac.id_fakultas_unit,
		fu.nama_fak_prod_unit AS nama_fakultas_unit
	`).
		Joins("LEFT JOIN v_fakultas_unit fu ON user_abac.id_fakultas_unit = fu.id").
		Where("(user_abac.id_user = ? OR user_abac.id_user = (SELECT id FROM users WHERE uuid = ? LIMIT 1)) AND (LOWER(user_abac.status) IN ('active', 'aktif') OR user_abac.status IS NULL OR user_abac.status = '')", idUser, idStr).
		Group("user_abac.tahun, role, user_abac.id_fakultas_unit, fu.nama_fak_prod_unit").
		Scan(&abacRoles).Error

	extraRoles = append(extraRoles, abacRoles...)

	return extraRoles, nil
}
