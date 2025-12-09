package presentation

import (
    "context"
    "github.com/gofiber/fiber/v2"
    "github.com/mehdihadeli/go-mediatr"
    "strings"

    commoninfra "UnpakSiamida/common/infrastructure"
    commondomain "UnpakSiamida/common/domain"
    JenisFiledomain "UnpakSiamida/modules/jenisfile/domain"

    GetJenisFile "UnpakSiamida/modules/jenisfile/application/GetJenisFile"
    GetAllJenisFiles "UnpakSiamida/modules/jenisfile/application/GetAllJenisFiles"
)

func ModuleJenisFile(app *fiber.App) {

    // ====================================================================
    // GET BY UUID
    // ====================================================================
    app.Get("/JenisFile/:uuid", func(c *fiber.Ctx) error {
        uuid := c.Params("uuid")

        query := GetJenisFile.GetJenisFileByUuidQuery{Uuid: uuid}

        JenisFile, err := mediatr.Send[
            GetJenisFile.GetJenisFileByUuidQuery,
            *JenisFiledomain.JenisFile,
        ](context.Background(), query)

        if err != nil {
            return commoninfra.HandleError(c, err)
        }

        if JenisFile == nil {
            return c.Status(404).JSON(fiber.Map{"error": "JenisFile not found"})
        }

        return c.JSON(JenisFile)
    })

    // ====================================================================
    // GET ALL
    // ====================================================================
    app.Get("/JenisFiles", func(c *fiber.Ctx) error {
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

        query := GetAllJenisFiles.GetAllJenisFilesQuery{
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
            GetAllJenisFiles.GetAllJenisFilesQuery,
            JenisFiledomain.PagedJenisFiles,
        ](context.Background(), query)

        if err != nil {
            return commoninfra.HandleError(c, err)
        }

        return adapter.Send(c, result)
    })
}