package application

import (
	domainrenstra "UnpakSiamida/modules/renstra/domain"
	domainrenstranilai "UnpakSiamida/modules/renstranilai/domain"
	"context"
	"errors"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
	"golang.org/x/sync/errgroup"
)

type UpdateRenstraNilaiCommandHandler struct {
	Repo        domainrenstranilai.IRenstraNilaiRepository
	RepoRenstra domainrenstra.IRenstraRepository
}

func (h *UpdateRenstraNilaiCommandHandler) Handle(
	ctx context.Context,
	cmd UpdateRenstraNilaiCommand,
) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	renstraNilaiUUID, err := uuid.Parse(cmd.Uuid)
	if err != nil {
		return "", domainrenstranilai.InvalidUuid()
	}

	renstraUUID, err := uuid.Parse(cmd.UuidRenstra)
	if err != nil {
		return "", domainrenstranilai.InvalidRenstra()
	}

	//paralel
	var (
		prev    *domainrenstranilai.RenstraNilai
		renstra *domainrenstra.Renstra
	)

	g, gctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		var err error
		prev, err = h.Repo.GetByUuid(gctx, renstraNilaiUUID)
		if err != nil {
			return domainrenstranilai.NotFound(cmd.Uuid)
		}
		return nil
	})

	g.Go(func() error {
		var err error
		renstra, err = h.RepoRenstra.GetByUuid(gctx, renstraUUID)
		if err != nil {
			return domainrenstranilai.NotFoundRenstra(cmd.UuidRenstra)
		}
		return nil
	})

	if err := g.Wait(); err != nil {
		return "", err
	}

	result := domainrenstranilai.UpdateRenstraNilai(
		prev,
		renstra,
		renstraNilaiUUID,
		renstraUUID,
		cmd.Tahun,
		cmd.Mode,
		cmd.Granted,
		cmd.Capaian,
		cmd.Catatan,
		cmd.LinkBukti,
		cmd.CapaianAuditor,
		cmd.CatatanAuditor,
	)

	if !result.IsSuccess {
		return "", result.Error
	}

	renstraNilai := result.Value

	if err := h.Repo.Update(ctx, renstraNilai); err != nil {
		//[pr] ini harusnya ada di generate renstra bagian create + update
		//ini rule mustahil masuk, tapi memungkinan jika datanya hardcode oleh developer untuk duplicate lalu di update baru kena ini
		//ini tidak dapat di tes karena di domainnya ada rule pengecekan prev data
		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) {
			if mysqlErr.Number == 1062 {
				return "", domainrenstranilai.DuplicateData()
			}
		}

		return "", err
	}

	return renstraNilai.UUID.String(), nil
}
