package presentation

import (
	"context"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/mehdihadeli/go-mediatr"

	commondomain "UnpakSiamida/common/domain"
	commoninfra "UnpakSiamida/common/infrastructure"
	commonpresentation "UnpakSiamida/common/presentation"

	DeleteRenstraNilai "UnpakSiamida/modules/renstranilai/application/DeleteRenstraNilai"
	GetAllRenstraNilais "UnpakSiamida/modules/renstranilai/application/GetAllRenstraNilais"
	GetRenstraNilai "UnpakSiamida/modules/renstranilai/application/GetRenstraNilai"
	SetupUuidRenstraNilai "UnpakSiamida/modules/renstranilai/application/SetupUuidRenstraNilai"
	UpdateRenstraNilai "UnpakSiamida/modules/renstranilai/application/UpdateRenstraNilai"
	renstranilaidomain "UnpakSiamida/modules/renstranilai/domain"
)

// =======================================================
// PUT /renstranilai/{uuid}
// =======================================================

// UpdateRenstraNilaiHandler godoc
// @Summary Update existing Renstra Nilai
// @Description Update nilai Renstra berdasarkan role (auditee / auditor)
// @Tags RenstraNilai
// @Accept multipart/form-data
// @Produce json
//
// @Param uuid path string true "RenstraNilai UUID" format(uuid)
//
// @Param uuidRenstra formData string true "Renstra UUID" format(uuid)
// @Param mode formData string true "Mode akses" Enums(auditee,auditor1,auditor2)
//
// @Param capaian formData string false "Capaian (khusus auditee)"
// @Param catatan formData string false "Catatan (khusus auditee)"
// @Param linkBukti formData string false "Link bukti (khusus auditee)"
//
// @Param capaianAuditor formData string false "Capaian auditor (khusus auditor1 / auditor2)"
// @Param catatanAuditor formData string false "Catatan auditor (khusus auditor1 / auditor2)"
//
// @Success 200 {object} map[string]string "uuid of updated RenstraNilai"
// @Failure 400 {object} commoninfra.ResponseError
// @Failure 404 {object} commoninfra.ResponseError
// @Failure 409 {object} commoninfra.ResponseError
// @Failure 500 {object} commoninfra.ResponseError
// @Router /renstranilai/{uuid} [put]
func UpdateRenstraNilaiHandlerfunc(c *fiber.Ctx) error {
	uuid := c.Params("uuid")
	uuidRenstra := c.FormValue("uuidRenstra")
	tahun := c.Params("tahun")
	mode := c.FormValue("mode")
	granted := c.FormValue("grantedaccess") //dari middleware
	capaian := nullableString(c.FormValue("capaian"))
	catatan := nullableString(c.FormValue("catatan"))
	linkBukti := nullableString(c.FormValue("linkBukti"))
	capaianAuditor := nullableString(c.FormValue("capaianAuditor"))
	catatanAuditor := nullableString(c.FormValue("catatanAuditor"))

	cmd := UpdateRenstraNilai.UpdateRenstraNilaiCommand{
		Uuid:           uuid,
		UuidRenstra:    uuidRenstra,
		Tahun:          tahun,
		Mode:           mode,
		Granted:        granted,
		Capaian:        capaian,
		Catatan:        catatan,
		LinkBukti:      linkBukti,
		CapaianAuditor: capaianAuditor,
		CatatanAuditor: catatanAuditor,
	}

	updatedID, err := mediatr.Send[
		UpdateRenstraNilai.UpdateRenstraNilaiCommand,
		string,
	](context.Background(), cmd)

	if err != nil {
		return commoninfra.HandleError(c, err)
	}

	return c.JSON(fiber.Map{"uuid": updatedID})
}

// =======================================================
// DELETE /renstranilai/{uuid}
// =======================================================

// DeleteRenstraNilaiHandler godoc
// @Summary Delete an RenstraNilai
// @Tags RenstraNilai
// @Param uuid path string true "RenstraNilai UUID" format(uuid)
// @Produce json
// @Success 200 {object} map[string]string "uuid of deleted RenstraNilai"
// @Failure 400 {object} commoninfra.ResponseError
// @Failure 404 {object} commoninfra.ResponseError
// @Failure 409 {object} commoninfra.ResponseError
// @Failure 500 {object} commoninfra.ResponseError
//
//	@Router /renstranilai/{uuid} [delete]
func DeleteRenstraNilaiHandlerfunc(c *fiber.Ctx) error {
	uuid := c.Params("uuid")

	cmd := DeleteRenstraNilai.DeleteRenstraNilaiCommand{Uuid: uuid}

	deletedID, err := mediatr.Send[
		DeleteRenstraNilai.DeleteRenstraNilaiCommand,
		string,
	](context.Background(), cmd)

	if err != nil {
		return commoninfra.HandleError(c, err)
	}

	return c.JSON(fiber.Map{"uuid": deletedID})
}

