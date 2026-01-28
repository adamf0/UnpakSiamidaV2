package presentation

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/mehdihadeli/go-mediatr"

	commoninfra "UnpakSiamida/common/infrastructure"
	commonpresentation "UnpakSiamida/common/presentation"

	// commondomain "UnpakSiamida/common/domain"

	GetPreviewTemplate "UnpakSiamida/modules/previewtemplate/application/GetPreviewTemplate"
	previewtemplatedomain "UnpakSiamida/modules/previewtemplate/domain"
)

// =======================================================
// GET /preview/audit/{tahun}/{fakultasUnit}
// =======================================================

// GetPreviewTemplateHandler godoc
// @Summary Get preview template by tahun dan fakultas unit
// @Description Mengambil data preview template berdasarkan tahun dan fakultas unit
// @Tags PreviewTemplate
// @Param tipe path string true "Tipe"
// @Param tahun path string true "Tahun Renstra"
// @Param fakultasUnit path string true "Fakultas Unit ID / UUID"
// @Produce json
// @Success 200 {array} previewtemplatedomain.PreviewTemplate
// @Failure 400 {object} commoninfra.ResponseError
// @Failure 404 {object} commoninfra.ResponseError
// @Failure 409 {object} commoninfra.ResponseError
// @Failure 500 {object} commoninfra.ResponseError
// @Router /preview/audit/{tahun}/{fakultasUnit} [get]
func GetPreviewTemplateHandler(c *fiber.Ctx) error {
	tipe := c.Params("tipe")
	tahun := c.Params("tahun")
	fakultasUnit := c.Params("fakultasUnit")

	query := GetPreviewTemplate.GetPreviewTemplateByTahunFakultasUnitQuery{
		Tipe:         tipe,
		Tahun:        tahun,
		FakultasUnit: fakultasUnit,
	}

	preview, err := mediatr.Send[
		GetPreviewTemplate.GetPreviewTemplateByTahunFakultasUnitQuery,
		[]previewtemplatedomain.PreviewTemplate,
	](context.Background(), query)

	if err != nil {
		return commoninfra.HandleError(c, err)
	}

	return c.JSON(preview)
}

func ModulePreviewTemplate(app *fiber.App) {
	admin := []string{"admin", "auditee", "auditor1", "auditor2"}
	whoamiURL := "http://localhost:3000/whoami"

	app.Get("/preview/audit/:tipe/:tahun/:fakultasUnit", commonpresentation.SmartCompress(), commonpresentation.JWTMiddleware(), commonpresentation.RBACMiddleware(admin, whoamiURL), GetPreviewTemplateHandler)
}
