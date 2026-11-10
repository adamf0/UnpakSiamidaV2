package application

import (
	"context"
	"errors"
	"strings"
	"time"

	domainaccount "UnpakSiamida/modules/account/domain"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UpdateProfileCommandHandler struct {
	Repo domainaccount.IAccountRepository
}

func (h *UpdateProfileCommandHandler) Handle(
	ctx context.Context,
	cmd UpdateProfileCommand,
) (*domainaccount.Account, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	if _, err := uuid.Parse(cmd.SID); err != nil {
		return nil, domainaccount.InvalidUuid()
	}

	existing, err := h.Repo.Get(ctx, cmd.SID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domainaccount.NotFound(cmd.SID)
		}
		return nil, err
	}

	newUsername := strings.TrimSpace(cmd.Username)
	newName := strings.TrimSpace(cmd.Name)
	trimmedEmail := strings.TrimSpace(cmd.Email)

	// Guard duplicate username against other rows only when changed (preserves legacy duplicate auditor accounts updating email/name).
	if !strings.EqualFold(newUsername, strings.TrimSpace(existing.NidnUsername)) {
		taken, err := h.Repo.IsUsernameTakenByOther(ctx, newUsername, cmd.SID)
		if err != nil {
			return nil, err
		}
		if taken {
			return nil, domainaccount.UsernameDuplicate(newUsername)
		}
	}

	var emailPtr *string
	if trimmedEmail != "" {
		emailPtr = &trimmedEmail
	}

	var passwordPtr *string
	if trimmedNewPwd := strings.TrimSpace(cmd.NewPassword); trimmedNewPwd != "" {
		if cmd.CurrentPassword != existing.Password {
			return nil, domainaccount.InvalidCurrentPassword()
		}
		passwordPtr = &trimmedNewPwd
	}

	if err := h.Repo.UpdateProfile(ctx, cmd.SID, newUsername, newName, emailPtr, passwordPtr); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domainaccount.NotFound(cmd.SID)
		}
		return nil, err
	}

	updated, err := h.Repo.Get(ctx, cmd.SID)
	if err != nil {
		return nil, err
	}
	if updated.ExtraRole == nil {
		updated.ExtraRole = []domainaccount.ExtraRole{}
	}

	return updated, nil
}
