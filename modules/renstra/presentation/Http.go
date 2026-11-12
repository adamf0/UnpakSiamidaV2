package presentation

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/mehdihadeli/go-mediatr"
	"gorm.io/gorm"

	commondomain "UnpakSiamida/common/domain"
	commoninfra "UnpakSiamida/common/infrastructure"
	commonpresentation "UnpakSiamida/common/presentation"
	renstradomain "UnpakSiamida/modules/renstra/domain"
	eventRenstra "UnpakSiamida/modules/renstra/event"

	CreateRenstra "UnpakSiamida/modules/renstra/application/CreateRenstra"
	DeleteRenstra "UnpakSiamida/modules/renstra/application/DeleteRenstra"
	GetAllRenstras "UnpakSiamida/modules/renstra/application/GetAllRenstras"
	GetRenstraDefault "UnpakSiamida/modules/renstra/application/GetRenstraDefault"
	GiveCodeRenstra "UnpakSiamida/modules/renstra/application/GiveCodeAccessRenstra"
	SetupUuidRenstra "UnpakSiamida/modules/renstra/application/SetupUuidRenstra"
	UpdateRenstra "UnpakSiamida/modules/renstra/application/UpdateRenstra"
)

var renstraDB *gorm.DB


// =======================================================
// POST /renstra
// =======================================================

// CreateRenstraHandler godoc
// @Summary Create new Renstra
// @Tags Renstra
// @Param tahun formData string true "Tahun"
// @Param fakultas_unit formData string true "Fakultas/Unit ID"
// @Param periode_upload_mulai formData string true "Periode Upload Mulai"
// @Param periode_upload_akhir formData string true "Periode Upload Akhir"
// @Param periode_assesment_dokumen_mulai formData string true "Periode Assesment Dokumen Mulai"
// @Param periode_assesment_dokumen_akhir formData string true "Periode Assesment Dokumen Akhir"
// @Param periode_assesment_lapangan_mulai formData string true "Periode Assesment Lapangan Mulai"
// @Param periode_assesment_lapangan_akhir formData string true "Periode Assesment Lapangan Akhir"
// @Param auditee formData string true "Auditee"
// @Param auditor1 formData string true "Auditor 1"
// @Param auditor2 formData string true "Auditor 2"
// @Produce json
// @Success 200 {object} map[string]string "uuid of created Renstra"
// @Failure 400 {object} commoninfra.ResponseError
// @Failure 404 {object} commoninfra.ResponseError
// @Failure 409 {object} commoninfra.ResponseError
// @Failure 500 {object} commoninfra.ResponseError
// @Router /renstra [post]
func CreateRenstraHandlerfunc(c *fiber.Ctx) error {
	cmd := CreateRenstra.CreateRenstraCommand{
		Tahun:                         c.FormValue("tahun"),
		FakultasUnit:                  c.FormValue("fakultas_unit"),
		PeriodeUploadMulai:            c.FormValue("periode_upload_mulai"),
		PeriodeUploadAkhir:            c.FormValue("periode_upload_akhir"),
		PeriodeAssesmentDokumenMulai:  c.FormValue("periode_assesment_dokumen_mulai"),
		PeriodeAssesmentDokumenAkhir:  c.FormValue("periode_assesment_dokumen_akhir"),
		PeriodeAssesmentLapanganMulai: c.FormValue("periode_assesment_lapangan_mulai"),
		PeriodeAssesmentLapanganAkhir: c.FormValue("periode_assesment_lapangan_akhir"),
		Auditee:                       c.FormValue("auditee"),
		Auditor1:                      c.FormValue("auditor1"),
		Auditor2:                      c.FormValue("auditor2"),
	}

	// Kirim ke mediator
	uuid, err := mediatr.Send[
		CreateRenstra.CreateRenstraCommand,
		string,
	](context.Background(), cmd)

	if err != nil {
		return commoninfra.HandleError(c, err)
	}

	return commonpresentation.JsonUUID(c, uuid)
}

// =======================================================
// PUT /renstra/{uuid}
// =======================================================

