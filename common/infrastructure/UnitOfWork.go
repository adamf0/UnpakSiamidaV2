package infrastructure

import (
	"encoding/json"

	"UnpakSiamida/common/domain"

	"gorm.io/gorm"
)

type UnitOfWork struct {
	DB *gorm.DB
}

func NewUnitOfWork(db *gorm.DB) *UnitOfWork {
	return &UnitOfWork{DB: db}
}

func (u *UnitOfWork) Save(entity *domain.Entity) error {
	return u.DB.Transaction(func(tx *gorm.DB) error {

		for _, event := range entity.DomainEvents() {

			payload, err := json.Marshal(event)
			if err != nil {
				return err
			}

			outbox := OutboxMessage{
				ID:            event.ID(),
				Type:          CanonicalTypeName(event), // ✅ FIX UTAMA
				Payload:       string(payload),
				OccurredOnUTC: event.OccurredOnUTC(),
			}

			if err := tx.Create(&outbox).Error; err != nil {
				return err
			}
		}

		entity.ClearDomainEvents()
		return nil
	})
}
