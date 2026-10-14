package presentation

import (
	"bytes"
	"context"
	"io"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/mehdihadeli/go-mediatr"

	commondomain "UnpakSiamida/common/domain"
	"UnpakSiamida/common/helper"
	commoninfra "UnpakSiamida/common/infrastructure"
	commonpresentation "UnpakSiamida/common/presentation"
	Ktsdomain "UnpakSiamida/modules/kts/domain"

	DeleteKts "UnpakSiamida/modules/kts/application/DeleteKts"
	ExportKts "UnpakSiamida/modules/kts/application/ExportKts"
	GenerateKts "UnpakSiamida/modules/kts/application/GenerateKts"
	GetAllKtss "UnpakSiamida/modules/kts/application/GetAllKtss"
	GetKts "UnpakSiamida/modules/kts/application/GetKts"
	GetTemuanKts "UnpakSiamida/modules/kts/application/GetTemuanKts"
	SetupUuidKts "UnpakSiamida/modules/kts/application/SetupUuidKts"
	UpdateKts "UnpakSiamida/modules/kts/application/UpdateKts"

	"strconv"

	"github.com/google/uuid"
)

// =======================================================
// PUT /kts/{uuid}
// =======================================================

// UpdateKtsHandler godoc
// @Summary Update existing Kts
// @Tags Kts
// @Param uuid path string true "Kts UUID" format(uuid)
// @Param nama formData string true "Nama Kts"
// @Produce json
// @Success 200 {object} map[string]string "uuid of updated Kts"
// @Failure 400 {object} commoninfra.ResponseError
// @Failure 404 {object} commoninfra.ResponseError
// @Failure 409 {object} commoninfra.ResponseError
// @Failure 500 {object} commoninfra.ResponseError
// @Router /kts/{uuid} [put]
func UpdateKtsHandlerfunc(c *fiber.Ctx) error {

	uuid := c.Params("uuid")
	tahun := c.Params("tahun")
	step := c.FormValue("step")
	nomorLaporan := helper.StrPtr(c.FormValue("nomorLaporan"))
	tanggalLaporan := helper.StrPtr(c.FormValue("tanggalLaporan"))
	uraianKetidaksesuaianP := helper.StrPtr(c.FormValue("uraianKetidaksesuaianP"))
	uraianKetidaksesuaianL := helper.StrPtr(c.FormValue("uraianKetidaksesuaianL"))
	uraianKetidaksesuaianO := helper.StrPtr(c.FormValue("uraianKetidaksesuaianO"))
	uraianKetidaksesuaianR := helper.StrPtr(c.FormValue("uraianKetidaksesuaianR"))
	akarMasalah := helper.StrPtr(c.FormValue("akarMasalah"))
	tindakanKoreksi := helper.StrPtr(c.FormValue("tindakanKoreksi"))

	statusAccAuditee := helper.StrPtr(c.FormValue("statusAccAuditee"))
	keteranganTolak := helper.StrPtr(c.FormValue("keteranganTolak"))
	tindakanPerbaikan := helper.StrPtr(c.FormValue("tindakanPerbaikan"))

	tanggalpenyelesaian := helper.StrPtr(c.FormValue("tanggalPenyelesaian"))

	tinjauanTindakanPerbaikan := helper.StrPtr(c.FormValue("tinjauanTindakanPerbaikan"))
	tanggalClosing := helper.StrPtr(c.FormValue("tanggalClosing"))

	tanggalClosingFinal := helper.StrPtr(c.FormValue("tanggalClosingFinal"))
	wmmUpmfUpmps := helper.StrPtr(c.FormValue("wmmUpmfUpmps"))

	sid := c.FormValue("sid")

	cmd := UpdateKts.UpdateKtsCommand{
		Uuid:                   uuid,
		NomorLaporan:           nomorLaporan,
		TanggalLaporan:         tanggalLaporan,
		UraianKetidaksesuaianP: uraianKetidaksesuaianP,
		UraianKetidaksesuaianL: uraianKetidaksesuaianL,
		UraianKetidaksesuaianO: uraianKetidaksesuaianO,
		UraianKetidaksesuaianR: uraianKetidaksesuaianR,
		AkarMasalah:            akarMasalah,
		TindakanKoreksi:        tindakanKoreksi,

		StatusAccAuditee:  statusAccAuditee,
		KeteranganTolak:   keteranganTolak,
		TindakanPerbaikan: tindakanPerbaikan,

		TanggalPenyelesaian: tanggalpenyelesaian,

		TinjauanTindakanPerbaikan: tinjauanTindakanPerbaikan,
		TanggalClosing:            tanggalClosing,

		TanggalClosingFinal: tanggalClosingFinal,
		WmmUpmfUpmps:        wmmUpmfUpmps,

		Acc:   sid,
		Tahun: tahun,
		Step:  step,
	}

	updatedID, err := mediatr.Send[UpdateKts.UpdateKtsCommand, string](context.Background(), cmd)
	if err != nil {
		return commoninfra.HandleError(c, err)
	}

	return c.JSON(fiber.Map{"uuid": updatedID})
}

