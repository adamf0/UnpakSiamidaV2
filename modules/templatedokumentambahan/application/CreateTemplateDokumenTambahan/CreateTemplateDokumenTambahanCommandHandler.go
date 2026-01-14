package application

import (
	"context"
	"github.com/google/uuid"
	
	domaintemplatedokumentambahan "UnpakSiamida/modules/templatedokumentambahan/domain"
	domainjenisfile "UnpakSiamida/modules/jenisfile/domain"
	"errors"
    "gorm.io/gorm"
	"encoding/json"
	"fmt"
	"time"
)

type CreateTemplateDokumenTambahanCommandHandler struct {
	Repo                	domaintemplatedokumentambahan.ITemplateDokumenTambahanRepository
	JenisFileRepo    		domainjenisfile.IJenisFileRepository
}

func (h *CreateTemplateDokumenTambahanCommandHandler) Handle(
	ctx context.Context,
	cmd CreateTemplateDokumenTambahanCommand,
) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	
	x, _ := json.MarshalIndent(cmd, "", "  ")
	fmt.Println("DEBUG cmd:", string(x))

	uuidJenisFile, err := uuid.Parse(cmd.JenisFile)
	if err != nil {
		return "", domaintemplatedokumentambahan.JenisFileNotFound()
	}

	jenisfile, err := h.JenisFileRepo.GetDefaultByUuid(ctx, uuidJenisFile)
	b, _ := json.MarshalIndent(jenisfile, "", "  ")
	fmt.Println("DEBUG jenisfile:", string(b))
	
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", domaintemplatedokumentambahan.JenisFileNotFound()
		}

		return "", err;
	}
	
	// c, _ := json.MarshalIndent(jenisfile, "", "  ")
	// fmt.Println("DEBUG jenisfile:", string(c))

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

	// b, _ := json.MarshalIndent(templateDokumenTambahan, "", "  ")
	// fmt.Println("DEBUG templateDokumenTambahan:", string(b))

	// --------------------------
	// SAVE REPOSITORY
	// --------------------------
	if err := h.Repo.Create(ctx, templateDokumenTambahan); err != nil {
		return "", err
	}

	return templateDokumenTambahan.UUID.String(), nil
}
