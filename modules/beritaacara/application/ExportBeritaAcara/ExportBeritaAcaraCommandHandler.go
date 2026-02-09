package application

import (
	"context"
	"os"
	"path/filepath"
	"strings"

	commondomain "UnpakSiamida/common/domain"
	"UnpakSiamida/common/helper"
	commoninfra "UnpakSiamida/common/infrastructure"
	domainberitaacara "UnpakSiamida/modules/beritaacara/domain"
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ExportBeritaAcaraCommandHandler struct {
	Repo  domainberitaacara.IBeritaAcaraRepository
	Redis commondomain.IRedisStore
}

func (h *ExportBeritaAcaraCommandHandler) Handle(
	ctx context.Context,
	cmd ExportBeritaAcaraCommand,
) ([]byte, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	uuidBeritaAcara, err := uuid.Parse(cmd.Uuid)
	if err != nil {
		return []byte{}, domainberitaacara.InvalidUuid()
	}

	beritaacara, err := h.Repo.GetDefaultByUuid(ctx, uuidBeritaAcara)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return []byte{}, domainberitaacara.NotFound(cmd.Uuid)
		}
		return []byte{}, err
	}

	audit := []string{"auditee", "auditor1", "auditor2"}
	audituuid := []string{beritaacara.AuditeeUuid.String(), beritaacara.Auditor1Uuid.String(), beritaacara.Auditor2Uuid.String()}
	other := []string{"admin", "fakultas"}

	if cmd.SID != "preview" {
		key := "pdf_berita_acara:" + cmd.Uuid + ":" + cmd.Token
		exists, err := h.Redis.Exists(ctx, key)
		if err != nil || !exists {
			return []byte{}, domainberitaacara.NotPushDownload()
		}
	}

	wd, err := os.Getwd()
	if err != nil {
		return []byte{}, domainberitaacara.NoPermission()
	}

	logo1 := filepath.Join(wd, "assets", "Logo1.jpeg")
	logo2 := filepath.Join(wd, "assets", "Logo2.jpeg")
	template := filepath.Join(wd, "modules/beritaacara/templates", "berita_acara.html")

	if err := helper.CheckFilesExist(logo1, logo2, template); err != nil {
		return []byte{}, domainberitaacara.NoResource()
	}

	data := buildBeritaAcaraData(beritaacara.Tahun, beritaacara.FakultasUnit, *beritaacara.Auditor1, *beritaacara.Auditor2, *beritaacara.Auditee, logo1, logo2, beritaacara.Tanggal)

	pdfGen := commoninfra.NewWkhtmlPdfGenerator()
	pdf, err := pdfGen.Generate(template, data)
	if err != nil {
		return []byte{}, domainberitaacara.GeneratePDF(err.Error())
	}

	if (grantedContains(audit, cmd.Tahun, cmd.Granted, false) && contains(audituuid, cmd.SID)) || grantedContains(other, cmd.Tahun, cmd.Granted, true) || cmd.SID == "preview" {
		return pdf, nil
	}

	return []byte{}, domainberitaacara.NotGranted()
}

func buildBeritaAcaraData(tahun, target, auditor1, auditor2, auditee, logo1, logo2 string, tanggal time.Time) domainberitaacara.BeritaAcaraPDF {
	qrAuditor, _ := helper.GenerateQRBase64("Nama : "+auditor1, 120)
	qrAuditee, _ := helper.GenerateQRBase64("Nama : "+auditee, 120)

	var hariID = map[time.Weekday]string{
		time.Monday:    "Senin",
		time.Tuesday:   "Selasa",
		time.Wednesday: "Rabu",
		time.Thursday:  "Kamis",
		time.Friday:    "Jumat",
		time.Saturday:  "Sabtu",
		time.Sunday:    "Minggu",
	}

	return domainberitaacara.BeritaAcaraPDF{
		Tahun:     tahun,
		Hari:      hariID[tanggal.Weekday()],
		Tanggal:   tanggal.Format("02 Januari 2006"),
		Target:    target,
		Auditor1:  auditor1,
		Auditor2:  auditor2,
		Auditee:   auditee,
		Logo1:     logo1,
		Logo2:     logo2,
		QrAuditee: qrAuditee,
		QrAuditor: qrAuditor,
	}
}

func contains(list []string, value string) bool {
	for _, v := range list {
		if v == value {
			return true
		}
	}
	return false
}

func grantedContains(audit []string, tahun string, granted string, isother bool) bool {
	entries := strings.Split(granted, ",")
	for _, e := range entries {
		parts := strings.Split(e, "#")
		if len(parts) != 2 {
			return false
		}

		year := strings.TrimSpace(parts[0])
		level := strings.TrimSpace(parts[1])

		if isother {
			if contains(audit, level) {
				return true
			}
		} else {
			if year == tahun && contains(audit, level) {
				return true
			}
		}
	}
	return false
}
