package presentation

import (
	"context"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/mehdihadeli/go-mediatr"

	commondomain "UnpakSiamida/common/domain"
	commoninfra "UnpakSiamida/common/infrastructure"
	commonpresentation "UnpakSiamida/common/presentation"

	GetAllFakultasUnits "UnpakSiamida/modules/fakultasunit/application/GetAllFakultasUnits"
	GetFakultasUnit "UnpakSiamida/modules/fakultasunit/application/GetFakultasUnit"
	SetupUuidFakultasUnit "UnpakSiamida/modules/fakultasunit/application/SetupUuidFakultasUnit"
	fakultasunitdomain "UnpakSiamida/modules/fakultasunit/domain"
)

// =======================================================
// GET /fakultasunit/{uuid}
// =======================================================

// GetFakultasUnitHandler godoc
// @Summary Get FakultasUnit by UUID
// @Tags FakultasUnit
// @Param uuid path string true "FakultasUnit UUID" format(uuid)
// @Produce json
// @Success 200 {object} fakultasunitdomain.FakultasUnit
// @Failure 400 {object} commoninfra.ResponseError
// @Failure 404 {object} commoninfra.ResponseError
// @Failure 409 {object} commoninfra.ResponseError
// @Failure 500 {object} commoninfra.ResponseError
// @Router /fakultasunit/{uuid} [get]
func GetFakultasUnitHandler(c *fiber.Ctx) error {
	uuid := c.Params("uuid")

	query := GetFakultasUnit.GetFakultasUnitByUuidQuery{
		Uuid: uuid,
	}

	fakultasunit, err := mediatr.Send[
		GetFakultasUnit.GetFakultasUnitByUuidQuery,
		*fakultasunitdomain.FakultasUnit,
	](context.Background(), query)

	if err != nil {
		return commoninfra.HandleError(c, err)
	}

	if fakultasunit == nil {
		return c.Status(404).JSON(fiber.Map{
			"error": "FakultasUnit not found",
		})
	}

	return c.JSON(fakultasunit)
}

// =======================================================
// GET /fakultasunits
// =======================================================

// GetAllFakultasUnitsHandler godoc
// @Summary Get All FakultasUnit
// @Tags FakultasUnit
// @Param mode query string false "paging | all | ndjson | sse"
// @Param page query int false "Page"
// @Param limit query int false "Limit"
// @Param search query string false "Search keyword"
// @Produce json
// @Success 200 {object} fakultasunitdomain.PagedFakultasUnits
// @Router /fakultasunits [get]
func GetAllFakultasUnitsHandler(c *fiber.Ctx) error {
	mode := c.Query("mode", "paging")
	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 10)
	search := c.Query("search", "")

	// Parse filters
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

			var valuePtr *string
			if rawValue != "" && rawValue != "null" {
				valuePtr = &rawValue
			}

			filters = append(filters, commondomain.SearchFilter{
				Field:    field,
				Operator: op,
				Value:    valuePtr,
			})
		}
	}

	query := GetAllFakultasUnits.GetAllFakultasUnitsQuery{
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
		GetAllFakultasUnits.GetAllFakultasUnitsQuery,
		fakultasunitdomain.PagedFakultasUnits,
	](context.Background(), query)

	if err != nil {
		return commoninfra.HandleError(c, err)
	}

	return adapter.Send(c, result)
}

func SetupUuidFakultasUnitsHandlerfunc(c *fiber.Ctx) error {
	cmd := SetupUuidFakultasUnit.SetupUuidFakultasUnitCommand{}

	message, err := mediatr.Send[SetupUuidFakultasUnit.SetupUuidFakultasUnitCommand, string](context.Background(), cmd)
	if err != nil {
		return commoninfra.HandleError(c, err)
	}

	return c.JSON(fiber.Map{"message": message})
}

func ModuleFakultasUnit(app *fiber.App) {
	admin := []string{"admin"}
	whoamiURL := "http://localhost:3000/whoami"

	app.Get("/fakultasunit/setupuuid", commonpresentation.JWTMiddleware(), commonpresentation.RBACMiddleware(admin, whoamiURL), SetupUuidFakultasUnitsHandlerfunc)

	app.Get("/fakultasunit/:uuid", commonpresentation.SmartCompress(), commonpresentation.JWTMiddleware(), GetFakultasUnitHandler)
	app.Get("/fakultasunits", commonpresentation.SmartCompress(), commonpresentation.JWTMiddleware(), GetAllFakultasUnitsHandler)
}
