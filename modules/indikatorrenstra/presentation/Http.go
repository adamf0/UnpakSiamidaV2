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
    // GetTreeIndikatorRenstraByTahun "UnpakSiamida/modules/indikatorrenstra/application/GetTreeIndikatorRenstraByTahun"
    GetAllIndikatorRenstras "UnpakSiamida/modules/indikatorrenstra/application/GetAllIndikatorRenstras"
    SetupUuidIndikatorRenstra "UnpakSiamida/modules/indikatorrenstra/application/SetupUuidIndikatorRenstra"
)

// CreateIndikatorRenstraHandler godoc
// @Summary Create new IndikatorRenstra
// @Tags IndikatorRenstra
// @Param standar_renstra formData string true "Standar Renstra ID"
// @Param indikator formData string true "Indikator Name"
// @Param parent formData string false "Parent ID"
// @Param tahun formData string true "Tahun"
// @Param tipe_target formData string true "Tipe Target"
// @Param operator formData string false "Operator"
// @Produce json
// @Success 200 {object} map[string]string "uuid of created IndikatorRenstra"
// @Failure 400 {object} commondomain.Error
// @Router /indikatorrenstra [post]
func CreateIndikatorRenstraHandlerfunc(c *fiber.Ctx) error {
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
}

// =======================================================
// PUT /indikatorrenstra/{uuid}
// =======================================================

// UpdateIndikatorRenstraHandler godoc
// @Summary Update existing IndikatorRenstra
// @Tags IndikatorRenstra
// @Param uuid path string true "IndikatorRenstra UUID" format(uuid)
// @Param standar_renstra formData string true "Standar Renstra ID"
// @Param indikator formData string true "Indikator Name"
// @Param parent formData string false "Parent ID"
// @Param tahun formData string true "Tahun"
// @Param tipe_target formData string true "Tipe Target"
// @Param operator formData string false "Operator"
// @Produce json
// @Success 200 {object} map[string]string "uuid of updated IndikatorRenstra"
// @Failure 400 {object} commondomain.Error
// @Router /indikatorrenstra/{uuid} [put]
func UpdateIndikatorRenstraHandlerfunc(c *fiber.Ctx) error {
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
}

// =======================================================
// DELETE /indikatorrenstra/{uuid}
// =======================================================

// DeleteIndikatorRenstraHandler godoc
// @Summary Delete an IndikatorRenstra
// @Tags IndikatorRenstra
// @Param uuid path string true "IndikatorRenstra UUID" format(uuid)
// @Produce json
// @Success 200 {object} map[string]string "uuid of deleted IndikatorRenstra"
// @Failure 404 {object} commondomain.Error
// @Router /indikatorrenstra/{uuid} [delete]
func DeleteIndikatorRenstraHandlerfunc(c *fiber.Ctx) error {
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
}

// =======================================================
// GET /indikatorrenstra/{uuid}
// =======================================================

// GetIndikatorRenstraHandler godoc
// @Summary Get IndikatorRenstra by UUID
// @Tags IndikatorRenstra
// @Param uuid path string true "IndikatorRenstra UUID" format(uuid)
// @Produce json
// @Success 200 {object} indikatorrenstradomain.IndikatorRenstra
// @Failure 404 {object} commondomain.Error
// @Router /indikatorrenstra/{uuid} [get]
func GetIndikatorRenstraHandlerfunc(c *fiber.Ctx) error {
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
}

// =======================================================
// GET /indikatorrenstra/tree/{tahun}
// =======================================================

// GetTreeIndikatorRenstraByTahunHandler godoc
// @Summary Get IndikatorRenstra by UUID
// @Tags IndikatorRenstra
// @Param uuidTahun path string true "tahun UUID" format(uuid)
// @Produce json
// @Success 200 {object} indikatorrenstradomain.IndikatorRenstra
// @Failure 404 {object} commondomain.Error
// @Router /indikatorrenstra/{uuid} [get]
// func GetTreeIndikatorRenstraByTahunHandler(c *fiber.Ctx) error {
//     uuidTahun := c.Params("uuidTahun")

//     query := GetTreeIndikatorRenstraByTahun.GetTreeIndikatorRenstraByTahunByUuidQuery{UuidTahun: uuidTahun}

//     indikatorrenstra, err := mediatr.Send[
//         GetTreeIndikatorRenstraByTahun.GetTreeIndikatorRenstraByTahunByUuidQuery,
//         *indikatorrenstradomain.IndikatorRenstra,
//     ](context.Background(), query)

//     if err != nil {
//         return commoninfra.HandleError(c, err)
//     }

//     if indikatorrenstra == nil {
//         return c.Status(404).JSON(fiber.Map{"error": "IndikatorRenstra not found"})
//     }

//     return c.JSON(indikatorrenstra)
// }

// GetAllIndikatorRenstrasHandler godoc
// @Summary Get All IndikatorRenstras
// @Tags IndikatorRenstra
// @Param mode query string false "paging | all | ndjson | sse"
// @Param page query int false "Page number"
// @Param limit query int false "Limit per page"
// @Param search query string false "Search keyword"
// @Produce json
// @Success 200 {object} indikatorrenstradomain.PagedIndikatorRenstras
// @Router /indikatorrenstras [get]
func GetAllIndikatorRenstrasHandlerfunc(c *fiber.Ctx) error {
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
}

func SetupUuidIndikatorRenstrasHandlerfunc(c *fiber.Ctx) error {
    cmd := SetupUuidIndikatorRenstra.SetupUuidIndikatorRenstraCommand{}

    message, err := mediatr.Send[SetupUuidIndikatorRenstra.SetupUuidIndikatorRenstraCommand, string](context.Background(), cmd)
    if err != nil {
        return commoninfra.HandleError(c, err)
    }

    return c.JSON(fiber.Map{"message": message})
}

func ModuleIndikatorRenstra(app *fiber.App) {
    app.Get("/indikatorrenstra/setupuuid", SetupUuidIndikatorRenstrasHandlerfunc)

    app.Post("/indikatorrenstra", CreateIndikatorRenstraHandlerfunc)
    app.Put("/indikatorrenstra/:uuid", UpdateIndikatorRenstraHandlerfunc)
    app.Delete("/indikatorrenstra/:uuid", DeleteIndikatorRenstraHandlerfunc)
    app.Get("/indikatorrenstra/:uuid", GetIndikatorRenstraHandlerfunc)
    app.Get("/indikatorrenstras", GetAllIndikatorRenstrasHandlerfunc)

    // app.Get("/indikatorrenstra/tree/:uuidTahun", GetTreeIndikatorRenstraByTahunHandlerfunc)
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
