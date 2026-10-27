package presentation

import (
	"context"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/mehdihadeli/go-mediatr"

	commondomain "UnpakSiamida/common/domain"
	"UnpakSiamida/common/helper"
	commoninfra "UnpakSiamida/common/infrastructure"
	commonpresentation "UnpakSiamida/common/presentation"
	CreateRisk "UnpakSiamida/modules/riskassessment/application/CreateRiskAssessment"
	DeleteRisk "UnpakSiamida/modules/riskassessment/application/DeleteRiskAssessment"
	GetAllRisk "UnpakSiamida/modules/riskassessment/application/GetAllRiskAssessments"
	GetRisk "UnpakSiamida/modules/riskassessment/application/GetRiskAssessment"
	getOptions "UnpakSiamida/modules/riskassessment/application/GetRiskOptions"
	UpdateRisk "UnpakSiamida/modules/riskassessment/application/UpdateRiskAssessment"
	domainRisk "UnpakSiamida/modules/riskassessment/domain"

	"gorm.io/gorm"
)


// =======================================================
// POST /risk_assessment
// =======================================================
func CreateRiskAssessmentHandlerfunc(c *fiber.Ctx) error {
	level, userFU := commonpresentation.GetUserInfo(c)

	var cmd CreateRisk.CreateRiskAssessmentCommand

	cmd.FakultasUuid = strings.TrimSpace(c.FormValue("fakultas_unit"))
	if cmd.FakultasUuid == "" {
		cmd.FakultasUuid = strings.TrimSpace(c.FormValue("id_fakultas_unit"))
	}
	cmd.IdFakultasUnit = helper.ParseOptionalInt(c.FormValue("id_fakultas_unit"))

	cmd.TahunUuid = strings.TrimSpace(c.FormValue("tahun_uuid"))
	if cmd.TahunUuid == "" {
		cmd.TahunUuid = strings.TrimSpace(c.FormValue("id_master_tahun"))
	}
	cmd.IdMasterTahun = helper.ParseOptionalInt(c.FormValue("id_master_tahun"))
	cmd.Tahun = strings.TrimSpace(c.FormValue("tahun"))

	if idMP := helper.ParseOptionalInt(c.FormValue("id_mata_program")); idMP != nil {
		cmd.IdMataProgram = *idMP
	}
	cmd.MataProgramUuid = strings.TrimSpace(c.FormValue("mata_program_uuid"))
	if cmd.MataProgramUuid == "" {
		cmd.MataProgramUuid = strings.TrimSpace(c.FormValue("mata_program"))
	}

	cmd.IsuInternalEksternal = strings.TrimSpace(c.FormValue("isu_internal_eksternal"))
	if cmd.IsuInternalEksternal == "" {
		cmd.IsuInternalEksternal = strings.TrimSpace(c.FormValue("penyebab"))
	}
	cmd.Penyebab = cmd.IsuInternalEksternal

	cmd.DeskripsiRisiko = strings.TrimSpace(c.FormValue("deskripsi_risiko"))
	if cmd.DeskripsiRisiko == "" {
		cmd.DeskripsiRisiko = strings.TrimSpace(c.FormValue("pernyataan_risiko"))
	}
	cmd.PernyataanRisiko = cmd.DeskripsiRisiko

	kAwal := helper.ParseOptionalInt(c.FormValue("kemungkinan_awal"))
	if kAwal == nil {
		kAwal = helper.ParseOptionalInt(c.FormValue("kemungkinan"))
	}
	if kAwal != nil && *kAwal > 0 {
		cmd.KemungkinanAwal = *kAwal
	} else {
		cmd.KemungkinanAwal = 3
	}
	cmd.Kemungkinan = &cmd.KemungkinanAwal

	dAwal := helper.ParseOptionalInt(c.FormValue("dampak_awal"))
	if dAwal == nil {
		dAwal = helper.ParseOptionalInt(c.FormValue("keparahan"))
	}
	if dAwal != nil && *dAwal > 0 {
		cmd.DampakAwal = *dAwal
	} else {
		cmd.DampakAwal = 3
	}
	cmd.Keparahan = &cmd.DampakAwal

	cmd.TingkatRisikoAwal = strings.TrimSpace(c.FormValue("tingkat_risiko_awal"))
	if cmd.TingkatRisikoAwal == "" {
		cmd.TingkatRisikoAwal = strings.TrimSpace(c.FormValue("level_risiko"))
	}
	cmd.LevelRisiko = cmd.TingkatRisikoAwal

	cmd.TindakanMitigasi = strings.TrimSpace(c.FormValue("tindakan_mitigasi"))
	if cmd.TindakanMitigasi == "" {
		cmd.TindakanMitigasi = strings.TrimSpace(c.FormValue("rencana_mitigasi"))
	}
	cmd.RencanaMitigasi = cmd.TindakanMitigasi

	cmd.Peluang = strings.TrimSpace(c.FormValue("peluang"))
	cmd.KenyataanResidu = strings.TrimSpace(c.FormValue("kenyataan_residu"))
	cmd.KemungkinanResidu = helper.ParseOptionalInt(c.FormValue("kemungkinan_residu"))
	cmd.DampakResidu = helper.ParseOptionalInt(c.FormValue("dampak_residu"))
	if val := strings.TrimSpace(c.FormValue("tingkat_risiko_residu")); val != "" {
		cmd.TingkatRisikoResidu = &val
	}
	cmd.RencanaTindakLanjut = strings.TrimSpace(c.FormValue("rencana_tindak_lanjut"))
	cmd.KodeRisiko = strings.TrimSpace(c.FormValue("kode_risiko"))
	cmd.Kategori = strings.TrimSpace(c.FormValue("kategori"))
	cmd.Dampak = strings.TrimSpace(c.FormValue("dampak"))
	cmd.Pengendalian = strings.TrimSpace(c.FormValue("pengendalian"))
	cmd.PenanggungJawab = strings.TrimSpace(c.FormValue("penanggung_jawab"))
	cmd.TargetWaktu = strings.TrimSpace(c.FormValue("target_waktu"))
	cmd.Status = strings.TrimSpace(c.FormValue("status"))

	// Scoping: Non-admin (auditee) dipaksa ke unit miliknya sendiri
	if level != "admin" {
		cmd.FakultasUuid = userFU
		cmd.IdFakultasUnit = nil
	}

	createdUUID, err := mediatr.Send[CreateRisk.CreateRiskAssessmentCommand, string](context.Background(), cmd)
	if err != nil {
		return commoninfra.HandleError(c, err)
	}

	return commonpresentation.JsonUUID(c, createdUUID)
}

