package presentation

import (
	"context"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/mehdihadeli/go-mediatr"

	commondomain "UnpakSiamida/common/domain"
	commoninfra "UnpakSiamida/common/infrastructure"
	commonpresentation "UnpakSiamida/common/presentation"
	CreateRisk "UnpakSiamida/modules/riskassessment/application/CreateRiskAssessment"
	DeleteRisk "UnpakSiamida/modules/riskassessment/application/DeleteRiskAssessment"
	GetAllRisk "UnpakSiamida/modules/riskassessment/application/GetAllRiskAssessments"
	GetRisk "UnpakSiamida/modules/riskassessment/application/GetRiskAssessment"
	UpdateRisk "UnpakSiamida/modules/riskassessment/application/UpdateRiskAssessment"
	domainRisk "UnpakSiamida/modules/riskassessment/domain"
	"gorm.io/gorm"
)

func parseOptionalInt(val string) *int {
	val = strings.TrimSpace(val)
	if val == "" {
		return nil
	}
	if i, err := strconv.Atoi(val); err == nil {
		return &i
	}
	return nil
}

func getUserInfo(c *fiber.Ctx) (level string, fakultasUnit string) {
	if lvl, ok := c.Locals("level").(string); ok && lvl != "" {
		level = strings.ToLower(strings.TrimSpace(lvl))
	}
	if fu, ok := c.Locals("fakultas_unit").(string); ok && fu != "" {
		fakultasUnit = strings.TrimSpace(fu)
	}

	if level == "" {
		level = strings.ToLower(strings.TrimSpace(c.FormValue("level")))
		if level == "" {
			level = strings.ToLower(strings.TrimSpace(c.Query("level")))
		}
	}

	if fakultasUnit == "" {
		fakultasUnit = strings.TrimSpace(c.FormValue("fakultas_unit"))
		if fakultasUnit == "" {
			fakultasUnit = strings.TrimSpace(c.FormValue("target"))
		}
		if fakultasUnit == "" {
			fakultasUnit = strings.TrimSpace(c.Query("filter_fakultas_unit"))
		}
		if fakultasUnit == "" {
			fakultasUnit = strings.TrimSpace(c.Query("fakultas_unit"))
		}
		if fakultasUnit == "" {
			fakultasUnit = strings.TrimSpace(c.Query("target"))
		}
	}

	return level, fakultasUnit
}

func matchUnit(itemUnitID *int, itemUnitUUID *string, userFU string) bool {
	if userFU == "" {
		return false
	}
	if itemUnitUUID != nil && strings.EqualFold(strings.TrimSpace(*itemUnitUUID), strings.TrimSpace(userFU)) {
		return true
	}
	if itemUnitID != nil && strconv.Itoa(*itemUnitID) == strings.TrimSpace(userFU) {
		return true
	}
	return false
}

