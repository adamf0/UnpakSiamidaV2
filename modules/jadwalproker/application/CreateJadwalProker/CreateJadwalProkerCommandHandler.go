package application

import (
	"context"
	"errors"

	domainfakultasunit "UnpakSiamida/modules/fakultasunit/domain"
	domainjadwalproker "UnpakSiamida/modules/jadwalproker/domain"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CreateJadwalProkerCommandHandler struct {
	Repo             domainjadwalproker.IJadwalProkerRepository
	RepoFakultasUnit domainfakultasunit.IFakultasUnitRepository
}

func (h *CreateJadwalProkerCommandHandler) Handle(
	ctx context.Context,
	cmd CreateJadwalProkerCommand,
) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	uuidFakultas, err := uuid.Parse(cmd.FakultasUuid)
	if err != nil {
		return "", domainjadwalproker.InvalidFakultas()
	}
	fakultas, err := h.RepoFakultasUnit.GetDefaultByUuid(ctx, uuidFakultas)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", domainjadwalproker.NotFoundFakultas()
		}
		return "", err
	}

	result := domainjadwalproker.NewJadwalProker(
		fakultas.ID,
		cmd.TanggalTutupEntry,
		cmd.TanggalTutupDokumen,
	)

	if !result.IsSuccess {
		return "", result.Error
	}

	createJadwalProker := result.Value
	if err := h.Repo.Create(ctx, createJadwalProker); err != nil {
		return "", err
	}

	return result.Value.UUID.String(), nil
}
