package application

import (
	"context"
	"strconv"
	"strings"

	domainuser "UnpakSiamida/modules/user/domain"
	"time"
)

type CreateUserCommandHandler struct{
	Repo domainuser.IUserRepository
}

func (h *CreateUserCommandHandler) Handle(
	ctx context.Context,
	cmd CreateUserCommand,
) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var fu *int = nil

	if cmd.FakultasUnit != nil {
		raw := strings.TrimSpace(*cmd.FakultasUnit)

		if raw != "" {
			parsed, err := strconv.ParseInt(raw, 10, 64)
			if err != nil {
				return "", err
			}

			v := int(parsed)
			fu = &v
		}
	}

	result := domainuser.NewUser(
		cmd.Username,
		cmd.Password,
		cmd.Name,
		cmd.Email,
		cmd.Level,
		fu,
	)

	if !result.IsSuccess {
		return "", result.Error
	}

	createUser := result.Value
	if err := h.Repo.Create(ctx, createUser); err != nil {
		return "", err
	}

	return result.Value.UUID.String(), nil
}
