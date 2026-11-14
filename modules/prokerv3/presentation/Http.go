package presentation

import (
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/mehdihadeli/go-mediatr"
	"gorm.io/gorm"

	commondomain "UnpakSiamida/common/domain"
	"UnpakSiamida/common/helper"
	commoninfra "UnpakSiamida/common/infrastructure"
	commonpresentation "UnpakSiamida/common/presentation"

	checkJadwal "UnpakSiamida/modules/prokerv3/application/CheckJadwalProkerV3"
	createJadwal "UnpakSiamida/modules/prokerv3/application/CreateJadwalProkerV3"
	createDokumen "UnpakSiamida/modules/prokerv3/application/CreateProkerDokumen"
	createProker "UnpakSiamida/modules/prokerv3/application/CreateProkerV3"
	deleteDokumen "UnpakSiamida/modules/prokerv3/application/DeleteDokumenV3"
	deleteJadwal "UnpakSiamida/modules/prokerv3/application/DeleteJadwalProkerV3"
	deletePagu "UnpakSiamida/modules/prokerv3/application/DeletePaguAnggaran"
	deleteProker "UnpakSiamida/modules/prokerv3/application/DeleteProkerV3"
	getAllJadwal "UnpakSiamida/modules/prokerv3/application/GetAllJadwalProkerV3"
	getAllPagu "UnpakSiamida/modules/prokerv3/application/GetAllPaguAnggaran"
	getAllProker "UnpakSiamida/modules/prokerv3/application/GetAllProkerV3"
	getCoa "UnpakSiamida/modules/prokerv3/application/GetCoa"
	getDashboard "UnpakSiamida/modules/prokerv3/application/GetDashboardProkerMonitoring"
	getDokumen "UnpakSiamida/modules/prokerv3/application/GetDokumenV3"
	getHistory "UnpakSiamida/modules/prokerv3/application/GetHistoryStateV3"
	getJadwal "UnpakSiamida/modules/prokerv3/application/GetJadwalProkerV3"
	GetProkerBatalSetting "UnpakSiamida/modules/prokerv3/application/GetProkerBatalSetting"
	getProker "UnpakSiamida/modules/prokerv3/application/GetProkerV3"
	getOptions "UnpakSiamida/modules/prokerv3/application/GetProkerV3Options"
	getSankey "UnpakSiamida/modules/prokerv3/application/GetSankeyData"
	getTargetUnpak "UnpakSiamida/modules/prokerv3/application/GetTargetUnpak"
	getUnitPagu "UnpakSiamida/modules/prokerv3/application/GetUnitPaguSummary"
	savePagu "UnpakSiamida/modules/prokerv3/application/SavePaguAnggaran"
	updateBatal "UnpakSiamida/modules/prokerv3/application/UpdateBatalProker"
	updateBeban "UnpakSiamida/modules/prokerv3/application/UpdateBebanYayasan"
	updateJadwal "UnpakSiamida/modules/prokerv3/application/UpdateJadwalProkerV3"
	UpdateProkerBatalSetting "UnpakSiamida/modules/prokerv3/application/UpdateProkerBatalSetting"
	updateProker "UnpakSiamida/modules/prokerv3/application/UpdateProkerV3"
	updateStatus "UnpakSiamida/modules/prokerv3/application/UpdateStatusV3"
	updateVerifikasiDokumen "UnpakSiamida/modules/prokerv3/application/UpdateVerifikasiDokumenV3"
	domain "UnpakSiamida/modules/prokerv3/domain"
)

var prokerDb *gorm.DB

// Helper pengecekan hak akses unit
func canAccessProkerUnit(c *fiber.Ctx, itemUnitID *int, itemUnitUUID string, itemUnitName string, parentFakultas ...string) bool {
	level, userFU := commonpresentation.GetUserInfo(c)
	level = strings.ToLower(strings.TrimSpace(level))

	// Roles dengan hak akses global ke seluruh unit
	if level == "admin" || level == "anggaran" || level == "keuangan" || level == "evaluator" || level == "lpm" {
		return true
	}

	pFaculty := ""
	if len(parentFakultas) > 0 {
		pFaculty = strings.TrimSpace(parentFakultas[0])
	}

	u, _ := c.Locals("user").(*commonpresentation.Account)
	activeTahun := strings.TrimSpace(c.Get("X-Active-Year"))
	if activeTahun == "" {
		activeTahun = strings.TrimSpace(c.Get("X-Position-Year"))
	}
	if activeTahun == "" {
		activeTahun = strings.TrimSpace(c.Get("X-Tahun"))
	}
	if activeTahun == "" {
		activeTahun = strings.TrimSpace(c.Query("filter_tahun"))
	}
	if activeTahun == "" {
		activeTahun = strings.TrimSpace(c.Query("tahun"))
	}

	// Kumpulkan penugasan unit untuk role aktif dan tahun aktif
	var assignedUnitIDs []uint
	var assignedUnitNames []string

	if u != nil {
		for _, r := range u.ExtraRole {
			if level != "" && !strings.EqualFold(r.Role, level) && level != "user" {
				continue
			}
			if activeTahun != "" && r.Tahun != "" && r.Tahun != activeTahun {
				continue
			}
			if r.IdFakultasUnit > 0 {
				assignedUnitIDs = append(assignedUnitIDs, r.IdFakultasUnit)
			}
			if strings.TrimSpace(r.NamaFakultasUnit) != "" {
				assignedUnitNames = append(assignedUnitNames, strings.TrimSpace(r.NamaFakultasUnit))
			}
		}
	}

	// Fallback ke userFU / u.FakultasUnit jika penugasan ExtraRole tidak terdaftar
	if len(assignedUnitIDs) == 0 && userFU != "" {
		if id, err := strconv.Atoi(userFU); err == nil && id > 0 {
			assignedUnitIDs = append(assignedUnitIDs, uint(id))
		}
	}
	if len(assignedUnitIDs) == 0 && u != nil && u.FakultasUnit != "" {
		if id, err := strconv.Atoi(u.FakultasUnit); err == nil && id > 0 {
			assignedUnitIDs = append(assignedUnitIDs, uint(id))
		}
	}

	// 1. Cek kecocokan langsung (Direct Match) dengan unit penugasan
	if itemUnitID != nil && *itemUnitID > 0 {
		for _, aid := range assignedUnitIDs {
			if uint(*itemUnitID) == aid {
				return true
			}
		}
	}
	if itemUnitName != "" {
		for _, aname := range assignedUnitNames {
			if strings.EqualFold(strings.TrimSpace(itemUnitName), strings.TrimSpace(aname)) {
				return true
			}
		}
	}
	if len(assignedUnitIDs) == 0 && helper.MatchUnit(itemUnitID, itemUnitUUID, userFU) {
		return true
	}

	// 2. Jika bukan direct match:
	// UPML, UPMS, dan USER biasa strictly scoped ke unit penugasan masing-masing
	if level == "upml" || level == "upms" || level == "user" {
		return false
	}

	// 3. UPMF: Dapat mengakses seluruh prodi/unit yang berada di bawah fakultas penugasannya (1 tipe fakultas)
	if level == "upmf" {
		if pFaculty != "" {
			for _, aname := range assignedUnitNames {
				if strings.EqualFold(pFaculty, aname) {
					return true
				}
			}
		}

		if prokerDb != nil && len(assignedUnitIDs) > 0 {
			var matchCount int64
			var q *gorm.DB
			if itemUnitID != nil && *itemUnitID > 0 {
				q = prokerDb.Table("v_fakultas_unit fu").
					Where("fu.id = ? AND (fu.id IN (?) OR (fu.fakultas IS NOT NULL AND fu.fakultas != '' AND LOWER(TRIM(fu.fakultas)) IN (SELECT LOWER(TRIM(COALESCE(nama_fak_prod_unit, ''))) FROM v_fakultas_unit WHERE id IN (?) AND type = 'fakultas')))", *itemUnitID, assignedUnitIDs, assignedUnitIDs)
			} else if itemUnitUUID != "" {
				q = prokerDb.Table("v_fakultas_unit fu").
					Where("fu.uuid = ? AND (fu.id IN (?) OR (fu.fakultas IS NOT NULL AND fu.fakultas != '' AND LOWER(TRIM(fu.fakultas)) IN (SELECT LOWER(TRIM(COALESCE(nama_fak_prod_unit, ''))) FROM v_fakultas_unit WHERE id IN (?) AND type = 'fakultas')))", itemUnitUUID, assignedUnitIDs, assignedUnitIDs)
			}
			if q != nil {
				if err := q.Count(&matchCount).Error; err == nil && matchCount > 0 {
					return true
				}
			}
		}
	}

	return false
}

func resolveActorForDelete(c *fiber.Ctx, level string) (role string, actorID *int, actorName string, actorUsername string) {
	role = strings.ToLower(strings.TrimSpace(level))
	if role == "" {
		role = "user"
	}

	if u, ok := c.Locals("user").(*commonpresentation.Account); ok && u != nil {
		actorUsername = strings.TrimSpace(u.NidnUsername)
		actorName = strings.TrimSpace(u.Name)
		if actorName == "" {
			actorName = actorUsername
		}
	}
	if actorName == "" {
		if rawName, ok := c.Locals("name").(string); ok && strings.TrimSpace(rawName) != "" {
			actorName = strings.TrimSpace(rawName)
		}
	}
	if actorUsername == "" {
		if rawUsername, ok := c.Locals("username").(string); ok && strings.TrimSpace(rawUsername) != "" {
			actorUsername = strings.TrimSpace(rawUsername)
		}
	}
	if actorID == nil {
		if rawUID, ok := c.Locals("user_id").(string); ok && strings.TrimSpace(rawUID) != "" {
			if parsedID, err := strconv.Atoi(strings.TrimSpace(rawUID)); err == nil && parsedID > 0 {
				actorID = &parsedID
			}
		}
	}
	return
}