// UpdateRenstraHandler godoc
// @Summary Update existing Renstra
// @Tags Renstra
// @Param uuid path string true "Renstra UUID" format(uuid)
// @Param tahun formData string true "Tahun"
// @Param fakultas_unit formData string true "Fakultas/Unit ID"
// @Param periode_upload_mulai formData string true "Periode Upload Mulai"
// @Param periode_upload_akhir formData string true "Periode Upload Akhir"
// @Param periode_assesment_dokumen_mulai formData string true "Periode Assesment Dokumen Mulai"
// @Param periode_assesment_dokumen_akhir formData string true "Periode Assesment Dokumen Akhir"
// @Param periode_assesment_lapangan_mulai formData string true "Periode Assesment Lapangan Mulai"
// @Param periode_assesment_lapangan_akhir formData string true "Periode Assesment Lapangan Akhir"
// @Param auditee formData string true "Auditee"
// @Param auditor1 formData string true "Auditor 1"
// @Param auditor2 formData string true "Auditor 2"
// @Produce json
// @Success 200 {object} map[string]string "uuid of updated Renstra"
// @Failure 400 {object} commoninfra.ResponseError
// @Failure 404 {object} commoninfra.ResponseError
// @Failure 409 {object} commoninfra.ResponseError
// @Failure 500 {object} commoninfra.ResponseError
// @Router /renstra/{uuid} [put]
func UpdateRenstraHandlerfunc(c *fiber.Ctx) error {
	cmd := UpdateRenstra.UpdateRenstraCommand{
		Uuid:                          c.Params("uuid"),
		Tahun:                         c.FormValue("tahun"),
		FakultasUnit:                  c.FormValue("fakultas_unit"),
		PeriodeUploadMulai:            c.FormValue("periode_upload_mulai"),
		PeriodeUploadAkhir:            c.FormValue("periode_upload_akhir"),
		PeriodeAssesmentDokumenMulai:  c.FormValue("periode_assesment_dokumen_mulai"),
		PeriodeAssesmentDokumenAkhir:  c.FormValue("periode_assesment_dokumen_akhir"),
		PeriodeAssesmentLapanganMulai: c.FormValue("periode_assesment_lapangan_mulai"),
		PeriodeAssesmentLapanganAkhir: c.FormValue("periode_assesment_lapangan_akhir"),
		Auditee:                       c.FormValue("auditee"),
		Auditor1:                      c.FormValue("auditor1"),
		Auditor2:                      c.FormValue("auditor2"),
	}

	updatedID, err := mediatr.Send[
		UpdateRenstra.UpdateRenstraCommand,
		string,
	](context.Background(), cmd)

	if err != nil {
		return commoninfra.HandleError(c, err)
	}

	return c.JSON(fiber.Map{"uuid": updatedID})
}

// =======================================================
// PUT /renstra/{uuid}/code_access
// =======================================================

// GiveCodeAccessRenstraHandler godoc
// @Summary Update Renstra Kode Akses
// @Tags Renstra
// @Param uuid path string true "Renstra UUID" format(uuid)
// @Param kodeAkses formData string true "Kode Akses baru"
// @Produce json
// @Success 200 {object} map[string]string "uuid of updated Renstra"
// @Failure 400 {object} commoninfra.ResponseError
// @Failure 404 {object} commoninfra.ResponseError
// @Failure 409 {object} commoninfra.ResponseError
// @Failure 500 {object} commoninfra.ResponseError
// @Router /renstra/{uuid}/code_access [put]
func GiveCodeAccessRenstraHandlerfunc(c *fiber.Ctx) error {
	cmd := GiveCodeRenstra.GiveCodeAccessRenstraCommand{
		Uuid:      c.Params("uuid"),
		KodeAkses: c.FormValue("kodeAkses"),
	}

	updatedID, err := mediatr.Send[
		GiveCodeRenstra.GiveCodeAccessRenstraCommand,
		string,
	](context.Background(), cmd)

	if err != nil {
		return commoninfra.HandleError(c, err)
	}

	return c.JSON(fiber.Map{"uuid": updatedID})
}

// =======================================================
// DELETE /renstra/{uuid}
// =======================================================

