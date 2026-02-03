package presentation

import (
	"context"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/mehdihadeli/go-mediatr"

	commondomain "UnpakSiamida/common/domain"
	commoninfra "UnpakSiamida/common/infrastructure"
	commonpresentation "UnpakSiamida/common/presentation"
	TahunProkerdomain "UnpakSiamida/modules/tahunproker/domain"

	CreateTahunProker "UnpakSiamida/modules/tahunproker/application/CreateTahunProker"
	DeleteTahunProker "UnpakSiamida/modules/tahunproker/application/DeleteTahunProker"
	GetAllTahunProkers "UnpakSiamida/modules/tahunproker/application/GetAllTahunProkers"
	GetTahunProker "UnpakSiamida/modules/tahunproker/application/GetTahunProker"
	SetupUuidTahunProker "UnpakSiamida/modules/tahunproker/application/SetupUuidTahunProker"
	UpdateTahunProker "UnpakSiamida/modules/tahunproker/application/UpdateTahunProker"
)

// =======================================================
// POST /tahunproker
// =======================================================

// CreateTahunProkerHandler godoc
// @Summary Create new TahunProker
// @Tags TahunProker
// @Param nama formData string true "Nama TahunProker"
// @Produce json
// @Success 200 {object} map[string]string "uuid of created TahunProker"
// @Failure 400 {object} commoninfra.ResponseError
// @Failure 404 {object} commoninfra.ResponseError
// @Failure 409 {object} commoninfra.ResponseError
// @Failure 500 {object} commoninfra.ResponseError
// @Router /tahunproker [post]
func CreateTahunProkerHandlerfunc(c *fiber.Ctx) error {

	tahun := c.FormValue("tahun")
	status := c.FormValue("status")

	cmd := CreateTahunProker.CreateTahunProkerCommand{
		Tahun:  tahun,
		Status: status,
	}

	uuid, err := mediatr.Send[CreateTahunProker.CreateTahunProkerCommand, string](context.Background(), cmd)
	if err != nil {
		return commoninfra.HandleError(c, err)
	}

	return c.JSON(fiber.Map{"uuid": uuid})
}

// =======================================================
// PUT /tahunproker/{uuid}
// =======================================================

// UpdateTahunProkerHandler godoc
// @Summary Update existing TahunProker
// @Tags TahunProker
// @Param uuid path string true "TahunProker UUID" format(uuid)
// @Param nama formData string true "Nama TahunProker"
// @Produce json
// @Success 200 {object} map[string]string "uuid of updated TahunProker"
// @Failure 400 {object} commoninfra.ResponseError
// @Failure 404 {object} commoninfra.ResponseError
// @Failure 409 {object} commoninfra.ResponseError
// @Failure 500 {object} commoninfra.ResponseError
// @Router /tahunproker/{uuid} [put]
func UpdateTahunProkerHandlerfunc(c *fiber.Ctx) error {

	uuid := c.Params("uuid")

	tahun := c.FormValue("tahun")
	status := c.FormValue("status")

	cmd := UpdateTahunProker.UpdateTahunProkerCommand{
		Uuid:   uuid,
		Tahun:  tahun,
		Status: status,
	}

	updatedID, err := mediatr.Send[UpdateTahunProker.UpdateTahunProkerCommand, string](context.Background(), cmd)
	if err != nil {
		return commoninfra.HandleError(c, err)
	}

	return c.JSON(fiber.Map{"uuid": updatedID})
}

// =======================================================
// DELETE /tahunproker/{uuid}
// =======================================================

// DeleteTahunProkerHandler godoc
// @Summary Delete a TahunProker
// @Tags TahunProker
// @Param uuid path string true "TahunProker UUID" format(uuid)
// @Produce json
// @Success 200 {object} map[string]string "uuid of deleted TahunProker"
// @Failure 400 {object} commoninfra.ResponseError
// @Failure 404 {object} commoninfra.ResponseError
// @Failure 409 {object} commoninfra.ResponseError
// @Failure 500 {object} commoninfra.ResponseError
//
//	@Router /tahunproker/{uuid} [delete]
func DeleteTahunProkerHandlerfunc(c *fiber.Ctx) error {

	uuid := c.Params("uuid")

	cmd := DeleteTahunProker.DeleteTahunProkerCommand{
		Uuid: uuid,
	}

	deletedID, err := mediatr.Send[DeleteTahunProker.DeleteTahunProkerCommand, string](context.Background(), cmd)
	if err != nil {
		return commoninfra.HandleError(c, err)
	}

	return c.JSON(fiber.Map{"uuid": deletedID})
}