// =======================================================
// PUT /risk_assessment/:uuid
// =======================================================
func UpdateRiskAssessmentHandlerfunc(c *fiber.Ctx) error {
	rawUuid := c.Params("uuid")
	parsedUuid, err := uuid.Parse(rawUuid)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid UUID"})
	}

	level, userFU := commonpresentation.GetUserInfo(c)

	// Scoping: non-admin harus memiliki item ini
	if level != "admin" {
		existing, err := mediatr.Send[GetRisk.GetRiskAssessmentByUuidQuery, *domainRisk.RiskAssessmentDefault](
			context.Background(),
			GetRisk.GetRiskAssessmentByUuidQuery{UUID: parsedUuid},
		)
		if err != nil {
			return commoninfra.HandleError(c, err)
		}
		if existing == nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Risk assessment not found"})
		}
		if !helper.MatchUnit(existing.IdFakultasUnit, existing.UuidFakultasUnit, userFU) {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Akses ditolak: bukan unit anda"})
		}
	}

	var cmd UpdateRisk.UpdateRiskAssessmentCommand
	cmd.UUID = rawUuid

	cmd.FakultasUuid = strings.TrimSpace(c.FormValue("fakultas_unit"))
	if cmd.FakultasUuid == "" {
		cmd.FakultasUuid = strings.TrimSpace(c.FormValue("id_fakultas_unit"))
	}
	cmd.IdFakultasUnit = helper.ParseOptionalInt(c.FormValue("id_fakultas_unit"))

	cmd.TahunUuid = strings.TrimSpace(c.FormValue("tahun_uuid"))
	if cmd.TahunUuid == "" {
		cmd.TahunUuid = strings.TrimSpace(c.FormValue("id_master_tahun"))
	}
	cmd.IdMasterTahun = helper.ParseOptionalInt(c.FormValue("id_master_tahun"))
	cmd.Tahun = strings.TrimSpace(c.FormValue("tahun"))

	if idMP := helper.ParseOptionalInt(c.FormValue("id_mata_program")); idMP != nil {
		cmd.IdMataProgram = *idMP
	}
	cmd.MataProgramUuid = strings.TrimSpace(c.FormValue("mata_program_uuid"))
	if cmd.MataProgramUuid == "" {
		cmd.MataProgramUuid = strings.TrimSpace(c.FormValue("mata_program"))
	}

	cmd.IsuInternalEksternal = strings.TrimSpace(c.FormValue("isu_internal_eksternal"))
	if cmd.IsuInternalEksternal == "" {
		cmd.IsuInternalEksternal = strings.TrimSpace(c.FormValue("penyebab"))
	}
	cmd.Penyebab = cmd.IsuInternalEksternal

	cmd.DeskripsiRisiko = strings.TrimSpace(c.FormValue("deskripsi_risiko"))
	if cmd.DeskripsiRisiko == "" {
		cmd.DeskripsiRisiko = strings.TrimSpace(c.FormValue("pernyataan_risiko"))
	}
	cmd.PernyataanRisiko = cmd.DeskripsiRisiko

	kAwal := helper.ParseOptionalInt(c.FormValue("kemungkinan_awal"))
	if kAwal == nil {
		kAwal = helper.ParseOptionalInt(c.FormValue("kemungkinan"))
	}
	if kAwal != nil && *kAwal > 0 {
		cmd.KemungkinanAwal = *kAwal
	} else {
		cmd.KemungkinanAwal = 3
	}
	cmd.Kemungkinan = &cmd.KemungkinanAwal

	dAwal := helper.ParseOptionalInt(c.FormValue("dampak_awal"))
	if dAwal == nil {
		dAwal = helper.ParseOptionalInt(c.FormValue("keparahan"))
	}
	if dAwal != nil && *dAwal > 0 {
		cmd.DampakAwal = *dAwal
	} else {
		cmd.DampakAwal = 3
	}
	cmd.Keparahan = &cmd.DampakAwal

	cmd.TingkatRisikoAwal = strings.TrimSpace(c.FormValue("tingkat_risiko_awal"))
	if cmd.TingkatRisikoAwal == "" {
		cmd.TingkatRisikoAwal = strings.TrimSpace(c.FormValue("level_risiko"))
	}
	cmd.LevelRisiko = cmd.TingkatRisikoAwal

	cmd.TindakanMitigasi = strings.TrimSpace(c.FormValue("tindakan_mitigasi"))
	if cmd.TindakanMitigasi == "" {
		cmd.TindakanMitigasi = strings.TrimSpace(c.FormValue("rencana_mitigasi"))
	}
	cmd.RencanaMitigasi = cmd.TindakanMitigasi

	cmd.Peluang = strings.TrimSpace(c.FormValue("peluang"))
	cmd.KenyataanResidu = strings.TrimSpace(c.FormValue("kenyataan_residu"))
	cmd.KemungkinanResidu = helper.ParseOptionalInt(c.FormValue("kemungkinan_residu"))
	cmd.DampakResidu = helper.ParseOptionalInt(c.FormValue("dampak_residu"))
	if val := strings.TrimSpace(c.FormValue("tingkat_risiko_residu")); val != "" {
		cmd.TingkatRisikoResidu = &val
	}
	cmd.RencanaTindakLanjut = strings.TrimSpace(c.FormValue("rencana_tindak_lanjut"))
	cmd.KodeRisiko = strings.TrimSpace(c.FormValue("kode_risiko"))
	cmd.Kategori = strings.TrimSpace(c.FormValue("kategori"))
	cmd.Dampak = strings.TrimSpace(c.FormValue("dampak"))
	cmd.Pengendalian = strings.TrimSpace(c.FormValue("pengendalian"))
	cmd.PenanggungJawab = strings.TrimSpace(c.FormValue("penanggung_jawab"))
	cmd.TargetWaktu = strings.TrimSpace(c.FormValue("target_waktu"))
	cmd.Status = strings.TrimSpace(c.FormValue("status"))

	// Scoping: Non-admin tidak boleh memindahkan unit
	if level != "admin" {
		cmd.FakultasUuid = userFU
		cmd.IdFakultasUnit = nil
	}

	updatedUUID, err := mediatr.Send[UpdateRisk.UpdateRiskAssessmentCommand, string](context.Background(), cmd)
	if err != nil {
		return commoninfra.HandleError(c, err)
	}

	return commonpresentation.JsonUUID(c, updatedUUID)
}

