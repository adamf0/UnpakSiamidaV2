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
	BeritaAcaradomain "UnpakSiamida/modules/beritaacara/domain"

	CreateBeritaAcara "UnpakSiamida/modules/beritaacara/application/CreateBeritaAcara"
	DeleteBeritaAcara "UnpakSiamida/modules/beritaacara/application/DeleteBeritaAcara"
	GetAllBeritaAcaras "UnpakSiamida/modules/beritaacara/application/GetAllBeritaAcaras"
	GetBeritaAcara "UnpakSiamida/modules/beritaacara/application/GetBeritaAcara"
	SetupUuidBeritaAcara "UnpakSiamida/modules/beritaacara/application/SetupUuidBeritaAcara"
	UpdateBeritaAcara "UnpakSiamida/modules/beritaacara/application/UpdateBeritaAcara"
)

// =======================================================
// POST /beritaacara
// =======================================================

// CreateBeritaAcaraHandler godoc
// @Summary Create new BeritaAcara
// @Tags BeritaAcara
// @Param tahun formData string true "Nama Tahun"
// @Param fakultasunit formData string true "Fakultas Unit"
// @Param tanggal formData string true "Tanggal"
// @Param auditee formData string true "Auditee"
// @Param auditor1 formData string true "Auditor1"
// @Param auditor2 formData string true "Auditor2"
// @Produce json
// @Success 200 {object} map[string]string "uuid of created BeritaAcara"
// @Failure 400 {object} commoninfra.ResponseError
// @Failure 404 {object} commoninfra.ResponseError
// @Failure 409 {object} commoninfra.ResponseError
// @Failure 500 {object} commoninfra.ResponseError
// @Router /beritaacara [post]
func CreateBeritaAcaraHandlerfunc(c *fiber.Ctx) error {

	cmd := CreateBeritaAcara.CreateBeritaAcaraCommand{
		Tahun:        c.FormValue("tahun"),
		FakultasUnit: parseInt(c.FormValue("fakultasunit")),
		Tanggal:      c.FormValue("tanggal"),
		Auditee:      parseOptionalInt(c.FormValue("auditee")),
		Auditor1:     parseOptionalInt(c.FormValue("auditor1")),
		Auditor2:     parseOptionalInt(c.FormValue("auditor2")),
	}

	uuid, err := mediatr.Send[CreateBeritaAcara.CreateBeritaAcaraCommand, string](context.Background(), cmd)
	if err != nil {
		return commoninfra.HandleError(c, err)
	}

	return c.JSON(fiber.Map{"uuid": uuid})
}

// =======================================================
// PUT /beritaacara/{uuid}
// =======================================================

// UpdateBeritaAcaraHandler godoc
// @Summary Update existing BeritaAcara
// @Tags BeritaAcara
// @Param uuid path string true "BeritaAcara UUID" format(uuid)
// @Param tahun formData string true "Nama Tahun"
// @Param fakultasunit formData string true "Fakultas Unit"
// @Param tanggal formData string true "Tanggal"
// @Param auditee formData string true "Auditee"
// @Param auditor1 formData string true "Auditor1"
// @Param auditor2 formData string true "Auditor2"
// @Produce json
// @Success 200 {object} map[string]string "uuid of updated BeritaAcara"
// @Failure 400 {object} commoninfra.ResponseError
// @Failure 404 {object} commoninfra.ResponseError
// @Failure 409 {object} commoninfra.ResponseError
// @Failure 500 {object} commoninfra.ResponseError
// @Router /beritaacara/{uuid} [put]
func UpdateBeritaAcaraHandlerfunc(c *fiber.Ctx) error {

	cmd := UpdateBeritaAcara.UpdateBeritaAcaraCommand{
		Uuid:         c.Params("uuid"),
		Tahun:        c.FormValue("tahun"),
		FakultasUnit: parseInt(c.FormValue("fakultasunit")),
		Tanggal:      c.FormValue("tanggal"),
		Auditee:      parseOptionalInt(c.FormValue("auditee")),
		Auditor1:     parseOptionalInt(c.FormValue("auditor1")),
		Auditor2:     parseOptionalInt(c.FormValue("auditor2")),
	}

	updatedID, err := mediatr.Send[UpdateBeritaAcara.UpdateBeritaAcaraCommand, string](context.Background(), cmd)
	if err != nil {
		return commoninfra.HandleError(c, err)
	}

	return c.JSON(fiber.Map{"uuid": updatedID})
}

// =======================================================
// DELETE /beritaacara/{uuid}
// =======================================================

