package application

import (
	"context"
	"os"
	"path/filepath"

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
	audituuid := []string{helper.UUIDString(beritaacara.AuditeeUuid), helper.UUIDString(beritaacara.Auditor1Uuid), helper.UUIDString(beritaacara.Auditor2Uuid)}
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

	data := buildBeritaAcaraData(beritaacara.Tahun, beritaacara.FakultasUnit, beritaacara.Auditor1, beritaacara.Auditor2, beritaacara.Auditee, logo1, logo2, beritaacara.Tanggal)

	pdfGen := commoninfra.NewWkhtmlPdfGenerator()
	pdf, err := pdfGen.Generate(template, data)
	if err != nil {
		return []byte{}, domainberitaacara.GeneratePDF(err.Error())
	}

	if cmd.SID == "preview" {
		return pdf, nil
	}

	if (helper.GrantedContains(audit, cmd.Tahun, cmd.Granted, false) && helper.Contains(audituuid, cmd.SID) && cmd.Tahun == data.Tahun) ||
		helper.GrantedContains(other, cmd.Tahun, cmd.Granted, true) {
		return pdf, nil
	}

	return []byte{}, domainberitaacara.NotGranted()
}

func buildBeritaAcaraData(tahun string, target string, auditor1 *string, auditor2 *string, auditee *string, logo1 string, logo2 string, tanggal time.Time) domainberitaacara.BeritaAcaraPDF {
	var qrAuditor = ""
	var auditor = ""
	if auditor1 != nil {
		qrAuditor, _ = helper.GenerateQRBase64("Nama : "+helper.NullableString(auditor1), 120)
		auditor = helper.NullableString(auditor1)
	} else {
		qrAuditor, _ = helper.GenerateQRBase64("Nama : "+helper.NullableString(auditor2), 120)
		auditor = helper.NullableString(auditor2)
	}

	qrAuditee, _ := helper.GenerateQRBase64("Nama : "+helper.NullableString(auditee), 120)

	ctx := helper.DateContext{}
	ctx.SetStrategy(helper.IndonesianDateFormatter{})

	return domainberitaacara.BeritaAcaraPDF{
		Tahun:     tahun,
		Hari:      ctx.NameDay(tanggal),
		Tanggal:   ctx.Format(tanggal),
		Target:    target, //[pr] belum masuk type fakultas, unit, prodi
		Auditor:   auditor,
		Auditor1:  helper.NullableString(auditor1),
		Auditor2:  helper.NullableString(auditor2),
		Auditee:   helper.NullableString(auditee),
		Logo1:     logo1,
		Logo2:     logo2,
		QrAuditee: qrAuditee,
		QrAuditor: qrAuditor,
	}
}