// DeleteRenstraHandler godoc
// @Summary Delete an Renstra
// @Tags Renstra
// @Param uuid path string true "Renstra UUID" format(uuid)
// @Produce json
// @Success 200 {object} map[string]string "uuid of deleted Renstra"
// @Failure 400 {object} commoninfra.ResponseError
// @Failure 404 {object} commoninfra.ResponseError
// @Failure 409 {object} commoninfra.ResponseError
// @Failure 500 {object} commoninfra.ResponseError
//
//	@Router /renstra/{uuid} [delete]
func DeleteRenstraHandlerfunc(c *fiber.Ctx) error {
	uuid := c.Params("uuid")

	cmd := DeleteRenstra.DeleteRenstraCommand{Uuid: uuid}

	deletedID, err := mediatr.Send[
		DeleteRenstra.DeleteRenstraCommand,
		string,
	](context.Background(), cmd)

	if err != nil {
		return commoninfra.HandleError(c, err)
	}

	return c.JSON(fiber.Map{"uuid": deletedID})
}

// =======================================================
// GET /renstra/{uuid}
// =======================================================

// GetRenstraHandler godoc
// @Summary Get Renstra by UUID
// @Tags Renstra
// @Param uuid path string true "Renstra UUID" format(uuid)
// @Produce json
// @Success 200 {object} renstradomain.RenstraDefault
// @Failure 400 {object} commoninfra.ResponseError
// @Failure 404 {object} commoninfra.ResponseError
// @Failure 409 {object} commoninfra.ResponseError
// @Failure 500 {object} commoninfra.ResponseError
//
//	@Router /renstra/{uuid} [get]
func GetRenstraHandlerfunc(c *fiber.Ctx) error {
	uuid := c.Params("uuid")

	query := GetRenstraDefault.GetRenstraDefaultByUuidQuery{Uuid: uuid}

	renstra, err := mediatr.Send[
		GetRenstraDefault.GetRenstraDefaultByUuidQuery,
		*renstradomain.RenstraDefault,
	](context.Background(), query)

	if err != nil {
		return commoninfra.HandleError(c, err)
	}

	if renstra == nil {
		return c.Status(404).JSON(fiber.Map{"error": "Renstra not found"})
	}

	return c.JSON(renstra)
}

// =======================================================
// GET /renstra/{tahun}
// =======================================================

// GetRenstraHandler godoc
// @Summary Get Renstra by Tahun
// @Tags Renstra
// @Param tahun path string true "Tahun"
// @Produce json
// @Success 200 {object} commondomain.Paged[renstradomain.RenstraDefault]
// @Failure 400 {object} commoninfra.ResponseError
// @Failure 404 {object} commoninfra.ResponseError
// @Failure 409 {object} commoninfra.ResponseError
// @Failure 500 {object} commoninfra.ResponseError
//
//	@Router /renstra/{tahun} [get]
func GetRenstraByTahunTargetHandlerfunc(c *fiber.Ctx) error {
	mode := c.Query("mode", "paging")
	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 10)
	search := c.Query("search", "")

	filtersRaw := c.Query("filters", "")
	tahun := c.Params("tahun")
	if tahun != "" {
		if filtersRaw != "" {
			filtersRaw += ";"
		}
		filtersRaw += fmt.Sprintf("tahun:eq:%s", tahun)
	}

	sid := c.FormValue("sid")
	if sid != "" {
		if filtersRaw != "" {
			filtersRaw += ";"
		}
		filtersRaw += fmt.Sprintf("uuidauditee:eq:%s;uuidauditor1:eq:%s;uuidauditor2:eq:%s", sid, sid, sid)
	}

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

	query := GetAllRenstras.GetAllRenstrasQuery{
		Scope:         "audit",
		Search:        search,
		SearchFilters: filters,
	}

	var adapter commonpresentation.OutputAdapter[renstradomain.RenstraDefault]
	switch mode {
	case "all":
		adapter = &commonpresentation.AllAdapter[renstradomain.RenstraDefault]{}
	case "ndjson":
		adapter = &commonpresentation.NDJSONAdapter[renstradomain.RenstraDefault]{}
	case "sse":
		adapter = &commonpresentation.SSEAdapter[renstradomain.RenstraDefault]{}
	default:
		query.Page = &page
		query.Limit = &limit
		adapter = &commonpresentation.PagingAdapter[renstradomain.RenstraDefault]{}
	}

	result, err := mediatr.Send[
		GetAllRenstras.GetAllRenstrasQuery,
		commondomain.Paged[renstradomain.RenstraDefault],
	](context.Background(), query)

	if err != nil {
		return commoninfra.HandleError(c, err)
	}

	return adapter.Send(c, result)
}

