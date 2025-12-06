package presentation

import (
    "context"
    "github.com/gofiber/fiber/v2"
    "github.com/mehdihadeli/go-mediatr"
    "strings"

    commoninfra "UnpakSiamida/common/infrastructure"
    commondomain "UnpakSiamida/common/domain"
    indikatorrenstradomain "UnpakSiamida/modules/indikatorrenstra/domain"

    CreateIndikatorRenstra "UnpakSiamida/modules/indikatorrenstra/application/CreateIndikatorRenstra"
    UpdateIndikatorRenstra "UnpakSiamida/modules/indikatorrenstra/application/UpdateIndikatorRenstra"
    DeleteIndikatorRenstra "UnpakSiamida/modules/indikatorrenstra/application/DeleteIndikatorRenstra"
    GetIndikatorRenstra "UnpakSiamida/modules/indikatorrenstra/application/GetIndikatorRenstra"
    GetAllIndikatorRenstras "UnpakSiamida/modules/indikatorrenstra/application/GetAllIndikatorRenstras"
)

func ModuleIndikatorRenstra(app *fiber.App) {

    // ====================================================================
    // CREATE (POST /indikatorrenstra)
    // ====================================================================
    app.Post("/indikatorrenstra", func(c *fiber.Ctx) error {

        standarRenstra := c.FormValue("standar_renstra")
        indikator := c.FormValue("indikator")
        parent := nullableString(c.FormValue("parent"))
        tahun := c.FormValue("tahun")
        tipeTarget := c.FormValue("tipe_target")
        operator := nullableString(c.FormValue("operator"))

        cmd := CreateIndikatorRenstra.CreateIndikatorRenstraCommand{
            StandarRenstra: standarRenstra,
            Indikator:      indikator,
            Parent:         parent,
            Tahun:          tahun,
            TipeTarget:     tipeTarget,
            Operator:       operator,
        }

        uuid, err := mediatr.Send[
            CreateIndikatorRenstra.CreateIndikatorRenstraCommand,
            string,
        ](context.Background(), cmd)

        if err != nil {
            return commoninfra.HandleError(c, err)
        }

        return c.JSON(fiber.Map{"uuid": uuid})
    })

    // ====================================================================
    // UPDATE (PUT /indikatorrenstra/:uuid)
    // ====================================================================
    app.Put("/indikatorrenstra/:uuid", func(c *fiber.Ctx) error {

        uuid := c.Params("uuid")

        standarRenstra := c.FormValue("standar_renstra")
        indikator := c.FormValue("indikator")
        parent := nullableString(c.FormValue("parent"))
        tahun := c.FormValue("tahun")
        tipeTarget := c.FormValue("tipe_target")
        operator := nullableString(c.FormValue("operator"))

        cmd := UpdateIndikatorRenstra.UpdateIndikatorRenstraCommand{
            Uuid:           uuid,
            StandarRenstra: standarRenstra,
            Indikator:      indikator,
            Parent:         parent,
            Tahun:          tahun,
            TipeTarget:     tipeTarget,
            Operator:       operator,
        }

        updatedID, err := mediatr.Send[
            UpdateIndikatorRenstra.UpdateIndikatorRenstraCommand,
            string,
        ](context.Background(), cmd)

        if err != nil {
            return commoninfra.HandleError(c, err)
        }

        return c.JSON(fiber.Map{"uuid": updatedID})
    })

    // ====================================================================
    // DELETE (DELETE /indikatorrenstra/:uuid)
    // ====================================================================
    app.Delete("/indikatorrenstra/:uuid", func(c *fiber.Ctx) error {
        uuid := c.Params("uuid")

        cmd := DeleteIndikatorRenstra.DeleteIndikatorRenstraCommand{Uuid: uuid}

        deletedID, err := mediatr.Send[
            DeleteIndikatorRenstra.DeleteIndikatorRenstraCommand,
            string,
        ](context.Background(), cmd)

        if err != nil {
            return commoninfra.HandleError(c, err)
        }

        return c.JSON(fiber.Map{"uuid": deletedID})
    })

    // ====================================================================
    // GET BY UUID
    // ====================================================================
    app.Get("/indikatorrenstra/:uuid", func(c *fiber.Ctx) error {
        uuid := c.Params("uuid")

        query := GetIndikatorRenstra.GetIndikatorRenstraByUuidQuery{Uuid: uuid}

        indikatorrenstra, err := mediatr.Send[
            GetIndikatorRenstra.GetIndikatorRenstraByUuidQuery,
            *indikatorrenstradomain.IndikatorRenstra,
        ](context.Background(), query)

        if err != nil {
            return commoninfra.HandleError(c, err)
        }

        if indikatorrenstra == nil {
            return c.Status(404).JSON(fiber.Map{"error": "IndikatorRenstra not found"})
        }

        return c.JSON(indikatorrenstra)
    })

    // ====================================================================
    // GET ALL
    // ====================================================================
    app.Get("/indikatorrenstras", func(c *fiber.Ctx) error {
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

        query := GetAllIndikatorRenstras.GetAllIndikatorRenstrasQuery{
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
            GetAllIndikatorRenstras.GetAllIndikatorRenstrasQuery,
            indikatorrenstradomain.PagedIndikatorRenstras,
        ](context.Background(), query)

        if err != nil {
            return commoninfra.HandleError(c, err)
        }

        return adapter.Send(c, result)
    })
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
