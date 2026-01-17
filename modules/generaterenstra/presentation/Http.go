package presentation

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/mehdihadeli/go-mediatr"

	commoninfra "UnpakSiamida/common/infrastructure"
	commonpresentation "UnpakSiamida/common/presentation"

	DeleteGenerateRenstra "UnpakSiamida/modules/generaterenstra/application/DeleteGenerateRenstra"
	GenerateRenstra "UnpakSiamida/modules/generaterenstra/application/GenerateRenstra"
)

// GenerateRenstraHandler godoc
// @Summary  new GenerateRenstra
// @Tags GenerateRenstra
// @Param tahun formData string true "Tahun"
// @Param renstra formData string true "Renstra Uuid"
// @Param fakultas_unit formData string true "Fakultas/Unit Uuid"
// @Produce json
// @Success 200 {object} map[string]string "uuid of created GenerateRenstra"
// @Failure 400 {object} commoninfra.ResponseError
// @Failure 404 {object} commoninfra.ResponseError
// @Failure 409 {object} commoninfra.ResponseError
// @Failure 500 {object} commoninfra.ResponseError
// @Router /renstra/generate [post]
func GenerateRenstraHandlerfunc(c *fiber.Ctx) error {
	cmd := GenerateRenstra.GenerateRenstraCommand{
		Tahun:            c.FormValue("tahun"),
		UuidRenstra:      c.FormValue("renstra"),
		UuidFakultasUnit: c.FormValue("fakultas_unit"),
	}

	// Kirim ke mediator
	uuid, err := mediatr.Send[
		GenerateRenstra.GenerateRenstraCommand,
		string,
	](context.Background(), cmd)

	if err != nil {
		return commoninfra.HandleError(c, err)
	}

	return c.JSON(fiber.Map{"uuid": uuid}) //[pr] dibuat struct common khusus uuid
}

// GenerateRenstraHandler godoc
// @Summary  new GenerateRenstra
// @Tags GenerateRenstra
// @Param tahun formData string true "Tahun"
// @Param renstra formData string true "Renstra Uuid"
// @Param fakultas_unit formData string true "Fakultas/Unit Uuid"
// @Produce json
// @Success 200 {object} map[string]string "uuid of created GenerateRenstra"
// @Failure 400 {object} commoninfra.ResponseError
// @Failure 404 {object} commoninfra.ResponseError
// @Failure 409 {object} commoninfra.ResponseError
// @Failure 500 {object} commoninfra.ResponseError
// @Router /renstra/generate [post]
func DeleteRenstraQuestionHandlerfunc(c *fiber.Ctx) error {
	cmd := DeleteGenerateRenstra.DeleteGenerateRenstraCommand{
		Uuid:        c.FormValue("uuid"),
		UuidRenstra: c.FormValue("renstra"),
		Type:        c.FormValue("type"),
	}

	// Kirim ke mediator
	uuid, err := mediatr.Send[
		DeleteGenerateRenstra.DeleteGenerateRenstraCommand,
		string,
	](context.Background(), cmd)

	if err != nil {
		return commoninfra.HandleError(c, err)
	}

	return c.JSON(fiber.Map{"uuid": uuid})
}

func ModuleGenerateRenstra(app *fiber.App) {
	admin := []string{"admin"}
	whoamiURL := "http://localhost:3000/whoami"

	app.Post("/renstra/generate", commonpresentation.JWTMiddleware(), commonpresentation.RBACMiddleware(admin, whoamiURL), GenerateRenstraHandlerfunc)
	app.Delete("/renstra/generate/question", commonpresentation.JWTMiddleware(), commonpresentation.RBACMiddleware(admin, whoamiURL), DeleteRenstraQuestionHandlerfunc)
}