// =======================================================
// POST /proker_v3
// =======================================================
func CreateProkerV3Handlerfunc(c *fiber.Ctx) error {
	level, userFU := commonpresentation.GetUserInfo(c)

	var cmd createProker.CreateProkerV3Command

	if val := c.FormValue("tahun"); val != "" {
		cmd.Tahun = strings.TrimSpace(val)
	}
	if val := c.FormValue("fakultas_unit"); val != "" {
		cmd.FakultasUuid = strings.TrimSpace(val)
	}
	if val := c.FormValue("nama_program"); val != "" {
		cmd.NamaProgram = strings.TrimSpace(val)
	}
	if val := c.FormValue("id_mata_program"); val != "" {
		if id, err := strconv.Atoi(val); err == nil {
			cmd.IdMataProgram = &id
		}
	}
	if val := c.FormValue("id_indikator"); val != "" {
		if id, err := strconv.Atoi(val); err == nil {
			cmd.IdIndikator = &id
		}
	}
	if val := c.FormValue("id_fakultas_unit"); val != "" {
		if id, err := strconv.Atoi(val); err == nil {
			cmd.IdFakultasUnit = &id
		}
	}
	if val := c.FormValue("no_proker"); val != "" {
		if np, err := strconv.ParseInt(val, 10, 64); err == nil {
			cmd.NoProker = &np
		}
	}
	if val := c.FormValue("standar_sub_standar"); val != "" {
		cmd.StandarSubStandar = strings.TrimSpace(val)
	}
	if val := c.FormValue("indikator_kinerja"); val != "" {
		cmd.IndikatorKinerja = strings.TrimSpace(val)
	}
	if val := c.FormValue("target_unpak"); val != "" {
		cmd.TargetUnpak = strings.TrimSpace(val)
	}
	if val := c.FormValue("capaian_terdahulu"); val != "" {
		cmd.CapaianUnitLalu = strings.TrimSpace(val)
	}
	if val := c.FormValue("capaian_unit_lalu"); val != "" {
		cmd.CapaianUnitLalu = strings.TrimSpace(val)
	}
	if val := c.FormValue("target_unit"); val != "" {
		cmd.TargetUnit = strings.TrimSpace(val)
	}
	if val := c.FormValue("output"); val != "" {
		cmd.Output = strings.TrimSpace(val)
	}
	if val := c.FormValue("outcome"); val != "" {
		cmd.Outcome = strings.TrimSpace(val)
	}
	if val := c.FormValue("pic"); val != "" {
		cmd.PIC = strings.TrimSpace(val)
	}
	if val := c.FormValue("anggaran"); val != "" {
		if ang, err := strconv.ParseFloat(val, 64); err == nil {
			cmd.Anggaran = ang
		}
	}
	if val := c.FormValue("dokumen_referensi"); val != "" {
		cmd.DokumenReferensi = strings.TrimSpace(val)
	}
	if val := c.FormValue("luaran_produk"); val != "" {
		cmd.LuaranProduk = strings.TrimSpace(val)
	}
	if val := c.FormValue("aktivitas"); val != "" {
		cmd.Aktivitas = strings.TrimSpace(val)
	}
	if val := c.FormValue("status"); val != "" {
		cmd.Status = strings.TrimSpace(val)
	}
	if val := c.FormValue("beban_yayasan"); val != "" {
		cmd.BebanYayasan = val == "1" || strings.ToLower(val) == "true"
	}
	if val := c.FormValue("scopes"); val != "" {
		cmd.ScopesRaw = strings.TrimSpace(val)
		var scopes []createProker.ScopeItemInput
		if err := json.Unmarshal([]byte(val), &scopes); err == nil {
			cmd.Scopes = scopes
		}
	} else if cmd.ScopesRaw != "" && len(cmd.Scopes) == 0 {
		var scopes []createProker.ScopeItemInput
		if err := json.Unmarshal([]byte(cmd.ScopesRaw), &scopes); err == nil {
			cmd.Scopes = scopes
		}
	}

	if level != "admin" {
		if userFU != "" && cmd.IdFakultasUnit == nil {
			if parsed, err := strconv.Atoi(userFU); err == nil {
				cmd.IdFakultasUnit = &parsed
			}
		}
	}

	uuidResult, err := mediatr.Send[createProker.CreateProkerV3Command, string](c.Context(), cmd)
	if err != nil {
		return commoninfra.HandleError(c, err)
	}

	return commonpresentation.JsonUUID(c, uuidResult)
}

// =======================================================
// PUT /proker_v3/:uuid
// =======================================================
func UpdateProkerV3Handlerfunc(c *fiber.Ctx) error {
	rawUuid := c.Params("uuid")
	if rawUuid == "" {
		rawUuid = c.Query("uuid")
	}

	var cmd updateProker.UpdateProkerV3Command

	if val := c.FormValue("tahun"); val != "" {
		cmd.Tahun = strings.TrimSpace(val)
	}
	if val := c.FormValue("fakultas_unit"); val != "" {
		cmd.FakultasUuid = strings.TrimSpace(val)
	}
	if val := c.FormValue("nama_program"); val != "" {
		cmd.NamaProgram = strings.TrimSpace(val)
	}
	if val := c.FormValue("id_mata_program"); val != "" {
		if id, err := strconv.Atoi(val); err == nil {
			cmd.IdMataProgram = &id
		}
	}
	if val := c.FormValue("id_indikator"); val != "" {
		if id, err := strconv.Atoi(val); err == nil {
			cmd.IdIndikator = &id
		}
	}
	if val := c.FormValue("id_fakultas_unit"); val != "" {
		if id, err := strconv.Atoi(val); err == nil {
			cmd.IdFakultasUnit = &id
		}
	}
	if val := c.FormValue("no_proker"); val != "" {
		if np, err := strconv.ParseInt(val, 10, 64); err == nil {
			cmd.NoProker = &np
		}
	}
	if val := c.FormValue("standar_sub_standar"); val != "" {
		cmd.StandarSubStandar = strings.TrimSpace(val)
	}
	if val := c.FormValue("indikator_kinerja"); val != "" {
		cmd.IndikatorKinerja = strings.TrimSpace(val)
	}
	if val := c.FormValue("target_unpak"); val != "" {
		cmd.TargetUnpak = strings.TrimSpace(val)
	}
	if val := c.FormValue("capaian_terdahulu"); val != "" {
		cmd.CapaianUnitLalu = strings.TrimSpace(val)
	}
	if val := c.FormValue("capaian_unit_lalu"); val != "" {
		cmd.CapaianUnitLalu = strings.TrimSpace(val)
	}
	if val := c.FormValue("target_unit"); val != "" {
		cmd.TargetUnit = strings.TrimSpace(val)
	}
	if val := c.FormValue("output"); val != "" {
		cmd.Output = strings.TrimSpace(val)
	}
	if val := c.FormValue("outcome"); val != "" {
		cmd.Outcome = strings.TrimSpace(val)
	}
	if val := c.FormValue("pic"); val != "" {
		cmd.PIC = strings.TrimSpace(val)
	}
	if val := c.FormValue("anggaran"); val != "" {
		if ang, err := strconv.ParseFloat(val, 64); err == nil {
			cmd.Anggaran = ang
		}
	}
	if val := c.FormValue("dokumen_referensi"); val != "" {
		cmd.DokumenReferensi = strings.TrimSpace(val)
	}
	if val := c.FormValue("luaran_produk"); val != "" {
		cmd.LuaranProduk = strings.TrimSpace(val)
	}
	if val := c.FormValue("aktivitas"); val != "" {
		cmd.Aktivitas = strings.TrimSpace(val)
	}
	if val := c.FormValue("status"); val != "" {
		cmd.Status = strings.TrimSpace(val)
	}
	if val := c.FormValue("beban_yayasan"); val != "" {
		b := val == "1" || strings.ToLower(val) == "true"
		cmd.BebanYayasan = &b
	}
	if val := c.FormValue("scopes"); val != "" {
		cmd.ScopesRaw = strings.TrimSpace(val)
		var scopes []createProker.ScopeItemInput
		if err := json.Unmarshal([]byte(val), &scopes); err == nil {
			cmd.Scopes = scopes
		}
	} else if cmd.ScopesRaw != "" && len(cmd.Scopes) == 0 {
		var scopes []createProker.ScopeItemInput
		if err := json.Unmarshal([]byte(cmd.ScopesRaw), &scopes); err == nil {
			cmd.Scopes = scopes
		}
	}
	cmd.UUID = rawUuid

	level, _ := commonpresentation.GetUserInfo(c)
	if level != "admin" {
		parsedUuid, err := uuid.Parse(rawUuid)
		if err == nil {
			query := getProker.GetProkerV3ByUuidQuery{UUID: parsedUuid}
			existing, err := mediatr.Send[getProker.GetProkerV3ByUuidQuery, *domain.ProkerV3Default](c.Context(), query)
			if err == nil && existing != nil {
				if existing.DeletedAt != nil {
					return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Program kerja sudah dihapus dan tidak dapat diubah"})
				}
				if !canAccessProkerUnit(c, existing.IdFakultasUnit, existing.UuidFakultasUnit, existing.NamaFakultasUnit, existing.ParentFakultas) {
					return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"code": "FORBIDDEN", "error": "Akses ditolak: bukan unit anda"})
				}
			}
		}
	}

	uuidResult, err := mediatr.Send[updateProker.UpdateProkerV3Command, string](c.Context(), cmd)
	if err != nil {
		return commoninfra.HandleError(c, err)
	}

	return commonpresentation.JsonUUID(c, uuidResult)
}

