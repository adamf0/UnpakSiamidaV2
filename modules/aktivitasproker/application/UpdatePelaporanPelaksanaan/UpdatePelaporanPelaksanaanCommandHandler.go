package application

import (
	"context"
	"errors"
	"strings"
	"time"

	domainaktivitasproker "UnpakSiamida/modules/aktivitasproker/domain"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UpdatePelaporanPelaksanaanCommandHandler struct {
	Repo domainaktivitasproker.IAktivitasProkerRepository
}

func (h *UpdatePelaporanPelaksanaanCommandHandler) Handle(
	ctx context.Context,
	cmd UpdatePelaporanPelaksanaanCommand,
) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	uuidAktivitasProker, err := uuid.Parse(strings.TrimSpace(cmd.Uuid))
	if err != nil {
		return "", domainaktivitasproker.InvalidUuid()
	}

	existingAktivitasProker, err := h.Repo.GetByUuid(ctx, uuidAktivitasProker)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", domainaktivitasproker.NotFound(cmd.Uuid)
		}
		return "", err
	}

	result := domainaktivitasproker.UpdatePelaporanPelaksanaan(
		existingAktivitasProker,
		uuidAktivitasProker,
		strings.TrimSpace(cmd.StatusRealisasi),
		strings.TrimSpace(cmd.TanggalRAwal),
		strings.TrimSpace(cmd.TanggalRAkhir),
		strings.TrimSpace(cmd.Detail),
	)

	if !result.IsSuccess {
		return "", result.Error
	}

	updated := result.Value
	if err := h.Repo.Update(ctx, updated); err != nil {
		return "", err
	}

	return updated.UUID.String(), nil
}
