package presentation

import (
    "context"
    "github.com/gofiber/fiber/v2"
    "github.com/mehdihadeli/go-mediatr"
    "strings"
    
    // "UnpakSiamida/common/domain"
    commoninfra "UnpakSiamida/common/infrastructure"
    commondomain "UnpakSiamida/common/domain"
    templaterenstradomain "UnpakSiamida/modules/templaterenstra/domain"
    CreateTemplateRenstra "UnpakSiamida/modules/templaterenstra/application/CreateTemplateRenstra"
    UpdateTemplateRenstra "UnpakSiamida/modules/templaterenstra/application/UpdateTemplateRenstra"
    DeleteTemplateRenstra "UnpakSiamida/modules/templaterenstra/application/DeleteTemplateRenstra"
    GetTemplateRenstra "UnpakSiamida/modules/templaterenstra/application/GetTemplateRenstra"
    GetAllTemplateRenstras "UnpakSiamida/modules/templaterenstra/application/GetAllTemplateRenstras"
)

func strPtr(s string) *string {
    if s == "" {
        return nil
    }
    return &s
}

func ModuleTemplateRenstra(app *fiber.App) {

    // ------------------------------------------------------
    // CREATE templaterenstra (POST /templaterenstra)
    // ------------------------------------------------------
    app.Post("/templaterenstra", func(c *fiber.Ctx) error {

        tahun := c.FormValue("tahun")
        indikator := c.FormValue("indikator")
        isPertanyaan := c.FormValue("isPertanyaan")
        fakultasUnit := c.FormValue("fakultasUnit")
        kategori := c.FormValue("kategori")
        klasifikasi := c.FormValue("klasifikasi")
        satuan := c.FormValue("satuan")
        target := c.FormValue("target")
        targetMin := c.FormValue("targetMin")
        targetMax := c.FormValue("targetMax")
        tugas := c.FormValue("tugas")

        cmd := CreateTemplateRenstra.CreateTemplateRenstraCommand{
            Tahun : tahun,
            Indikator : indikator,
            IsPertanyaan : isPertanyaan,
            FakultasUnit : fakultasUnit,
            Kategori : kategori,
            Klasifikasi : klasifikasi,
            Satuan:        strPtr(satuan),
            Target:        strPtr(target),
            TargetMin:     strPtr(targetMin),
            TargetMax:     strPtr(targetMax),
            Tugas : tugas,
        }

        uuid, err := mediatr.Send[CreateTemplateRenstra.CreateTemplateRenstraCommand, string](context.Background(), cmd)
        if err != nil {
            return commoninfra.HandleError(c, err)
        }

        return c.JSON(fiber.Map{"uuid": uuid})
    })

    // ------------------------------------------------------
    // UPDATE templaterenstra (PUT /templaterenstra/:uuid)
    // ------------------------------------------------------
    app.Put("/templaterenstra/:uuid", func(c *fiber.Ctx) error {

        uuid := c.Params("uuid")
        tahun := c.FormValue("tahun")
        indikator := c.FormValue("indikator")
        isPertanyaan := c.FormValue("isPertanyaan")
        fakultasUnit := c.FormValue("fakultasUnit")
        kategori := c.FormValue("kategori")
        klasifikasi := c.FormValue("klasifikasi")
        satuan := c.FormValue("satuan")
        target := c.FormValue("target")
        targetMin := c.FormValue("targetMin")
        targetMax := c.FormValue("targetMax")
        tugas := c.FormValue("tugas")

        cmd := UpdateTemplateRenstra.UpdateTemplateRenstraCommand{
            Uuid:         uuid,
            Tahun : tahun,
            Indikator : indikator,
            IsPertanyaan : isPertanyaan,
            FakultasUnit : fakultasUnit,
            Kategori : kategori,
            Klasifikasi : klasifikasi,
            Satuan:        strPtr(satuan),
            Target:        strPtr(target),
            TargetMin:     strPtr(targetMin),
            TargetMax:     strPtr(targetMax),
            Tugas : tugas,
        }

        updatedID, err := mediatr.Send[UpdateTemplateRenstra.UpdateTemplateRenstraCommand, string](context.Background(), cmd)
        if err != nil {
            return commoninfra.HandleError(c, err)
        }

        return c.JSON(fiber.Map{"uuid": updatedID})
    })

    // ------------------------------------------------------ 
    // DELETE templaterenstra (DELETE /templaterenstra/:uuid)
    // ------------------------------------------------------
    app.Delete("/templaterenstra/:uuid", func(c *fiber.Ctx) error {

        uuid := c.Params("uuid")

        cmd := DeleteTemplateRenstra.DeleteTemplateRenstraCommand{
            Uuid: uuid,
        }

        deletedID, err := mediatr.Send[DeleteTemplateRenstra.DeleteTemplateRenstraCommand, string](context.Background(), cmd)
        if err != nil {
            return commoninfra.HandleError(c, err)
        }

        return c.JSON(fiber.Map{"uuid": deletedID})
    })

    // ------------------------------------------------------ 
    // Get templaterenstra
    // ------------------------------------------------------
    app.Get("/templaterenstra/:uuid", func(c *fiber.Ctx) error {
        uuid := c.Params("uuid")

        query := GetTemplateRenstra.GetTemplateRenstraByUuidQuery{
            Uuid: uuid,
        }

        templaterenstra, err := mediatr.Send[GetTemplateRenstra.GetTemplateRenstraByUuidQuery, *templaterenstradomain.TemplateRenstra](context.Background(), query)
        if err != nil {
            return commoninfra.HandleError(c, err)
        }

        if templaterenstra == nil {
            return c.Status(404).JSON(fiber.Map{"error": "TemplateRenstra not found"})
        }

        return c.JSON(templaterenstra)
    })

    // ------------------------------------------------------ 
    // Get All templaterenstra
    // ------------------------------------------------------
    app.Get("/templaterenstras", func(c *fiber.Ctx) error {
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

        query := GetAllTemplateRenstras.GetAllTemplateRenstrasQuery{
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
        templaterenstras, err := mediatr.Send[GetAllTemplateRenstras.GetAllTemplateRenstrasQuery, templaterenstradomain.PagedTemplateRenstras](context.Background(), query)
        if err != nil {
            return commoninfra.HandleError(c, err)
        }

        return adapter.Send(c, templaterenstras)
    })
}