// GetAllRenstrasHandler godoc
// @Summary Get All Renstras
// @Tags Renstra
// @Param mode query string false "paging | all | ndjson | sse"
// @Param page query int false "Page number"
// @Param limit query int false "Limit per page"
// @Param search query string false "Search keyword"
// @Produce json
// @Success 200 {object} commondomain.Paged[renstradomain.RenstraDefault]
// @Router /renstras [get]
func GetAllRenstrasHandlerfunc(c *fiber.Ctx) error {
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

	query := GetAllRenstras.GetAllRenstrasQuery{
		Search:        search,
		SearchFilters: filters,
	}

	var adapter commonpresentation.OutputAdapter[renstradomain.RenstraDefault]
	switch mode {
	case "all":
		adapter = &commonpresentation.AllAdapter[renstradomain.RenstraDefault]{}
	case "ndjson":
		adapter = &commonpresentation.NDJSONAdapter[renstradomain.RenstraDefault]{}
	case "sse":
		adapter = &commonpresentation.SSEAdapter[renstradomain.RenstraDefault]{}
	default:
		query.Page = &page
		query.Limit = &limit
		adapter = &commonpresentation.PagingAdapter[renstradomain.RenstraDefault]{}
	}

	result, err := mediatr.Send[
		GetAllRenstras.GetAllRenstrasQuery,
		commondomain.Paged[renstradomain.RenstraDefault],
	](context.Background(), query)

	if err != nil {
		return commoninfra.HandleError(c, err)
	}

	return adapter.Send(c, result)
}

func SetupUuidRenstrasHandlerfunc(c *fiber.Ctx) error {
	cmd := SetupUuidRenstra.SetupUuidRenstraCommand{}

	message, err := mediatr.Send[SetupUuidRenstra.SetupUuidRenstraCommand, string](context.Background(), cmd)
	if err != nil {
		return commoninfra.HandleError(c, err)
	}

	return c.JSON(fiber.Map{"message": message})
}

func parseUUIDListFromRequest(c *fiber.Ctx) []string {
	var list []string
	if paramUUID := strings.TrimSpace(c.Params("uuid")); paramUUID != "" {
		list = append(list, paramUUID)
	}

	if formUUIDs := strings.TrimSpace(c.FormValue("uuids")); formUUIDs != "" {
		for _, part := range strings.Split(formUUIDs, ",") {
			if trimmed := strings.TrimSpace(part); trimmed != "" {
				list = append(list, trimmed)
			}
		}
	}
	if formUUID := strings.TrimSpace(c.FormValue("uuid")); formUUID != "" {
		list = append(list, formUUID)
	}

	seen := make(map[string]bool)
	var unique []string
	for _, item := range list {
		trimmed := strings.TrimSpace(item)
		if trimmed == "" || seen[trimmed] {
			continue
		}
		seen[trimmed] = true
		unique = append(unique, trimmed)
	}
	return unique
}