// =======================================================
// GET /Kts/{uuid}
// =======================================================

// GetKtsHandler godoc
// @Summary Get Kts by UUID
// @Tags Kts
// @Param uuid path string true "Kts UUID" format(uuid)
// @Produce json
// @Success 200 {object} Ktsdomain.Kts
// @Failure 400 {object} commoninfra.ResponseError
// @Failure 404 {object} commoninfra.ResponseError
// @Failure 409 {object} commoninfra.ResponseError
// @Failure 500 {object} commoninfra.ResponseError
//
//	@Router /Kts/{uuid} [get]
func GetKtsHandlerfunc(c *fiber.Ctx) error {
	uuid := c.Params("uuid")

	query := GetKts.GetKtsByUuidQuery{Uuid: uuid}

	Kts, err := mediatr.Send[
		GetKts.GetKtsByUuidQuery,
		*Ktsdomain.Kts,
	](context.Background(), query)

	if err != nil {
		return commoninfra.HandleError(c, err)
	}

	if Kts == nil {
		return c.Status(404).JSON(fiber.Map{"error": "Kts not found"})
	}

	return c.JSON(Kts)
}

// =======================================================
// GET /Ktss
// =======================================================

// GetAllKtssHandler godoc
// @Summary Get all Ktss
// @Tags Kts
// @Param mode query string false "paging | all | ndjson | sse"
// @Param page query int false "Page number"
// @Param limit query int false "Limit per page"
// @Param search query string false "Search keyword"
// @Produce json
// @Success 200 {object} commondomain.Paged[Ktsdomain.KtsDefault]
// @Router /Ktss [get]
func GetAllKtssHandlerfunc(c *fiber.Ctx) error {
	mode := c.Query("mode", "paging")
	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 10)
	search := c.Query("search", "")

	filtersRaw := c.Query("filters", "")
	var filters []commondomain.SearchFilter

	if filtersRaw != "" {
		parts := strings.Split(filtersRaw, ";")
		for _, p := range parts {
			tokens := strings.SplitN(p, ":", 3)
			if len(tokens) != 3 {
				continue
			}

			field := strings.TrimSpace(tokens[0])
			op := strings.TrimSpace(tokens[1])
			rawValue := strings.TrimSpace(tokens[2])

			var ptr *string
			if rawValue != "" && rawValue != "null" {
				ptr = &rawValue
			}

			filters = append(filters, commondomain.SearchFilter{
				Field:    field,
				Operator: op,
				Value:    ptr,
			})
		}
	}

	query := GetAllKtss.GetAllKtssQuery{
		Search:        search,
		SearchFilters: filters,
	}

	var adapter commonpresentation.OutputAdapter[Ktsdomain.KtsDefault]
	switch mode {
	case "all":
		adapter = &commonpresentation.AllAdapter[Ktsdomain.KtsDefault]{}
	case "ndjson":
		adapter = &commonpresentation.NDJSONAdapter[Ktsdomain.KtsDefault]{}
	case "sse":
		adapter = &commonpresentation.SSEAdapter[Ktsdomain.KtsDefault]{}
	default:
		query.Page = &page
		query.Limit = &limit
		adapter = &commonpresentation.PagingAdapter[Ktsdomain.KtsDefault]{}
	}

	result, err := mediatr.Send[
		GetAllKtss.GetAllKtssQuery,
		commondomain.Paged[Ktsdomain.KtsDefault],
	](context.Background(), query)

	if err != nil {
		return commoninfra.HandleError(c, err)
	}

	return adapter.Send(c, result)
}

