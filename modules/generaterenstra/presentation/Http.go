package presentation

import (
	"context"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/mehdihadeli/go-mediatr"
	"gorm.io/gorm"

	commoninfra "UnpakSiamida/common/infrastructure"
	commonpresentation "UnpakSiamida/common/presentation"

	DeleteGenerateRenstra "UnpakSiamida/modules/generaterenstra/application/DeleteGenerateRenstra"
	ExecuteScheduleSync "UnpakSiamida/modules/generaterenstra/application/ExecuteScheduleSync"
	GenerateRenstra "UnpakSiamida/modules/generaterenstra/application/GenerateRenstra"
	PreviewScheduleSync "UnpakSiamida/modules/generaterenstra/application/PreviewScheduleSync"
)

func parseGenerateUUIDs(c *fiber.Ctx) ([]string, string) {
	var rawList []string

	for _, key := range []string{"uuids", "uuid", "renstra"} {
		if qv := strings.TrimSpace(c.Query(key)); qv != "" {
			for _, part := range strings.Split(qv, ",") {
				if t := strings.TrimSpace(part); t != "" {
					rawList = append(rawList, t)
				}
			}
		}
		if fv := strings.TrimSpace(c.FormValue(key)); fv != "" {
			for _, part := range strings.Split(fv, ",") {
				if t := strings.TrimSpace(part); t != "" {
					rawList = append(rawList, t)
				}
			}
		}
	}

	confirmation := strings.TrimSpace(c.FormValue("confirmation"))

	seen := make(map[string]bool)
	var unique []string
	for _, u := range rawList {
		t := strings.TrimSpace(u)
		if t == "" || seen[t] {
			continue
		}
		seen[t] = true
		unique = append(unique, t)
	}
	return unique, confirmation
}

// POST /renstra/generate/preview & GET /renstra/generate/preview
func PreviewGenerateRenstraHandlerfunc(c *fiber.Ctx) error {
	uuids, _ := parseGenerateUUIDs(c)
	if len(uuids) == 0 {
		return c.Status(400).JSON(fiber.Map{
			"Code":    "GenerateRenstra.EmptySelection",
			"Message": "Pilih minimal 1 jadwal audit untuk melihat preview generate data.",
		})
	}

	query := PreviewScheduleSync.PreviewScheduleSyncQuery{
		UUIDs: uuids,
	}

	res, err := mediatr.Send[
		PreviewScheduleSync.PreviewScheduleSyncQuery,
		*PreviewScheduleSync.PreviewScheduleSyncResult,
	](c.Context(), query)

	if err != nil {
		return commoninfra.HandleError(c, err)
	}

	return c.JSON(res)
}

// POST /renstra/generate/execute
func ExecuteGenerateRenstraHandlerfunc(c *fiber.Ctx) error {
	uuids, confirmation := parseGenerateUUIDs(c)
	if len(uuids) == 0 {
		return c.Status(400).JSON(fiber.Map{
			"Code":    "GenerateRenstra.EmptySelection",
			"Message": "Pilih minimal 1 jadwal audit untuk dieksekusi.",
		})
	}

	cmd := ExecuteScheduleSync.ExecuteScheduleSyncCommand{
		UUIDs:        uuids,
		Confirmation: confirmation,
	}

	res, err := mediatr.Send[
		ExecuteScheduleSync.ExecuteScheduleSyncCommand,
		*ExecuteScheduleSync.ExecuteScheduleSyncResult,
	](c.Context(), cmd)

	if err != nil {
		return commoninfra.HandleError(c, err)
	}

	return c.JSON(res)
}

// GenerateRenstraHandler godoc
// @Summary  new GenerateRenstra
// @Tags GenerateRenstra
// @Param tahun formData string true "Tahun"
// @Param renstra formData string true "Renstra Uuid"
// @Param fakultas_unit formData string true "Fakultas/Unit Uuid"
// @Produce json
// @Success 200 {object} map[string]string "uuid of created GenerateRenstra"
// @Router /renstra/generate [post]
func GenerateRenstraHandlerfunc(c *fiber.Ctx) error {
	cmd := GenerateRenstra.GenerateRenstraCommand{
		Tahun:            c.FormValue("tahun"),
		UuidRenstra:      c.FormValue("renstra"),
		UuidFakultasUnit: c.FormValue("fakultas_unit"),
	}

	uuid, err := mediatr.Send[
		GenerateRenstra.GenerateRenstraCommand,
		string,
	](context.Background(), cmd)

	if err != nil {
		return commoninfra.HandleError(c, err)
	}

	return commonpresentation.JsonUUID(c, uuid)
}

func DeleteRenstraQuestionHandlerfunc(c *fiber.Ctx) error {
	cmd := DeleteGenerateRenstra.DeleteGenerateRenstraCommand{
		Uuid:        c.FormValue("uuid"),
		UuidRenstra: c.FormValue("renstra"),
		Type:        c.FormValue("type"),
	}

	uuid, err := mediatr.Send[
		DeleteGenerateRenstra.DeleteGenerateRenstraCommand,
		string,
	](context.Background(), cmd)

	if err != nil {
		return commoninfra.HandleError(c, err)
	}

	return commonpresentation.JsonUUID(c, uuid)
}

func ModuleGenerateRenstra(app *fiber.App, dbs ...*gorm.DB) {
	admin := []string{"admin"}
	auditRoles := []string{"admin", "auditee", "auditor1", "auditor2"}
	whoamiURL := "http://localhost:3000/whoami"

	app.Get("/renstra/generate/preview", commonpresentation.SmartCompress(), commonpresentation.JWTMiddleware(), commonpresentation.RBACMiddleware(auditRoles, whoamiURL), PreviewGenerateRenstraHandlerfunc)
	app.Post("/renstra/generate/preview", commonpresentation.SmartCompress(), commonpresentation.JWTMiddleware(), commonpresentation.RBACMiddleware(auditRoles, whoamiURL), PreviewGenerateRenstraHandlerfunc)
	app.Post("/renstra/generate/execute", commonpresentation.JWTMiddleware(), commonpresentation.RBACMiddleware(admin, whoamiURL), ExecuteGenerateRenstraHandlerfunc)

	app.Post("/renstra/generate", commonpresentation.JWTMiddleware(), commonpresentation.RBACMiddleware(auditRoles, whoamiURL), GenerateRenstraHandlerfunc)
	app.Delete("/renstra/generate/question", commonpresentation.JWTMiddleware(), commonpresentation.RBACMiddleware(admin, whoamiURL), DeleteRenstraQuestionHandlerfunc)
}