// =======================================================
// POST /risk_assessment
// =======================================================
func CreateRiskAssessmentHandlerfunc(c *fiber.Ctx) error {
	level, userFU := getUserInfo(c)

	var req struct {
		FakultasUuid   string `json:"fakultas_unit"`
		IdFakultasUnit *int   `json:"id_fakultas_unit"`
		TahunUuid      string `json:"tahun_uuid"`
		IdMasterTahun  *int   `json:"id_master_tahun"`
		Tahun          string `json:"tahun"`

		// Excel Fields
		IdMataProgram        *int    `json:"id_mata_program"`
		MataProgramUuid      string  `json:"mata_program_uuid"`
		IsuInternalEksternal string  `json:"isu_internal_eksternal"`
		DeskripsiRisiko      string  `json:"deskripsi_risiko"`
		KemungkinanAwal      *int    `json:"kemungkinan_awal"`
		DampakAwal           *int    `json:"dampak_awal"`
		TingkatRisikoAwal    string  `json:"tingkat_risiko_awal"`
		TindakanMitigasi     string  `json:"tindakan_mitigasi"`
		Peluang              string  `json:"peluang"`
		KenyataanResidu      string  `json:"kenyataan_residu"`
		KemungkinanResidu    *int    `json:"kemungkinan_residu"`
		DampakResidu         *int    `json:"dampak_residu"`
		TingkatRisikoResidu  *string `json:"tingkat_risiko_residu"`
		RencanaTindakLanjut  string  `json:"rencana_tindak_lanjut"`

		// Fallback & Legacy fields
		KodeRisiko       string `json:"kode_risiko"`
		Kategori         string `json:"kategori"`
		PernyataanRisiko string `json:"pernyataan_risiko"`
		Penyebab         string `json:"penyebab"`
		Dampak           string `json:"dampak"`
		Kemungkinan      *int   `json:"kemungkinan"`
		Keparahan        *int   `json:"keparahan"`
		TingkatRisiko    *int   `json:"tingkat_risiko"`
		LevelRisiko      string `json:"level_risiko"`
		RencanaMitigasi  string `json:"rencana_mitigasi"`
		Pengendalian     string `json:"pengendalian"`
		PenanggungJawab  string `json:"penanggung_jawab"`
		TargetWaktu      string `json:"target_waktu"`
		Status           string `json:"status"`
	}

	_ = c.BodyParser(&req)

	if req.IdMataProgram == nil && c.FormValue("id_mata_program") != "" {
		req.IdMataProgram = parseOptionalInt(c.FormValue("id_mata_program"))
	}
	if req.MataProgramUuid == "" {
		req.MataProgramUuid = c.FormValue("mata_program_uuid")
		if req.MataProgramUuid == "" {
			req.MataProgramUuid = c.FormValue("mata_program")
		}
	}
	if req.IsuInternalEksternal == "" {
		req.IsuInternalEksternal = c.FormValue("isu_internal_eksternal")
		if req.IsuInternalEksternal == "" {
			req.IsuInternalEksternal = c.FormValue("penyebab")
		}
	}
	if req.DeskripsiRisiko == "" {
		req.DeskripsiRisiko = c.FormValue("deskripsi_risiko")
		if req.DeskripsiRisiko == "" {
			req.DeskripsiRisiko = c.FormValue("pernyataan_risiko")
		}
	}
	if req.KemungkinanAwal == nil && c.FormValue("kemungkinan_awal") != "" {
		req.KemungkinanAwal = parseOptionalInt(c.FormValue("kemungkinan_awal"))
	}
	if req.KemungkinanAwal == nil && c.FormValue("kemungkinan") != "" {
		req.KemungkinanAwal = parseOptionalInt(c.FormValue("kemungkinan"))
	}
	if req.DampakAwal == nil && c.FormValue("dampak_awal") != "" {
		req.DampakAwal = parseOptionalInt(c.FormValue("dampak_awal"))
	}
	if req.DampakAwal == nil && c.FormValue("keparahan") != "" {
		req.DampakAwal = parseOptionalInt(c.FormValue("keparahan"))
	}
	if req.TingkatRisikoAwal == "" {
		req.TingkatRisikoAwal = c.FormValue("tingkat_risiko_awal")
		if req.TingkatRisikoAwal == "" {
			req.TingkatRisikoAwal = c.FormValue("level_risiko")
		}
	}
	if req.TindakanMitigasi == "" {
		req.TindakanMitigasi = c.FormValue("tindakan_mitigasi")
		if req.TindakanMitigasi == "" {
			req.TindakanMitigasi = c.FormValue("rencana_mitigasi")
		}
	}
	if req.Peluang == "" {
		req.Peluang = c.FormValue("peluang")
	}
	if req.KenyataanResidu == "" {
		req.KenyataanResidu = c.FormValue("kenyataan_residu")
	}
	if req.KemungkinanResidu == nil && c.FormValue("kemungkinan_residu") != "" {
		req.KemungkinanResidu = parseOptionalInt(c.FormValue("kemungkinan_residu"))
	}
	if req.DampakResidu == nil && c.FormValue("dampak_residu") != "" {
		req.DampakResidu = parseOptionalInt(c.FormValue("dampak_residu"))
	}
	if req.TingkatRisikoResidu == nil && c.FormValue("tingkat_risiko_residu") != "" {
		val := c.FormValue("tingkat_risiko_residu")
		req.TingkatRisikoResidu = &val
	}
	if req.RencanaTindakLanjut == "" {
		req.RencanaTindakLanjut = c.FormValue("rencana_tindak_lanjut")
	}
	if req.TahunUuid == "" {
		req.TahunUuid = c.FormValue("tahun_uuid")
		if req.TahunUuid == "" {
			req.TahunUuid = c.FormValue("id_master_tahun")
		}
	}
	if req.Tahun == "" {
		req.Tahun = c.FormValue("tahun")
	}
	if req.Status == "" {
		req.Status = c.FormValue("status")
	}

	// Scoping: Non-admin (auditee) dipaksa ke unit miliknya sendiri
	if level != "admin" {
		req.FakultasUuid = userFU
		req.IdFakultasUnit = nil
	} else {
		if req.FakultasUuid == "" {
			req.FakultasUuid = c.FormValue("fakultas_unit")
			if req.FakultasUuid == "" {
				req.FakultasUuid = c.FormValue("id_fakultas_unit")
			}
		}
	}

	idMP := 0
	if req.IdMataProgram != nil {
		idMP = *req.IdMataProgram
	}
	kAwal := 3
	if req.KemungkinanAwal != nil && *req.KemungkinanAwal > 0 {
		kAwal = *req.KemungkinanAwal
	}
	dAwal := 3
	if req.DampakAwal != nil && *req.DampakAwal > 0 {
		dAwal = *req.DampakAwal
	}

	cmd := CreateRisk.CreateRiskAssessmentCommand{
		FakultasUuid:         req.FakultasUuid,
		IdFakultasUnit:       req.IdFakultasUnit,
		TahunUuid:            req.TahunUuid,
		IdMasterTahun:        req.IdMasterTahun,
		Tahun:                req.Tahun,
		IdMataProgram:        idMP,
		MataProgramUuid:      req.MataProgramUuid,
		IsuInternalEksternal: req.IsuInternalEksternal,
		DeskripsiRisiko:      req.DeskripsiRisiko,
		KemungkinanAwal:      kAwal,
		DampakAwal:           dAwal,
		TingkatRisikoAwal:    req.TingkatRisikoAwal,
		TindakanMitigasi:     req.TindakanMitigasi,
		Peluang:              req.Peluang,
		KenyataanResidu:      req.KenyataanResidu,
		KemungkinanResidu:    req.KemungkinanResidu,
		DampakResidu:         req.DampakResidu,
		TingkatRisikoResidu:  req.TingkatRisikoResidu,
		RencanaTindakLanjut:  req.RencanaTindakLanjut,
		KodeRisiko:           req.KodeRisiko,
		Kategori:             req.Kategori,
		PernyataanRisiko:     req.PernyataanRisiko,
		Status:               req.Status,
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

	level, userFU := getUserInfo(c)

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
		if !matchUnit(existing.IdFakultasUnit, existing.UuidFakultasUnit, userFU) {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Akses ditolak: bukan unit anda"})
		}
	}

	var req struct {
		FakultasUuid   string `json:"fakultas_unit"`
		IdFakultasUnit *int   `json:"id_fakultas_unit"`
		TahunUuid      string `json:"tahun_uuid"`
		IdMasterTahun  *int   `json:"id_master_tahun"`
		Tahun          string `json:"tahun"`

		// Excel Fields
		IdMataProgram        *int    `json:"id_mata_program"`
		MataProgramUuid      string  `json:"mata_program_uuid"`
		IsuInternalEksternal string  `json:"isu_internal_eksternal"`
		DeskripsiRisiko      string  `json:"deskripsi_risiko"`
		KemungkinanAwal      *int    `json:"kemungkinan_awal"`
		DampakAwal           *int    `json:"dampak_awal"`
		TingkatRisikoAwal    string  `json:"tingkat_risiko_awal"`
		TindakanMitigasi     string  `json:"tindakan_mitigasi"`
		Peluang              string  `json:"peluang"`
		KenyataanResidu      string  `json:"kenyataan_residu"`
		KemungkinanResidu    *int    `json:"kemungkinan_residu"`
		DampakResidu         *int    `json:"dampak_residu"`
		TingkatRisikoResidu  *string `json:"tingkat_risiko_residu"`
		RencanaTindakLanjut  string  `json:"rencana_tindak_lanjut"`

		// Fallback & Legacy fields
		KodeRisiko       string `json:"kode_risiko"`
		Kategori         string `json:"kategori"`
		PernyataanRisiko string `json:"pernyataan_risiko"`
		Penyebab         string `json:"penyebab"`
		Dampak           string `json:"dampak"`
		Kemungkinan      *int   `json:"kemungkinan"`
		Keparahan        *int   `json:"keparahan"`
		TingkatRisiko    *int   `json:"tingkat_risiko"`
		LevelRisiko      string `json:"level_risiko"`
		RencanaMitigasi  string `json:"rencana_mitigasi"`
		Pengendalian     string `json:"pengendalian"`
		PenanggungJawab  string `json:"penanggung_jawab"`
		TargetWaktu      string `json:"target_waktu"`
		Status           string `json:"status"`
	}

	_ = c.BodyParser(&req)

	if req.IdMataProgram == nil && c.FormValue("id_mata_program") != "" {
		req.IdMataProgram = parseOptionalInt(c.FormValue("id_mata_program"))
	}
	if req.MataProgramUuid == "" {
		req.MataProgramUuid = c.FormValue("mata_program_uuid")
		if req.MataProgramUuid == "" {
			req.MataProgramUuid = c.FormValue("mata_program")
		}
	}
	if req.IsuInternalEksternal == "" {
		req.IsuInternalEksternal = c.FormValue("isu_internal_eksternal")
		if req.IsuInternalEksternal == "" {
			req.IsuInternalEksternal = c.FormValue("penyebab")
		}
	}
	if req.DeskripsiRisiko == "" {
		req.DeskripsiRisiko = c.FormValue("deskripsi_risiko")
		if req.DeskripsiRisiko == "" {
			req.DeskripsiRisiko = c.FormValue("pernyataan_risiko")
		}
	}
	if req.KemungkinanAwal == nil && c.FormValue("kemungkinan_awal") != "" {
		req.KemungkinanAwal = parseOptionalInt(c.FormValue("kemungkinan_awal"))
	}
	if req.KemungkinanAwal == nil && c.FormValue("kemungkinan") != "" {
		req.KemungkinanAwal = parseOptionalInt(c.FormValue("kemungkinan"))
	}
	if req.DampakAwal == nil && c.FormValue("dampak_awal") != "" {
		req.DampakAwal = parseOptionalInt(c.FormValue("dampak_awal"))
	}
	if req.DampakAwal == nil && c.FormValue("keparahan") != "" {
		req.DampakAwal = parseOptionalInt(c.FormValue("keparahan"))
	}
	if req.TingkatRisikoAwal == "" {
		req.TingkatRisikoAwal = c.FormValue("tingkat_risiko_awal")
		if req.TingkatRisikoAwal == "" {
			req.TingkatRisikoAwal = c.FormValue("level_risiko")
		}
	}
	if req.TindakanMitigasi == "" {
		req.TindakanMitigasi = c.FormValue("tindakan_mitigasi")
		if req.TindakanMitigasi == "" {
			req.TindakanMitigasi = c.FormValue("rencana_mitigasi")
		}
	}
	if req.Peluang == "" {
		req.Peluang = c.FormValue("peluang")
	}
	if req.KenyataanResidu == "" {
		req.KenyataanResidu = c.FormValue("kenyataan_residu")
	}
	if req.KemungkinanResidu == nil && c.FormValue("kemungkinan_residu") != "" {
		req.KemungkinanResidu = parseOptionalInt(c.FormValue("kemungkinan_residu"))
	}
	if req.DampakResidu == nil && c.FormValue("dampak_residu") != "" {
		req.DampakResidu = parseOptionalInt(c.FormValue("dampak_residu"))
	}
	if req.TingkatRisikoResidu == nil && c.FormValue("tingkat_risiko_residu") != "" {
		val := c.FormValue("tingkat_risiko_residu")
		req.TingkatRisikoResidu = &val
	}
	if req.RencanaTindakLanjut == "" {
		req.RencanaTindakLanjut = c.FormValue("rencana_tindak_lanjut")
	}
	if req.TahunUuid == "" {
		req.TahunUuid = c.FormValue("tahun_uuid")
		if req.TahunUuid == "" {
			req.TahunUuid = c.FormValue("id_master_tahun")
		}
	}
	if req.Tahun == "" {
		req.Tahun = c.FormValue("tahun")
	}
	if req.Status == "" {
		req.Status = c.FormValue("status")
	}

	// Scoping: Non-admin tidak boleh memindahkan unit
	if level != "admin" {
		req.FakultasUuid = userFU
		req.IdFakultasUnit = nil
	} else {
		if req.FakultasUuid == "" {
			req.FakultasUuid = c.FormValue("fakultas_unit")
			if req.FakultasUuid == "" {
				req.FakultasUuid = c.FormValue("id_fakultas_unit")
			}
		}
	}

	idMP := 0
	if req.IdMataProgram != nil {
		idMP = *req.IdMataProgram
	}
	kAwal := 0
	if req.KemungkinanAwal != nil {
		kAwal = *req.KemungkinanAwal
	}
	dAwal := 0
	if req.DampakAwal != nil {
		dAwal = *req.DampakAwal
	}

	cmd := UpdateRisk.UpdateRiskAssessmentCommand{
		UUID:                 rawUuid,
		FakultasUuid:         req.FakultasUuid,
		IdFakultasUnit:       req.IdFakultasUnit,
		TahunUuid:            req.TahunUuid,
		IdMasterTahun:        req.IdMasterTahun,
		Tahun:                req.Tahun,
		IdMataProgram:        idMP,
		MataProgramUuid:      req.MataProgramUuid,
		IsuInternalEksternal: req.IsuInternalEksternal,
		DeskripsiRisiko:      req.DeskripsiRisiko,
		KemungkinanAwal:      kAwal,
		DampakAwal:           dAwal,
		TingkatRisikoAwal:    req.TingkatRisikoAwal,
		TindakanMitigasi:     req.TindakanMitigasi,
		Peluang:              req.Peluang,
		KenyataanResidu:      req.KenyataanResidu,
		KemungkinanResidu:    req.KemungkinanResidu,
		DampakResidu:         req.DampakResidu,
		TingkatRisikoResidu:  req.TingkatRisikoResidu,
		RencanaTindakLanjut:  req.RencanaTindakLanjut,
		KodeRisiko:           req.KodeRisiko,
		Kategori:             req.Kategori,
		PernyataanRisiko:     req.PernyataanRisiko,
		Status:               req.Status,
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

	level, userFU := getUserInfo(c)

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
		if !matchUnit(existing.IdFakultasUnit, existing.UuidFakultasUnit, userFU) {
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

	level, userFU := getUserInfo(c)

	query := GetRisk.GetRiskAssessmentByUuidQuery{UUID: parsedUuid}
	result, err := mediatr.Send[GetRisk.GetRiskAssessmentByUuidQuery, *domainRisk.RiskAssessmentDefault](context.Background(), query)
	if err != nil {
		return commoninfra.HandleError(c, err)
	}
	if result == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Risk assessment not found"})
	}

	if level != "admin" && !matchUnit(result.IdFakultasUnit, result.UuidFakultasUnit, userFU) {
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

	level, userFU := getUserInfo(c)

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
// Master Data Structs & Handlers
// =======================================================
type RiskCriteriaItem struct {
	ID        int    `json:"id" gorm:"column:id"`
	Tipe      string `json:"tipe" gorm:"column:tipe"`
	Skor      int    `json:"skor" gorm:"column:skor"`
	Label     string `json:"label" gorm:"column:label"`
	Deskripsi string `json:"deskripsi" gorm:"column:deskripsi"`
}

type RiskLevelItem struct {
	ID         int    `json:"id" gorm:"column:id"`
	Kode       string `json:"kode" gorm:"column:kode"`
	NamaLevel  string `json:"nama_level" gorm:"column:nama_level"`
	WarnaBadge string `json:"warna_badge" gorm:"column:warna_badge"`
	Keterangan string `json:"keterangan" gorm:"column:keterangan"`
	MinSkor    int    `json:"min_skor" gorm:"column:min_skor"`
	MaxSkor    int    `json:"max_skor" gorm:"column:max_skor"`
}

type RiskCategoryItem struct {
	ID        int    `json:"id" gorm:"column:id"`
	Nama      string `json:"nama" gorm:"column:nama"`
	Deskripsi string `json:"deskripsi" gorm:"column:deskripsi"`
}

type RiskStatusItem struct {
	ID        int    `json:"id" gorm:"column:id"`
	Nama      string `json:"nama" gorm:"column:nama"`
	Deskripsi string `json:"deskripsi" gorm:"column:deskripsi"`
	Urutan    int    `json:"urutan" gorm:"column:urutan"`
}

type RiskMataProgramItem struct {
	ID             int     `json:"id" gorm:"column:id"`
	UUID           string  `json:"uuid" gorm:"column:uuid"`
	MataProgram    string  `json:"mata_program" gorm:"column:mata_program"`
	IdFakultasUnit *int    `json:"id_fakultas_unit" gorm:"column:id_fakultas_unit"`
	IdMasterTahun  *int    `json:"id_master_tahun" gorm:"column:id_master_tahun"`
	Tahun          *string `json:"tahun" gorm:"column:tahun"`
}

func GetRiskCriteriaHandler(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		tipe := c.Query("tipe")
		var items []RiskCriteriaItem
		q := db.Table("risk_criteria").Order("tipe ASC, skor ASC")
		if tipe != "" {
			q = q.Where("tipe = ?", strings.ToUpper(tipe))
		}
		if err := q.Find(&items).Error; err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(items)
	}
}

func GetRiskLevelsHandler(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var items []RiskLevelItem
		if err := db.Table("risk_levels").Order("min_skor ASC, id ASC").Find(&items).Error; err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(items)
	}
}

func GetRiskCategoriesHandler(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var items []RiskCategoryItem
		if err := db.Table("risk_categories").Order("id ASC").Find(&items).Error; err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(items)
	}
}