// POST /renstra/send-email and POST /renstra/:uuid/send-email
// Enqueues RenstraSendEmailEvent into outbox_messages table (Outbox Pattern)
func SendEmailRenstraHandlerfunc(c *fiber.Ctx) error {
	if renstraDB == nil {
		return c.Status(500).JSON(fiber.Map{"Code": "INTERNAL_SERVER_ERROR", "Message": "Database not initialized"})
	}

	uuids := parseUUIDListFromRequest(c)
	if len(uuids) == 0 {
		return c.Status(400).JSON(fiber.Map{
			"Code":    " Renstra.EmptySelection",
			"Message": "Pilih minimal 1 jadwal audit untuk dikirim email kode akses.",
		})
	}

	type renstraEmailRow struct {
		UUID         string  `gorm:"column:uuid"`
		Tahun        string  `gorm:"column:tahun"`
		KodeAkses    *string `gorm:"column:kodeAkses"`
		TargetAudit  string  `gorm:"column:target_audit"`
		AuditeeName  *string `gorm:"column:auditee_name"`
		AuditeeEmail *string `gorm:"column:auditee_email"`
		Auditor1Name *string `gorm:"column:auditor1_name"`
		Auditor1Mail *string `gorm:"column:auditor1_email"`
		Auditor2Name *string `gorm:"column:auditor2_name"`
		Auditor2Mail *string `gorm:"column:auditor2_email"`
	}

	var rows []renstraEmailRow
	err := renstraDB.Raw(`
		SELECT 
			r.uuid AS uuid,
			r.tahun AS tahun,
			r.kodeAkses AS kodeAkses,
			COALESCE(fu.nama_fak_prod_unit, '-') AS target_audit,
			u1.name AS auditee_name,
			u1.email AS auditee_email,
			u2.name AS auditor1_name,
			u2.email AS auditor1_email,
			u3.name AS auditor2_name,
			u3.email AS auditor2_email
		FROM renstra r
		LEFT JOIN v_fakultas_unit fu ON r.fakultas_unit = fu.id
		LEFT JOIN usersv3 u1 ON r.auditee = u1.id
		LEFT JOIN usersv3 u2 ON r.auditor1 = u2.id
		LEFT JOIN usersv3 u3 ON r.auditor2 = u3.id
		WHERE r.uuid IN ?
	`, uuids).Scan(&rows).Error
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"Code": "INTERNAL_SERVER_ERROR", "Message": err.Error()})
	}

	if len(rows) == 0 {
		return c.Status(404).JSON(fiber.Map{"Code": "Renstra.NotFound", "Message": "Jadwal audit tidak ditemukan."})
	}

	var queuedUUIDs []string
	var outboxIDs []string
	var skippedTargets []string

	tx := renstraDB.Begin()
	if tx.Error != nil {
		return c.Status(500).JSON(fiber.Map{"Code": "INTERNAL_SERVER_ERROR", "Message": tx.Error.Error()})
	}

	for _, row := range rows {
		kode := ""
		if row.KodeAkses != nil {
			kode = strings.TrimSpace(*row.KodeAkses)
		}
		if kode == "" {
			skippedTargets = append(skippedTargets, fmt.Sprintf("%s (%s)", row.TargetAudit, row.Tahun))
			continue
		}

		var recipients []eventRenstra.RenstraEmailRecipient
		addRec := func(role string, namePtr, emailPtr *string) {
			name := ""
			if namePtr != nil {
				name = strings.TrimSpace(*namePtr)
			}
			email := ""
			if emailPtr != nil {
				email = strings.TrimSpace(*emailPtr)
			}
			if name == "" && email == "" {
				return
			}
			recipients = append(recipients, eventRenstra.RenstraEmailRecipient{
				Role:  role,
				Name:  name,
				Email: email,
			})
		}

		addRec("Auditee", row.AuditeeName, row.AuditeeEmail)
		addRec("Auditor 1", row.Auditor1Name, row.Auditor1Mail)
		addRec("Auditor 2", row.Auditor2Name, row.Auditor2Mail)

		parsedRenstraUUID, _ := uuid.Parse(row.UUID)
		evt := eventRenstra.RenstraSendEmailEvent{
			EventID:     uuid.New(),
			RenstraUUID: parsedRenstraUUID,
			Tahun:       row.Tahun,
			TargetAudit: row.TargetAudit,
			KodeAkses:   kode,
			Recipients:  recipients,
			OccurredOn:  time.Now().UTC(),
		}

		payload, err := json.Marshal(evt)
		if err != nil {
			tx.Rollback()
			return c.Status(500).JSON(fiber.Map{"Code": "INTERNAL_SERVER_ERROR", "Message": err.Error()})
		}

		outbox := commoninfra.OutboxMessage{
			ID:            evt.ID(),
			Type:          commoninfra.CanonicalTypeName(evt),
			Payload:       string(payload),
			OccurredOnUTC: evt.OccurredOnUTC(),
		}

		if err := tx.Create(&outbox).Error; err != nil {
			tx.Rollback()
			return c.Status(500).JSON(fiber.Map{"Code": "INTERNAL_SERVER_ERROR", "Message": err.Error()})
		}

		queuedUUIDs = append(queuedUUIDs, row.UUID)
		outboxIDs = append(outboxIDs, outbox.ID)
	}

	if len(queuedUUIDs) == 0 {
		tx.Rollback()
		return c.Status(400).JSON(fiber.Map{
			"Code":           "Renstra.MissingKodeAkses",
			"Message":        "Gagal mengirim email: tidak ada kode akses pada jadwal yang dipilih. Silakan atur kode akses terlebih dahulu.",
			"skippedTargets": skippedTargets,
		})
	}

	if err := tx.Commit().Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"Code": "INTERNAL_SERVER_ERROR", "Message": err.Error()})
	}

	return c.JSON(fiber.Map{
		"message":        fmt.Sprintf("Berhasil memasukkan %d pengiriman email ke antrean background job (Outbox).", len(queuedUUIDs)),
		"queuedCount":    len(queuedUUIDs),
		"queuedUuids":    queuedUUIDs,
		"outboxIds":      outboxIDs,
		"skippedCount":   len(skippedTargets),
		"skippedTargets": skippedTargets,
	})
}

