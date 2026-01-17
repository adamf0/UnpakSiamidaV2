package presentation

import (
	"context"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/mehdihadeli/go-mediatr"

	commondomain "UnpakSiamida/common/domain"
	commoninfra "UnpakSiamida/common/infrastructure"
	commonpresentation "UnpakSiamida/common/presentation"

	DeleteDokumenTambahan "UnpakSiamida/modules/dokumentambahan/application/DeleteDokumenTambahan"
	GetAllDokumenTambahans "UnpakSiamida/modules/dokumentambahan/application/GetAllDokumenTambahans"
	GetDokumenTambahan "UnpakSiamida/modules/dokumentambahan/application/GetDokumenTambahan"
	SetupUuidDokumenTambahan "UnpakSiamida/modules/dokumentambahan/application/SetupUuidDokumenTambahan"
	UpdateDokumenTambahan "UnpakSiamida/modules/dokumentambahan/application/UpdateDokumenTambahan"
	dokumentambahandomain "UnpakSiamida/modules/dokumentambahan/domain"
)

// =======================================================
// PUT /dokumentambahan/{uuid}
// =======================================================

// UpdateDokumenTambahanHandler godoc
// @Summary Update existing Renstra Nilai
// @Description Update nilai Renstra berdasarkan role (auditee / auditor)
// @Tags DokumenTambahan
// @Accept multipart/form-data
// @Produce json
//
// @Param uuid path string true "DokumenTambahan UUID" format(uuid)
//
// @Param uuidRenstra formData string true "Renstra UUID" format(uuid)
// @Param mode formData string true "Mode akses" Enums(auditee,auditor2)
//
// @Param capaian formData string false "Capaian (khusus auditee)"
// @Param catatan formData string false "Catatan (khusus auditee)"
// @Param linkBukti formData string false "Link bukti (khusus auditee)"
//
// @Param capaianAuditor formData string false "Capaian auditor (khusus auditor2)"
// @Param catatanAuditor formData string false "Catatan auditor (khusus auditor2)"
//
// @Success 200 {object} map[string]string "uuid of updated DokumenTambahan"
// @Failure 400 {object} commoninfra.ResponseError
// @Failure 404 {object} commoninfra.ResponseError
// @Failure 409 {object} commoninfra.ResponseError
// @Failure 500 {object} commoninfra.ResponseError
// @Router /dokumentambahan/{uuid} [put]
func UpdateDokumenTambahanHandlerfunc(c *fiber.Ctx) error {
	uuid := c.Params("uuid")
	uuidRenstra := c.FormValue("uuidRenstra")
	tahun := c.Params("tahun")
	mode := c.FormValue("mode")
	granted := c.FormValue("grantedaccess") //dari middleware
	link := nullableString(c.FormValue("link"))
	capaianAuditor := nullableString(c.FormValue("capaianAuditor"))
	catatanAuditor := nullableString(c.FormValue("catatanAuditor"))

	cmd := UpdateDokumenTambahan.UpdateDokumenTambahanCommand{
		Uuid:           uuid,
		UuidRenstra:    uuidRenstra,
		Tahun:          tahun,
		Mode:           mode,
		Granted:        granted,
		Link:           link,
		CapaianAuditor: capaianAuditor,
		CatatanAuditor: catatanAuditor,
	}

	updatedID, err := mediatr.Send[
		UpdateDokumenTambahan.UpdateDokumenTambahanCommand,
		string,
	](context.Background(), cmd)

	if err != nil {
		return commoninfra.HandleError(c, err)
	}

	return c.JSON(fiber.Map{"uuid": updatedID})
}

// =======================================================
// DELETE /dokumentambahan/{uuid}
// =======================================================

// DeleteDokumenTambahanHandler godoc
// @Summary Delete an DokumenTambahan
// @Tags DokumenTambahan
// @Param uuid path string true "DokumenTambahan UUID" format(uuid)
// @Produce json
// @Success 200 {object} map[string]string "uuid of deleted DokumenTambahan"
// @Failure 404 {object} commondomain.Error
// @Router /dokumentambahan/{uuid} [delete]
func DeleteDokumenTambahanHandlerfunc(c *fiber.Ctx) error {
	uuid := c.Params("uuid")

	cmd := DeleteDokumenTambahan.DeleteDokumenTambahanCommand{Uuid: uuid}

	deletedID, err := mediatr.Send[
		DeleteDokumenTambahan.DeleteDokumenTambahanCommand,
		string,
	](context.Background(), cmd)

	if err != nil {
		return commoninfra.HandleError(c, err)
	}

	return c.JSON(fiber.Map{"uuid": deletedID})
}