// =======================================================
// DELETE /risk_assessment/:uuid
// =======================================================
func DeleteRiskAssessmentHandlerfunc(c *fiber.Ctx) error {
	rawUuid := c.Params("uuid")
	parsedUuid, err := uuid.Parse(rawUuid)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid UUID"})
	}

	level, userFU := commonpresentation.GetUserInfo(c)

	if level != "admin" {
		existing, err := mediatr.Send[GetRisk.GetRiskAssessmentByUuidQuery, *domainRisk.RiskAssessmentDefault](
			context.Background(),
			GetRisk.GetRiskAssessmentByUuidQuery{UUID: parsedUuid},
		)
		if err != nil {
			return commoninfra.HandleError(c, err)
		}
		if existing == nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Risk assessment not found"})
		}
		if !helper.MatchUnit(existing.IdFakultasUnit, existing.UuidFakultasUnit, userFU) {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Akses ditolak: bukan unit anda"})
		}
	}

	cmd := DeleteRisk.DeleteRiskAssessmentCommand{UUID: rawUuid}
	deletedUUID, err := mediatr.Send[DeleteRisk.DeleteRiskAssessmentCommand, string](context.Background(), cmd)
	if err != nil {
		return commoninfra.HandleError(c, err)
	}

	return commonpresentation.JsonUUID(c, deletedUUID)
}

