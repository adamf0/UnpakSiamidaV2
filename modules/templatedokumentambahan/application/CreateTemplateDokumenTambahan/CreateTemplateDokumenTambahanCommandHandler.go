package application

import (
	"context"

	"github.com/google/uuid"

	domainjenisfile "UnpakSiamida/modules/jenisfile/domain"
	domaintemplatedokumentambahan "UnpakSiamida/modules/templatedokumentambahan/domain"
	"errors"
	"time"

	"gorm.io/gorm"
)

type CreateTemplateDokumenTambahanCommandHandler struct {
	Repo          domaintemplatedokumentambahan.ITemplateDokumenTambahanRepository
	JenisFileRepo domainjenisfile.IJenisFileRepository
}

func (h *CreateTemplateDokumenTambahanCommandHandler) Handle(
	ctx context.Context,
	cmd CreateTemplateDokumenTambahanCommand,
) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	uuidJenisFile, err := uuid.Parse(cmd.JenisFile)
	if err != nil {
		return "", domaintemplatedokumentambahan.JenisFileNotFound()
	}

	jenisfile, err := h.JenisFileRepo.GetDefaultByUuid(ctx, uuidJenisFile)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", domaintemplatedokumentambahan.JenisFileNotFound()
		}

		return "", err
	}

	result := domaintemplatedokumentambahan.NewTemplateDokumenTambahan(
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

	templateDokumenTambahan := result.Value

	// --------------------------
	// SAVE REPOSITORY
	// --------------------------
	if err := h.Repo.Create(ctx, templateDokumenTambahan); err != nil {
		// var mysqlErr *mysql.MySQLError
		// if errors.As(err, &mysqlErr) { //gagal masuk kesini
		// 	if mysqlErr.Number == 1062 {
		// 		return "", domaintemplatedokumentambahan.DuplicateData()
		// 	}
		// }
		// if strings.Contains(err.Error(), "Duplicate entry") { //gagal masuk kesini
		// 	return "", domaintemplatedokumentambahan.DuplicateData()
		// }

		return "", err
	}

	return templateDokumenTambahan.UUID.String(), nil
}