// =======================================================
// GET /renstranilai/{uuid}
// =======================================================

// GetRenstraNilaiHandler godoc
// @Summary Get RenstraNilai by UUID
// @Tags RenstraNilai
// @Param uuid path string true "RenstraNilai UUID" format(uuid)
// @Produce json
// @Success 200 {object} renstranilaidomain.RenstraNilai
// @Failure 400 {object} commoninfra.ResponseError
// @Failure 404 {object} commoninfra.ResponseError
// @Failure 409 {object} commoninfra.ResponseError
// @Failure 500 {object} commoninfra.ResponseError
//
//	@Router /renstranilai/{uuid} [get]
func GetRenstraNilaiHandlerfunc(c *fiber.Ctx) error {
	uuid := c.Params("uuid")

	query := GetRenstraNilai.GetRenstraNilaiByUuidQuery{Uuid: uuid}

	renstranilai, err := mediatr.Send[
		GetRenstraNilai.GetRenstraNilaiByUuidQuery,
		*renstranilaidomain.RenstraNilai,
	](context.Background(), query)

	if err != nil {
		return commoninfra.HandleError(c, err)
	}

	if renstranilai == nil {
		return c.Status(404).JSON(fiber.Map{"error": "RenstraNilai not found"})
	}

	return c.JSON(renstranilai)
}

// GetAllRenstraNilaisHandler godoc
// @Summary Get All RenstraNilais
// @Tags RenstraNilai
// @Param mode query string false "paging | all | ndjson | sse"
// @Param page query int false "Page number"
// @Param limit query int false "Limit per page"
// @Param search query string false "Search keyword"
// @Produce json
// @Success 200 {object} renstranilaidomain.PagedRenstraNilais
// @Router /renstranilais [get]
func GetAllRenstraNilaisHandlerfunc(c *fiber.Ctx) error {
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

	query := GetAllRenstraNilais.GetAllRenstraNilaisQuery{
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
		GetAllRenstraNilais.GetAllRenstraNilaisQuery,
		renstranilaidomain.PagedRenstraNilais,
	](context.Background(), query)

	if err != nil {
		return commoninfra.HandleError(c, err)
	}

	return adapter.Send(c, result)
}

func SetupUuidRenstraNilaisHandlerfunc(c *fiber.Ctx) error {
	cmd := SetupUuidRenstraNilai.SetupUuidRenstraNilaiCommand{}

	message, err := mediatr.Send[SetupUuidRenstraNilai.SetupUuidRenstraNilaiCommand, string](context.Background(), cmd)
	if err != nil {
		return commoninfra.HandleError(c, err)
	}

	return c.JSON(fiber.Map{"message": message})
}

func ModuleRenstraNilai(app *fiber.App) {
	admin := []string{"admin"}
	audit := []string{"auditee", "auditor1", "auditor2"}
	whoamiURL := "http://localhost:3000/whoami"

	app.Get("/renstranilai/setupuuid", commonpresentation.JWTMiddleware(), commonpresentation.RBACMiddleware(admin, whoamiURL), SetupUuidRenstraNilaisHandlerfunc)

	//hanya admin
	app.Delete("/renstranilai/:uuid", commonpresentation.JWTMiddleware(), commonpresentation.RBACMiddleware(admin, whoamiURL), DeleteRenstraNilaiHandlerfunc)

	//admin & audit
	app.Put("/renstranilai/:tahun/:uuid", commonpresentation.JWTMiddleware(), commonpresentation.RBACMiddleware(audit, whoamiURL), UpdateRenstraNilaiHandlerfunc)

	//private
	app.Get("/renstranilai/:uuid", commonpresentation.SmartCompress(), commonpresentation.JWTMiddleware(), GetRenstraNilaiHandlerfunc)
	app.Get("/renstranilais", commonpresentation.SmartCompress(), commonpresentation.JWTMiddleware(), GetAllRenstraNilaisHandlerfunc)
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
