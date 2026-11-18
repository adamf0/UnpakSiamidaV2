package CreateProkerDokumen

import (
	domain "UnpakSiamida/modules/prokerv3/domain"
	"context"
)

type CreateProkerDokumenCommandHandler struct {
	Repo domain.IProkerV3Repository
}

func (h *CreateProkerDokumenCommandHandler) Handle(ctx context.Context, cmd CreateProkerDokumenCommand) (*domain.MatrixProkerDokumen, error) {
	if cmd.Dokumen == nil {
		return nil, domain.DokumenEmpty()
	}
	if cmd.Dokumen.StatusVerifikasi == "" {
		cmd.Dokumen.StatusVerifikasi = "belum_terverifikasi"
	}
	if cmd.Dokumen.StatusVerifikasiFakultas == "" {
		cmd.Dokumen.StatusVerifikasiFakultas = "belum_terverifikasi"
	}
	if cmd.Dokumen.StatusVerifikasiAdmin == "" {
		cmd.Dokumen.StatusVerifikasiAdmin = "belum_terverifikasi"
	}
	if err := h.Repo.CreateDokumen(ctx, cmd.Dokumen); err != nil {
		return nil, err
	}
	return cmd.Dokumen, nil
}
