package presentation

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"github.com/mehdihadeli/go-mediatr"

	commoninfra "UnpakSiamida/common/infrastructure"
	// commondomain "UnpakSiamida/common/domain"

	previewtemplatedomain "UnpakSiamida/modules/previewtemplate/domain"
	GetPreviewTemplate "UnpakSiamida/modules/previewtemplate/application/GetPreviewTemplate"
)

// =======================================================
// GET /preview/audit/{tahun}/{fakultasUnit}
// =======================================================

// GetPreviewTemplateHandler godoc
// @Summary Get preview template by tahun dan fakultas unit
// @Description Mengambil data preview template berdasarkan tahun dan fakultas unit
// @Tags PreviewTemplate
// @Param tahun path string true "Tahun Renstra"
// @Param fakultasUnit path string true "Fakultas Unit ID / UUID"
// @Produce json
// @Success 200 {array} previewtemplatedomain.PreviewTemplate
// @Failure 400 {object} commoninfra.ErrorResponse
// @Failure 404 {object} commoninfra.ErrorResponse
// @Failure 500 {object} commoninfra.ErrorResponse
// @Router /preview/audit/{tahun}/{fakultasUnit} [get]
func GetPreviewTemplateHandler(c *fiber.Ctx) error {
	tahun := c.Params("tahun")
	fakultasUnit := c.Params("fakultasUnit")

	query := GetPreviewTemplate.GetPreviewTemplateByTahunFakultasUnitQuery{
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
	app.Get("/preview/audit/:tahun/:fakultasUnit", GetPreviewTemplateHandler)
}