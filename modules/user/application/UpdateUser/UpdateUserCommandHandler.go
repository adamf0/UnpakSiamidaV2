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
	fakultasUnitUUID := uuid.New()
	if cmd.UuidFakultasUnit != nil {
		parsed, err := uuid.Parse(*cmd.UuidFakultasUnit)
		if err != nil {
			return "", domainuser.InvalidParsing("Fakultas Unit")
		}
		fakultasUnitUUID = parsed
	}

	var fu *int = nil
	target, err := h.RepoFakultasUnit.GetDefaultByUuid(ctx, fakultasUnitUUID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", domainuser.NotFoundFakultasUnit(*cmd.UuidFakultasUnit)
		}
		return "", err
	}
	fuId := int(target.ID)
	fu = &fuId

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
		&target.Nama,
		&target.Type,
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

	// -------------------------
	// SAVE TO REPOSITORY
	// -------------------------
	// if err := h.Repo.Update(ctx, updatedUser); err != nil {
	// 	return "", err
	// }

	return updatedUser.UUID.String(), nil
}
