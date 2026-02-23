package application

import (
	"context"
	"os"
	"path/filepath"

	commondomain "UnpakSiamida/common/domain"
	"UnpakSiamida/common/helper"
	commoninfra "UnpakSiamida/common/infrastructure"
	domainkts "UnpakSiamida/modules/kts/domain"
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ExportKtsCommandHandler struct {
	Repo  domainkts.IKtsRepository
	Redis commondomain.IRedisStore
}

func (h *ExportKtsCommandHandler) Handle(
	ctx context.Context,
	cmd ExportKtsCommand,
) ([]byte, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	uuidKts, err := uuid.Parse(cmd.Uuid)
	if err != nil {
		return []byte{}, domainkts.InvalidUuid()
	}

	kts, err := h.Repo.GetDefaultByUuid(ctx, uuidKts)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return []byte{}, domainkts.NotFound(cmd.Uuid)
		}
		return []byte{}, err
	}

	audit := []string{"auditee", "auditor1", "auditor2"}
	audituuid := []string{helper.UUIDString(kts.AuditeeUuid), helper.UUIDString(kts.AuditorUuid)}
	other := []string{"admin", "fakultas"}

	if cmd.SID != "preview" {
		key := "pdf_kts:" + cmd.Uuid + ":" + cmd.Token
		exists, err := h.Redis.Exists(ctx, key)
		if err != nil || !exists {
			return []byte{}, domainkts.NotPushDownload()
		}
	}

	wd, err := os.Getwd()
	if err != nil {
		return []byte{}, domainkts.NoPermission()
	}

	template := filepath.Join(wd, "modules/kts/templates", "kts.html")

	if err := helper.CheckFilesExist(template); err != nil {
		return []byte{}, domainkts.NoResource()
	}

	data := buildKtsData(kts)

	pdfGen := commoninfra.NewWkhtmlPdfGenerator()
	pdf, err := pdfGen.Generate(template, data)
	if err != nil {
		return []byte{}, domainkts.GeneratePDF(err.Error())
	}

	if cmd.SID == "preview" {
		return pdf, nil
	}

	if (helper.GrantedContains(audit, cmd.Tahun, cmd.Granted, false) && helper.Contains(audituuid, cmd.SID) && cmd.Tahun == helper.NullableString(kts.Tahun)) ||
		helper.GrantedContains(other, cmd.Tahun, cmd.Granted, true) {
		return pdf, nil
	}

	return []byte{}, domainkts.NotGranted()
}

func buildKtsData(data *domainkts.KtsDefault) domainkts.KtsPDF {
	ctx := helper.DateContext{}
	ctx.SetStrategy(helper.IndonesianDateFormatter{})

	stz := helper.SanitizerContext{}
	stz.SetStrategy(helper.NewDefaultSanitizer())

	payload := domainkts.KtsPDF{
		Nomor:   helper.StringValue(data.NomorLaporan),
		Tanggal: ctx.FormatDefault(data.TanggalLaporan),
		Auditee: domainkts.Auditee{
			Nama: helper.StringValue(data.NamaAuditee),
		},
		Auditor1: helper.StringValue(data.NamaAuditor),
		AccAuditor1: domainkts.AccAuditor{
			Nama: helper.StringValue(data.NamaAuditor),
		},
		AccAuditor1Final: domainkts.AccAuditor{
			Nama: helper.StringValue(data.NamaAuditor),
			Tgl:  ctx.FormatDefault(data.TanggalClosing),
		},
		P:                         stz.Sanitize(helper.StringValue(data.KetidaksesuaianP)),
		L:                         stz.Sanitize(helper.StringValue(data.KetidaksesuaianL)),
		O:                         stz.Sanitize(helper.StringValue(data.KetidaksesuaianO)),
		R:                         stz.Sanitize(helper.StringValue(data.KetidaksesuaianR)),
		AkarMasalah:               stz.Sanitize(helper.StringValue(data.AkarMasalah)),
		TindakanKoreksi:           stz.Sanitize(helper.StringValue(data.TindakanKoreksi)),
		Referensi:                 helper.StringValue(data.Referensi),
		HasilTemuan:               helper.StringValue(data.HasilTemuan),
		Detail:                    []domainkts.DetailPoint{},
		TindakanPerbaikan:         helper.StringValue(data.TindakanPerbaikan),
		TanggalPenyelesaian:       ctx.FormatDefault(data.TanggalPenyelesaian),
		TinjauanTindakanPerbaikan: helper.StringValue(data.TinjauanTindakanPerbaikan),
		Close: domainkts.CloseData{
			WmmUpmfUpmps:   helper.StringValue(data.WmmUpmfUpmps),
			TanggalClosing: ctx.FormatDefault(data.TanggalClosingFinal),
		},
	}

	if data.AccAuditor != nil && payload.AccAuditor1.Nama != "" { //check AccAuditor
		payload.QRAccAuditor1, _ = helper.GenerateQRBase64("Nama : "+payload.AccAuditor1.Nama, 100)
	}
	if data.StatusAccAuditee != nil && payload.Auditee.Nama != "" { //check StatusAccAuditee
		payload.QRAuditee, _ = helper.GenerateQRBase64("Nama : "+payload.Auditee.Nama, 100)
	}
	if data.AccFinal != nil && payload.AccAuditor1Final.Nama != "" { //check AccFinal
		payload.QRAccAuditor1Final, _ = helper.GenerateQRBase64("Nama : "+payload.AccAuditor1Final.Nama, 100)
	}
	if payload.Close.WmmUpmfUpmps != "" {
		payload.QRClose, _ = helper.GenerateQRBase64("Nama : "+payload.Close.WmmUpmfUpmps, 100)
	}

	return payload
}
