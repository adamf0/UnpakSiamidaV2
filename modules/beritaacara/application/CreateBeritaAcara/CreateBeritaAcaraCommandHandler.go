package application

import (
	"context"

	domainberitaacara "UnpakSiamida/modules/beritaacara/domain"
	"time"
)

type CreateBeritaAcaraCommandHandler struct {
	Repo domainberitaacara.IBeritaAcaraRepository
}

func (h *CreateBeritaAcaraCommandHandler) Handle(
	ctx context.Context,
	cmd CreateBeritaAcaraCommand,
) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	tanggal, err := time.Parse("2006-01-02", cmd.Tanggal)
	if err != nil {
		return "", domainberitaacara.InvalidTanggal()
	}

	//[pr] bagian auditee, auditor masih dalam bentuk int harusnya uuid
	result := domainberitaacara.NewBeritaAcara(
		cmd.Tahun,
		cmd.FakultasUnit,
		tanggal,
		cmd.Auditee,
		cmd.Auditor1,
		cmd.Auditor2,
	)

	if !result.IsSuccess {
		return "", result.Error
	}

	createBeritaAcara := result.Value
	if err := h.Repo.Create(ctx, createBeritaAcara); err != nil {
		return "", err
	}

	return result.Value.UUID.String(), nil
}
