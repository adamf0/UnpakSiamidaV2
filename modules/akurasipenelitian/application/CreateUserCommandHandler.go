package application

import (
    "context"
    "strconv"
    domain "UnpakSiamida/modules/akurasipenelitian/domain"
)

type CreateUserCommandHandler struct{}

func (h *CreateUserCommandHandler) Handle(
    ctx context.Context,
    cmd CreateUserCommand,
) (string, error) {
    skor := 0
	if parsed, err := strconv.Atoi(cmd.Skor); err == nil {
		skor = parsed
	}

    result := domain.NewAkurasiPenelitian(cmd.Nama, skor)

    if !result.IsSuccess {
        return "", result.Error
    }

    ap := result.Value
    return ap.UUID.String(), nil
}
