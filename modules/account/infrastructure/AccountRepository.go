package infrastructure

import (
	"context"
	"errors"

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
		Where("nidn_username = ? AND password = ?", username, password).
		First(&user).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, err
	}

	// Ambil extrarole dari renstra
	extraRoles, err := r.getExtraRole(ctx, user.NidnUsername)
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
func (r *AccountRepository) getExtraRole(ctx context.Context, idUser string) ([]domain.ExtraRole, error) {
	var extraRoles = make([]domain.ExtraRole, 0)

	db := r.db.WithContext(ctx).Table("renstra").
		Select(`
		tahun,
		CASE
			WHEN auditee = ? THEN 'auditee'
			WHEN auditor1 = ? THEN 'auditor1'
			WHEN auditor2 = ? THEN 'auditor2'
		END AS role
	`, idUser, idUser, idUser).
		Where("auditee = ? OR auditor1 = ? OR auditor2 = ?", idUser, idUser, idUser).
		Group("tahun, role")

	err := db.Scan(&extraRoles).Error
	if err != nil {
		return nil, err
	}

	return extraRoles, nil
}