// =======================================================
// PUT /proker_v3/:uuid/beban_yayasan
// =======================================================
func UpdateBebanYayasanHandler(c *fiber.Ctx) error {
	rawUuid := c.Params("uuid")
	if rawUuid == "" {
		rawUuid = c.Query("uuid")
	}
	parsedUuid, err := uuid.Parse(rawUuid)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid proker UUID"})
	}

	level, _ := commonpresentation.GetUserInfo(c)
	if level != "admin" {
		query := getProker.GetProkerV3ByUuidQuery{UUID: parsedUuid}
		existing, err := mediatr.Send[getProker.GetProkerV3ByUuidQuery, *domain.ProkerV3Default](c.Context(), query)
		if err == nil && existing != nil {
			if !canAccessProkerUnit(c, existing.IdFakultasUnit, existing.UuidFakultasUnit, existing.NamaFakultasUnit, existing.ParentFakultas) {
				return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"code": "FORBIDDEN", "error": "Akses ditolak: bukan unit anda"})
			}
		}
	}

	var newBebanPtr *bool
	if val := strings.TrimSpace(c.FormValue("beban_yayasan")); val != "" {
		b := val == "1" || strings.ToLower(val) == "true"
		newBebanPtr = &b
	}

	result, err := mediatr.Send[updateBeban.UpdateBebanYayasanCommand, *updateBeban.UpdateBebanYayasanResult](
		c.Context(),
		updateBeban.UpdateBebanYayasanCommand{
			UUID:         parsedUuid,
			BebanYayasan: newBebanPtr,
		},
	)
	if err != nil {
		return commoninfra.HandleError(c, err)
	}

	return c.JSON(result)
}

// =======================================================
// PUT /proker_v3/:uuid/batal
// =======================================================
func UpdateBatalProkerHandler(c *fiber.Ctx) error {
	rawUuid := c.Params("uuid")
	if rawUuid == "" {
		rawUuid = c.Query("uuid")
	}
	parsedUuid, err := uuid.Parse(rawUuid)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid proker UUID"})
	}

	level, _ := commonpresentation.GetUserInfo(c)
	if level != "admin" {
		query := getProker.GetProkerV3ByUuidQuery{UUID: parsedUuid}
		existing, err := mediatr.Send[getProker.GetProkerV3ByUuidQuery, *domain.ProkerV3Default](c.Context(), query)
		if err == nil && existing != nil {
			if !canAccessProkerUnit(c, existing.IdFakultasUnit, existing.UuidFakultasUnit, existing.NamaFakultasUnit, existing.ParentFakultas) {
				return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"code": "FORBIDDEN", "error": "Akses ditolak: bukan unit anda"})
			}
		}

		var targetUnitID *int
		if existing != nil {
			targetUnitID = existing.IdFakultasUnit
		}
		batalSetting, errSetting := mediatr.Send[
			GetProkerBatalSetting.GetProkerBatalSettingQuery,
			*GetProkerBatalSetting.GetProkerBatalSettingResult,
		](c.Context(), GetProkerBatalSetting.GetProkerBatalSettingQuery{
			TargetUnitID: targetUnitID,
		})
		if errSetting != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal memeriksa pengaturan Proker Batal"})
		}
		if batalSetting != nil && !batalSetting.IsActive {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": batalSetting.Message})
		}
	}

	var newBatalPtr *bool
	if val := strings.TrimSpace(c.FormValue("is_batal")); val != "" {
		b := val == "1" || strings.ToLower(val) == "true"
		newBatalPtr = &b
	}

	catatanBatal := strings.TrimSpace(c.FormValue("catatan_batal"))

	result, err := mediatr.Send[updateBatal.UpdateBatalProkerCommand, *updateBatal.UpdateBatalProkerResult](
		c.Context(),
		updateBatal.UpdateBatalProkerCommand{
			UUID:         parsedUuid,
			IsBatal:      newBatalPtr,
			CatatanBatal: catatanBatal,
		},
	)
	if err != nil {
		return commoninfra.HandleError(c, err)
	}

	return c.JSON(result)
}

// =======================================================
// DELETE /proker_v3/:uuid
// =======================================================
func DeleteProkerV3Handlerfunc(c *fiber.Ctx) error {
	rawUuid := c.Params("uuid")
	if rawUuid == "" {
		rawUuid = c.Query("uuid")
	}

	parsedUuid, err := uuid.Parse(strings.TrimSpace(rawUuid))
	if err != nil {
		return commoninfra.HandleError(c, domain.ProkerNotFound())
	}

	query := getProker.GetProkerV3ByUuidQuery{UUID: parsedUuid}
	existing, err := mediatr.Send[getProker.GetProkerV3ByUuidQuery, *domain.ProkerV3Default](c.Context(), query)
	if err != nil {
		return commoninfra.HandleError(c, err)
	}
	if existing == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Program kerja tidak ditemukan"})
	}

	level, _ := commonpresentation.GetUserInfo(c)
	isAdmin := strings.EqualFold(strings.TrimSpace(level), "admin")

	if !isAdmin {
		if existing.DeletedAt != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Program kerja sudah berada di tab Buang"})
		}
		if !canAccessProkerUnit(c, existing.IdFakultasUnit, existing.UuidFakultasUnit, existing.NamaFakultasUnit, existing.ParentFakultas) {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"code": "FORBIDDEN", "error": "Akses ditolak: bukan unit anda"})
		}
	}

	type deleteProkerBody struct {
		Catatan      string `json:"catatan"`
		Alasan       string `json:"alasan"`
		CatatanHapus string `json:"catatan_hapus"`
	}
	var body deleteProkerBody
	_ = c.BodyParser(&body)

	catatan := strings.TrimSpace(body.Catatan)
	if catatan == "" {
		catatan = strings.TrimSpace(body.Alasan)
	}
	if catatan == "" {
		catatan = strings.TrimSpace(body.CatatanHapus)
	}
	if catatan == "" {
		catatan = strings.TrimSpace(c.FormValue("catatan"))
	}
	if catatan == "" {
		catatan = strings.TrimSpace(c.FormValue("alasan"))
	}
	if catatan == "" {
		catatan = strings.TrimSpace(c.FormValue("catatan_hapus"))
	}
	if catatan == "" {
		catatan = strings.TrimSpace(c.Query("catatan"))
	}
	if catatan == "" {
		catatan = strings.TrimSpace(c.Query("alasan"))
	}
	if catatan == "" {
		catatan = strings.TrimSpace(c.Query("catatan_hapus"))
	}

	softDelete := !isAdmin || catatan != ""
	if softDelete && existing.DeletedAt != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Program kerja sudah berada di tab Buang"})
	}

	role, actorID, actorName, actorUsername := resolveActorForDelete(c, level)
	cmd := deleteProker.DeleteProkerV3Command{
		UUID:          parsedUuid.String(),
		SoftDelete:    softDelete,
		Catatan:       catatan,
		Role:          role,
		ActorID:       actorID,
		ActorName:     actorName,
		ActorUsername: actorUsername,
	}
	deletedUUID, err := mediatr.Send[deleteProker.DeleteProkerV3Command, string](c.Context(), cmd)
	if err != nil {
		return commoninfra.HandleError(c, err)
	}

	return commonpresentation.JsonUUID(c, deletedUUID)
}

// =======================================================
// GET /proker_v3/:uuid
// =======================================================
func GetProkerV3Handlerfunc(c *fiber.Ctx) error {
	rawUuid := c.Params("uuid")
	if rawUuid == "" {
		rawUuid = c.Query("uuid")
	}

	parsedUuid, err := uuid.Parse(strings.TrimSpace(rawUuid))
	if err != nil {
		return commoninfra.HandleError(c, domain.ProkerNotFound())
	}

	query := getProker.GetProkerV3ByUuidQuery{UUID: parsedUuid}
	res, err := mediatr.Send[getProker.GetProkerV3ByUuidQuery, *domain.ProkerV3Default](c.Context(), query)
	if err != nil {
		return commoninfra.HandleError(c, err)
	}
	if res == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Program kerja tidak ditemukan"})
	}

	level, _ := commonpresentation.GetUserInfo(c)
	if level != "admin" {
		if !canAccessProkerUnit(c, res.IdFakultasUnit, res.UuidFakultasUnit, res.NamaFakultasUnit, res.ParentFakultas) {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"code": "FORBIDDEN", "error": "Akses ditolak: bukan unit anda"})
		}
	}

	return c.JSON(res)
}

