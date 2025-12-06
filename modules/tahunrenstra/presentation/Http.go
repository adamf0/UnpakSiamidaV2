package presentation

import (
    "context"
    "github.com/gofiber/fiber/v2"
    "github.com/mehdihadeli/go-mediatr"
    "strings"
    
    // "UnpakSiamida/common/domain"
    commoninfra "UnpakSiamida/common/infrastructure"
    commondomain "UnpakSiamida/common/domain"
    TahunRenstradomain "UnpakSiamida/modules/tahunrenstra/domain"
    GetActiveTahunRenstra "UnpakSiamida/modules/tahunrenstra/application/GetActiveTahunRenstra"
    GetAllTahunRenstras "UnpakSiamida/modules/tahunrenstra/application/GetAllTahunRenstras"
)

func ModuleTahunRenstra(app *fiber.App) {

    // ------------------------------------------------------ 
    // Get Active TahunRenstra
    // ------------------------------------------------------
    app.Get("/tahunrenstra/active", func(c *fiber.Ctx) error {
        query := GetActiveTahunRenstra.GetActiveTahunRenstraQuery{}

        TahunRenstra, err := mediatr.Send[GetActiveTahunRenstra.GetActiveTahunRenstraQuery, *TahunRenstradomain.TahunRenstra](context.Background(), query)
        if err != nil {
            return commoninfra.HandleError(c, err)
        }

        if TahunRenstra == nil {
            return c.Status(404).JSON(fiber.Map{"error": "TahunRenstra not found"})
        }

        return c.JSON(TahunRenstra)
    })

    // ------------------------------------------------------ 
    // Get All TahunRenstra
    // ------------------------------------------------------
    app.Get("/tahunrenstras", func(c *fiber.Ctx) error {
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

        query := GetAllTahunRenstras.GetAllTahunRenstrasQuery{
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
        TahunRenstras, err := mediatr.Send[GetAllTahunRenstras.GetAllTahunRenstrasQuery, TahunRenstradomain.PagedTahunRenstras](context.Background(), query)
        if err != nil {
            return commoninfra.HandleError(c, err)
        }

        return adapter.Send(c, TahunRenstras)
    })
}

