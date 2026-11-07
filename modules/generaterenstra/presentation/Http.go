package presentation

import (
	"context"
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/mehdihadeli/go-mediatr"
	"gorm.io/gorm"

	commoninfra "UnpakSiamida/common/infrastructure"
	commonpresentation "UnpakSiamida/common/presentation"

	DeleteGenerateRenstra "UnpakSiamida/modules/generaterenstra/application/DeleteGenerateRenstra"
	GenerateRenstra "UnpakSiamida/modules/generaterenstra/application/GenerateRenstra"
)

var generateDB *gorm.DB

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
	if generateDB == nil {
		return c.Status(500).JSON(fiber.Map{"Code": "INTERNAL_SERVER_ERROR", "Message": "Database not initialized"})
	}

	uuids, _ := parseGenerateUUIDs(c)
	if len(uuids) == 0 {
		return c.Status(400).JSON(fiber.Map{
			"Code":    "GenerateRenstra.EmptySelection",
			"Message": "Pilih minimal 1 jadwal audit untuk melihat preview generate data.",
		})
	}

	var items []SchedulePreviewComparison
	var totalBeforeRows, totalAfterRows, totalInserted, totalDeleted, totalMerged, totalUpdatedTugas int

	for _, u := range uuids {
		comp, err := BuildSchedulePreview(generateDB, u)
		if err != nil {
			return c.Status(400).JSON(fiber.Map{
				"Code":    "GenerateRenstra.PreviewFailed",
				"Message": err.Error(),
			})
		}
		items = append(items, *comp)
		totalBeforeRows += comp.Summary.BeforeTotalRows
		totalAfterRows += comp.Summary.AfterTotalRows
		totalInserted += comp.Summary.InsertedCount
		totalDeleted += comp.Summary.DeletedCount
		totalMerged += comp.Summary.MergedGroupCount
		totalUpdatedTugas += comp.Summary.UpdatedTugasCount
	}

	return c.JSON(fiber.Map{
		"items":                items,
		"totalSchedules":       len(items),
		"totalBeforeRows":      totalBeforeRows,
		"totalAfterRows":       totalAfterRows,
		"totalInserted":        totalInserted,
		"totalDeleted":         totalDeleted,
		"totalMerged":          totalMerged,
		"totalUpdatedTugas":    totalUpdatedTugas,
		"requiredConfirmation": RequiredGenerateConfirmationPhrase,
	})
}

// POST /renstra/generate/execute
func ExecuteGenerateRenstraHandlerfunc(c *fiber.Ctx) error {
	if generateDB == nil {
		return c.Status(500).JSON(fiber.Map{"Code": "INTERNAL_SERVER_ERROR", "Message": "Database not initialized"})
	}

	uuids, confirmation := parseGenerateUUIDs(c)
	if len(uuids) == 0 {
		return c.Status(400).JSON(fiber.Map{
			"Code":    "GenerateRenstra.EmptySelection",
			"Message": "Pilih minimal 1 jadwal audit untuk dieksekusi.",
		})
	}

	if confirmation != RequiredGenerateConfirmationPhrase {
		return c.Status(400).JSON(fiber.Map{
			"Code":    "GenerateRenstra.InvalidConfirmation",
			"Message": fmt.Sprintf("Kalimat konfirmasi tidak sesuai. Ketik tepat: \"%s\"", RequiredGenerateConfirmationPhrase),
		})
	}

	tx := generateDB.Begin()
	if tx.Error != nil {
		return c.Status(500).JSON(fiber.Map{"Code": "INTERNAL_SERVER_ERROR", "Message": tx.Error.Error()})
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var processedTargets []string
	var processedUUIDs []string

	for _, u := range uuids {
		rec, err := ExecuteLegacySyncRenstra(tx, u)
		if err != nil {
			tx.Rollback()
			return c.Status(500).JSON(fiber.Map{
				"Code":    "GenerateRenstra.ExecuteFailed",
				"Message": err.Error(),
			})
		}
		processedUUIDs = append(processedUUIDs, rec.UUID)
		processedTargets = append(processedTargets, fmt.Sprintf("%s (%s)", rec.NamaFakProdUnit, rec.Tahun))
	}

	if err := tx.Commit().Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"Code": "INTERNAL_SERVER_ERROR", "Message": err.Error()})
	}

	return c.JSON(fiber.Map{
		"message":          fmt.Sprintf("Berhasil generate data pada %d jadwal audit.", len(processedUUIDs)),
		"processedCount":   len(processedUUIDs),
		"processedUuids":   processedUUIDs,
		"processedTargets": processedTargets,
	})
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
	renstraUUID := strings.TrimSpace(c.FormValue("renstra"))
	if renstraUUID != "" && generateDB != nil {
		tx := generateDB.Begin()
		if tx.Error != nil {
			return c.Status(500).JSON(fiber.Map{"Code": "INTERNAL_SERVER_ERROR", "Message": tx.Error.Error()})
		}
		rec, err := ExecuteLegacySyncRenstra(tx, renstraUUID)
		if err != nil {
			tx.Rollback()
			return c.Status(400).JSON(fiber.Map{"Code": "GenerateRenstra.Failed", "Message": err.Error()})
		}
		if err := tx.Commit().Error; err != nil {
			return c.Status(500).JSON(fiber.Map{"Code": "INTERNAL_SERVER_ERROR", "Message": err.Error()})
		}
		return commonpresentation.JsonUUID(c, rec.UUID)
	}

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
	if len(dbs) > 0 && dbs[0] != nil {
		generateDB = dbs[0]
	}

	admin := []string{"admin"}
	auditRoles := []string{"admin", "auditee", "auditor1", "auditor2"}
	whoamiURL := "http://localhost:3000/whoami"

	app.Get("/renstra/generate/preview", commonpresentation.SmartCompress(), commonpresentation.JWTMiddleware(), commonpresentation.RBACMiddleware(auditRoles, whoamiURL), PreviewGenerateRenstraHandlerfunc)
	app.Post("/renstra/generate/preview", commonpresentation.SmartCompress(), commonpresentation.JWTMiddleware(), commonpresentation.RBACMiddleware(auditRoles, whoamiURL), PreviewGenerateRenstraHandlerfunc)
	app.Post("/renstra/generate/execute", commonpresentation.JWTMiddleware(), commonpresentation.RBACMiddleware(admin, whoamiURL), ExecuteGenerateRenstraHandlerfunc)

	app.Post("/renstra/generate", commonpresentation.JWTMiddleware(), commonpresentation.RBACMiddleware(auditRoles, whoamiURL), GenerateRenstraHandlerfunc)
	app.Delete("/renstra/generate/question", commonpresentation.JWTMiddleware(), commonpresentation.RBACMiddleware(admin, whoamiURL), DeleteRenstraQuestionHandlerfunc)
}
