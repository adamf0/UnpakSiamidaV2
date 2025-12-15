package presentation

import (
    "context"
    "github.com/gofiber/fiber/v2"
    "github.com/mehdihadeli/go-mediatr"
    "strings"
    
    // "UnpakSiamida/common/domain"
    commoninfra "UnpakSiamida/common/infrastructure"
    commondomain "UnpakSiamida/common/domain"
    standarrenstradomain "UnpakSiamida/modules/standarrenstra/domain"
    CreateStandarRenstra "UnpakSiamida/modules/standarrenstra/application/CreateStandarRenstra"
    UpdateStandarRenstra "UnpakSiamida/modules/standarrenstra/application/UpdateStandarRenstra"
    DeleteStandarRenstra "UnpakSiamida/modules/standarrenstra/application/DeleteStandarRenstra"
    GetStandarRenstra "UnpakSiamida/modules/standarrenstra/application/GetStandarRenstra"
    GetAllStandarRenstras "UnpakSiamida/modules/standarrenstra/application/GetAllStandarRenstras"
)

// =======================================================
// POST /standarrenstra
// =======================================================

// CreateStandarRenstraHandler godoc
// @Summary Create new StandarRenstra
// @Tags StandarRenstra
// @Param nama formData string true "Nama StandarRenstra"
// @Produce json
// @Success 200 {object} map[string]string "uuid of created StandarRenstra"
// @Failure 400 {object} commondomain.Error
// @Router /standarrenstra [post]
func CreateStandarRenstraHandlerfunc(c *fiber.Ctx) error {

    nama := c.FormValue("nama")

    cmd := CreateStandarRenstra.CreateStandarRenstraCommand{
        Nama:         nama,
    }

    uuid, err := mediatr.Send[CreateStandarRenstra.CreateStandarRenstraCommand, string](context.Background(), cmd)
    if err != nil {
        return commoninfra.HandleError(c, err)
    }

    return c.JSON(fiber.Map{"uuid": uuid})
}

// =======================================================
// PUT /standarrenstra/{uuid}
// =======================================================

// UpdateStandarRenstraHandler godoc
// @Summary Update existing StandarRenstra
// @Tags StandarRenstra
// @Param uuid path string true "StandarRenstra UUID" format(uuid)
// @Param nama formData string true "Nama StandarRenstra"
// @Produce json
// @Success 200 {object} map[string]string "uuid of updated StandarRenstra"
// @Failure 400 {object} commondomain.Error
// @Router /standarrenstra/{uuid} [put]
func UpdateStandarRenstraHandlerfunc(c *fiber.Ctx) error {

    uuid := c.Params("uuid")

    nama := c.FormValue("nama")

    cmd := UpdateStandarRenstra.UpdateStandarRenstraCommand{
        Uuid:         uuid,
        Nama:         nama,
    }

    updatedID, err := mediatr.Send[UpdateStandarRenstra.UpdateStandarRenstraCommand, string](context.Background(), cmd)
    if err != nil {
        return commoninfra.HandleError(c, err)
    }

    return c.JSON(fiber.Map{"uuid": updatedID})
}

// =======================================================
// DELETE /standarrenstra/{uuid}
// =======================================================

// DeleteStandarRenstraHandler godoc
// @Summary Delete a StandarRenstra
// @Tags StandarRenstra
// @Param uuid path string true "StandarRenstra UUID" format(uuid)
// @Produce json
// @Success 200 {object} map[string]string "uuid of deleted StandarRenstra"
// @Failure 404 {object} commondomain.Error
// @Router /standarrenstra/{uuid} [delete]
func DeleteStandarRenstraHandlerfunc(c *fiber.Ctx) error {

    uuid := c.Params("uuid")

    cmd := DeleteStandarRenstra.DeleteStandarRenstraCommand{
        Uuid: uuid,
    }

    deletedID, err := mediatr.Send[DeleteStandarRenstra.DeleteStandarRenstraCommand, string](context.Background(), cmd)
    if err != nil {
        return commoninfra.HandleError(c, err)
    }

    return c.JSON(fiber.Map{"uuid": deletedID})
}

// =======================================================
// GET /standarrenstra/{uuid}
// =======================================================

// GetStandarRenstraHandler godoc
// @Summary Get StandarRenstra by UUID
// @Tags StandarRenstra
// @Param uuid path string true "StandarRenstra UUID" format(uuid)
// @Produce json
// @Success 200 {object} standarrenstradomain.StandarRenstra
// @Failure 404 {object} commondomain.Error
// @Router /standarrenstra/{uuid} [get]
func GetStandarRenstraHandlerfunc(c *fiber.Ctx) error {
    uuid := c.Params("uuid")

    query := GetStandarRenstra.GetStandarRenstraByUuidQuery{
        Uuid: uuid,
    }

    standarrenstra, err := mediatr.Send[GetStandarRenstra.GetStandarRenstraByUuidQuery, *standarrenstradomain.StandarRenstra](context.Background(), query)
    if err != nil {
        return commoninfra.HandleError(c, err)
    }

    if standarrenstra == nil {
        return c.Status(404).JSON(fiber.Map{"error": "StandarRenstra not found"})
    }

    return c.JSON(standarrenstra)
}

// =======================================================
// GET /standarrenstras
// =======================================================

// GetAllStandarRenstrasHandler godoc
// @Summary Get All StandarRenstras
// @Tags StandarRenstra
// @Param mode query string false "paging | all | ndjson | sse"
// @Param page query int false "Page number"
// @Param limit query int false "Limit per page"
// @Param search query string false "Search keyword"
// @Produce json
// @Success 200 {object} standarrenstradomain.PagedStandarRenstras
// @Router /standarrenstras [get]
func GetAllStandarRenstrasHandlerfunc(c *fiber.Ctx) error {
    mode := c.Query("mode", "paging") // default mode = paging
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

    query := GetAllStandarRenstras.GetAllStandarRenstrasQuery{
        Search:        search,
        SearchFilters: filters,
    }

    // Pilih adapter sesuai mode
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

    // Ambil data
    standarrenstras, err := mediatr.Send[GetAllStandarRenstras.GetAllStandarRenstrasQuery, standarrenstradomain.PagedStandarRenstras](context.Background(), query)
    if err != nil {
        return commoninfra.HandleError(c, err)
    }

    return adapter.Send(c, standarrenstras)
}

func ModuleStandarRenstra(app *fiber.App) {
    app.Post("/standarrenstra", CreateStandarRenstraHandlerfunc)
    app.Put("/standarrenstra/:uuid", UpdateStandarRenstraHandlerfunc)
    app.Delete("/standarrenstra/:uuid", DeleteStandarRenstraHandlerfunc)
    app.Get("/standarrenstra/:uuid", GetStandarRenstraHandlerfunc)
    app.Get("/standarrenstras", GetAllStandarRenstrasHandlerfunc)
}