// =======================================================
// DELETE /kts/{uuid}
// =======================================================

// DeleteKtsHandler godoc
// @Summary Delete an Kts
// @Tags Kts
// @Param uuid path string true "Kts UUID" format(uuid)
// @Produce json
// @Success 200 {object} map[string]string "uuid of deleted Kts"
// @Failure 400 {object} commoninfra.ResponseError
// @Failure 404 {object} commoninfra.ResponseError
// @Failure 409 {object} commoninfra.ResponseError
// @Failure 500 {object} commoninfra.ResponseError
//
//	@Router /kts/{uuid} [delete]
func DeleteKtsHandlerfunc(c *fiber.Ctx) error {
	uuid := c.Params("uuid")

	cmd := DeleteKts.DeleteKtsCommand{Uuid: uuid}

	deletedID, err := mediatr.Send[
		DeleteKts.DeleteKtsCommand,
		string,
	](context.Background(), cmd)

	if err != nil {
		return commoninfra.HandleError(c, err)
	}

	return c.JSON(fiber.Map{"uuid": deletedID})
}

// PublishKtsHandler godoc
// @Summary Publish existing Kts
// @Tags Kts
// @Param uuid path string true "Kts UUID" format(uuid)
// @Produce json
// @Success 200 {object} map[string]string "uuid of export Kts"
// @Failure 400 {object} commoninfra.ResponseError
// @Failure 404 {object} commoninfra.ResponseError
// @Failure 409 {object} commoninfra.ResponseError
// @Failure 500 {object} commoninfra.ResponseError
// @Router /beritaacara/{uuid} [post]
func PublishKtsHandlerfunc(c *fiber.Ctx) error {

	uuid := c.Params("uuid")
	token := c.FormValue("token")

	cmd := ExportKts.PublishKtsCommand{
		Uuid:  uuid,
		Token: token,
	}

	exportID, err := mediatr.Send[ExportKts.PublishKtsCommand, string](context.Background(), cmd)
	if err != nil {
		return commoninfra.HandleError(c, err)
	}

	return c.JSON(fiber.Map{"uuid": exportID})
}

// PreviewKtsHandler godoc
// @Summary Publish existing Kts
// @Tags Kts
// @Param uuid path string true "Kts UUID" format(uuid)
// @Produce json
// @Success 200 {object} map[string]string "uuid of export Kts"
// @Failure 400 {object} commoninfra.ResponseError
// @Failure 404 {object} commoninfra.ResponseError
// @Failure 409 {object} commoninfra.ResponseError
// @Failure 500 {object} commoninfra.ResponseError
// @Router /beritaacara/preview/:uuid [get]
func PreviewKtsHandlerfunc(c *fiber.Ctx) error {

	uuid := c.Params("uuid")
	token := c.Query("token")
	tahun := c.Query("ctxtahun")
	// sid := c.FormValue("sid")
	granted := c.FormValue("grantedaccess")

	cmd := ExportKts.ExportKtsCommand{
		Uuid:    uuid,
		Token:   token,
		SID:     "preview",
		Granted: granted,
		Tahun:   tahun,
	}

	data, err := mediatr.Send[ExportKts.ExportKtsCommand, []byte](context.Background(), cmd)
	if err != nil {
		return commoninfra.HandleError(c, err)
	}

	c.Set("Content-Type", "application/pdf")
	c.Set("Content-Disposition", "inline; filename=kts.pdf")
	c.Set("Cache-Control", "private, max-age=60")
	c.Set("X-Accel-Buffering", "no") // penting untuk nginx/Cloudflare agar streaming langsung
	return c.Send(data)
}