// POST /renstra/bulk-delete
func BulkDeleteRenstraHandlerfunc(c *fiber.Ctx) error {
	uuids := parseUUIDListFromRequest(c)
	if len(uuids) == 0 {
		return c.Status(400).JSON(fiber.Map{
			"Code":    "Renstra.EmptySelection",
			"Message": "Pilih minimal 1 jadwal audit untuk dihapus.",
		})
	}

	var deletedUUIDs []string
	for _, u := range uuids {
		cmd := DeleteRenstra.DeleteRenstraCommand{Uuid: u}
		deletedID, err := mediatr.Send[
			DeleteRenstra.DeleteRenstraCommand,
			string,
		](context.Background(), cmd)
		if err != nil {
			return commoninfra.HandleError(c, err)
		}
		deletedUUIDs = append(deletedUUIDs, deletedID)
	}

	return c.JSON(fiber.Map{
		"message":      fmt.Sprintf("Berhasil menghapus %d jadwal audit.", len(deletedUUIDs)),
		"deletedCount": len(deletedUUIDs),
		"deletedUuids": deletedUUIDs,
	})
}

func ModuleRenstra(app *fiber.App, dbs ...*gorm.DB) {
	if len(dbs) > 0 && dbs[0] != nil {
		renstraDB = dbs[0]
	}

	admin := []string{"admin"}
	audit := []string{"auditee", "auditor1", "auditor2"}
	whoamiURL := "http://localhost:3000/whoami"

	app.Get("/renstra/setupuuid", SetupUuidRenstrasHandlerfunc)

	app.Post("/renstra/send-email", commonpresentation.JWTMiddleware(), commonpresentation.RBACMiddleware(admin, whoamiURL), SendEmailRenstraHandlerfunc)
	app.Post("/renstra/:uuid/send-email", commonpresentation.JWTMiddleware(), commonpresentation.RBACMiddleware(admin, whoamiURL), SendEmailRenstraHandlerfunc)
	app.Post("/renstra/bulk-delete", commonpresentation.JWTMiddleware(), commonpresentation.RBACMiddleware(admin, whoamiURL), BulkDeleteRenstraHandlerfunc)

	app.Post("/renstra", commonpresentation.JWTMiddleware(), commonpresentation.RBACMiddleware(admin, whoamiURL), CreateRenstraHandlerfunc)
	app.Put("/renstra/:uuid", commonpresentation.JWTMiddleware(), commonpresentation.RBACMiddleware(admin, whoamiURL), UpdateRenstraHandlerfunc)
	app.Put("/renstra/:uuid/code_access", commonpresentation.JWTMiddleware(), commonpresentation.RBACMiddleware(admin, whoamiURL), GiveCodeAccessRenstraHandlerfunc)
	app.Delete("/renstra/:uuid", commonpresentation.JWTMiddleware(), commonpresentation.RBACMiddleware(admin, whoamiURL), DeleteRenstraHandlerfunc)
	app.Get("/renstra/:uuid", commonpresentation.SmartCompress(), commonpresentation.JWTMiddleware(), GetRenstraHandlerfunc)
	app.Get("/renstras", commonpresentation.SmartCompress(), commonpresentation.JWTMiddleware(), GetAllRenstrasHandlerfunc)

	//untuk audit
	app.Get("/renstra/audit/:tahun", commonpresentation.SmartCompress(), commonpresentation.JWTMiddleware(), commonpresentation.RBACMiddleware(audit, whoamiURL), GetRenstraByTahunTargetHandlerfunc)
}