// =======================================================
// GET /risk_assessment/:uuid
// =======================================================
func GetRiskAssessmentHandlerfunc(c *fiber.Ctx) error {
	rawUuid := c.Params("uuid")
	parsedUuid, err := uuid.Parse(rawUuid)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid UUID"})
	}

	level, userFU := commonpresentation.GetUserInfo(c)

	query := GetRisk.GetRiskAssessmentByUuidQuery{UUID: parsedUuid}
	result, err := mediatr.Send[GetRisk.GetRiskAssessmentByUuidQuery, *domainRisk.RiskAssessmentDefault](context.Background(), query)
	if err != nil {
		return commoninfra.HandleError(c, err)
	}
	if result == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Risk assessment not found"})
	}

	if level != "admin" && !helper.MatchUnit(result.IdFakultasUnit, result.UuidFakultasUnit, userFU) {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Akses ditolak: data milik unit lain"})
	}

	return c.JSON(result)
}

// =======================================================
// GET /risk_assessments
// =======================================================
func GetAllRiskAssessmentsHandlerfunc(c *fiber.Ctx) error {
	search := c.Query("search")
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	mode := c.Query("mode", "paging")

	level, userFU := commonpresentation.GetUserInfo(c)

	var filters []commondomain.SearchFilter

	var targetFU string
	if level != "admin" {
		// Non-admin (auditee/BPSI/unit/fakultas): batasi mutlak ke unit miliknya
		if userFU != "" && userFU != "0" {
			targetFU = userFU
		} else if qFU := strings.TrimSpace(c.Query("filter_fakultas_unit")); qFU != "" && qFU != "all" {
			targetFU = qFU
		} else {
			targetFU = "NONE"
		}
	} else {
		// Admin: tampilkan semua atau filter pilihan admin
		targetFU = strings.TrimSpace(c.Params("target"))
		if targetFU == "" {
			targetFU = strings.TrimSpace(c.Query("filter_fakultas_unit"))
		}
		if targetFU == "" {
			targetFU = strings.TrimSpace(c.Query("fakultas_unit"))
		}
		if targetFU == "" {
			targetFU = strings.TrimSpace(c.Query("id_fakultas_unit"))
		}
	}

	if targetFU != "" && targetFU != "all" && targetFU != "semua" {
		field := "id_fakultas_unit"
		operator := "eq"
		if strings.Contains(targetFU, "-") {
			field = "fakultasunituuid"
		} else if _, err := strconv.Atoi(targetFU); err != nil {
			field = "fakultas_unit"
			operator = "like"
		}
		filters = append(filters, commondomain.SearchFilter{
			Field:    field,
			Operator: operator,
			Value:    &targetFU,
		})
	}

	filterTahun := strings.TrimSpace(c.Query("filter_tahun"))
	if filterTahun == "" {
		filterTahun = strings.TrimSpace(c.Query("tahun"))
	}
	if filterTahun != "" && filterTahun != "all" && filterTahun != "semua" {
		field := "tahun"
		if strings.Contains(filterTahun, "-") {
			field = "tahunuuid"
		}
		filters = append(filters, commondomain.SearchFilter{
			Field:    field,
			Operator: "eq",
			Value:    &filterTahun,
		})
	}

	filterLevel := strings.TrimSpace(c.Query("filter_level_risiko"))
	if filterLevel == "" {
		filterLevel = strings.TrimSpace(c.Query("level_risiko"))
	}
	if filterLevel != "" && filterLevel != "all" && filterLevel != "semua" {
		filters = append(filters, commondomain.SearchFilter{
			Field:    "tingkat_risiko_awal",
			Operator: "eq",
			Value:    &filterLevel,
		})
	}

	filterProgram := strings.TrimSpace(c.Query("filter_program"))
	if filterProgram == "" {
		filterProgram = strings.TrimSpace(c.Query("program"))
	}
	if filterProgram != "" && filterProgram != "all" && filterProgram != "semua" {
		field := "id_mata_program"
		if strings.Contains(filterProgram, "-") {
			field = "mataprogramuuid"
		}
		filters = append(filters, commondomain.SearchFilter{
			Field:    field,
			Operator: "eq",
			Value:    &filterProgram,
		})
	}

	// Advanced SearchFilters / filters from RemoteTable
	filtersRaw := c.Query("filters")
	if filtersRaw == "" {
		filtersRaw = c.Query("searchFilters")
	}
	if filtersRaw != "" {
		parts := strings.Split(filtersRaw, ";")
		for _, p := range parts {
			tokens := strings.SplitN(p, ":", 3)
			if len(tokens) != 3 {
				continue
			}
			fld := strings.TrimSpace(tokens[0])
			op := strings.TrimSpace(tokens[1])
			val := strings.TrimSpace(tokens[2])

			var valPtr *string
			if val != "" && val != "null" {
				valPtr = &val
			}
			filters = append(filters, commondomain.SearchFilter{
				Field:    fld,
				Operator: op,
				Value:    valPtr,
			})
		}
	}

	query := GetAllRisk.GetAllRiskAssessmentsQuery{
		Search:        search,
		SearchFilters: filters,
	}

	var adapter commonpresentation.OutputAdapter[domainRisk.RiskAssessmentDefault]
	switch mode {
	case "all":
		adapter = &commonpresentation.AllAdapter[domainRisk.RiskAssessmentDefault]{}
	case "ndjson":
		adapter = &commonpresentation.NDJSONAdapter[domainRisk.RiskAssessmentDefault]{}
	case "sse":
		adapter = &commonpresentation.SSEAdapter[domainRisk.RiskAssessmentDefault]{}
	default:
		query.Page = &page
		query.Limit = &limit
		adapter = &commonpresentation.PagingAdapter[domainRisk.RiskAssessmentDefault]{}
	}

	result, err := mediatr.Send[
		GetAllRisk.GetAllRiskAssessmentsQuery,
		commondomain.Paged[domainRisk.RiskAssessmentDefault],
	](context.Background(), query)

	if err != nil {
		return commoninfra.HandleError(c, err)
	}

	return adapter.Send(c, result)
}

