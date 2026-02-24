package presentation

import (
	"context"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/mehdihadeli/go-mediatr"

	// "UnpakSiamida/common/domain"
	commondomain "UnpakSiamida/common/domain"
	"UnpakSiamida/common/helper"
	commoninfra "UnpakSiamida/common/infrastructure"
	commonpresentation "UnpakSiamida/common/presentation"

	CreateTemplateRenstra "UnpakSiamida/modules/templaterenstra/application/CreateTemplateRenstra"
	DeleteTemplateRenstra "UnpakSiamida/modules/templaterenstra/application/DeleteTemplateRenstra"
	GetAllTemplateRenstras "UnpakSiamida/modules/templaterenstra/application/GetAllTemplateRenstras"
	GetTemplateRenstra "UnpakSiamida/modules/templaterenstra/application/GetTemplateRenstra"
	SetupUuidTemplateRenstra "UnpakSiamida/modules/templaterenstra/application/SetupUuidTemplateRenstra"
	UpdateTemplateRenstra "UnpakSiamida/modules/templaterenstra/application/UpdateTemplateRenstra"
	templaterenstradomain "UnpakSiamida/modules/templaterenstra/domain"
)

// =======================================================
// POST /templaterenstra
// =======================================================

// CreateTemplateRenstraHandler godoc
// @Summary Create new TemplateRenstra
// @Tags TemplateRenstra
// @Param tahun formData string true "Tahun"
// @Param indikator formData string true "Indikator"
// @Param isPertanyaan formData string true "Is Pertanyaan"
// @Param fakultasUnit formData string true "Fakultas Unit"
// @Param kategori formData string true "Kategori"
// @Param klasifikasi formData string true "Klasifikasi"
// @Param satuan formData string false "Satuan"
// @Param target formData string false "Target"
// @Param targetMin formData string false "Target Min"
// @Param targetMax formData string false "Target Max"
// @Param tugas formData string true "Tugas"
// @Produce json
// @Success 200 {object} map[string]string "uuid of created TemplateRenstra"
// @Failure 400 {object} commoninfra.ResponseError
// @Failure 404 {object} commoninfra.ResponseError
// @Failure 409 {object} commoninfra.ResponseError
// @Failure 500 {object} commoninfra.ResponseError
// @Router /templaterenstra [post]
func CreateTemplateRenstraHandlerfunc(c *fiber.Ctx) error {

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
		Tahun:        tahun,
		Indikator:    indikator,
		IsPertanyaan: isPertanyaan,
		FakultasUnit: fakultasUnit,
		Kategori:     kategori,
		Klasifikasi:  klasifikasi,
		Satuan:       helper.StrPtr(satuan),
		Target:       helper.StrPtr(target),
		TargetMin:    helper.StrPtr(targetMin),
		TargetMax:    helper.StrPtr(targetMax),
		Tugas:        tugas,
	}

	uuid, err := mediatr.Send[CreateTemplateRenstra.CreateTemplateRenstraCommand, string](context.Background(), cmd)
	if err != nil {
		return commoninfra.HandleError(c, err)
	}

	return commonpresentation.JsonUUID(c, uuid)
}

// =======================================================
// PUT /templaterenstra/{uuid}
// =======================================================

// UpdateTemplateRenstraHandler godoc
// @Summary Update existing TemplateRenstra
// @Tags TemplateRenstra
// @Param uuid path string true "TemplateRenstra UUID" format(uuid)
// @Param tahun formData string true "Tahun"
// @Param indikator formData string true "Indikator"
// @Param isPertanyaan formData string true "Is Pertanyaan"
// @Param fakultasUnit formData string true "Fakultas Unit"
// @Param kategori formData string true "Kategori"
// @Param klasifikasi formData string true "Klasifikasi"
// @Param satuan formData string false "Satuan"
// @Param target formData string false "Target"
// @Param targetMin formData string false "Target Min"
// @Param targetMax formData string false "Target Max"
// @Param tugas formData string true "Tugas"
// @Produce json
// @Success 200 {object} map[string]string "uuid of updated TemplateRenstra"
// @Failure 400 {object} commoninfra.ResponseError
// @Failure 404 {object} commoninfra.ResponseError
// @Failure 409 {object} commoninfra.ResponseError
// @Failure 500 {object} commoninfra.ResponseError
// @Router /templaterenstra/{uuid} [put]
func UpdateTemplateRenstraHandlerfunc(c *fiber.Ctx) error {

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
		Tahun:        tahun,
		Indikator:    indikator,
		IsPertanyaan: isPertanyaan,
		FakultasUnit: fakultasUnit,
		Kategori:     kategori,
		Klasifikasi:  klasifikasi,
		Satuan:       helper.StrPtr(satuan),
		Target:       helper.StrPtr(target),
		TargetMin:    helper.StrPtr(targetMin),
		TargetMax:    helper.StrPtr(targetMax),
		Tugas:        tugas,
	}

	updatedID, err := mediatr.Send[UpdateTemplateRenstra.UpdateTemplateRenstraCommand, string](context.Background(), cmd)
	if err != nil {
		return commoninfra.HandleError(c, err)
	}

	return c.JSON(fiber.Map{"uuid": updatedID})
}

// =======================================================
// DELETE /templaterenstra/{uuid}
// =======================================================