// =======================================================
// GET /TahunProker/{uuid}
// =======================================================

// GetTahunProkerHandler godoc
// @Summary Get TahunProker by UUID
// @Tags TahunProker
// @Param uuid path string true "TahunProker UUID" format(uuid)
// @Produce json
// @Success 200 {object} TahunProkerdomain.TahunProker
// @Failure 400 {object} commoninfra.ResponseError
// @Failure 404 {object} commoninfra.ResponseError
// @Failure 409 {object} commoninfra.ResponseError
// @Failure 500 {object} commoninfra.ResponseError
//
//	@Router /TahunProker/{uuid} [get]
func GetTahunProkerHandlerfunc(c *fiber.Ctx) error {
	uuid := c.Params("uuid")

	query := GetTahunProker.GetTahunProkerByUuidQuery{Uuid: uuid}

	TahunProker, err := mediatr.Send[
		GetTahunProker.GetTahunProkerByUuidQuery,
		*TahunProkerdomain.TahunProker,
	](context.Background(), query)

	if err != nil {
		return commoninfra.HandleError(c, err)
	}

	if TahunProker == nil {
		return c.Status(404).JSON(fiber.Map{"error": "TahunProker not found"})
	}

	return c.JSON(TahunProker)
}

// =======================================================
// GET /TahunProkers
// =======================================================

// GetAllTahunProkersHandler godoc
// @Summary Get all TahunProkers
// @Tags TahunProker
// @Param mode query string false "paging | all | ndjson | sse"
// @Param page query int false "Page number"
// @Param limit query int false "Limit per page"
// @Param search query string false "Search keyword"
// @Produce json
// @Success 200 {object} TahunProkerdomain.PagedTahunProkers
// @Router /TahunProkers [get]
func GetAllTahunProkersHandlerfunc(c *fiber.Ctx) error {
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

	query := GetAllTahunProkers.GetAllTahunProkersQuery{
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
		GetAllTahunProkers.GetAllTahunProkersQuery,
		TahunProkerdomain.PagedTahunProkers,
	](context.Background(), query)

	if err != nil {
		return commoninfra.HandleError(c, err)
	}

	return adapter.Send(c, result)
}

func SetupUuidTahunProkersHandlerfunc(c *fiber.Ctx) error {
	cmd := SetupUuidTahunProker.SetupUuidTahunProkerCommand{}

	message, err := mediatr.Send[SetupUuidTahunProker.SetupUuidTahunProkerCommand, string](context.Background(), cmd)
	if err != nil {
		return commoninfra.HandleError(c, err)
	}

	return c.JSON(fiber.Map{"message": message})
}

func ModuleTahunProker(app *fiber.App) {
	admin := []string{"admin"}
	whoamiURL := "http://localhost:3000/whoami"

	app.Get("/tahunproker/setupuuid", commonpresentation.JWTMiddleware(), commonpresentation.RBACMiddleware(admin, whoamiURL), SetupUuidTahunProkersHandlerfunc)

	app.Post("/tahunproker", commonpresentation.JWTMiddleware(), commonpresentation.RBACMiddleware(admin, whoamiURL), CreateTahunProkerHandlerfunc)
	app.Put("/tahunproker/:uuid", commonpresentation.JWTMiddleware(), commonpresentation.RBACMiddleware(admin, whoamiURL), UpdateTahunProkerHandlerfunc)
	app.Delete("/tahunproker/:uuid", commonpresentation.JWTMiddleware(), commonpresentation.RBACMiddleware(admin, whoamiURL), DeleteTahunProkerHandlerfunc)
	app.Get("/tahunproker/:uuid", commonpresentation.SmartCompress(), commonpresentation.JWTMiddleware(), GetTahunProkerHandlerfunc)
	app.Get("/tahunprokers", commonpresentation.SmartCompress(), commonpresentation.JWTMiddleware(), GetAllTahunProkersHandlerfunc)
}
