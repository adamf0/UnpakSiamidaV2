package application

import (
	domaindokumentambahan "UnpakSiamida/modules/dokumentambahan/domain"
	domainrenstra "UnpakSiamida/modules/renstra/domain"
	"context"
	"errors"
	"strings"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
	"golang.org/x/sync/errgroup"
	"gorm.io/gorm"
)

type UpdateDokumenTambahanCommandHandler struct {
	Repo        domaindokumentambahan.IDokumenTambahanRepository
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
			if err == gorm.ErrRecordNotFound {
				return domaindokumentambahan.NotFound(cmd.Uuid)
			}
			return err
		}
		return nil
	})

	g.Go(func() error {
		var err error
		renstra, err = h.RepoRenstra.GetByUuid(gctx, renstraUUID)
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				return domaindokumentambahan.NotFoundRenstra(cmd.Uuid)
			}
			return err
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
		//[note] ini harusnya ada di generate renstra bagian create + update
		//ini rule mustahil masuk, tapi memungkinan jika datanya hardcode oleh developer untuk duplicate lalu di update baru kena ini
		//ini tidak dapat di tes karena di domainnya ada rule pengecekan prev data
		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) {
			if mysqlErr.Number == 1062 {
				return "", domaindokumentambahan.DuplicateData()
			}
		}
		if strings.Contains(err.Error(), "Duplicate entry") {
			return "", domaindokumentambahan.DuplicateData()
		}

		return "", err
	}

	return dokumenTambahan.UUID.String(), nil
}
