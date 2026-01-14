package presentation

import (
	"context"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/mehdihadeli/go-mediatr"

	commondomain "UnpakSiamida/common/domain"
	commoninfra "UnpakSiamida/common/infrastructure"
	commonpresentation "UnpakSiamida/common/presentation"
	Ktsdomain "UnpakSiamida/modules/kts/domain"

	DeleteKts "UnpakSiamida/modules/kts/application/DeleteKts"
	GetAllKtss "UnpakSiamida/modules/kts/application/GetAllKtss"
	GetKts "UnpakSiamida/modules/kts/application/GetKts"
	SetupUuidKts "UnpakSiamida/modules/kts/application/SetupUuidKts"
	UpdateKts "UnpakSiamida/modules/kts/application/UpdateKts"
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
// @Failure 400 {object} commondomain.Error
// @Router /kts/{uuid} [put]
func UpdateKtsHandlerfunc(c *fiber.Ctx) error {

	uuid := c.Params("uuid")
	tahun := c.Params("tahun")
	step := c.FormValue("step")
	nomorLaporan := ptr(c.FormValue("nomorLaporan"))
	tanggalLaporan := ptr(c.FormValue("tanggalLaporan"))
	uraianKetidaksesuaianP := ptr(c.FormValue("uraianKetidaksesuaianP"))
	uraianKetidaksesuaianL := ptr(c.FormValue("uraianKetidaksesuaianL"))
	uraianKetidaksesuaianO := ptr(c.FormValue("uraianKetidaksesuaianO"))
	uraianKetidaksesuaianR := ptr(c.FormValue("uraianKetidaksesuaianR"))
	akarMasalah := ptr(c.FormValue("akarMasalah"))
	tindakanKoreksi := ptr(c.FormValue("tindakanKoreksi"))

	statusAccAuditee := ptr(c.FormValue("statusAccAuditee"))
	keteranganTolak := ptr(c.FormValue("keteranganTolak"))
	tindakanPerbaikan := ptr(c.FormValue("tindakanPerbaikan"))

	tanggalpenyelesaian := ptr(c.FormValue("tanggalPenyelesaian"))

	tinjauanTindakanPerbaikan := ptr(c.FormValue("tinjauanTindakanPerbaikan"))
	tanggalClosing := ptr(c.FormValue("tanggalClosing"))

	tanggalClosingFinal := ptr(c.FormValue("tanggalClosingFinal"))
	wmmUpmfUpmps := ptr(c.FormValue("wmmUpmfUpmps"))

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
// @Failure 404 {object} commondomain.Error
// @Router /Kts/{uuid} [get]
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
// @Success 200 {object} Ktsdomain.PagedKtss
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

	var adapter OutputAdapter
	switch mode {
	case "all":
		adapter = &AllAdapter{}
	case "ndjson":
		adapter = &NDJSONAdapter{}
	case "sse":
		adapter = &SSEAdapter{}
	default:
		query.Page = &page
		query.Limit = &limit
		adapter = &PagingAdapter{}
	}

	result, err := mediatr.Send[
		GetAllKtss.GetAllKtssQuery,
		Ktsdomain.PagedKtss,
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
// @Failure 404 {object} commondomain.Error
// @Router /kts/{uuid} [delete]
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

func SetupUuidKtssHandlerfunc(c *fiber.Ctx) error {
	cmd := SetupUuidKts.SetupUuidKtsCommand{}

	message, err := mediatr.Send[SetupUuidKts.SetupUuidKtsCommand, string](context.Background(), cmd)
	if err != nil {
		return commoninfra.HandleError(c, err)
	}

	return c.JSON(fiber.Map{"message": message})
}

func ModuleKts(app *fiber.App) {
	admin := []string{"admin"}
	audit := []string{"auditee", "auditor1", "auditor2"}
	whoamiURL := "http://localhost:3000/whoami"

	app.Get("/kts/setupuuid", commonpresentation.JWTMiddleware(), commonpresentation.RBACMiddleware(admin, whoamiURL), SetupUuidKtssHandlerfunc)
	app.Delete("/ks/:uuid", commonpresentation.JWTMiddleware(), commonpresentation.RBACMiddleware(admin, whoamiURL), DeleteKtsHandlerfunc)

	app.Put("/kts/:tahun/:uuid", commonpresentation.JWTMiddleware(), commonpresentation.RBACMiddleware(audit, whoamiURL), UpdateKtsHandlerfunc)

	app.Get("/Kts/:uuid", commonpresentation.JWTMiddleware(), GetKtsHandlerfunc)
	app.Get("/Ktss", commonpresentation.JWTMiddleware(), GetAllKtssHandlerfunc)
}

func ptr(s string) *string {
	if strings.TrimSpace(s) == "" {
		return nil
	}
	return &s
}