// =======================================================
// GET /dokumentambahan/{uuid}
// =======================================================

// GetDokumenTambahanHandler godoc
// @Summary Get DokumenTambahan by UUID
// @Tags DokumenTambahan
// @Param uuid path string true "DokumenTambahan UUID" format(uuid)
// @Produce json
// @Success 200 {object} dokumentambahandomain.DokumenTambahan
// @Failure 404 {object} commondomain.Error
// @Router /dokumentambahan/{uuid} [get]
func GetDokumenTambahanHandlerfunc(c *fiber.Ctx) error {
	uuid := c.Params("uuid")

	query := GetDokumenTambahan.GetDokumenTambahanByUuidQuery{Uuid: uuid}

	dokumentambahan, err := mediatr.Send[
		GetDokumenTambahan.GetDokumenTambahanByUuidQuery,
		*dokumentambahandomain.DokumenTambahan,
	](context.Background(), query)

	if err != nil {
		return commoninfra.HandleError(c, err)
	}

	if dokumentambahan == nil {
		return c.Status(404).JSON(fiber.Map{"error": "DokumenTambahan not found"})
	}

	return c.JSON(dokumentambahan)
}

// GetAllDokumenTambahansHandler godoc
// @Summary Get All DokumenTambahans
// @Tags DokumenTambahan
// @Param mode query string false "paging | all | ndjson | sse"
// @Param page query int false "Page number"
// @Param limit query int false "Limit per page"
// @Param search query string false "Search keyword"
// @Produce json
// @Success 200 {object} dokumentambahandomain.PagedDokumenTambahans
// @Router /dokumentambahans [get]
func GetAllDokumenTambahansHandlerfunc(c *fiber.Ctx) error {
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

	query := GetAllDokumenTambahans.GetAllDokumenTambahansQuery{
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
		GetAllDokumenTambahans.GetAllDokumenTambahansQuery,
		dokumentambahandomain.PagedDokumenTambahans,
	](context.Background(), query)

	if err != nil {
		return commoninfra.HandleError(c, err)
	}

	return adapter.Send(c, result)
}

func SetupUuidDokumenTambahansHandlerfunc(c *fiber.Ctx) error {
	cmd := SetupUuidDokumenTambahan.SetupUuidDokumenTambahanCommand{}

	message, err := mediatr.Send[SetupUuidDokumenTambahan.SetupUuidDokumenTambahanCommand, string](context.Background(), cmd)
	if err != nil {
		return commoninfra.HandleError(c, err)
	}

	return c.JSON(fiber.Map{"message": message})
}

func ModuleDokumenTambahan(app *fiber.App) {
	admin := []string{"admin"}
	audit := []string{"auditee", "auditor2"}
	whoamiURL := "http://localhost:3000/whoami"

	app.Get("/dokumentambahan/setupuuid", commonpresentation.JWTMiddleware(), commonpresentation.RBACMiddleware(admin, whoamiURL), SetupUuidDokumenTambahansHandlerfunc)

	//hanya admin
	app.Delete("/dokumentambahan/:uuid", commonpresentation.JWTMiddleware(), commonpresentation.RBACMiddleware(admin, whoamiURL), DeleteDokumenTambahanHandlerfunc)

	//admin & audit
	app.Put("/dokumentambahan/:tahun/:uuid", commonpresentation.JWTMiddleware(), commonpresentation.RBACMiddleware(audit, whoamiURL), UpdateDokumenTambahanHandlerfunc)

	//private
	app.Get("/dokumentambahan/:uuid", commonpresentation.JWTMiddleware(), GetDokumenTambahanHandlerfunc)
	app.Get("/dokumentambahans", commonpresentation.JWTMiddleware(), GetAllDokumenTambahansHandlerfunc)
}

// ====================================================================
// Helper: convert empty string → nil
// ====================================================================
func nullableString(s string) *string {
	trimmed := strings.TrimSpace(s)
	if trimmed == "" {
		return nil
	}
	return &trimmed
}
