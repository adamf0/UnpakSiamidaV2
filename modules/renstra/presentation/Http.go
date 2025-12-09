package presentation

import (
    "context"
    "github.com/gofiber/fiber/v2"
    "github.com/mehdihadeli/go-mediatr"
    "strings"

    commoninfra "UnpakSiamida/common/infrastructure"
    commondomain "UnpakSiamida/common/domain"
    renstradomain "UnpakSiamida/modules/renstra/domain"

    CreateRenstra "UnpakSiamida/modules/renstra/application/CreateRenstra"
    UpdateRenstra "UnpakSiamida/modules/renstra/application/UpdateRenstra"
    GiveCodeRenstra "UnpakSiamida/modules/renstra/application/GiveCodeAccessRenstra"
    DeleteRenstra "UnpakSiamida/modules/renstra/application/DeleteRenstra"
    GetRenstraDefault "UnpakSiamida/modules/renstra/application/GetRenstraDefault"
    GetAllRenstras "UnpakSiamida/modules/renstra/application/GetAllRenstras"
)

func ModuleRenstra(app *fiber.App) {

    // ====================================================================
    // CREATE (POST /renstra)
    // ====================================================================
    app.Post("/renstra", func(c *fiber.Ctx) error {
        cmd := CreateRenstra.CreateRenstraCommand{
            Tahun:                          c.FormValue("tahun"),
            FakultasUnit:                   c.FormValue("fakultas_unit"),
            PeriodeUploadMulai:             c.FormValue("periode_upload_mulai"),
            PeriodeUploadAkhir:             c.FormValue("periode_upload_akhir"),
            PeriodeAssesmentDokumenMulai:   c.FormValue("periode_assesment_dokumen_mulai"),
            PeriodeAssesmentDokumenAkhir:   c.FormValue("periode_assesment_dokumen_akhir"),
            PeriodeAssesmentLapanganMulai:  c.FormValue("periode_assesment_lapangan_mulai"),
            PeriodeAssesmentLapanganAkhir:  c.FormValue("periode_assesment_lapangan_akhir"),
            Auditee:                        c.FormValue("auditee"),
            Auditor1:                       c.FormValue("auditor1"),
            Auditor2:                       c.FormValue("auditor2"),
        }

        // Kirim ke mediator
        uuid, err := mediatr.Send[
            CreateRenstra.CreateRenstraCommand,
            string,
        ](context.Background(), cmd)

        if err != nil {
            return commoninfra.HandleError(c, err)
        }

        return c.JSON(fiber.Map{"uuid": uuid})
    })

    // ====================================================================
    // UPDATE (PUT /renstra/:uuid)
    // ====================================================================
    app.Put("/renstra/:uuid", func(c *fiber.Ctx) error {
        cmd := UpdateRenstra.UpdateRenstraCommand{
            Uuid:                           c.Params("uuid"),
            Tahun:                          c.FormValue("tahun"),
            FakultasUnit:                   c.FormValue("fakultas_unit"),
            PeriodeUploadMulai:             c.FormValue("periode_upload_mulai"),
            PeriodeUploadAkhir:             c.FormValue("periode_upload_akhir"),
            PeriodeAssesmentDokumenMulai:   c.FormValue("periode_assesment_dokumen_mulai"),
            PeriodeAssesmentDokumenAkhir:   c.FormValue("periode_assesment_dokumen_akhir"),
            PeriodeAssesmentLapanganMulai:  c.FormValue("periode_assesment_lapangan_mulai"),
            PeriodeAssesmentLapanganAkhir:  c.FormValue("periode_assesment_lapangan_akhir"),
            Auditee:                        c.FormValue("auditee"),
            Auditor1:                       c.FormValue("auditor1"),
            Auditor2:                       c.FormValue("auditor2"),
        }

        updatedID, err := mediatr.Send[
            UpdateRenstra.UpdateRenstraCommand,
            string,
        ](context.Background(), cmd)

        if err != nil {
            return commoninfra.HandleError(c, err)
        }

        return c.JSON(fiber.Map{"uuid": updatedID})
    })

    // ====================================================================
    // UPDATE (PUT /renstra/:uuid/code_access)
    // ====================================================================
    app.Put("/renstra/:uuid/code_access", func(c *fiber.Ctx) error {
        cmd := GiveCodeRenstra.GiveCodeAccessRenstraCommand{
            Uuid:                          c.Params("uuid"),
            KodeAkses:                     c.FormValue("kodeAkses"),
        }

        updatedID, err := mediatr.Send[
            GiveCodeRenstra.GiveCodeAccessRenstraCommand,
            string,
        ](context.Background(), cmd)

        if err != nil {
            return commoninfra.HandleError(c, err)
        }

        return c.JSON(fiber.Map{"uuid": updatedID})
    })

    // ====================================================================
    // DELETE (DELETE /renstra/:uuid)
    // ====================================================================
    app.Delete("/renstra/:uuid", func(c *fiber.Ctx) error {
        uuid := c.Params("uuid")

        cmd := DeleteRenstra.DeleteRenstraCommand{Uuid: uuid}

        deletedID, err := mediatr.Send[
            DeleteRenstra.DeleteRenstraCommand,
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
    app.Get("/renstra/:uuid", func(c *fiber.Ctx) error {
        uuid := c.Params("uuid")

        query := GetRenstraDefault.GetRenstraDefaultByUuidQuery{Uuid: uuid}

        renstra, err := mediatr.Send[
            GetRenstraDefault.GetRenstraDefaultByUuidQuery,
            *renstradomain.RenstraDefault,
        ](context.Background(), query)

        if err != nil {
            return commoninfra.HandleError(c, err)
        }

        if renstra == nil {
            return c.Status(404).JSON(fiber.Map{"error": "Renstra not found"})
        }

        return c.JSON(renstra)
    })

    // ====================================================================
    // GET ALL
    // ====================================================================
    app.Get("/renstras", func(c *fiber.Ctx) error {
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

        query := GetAllRenstras.GetAllRenstrasQuery{
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
            GetAllRenstras.GetAllRenstrasQuery,
            renstradomain.PagedRenstras,
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
// func nullableString(s string) *string {
//     trimmed := strings.TrimSpace(s)
//     if trimmed == "" {
//         return nil
//     }
//     return &trimmed
// }