// =======================================================
// Master Data Handlers
// =======================================================
func GetRiskCriteriaHandler(c *fiber.Ctx) error {
	tipe := c.Query("tipe")
	items, err := mediatr.Send[getOptions.GetRiskCriteriaQuery, []domainRisk.RiskCriteriaItem](c.Context(), getOptions.GetRiskCriteriaQuery{Tipe: tipe})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(items)
}

func GetRiskLevelsHandler(c *fiber.Ctx) error {
	items, err := mediatr.Send[getOptions.GetRiskLevelsQuery, []domainRisk.RiskLevelItem](c.Context(), getOptions.GetRiskLevelsQuery{})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(items)
}

func GetRiskCategoriesHandler(c *fiber.Ctx) error {
	items, err := mediatr.Send[getOptions.GetRiskCategoriesQuery, []domainRisk.RiskCategoryItem](c.Context(), getOptions.GetRiskCategoriesQuery{})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(items)
}

func GetRiskStatusesHandler(c *fiber.Ctx) error {
	items, err := mediatr.Send[getOptions.GetRiskStatusesQuery, []domainRisk.RiskStatusItem](c.Context(), getOptions.GetRiskStatusesQuery{})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(items)
}

func GetRiskProgramsHandler(c *fiber.Ctx) error {
	level, userFU := commonpresentation.GetUserInfo(c)

	fuParam := strings.TrimSpace(c.Query("fakultas_unit"))
	if fuParam == "" {
		fuParam = strings.TrimSpace(c.Query("id_fakultas_unit"))
	}
	if fuParam == "" {
		fuParam = strings.TrimSpace(c.Query("target"))
	}
	if fuParam == "" && level != "admin" && userFU != "" {
		fuParam = userFU
	}

	items, err := mediatr.Send[getOptions.GetRiskProgramsQuery, []domainRisk.RiskMataProgramItem](c.Context(), getOptions.GetRiskProgramsQuery{
		FakultasUnit: fuParam,
		Tahun:        c.Query("tahun"),
		Search:       c.Query("search"),
	})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(items)
}

