package application

import (
	domainfakultasunit "UnpakSiamida/modules/fakultasunit/domain"
	"context"
	"errors"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UpdateFakultasUnitCommandHandler struct {
	Repo domainfakultasunit.IFakultasUnitRepository
}

func (h *UpdateFakultasUnitCommandHandler) Handle(
	ctx context.Context,
	cmd UpdateFakultasUnitCommand,
) (*domainfakultasunit.FakultasUnit, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	if cmd.Nama == nil && cmd.IsActive == nil {
		return nil, domainfakultasunit.EmptyUpdatePayload()
	}

	uid, err := uuid.Parse(cmd.Uuid)
	if err != nil {
		return nil, domainfakultasunit.InvalidUuid()
	}

	var trimmedNama *string
	if cmd.Nama != nil {
		v := strings.TrimSpace(*cmd.Nama)
		if v == "" || utf8.RuneCountInString(v) > 100 {
			return nil, domainfakultasunit.InvalidNama()
		}
		trimmedNama = &v
	}

	_, err = h.Repo.GetDefaultByUuid(ctx, uid)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domainfakultasunit.NotFound(cmd.Uuid)
		}
		return nil, err
	}

	if err := h.Repo.UpdateTarget(ctx, uid, trimmedNama, cmd.IsActive); err != nil {
		return nil, err
	}

	updated, err := h.Repo.GetDefaultByUuid(ctx, uid)
	if err != nil {
		return nil, err
	}

	return updated, nil
}
