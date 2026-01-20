package presentation

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/mehdihadeli/go-mediatr"

	commondomain "UnpakSiamida/common/domain"
	commoninfra "UnpakSiamida/common/infrastructure"
	commonpresentation "UnpakSiamida/common/presentation"
	renstradomain "UnpakSiamida/modules/renstra/domain"

	CreateRenstra "UnpakSiamida/modules/renstra/application/CreateRenstra"
	DeleteRenstra "UnpakSiamida/modules/renstra/application/DeleteRenstra"
	GetAllRenstras "UnpakSiamida/modules/renstra/application/GetAllRenstras"
	GetRenstraDefault "UnpakSiamida/modules/renstra/application/GetRenstraDefault"
	GiveCodeRenstra "UnpakSiamida/modules/renstra/application/GiveCodeAccessRenstra"
	SetupUuidRenstra "UnpakSiamida/modules/renstra/application/SetupUuidRenstra"
	UpdateRenstra "UnpakSiamida/modules/renstra/application/UpdateRenstra"
)

// =======================================================
// POST /renstra
// =======================================================

// CreateRenstraHandler godoc
// @Summary Create new Renstra
// @Tags Renstra
// @Param tahun formData string true "Tahun"
// @Param fakultas_unit formData string true "Fakultas/Unit ID"
// @Param periode_upload_mulai formData string true "Periode Upload Mulai"
// @Param periode_upload_akhir formData string true "Periode Upload Akhir"
// @Param periode_assesment_dokumen_mulai formData string true "Periode Assesment Dokumen Mulai"
// @Param periode_assesment_dokumen_akhir formData string true "Periode Assesment Dokumen Akhir"
// @Param periode_assesment_lapangan_mulai formData string true "Periode Assesment Lapangan Mulai"
// @Param periode_assesment_lapangan_akhir formData string true "Periode Assesment Lapangan Akhir"
// @Param auditee formData string true "Auditee"
// @Param auditor1 formData string true "Auditor 1"
// @Param auditor2 formData string true "Auditor 2"
// @Produce json
// @Success 200 {object} map[string]string "uuid of created Renstra"
// @Failure 400 {object} commoninfra.ResponseError
// @Failure 404 {object} commoninfra.ResponseError
// @Failure 409 {object} commoninfra.ResponseError
// @Failure 500 {object} commoninfra.ResponseError
// @Router /renstra [post]
func CreateRenstraHandlerfunc(c *fiber.Ctx) error {
	cmd := CreateRenstra.CreateRenstraCommand{
		Tahun:                         c.FormValue("tahun"),
		FakultasUnit:                  c.FormValue("fakultas_unit"),
		PeriodeUploadMulai:            c.FormValue("periode_upload_mulai"),
		PeriodeUploadAkhir:            c.FormValue("periode_upload_akhir"),
		PeriodeAssesmentDokumenMulai:  c.FormValue("periode_assesment_dokumen_mulai"),
		PeriodeAssesmentDokumenAkhir:  c.FormValue("periode_assesment_dokumen_akhir"),
		PeriodeAssesmentLapanganMulai: c.FormValue("periode_assesment_lapangan_mulai"),
		PeriodeAssesmentLapanganAkhir: c.FormValue("periode_assesment_lapangan_akhir"),
		Auditee:                       c.FormValue("auditee"),
		Auditor1:                      c.FormValue("auditor1"),
		Auditor2:                      c.FormValue("auditor2"),
	}

	// Kirim ke mediator
	uuid, err := mediatr.Send[
		CreateRenstra.CreateRenstraCommand,
		string,
	](context.Background(), cmd)

	if err != nil {
		return commoninfra.HandleError(c, err)
	}

	return c.JSON(fiber.Map{"uuid": uuid})
}

// =======================================================
// PUT /renstra/{uuid}
// =======================================================