func GetRiskStatusesHandler(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var items []RiskStatusItem
		if err := db.Table("risk_statuses").Order("urutan ASC, id ASC").Find(&items).Error; err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(items)
	}
}

func GetRiskProgramsHandler(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var programs []RiskMataProgramItem
		q := db.Table("mata_program mp").
			Select("mp.id, mp.uuid, mp.mata_program, mp.id_fakultas_unit, mp.id_master_tahun, mt.tahun").
			Joins("LEFT JOIN master_tahun mt ON mp.id_master_tahun = mt.id").
			Order("mp.mata_program ASC")

		if fu := c.Query("fakultas_unit"); fu != "" {
			if id, err := strconv.Atoi(fu); err == nil {
				q = q.Where("mp.id_fakultas_unit = ?", id)
			} else {
				q = q.Joins("LEFT JOIN sijamu_fakultas_unit fu ON mp.id_fakultas_unit = fu.id").
					Where("fu.uuid = ?", fu)
			}
		}
		if thn := c.Query("tahun"); thn != "" {
			if id, err := strconv.Atoi(thn); err == nil {
				q = q.Where("mp.id_master_tahun = ? OR mt.tahun = ?", id, thn)
			} else {
				q = q.Where("mt.uuid = ? OR mt.tahun = ?", thn, thn)
			}
		}
		if search := c.Query("search"); search != "" {
			q = q.Where("mp.mata_program LIKE ?", "%"+search+"%")
		}

		if err := q.Find(&programs).Error; err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(programs)
	}
}