// DeleteTemplateRenstraHandler godoc
// @Summary Delete a TemplateRenstra
// @Tags TemplateRenstra
// @Param uuid path string true "TemplateRenstra UUID" format(uuid)
// @Produce json
// @Success 200 {object} map[string]string "uuid of deleted TemplateRenstra"
// @Failure 400 {object} commoninfra.ResponseError
// @Failure 404 {object} commoninfra.ResponseError
// @Failure 409 {object} commoninfra.ResponseError
// @Failure 500 {object} commoninfra.ResponseError
//
//	@Router /templaterenstra/{uuid} [delete]
func DeleteTemplateRenstraHandlerfunc(c *fiber.Ctx) error {

	uuid := c.Params("uuid")

	cmd := DeleteTemplateRenstra.DeleteTemplateRenstraCommand{
		Uuid: uuid,
	}

	deletedID, err := mediatr.Send[DeleteTemplateRenstra.DeleteTemplateRenstraCommand, string](context.Background(), cmd)
	if err != nil {
		return commoninfra.HandleError(c, err)
	}

	return c.JSON(fiber.Map{"uuid": deletedID})
}

// =======================================================
// GET /templaterenstra/{uuid}
// =======================================================

// GetTemplateRenstraHandler godoc
// @Summary Get TemplateRenstra by UUID
// @Tags TemplateRenstra
// @Param uuid path string true "TemplateRenstra UUID" format(uuid)
// @Produce json
// @Success 200 {object} templaterenstradomain.TemplateRenstra
// @Failure 400 {object} commoninfra.ResponseError
// @Failure 404 {object} commoninfra.ResponseError
// @Failure 409 {object} commoninfra.ResponseError
// @Failure 500 {object} commoninfra.ResponseError
//
//	@Router /templaterenstra/{uuid} [get]
func GetTemplateRenstraHandlerfunc(c *fiber.Ctx) error {
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
}

// =======================================================
// GET /templaterenstras
// =======================================================

// GetAllTemplateRenstrasHandler godoc
// @Summary Get All TemplateRenstras
// @Tags TemplateRenstra
// @Param mode query string false "paging | all | ndjson | sse"
// @Param page query int false "Page number"
// @Param limit query int false "Limit per page"
// @Param search query string false "Search keyword"
// @Produce json
// @Success 200 {object} commondomain.Paged[templaterenstradomain.TemplateRenstraDefault]
// @Router /templaterenstras [get]
func GetAllTemplateRenstrasHandlerfunc(c *fiber.Ctx) error {
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
	var adapter commonpresentation.OutputAdapter[templaterenstradomain.TemplateRenstraDefault]
	switch mode {
	case "all":
		adapter = &commonpresentation.AllAdapter[templaterenstradomain.TemplateRenstraDefault]{}
	case "ndjson":
		adapter = &commonpresentation.NDJSONAdapter[templaterenstradomain.TemplateRenstraDefault]{}
	case "sse":
		adapter = &commonpresentation.SSEAdapter[templaterenstradomain.TemplateRenstraDefault]{}
	default:
		query.Page = &page
		query.Limit = &limit
		adapter = &commonpresentation.PagingAdapter[templaterenstradomain.TemplateRenstraDefault]{}
	}

	// Ambil data
	templaterenstras, err := mediatr.Send[GetAllTemplateRenstras.GetAllTemplateRenstrasQuery, commondomain.Paged[templaterenstradomain.TemplateRenstraDefault]](context.Background(), query)
	if err != nil {
		return commoninfra.HandleError(c, err)
	}

	return adapter.Send(c, templaterenstras)
}

func SetupUuidTemplateRenstrasHandlerfunc(c *fiber.Ctx) error {
	cmd := SetupUuidTemplateRenstra.SetupUuidTemplateRenstraCommand{}

	message, err := mediatr.Send[SetupUuidTemplateRenstra.SetupUuidTemplateRenstraCommand, string](context.Background(), cmd)
	if err != nil {
		return commoninfra.HandleError(c, err)
	}

	return c.JSON(fiber.Map{"message": message})
}

func ModuleTemplateRenstra(app *fiber.App) {
	admin := []string{"admin"}
	whoamiURL := "http://localhost:3000/whoami"

	app.Get("/templaterenstra/setupuuid", commonpresentation.JWTMiddleware(), commonpresentation.RBACMiddleware(admin, whoamiURL), SetupUuidTemplateRenstrasHandlerfunc)

	app.Post("/templaterenstra", commonpresentation.JWTMiddleware(), commonpresentation.RBACMiddleware(admin, whoamiURL), CreateTemplateRenstraHandlerfunc)
	// app.Put("/templaterenstra/:uuid", commonpresentation.JWTMiddleware(), UpdateTemplateRenstraHandlerfunc)
	app.Delete("/templaterenstra/:uuid", commonpresentation.JWTMiddleware(), commonpresentation.RBACMiddleware(admin, whoamiURL), DeleteTemplateRenstraHandlerfunc)
	app.Get("/templaterenstra/:uuid", commonpresentation.SmartCompress(), commonpresentation.JWTMiddleware(), GetTemplateRenstraHandlerfunc)
	app.Get("/templaterenstras", commonpresentation.SmartCompress(), commonpresentation.JWTMiddleware(), GetAllTemplateRenstrasHandlerfunc)
}