// DeleteBeritaAcaraHandler godoc
// @Summary Delete a BeritaAcara
// @Tags BeritaAcara
// @Param uuid path string true "BeritaAcara UUID" format(uuid)
// @Produce json
// @Success 200 {object} map[string]string "uuid of deleted BeritaAcara"
// @Failure 404 {object} commondomain.Error
// @Router /beritaacara/{uuid} [delete]
func DeleteBeritaAcaraHandlerfunc(c *fiber.Ctx) error {

	uuid := c.Params("uuid")

	cmd := DeleteBeritaAcara.DeleteBeritaAcaraCommand{
		Uuid: uuid,
	}

	deletedID, err := mediatr.Send[DeleteBeritaAcara.DeleteBeritaAcaraCommand, string](context.Background(), cmd)
	if err != nil {
		return commoninfra.HandleError(c, err)
	}

	return c.JSON(fiber.Map{"uuid": deletedID})
}

// =======================================================
// GET /BeritaAcara/{uuid}
// =======================================================

// GetBeritaAcaraHandler godoc
// @Summary Get BeritaAcara by UUID
// @Tags BeritaAcara
// @Param uuid path string true "BeritaAcara UUID" format(uuid)
// @Produce json
// @Success 200 {object} BeritaAcaradomain.BeritaAcara
// @Failure 404 {object} commondomain.Error
// @Router /BeritaAcara/{uuid} [get]
func GetBeritaAcaraHandlerfunc(c *fiber.Ctx) error {

	query := GetBeritaAcara.GetBeritaAcaraByUuidQuery{Uuid: c.Params("uuid")}

	BeritaAcara, err := mediatr.Send[
		GetBeritaAcara.GetBeritaAcaraByUuidQuery,
		*BeritaAcaradomain.BeritaAcara,
	](context.Background(), query)

	if err != nil {
		return commoninfra.HandleError(c, err)
	}

	if BeritaAcara == nil {
		return c.Status(404).JSON(fiber.Map{"error": "BeritaAcara not found"})
	}

	return c.JSON(BeritaAcara)
}

// =======================================================
// GET /BeritaAcaras
// =======================================================

// GetAllBeritaAcarasHandler godoc
// @Summary Get all BeritaAcaras
// @Tags BeritaAcara
// @Param mode query string false "paging | all | ndjson | sse"
// @Param page query int false "Page number"
// @Param limit query int false "Limit per page"
// @Param search query string false "Search keyword"
// @Produce json
// @Success 200 {object} BeritaAcaradomain.PagedBeritaAcaras
// @Router /BeritaAcaras [get]
func GetAllBeritaAcarasHandlerfunc(c *fiber.Ctx) error {
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

	query := GetAllBeritaAcaras.GetAllBeritaAcarasQuery{
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
		GetAllBeritaAcaras.GetAllBeritaAcarasQuery,
		BeritaAcaradomain.PagedBeritaAcaras,
	](context.Background(), query)

	if err != nil {
		return commoninfra.HandleError(c, err)
	}

	return adapter.Send(c, result)
}

func SetupUuidBeritaAcarasHandlerfunc(c *fiber.Ctx) error {
	cmd := SetupUuidBeritaAcara.SetupUuidBeritaAcaraCommand{}

	message, err := mediatr.Send[SetupUuidBeritaAcara.SetupUuidBeritaAcaraCommand, string](context.Background(), cmd)
	if err != nil {
		return commoninfra.HandleError(c, err)
	}

	return c.JSON(fiber.Map{"message": message})
}

func ModuleBeritaAcara(app *fiber.App) {
	admin := []string{"admin"}
	whoamiURL := "http://localhost:3000/whoami"

	app.Get("/beritaacara/setupuuid", commonpresentation.JWTMiddleware(), commonpresentation.RBACMiddleware(admin, whoamiURL), SetupUuidBeritaAcarasHandlerfunc)

	app.Post("/beritaacara", commonpresentation.JWTMiddleware(), commonpresentation.RBACMiddleware(admin, whoamiURL), CreateBeritaAcaraHandlerfunc)
	app.Put("/beritaacara/:uuid", commonpresentation.JWTMiddleware(), commonpresentation.RBACMiddleware(admin, whoamiURL), UpdateBeritaAcaraHandlerfunc)
	app.Delete("/beritaacara/:uuid", commonpresentation.JWTMiddleware(), commonpresentation.RBACMiddleware(admin, whoamiURL), DeleteBeritaAcaraHandlerfunc)
	app.Get("/BeritaAcara/:uuid", commonpresentation.SmartCompress(), commonpresentation.JWTMiddleware(), GetBeritaAcaraHandlerfunc)
	app.Get("/BeritaAcaras", commonpresentation.SmartCompress(), commonpresentation.JWTMiddleware(), GetAllBeritaAcarasHandlerfunc)
}

func parseInt(val string) int {
	if val == "" {
		return 0
	}
	i, err := strconv.Atoi(val)
	if err != nil {
		return 0
	}
	return i
}

func parseOptionalInt(val string) *int {
	if val == "" {
		return nil
	}
	i, err := strconv.Atoi(val)
	if err != nil {
		return nil
	}
	return &i
}