// ExportKtsHandler godoc
// @Summary Publish existing Kts
// @Tags Kts
// @Param uuid path string true "Kts UUID" format(uuid)
// @Produce json
// @Success 200 {object} map[string]string "uuid of export Kts"
// @Failure 400 {object} commoninfra.ResponseError
// @Failure 404 {object} commoninfra.ResponseError
// @Failure 409 {object} commoninfra.ResponseError
// @Failure 500 {object} commoninfra.ResponseError
// @Router /beritaacara/export/:uuid [get]
func ExportKtsHandlerfunc(c *fiber.Ctx) error {

	uuid := c.Params("uuid")
	token := c.FormValue("token")
	tahun := c.Query("ctxtahun")
	sid := c.FormValue("sid")
	granted := c.FormValue("grantedaccess")

	cmd := ExportKts.ExportKtsCommand{
		Uuid:    uuid,
		Token:   token,
		SID:     sid,
		Granted: granted,
		Tahun:   tahun,
	}

	data, err := mediatr.Send[ExportKts.ExportKtsCommand, []byte](context.Background(), cmd)
	if err != nil {
		return commoninfra.HandleError(c, err)
	}

	c.Set("Content-Type", "application/pdf")
	c.Set("Content-Disposition", "attachment; filename=berita_acara.pdf")
	c.Set("Cache-Control", "private, max-age=60")
	c.Set("X-Accel-Buffering", "no")

	dataReader := bytes.NewReader(data)
	buf := make([]byte, 32*1024) // 32 KB per chunk

	for {
		n, err := dataReader.Read(buf)
		if n > 0 {
			if _, writeErr := c.Write(buf[:n]); writeErr != nil {
				return writeErr
			}
		}
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
	}

	return nil
}

func SetupUuidKtssHandlerfunc(c *fiber.Ctx) error {
	cmd := SetupUuidKts.SetupUuidKtsCommand{}

	message, err := mediatr.Send[SetupUuidKts.SetupUuidKtsCommand, string](context.Background(), cmd)
	if err != nil {
		return commoninfra.HandleError(c, err)
	}

	return c.JSON(fiber.Map{"message": message})
}

// GetTemuanKtsHandlerfunc godoc
// @Summary Get Temuan Kts for Renstra
// @Tags Kts
// @Param renstra_uuid path string false "Renstra UUID"
// @Param tugas query string false "auditor1 | auditor2"
// @Param tahun query string false "Tahun"
// @Produce json
// @Router /kts/temuan/{renstra_uuid} [get]
func GetTemuanKtsHandlerfunc(c *fiber.Ctx) error {
	renstraUuid := c.Params("renstra_uuid")
	if renstraUuid == "" {
		renstraUuid = c.Query("renstra_uuid")
	}
	tugas := c.Query("tugas")
	tahun := c.Query("tahun")

	query := GetTemuanKts.GetTemuanKtsQuery{
		RenstraUUID: renstraUuid,
		Tugas:       tugas,
		Tahun:       tahun,
	}

	results, err := mediatr.Send[
		GetTemuanKts.GetTemuanKtsQuery,
		[]Ktsdomain.TemuanKts,
	](context.Background(), query)

	if err != nil {
		return commoninfra.HandleError(c, err)
	}

	return c.JSON(results)
}

