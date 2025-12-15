package presentation

import (
    "context"
    "github.com/gofiber/fiber/v2"
    "github.com/mehdihadeli/go-mediatr"
    "strings"
    
    // "UnpakSiamida/common/domain"
    commoninfra "UnpakSiamida/common/infrastructure"
    commondomain "UnpakSiamida/common/domain"
    templatedokumentambahandomain "UnpakSiamida/modules/templatedokumentambahan/domain"
    CreateTemplateDokumenTambahan "UnpakSiamida/modules/templatedokumentambahan/application/CreateTemplateDokumenTambahan"
    UpdateTemplateDokumenTambahan "UnpakSiamida/modules/templatedokumentambahan/application/UpdateTemplateDokumenTambahan"
    DeleteTemplateDokumenTambahan "UnpakSiamida/modules/templatedokumentambahan/application/DeleteTemplateDokumenTambahan"
    GetTemplateDokumenTambahan "UnpakSiamida/modules/templatedokumentambahan/application/GetTemplateDokumenTambahan"
    GetAllTemplateDokumenTambahans "UnpakSiamida/modules/templatedokumentambahan/application/GetAllTemplateDokumenTambahans"
)

func strPtr(s string) *string {
    if s == "" {
        return nil
    }
    return &s
}

// =======================================================
// POST /templatedokumentambahan
// =======================================================

// CreateTemplateDokumenTambahanHandler godoc
// @Summary Create new TemplateDokumenTambahan
// @Tags TemplateDokumenTambahan
// @Param tahun formData string true "Tahun"
// @Param jenisFile formData string true "Jenis File"
// @Param pertanyaan formData string true "Pertanyaan"
// @Param klasifikasi formData string true "Klasifikasi"
// @Param kategori formData string true "Kategori"
// @Param tugas formData string true "Tugas"
// @Produce json
// @Success 200 {object} map[string]string "uuid of created TemplateDokumenTambahan"
// @Failure 400 {object} commondomain.Error
// @Router /templatedokumentambahan [post]
func CreateTemplateDokumenTambahanHandlerfunc(c *fiber.Ctx) error {

    tahun := c.FormValue("tahun")
    jenisFile := c.FormValue("jenisFile")
    pertanyaan := c.FormValue("pertanyaan")
    klasifikasi := c.FormValue("klasifikasi")
    kategori := c.FormValue("kategori")
    tugas := c.FormValue("tugas")

    cmd := CreateTemplateDokumenTambahan.CreateTemplateDokumenTambahanCommand{
        Tahun : tahun,
        Pertanyaan : pertanyaan,
        JenisFile : jenisFile,
        Klasifikasi : klasifikasi,
        Kategori : kategori,
        Tugas : tugas,
    }

    uuid, err := mediatr.Send[CreateTemplateDokumenTambahan.CreateTemplateDokumenTambahanCommand, string](context.Background(), cmd)
    if err != nil {
        return commoninfra.HandleError(c, err)
    }

    return c.JSON(fiber.Map{"uuid": uuid})
}

// =======================================================
// PUT /templatedokumentambahan/{uuid}
// =======================================================

// UpdateTemplateDokumenTambahanHandler godoc
// @Summary Update existing TemplateDokumenTambahan
// @Tags TemplateDokumenTambahan
// @Param uuid path string true "TemplateDokumenTambahan UUID" format(uuid)
// @Param tahun formData string true "Tahun"
// @Param jenisFile formData string true "Jenis File"
// @Param pertanyaan formData string true "Pertanyaan"
// @Param klasifikasi formData string true "Klasifikasi"
// @Param kategori formData string true "Kategori"
// @Param tugas formData string true "Tugas"
// @Produce json
// @Success 200 {object} map[string]string "uuid of updated TemplateDokumenTambahan"
// @Failure 400 {object} commondomain.Error
// @Router /templatedokumentambahan/{uuid} [put]
func UpdateTemplateDokumenTambahanHandlerfunc(c *fiber.Ctx) error {

    uuid := c.Params("uuid")
    tahun := c.FormValue("tahun")
    jenisFile := c.FormValue("jenisFile")
    pertanyaan := c.FormValue("pertanyaan")
    klasifikasi := c.FormValue("klasifikasi")
    kategori := c.FormValue("kategori")
    tugas := c.FormValue("tugas")

    cmd := UpdateTemplateDokumenTambahan.UpdateTemplateDokumenTambahanCommand{
        Uuid:         uuid,
        Tahun : tahun,
        Pertanyaan : pertanyaan,
        JenisFile : jenisFile,
        Klasifikasi : klasifikasi,
        Kategori : kategori,
        Tugas : tugas,
    }

    updatedID, err := mediatr.Send[UpdateTemplateDokumenTambahan.UpdateTemplateDokumenTambahanCommand, string](context.Background(), cmd)
    if err != nil {
        return commoninfra.HandleError(c, err)
    }

    return c.JSON(fiber.Map{"uuid": updatedID})
}