func GetRiskOptionsHandler(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var criteria []RiskCriteriaItem
		_ = db.Table("risk_criteria").Order("tipe ASC, skor ASC").Find(&criteria)

		var levels []RiskLevelItem
		_ = db.Table("risk_levels").Order("min_skor ASC, id ASC").Find(&levels)

		var categories []RiskCategoryItem
		_ = db.Table("risk_categories").Order("id ASC").Find(&categories)

		var statuses []RiskStatusItem
		_ = db.Table("risk_statuses").Order("urutan ASC, id ASC").Find(&statuses)

		var programs []RiskMataProgramItem
		qMP := db.Table("mata_program mp").
			Select("mp.id, mp.uuid, mp.mata_program, mp.id_fakultas_unit, mp.id_master_tahun, mt.tahun").
			Joins("LEFT JOIN master_tahun mt ON mp.id_master_tahun = mt.id").
			Order("mp.mata_program ASC")

		if fu := c.Query("fakultas_unit"); fu != "" {
			if id, err := strconv.Atoi(fu); err == nil {
				qMP = qMP.Where("mp.id_fakultas_unit = ?", id)
			}
		}
		if thn := c.Query("tahun"); thn != "" {
			if id, err := strconv.Atoi(thn); err == nil {
				qMP = qMP.Where("mp.id_master_tahun = ? OR mt.tahun = ?", id, thn)
			} else {
				qMP = qMP.Where("mt.uuid = ? OR mt.tahun = ?", thn, thn)
			}
		}
		_ = qMP.Find(&programs)

		var kemungkinan []RiskCriteriaItem
		var dampak []RiskCriteriaItem
		for _, cr := range criteria {
			if strings.EqualFold(cr.Tipe, "KEMUNGKINAN") {
				kemungkinan = append(kemungkinan, cr)
			} else if strings.EqualFold(cr.Tipe, "DAMPAK") {
				dampak = append(dampak, cr)
			}
		}

		return c.JSON(fiber.Map{
			"criteria":    criteria,
			"kemungkinan": kemungkinan,
			"dampak":      dampak,
			"levels":      levels,
			"categories":  categories,
			"statuses":    statuses,
			"programs":    programs,
		})
	}
}

func ModuleRiskAssessment(app *fiber.App, db *gorm.DB) {
	roles := []string{"admin", "auditee", "fakultas", "unit", "upml", "upmf", "upms"}
	whoamiURL := "http://localhost:3000/whoami"

	// Dynamic Master Data Endpoints
	app.Get("/risk_options", commonpresentation.SmartCompress(), commonpresentation.JWTMiddleware(), GetRiskOptionsHandler(db))
	app.Get("/risk_criteria", commonpresentation.SmartCompress(), commonpresentation.JWTMiddleware(), GetRiskCriteriaHandler(db))
	app.Get("/risk_levels", commonpresentation.SmartCompress(), commonpresentation.JWTMiddleware(), GetRiskLevelsHandler(db))
	app.Get("/risk_categories", commonpresentation.SmartCompress(), commonpresentation.JWTMiddleware(), GetRiskCategoriesHandler(db))
	app.Get("/risk_statuses", commonpresentation.SmartCompress(), commonpresentation.JWTMiddleware(), GetRiskStatusesHandler(db))
	app.Get("/risk_programs", commonpresentation.SmartCompress(), commonpresentation.JWTMiddleware(), GetRiskProgramsHandler(db))

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