// GenerateKtsHandlerfunc godoc
// @Summary Generate Kts from Temuan
// @Tags Kts
// @Param renstra_nilai_id formData int false "ID Renstra Nilai"
// @Param dokumen_tambahan_id formData int false "ID Dokumen Tambahan"
// @Param auditor_id formData string false "Auditor ID"
// @Produce json
// @Success 200 {object} map[string]string "uuid of generated Kts"
// @Router /kts/generate [post]
func GenerateKtsHandlerfunc(c *fiber.Ctx) error {
	rnInput := strings.TrimSpace(c.FormValue("renstra_nilai_uuid"))
	if rnInput == "" {
		rnInput = strings.TrimSpace(c.FormValue("renstra_nilai_id"))
	}
	if rnInput == "" {
		rnInput = strings.TrimSpace(c.FormValue("renstra_nilai"))
	}

	dtInput := strings.TrimSpace(c.FormValue("dokumen_tambahan_uuid"))
	if dtInput == "" {
		dtInput = strings.TrimSpace(c.FormValue("dokumen_tambahan_id"))
	}
	if dtInput == "" {
		dtInput = strings.TrimSpace(c.FormValue("dokumen_tambahan"))
	}

	var rnIDPtr *uint
	var rnUUIDPtr *string
	if rnInput != "" {
		if val, err := strconv.ParseUint(rnInput, 10, 32); err == nil && val > 0 {
			uVal := uint(val)
			rnIDPtr = &uVal
		} else {
			rnUUIDPtr = &rnInput
		}
	}

	var dtIDPtr *uint
	var dtUUIDPtr *string
	if dtInput != "" {
		if val, err := strconv.ParseUint(dtInput, 10, 32); err == nil && val > 0 {
			uVal := uint(val)
			dtIDPtr = &uVal
		} else {
			dtUUIDPtr = &dtInput
		}
	}

	auditorID := strings.TrimSpace(c.FormValue("auditor_id"))
	if auditorID == "" {
		auditorID = strings.TrimSpace(c.FormValue("sid"))
	}

	cmd := GenerateKts.GenerateKtsCommand{
		RenstraNilaiID:      rnIDPtr,
		RenstraNilaiUUID:    rnUUIDPtr,
		DokumenTambahanID:   dtIDPtr,
		DokumenTambahanUUID: dtUUIDPtr,
		AuditorID:           auditorID,
	}

	resUUID, err := mediatr.Send[
		GenerateKts.GenerateKtsCommand,
		*uuid.UUID,
	](context.Background(), cmd)

	if err != nil {
		return commoninfra.HandleError(c, err)
	}

	return c.JSON(fiber.Map{
		"uuid": resUUID.String(),
	})
}

func ModuleKts(app *fiber.App) {
	admin := []string{"admin"}
	audit := []string{"auditee", "auditor1", "auditor2"}
	whoamiURL := "http://localhost:3000/whoami"

	app.Get("/kts/setupuuid", SetupUuidKtssHandlerfunc)
	app.Delete("/ks/:uuid", commonpresentation.JWTMiddleware(), commonpresentation.RBACMiddleware(admin, whoamiURL), DeleteKtsHandlerfunc)

	app.Put("/kts/:tahun/:uuid", commonpresentation.JWTMiddleware(), commonpresentation.RBACMiddleware(audit, whoamiURL), UpdateKtsHandlerfunc)

	// Sub-routes must be registered BEFORE generic /Kts/:uuid parameterized route
	app.Get("/kts/temuan", commonpresentation.SmartCompress(), commonpresentation.JWTMiddleware(), GetTemuanKtsHandlerfunc)
	app.Get("/kts/temuan/:renstra_uuid", commonpresentation.SmartCompress(), commonpresentation.JWTMiddleware(), GetTemuanKtsHandlerfunc)
	app.Post("/kts/generate", commonpresentation.JWTMiddleware(), commonpresentation.RBACMiddleware(audit, whoamiURL), GenerateKtsHandlerfunc)

	app.Get("/kts/publish/:uuid", commonpresentation.JWTMiddleware(), PublishKtsHandlerfunc)
	app.Get("/kts/preview/:uuid", PreviewKtsHandlerfunc) //private network
	app.Get("/kts/export/:uuid", commonpresentation.JWTMiddleware(), commonpresentation.RBACMiddleware(audit, whoamiURL), ExportKtsHandlerfunc)

	app.Get("/Ktss", commonpresentation.SmartCompress(), commonpresentation.JWTMiddleware(), GetAllKtssHandlerfunc)
	app.Get("/Kts/:uuid", commonpresentation.SmartCompress(), commonpresentation.JWTMiddleware(), GetKtsHandlerfunc)
}