// =======================================================
// DELETE /templatedokumentambahan/{uuid}
// =======================================================

// DeleteTemplateDokumenTambahanHandler godoc
// @Summary Delete a TemplateDokumenTambahan
// @Tags TemplateDokumenTambahan
// @Param uuid path string true "TemplateDokumenTambahan UUID" format(uuid)
// @Produce json
// @Success 200 {object} map[string]string "uuid of deleted TemplateDokumenTambahan"
// @Failure 404 {object} commondomain.Error
// @Router /templatedokumentambahan/{uuid} [delete]
func DeleteTemplateDokumenTambahanHandlerfunc(c *fiber.Ctx) error {

    uuid := c.Params("uuid")

    cmd := DeleteTemplateDokumenTambahan.DeleteTemplateDokumenTambahanCommand{
        Uuid: uuid,
    }

    deletedID, err := mediatr.Send[DeleteTemplateDokumenTambahan.DeleteTemplateDokumenTambahanCommand, string](context.Background(), cmd)
    if err != nil {
        return commoninfra.HandleError(c, err)
    }

    return c.JSON(fiber.Map{"uuid": deletedID})
}

// =======================================================
// GET /templatedokumentambahan/{uuid}
// =======================================================

// GetTemplateDokumenTambahanHandler godoc
// @Summary Get TemplateDokumenTambahan by UUID
// @Tags TemplateDokumenTambahan
// @Param uuid path string true "TemplateDokumenTambahan UUID" format(uuid)
// @Produce json
// @Success 200 {object} templatedokumentambahandomain.TemplateDokumenTambahan
// @Failure 404 {object} commondomain.Error
// @Router /templatedokumentambahan/{uuid} [get]
func GetTemplateDokumenTambahanHandlerfunc(c *fiber.Ctx) error {
    uuid := c.Params("uuid")

    query := GetTemplateDokumenTambahan.GetTemplateDokumenTambahanByUuidQuery{
        Uuid: uuid,
    }

    templatedokumentambahan, err := mediatr.Send[GetTemplateDokumenTambahan.GetTemplateDokumenTambahanByUuidQuery, *templatedokumentambahandomain.TemplateDokumenTambahan](context.Background(), query)
    if err != nil {
        return commoninfra.HandleError(c, err)
    }

    if templatedokumentambahan == nil {
        return c.Status(404).JSON(fiber.Map{"error": "TemplateDokumenTambahan not found"})
    }

    return c.JSON(templatedokumentambahan)
}

// =======================================================
// GET /templatedokumentambahans
// =======================================================

// GetAllTemplateDokumenTambahansHandler godoc
// @Summary Get All TemplateDokumenTambahans
// @Tags TemplateDokumenTambahan
// @Param mode query string false "paging | all | ndjson | sse"
// @Param page query int false "Page number"
// @Param limit query int false "Limit per page"
// @Param search query string false "Search keyword"
// @Produce json
// @Success 200 {object} templatedokumentambahandomain.PagedTemplateDokumenTambahans
// @Router /templatedokumentambahans [get]
func GetAllTemplateDokumenTambahansHandlerfunc(c *fiber.Ctx) error {
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

    query := GetAllTemplateDokumenTambahans.GetAllTemplateDokumenTambahansQuery{
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
    templatedokumentambahans, err := mediatr.Send[GetAllTemplateDokumenTambahans.GetAllTemplateDokumenTambahansQuery, templatedokumentambahandomain.PagedTemplateDokumenTambahans](context.Background(), query)
    if err != nil {
        return commoninfra.HandleError(c, err)
    }

    return adapter.Send(c, templatedokumentambahans)
}

func ModuleTemplateDokumenTambahan(app *fiber.App) {
    app.Post("/templatedokumentambahan", CreateTemplateDokumenTambahanHandlerfunc)
    app.Put("/templatedokumentambahan/:uuid", UpdateTemplateDokumenTambahanHandlerfunc)
    app.Delete("/templatedokumentambahan/:uuid", DeleteTemplateDokumenTambahanHandlerfunc)
    app.Get("/templatedokumentambahan/:uuid", GetTemplateDokumenTambahanHandlerfunc)
    app.Get("/templatedokumentambahans", GetAllTemplateDokumenTambahansHandlerfunc)
}

