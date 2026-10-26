package presentation

import (
	"context"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/mehdihadeli/go-mediatr"

	commondomain "UnpakSiamida/common/domain"
	commoninfra "UnpakSiamida/common/infrastructure"
	commonpresentation "UnpakSiamida/common/presentation"
	MataProgramdomain "UnpakSiamida/modules/mataprogram/domain"

	CreateMataProgram "UnpakSiamida/modules/mataprogram/application/CreateMataProgram"
	DeleteMataProgram "UnpakSiamida/modules/mataprogram/application/DeleteMataProgram"
	GetAllMataPrograms "UnpakSiamida/modules/mataprogram/application/GetAllMataPrograms"
	GetMataProgram "UnpakSiamida/modules/mataprogram/application/GetMataProgram"
	SetupUuidMataProgram "UnpakSiamida/modules/mataprogram/application/SetupUuidMataProgram"
	UpdateMataProgram "UnpakSiamida/modules/mataprogram/application/UpdateMataProgram"
)

// =======================================================
// POST /mataprogram
// =======================================================

// CreateMataProgramHandler godoc
// @Summary Create new MataProgram
// @Tags MataProgram
// @Param tahun path string true "Tahun UUID" format(uuid)
// @Param mataprogram formData string true "Mata Program"
// @Produce json
// @Success 200 {object} map[string]string "uuid of created MataProgram"
// @Failure 400 {object} commoninfra.ResponseError
// @Failure 404 {object} commoninfra.ResponseError
// @Failure 409 {object} commoninfra.ResponseError
// @Failure 500 {object} commoninfra.ResponseError
// @Router /mataprogram [post]
func CreateMataProgramHandlerfunc(c *fiber.Ctx) error {

	tahun := c.FormValue("tahun")
	if tahun == "" {
		tahun = c.FormValue("id_master_tahun")
	}
	mataprogram := c.FormValue("mataprogram")
	if mataprogram == "" {
		mataprogram = c.FormValue("mata_program")
	}

	target := c.FormValue("id_fakultas_unit")
	if target == "" {
		target = c.FormValue("fakultas_unit")
	}
	if target == "" {
		target = c.FormValue("target")
	}
	if target == "" {
		target = c.FormValue("fakultasunit")
	}
	if target == "" {
		if fu, ok := c.Locals("fakultas_unit").(string); ok && fu != "" {
			target = fu
		}
	}
	if target == "" {
		target = strings.TrimSpace(c.Get("X-Active-Unit"))
	}

	var fuID uint = 0
	if target != "" {
		if id, err := strconv.Atoi(target); err == nil && id > 0 {
			fuID = uint(id)
		}
	}

	cmd := CreateMataProgram.CreateMataProgramCommand{
		TahunUuid:        tahun,
		MataProgram:      mataprogram,
		FakultasUnitUuid: target,
		IdFakultasUnit:   fuID,
	}

	uuid, err := mediatr.Send[CreateMataProgram.CreateMataProgramCommand, string](context.Background(), cmd)
	if err != nil {
		return commoninfra.HandleError(c, err)
	}

	return commonpresentation.JsonUUID(c, uuid)
}

// =======================================================
// PUT /mataprogram/{uuid}
// =======================================================

// UpdateMataProgramHandler godoc
// @Summary Update existing MataProgram
// @Tags MataProgram
// @Param uuid path string true "MataProgram UUID" format(uuid)
// @Param tahun path string true "Tahun UUID" format(uuid)
// @Param mataprogram formData string true "Mata Program"
// @Produce json
// @Success 200 {object} map[string]string "uuid of updated MataProgram"
// @Failure 400 {object} commoninfra.ResponseError
// @Failure 404 {object} commoninfra.ResponseError
// @Failure 409 {object} commoninfra.ResponseError
// @Failure 500 {object} commoninfra.ResponseError
// @Router /mataprogram/{uuid} [put]
func UpdateMataProgramHandlerfunc(c *fiber.Ctx) error {

	uuid := c.Params("uuid")

	tahun := c.FormValue("tahun")
	if tahun == "" {
		tahun = c.FormValue("id_master_tahun")
	}
	mataprogram := c.FormValue("mataprogram")
	if mataprogram == "" {
		mataprogram = c.FormValue("mata_program")
	}

	target := c.FormValue("id_fakultas_unit")
	if target == "" {
		target = c.FormValue("fakultas_unit")
	}
	if target == "" {
		target = c.FormValue("target")
	}
	if target == "" {
		target = c.FormValue("fakultasunit")
	}
	if target == "" {
		if fu, ok := c.Locals("fakultas_unit").(string); ok && fu != "" {
			target = fu
		}
	}
	if target == "" {
		target = strings.TrimSpace(c.Get("X-Active-Unit"))
	}

	var fuID uint = 0
	if target != "" {
		if id, err := strconv.Atoi(target); err == nil && id > 0 {
			fuID = uint(id)
		}
	}

	cmd := UpdateMataProgram.UpdateMataProgramCommand{
		Uuid:             uuid,
		TahunUuid:        tahun,
		MataProgram:      mataprogram,
		FakultasUnitUuid: target,
		IdFakultasUnit:   fuID,
	}

	updatedID, err := mediatr.Send[UpdateMataProgram.UpdateMataProgramCommand, string](context.Background(), cmd)
	if err != nil {
		return commoninfra.HandleError(c, err)
	}

	return c.JSON(fiber.Map{"uuid": updatedID})
}

