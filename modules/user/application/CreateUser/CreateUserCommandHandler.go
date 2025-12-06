package application

import (
	"context"
	"strconv"
	"strings"

	domain "UnpakSiamida/modules/user/domain"
)

type CreateUserCommandHandler struct{
	Repo domain.IUserRepository
}

func (h *CreateUserCommandHandler) Handle(
	ctx context.Context,
	cmd CreateUserCommand,
) (string, error) {

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

	result := domain.NewUser(
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
