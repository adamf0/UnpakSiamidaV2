package presentation

import (
    "context"
    "github.com/gofiber/fiber/v2"
    "github.com/mehdihadeli/go-mediatr"
    "strings"
    
    // "UnpakSiamida/common/domain"
    commoninfra "UnpakSiamida/common/infrastructure"
    commondomain "UnpakSiamida/common/domain"
    fakultasunitdomain "UnpakSiamida/modules/fakultasunit/domain"
    GetFakultasUnit "UnpakSiamida/modules/fakultasunit/application/GetFakultasUnit"
    GetAllFakultasUnits "UnpakSiamida/modules/fakultasunit/application/GetAllFakultasUnits"
)

func ModuleFakultasUnit(app *fiber.App) {

    // ------------------------------------------------------ 
    // Get fakultasunit
    // ------------------------------------------------------
    app.Get("/fakultasunit/:uuid", func(c *fiber.Ctx) error {
        uuid := c.Params("uuid")

        query := GetFakultasUnit.GetFakultasUnitByUuidQuery{
            Uuid: uuid,
        }

        fakultasunit, err := mediatr.Send[GetFakultasUnit.GetFakultasUnitByUuidQuery, *fakultasunitdomain.FakultasUnit](context.Background(), query)
        if err != nil {
            return commoninfra.HandleError(c, err)
        }

        if fakultasunit == nil {
            return c.Status(404).JSON(fiber.Map{"error": "FakultasUnit not found"})
        }

        return c.JSON(fakultasunit)
    })

    // ------------------------------------------------------ 
    // Get All fakultasunit
    // ------------------------------------------------------
    app.Get("/fakultasunits", func(c *fiber.Ctx) error {
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

        query := GetAllFakultasUnits.GetAllFakultasUnitsQuery{
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
        fakultasunits, err := mediatr.Send[GetAllFakultasUnits.GetAllFakultasUnitsQuery, fakultasunitdomain.PagedFakultasUnits](context.Background(), query)
        if err != nil {
            return commoninfra.HandleError(c, err)
        }

        return adapter.Send(c, fakultasunits)
    })
}