// =======================================================
// GET /proker_v3
// =======================================================
func GetAllProkerV3Handlerfunc(c *fiber.Ctx) error {
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

	tab := strings.TrimSpace(c.Query("tab"))
	if tab == "" {
		tab = strings.TrimSpace(c.Query("status_tab"))
	}
	if tab != "" && tab != "all" && tab != "semua" {
		filters = append(filters, commondomain.SearchFilter{
			Field:    "tab",
			Operator: "eq",
			Value:    &tab,
		})
	}

	target := strings.TrimSpace(c.Params("target"))
	if target == "" {
		target = strings.TrimSpace(c.Query("target"))
	}
	if target == "" {
		target = strings.TrimSpace(c.Query("id_fakultas_unit"))
	}
	if target == "" {
		filterFakultasUnit := strings.TrimSpace(c.Query("filter_fakultas_unit"))
		if filterFakultasUnit == "" {
			filterFakultasUnit = strings.TrimSpace(c.Query("fakultas_unit"))
		}
		if filterFakultasUnit != "" && filterFakultasUnit != "all" && filterFakultasUnit != "semua" {
			target = filterFakultasUnit
		}
	}
	if target == "" {
		target = strings.TrimSpace(c.Get("X-Active-Unit"))
	}

	currentRole := strings.ToLower(strings.TrimSpace(c.Get("X-Active-Role")))
	if currentRole == "" {
		currentRole = strings.ToLower(strings.TrimSpace(c.Get("X-Level")))
	}
	if currentRole == "" {
		if lvl, ok := c.Locals("level").(string); ok {
			currentRole = strings.ToLower(strings.TrimSpace(lvl))
		}
	}
	if currentRole == "" {
		currentRole = strings.ToLower(strings.TrimSpace(c.FormValue("level")))
	}

	fakultasUnit := strings.TrimSpace(c.Get("X-Active-Unit"))
	if fakultasUnit == "" {
		if fu, ok := c.Locals("fakultas_unit").(string); ok {
			fakultasUnit = fu
		}
	}
	if fakultasUnit == "" {
		fakultasUnit = strings.TrimSpace(c.FormValue("fakultas_unit"))
	}
	if fakultasUnit == "" {
		fakultasUnit = strings.TrimSpace(c.FormValue("target"))
	}

	isUnrestricted := currentRole == "admin" || currentRole == "anggaran" || currentRole == "keuangan" || currentRole == "evaluator"
	if !isUnrestricted && target == "" && fakultasUnit != "" && fakultasUnit != "0" {
		target = fakultasUnit
	}

	if target != "" && target != "all" && target != "semua" {
		field := "id_fakultas_unit"
		if strings.Contains(target, "-") {
			field = "fakultasunituuid"
		}
		filters = append(filters, commondomain.SearchFilter{
			Field:    field,
			Operator: "eq",
			Value:    &target,
		})
	}

	filterProdi := strings.TrimSpace(c.Query("filter_prodi"))
	if filterProdi == "" {
		filterProdi = strings.TrimSpace(c.Query("id_prodi"))
	}
	if filterProdi == "" {
		filterProdi = strings.TrimSpace(c.Query("prodi"))
	}
	if filterProdi != "" && filterProdi != "all" && filterProdi != "semua" {
		filters = append(filters, commondomain.SearchFilter{
			Field:    "filter_prodi",
			Operator: "eq",
			Value:    &filterProdi,
		})
	}

	if currentRole != "" {
		filters = append(filters, commondomain.SearchFilter{
			Field:    "active_role",
			Operator: "eq",
			Value:    &currentRole,
		})
	}

	filterTahun := strings.TrimSpace(c.Query("filter_tahun"))
	if filterTahun == "" {
		filterTahun = strings.TrimSpace(c.Query("tahun"))
	}
	if filterTahun == "" {
		filterTahun = strings.TrimSpace(c.Get("X-Active-Year"))
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

	query := getAllProker.GetAllProkerV3Query{
		Search:  search,
		Filters: filters,
	}

	var adapter commonpresentation.OutputAdapter[domain.ProkerV3Default]
	switch mode {
	case "all":
		adapter = &commonpresentation.AllAdapter[domain.ProkerV3Default]{}
	case "ndjson":
		adapter = &commonpresentation.NDJSONAdapter[domain.ProkerV3Default]{}
	case "sse":
		adapter = &commonpresentation.SSEAdapter[domain.ProkerV3Default]{}
	default:
		query.Page = &page
		query.Limit = &limit
		adapter = &commonpresentation.PagingAdapter[domain.ProkerV3Default]{}
	}

	result, err := mediatr.Send[
		getAllProker.GetAllProkerV3Query,
		commondomain.Paged[domain.ProkerV3Default],
	](c.Context(), query)

	if err != nil {
		return commoninfra.HandleError(c, err)
	}

	return adapter.Send(c, result)
}

// =======================================================
// PUT /proker_v3/:uuid/status
// =======================================================
func UpdateProkerStatusHandler(c *fiber.Ctx) error {
	rawUuid := c.Params("uuid")
	if rawUuid == "" {
		rawUuid = c.Query("uuid")
	}

	var cmd updateStatus.UpdateStatusV3Command
	_ = c.BodyParser(&cmd)

	if val := c.FormValue("status"); val != "" {
		cmd.Status = val
	}
	if val := c.FormValue("role"); val != "" {
		cmd.Role = val
	}
	if val := c.FormValue("catatan_penolakan"); val != "" {
		cmd.CatatanPenolakan = val
	}
	if val := c.FormValue("alasan"); val != "" {
		cmd.Alasan = val
	}
	if val := c.FormValue("catatan"); val != "" {
		cmd.Catatan = val
	}

	level, userFU := commonpresentation.GetUserInfo(c)
	cmd.UUID = strings.TrimSpace(rawUuid)
	cmd.Level = level
	cmd.UserFU = userFU

	if parsedUuid, err := uuid.Parse(strings.TrimSpace(rawUuid)); err == nil && level != "admin" {
		query := getProker.GetProkerV3ByUuidQuery{UUID: parsedUuid}
		if existing, err := mediatr.Send[getProker.GetProkerV3ByUuidQuery, *domain.ProkerV3Default](c.Context(), query); err == nil && existing != nil {
			if !canAccessProkerUnit(c, existing.IdFakultasUnit, existing.UuidFakultasUnit, existing.NamaFakultasUnit, existing.ParentFakultas) {
				return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"code": "FORBIDDEN", "error": "Akses ditolak: bukan unit penugasan anda"})
			}
		}
	}

	if u, ok := c.Locals("user").(*commonpresentation.Account); ok && u != nil {
		cmd.UserName = u.Name
		if cmd.UserName == "" {
			cmd.UserName = u.NidnUsername
		}
	}
	if uid, ok := c.Locals("user_id").(int); ok && uid > 0 {
		cmd.UserID = &uid
	}

	res, err := mediatr.Send[updateStatus.UpdateStatusV3Command, *domain.UpdateStatusResult](c.Context(), cmd)
	if err != nil {
		return commoninfra.HandleError(c, err)
	}

	return c.JSON(res)
}

// =======================================================
// GET /proker_v3/:uuid/history
// =======================================================
func GetHistoryStateHandler(c *fiber.Ctx) error {
	rawUuid := c.Params("uuid")
	if rawUuid == "" {
		rawUuid = c.Params("aktivitas_uuid")
	}
	if rawUuid == "" {
		rawUuid = c.Query("uuid")
	}

	res, err := mediatr.Send[getHistory.GetHistoryStateV3Query, []domain.HistoryState](
		c.Context(),
		getHistory.GetHistoryStateV3Query{UUID: strings.TrimSpace(rawUuid)},
	)
	if err != nil {
		return commoninfra.HandleError(c, err)
	}

	return c.JSON(fiber.Map{
		"data":  res,
		"total": len(res),
	})
}

// =======================================================
// GET /proker_v3_options
// =======================================================
func GetProkerV3OptionsHandler(c *fiber.Ctx) error {
	res, err := mediatr.Send[getOptions.GetProkerV3OptionsQuery, *domain.ProkerOptionsResult](
		c.Context(),
		getOptions.GetProkerV3OptionsQuery{},
	)
	if err != nil {
		return commoninfra.HandleError(c, err)
	}

	return c.JSON(res)
}

// =======================================================
// GET /proker_v3/coa
// =======================================================
func GetCoaHandler(c *fiber.Ctx) error {
	search := strings.TrimSpace(c.Query("search"))
	if search == "" {
		search = strings.TrimSpace(c.Query("q"))
	}

	items, err := mediatr.Send[getCoa.GetCoaQuery, []domain.CoaItem](
		c.Context(),
		getCoa.GetCoaQuery{Search: search},
	)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Gagal mengambil data akun COA: " + err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"status": "success",
		"data":   items,
		"total":  len(items),
	})
}

// =======================================================
// GET /proker_v3/target_unpak
// =======================================================
func GetTargetUnpakHandler(c *fiber.Ctx) error {
	tahun := strings.TrimSpace(c.Query("tahun"))
	if tahun == "" {
		tahun = strings.TrimSpace(c.Query("filter_tahun"))
	}
	indikator := strings.TrimSpace(c.Query("indikator"))
	if indikator == "" {
		indikator = strings.TrimSpace(c.Query("id_indikator"))
	}
	unit := strings.TrimSpace(c.Query("unit"))
	if unit == "" {
		unit = strings.TrimSpace(c.Query("id_fakultas_unit"))
	}

	res, err := mediatr.Send[getTargetUnpak.GetTargetUnpakQuery, *domain.TargetUnpakResult](
		c.Context(),
		getTargetUnpak.GetTargetUnpakQuery{
			Tahun:          tahun,
			IndikatorParam: indikator,
			UnitParam:      unit,
		},
	)
	if err != nil {
		return commoninfra.HandleError(c, err)
	}

	return c.JSON(res)
}

// =======================================================
// GET /proker_v3/:uuid/dokumen
// =======================================================
func GetProkerDokumenHandler(c *fiber.Ctx) error {
	rawUuid := c.Params("uuid")
	parsedUuid, err := uuid.Parse(rawUuid)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid matrix UUID"})
	}

	res, err := mediatr.Send[getDokumen.GetDokumenV3Query, *domain.ProkerDokumenResult](
		c.Context(),
		getDokumen.GetDokumenV3Query{ProkerUUID: parsedUuid},
	)
	if err != nil {
		return commoninfra.HandleError(c, err)
	}

	return c.JSON(res)
}

