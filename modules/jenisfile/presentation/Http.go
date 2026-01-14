package presentation

import (
    "context"
    "github.com/gofiber/fiber/v2"
    "github.com/mehdihadeli/go-mediatr"
    "strings"

    commoninfra "UnpakSiamida/common/infrastructure"
    commonpresentation "UnpakSiamida/common/presentation"
    commondomain "UnpakSiamida/common/domain"
    JenisFiledomain "UnpakSiamida/modules/jenisfile/domain"

    CreateJenisFile "UnpakSiamida/modules/jenisfile/application/CreateJenisFile"
    UpdateJenisFile "UnpakSiamida/modules/jenisfile/application/UpdateJenisFile"
    DeleteJenisFile "UnpakSiamida/modules/jenisfile/application/DeleteJenisFile"
    GetJenisFile "UnpakSiamida/modules/jenisfile/application/GetJenisFile"
    GetAllJenisFiles "UnpakSiamida/modules/jenisfile/application/GetAllJenisFiles"
    SetupUuidJenisFile "UnpakSiamida/modules/jenisfile/application/SetupUuidJenisFile"
)

// =======================================================
// POST /jenisfile
// =======================================================

// CreateJenisFileHandler godoc
// @Summary Create new JenisFile
// @Tags JenisFile
// @Param nama formData string true "Nama JenisFile"
// @Produce json
// @Success 200 {object} map[string]string "uuid of created JenisFile"
// @Failure 400 {object} commondomain.Error
// @Router /jenisfile [post]
func CreateJenisFileHandlerfunc(c *fiber.Ctx) error {

    nama := c.FormValue("nama")

    cmd := CreateJenisFile.CreateJenisFileCommand{
        Nama:         nama,
    }

    uuid, err := mediatr.Send[CreateJenisFile.CreateJenisFileCommand, string](context.Background(), cmd)
    if err != nil {
        return commoninfra.HandleError(c, err)
    }

    return c.JSON(fiber.Map{"uuid": uuid})
}

// =======================================================
// PUT /jenisfile/{uuid}
// =======================================================

// UpdateJenisFileHandler godoc
// @Summary Update existing JenisFile
// @Tags JenisFile
// @Param uuid path string true "JenisFile UUID" format(uuid)
// @Param nama formData string true "Nama JenisFile"
// @Produce json
// @Success 200 {object} map[string]string "uuid of updated JenisFile"
// @Failure 400 {object} commondomain.Error
// @Router /jenisfile/{uuid} [put]
func UpdateJenisFileHandlerfunc(c *fiber.Ctx) error {

    uuid := c.Params("uuid")

    nama := c.FormValue("nama")

    cmd := UpdateJenisFile.UpdateJenisFileCommand{
        Uuid:         uuid,
        Nama:         nama,
    }

    updatedID, err := mediatr.Send[UpdateJenisFile.UpdateJenisFileCommand, string](context.Background(), cmd)
    if err != nil {
        return commoninfra.HandleError(c, err)
    }

    return c.JSON(fiber.Map{"uuid": updatedID})
}

// =======================================================
// DELETE /jenisfile/{uuid}
// =======================================================

// DeleteJenisFileHandler godoc
// @Summary Delete a JenisFile
// @Tags JenisFile
// @Param uuid path string true "JenisFile UUID" format(uuid)
// @Produce json
// @Success 200 {object} map[string]string "uuid of deleted JenisFile"
// @Failure 404 {object} commondomain.Error
// @Router /jenisfile/{uuid} [delete]
func DeleteJenisFileHandlerfunc(c *fiber.Ctx) error {

    uuid := c.Params("uuid")

    cmd := DeleteJenisFile.DeleteJenisFileCommand{
        Uuid: uuid,
    }

    deletedID, err := mediatr.Send[DeleteJenisFile.DeleteJenisFileCommand, string](context.Background(), cmd)
    if err != nil {
        return commoninfra.HandleError(c, err)
    }

    return c.JSON(fiber.Map{"uuid": deletedID})
}

// =======================================================
// GET /JenisFile/{uuid}
// =======================================================

// GetJenisFileHandler godoc
// @Summary Get JenisFile by UUID
// @Tags JenisFile
// @Param uuid path string true "JenisFile UUID" format(uuid)
// @Produce json
// @Success 200 {object} JenisFiledomain.JenisFile
// @Failure 404 {object} commondomain.Error
// @Router /JenisFile/{uuid} [get]
func GetJenisFileHandlerfunc(c *fiber.Ctx) error {
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
}

// =======================================================
// GET /JenisFiles
// =======================================================

// GetAllJenisFilesHandler godoc
// @Summary Get all JenisFiles
// @Tags JenisFile
// @Param mode query string false "paging | all | ndjson | sse"
// @Param page query int false "Page number"
// @Param limit query int false "Limit per page"
// @Param search query string false "Search keyword"
// @Produce json
// @Success 200 {object} JenisFiledomain.PagedJenisFiles
// @Router /JenisFiles [get]
func GetAllJenisFilesHandlerfunc(c *fiber.Ctx) error {
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
}

func SetupUuidJenisFilesHandlerfunc(c *fiber.Ctx) error {
    cmd := SetupUuidJenisFile.SetupUuidJenisFileCommand{}

    message, err := mediatr.Send[SetupUuidJenisFile.SetupUuidJenisFileCommand, string](context.Background(), cmd)
    if err != nil {
        return commoninfra.HandleError(c, err)
    }

    return c.JSON(fiber.Map{"message": message})
}

func ModuleJenisFile(app *fiber.App) {
    admin := []string{"admin"}
	whoamiURL := "http://localhost:3000/whoami"

    app.Get("/jenisfile/setupuuid", commonpresentation.JWTMiddleware(), commonpresentation.RBACMiddleware(admin, whoamiURL), SetupUuidJenisFilesHandlerfunc)

    app.Post("/jenisfile", commonpresentation.JWTMiddleware(), commonpresentation.RBACMiddleware(admin, whoamiURL), CreateJenisFileHandlerfunc)
    app.Put("/jenisfile/:uuid", commonpresentation.JWTMiddleware(), commonpresentation.RBACMiddleware(admin, whoamiURL), UpdateJenisFileHandlerfunc)
    app.Delete("/jenisfile/:uuid", commonpresentation.JWTMiddleware(), commonpresentation.RBACMiddleware(admin, whoamiURL), DeleteJenisFileHandlerfunc)
    app.Get("/JenisFile/:uuid", commonpresentation.JWTMiddleware(), GetJenisFileHandlerfunc)
    app.Get("/JenisFiles", commonpresentation.JWTMiddleware(), GetAllJenisFilesHandlerfunc)
}