func GetRiskOptionsHandler(c *fiber.Ctx) error {
	level, userFU := commonpresentation.GetUserInfo(c)

	fuParam := strings.TrimSpace(c.Query("fakultas_unit"))
	if fuParam == "" {
		fuParam = strings.TrimSpace(c.Query("id_fakultas_unit"))
	}
	if fuParam == "" {
		fuParam = strings.TrimSpace(c.Query("target"))
	}
	if fuParam == "" && level != "admin" && userFU != "" {
		fuParam = userFU
	}

	res, err := mediatr.Send[getOptions.GetRiskOptionsQuery, *domainRisk.RiskOptionsResult](c.Context(), getOptions.GetRiskOptionsQuery{
		FakultasUnit: fuParam,
		Tahun:        c.Query("tahun"),
	})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{
		"criteria":    res.Criteria,
		"kemungkinan": res.Kemungkinan,
		"dampak":      res.Dampak,
		"levels":      res.Levels,
		"categories":  res.Categories,
		"statuses":    res.Statuses,
		"programs":    res.Programs,
	})
}

func ModuleRiskAssessment(app *fiber.App, dbs ...*gorm.DB) {
	roles := []string{"admin", "auditee", "fakultas", "unit", "upml", "upmf", "upms"}
	whoamiURL := "http://localhost:3000/whoami"

	// Dynamic Master Data Endpoints
	app.Get("/risk_options", commonpresentation.SmartCompress(), commonpresentation.JWTMiddleware(), GetRiskOptionsHandler)
	app.Get("/risk_criteria", commonpresentation.SmartCompress(), commonpresentation.JWTMiddleware(), GetRiskCriteriaHandler)
	app.Get("/risk_levels", commonpresentation.SmartCompress(), commonpresentation.JWTMiddleware(), GetRiskLevelsHandler)
	app.Get("/risk_categories", commonpresentation.SmartCompress(), commonpresentation.JWTMiddleware(), GetRiskCategoriesHandler)
	app.Get("/risk_statuses", commonpresentation.SmartCompress(), commonpresentation.JWTMiddleware(), GetRiskStatusesHandler)
	app.Get("/risk_programs", commonpresentation.SmartCompress(), commonpresentation.JWTMiddleware(), GetRiskProgramsHandler)

	app.Post("/risk_assessment", commonpresentation.JWTMiddleware(), commonpresentation.RBACMiddleware(roles, whoamiURL), CreateRiskAssessmentHandlerfunc)
	app.Put("/risk_assessment/:uuid", commonpresentation.JWTMiddleware(), commonpresentation.RBACMiddleware(roles, whoamiURL), UpdateRiskAssessmentHandlerfunc)
	app.Delete("/risk_assessment/:uuid", commonpresentation.JWTMiddleware(), commonpresentation.RBACMiddleware(roles, whoamiURL), DeleteRiskAssessmentHandlerfunc)

	app.Get("/risk_assessment/:uuid", commonpresentation.SmartCompress(), commonpresentation.JWTMiddleware(), commonpresentation.RBACMiddleware(roles, whoamiURL), GetRiskAssessmentHandlerfunc)
	app.Get("/RiskAssessment/:uuid", commonpresentation.SmartCompress(), commonpresentation.JWTMiddleware(), commonpresentation.RBACMiddleware(roles, whoamiURL), GetRiskAssessmentHandlerfunc)

	app.Get("/risk_assessments", commonpresentation.SmartCompress(), commonpresentation.JWTMiddleware(), commonpresentation.RBACMiddleware(roles, whoamiURL), GetAllRiskAssessmentsHandlerfunc)
	app.Get("/RiskAssessments", commonpresentation.SmartCompress(), commonpresentation.JWTMiddleware(), commonpresentation.RBACMiddleware(roles, whoamiURL), GetAllRiskAssessmentsHandlerfunc)

	app.Get("/risk_assessments/target/:target", commonpresentation.SmartCompress(), commonpresentation.JWTMiddleware(), commonpresentation.RBACMiddleware(roles, whoamiURL), GetAllRiskAssessmentsHandlerfunc)
	app.Get("/RiskAssessments/target/:target", commonpresentation.SmartCompress(), commonpresentation.JWTMiddleware(), commonpresentation.RBACMiddleware(roles, whoamiURL), GetAllRiskAssessmentsHandlerfunc)
}
