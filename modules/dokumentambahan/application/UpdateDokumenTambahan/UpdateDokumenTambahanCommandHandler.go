package application

import (
	"context"
	"golang.org/x/sync/errgroup"
	"github.com/google/uuid"
	domainrenstra "UnpakSiamida/modules/renstra/domain"
	domaindokumentambahan "UnpakSiamida/modules/dokumentambahan/domain"
	"time"
)

type UpdateDokumenTambahanCommandHandler struct{
	Repo domaindokumentambahan.IDokumenTambahanRepository
	RepoRenstra domainrenstra.IRenstraRepository
}

func (h *UpdateDokumenTambahanCommandHandler) Handle(
	ctx context.Context,
	cmd UpdateDokumenTambahanCommand,
) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	dokumenTambahanUUID, err := uuid.Parse(cmd.Uuid)
	if err != nil {
		return "", domaindokumentambahan.InvalidUuid()
	}

	renstraUUID, err := uuid.Parse(cmd.UuidRenstra)
	if err != nil {
		return "", domaindokumentambahan.InvalidRenstra()
	}

	//paralel
	var (
		prev    *domaindokumentambahan.DokumenTambahan
		renstra *domainrenstra.Renstra
	)

	g, gctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		var err error
		prev, err = h.Repo.GetByUuid(gctx, dokumenTambahanUUID)
		if err != nil {
			return domaindokumentambahan.NotFound(cmd.Uuid)
		}
		return nil
	})

	g.Go(func() error {
		var err error
		renstra, err = h.RepoRenstra.GetByUuid(gctx, renstraUUID)
		if err != nil {
			return domaindokumentambahan.NotFoundRenstra(cmd.UuidRenstra)
		}
		return nil
	})

	if err := g.Wait(); err != nil {
		return "", err
	}

	result := domaindokumentambahan.UpdateDokumenTambahan(
		prev,
		renstra,
		dokumenTambahanUUID,
		renstraUUID,
		cmd.Tahun,
		cmd.Mode,
		cmd.Granted,
		cmd.Link,
		cmd.CapaianAuditor,
		cmd.CatatanAuditor,
	)

	if !result.IsSuccess {
		return "", result.Error
	}

	dokumenTambahan := result.Value

	if err := h.Repo.Update(ctx, dokumenTambahan); err != nil {
		return "", err
	}

	return dokumenTambahan.UUID.String(), nil
}