// =======================================================
// DELETE /mataprogram/{uuid}
// =======================================================

// DeleteMataProgramHandler godoc
// @Summary Delete a MataProgram
// @Tags MataProgram
// @Param uuid path string true "MataProgram UUID" format(uuid)
// @Produce json
// @Success 200 {object} map[string]string "uuid of deleted MataProgram"
// @Failure 400 {object} commoninfra.ResponseError
// @Failure 404 {object} commoninfra.ResponseError
// @Failure 409 {object} commoninfra.ResponseError
// @Failure 500 {object} commoninfra.ResponseError
//
//	@Router /mataprogram/{uuid} [delete]
func DeleteMataProgramHandlerfunc(c *fiber.Ctx) error {

	uuid := c.Params("uuid")

	cmd := DeleteMataProgram.DeleteMataProgramCommand{
		Uuid: uuid,
	}

	deletedID, err := mediatr.Send[DeleteMataProgram.DeleteMataProgramCommand, string](context.Background(), cmd)
	if err != nil {
		return commoninfra.HandleError(c, err)
	}

	return c.JSON(fiber.Map{"uuid": deletedID})
}

// =======================================================
// GET /MataProgram/{uuid}
// =======================================================

// GetMataProgramHandler godoc
// @Summary Get MataProgram by UUID
// @Tags MataProgram
// @Param uuid path string true "MataProgram UUID" format(uuid)
// @Produce json
// @Success 200 {object} MataProgramdomain.MataProgram
// @Failure 400 {object} commoninfra.ResponseError
// @Failure 404 {object} commoninfra.ResponseError
// @Failure 409 {object} commoninfra.ResponseError
// @Failure 500 {object} commoninfra.ResponseError
//
//	@Router /MataProgram/{uuid} [get]
func GetMataProgramHandlerfunc(c *fiber.Ctx) error {
	uuid := c.Params("uuid")

	query := GetMataProgram.GetMataProgramByUuidQuery{Uuid: uuid}

	MataProgram, err := mediatr.Send[
		GetMataProgram.GetMataProgramByUuidQuery,
		*MataProgramdomain.MataProgram,
	](context.Background(), query)

	if err != nil {
		return commoninfra.HandleError(c, err)
	}

	if MataProgram == nil {
		return c.Status(404).JSON(fiber.Map{"error": "MataProgram not found"})
	}

	return c.JSON(MataProgram)
}

// =======================================================
// GET /MataPrograms
// =======================================================

