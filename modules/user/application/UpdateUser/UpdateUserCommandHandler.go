package application

import (
	"context"
	"strconv"
	"strings"

	domainuser "UnpakSiamida/modules/user/domain"
	"github.com/google/uuid"
	"time"
)

type UpdateUserCommandHandler struct {
	Repo domainuser.IUserRepository
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
		return "", err
	}
	if existingUser == nil {
		return "", domainuser.NotFound(cmd.Uuid)
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
	var fakultasUnit *int = nil

	if cmd.FakultasUnit != nil {
		raw := strings.TrimSpace(*cmd.FakultasUnit)

		if raw != "" {
			parsed, err := strconv.ParseInt(raw, 10, 64)
			if err != nil {
				return "", err
			}

			v := int(parsed)
			fakultasUnit = &v
		}
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
		fakultasUnit,
	)

	if !result.IsSuccess {
		return "", result.Error
	}

	updatedUser := result.Value

	// -------------------------
	// SAVE TO REPOSITORY
	// -------------------------
	if err := h.Repo.Update(ctx, updatedUser); err != nil {
		return "", err
	}

	return updatedUser.UUID.String(), nil
}