// =======================================================
// DELETE /proker_v3/dokumen/:doc_uuid
// =======================================================
func DeleteProkerDokumenHandler(c *fiber.Ctx) error {
	rawUuid := c.Params("doc_uuid")
	_, err := uuid.Parse(rawUuid)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid dokumen UUID"})
	}

	msg, err := mediatr.Send[deleteDokumen.DeleteDokumenV3Command, string](
		c.Context(),
		deleteDokumen.DeleteDokumenV3Command{DocUUID: rawUuid},
	)
	if err != nil {
		return commoninfra.HandleError(c, err)
	}

	return c.JSON(fiber.Map{"message": msg, "uuid": rawUuid})
}

// =======================================================
// PUT /proker_v3/dokumen/:doc_uuid/verifikasi
// =======================================================
func UpdateVerifikasiDokumenHandler(c *fiber.Ctx) error {
	rawUuid := c.Params("doc_uuid")
	_, err := uuid.Parse(rawUuid)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid dokumen UUID"})
	}

	var body struct {
		StatusVerifikasi string `json:"status_verifikasi"`
		Catatan          string `json:"catatan"`
		CatatanPenolakan string `json:"catatan_penolakan"`
	}
	_ = c.BodyParser(&body)

	status := strings.TrimSpace(body.StatusVerifikasi)
	if status == "" {
		status = strings.TrimSpace(c.FormValue("status_verifikasi"))
	}
	catatan := strings.TrimSpace(body.Catatan)
	if catatan == "" {
		catatan = strings.TrimSpace(c.FormValue("catatan"))
	}
	catatanPenolakan := strings.TrimSpace(body.CatatanPenolakan)
	if catatanPenolakan == "" {
		catatanPenolakan = strings.TrimSpace(c.FormValue("catatan_penolakan"))
	}

	res, err := mediatr.Send[updateVerifikasiDokumen.UpdateVerifikasiDokumenV3Command, *domain.MatrixProkerDokumen](
		c.Context(),
		updateVerifikasiDokumen.UpdateVerifikasiDokumenV3Command{
			DocUUID:          rawUuid,
			StatusVerifikasi: status,
			Catatan:          catatan,
			CatatanPenolakan: catatanPenolakan,
		},
	)
	if err != nil {
		return commoninfra.HandleError(c, err)
	}

	return c.JSON(fiber.Map{
		"message": "Status verifikasi dokumen berhasil diperbarui",
		"data":    res,
	})
}

// =======================================================
// UPLOAD DOKUMEN HANDLERS
// =======================================================

func sanitizeFolderSegment(s string) string {
	s = strings.TrimSpace(s)
	if s == "" {
		return "LAINNYA"
	}
	invalid := []string{"/", "\\", ":", "*", "?", "\"", "<", ">", "|", "\n", "\r", "\t"}
	for _, char := range invalid {
		s = strings.ReplaceAll(s, char, "_")
	}
	for strings.Contains(s, "__") {
		s = strings.ReplaceAll(s, "__", "_")
	}
	s = strings.Trim(s, " ._")
	if s == "" {
		return "LAINNYA"
	}
	if len(s) > 100 {
		s = s[:100]
	}
	return s
}

func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	if _, err = io.Copy(out, in); err != nil {
		return err
	}
	return out.Sync()
}

func UploadProkerDokumenChunkHandler(c *fiber.Ctx) error {
	rawUuid := c.Params("uuid")
	parsedUuid, err := uuid.Parse(rawUuid)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid matrix UUID"})
	}

	proker, err := mediatr.Send[getProker.GetProkerV3ByUuidQuery, *domain.ProkerV3Default](
		c.Context(),
		getProker.GetProkerV3ByUuidQuery{UUID: parsedUuid},
	)
	if err != nil {
		return commoninfra.HandleError(c, err)
	}
	if proker == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Program kerja tidak ditemukan"})
	}
	if proker.DeletedAt != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Program kerja sudah dihapus dan tidak dapat mengunggah dokumen"})
	}

	uploadIDStr := strings.TrimSpace(c.FormValue("upload_id"))
	if uploadIDStr == "" {
		uploadIDStr = strings.TrimSpace(c.Get("X-Upload-ID"))
	}
	docUUID, err := uuid.Parse(uploadIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid or missing upload_id (must be UUID)"})
	}

	chunkIndex, _ := strconv.Atoi(c.FormValue("chunk_index"))
	totalChunks, _ := strconv.Atoi(c.FormValue("total_chunks"))
	chunkOffset, _ := strconv.ParseInt(c.FormValue("chunk_offset"), 10, 64)
	totalSize, _ := strconv.ParseInt(c.FormValue("total_size"), 10, 64)
	fileName := strings.TrimSpace(c.FormValue("file_name"))
	jenisDok := strings.TrimSpace(c.FormValue("jenis_dokumen"))
	if jenisDok == "" {
		jenisDok = "DOCUMENT PENDUKUNG"
	}
	catatan := strings.TrimSpace(c.FormValue("catatan"))
	keterangan := strings.TrimSpace(c.FormValue("keterangan"))
	if catatan == "" {
		catatan = keterangan
	}
	if keterangan == "" {
		keterangan = catatan
	}

	fileHeader, err := c.FormFile("chunk")
	if err != nil {
		fileHeader, err = c.FormFile("file")
	}
	if err != nil || fileHeader == nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Tidak ada binary chunk yang dikirim"})
	}

	tmpDir := "./upload/tmp"
	if err := os.MkdirAll(tmpDir, 0755); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal menyiapkan direktori temporary: " + err.Error()})
	}

	partPath := filepath.Join(tmpDir, fmt.Sprintf("%s.part", docUUID.String()))
	flags := os.O_CREATE | os.O_WRONLY
	if chunkIndex == 0 && chunkOffset == 0 {
		flags = os.O_CREATE | os.O_WRONLY | os.O_TRUNC
	}

	partFile, err := os.OpenFile(partPath, flags, 0644)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal membuka file temporary: " + err.Error()})
	}

	if _, err := partFile.Seek(chunkOffset, io.SeekStart); err != nil {
		partFile.Close()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal memposisikan offset chunk: " + err.Error()})
	}

	src, err := fileHeader.Open()
	if err != nil {
		partFile.Close()
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Gagal membaca chunk: " + err.Error()})
	}

	written, err := io.Copy(partFile, src)
	src.Close()
	partFile.Close()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal menulis chunk ke temporary: " + err.Error()})
	}

	isFinalChunk := (totalChunks > 0 && chunkIndex >= totalChunks-1) ||
		(totalSize > 0 && (chunkOffset+written) >= totalSize)

	if isFinalChunk {
		namaUnit := sanitizeFolderSegment(proker.NamaFakultasUnit)
		tahun := sanitizeFolderSegment(proker.Tahun)
		standar := sanitizeFolderSegment(proker.StandarSubStandar)
		indikator := sanitizeFolderSegment(proker.IndikatorKinerja)
		jenisFolder := sanitizeFolderSegment(jenisDok)

		targetDir := filepath.Join("./upload/prokerv2", namaUnit, tahun, standar, indikator, jenisFolder)
		if err := os.MkdirAll(targetDir, 0755); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal membuat direktori upload: " + err.Error()})
		}

		ext := filepath.Ext(fileName)
		if ext == "" {
			ext = filepath.Ext(fileHeader.Filename)
		}
		savedFileName := fmt.Sprintf("%s%s", docUUID.String(), ext)
		destPath := filepath.Join(targetDir, savedFileName)

		if err := os.Rename(partPath, destPath); err != nil {
			if copyErr := copyFile(partPath, destPath); copyErr != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal memindahkan file dokumen: " + copyErr.Error()})
			}
			_ = os.Remove(partPath)
		}

		relURL := filepath.ToSlash(destPath)
		if !strings.HasPrefix(relURL, "/") {
			relURL = "/" + strings.TrimPrefix(relURL, "./")
		}

		doc := domain.MatrixProkerDokumen{
			UUID:             docUUID,
			IdMatrikProker:   proker.ID,
			JenisDokumen:     jenisDok,
			NamaFile:         fileName,
			File:             destPath,
			FileURL:          relURL,
			FileSize:         totalSize,
			FileType:         fileHeader.Header.Get("Content-Type"),
			StatusVerifikasi: "belum_terverifikasi",
			Catatan:          catatan,
			Keterangan:       keterangan,
		}

		createdDoc, err := mediatr.Send[createDokumen.CreateProkerDokumenCommand, *domain.MatrixProkerDokumen](
			c.Context(),
			createDokumen.CreateProkerDokumenCommand{Dokumen: &doc},
		)
		if err != nil {
			_ = os.Remove(destPath)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal menyimpan record dokumen: " + err.Error()})
		}

		return c.Status(fiber.StatusCreated).JSON(fiber.Map{
			"message":   "Dokumen berhasil diunggah secara utuh",
			"completed": true,
			"data":      createdDoc,
		})
	}

	return c.JSON(fiber.Map{
		"message":     "Chunk binary berhasil diterima",
		"completed":   false,
		"chunk_index": chunkIndex,
		"uploaded":    chunkOffset + fileHeader.Size,
	})
}