// UpdateRenstraHandler godoc
// @Summary Update existing Renstra
// @Tags Renstra
// @Param uuid path string true "Renstra UUID" format(uuid)
// @Param tahun formData string true "Tahun"
// @Param fakultas_unit formData string true "Fakultas/Unit ID"
// @Param periode_upload_mulai formData string true "Periode Upload Mulai"
// @Param periode_upload_akhir formData string true "Periode Upload Akhir"
// @Param periode_assesment_dokumen_mulai formData string true "Periode Assesment Dokumen Mulai"
// @Param periode_assesment_dokumen_akhir formData string true "Periode Assesment Dokumen Akhir"
// @Param periode_assesment_lapangan_mulai formData string true "Periode Assesment Lapangan Mulai"
// @Param periode_assesment_lapangan_akhir formData string true "Periode Assesment Lapangan Akhir"
// @Param auditee formData string true "Auditee"
// @Param auditor1 formData string true "Auditor 1"
// @Param auditor2 formData string true "Auditor 2"
// @Produce json
// @Success 200 {object} map[string]string "uuid of updated Renstra"
// @Failure 400 {object} commoninfra.ResponseError
// @Failure 404 {object} commoninfra.ResponseError
// @Failure 409 {object} commoninfra.ResponseError
// @Failure 500 {object} commoninfra.ResponseError
// @Router /renstra/{uuid} [put]
func UpdateRenstraHandlerfunc(c *fiber.Ctx) error {
	cmd := UpdateRenstra.UpdateRenstraCommand{
		Uuid:                          c.Params("uuid"),
		Tahun:                         c.FormValue("tahun"),
		FakultasUnit:                  c.FormValue("fakultas_unit"),
		PeriodeUploadMulai:            c.FormValue("periode_upload_mulai"),
		PeriodeUploadAkhir:            c.FormValue("periode_upload_akhir"),
		PeriodeAssesmentDokumenMulai:  c.FormValue("periode_assesment_dokumen_mulai"),
		PeriodeAssesmentDokumenAkhir:  c.FormValue("periode_assesment_dokumen_akhir"),
		PeriodeAssesmentLapanganMulai: c.FormValue("periode_assesment_lapangan_mulai"),
		PeriodeAssesmentLapanganAkhir: c.FormValue("periode_assesment_lapangan_akhir"),
		Auditee:                       c.FormValue("auditee"),
		Auditor1:                      c.FormValue("auditor1"),
		Auditor2:                      c.FormValue("auditor2"),
	}

	updatedID, err := mediatr.Send[
		UpdateRenstra.UpdateRenstraCommand,
		string,
	](context.Background(), cmd)

	if err != nil {
		return commoninfra.HandleError(c, err)
	}

	return c.JSON(fiber.Map{"uuid": updatedID})
}

// =======================================================
// PUT /renstra/{uuid}/code_access
// =======================================================

// GiveCodeAccessRenstraHandler godoc
// @Summary Update Renstra Kode Akses
// @Tags Renstra
// @Param uuid path string true "Renstra UUID" format(uuid)
// @Param kodeAkses formData string true "Kode Akses baru"
// @Produce json
// @Success 200 {object} map[string]string "uuid of updated Renstra"
// @Failure 400 {object} commoninfra.ResponseError
// @Failure 404 {object} commoninfra.ResponseError
// @Failure 409 {object} commoninfra.ResponseError
// @Failure 500 {object} commoninfra.ResponseError
// @Router /renstra/{uuid}/code_access [put]
func GiveCodeAccessRenstraHandlerfunc(c *fiber.Ctx) error {
	cmd := GiveCodeRenstra.GiveCodeAccessRenstraCommand{
		Uuid:      c.Params("uuid"),
		KodeAkses: c.FormValue("kodeAkses"),
	}

	updatedID, err := mediatr.Send[
		GiveCodeRenstra.GiveCodeAccessRenstraCommand,
		string,
	](context.Background(), cmd)

	if err != nil {
		return commoninfra.HandleError(c, err)
	}

	return c.JSON(fiber.Map{"uuid": updatedID})
}

// =======================================================
// DELETE /renstra/{uuid}
// =======================================================

// DeleteRenstraHandler godoc
// @Summary Delete an Renstra
// @Tags Renstra
// @Param uuid path string true "Renstra UUID" format(uuid)
// @Produce json
// @Success 200 {object} map[string]string "uuid of deleted Renstra"
// @Failure 404 {object} commondomain.Error
// @Router /renstra/{uuid} [delete]
func DeleteRenstraHandlerfunc(c *fiber.Ctx) error {
	uuid := c.Params("uuid")

	cmd := DeleteRenstra.DeleteRenstraCommand{Uuid: uuid}

	deletedID, err := mediatr.Send[
		DeleteRenstra.DeleteRenstraCommand,
		string,
	](context.Background(), cmd)

	if err != nil {
		return commoninfra.HandleError(c, err)
	}

	return c.JSON(fiber.Map{"uuid": deletedID})
}

// =======================================================
// GET /renstra/{uuid}
// =======================================================

// GetRenstraHandler godoc
// @Summary Get Renstra by UUID
// @Tags Renstra
// @Param uuid path string true "Renstra UUID" format(uuid)
// @Produce json
// @Success 200 {object} renstradomain.RenstraDefault
// @Failure 404 {object} commondomain.Error
// @Router /renstra/{uuid} [get]
func GetRenstraHandlerfunc(c *fiber.Ctx) error {
	uuid := c.Params("uuid")

	query := GetRenstraDefault.GetRenstraDefaultByUuidQuery{Uuid: uuid}

	renstra, err := mediatr.Send[
		GetRenstraDefault.GetRenstraDefaultByUuidQuery,
		*renstradomain.RenstraDefault,
	](context.Background(), query)

	if err != nil {
		return commoninfra.HandleError(c, err)
	}

	if renstra == nil {
		return c.Status(404).JSON(fiber.Map{"error": "Renstra not found"})
	}

	return c.JSON(renstra)
}