// GetAllMataProgramsHandler godoc
// @Summary Get all MataPrograms
// @Tags MataProgram
// @Param mode query string false "paging | all | ndjson | sse"
// @Param page query int false "Page number"
// @Param limit query int false "Limit per page"
// @Param search query string false "Search keyword"
// @Produce json
// @Success 200 {object} commondomain.Paged[MataProgramdomain.MataProgramDefault]
// @Router /MataPrograms [get]
func GetAllMataProgramsHandlerfunc(c *fiber.Ctx) error {
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

	target := strings.TrimSpace(c.Params("target"))
	if target == "" {
		target = strings.TrimSpace(c.Query("target"))
	}
	if target == "" {
		target = strings.TrimSpace(c.Query("id_fakultas_unit"))
	}
	if target == "" {
		filterFakultasUnit := strings.TrimSpace(c.Query("filter_fakultas_unit"))
		if filterFakultasUnit == "" {
			filterFakultasUnit = strings.TrimSpace(c.Query("fakultas_unit"))
		}
		if filterFakultasUnit != "" && filterFakultasUnit != "all" && filterFakultasUnit != "semua" {
			target = filterFakultasUnit
		}
	}
	if target == "" {
		target = strings.TrimSpace(c.Get("X-Active-Unit"))
	}

	currentRole := strings.ToLower(strings.TrimSpace(c.Get("X-Active-Role")))
	if currentRole == "" {
		currentRole = strings.ToLower(strings.TrimSpace(c.Get("X-Level")))
	}
	if currentRole == "" {
		if lvl, ok := c.Locals("level").(string); ok {
			currentRole = strings.ToLower(strings.TrimSpace(lvl))
		}
	}
	if currentRole == "" {
		currentRole = strings.ToLower(strings.TrimSpace(c.FormValue("level")))
	}

	fakultasUnit := strings.TrimSpace(c.FormValue("fakultas_unit"))
	if fakultasUnit == "" {
		fakultasUnit = strings.TrimSpace(c.FormValue("target"))
	}
	if fakultasUnit == "" {
		if fu, ok := c.Locals("fakultas_unit").(string); ok {
			fakultasUnit = fu
		}
	}
	if fakultasUnit == "" {
		fakultasUnit = strings.TrimSpace(c.Get("X-Active-Unit"))
	}

	// For non-admin, non-anggaran, non-evaluator, restrict to their own fakultas_unit forwarded from whoami or header
	isUnrestricted := currentRole == "admin" || currentRole == "anggaran" || currentRole == "keuangan" || currentRole == "evaluator"
	if !isUnrestricted && fakultasUnit != "" && fakultasUnit != "0" {
		target = fakultasUnit
	}

	if target != "" && target != "all" && target != "semua" {
		field := "id_fakultas_unit"
		if strings.Contains(target, "-") {
			field = "fakultasunituuid"
		}
		filters = append(filters, commondomain.SearchFilter{
			Field:    field,
			Operator: "eq",
			Value:    &target,
		})
	}

	filterTahun := strings.TrimSpace(c.Query("filter_tahun"))
	if filterTahun == "" {
		filterTahun = strings.TrimSpace(c.Query("tahun"))
	}
	if filterTahun == "" {
		filterTahun = strings.TrimSpace(c.Get("X-Active-Year"))
	}
	if filterTahun != "" && filterTahun != "all" && filterTahun != "semua" {
		field := "tahun"
		if strings.Contains(filterTahun, "-") {
			field = "tahunuuid"
		}
		filters = append(filters, commondomain.SearchFilter{
			Field:    field,
			Operator: "eq",
			Value:    &filterTahun,
		})
	}

	query := GetAllMataPrograms.GetAllMataProgramsQuery{
		Search:        search,
		SearchFilters: filters,
	}

	var adapter commonpresentation.OutputAdapter[MataProgramdomain.MataProgramDefault]
	switch mode {
	case "all":
		adapter = &commonpresentation.AllAdapter[MataProgramdomain.MataProgramDefault]{}
	case "ndjson":
		adapter = &commonpresentation.NDJSONAdapter[MataProgramdomain.MataProgramDefault]{}
	case "sse":
		adapter = &commonpresentation.SSEAdapter[MataProgramdomain.MataProgramDefault]{}
	default:
		query.Page = &page
		query.Limit = &limit
		adapter = &commonpresentation.PagingAdapter[MataProgramdomain.MataProgramDefault]{}
	}

	result, err := mediatr.Send[
		GetAllMataPrograms.GetAllMataProgramsQuery,
		commondomain.Paged[MataProgramdomain.MataProgramDefault],
	](context.Background(), query)

	if err != nil {
		return commoninfra.HandleError(c, err)
	}

	return adapter.Send(c, result)
}

func SetupUuidMataProgramsHandlerfunc(c *fiber.Ctx) error {
	cmd := SetupUuidMataProgram.SetupUuidMataProgramCommand{}

	message, err := mediatr.Send[SetupUuidMataProgram.SetupUuidMataProgramCommand, string](context.Background(), cmd)
	if err != nil {
		return commoninfra.HandleError(c, err)
	}

	return c.JSON(fiber.Map{"message": message})
}

func ModuleMataProgram(app *fiber.App) {
	roles := []string{"admin", "auditee", "fakultas", "unit", "user", "upml", "upmf", "upms"}
	whoamiURL := "http://localhost:3000/whoami"

	app.Get("/mataprogram/setupuuid", SetupUuidMataProgramsHandlerfunc)

	app.Post("/mataprogram", commonpresentation.JWTMiddleware(), commonpresentation.RBACMiddleware(roles, whoamiURL), CreateMataProgramHandlerfunc)
	app.Put("/mataprogram/:uuid", commonpresentation.JWTMiddleware(), commonpresentation.RBACMiddleware(roles, whoamiURL), UpdateMataProgramHandlerfunc)
	app.Delete("/mataprogram/:uuid", commonpresentation.JWTMiddleware(), commonpresentation.RBACMiddleware(roles, whoamiURL), DeleteMataProgramHandlerfunc)
	app.Get("/mataprogram/:uuid", commonpresentation.SmartCompress(), commonpresentation.JWTMiddleware(), GetMataProgramHandlerfunc)
	app.Get("/mataprograms", commonpresentation.SmartCompress(), commonpresentation.JWTMiddleware(), GetAllMataProgramsHandlerfunc)
	app.Get("/mataprograms/target/:target", commonpresentation.SmartCompress(), commonpresentation.JWTMiddleware(), GetAllMataProgramsHandlerfunc)
	app.Get("/MataProgram/:uuid", commonpresentation.SmartCompress(), commonpresentation.JWTMiddleware(), GetMataProgramHandlerfunc)
	app.Get("/MataPrograms", commonpresentation.SmartCompress(), commonpresentation.JWTMiddleware(), GetAllMataProgramsHandlerfunc)
}