func CancelProkerDokumenChunkHandler(c *fiber.Ctx) error {
	uploadIDStr := strings.TrimSpace(c.Params("upload_id"))
	if uploadIDStr == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Missing upload_id"})
	}
	partPath := filepath.Join("./upload/tmp", fmt.Sprintf("%s.part", uploadIDStr))
	_ = os.Remove(partPath)
	return c.JSON(fiber.Map{"message": "Upload chunk dibatalkan dan file temporary dihapus"})
}

func UploadProkerDokumenHandler(c *fiber.Ctx) error {
	rawUuid := c.Params("uuid")
	parsedUuid, err := uuid.Parse(rawUuid)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid matrix UUID"})
	}

	proker, err := mediatr.Send[getProker.GetProkerV3ByUuidQuery, *domain.ProkerV3Default](
		c.Context(),
		getProker.GetProkerV3ByUuidQuery{UUID: parsedUuid},
	)
	if err != nil {
		return commoninfra.HandleError(c, err)
	}
	if proker == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Program kerja tidak ditemukan"})
	}
	if proker.DeletedAt != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Program kerja sudah dihapus dan tidak dapat mengunggah dokumen"})
	}

	form, err := c.MultipartForm()
	if err != nil || form == nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Gagal membaca form upload: " + err.Error()})
	}

	var fileHeaders []*multipart.FileHeader
	if len(form.File["file"]) > 0 {
		fileHeaders = append(fileHeaders, form.File["file"]...)
	}
	if len(form.File["files"]) > 0 {
		fileHeaders = append(fileHeaders, form.File["files"]...)
	}
	if len(form.File["files[]"]) > 0 {
		fileHeaders = append(fileHeaders, form.File["files[]"]...)
	}
	if len(fileHeaders) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Tidak ada file yang dipilih untuk diunggah"})
	}

	jenisDok := strings.TrimSpace(c.FormValue("jenis_dokumen"))
	if jenisDok == "" {
		jenisDok = "DOCUMENT PENDUKUNG"
	}
	catatan := strings.TrimSpace(c.FormValue("catatan"))
	keterangan := strings.TrimSpace(c.FormValue("keterangan"))
	if catatan == "" {
		catatan = keterangan
	}
	if keterangan == "" {
		keterangan = catatan
	}

	namaUnit := sanitizeFolderSegment(proker.NamaFakultasUnit)
	tahun := sanitizeFolderSegment(proker.Tahun)
	standar := sanitizeFolderSegment(proker.StandarSubStandar)
	indikator := sanitizeFolderSegment(proker.IndikatorKinerja)
	jenisFolder := sanitizeFolderSegment(jenisDok)

	targetDir := filepath.Join("./upload/prokerv2", namaUnit, tahun, standar, indikator, jenisFolder)
	if err := os.MkdirAll(targetDir, 0755); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal membuat direktori upload: " + err.Error()})
	}

	var createdDocs []domain.MatrixProkerDokumen
	for _, fh := range fileHeaders {
		docUUID := uuid.New()
		ext := filepath.Ext(fh.Filename)
		savedFileName := fmt.Sprintf("%s%s", docUUID.String(), ext)
		destPath := filepath.Join(targetDir, savedFileName)

		if err := c.SaveFile(fh, destPath); err != nil {
			continue
		}

		relURL := filepath.ToSlash(destPath)
		if !strings.HasPrefix(relURL, "/") {
			relURL = "/" + strings.TrimPrefix(relURL, "./")
		}

		doc := domain.MatrixProkerDokumen{
			UUID:             docUUID,
			IdMatrikProker:   proker.ID,
			JenisDokumen:     jenisDok,
			NamaFile:         fh.Filename,
			File:             destPath,
			FileURL:          relURL,
			FileSize:         fh.Size,
			FileType:         fh.Header.Get("Content-Type"),
			StatusVerifikasi: "belum_terverifikasi",
			Catatan:          catatan,
			Keterangan:       keterangan,
		}

		createdDoc, err := mediatr.Send[createDokumen.CreateProkerDokumenCommand, *domain.MatrixProkerDokumen](
			c.Context(),
			createDokumen.CreateProkerDokumenCommand{Dokumen: &doc},
		)
		if err != nil {
			_ = os.Remove(destPath)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Gagal menyimpan record dokumen: " + err.Error(),
			})
		}
		if createdDoc != nil {
			createdDocs = append(createdDocs, *createdDoc)
		}
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": fmt.Sprintf("Berhasil mengunggah %d dokumen", len(createdDocs)),
		"data":    createdDocs,
	})
}

// =======================================================
// SETTINGS PROKER BATAL
// =======================================================
func GetProkerBatalSettingHandler(c *fiber.Ctx) error {
	level, userFU := commonpresentation.GetUserInfo(c)
	isAdmin := strings.EqualFold(strings.TrimSpace(level), "admin")

	rawUnitID := strings.TrimSpace(c.Query("unit_id"))
	if rawUnitID != "" {
		if _, err := strconv.Atoi(rawUnitID); err == nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"code":    "BAD_REQUEST",
				"message": "salah format unit_id",
			})
		}
	}

	unitUUID := strings.TrimSpace(c.Query("unit_uuid"))
	if unitUUID == "" {
		unitUUID = strings.TrimSpace(c.Query("unit"))
	}
	if unitUUID == "" && rawUnitID != "" {
		unitUUID = rawUnitID
	}

	if unitUUID != "" {
		if _, err := uuid.Parse(unitUUID); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"code":    "BAD_REQUEST",
				"message": "Format unit_uuid tidak valid, harus berupa UUID",
			})
		}
	}

	var targetUnitID *int
	if !isAdmin && userFU != "" && unitUUID == "" {
		if id, err := strconv.Atoi(userFU); err == nil && id > 0 {
			targetUnitID = &id
		}
	}

	result, err := mediatr.Send[
		GetProkerBatalSetting.GetProkerBatalSettingQuery,
		*GetProkerBatalSetting.GetProkerBatalSettingResult,
	](c.Context(), GetProkerBatalSetting.GetProkerBatalSettingQuery{
		TargetUnitID:   targetUnitID,
		TargetUnitUUID: unitUUID,
		UserUnit:       userFU,
		UserLevel:      level,
	})
	if err != nil {
		return commoninfra.HandleError(c, err)
	}

	return c.JSON(result)
}

func UpdateProkerBatalSettingHandler(c *fiber.Ctx) error {
	level, _ := commonpresentation.GetUserInfo(c)
	isAdmin := strings.EqualFold(strings.TrimSpace(level), "admin")
	if !isAdmin {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Hanya admin yang dapat mengubah pengaturan pembatalan proker"})
	}

	var cmd UpdateProkerBatalSetting.UpdateProkerBatalSettingCommand
	_ = c.BodyParser(&cmd)
	cmd.IsAdmin = isAdmin

	if u, ok := c.Locals("user").(*commonpresentation.Account); ok && u != nil {
		cmd.ActorSID = u.NidnUsername
	}

	if cmd.GlobalDate == nil && c.FormValue("global_date") != "" {
		val := strings.TrimSpace(c.FormValue("global_date"))
		cmd.GlobalDate = &val
	}
	if cmd.Tanggal == nil && c.FormValue("tanggal") != "" {
		val := strings.TrimSpace(c.FormValue("tanggal"))
		cmd.Tanggal = &val
	}
	if cmd.Action == "" && c.FormValue("action") != "" {
		cmd.Action = strings.TrimSpace(c.FormValue("action"))
	}

	result, err := mediatr.Send[
		UpdateProkerBatalSetting.UpdateProkerBatalSettingCommand,
		*UpdateProkerBatalSetting.UpdateProkerBatalSettingResult,
	](c.Context(), cmd)
	if err != nil {
		return commoninfra.HandleError(c, err)
	}

	return c.JSON(result)
}

// =======================================================
// JADWAL PROKER V3 HANDLERS
// =======================================================
func GetAllJadwalProkerV3Handler(c *fiber.Ctx) error {
	filter := domain.JadwalFilter{
		Search:     strings.TrimSpace(c.Query("search")),
		Tahun:      strings.TrimSpace(c.Query("tahun")),
		FilterUnit: strings.TrimSpace(c.Query("filter_unit")),
	}

	result, err := mediatr.Send[
		getAllJadwal.GetAllJadwalProkerV3Query,
		*getAllJadwal.GetAllJadwalProkerV3Result,
	](c.Context(), getAllJadwal.GetAllJadwalProkerV3Query{Filter: filter})
	if err != nil {
		return commoninfra.HandleError(c, err)
	}

	return c.JSON(result)
}

func GetJadwalProkerV3Handler(c *fiber.Ctx) error {
	rawUuid := c.Params("uuid")
	parsedUuid, err := uuid.Parse(rawUuid)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid UUID"})
	}

	result, err := mediatr.Send[
		getJadwal.GetJadwalProkerV3Query,
		*domain.JadwalProkerV3DTO,
	](c.Context(), getJadwal.GetJadwalProkerV3Query{UUID: parsedUuid})
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Jadwal proker tidak ditemukan"})
	}

	return c.JSON(result)
}

func CheckJadwalProkerV3Handler(c *fiber.Ctx) error {
	tahun := strings.TrimSpace(c.Query("tahun"))
	unitParam := strings.TrimSpace(c.Query("id_fakultas_unit"))
	if unitParam == "" {
		unitParam = strings.TrimSpace(c.Query("fakultas_unit"))
	}

	result, err := mediatr.Send[
		checkJadwal.CheckJadwalProkerV3Query,
		*domain.CheckJadwalResult,
	](c.Context(), checkJadwal.CheckJadwalProkerV3Query{
		Request: domain.CheckJadwalRequest{
			Tahun:     tahun,
			UnitParam: unitParam,
		},
	})
	if err != nil {
		return commoninfra.HandleError(c, err)
	}

	return c.JSON(result)
}

