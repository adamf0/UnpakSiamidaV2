package infrastructure

import (
	"context"
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
	var extraRoles []domain.ExtraRole

	err := r.db.WithContext(ctx).
		Raw(`SELECT 
				tahun, 
				(
				case 
					when auditee=? then "auditee"
					when auditor1=? then "auditor1"
					when auditor2=? then "auditor2"
				end
				) as role 
			FROM renstra 
			WHERE auditee = ? OR auditor1 = ? OR auditor2 = ? group by tahun, role`, idUser, idUser, idUser, idUser, idUser, idUser).
		Scan(&extraRoles).Error
	if err != nil {
		return nil, err
	}

	return extraRoles, nil
}
