package application

import (
	"context"
	"encoding/json"
	"errors"
	"reflect"
	"strings"

	commoninfra "UnpakSiamida/common/infrastructure"
	domainfakultasunit "UnpakSiamida/modules/fakultasunit/domain"
	domainuser "UnpakSiamida/modules/user/domain"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UpdateUserCommandHandler struct {
	Repo             domainuser.IUserRepository
	RepoFakultasUnit domainfakultasunit.IFakultasUnitRepository
}

func (h *UpdateUserCommandHandler) Handle(
	ctx context.Context,
	cmd UpdateUserCommand,
) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// -------------------------
	// VALIDATE UUID
	// -------------------------
	userUUID, err := uuid.Parse(cmd.Uuid)
	if err != nil {
		return "", domainuser.InvalidUuid()
	}

	// -------------------------
	// GET EXISTING USER
	// -------------------------
	existingUser, err := h.Repo.GetByUuid(ctx, userUUID) // ← memastikan pakai nama interface yg benar
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", domainuser.NotFound(cmd.Uuid)
		}
		return "", err
	}

	// -------------------------
	// HANDLE PASSWORD
	// -------------------------
	var password *string = nil

	if cmd.Password != nil { // ← FIX: sebelumnya cmd.password (typo)
		raw := strings.TrimSpace(*cmd.Password)
		if raw != "" {
			password = &raw
		}
	}

	// -------------------------
	// HANDLE FakultasUnit (string → int pointer)
	// -------------------------
	var fu *int
	var targetNama *string
	var targetType *string

	if cmd.UuidFakultasUnit != nil && strings.TrimSpace(*cmd.UuidFakultasUnit) != "" {
		parsed, err := uuid.Parse(strings.TrimSpace(*cmd.UuidFakultasUnit))
		if err != nil {
			return "", domainuser.InvalidParsing("Fakultas Unit")
		}
		target, err := h.RepoFakultasUnit.GetDefaultByUuid(ctx, parsed)
		if err == nil && target != nil {
			fuId := int(target.ID)
			fu = &fuId
			targetNama = &target.Nama
			targetType = &target.Type
		}
	} else if cmd.Level == "fakultas" || cmd.Level == "auditee" {
		fu = existingUser.FakultasUnit
	}

	// -------------------------
	// AGGREGATE ROOT LOGIC
	// -------------------------
	result := domainuser.UpdateUser(
		existingUser,
		userUUID,
		cmd.Username,
		password,
		cmd.Name,
		cmd.Email,
		cmd.Level,
		fu,
		targetNama,
		targetType,
	)

	if !result.IsSuccess {
		return "", result.Error
	}

	updatedUser := result.Value

	err = h.Repo.WithTx(ctx, func(txRepo domainuser.IUserRepositoryTx) error {

		if err := txRepo.Update(ctx, updatedUser); err != nil {
			return err
		}

		for _, event := range updatedUser.DomainEvents() {
			payload, err := json.Marshal(event)
			if err != nil {
				return err
			}

			outbox := commoninfra.OutboxMessage{
				ID:            event.ID(),
				Type:          reflect.TypeOf(event).String(),
				Payload:       string(payload),
				OccurredOnUTC: event.OccurredOnUTC(),
			}

			if err := txRepo.InsertOutbox(ctx, &outbox); err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return "", err
	}

	// 🧹 CLEAR AFTER COMMIT
	updatedUser.ClearDomainEvents()

	if updatedUser.FakultasUnit != nil && *updatedUser.FakultasUnit > 0 {
		fuUint := uint(*updatedUser.FakultasUnit)
		roleName := strings.ToLower(strings.TrimSpace(updatedUser.Level))
		if roleName == "" {
			roleName = "user"
		}
		for _, yr := range []string{"2026", "2027"} {
			exists, _ := h.Repo.CheckUserAbacExists(ctx, updatedUser.ID, yr, &fuUint)
			if !exists {
				now := time.Now()
				_ = h.Repo.CreateUserAbac(ctx, &domainuser.UserAbac{
					UUID:           uuid.New(),
					IdUser:         updatedUser.ID,
					Tahun:          yr,
					IdFakultasUnit: &fuUint,
					Role:           roleName,
					Status:         "active",
					CreatedAt:      now,
					UpdatedAt:      now,
				})
			}
		}
	}

	// -------------------------
	// SAVE TO REPOSITORY
	// -------------------------
	// if err := h.Repo.Update(ctx, updatedUser); err != nil {
	// 	return "", err
	// }

	return updatedUser.UUID.String(), nil
}