func CreateJadwalProkerV3Handler(c *fiber.Ctx) error {
	tahun := strings.TrimSpace(c.FormValue("tahun"))
	tanggalMulai := strings.TrimSpace(c.FormValue("tanggal_mulai"))
	tanggalAkhir := strings.TrimSpace(c.FormValue("tanggal_akhir"))

	tMulai, errM := time.Parse("2006-01-02", tanggalMulai)
	if errM != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Format tanggal mulai harus YYYY-MM-DD"})
	}
	tAkhir, errA := time.Parse("2006-01-02", tanggalAkhir)
	if errA != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Format tanggal akhir harus YYYY-MM-DD"})
	}

	var fuIdPtr *uint64
	rawFu := strings.TrimSpace(c.FormValue("id_fakultas_unit"))
	if rawFu != "" && rawFu != "0" && rawFu != "all" && rawFu != "null" {
		if val, err := strconv.ParseUint(rawFu, 10, 64); err == nil && val > 0 {
			fuIdPtr = &val
		}
	}

	result, err := mediatr.Send[
		createJadwal.CreateJadwalProkerV3Command,
		*createJadwal.CreateJadwalProkerV3Result,
	](c.Context(), createJadwal.CreateJadwalProkerV3Command{
		Tahun:          tahun,
		TanggalMulai:   tMulai,
		TanggalAkhir:   tAkhir,
		IdFakultasUnit: fuIdPtr,
	})
	if err != nil {
		return commoninfra.HandleError(c, err)
	}

	return c.Status(fiber.StatusCreated).JSON(result)
}

func UpdateJadwalProkerV3Handler(c *fiber.Ctx) error {
	rawUuid := c.Params("uuid")
	parsedUuid, err := uuid.Parse(rawUuid)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid UUID"})
	}

	tahun := strings.TrimSpace(c.FormValue("tahun"))
	tanggalMulai := strings.TrimSpace(c.FormValue("tanggal_mulai"))
	tanggalAkhir := strings.TrimSpace(c.FormValue("tanggal_akhir"))

	var tMulaiPtr, tAkhirPtr *time.Time
	if tanggalMulai != "" {
		if t, err := time.Parse("2006-01-02", tanggalMulai); err == nil {
			tMulaiPtr = &t
		} else {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Format tanggal mulai harus YYYY-MM-DD"})
		}
	}
	if tanggalAkhir != "" {
		if t, err := time.Parse("2006-01-02", tanggalAkhir); err == nil {
			tAkhirPtr = &t
		} else {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Format tanggal akhir harus YYYY-MM-DD"})
		}
	}

	var fuIdPtr *uint64
	_, hasFuKey := c.Context().PostArgs().PeekBytes([]byte("id_fakultas_unit")), true
	if hasFuKey {
		rawFu := strings.TrimSpace(c.FormValue("id_fakultas_unit"))
		if rawFu != "" && rawFu != "0" && rawFu != "all" && rawFu != "null" {
			if val, err := strconv.ParseUint(rawFu, 10, 64); err == nil && val > 0 {
				fuIdPtr = &val
			}
		}
	}

	result, err := mediatr.Send[
		updateJadwal.UpdateJadwalProkerV3Command,
		*updateJadwal.UpdateJadwalProkerV3Result,
	](c.Context(), updateJadwal.UpdateJadwalProkerV3Command{
		UUID:           parsedUuid,
		Tahun:          tahun,
		TanggalMulai:   tMulaiPtr,
		TanggalAkhir:   tAkhirPtr,
		HasFuKey:       hasFuKey,
		IdFakultasUnit: fuIdPtr,
	})
	if err != nil {
		return commoninfra.HandleError(c, err)
	}

	return c.JSON(result)
}

func DeleteJadwalProkerV3Handler(c *fiber.Ctx) error {
	rawUuid := c.Params("uuid")
	parsedUuid, err := uuid.Parse(rawUuid)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid UUID"})
	}

	result, err := mediatr.Send[
		deleteJadwal.DeleteJadwalProkerV3Command,
		*deleteJadwal.DeleteJadwalProkerV3Result,
	](c.Context(), deleteJadwal.DeleteJadwalProkerV3Command{UUID: parsedUuid})
	if err != nil {
		return commoninfra.HandleError(c, err)
	}

	return c.JSON(result)
}

// =======================================================
// PAGU ANGGARAN HANDLERS
// =======================================================
func GetAllPaguAnggaranHandler(c *fiber.Ctx) error {
	tahun := c.Query("tahun")
	idFu := c.Query("id_fakultas_unit")

	result, err := mediatr.Send[
		getAllPagu.GetAllPaguAnggaranQuery,
		*getAllPagu.GetAllPaguAnggaranResult,
	](c.Context(), getAllPagu.GetAllPaguAnggaranQuery{
		Tahun:          tahun,
		IdFakultasUnit: idFu,
	})
	if err != nil {
		return commoninfra.HandleError(c, err)
	}

	return c.JSON(result)
}

func SavePaguAnggaranHandler(c *fiber.Ctx) error {
	var req domain.SavePaguAnggaranRequest
	req.Tahun = c.FormValue("tahun")
	req.FakultasUnit = c.FormValue("fakultas_unit")
	if c.FormValue("id_fakultas_unit") != "" {
		if id, err := strconv.Atoi(c.FormValue("id_fakultas_unit")); err == nil {
			req.IdFakultasUnit = &id
		}
	}
	if c.FormValue("pagu") != "" {
		if p, err := strconv.ParseFloat(c.FormValue("pagu"), 64); err == nil {
			req.Pagu = &p
		}
	}
	if c.FormValue("id") != "" {
		if id, err := strconv.ParseUint(c.FormValue("id"), 10, 64); err == nil {
			req.ID = uint(id)
		}
	}

	result, err := mediatr.Send[
		savePagu.SavePaguAnggaranCommand,
		*savePagu.SavePaguAnggaranResult,
	](c.Context(), savePagu.SavePaguAnggaranCommand{Request: req})
	if err != nil {
		return commoninfra.HandleError(c, err)
	}

	return c.JSON(result)
}

func DeletePaguAnggaranHandler(c *fiber.Ctx) error {
	id := c.Params("id")
	result, err := mediatr.Send[
		deletePagu.DeletePaguAnggaranCommand,
		*deletePagu.DeletePaguAnggaranResult,
	](c.Context(), deletePagu.DeletePaguAnggaranCommand{ID: id})
	if err != nil {
		return commoninfra.HandleError(c, err)
	}

	return c.JSON(result)
}

func GetUnitPaguSummaryHandler(c *fiber.Ctx) error {
	tahun := strings.TrimSpace(c.Query("tahun"))
	if tahun == "" || tahun == "all" {
		tahun = strings.TrimSpace(c.Get("X-Position-Year"))
	}
	if tahun == "" {
		tahun = strings.TrimSpace(c.Get("X-Tahun"))
	}
	if tahun == "" {
		tahun = "2027"
	}

	_, userFU := commonpresentation.GetUserInfo(c)
	reqFu := strings.TrimSpace(c.Query("id_fakultas_unit"))

	result, err := mediatr.Send[
		getUnitPagu.GetUnitPaguSummaryQuery,
		*domain.UnitPaguSummaryResult,
	](c.Context(), getUnitPagu.GetUnitPaguSummaryQuery{
		Tahun:  tahun,
		ReqFU:  reqFu,
		UserFU: userFU,
	})
	if err != nil {
		return commoninfra.HandleError(c, err)
	}

	return c.JSON(result)
}

// =======================================================
// DASHBOARD MONITORING & SANKEY HANDLERS
// =======================================================
func GetDashboardProkerMonitoringHandler(c *fiber.Ctx) error {
	sid := c.FormValue("sid")
	level, userFU := commonpresentation.GetUserInfo(c)

	filter := domain.DashboardMonitoringFilter{
		SID:              sid,
		Level:            level,
		UserFU:           userFU,
		Tahun:            strings.TrimSpace(c.Query("tahun")),
		FilterUnit:       strings.TrimSpace(c.Query("id_fakultas_unit")),
		FilterWadekProdi: strings.TrimSpace(c.Query("filter_wadek_prodi")),
		FilterStatus:     strings.TrimSpace(c.Query("filter_status")),
		Search:           strings.TrimSpace(c.Query("search")),
	}

	result, err := mediatr.Send[
		getDashboard.GetDashboardProkerMonitoringQuery,
		*domain.DashboardMonitoringResult,
	](c.Context(), getDashboard.GetDashboardProkerMonitoringQuery{Filter: filter})
	if err != nil {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"allowed": false,
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"allowed":      result.Meta["allowed"],
		"role_type":    result.Meta["role_type"],
		"user_name":    result.Meta["user_name"],
		"unit_options": result.Meta["unit_options"],
		"data":         result.Data,
		"summary":      result.Summary,
	})
}

func GetSankeyDataHandler(c *fiber.Ctx) error {
	tahun := strings.TrimSpace(c.Query("tahun"))

	result, err := mediatr.Send[
		getSankey.GetSankeyDataQuery,
		*domain.SankeyDataResult,
	](c.Context(), getSankey.GetSankeyDataQuery{Tahun: tahun})
	if err != nil {
		return commoninfra.HandleError(c, err)
	}

	return c.JSON(fiber.Map{
		"tahun":         result.Tahun,
		"monitoring":    result.Monitoring,
		"master_ikus":   result.MasterIkus,
		"fakultas_unit": result.FakultasUnit,
	})
}

