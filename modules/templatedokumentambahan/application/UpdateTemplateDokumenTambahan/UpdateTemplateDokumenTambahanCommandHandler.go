package application

import (
	"context"

	"github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
	"golang.org/x/sync/errgroup"

	domainjenisfile "UnpakSiamida/modules/jenisfile/domain"
	domaintemplatedokumentambahan "UnpakSiamida/modules/templatedokumentambahan/domain"
	"errors"
	"time"

	"gorm.io/gorm"
)

type UpdateTemplateDokumenTambahanCommandHandler struct {
	Repo          domaintemplatedokumentambahan.ITemplateDokumenTambahanRepository
	JenisFileRepo domainjenisfile.IJenisFileRepository
}

func (h *UpdateTemplateDokumenTambahanCommandHandler) Handle(
	ctx context.Context,
	cmd UpdateTemplateDokumenTambahanCommand,
) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// -------------------------
	// VALIDATE UUID
	// -------------------------
	templatedokumentambahanUUID, err := uuid.Parse(cmd.Uuid)
	if err != nil {
		return "", domaintemplatedokumentambahan.InvalidUuid()
	}
	uuidJenisFile, err := uuid.Parse(cmd.JenisFile)
	if err != nil {
		return "", domaintemplatedokumentambahan.JenisFileNotFound()
	}

	// -------------------------
	// GET EXISTING templatedokumentambahan
	// -------------------------
	var (
		jenisfile                       *domainjenisfile.JenisFileDefault
		existingTemplateDokumenTambahan *domaintemplatedokumentambahan.TemplateDokumenTambahan
	)

	g, gctx := errgroup.WithContext(context.Background())

	g.Go(func() error {
		r, err := h.JenisFileRepo.GetDefaultByUuid(gctx, uuidJenisFile)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return domaintemplatedokumentambahan.JenisFileNotFound()
			}
			return err
		}
		jenisfile = r
		return nil
	})

	g.Go(func() error {
		r, err := h.Repo.GetByUuid(ctx, templatedokumentambahanUUID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return domaintemplatedokumentambahan.NotFound(cmd.Uuid)
			}
			return err
		}
		existingTemplateDokumenTambahan = r
		return nil
	})

	if err := g.Wait(); err != nil {
		return "", err
	}

	// -------------------------
	// AGGREGATE ROOT LOGIC
	// -------------------------
	result := domaintemplatedokumentambahan.UpdateTemplateDokumenTambahan(
		existingTemplateDokumenTambahan,
		templatedokumentambahanUUID,
		cmd.Tahun,
		jenisfile.Id,
		cmd.Pertanyaan,
		cmd.Klasifikasi,
		cmd.Kategori,
		cmd.Tugas,
	)

	if !result.IsSuccess {
		return "", result.Error
	}

	updatedTemplateDokumenTambahan := result.Value

	// -------------------------
	// SAVE TO REPOSITORY
	// -------------------------
	if err := h.Repo.Update(ctx, updatedTemplateDokumenTambahan); err != nil {
		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) {
			if mysqlErr.Number == 1062 {
				return "", domaintemplatedokumentambahan.DuplicateData()
			}
		}

		return "", err
	}

	return updatedTemplateDokumenTambahan.UUID.String(), nil
}
