package application

import (
	domainSwot "UnpakSiamida/modules/swot/domain"
	"context"
	"strings"
	"time"

	"github.com/google/uuid"
)

type DeleteSwotCommandHandler struct {
	Repo domainSwot.ISwotRepository
}

func (h *DeleteSwotCommandHandler) Handle(
	ctx context.Context,
	cmd DeleteSwotCommand,
) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	parsedUID, err := uuid.Parse(strings.TrimSpace(cmd.UUID))
	if err != nil {
		return "", domainSwot.InvalidData()
	}

	if err := h.Repo.Delete(ctx, parsedUID); err != nil {
		return "", err
	}

	return cmd.UUID, nil
}
