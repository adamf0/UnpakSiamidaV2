package presentation

import (
	"context"
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
	mataprogram := c.FormValue("mataprogram")

	cmd := CreateMataProgram.CreateMataProgramCommand{
		TahunUuid:   tahun,
		MataProgram: mataprogram,
	}

	uuid, err := mediatr.Send[CreateMataProgram.CreateMataProgramCommand, string](context.Background(), cmd)
	if err != nil {
		return commoninfra.HandleError(c, err)
	}

	return c.JSON(fiber.Map{"uuid": uuid})
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
	mataprogram := c.FormValue("mataprogram")

	cmd := UpdateMataProgram.UpdateMataProgramCommand{
		Uuid:        uuid,
		TahunUuid:   tahun,
		MataProgram: mataprogram,
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
// @Success 200 {object} MataProgramdomain.PagedMataPrograms
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

	query := GetAllMataPrograms.GetAllMataProgramsQuery{
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
		GetAllMataPrograms.GetAllMataProgramsQuery,
		MataProgramdomain.PagedMataPrograms,
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
	admin := []string{"admin"}
	whoamiURL := "http://localhost:3000/whoami"

	app.Get("/mataprogram/setupuuid", commonpresentation.JWTMiddleware(), commonpresentation.RBACMiddleware(admin, whoamiURL), SetupUuidMataProgramsHandlerfunc)

	app.Post("/mataprogram", commonpresentation.JWTMiddleware(), commonpresentation.RBACMiddleware(admin, whoamiURL), CreateMataProgramHandlerfunc)
	app.Put("/mataprogram/:uuid", commonpresentation.JWTMiddleware(), commonpresentation.RBACMiddleware(admin, whoamiURL), UpdateMataProgramHandlerfunc)
	app.Delete("/mataprogram/:uuid", commonpresentation.JWTMiddleware(), commonpresentation.RBACMiddleware(admin, whoamiURL), DeleteMataProgramHandlerfunc)
	app.Get("/MataProgram/:uuid", commonpresentation.SmartCompress(), commonpresentation.JWTMiddleware(), GetMataProgramHandlerfunc)
	app.Get("/MataPrograms", commonpresentation.SmartCompress(), commonpresentation.JWTMiddleware(), GetAllMataProgramsHandlerfunc)
}