// =======================================================
// ModuleProkerV3 route initialization
// =======================================================
func ModuleProkerV3(app *fiber.App, args ...interface{}) {
	if len(args) > 0 {
		if db, ok := args[0].(*gorm.DB); ok {
			prokerDb = db
		}
	}
	roles := []string{"admin", "auditee", "fakultas", "unit", "user", "upml", "upmf", "upms", "anggaran", "keuangan", "evaluator"}
	whoamiURL := "http://localhost:3000/whoami"

	// App Settings (Proker Batal activation date)
	adminRoles := []string{"admin"}
	app.Get("/settings/proker-batal", commonpresentation.SmartCompress(), commonpresentation.JWTMiddleware(), GetProkerBatalSettingHandler)
	app.Put("/settings/proker-batal", commonpresentation.JWTMiddleware(), commonpresentation.RBACMiddleware(adminRoles, whoamiURL), UpdateProkerBatalSettingHandler)
	app.Patch("/settings/proker-batal", commonpresentation.JWTMiddleware(), commonpresentation.RBACMiddleware(adminRoles, whoamiURL), UpdateProkerBatalSettingHandler)

	// 1. Options & Target Unpak & COA
	app.Get("/proker_v3_options", commonpresentation.SmartCompress(), commonpresentation.JWTMiddleware(), GetProkerV3OptionsHandler)
	app.Get("/proker_v3/target_unpak", commonpresentation.SmartCompress(), commonpresentation.JWTMiddleware(), GetTargetUnpakHandler)
	app.Get("/proker_v3/coa", commonpresentation.SmartCompress(), commonpresentation.JWTMiddleware(), GetCoaHandler)

	// Jadwal Proker V3 CRUD & Status Check
	app.Get("/proker_v3/jadwal", commonpresentation.SmartCompress(), commonpresentation.JWTMiddleware(), GetAllJadwalProkerV3Handler)
	app.Get("/proker_v3/jadwal_status", commonpresentation.SmartCompress(), commonpresentation.JWTMiddleware(), CheckJadwalProkerV3Handler)
	app.Get("/proker_v3/jadwal/:uuid", commonpresentation.SmartCompress(), commonpresentation.JWTMiddleware(), GetJadwalProkerV3Handler)
	app.Post("/proker_v3/jadwal", commonpresentation.JWTMiddleware(), commonpresentation.RBACMiddleware(adminRoles, whoamiURL), CreateJadwalProkerV3Handler)
	app.Put("/proker_v3/jadwal/:uuid", commonpresentation.JWTMiddleware(), commonpresentation.RBACMiddleware(adminRoles, whoamiURL), UpdateJadwalProkerV3Handler)
	app.Delete("/proker_v3/jadwal/:uuid", commonpresentation.JWTMiddleware(), commonpresentation.RBACMiddleware(adminRoles, whoamiURL), DeleteJadwalProkerV3Handler)

	// Pagu Anggaran CRUD
	app.Get("/pagu_anggaran", commonpresentation.SmartCompress(), commonpresentation.JWTMiddleware(), GetAllPaguAnggaranHandler)
	app.Post("/pagu_anggaran", commonpresentation.JWTMiddleware(), commonpresentation.RBACMiddleware(roles, whoamiURL), SavePaguAnggaranHandler)
	app.Delete("/pagu_anggaran/:id", commonpresentation.JWTMiddleware(), commonpresentation.RBACMiddleware(roles, whoamiURL), DeletePaguAnggaranHandler)

	// Dashboard Monitoring Proker, Aktivitas & RAB (Wakil Rektor 1-4 & Fakultas)
	app.Get("/proker_v3/dashboard_monitoring", commonpresentation.SmartCompress(), commonpresentation.JWTMiddleware(), GetDashboardProkerMonitoringHandler)

	// Diagram Alir Alokasi (Sankey) Data dari Database
	app.Get("/dashboard/sankey", commonpresentation.SmartCompress(), commonpresentation.JWTMiddleware(), GetSankeyDataHandler)
	app.Get("/proker_v3/sankey", commonpresentation.SmartCompress(), commonpresentation.JWTMiddleware(), GetSankeyDataHandler)

	// Summary Pagu & Realisasi Unit
	app.Get("/proker_v3/unit_pagu_summary", commonpresentation.SmartCompress(), commonpresentation.JWTMiddleware(), GetUnitPaguSummaryHandler)

	// 2. Program Kerja (matrix_proker) CRUD
	app.Post("/proker_v3", commonpresentation.JWTMiddleware(), commonpresentation.RBACMiddleware(roles, whoamiURL), CreateProkerV3Handlerfunc)
	app.Get("/proker_v3/:uuid", commonpresentation.SmartCompress(), commonpresentation.JWTMiddleware(), commonpresentation.RBACMiddleware(roles, whoamiURL), GetProkerV3Handlerfunc)
	app.Put("/proker_v3/:uuid", commonpresentation.JWTMiddleware(), commonpresentation.RBACMiddleware(roles, whoamiURL), UpdateProkerV3Handlerfunc)
	app.Put("/proker_v3/:uuid/beban_yayasan", commonpresentation.JWTMiddleware(), commonpresentation.RBACMiddleware(roles, whoamiURL), UpdateBebanYayasanHandler)
	app.Patch("/proker_v3/:uuid/beban_yayasan", commonpresentation.JWTMiddleware(), commonpresentation.RBACMiddleware(roles, whoamiURL), UpdateBebanYayasanHandler)
	app.Put("/proker_v3/:uuid/batal", commonpresentation.JWTMiddleware(), commonpresentation.RBACMiddleware(roles, whoamiURL), UpdateBatalProkerHandler)
	app.Patch("/proker_v3/:uuid/batal", commonpresentation.JWTMiddleware(), commonpresentation.RBACMiddleware(roles, whoamiURL), UpdateBatalProkerHandler)
	app.Delete("/proker_v3/:uuid", commonpresentation.JWTMiddleware(), commonpresentation.RBACMiddleware(roles, whoamiURL), DeleteProkerV3Handlerfunc)

	app.Get("/proker_v3", commonpresentation.SmartCompress(), commonpresentation.JWTMiddleware(), commonpresentation.RBACMiddleware(roles, whoamiURL), GetAllProkerV3Handlerfunc)
	app.Get("/proker_v3s", commonpresentation.SmartCompress(), commonpresentation.JWTMiddleware(), commonpresentation.RBACMiddleware(roles, whoamiURL), GetAllProkerV3Handlerfunc)

	// Status Proker & Verification
	app.Put("/proker_v3/:uuid/status", commonpresentation.JWTMiddleware(), commonpresentation.RBACMiddleware(roles, whoamiURL), UpdateProkerStatusHandler)
	app.Put("/proker_v3/:uuid/verifikasi", commonpresentation.JWTMiddleware(), commonpresentation.RBACMiddleware(roles, whoamiURL), UpdateProkerStatusHandler)

	// 5. Matrix Proker Dokumen
	app.Get("/proker_v3/:uuid/dokumen", commonpresentation.SmartCompress(), commonpresentation.JWTMiddleware(), commonpresentation.RBACMiddleware(roles, whoamiURL), GetProkerDokumenHandler)
	app.Post("/proker_v3/:uuid/dokumen/chunk", commonpresentation.JWTMiddleware(), commonpresentation.RBACMiddleware(roles, whoamiURL), UploadProkerDokumenChunkHandler)
	app.Delete("/proker_v3/:uuid/dokumen/chunk/:upload_id", commonpresentation.JWTMiddleware(), commonpresentation.RBACMiddleware(roles, whoamiURL), CancelProkerDokumenChunkHandler)
	app.Post("/proker_v3/:uuid/dokumen", commonpresentation.JWTMiddleware(), commonpresentation.RBACMiddleware(roles, whoamiURL), UploadProkerDokumenHandler)
	app.Delete("/proker_v3/dokumen/:doc_uuid", commonpresentation.JWTMiddleware(), commonpresentation.RBACMiddleware(roles, whoamiURL), DeleteProkerDokumenHandler)
	app.Put("/proker_v3/dokumen/:doc_uuid/verifikasi", commonpresentation.JWTMiddleware(), commonpresentation.RBACMiddleware(roles, whoamiURL), UpdateVerifikasiDokumenHandler)
	app.Put("/proker_v3/dokumen/:doc_uuid/status", commonpresentation.JWTMiddleware(), commonpresentation.RBACMiddleware(roles, whoamiURL), UpdateVerifikasiDokumenHandler)
	app.Put("/proker_v3/:uuid/dokumen/:doc_uuid/verifikasi", commonpresentation.JWTMiddleware(), commonpresentation.RBACMiddleware(roles, whoamiURL), UpdateVerifikasiDokumenHandler)

	// 6. History State
	app.Get("/proker_v3/:uuid/history", commonpresentation.SmartCompress(), commonpresentation.JWTMiddleware(), commonpresentation.RBACMiddleware(roles, whoamiURL), GetHistoryStateHandler)
	app.Get("/proker_v3/history/:uuid", commonpresentation.SmartCompress(), commonpresentation.JWTMiddleware(), commonpresentation.RBACMiddleware(roles, whoamiURL), GetHistoryStateHandler)
	app.Get("/proker_v3/aktivitas/:aktivitas_uuid/history", commonpresentation.SmartCompress(), commonpresentation.JWTMiddleware(), commonpresentation.RBACMiddleware(roles, whoamiURL), GetHistoryStateHandler)
}