// =======================================================
// GET /renstra/{tahun}
// =======================================================

// GetRenstraHandler godoc
// @Summary Get Renstra by Tahun
// @Tags Renstra
// @Param tahun path string true "Tahun"
// @Produce json
// @Success 200 {object} renstradomain.RenstraDefault
// @Failure 404 {object} commondomain.Error
// @Router /renstra/{tahun} [get]
func GetRenstraByTahunTargetHandlerfunc(c *fiber.Ctx) error {
	mode := c.Query("mode", "paging")
	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 10)
	search := c.Query("search", "")

	filtersRaw := c.Query("filters", "")
	tahun := c.Params("tahun")
	if tahun != "" {
		if filtersRaw != "" {
			filtersRaw += ";"
		}
		filtersRaw += fmt.Sprintf("tahun:eq:%s", tahun)
	}

	sid := c.FormValue("sid")
	if sid != "" {
		if filtersRaw != "" {
			filtersRaw += ";"
		}
		filtersRaw += fmt.Sprintf("uuidauditee:eq:%s;uuidauditor1:eq:%s;uuidauditor2:eq:%s", sid, sid, sid)
	}
	log.Printf("[filter] audit: %s", filtersRaw)

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

	query := GetAllRenstras.GetAllRenstrasQuery{
		Scope:         "audit",
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
		GetAllRenstras.GetAllRenstrasQuery,
		renstradomain.PagedRenstras,
	](context.Background(), query)

	if err != nil {
		return commoninfra.HandleError(c, err)
	}

	return adapter.Send(c, result)
}

// GetAllRenstrasHandler godoc
// @Summary Get All Renstras
// @Tags Renstra
// @Param mode query string false "paging | all | ndjson | sse"
// @Param page query int false "Page number"
// @Param limit query int false "Limit per page"
// @Param search query string false "Search keyword"
// @Produce json
// @Success 200 {object} renstradomain.PagedRenstras
// @Router /renstras [get]
func GetAllRenstrasHandlerfunc(c *fiber.Ctx) error {
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

	query := GetAllRenstras.GetAllRenstrasQuery{
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
		GetAllRenstras.GetAllRenstrasQuery,
		renstradomain.PagedRenstras,
	](context.Background(), query)

	if err != nil {
		return commoninfra.HandleError(c, err)
	}

	return adapter.Send(c, result)
}

func SetupUuidRenstrasHandlerfunc(c *fiber.Ctx) error {
	cmd := SetupUuidRenstra.SetupUuidRenstraCommand{}

	message, err := mediatr.Send[SetupUuidRenstra.SetupUuidRenstraCommand, string](context.Background(), cmd)
	if err != nil {
		return commoninfra.HandleError(c, err)
	}

	return c.JSON(fiber.Map{"message": message})
}

func ModuleRenstra(app *fiber.App) {
	admin := []string{"admin"}
	audit := []string{"auditee", "auditor1", "auditor2"}
	whoamiURL := "http://localhost:3000/whoami"

	app.Get("/renstra/setupuuid", commonpresentation.JWTMiddleware(), commonpresentation.RBACMiddleware(admin, whoamiURL), SetupUuidRenstrasHandlerfunc)

	app.Post("/renstra", commonpresentation.JWTMiddleware(), commonpresentation.RBACMiddleware(admin, whoamiURL), CreateRenstraHandlerfunc)
	app.Put("/renstra/:uuid", commonpresentation.JWTMiddleware(), commonpresentation.RBACMiddleware(admin, whoamiURL), UpdateRenstraHandlerfunc)
	app.Put("/renstra/:uuid/code_access", commonpresentation.JWTMiddleware(), commonpresentation.RBACMiddleware(admin, whoamiURL), GiveCodeAccessRenstraHandlerfunc)
	app.Delete("/renstra/:uuid", commonpresentation.JWTMiddleware(), commonpresentation.RBACMiddleware(admin, whoamiURL), DeleteRenstraHandlerfunc)
	app.Get("/renstra/:uuid", commonpresentation.JWTMiddleware(), GetRenstraHandlerfunc)
	app.Get("/renstras", commonpresentation.JWTMiddleware(), GetAllRenstrasHandlerfunc)

	//untuk audit
	app.Get("/renstra/audit/:tahun", commonpresentation.JWTMiddleware(), commonpresentation.RBACMiddleware(audit, whoamiURL), GetRenstraByTahunTargetHandlerfunc)
}
