package presentation

import (
	"context"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/mehdihadeli/go-mediatr"

	commondomain "UnpakSiamida/common/domain"
	commoninfra "UnpakSiamida/common/infrastructure"
	commonpresentation "UnpakSiamida/common/presentation"
	DokumenProkerdomain "UnpakSiamida/modules/dokumenproker/domain"

	CreateDokumenProker "UnpakSiamida/modules/dokumenproker/application/CreateDokumenProker"
	DeleteDokumenProker "UnpakSiamida/modules/dokumenproker/application/DeleteDokumenProker"
	GetAllDokumenProkers "UnpakSiamida/modules/dokumenproker/application/GetAllDokumenProkers"
	GetDokumenProker "UnpakSiamida/modules/dokumenproker/application/GetDokumenProker"
	SetupUuidDokumenProker "UnpakSiamida/modules/dokumenproker/application/SetupUuidDokumenProker"
	UpdateDokumenProker "UnpakSiamida/modules/dokumenproker/application/UpdateDokumenProker"
)

// =======================================================
// POST /dokumenproker
// =======================================================

// CreateDokumenProkerHandler godoc
// @Summary Create new DokumenProker
// @Tags DokumenProker
// @Param fakultas_unit formData string true "Fakultas/Unit Uuid"
// @Param mata_program formData string true "Mata Program Uuid"
// @Param jenis_dokumen formData string true "Jenis Dokumen"
// @Param file formData string true "file"
// @Param status formData string true "Status"
// @Param catatan formData string false "Catatan"
// @Produce json
// @Success 200 {object} map[string]string "uuid of created DokumenProker"
// @Failure 400 {object} commoninfra.ResponseError
// @Failure 404 {object} commoninfra.ResponseError
// @Failure 409 {object} commoninfra.ResponseError
// @Failure 500 {object} commoninfra.ResponseError
// @Router /dokumenproker [post]
func CreateDokumenProkerHandlerfunc(c *fiber.Ctx) error {

	fakultas_unit := c.FormValue("fakultas_unit")
	mata_program := c.FormValue("mata_program")
	jenis_dokumen := c.FormValue("jenis_dokumen")
	file := c.FormValue("file")
	status := c.FormValue("status")
	catatan := strPtr(c.FormValue("catatan"))

	cmd := CreateDokumenProker.CreateDokumenProkerCommand{
		FakultasUuid:    fakultas_unit,
		MataProgramUuid: mata_program,
		JenisDokumen:    jenis_dokumen,
		File:            file,
		Status:          status,
		Catatan:         catatan,
	}

	uuid, err := mediatr.Send[CreateDokumenProker.CreateDokumenProkerCommand, string](context.Background(), cmd)
	if err != nil {
		return commoninfra.HandleError(c, err)
	}

	return c.JSON(fiber.Map{"uuid": uuid})
}

// =======================================================
// PUT /dokumenproker/{uuid}
// =======================================================

// UpdateDokumenProkerHandler godoc
// @Summary Update existing DokumenProker
// @Tags DokumenProker
// @Param uuid path string true "DokumenProker UUID" format(uuid)
// @Param fakultas_unit formData string true "Fakultas/Unit Uuid"
// @Param mata_program formData string true "Mata Program Uuid"
// @Param jenis_dokumen formData string true "Jenis Dokumen"
// @Param file formData string true "file"
// @Param status formData string true "Status"
// @Param catatan formData string false "Catatan"
// @Produce json
// @Success 200 {object} map[string]string "uuid of updated DokumenProker"
// @Failure 400 {object} commoninfra.ResponseError
// @Failure 404 {object} commoninfra.ResponseError
// @Failure 409 {object} commoninfra.ResponseError
// @Failure 500 {object} commoninfra.ResponseError
// @Router /dokumenproker/{uuid} [put]
func UpdateDokumenProkerHandlerfunc(c *fiber.Ctx) error {

	uuid := c.Params("uuid")
	fakultas_unit := c.FormValue("fakultas_unit")
	mata_program := c.FormValue("mata_program")
	jenis_dokumen := c.FormValue("jenis_dokumen")
	file := c.FormValue("file")
	status := c.FormValue("status")
	catatan := strPtr(c.FormValue("catatan"))

	cmd := UpdateDokumenProker.UpdateDokumenProkerCommand{
		Uuid:            uuid,
		FakultasUuid:    fakultas_unit,
		MataProgramUuid: mata_program,
		JenisDokumen:    jenis_dokumen,
		File:            file,
		Status:          status,
		Catatan:         catatan,
	}

	updatedID, err := mediatr.Send[UpdateDokumenProker.UpdateDokumenProkerCommand, string](context.Background(), cmd)
	if err != nil {
		return commoninfra.HandleError(c, err)
	}

	return c.JSON(fiber.Map{"uuid": updatedID})
}

// =======================================================
// DELETE /dokumenproker/{uuid}
// =======================================================

