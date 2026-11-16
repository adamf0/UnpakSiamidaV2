package presentation

import (
	"fmt"
	"strings"
	"time"

	commonpresentation "UnpakSiamida/common/presentation"
	"UnpakSiamida/modules/audit/domain"
	"UnpakSiamida/modules/audit/infrastructure"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

var (
	auditService *infrastructure.AuditService
	ruleEngine   *infrastructure.RuleEngine
	aiGuardian   *infrastructure.AIGuardianService
	auditDb      *gorm.DB
)

func ModuleAudit(app *fiber.App, db *gorm.DB) {
	auditDb = db
	auditService = infrastructure.NewAuditService(db)
	ruleEngine = infrastructure.NewRuleEngine(db, auditService)
	aiGuardian = infrastructure.NewAIGuardianService(db, auditService, ruleEngine)

	adminRoles := []string{"admin"}
	whoamiURL := "http://localhost:3000/whoami"

	// 1. Settings Endpoints
	app.Get("/settings/audit-guardian", commonpresentation.JWTMiddleware(), GetGuardianConfigHandler)
	app.Put("/settings/audit-guardian", commonpresentation.JWTMiddleware(), commonpresentation.RBACMiddleware(adminRoles, whoamiURL), UpdateGuardianConfigHandler)

	// 2. Audit Logs Endpoints
	app.Get("/audit/logs", commonpresentation.JWTMiddleware(), GetAuditLogsHandler)
	app.Get("/audit/anomalies", commonpresentation.JWTMiddleware(), GetAuditAnomaliesHandler)

	// 3. Scan & Baseline Trigger Endpoints
	app.Post("/audit/scan", commonpresentation.JWTMiddleware(), commonpresentation.RBACMiddleware(adminRoles, whoamiURL), TriggerAuditScanHandler)
	app.Post("/audit/baseline", commonpresentation.JWTMiddleware(), commonpresentation.RBACMiddleware(adminRoles, whoamiURL), CreateBaselineHandler)
	app.Post("/audit/anomalies/:id/reconcile", commonpresentation.JWTMiddleware(), commonpresentation.RBACMiddleware(adminRoles, whoamiURL), ReconcileAnomalyHandler)
}

func GetGuardianConfigHandler(c *fiber.Ctx) error {
	cfg, err := auditService.GetGuardianConfig()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	var totalLogs int64
	var totalAnomalies int64
	var totalCheckpoints int64

	auditDb.Model(&domain.AuditLog{}).Count(&totalLogs)
	auditDb.Model(&domain.AuditAnomaly{}).Where("status = 'DETECTED'").Count(&totalAnomalies)
	auditDb.Model(&domain.AuditCheckpoint{}).Count(&totalCheckpoints)

	chainValid, brokenSeq, _, _ := auditService.VerifyHashChain()

	// Mask API key for security
	maskedKey := cfg.AIAPIKey
	if len(maskedKey) > 8 {
		maskedKey = maskedKey[:4] + "••••••••" + maskedKey[len(maskedKey)-4:]
	}

	return c.JSON(fiber.Map{
		"config": cfg,
		"masked_api_key": maskedKey,
		"stats": fiber.Map{
			"total_audit_logs":   totalLogs,
			"active_anomalies":   totalAnomalies,
			"total_checkpoints":  totalCheckpoints,
			"hash_chain_valid":   chainValid,
			"broken_sequence_no": brokenSeq,
		},
	})
}

func UpdateGuardianConfigHandler(c *fiber.Ctx) error {
	var input domain.AuditGuardianConfig
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid JSON payload"})
	}

	// Jika input API key adalah placeholder masked, pertahankan key lama
	if strings.Contains(input.AIAPIKey, "••••") {
		existing, _ := auditService.GetGuardianConfig()
		if existing != nil {
			input.AIAPIKey = existing.AIAPIKey
		}
	}

	userID := 1
	if u, ok := c.Locals("user").(*commonpresentation.Account); ok && u != nil {
		// Logged in user info
		_ = u
	}

	if err := auditService.SaveGuardianConfig(&input, userID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"message": "Pengaturan Audit & Anomaly Guardian berhasil disimpan",
		"config":  input,
	})
}

func GetAuditLogsHandler(c *fiber.Ctx) error {
	entityType := c.Query("entity_type")
	entityID := c.Query("entity_id")

	logs, err := auditService.GetAuditLogs(entityType, entityID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"data":  logs,
		"total": len(logs),
	})
}

func GetAuditAnomaliesHandler(c *fiber.Ctx) error {
	var anomalies []domain.AuditAnomaly
	status := c.Query("status", "DETECTED")

	q := auditDb.Model(&domain.AuditAnomaly{}).Order("detected_at desc")
	if status != "ALL" {
		q = q.Where("status = ?", status)
	}

	if err := q.Find(&anomalies).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"data":  anomalies,
		"total": len(anomalies),
	})
}

func TriggerAuditScanHandler(c *fiber.Ctx) error {
	cfg, err := auditService.GetGuardianConfig()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	mode := strings.ToLower(cfg.Mode)
	var findings []domain.AuditAnomaly

	switch mode {
	case "ai_guardian":
		f, err := aiGuardian.RunAIGuardianScan(c.Context(), cfg)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		findings = f

	case "hybrid":
		ruleFindings, _ := ruleEngine.RunDeterministicScan()
		aiFindings, err := aiGuardian.RunAIGuardianScan(c.Context(), cfg)
		if err == nil {
			findings = append(ruleFindings, aiFindings...)
		} else {
			findings = ruleFindings
		}

	default: // "rule_engine" (Tanpa AI)
		f, err := ruleEngine.RunDeterministicScan()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		findings = f
	}

	return c.JSON(fiber.Map{
		"message":   "Audit scan selesai dijalankan",
		"mode_used": cfg.Mode,
		"findings":  findings,
		"total":     len(findings),
	})
}

func CreateBaselineHandler(c *fiber.Ctx) error {
	actorName := "System Admin"
	if u, ok := c.Locals("user").(*commonpresentation.Account); ok && u != nil {
		actorName = u.Name
	}

	checkpoints, err := auditService.CreateBaselineCheckpoint(actorName)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"message":     fmt.Sprintf("Checkpoint baseline data aktif (%d tabel: aktivitas_v2, rab, timeline, file, matrix_proker) berhasil dibuat", len(checkpoints)),
		"checkpoints": checkpoints,
		"total":       len(checkpoints),
	})
}

func ReconcileAnomalyHandler(c *fiber.Ctx) error {
	id := c.Params("id")
	actorName := "Admin"
	if u, ok := c.Locals("user").(*commonpresentation.Account); ok && u != nil {
		actorName = u.Name
	}

	var body struct {
		Notes string `json:"notes"`
	}
	_ = c.BodyParser(&body)

	now := time.Now().UTC()
	err := auditDb.Model(&domain.AuditAnomaly{}).
		Where("id = ?", id).
		Updates(map[string]any{
			"status":           "RECONCILED",
			"resolved_at":      now,
			"resolved_by":      actorName,
			"resolution_notes": body.Notes,
		}).Error

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "Anomali berhasil diselesaikan / direkonsiliasi"})
}