// DeleteDokumenProkerHandler godoc
// @Summary Delete a DokumenProker
// @Tags DokumenProker
// @Param uuid path string true "DokumenProker UUID" format(uuid)
// @Produce json
// @Success 200 {object} map[string]string "uuid of deleted DokumenProker"
// @Failure 400 {object} commoninfra.ResponseError
// @Failure 404 {object} commoninfra.ResponseError
// @Failure 409 {object} commoninfra.ResponseError
// @Failure 500 {object} commoninfra.ResponseError
// @Router /dokumenproker/{uuid} [delete]
func DeleteDokumenProkerHandlerfunc(c *fiber.Ctx) error {

	uuid := c.Params("uuid")

	cmd := DeleteDokumenProker.DeleteDokumenProkerCommand{
		Uuid: uuid,
	}

	deletedID, err := mediatr.Send[DeleteDokumenProker.DeleteDokumenProkerCommand, string](context.Background(), cmd)
	if err != nil {
		return commoninfra.HandleError(c, err)
	}

	return c.JSON(fiber.Map{"uuid": deletedID})
}

// =======================================================
// GET /DokumenProker/{uuid}
// =======================================================

// GetDokumenProkerHandler godoc
// @Summary Get DokumenProker by UUID
// @Tags DokumenProker
// @Param uuid path string true "DokumenProker UUID" format(uuid)
// @Produce json
// @Success 200 {object} DokumenProkerdomain.DokumenProker
// @Failure 400 {object} commoninfra.ResponseError
// @Failure 404 {object} commoninfra.ResponseError
// @Failure 409 {object} commoninfra.ResponseError
// @Failure 500 {object} commoninfra.ResponseError
// @Router /DokumenProker/{uuid} [get]
func GetDokumenProkerHandlerfunc(c *fiber.Ctx) error {
	uuid := c.Params("uuid")

	query := GetDokumenProker.GetDokumenProkerByUuidQuery{Uuid: uuid}

	DokumenProker, err := mediatr.Send[
		GetDokumenProker.GetDokumenProkerByUuidQuery,
		*DokumenProkerdomain.DokumenProker,
	](context.Background(), query)

	if err != nil {
		return commoninfra.HandleError(c, err)
	}

	if DokumenProker == nil {
		return c.Status(404).JSON(fiber.Map{"error": "DokumenProker not found"})
	}

	return c.JSON(DokumenProker)
}

// =======================================================
// GET /DokumenProkers
// =======================================================

// GetAllDokumenProkersHandler godoc
// @Summary Get all DokumenProkers
// @Tags DokumenProker
// @Param mode query string false "paging | all | ndjson | sse"
// @Param page query int false "Page number"
// @Param limit query int false "Limit per page"
// @Param search query string false "Search keyword"
// @Produce json
// @Success 200 {object} DokumenProkerdomain.PagedDokumenProkers
// @Router /DokumenProkers [get]
func GetAllDokumenProkersHandlerfunc(c *fiber.Ctx) error {
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

	query := GetAllDokumenProkers.GetAllDokumenProkersQuery{
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
		GetAllDokumenProkers.GetAllDokumenProkersQuery,
		DokumenProkerdomain.PagedDokumenProkers,
	](context.Background(), query)

	if err != nil {
		return commoninfra.HandleError(c, err)
	}

	return adapter.Send(c, result)
}

func SetupUuidDokumenProkersHandlerfunc(c *fiber.Ctx) error {
	cmd := SetupUuidDokumenProker.SetupUuidDokumenProkerCommand{}

	message, err := mediatr.Send[SetupUuidDokumenProker.SetupUuidDokumenProkerCommand, string](context.Background(), cmd)
	if err != nil {
		return commoninfra.HandleError(c, err)
	}

	return c.JSON(fiber.Map{"message": message})
}

func ModuleDokumenProker(app *fiber.App) {
	admin := []string{"admin"}
	whoamiURL := "http://localhost:3000/whoami"

	app.Get("/dokumenproker/setupuuid", commonpresentation.JWTMiddleware(), commonpresentation.RBACMiddleware(admin, whoamiURL), SetupUuidDokumenProkersHandlerfunc)

	app.Post("/dokumenproker", commonpresentation.JWTMiddleware(), commonpresentation.RBACMiddleware(admin, whoamiURL), CreateDokumenProkerHandlerfunc)
	app.Put("/dokumenproker/:uuid", commonpresentation.JWTMiddleware(), commonpresentation.RBACMiddleware(admin, whoamiURL), UpdateDokumenProkerHandlerfunc)
	app.Delete("/dokumenproker/:uuid", commonpresentation.JWTMiddleware(), commonpresentation.RBACMiddleware(admin, whoamiURL), DeleteDokumenProkerHandlerfunc)
	app.Get("/DokumenProker/:uuid", commonpresentation.SmartCompress(), commonpresentation.JWTMiddleware(), GetDokumenProkerHandlerfunc)
	app.Get("/DokumenProkers", commonpresentation.SmartCompress(), commonpresentation.JWTMiddleware(), GetAllDokumenProkersHandlerfunc)
}

func strPtr(v string) *string {
	return &v
}
