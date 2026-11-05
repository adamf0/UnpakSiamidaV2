package presentation

import (
	"context"
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
	createAktivitas "UnpakSiamida/modules/prokerv3/application/CreateAktivitasV3"
	createProker "UnpakSiamida/modules/prokerv3/application/CreateProkerV3"
	createTimeline "UnpakSiamida/modules/prokerv3/application/CreateTimelineV3"
	deleteAktivitas "UnpakSiamida/modules/prokerv3/application/DeleteAktivitasV3"
	deleteDokumen "UnpakSiamida/modules/prokerv3/application/DeleteDokumenV3"
	deleteProker "UnpakSiamida/modules/prokerv3/application/DeleteProkerV3"
	deleteTimelineFile "UnpakSiamida/modules/prokerv3/application/DeleteTimelineFileV3"
	deleteTimeline "UnpakSiamida/modules/prokerv3/application/DeleteTimelineV3"
	getAktivitasByMatrix "UnpakSiamida/modules/prokerv3/application/GetAktivitasByMatrix"
	getAktivitasDetail "UnpakSiamida/modules/prokerv3/application/GetAktivitasDetail"
	getAllProker "UnpakSiamida/modules/prokerv3/application/GetAllProkerV3"
	getDokumen "UnpakSiamida/modules/prokerv3/application/GetDokumenV3"
	getHistory "UnpakSiamida/modules/prokerv3/application/GetHistoryStateV3"
	getProker "UnpakSiamida/modules/prokerv3/application/GetProkerV3"
	getOptions "UnpakSiamida/modules/prokerv3/application/GetProkerV3Options"
	getTargetUnpak "UnpakSiamida/modules/prokerv3/application/GetTargetUnpak"
	reorderTimeline "UnpakSiamida/modules/prokerv3/application/ReorderTimelineV3"
	updateAktivitas "UnpakSiamida/modules/prokerv3/application/UpdateAktivitasV3"
	updateProker "UnpakSiamida/modules/prokerv3/application/UpdateProkerV3"
	updateStatus "UnpakSiamida/modules/prokerv3/application/UpdateStatusV3"
	updateTimeline "UnpakSiamida/modules/prokerv3/application/UpdateTimelineV3"
	updateVerifikasiDokumen "UnpakSiamida/modules/prokerv3/application/UpdateVerifikasiDokumenV3"
	domain "UnpakSiamida/modules/prokerv3/domain"
	infra "UnpakSiamida/modules/prokerv3/infrastructure"
)

var (
	repoProker domain.IProkerV3Repository
	dbInstance *gorm.DB
)

func canAccessProkerUnit(c *fiber.Ctx, itemUnitID *int, itemUnitUUID string, itemUnitName string) bool {
	level, userFU := commonpresentation.GetUserInfo(c)
	if level == "admin" || level == "upms" || level == "anggaran" || level == "keuangan" || level == "evaluator" {
		return true
	}

	if helper.MatchUnit(itemUnitID, itemUnitUUID, userFU) {
		return true
	}
	if itemUnitName != "" && userFU != "" && strings.EqualFold(strings.TrimSpace(itemUnitName), strings.TrimSpace(userFU)) {
		return true
	}

	if u, ok := c.Locals("user").(*commonpresentation.Account); ok && u != nil {
		activeTahun := strings.TrimSpace(c.Get("X-Position-Year"))
		if activeTahun == "" {
			activeTahun = strings.TrimSpace(c.Get("X-Tahun"))
		}
		if activeTahun == "" {
			activeTahun = strings.TrimSpace(c.Query("tahun"))
		}

		for _, r := range u.ExtraRole {
			if level != "" && !strings.EqualFold(r.Role, level) && level != "user" {
				continue
			}
			if activeTahun != "" && r.Tahun != "" && r.Tahun != activeTahun {
				continue
			}

			if itemUnitID != nil && r.IdFakultasUnit > 0 && uint(*itemUnitID) == r.IdFakultasUnit {
				return true
			}
			if itemUnitName != "" && r.NamaFakultasUnit != "" && strings.EqualFold(strings.TrimSpace(itemUnitName), strings.TrimSpace(r.NamaFakultasUnit)) {
				return true
			}
			if userFU != "" && (strconv.Itoa(int(r.IdFakultasUnit)) == userFU || strings.EqualFold(r.NamaFakultasUnit, userFU)) {
				if itemUnitID != nil && uint(*itemUnitID) == r.IdFakultasUnit {
					return true
				}
			}
		}
	}

	if level == "upml" || level == "upmf" {
		return true
	}

	return false
}

// =======================================================
// POST /proker_v3 (Simpan ke matrix_proker)
// =======================================================
func CreateProkerV3Handlerfunc(c *fiber.Ctx) error {
	level, userFU := commonpresentation.GetUserInfo(c)

	var cmd createProker.CreateProkerV3Command
	_ = c.BodyParser(&cmd)

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
// PUT /proker_v3/:uuid (Update matrix_proker)
// =======================================================
func UpdateProkerV3Handlerfunc(c *fiber.Ctx) error {
	rawUuid := c.Params("uuid")
	if rawUuid == "" {
		rawUuid = c.Query("uuid")
	}

	var cmd updateProker.UpdateProkerV3Command
	_ = c.BodyParser(&cmd)

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
		by := val == "1" || strings.ToLower(val) == "true"
		cmd.BebanYayasan = &by
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
				if !canAccessProkerUnit(c, existing.IdFakultasUnit, existing.UuidFakultasUnit, existing.NamaFakultasUnit) {
					return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Akses ditolak: bukan unit anda"})
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
// PUT /proker_v3/:uuid/beban_yayasan (Toggle/Update Beban Yayasan)
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

	repo := repoProker
	if repo == nil && dbInstance != nil {
		repo = infra.NewProkerV3Repository(dbInstance)
	}
	proker, err := repo.GetByUuid(c.Context(), parsedUuid)
	if err != nil || proker == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Program kerja tidak ditemukan"})
	}

	level, _ := commonpresentation.GetUserInfo(c)
	if level != "admin" {
		query := getProker.GetProkerV3ByUuidQuery{UUID: parsedUuid}
		existing, err := mediatr.Send[getProker.GetProkerV3ByUuidQuery, *domain.ProkerV3Default](c.Context(), query)
		if err == nil && existing != nil {
			if !canAccessProkerUnit(c, existing.IdFakultasUnit, existing.UuidFakultasUnit, existing.NamaFakultasUnit) {
				return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Akses ditolak: bukan unit anda"})
			}
		}
	}

	var reqBody struct {
		BebanYayasan *bool `json:"beban_yayasan"`
	}
	_ = c.BodyParser(&reqBody)

	var newBeban bool
	if reqBody.BebanYayasan != nil {
		newBeban = *reqBody.BebanYayasan
	} else if val := c.FormValue("beban_yayasan"); val != "" {
		newBeban = val == "1" || strings.ToLower(val) == "true"
	} else {
		newBeban = !proker.BebanYayasan
	}

	proker.BebanYayasan = newBeban
	if err := repo.Update(c.Context(), proker); err != nil {
		return commoninfra.HandleError(c, err)
	}

	return c.JSON(fiber.Map{
		"message":       "Status beban yayasan berhasil diperbarui",
		"uuid":          proker.UUID.String(),
		"beban_yayasan": proker.BebanYayasan,
	})
}

// =======================================================
// PUT /proker_v3/:uuid/batal (Toggle Batal Proker)
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

	repo := repoProker
	if repo == nil && dbInstance != nil {
		repo = infra.NewProkerV3Repository(dbInstance)
	}
	proker, err := repo.GetByUuid(c.Context(), parsedUuid)
	if err != nil || proker == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Program kerja tidak ditemukan"})
	}

	level, _ := commonpresentation.GetUserInfo(c)
	if level != "admin" {
		query := getProker.GetProkerV3ByUuidQuery{UUID: parsedUuid}
		existing, err := mediatr.Send[getProker.GetProkerV3ByUuidQuery, *domain.ProkerV3Default](c.Context(), query)
		if err == nil && existing != nil {
			if !canAccessProkerUnit(c, existing.IdFakultasUnit, existing.UuidFakultasUnit, existing.NamaFakultasUnit) {
				return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Akses ditolak: bukan unit anda"})
			}
		}
	}

	var reqBody struct {
		IsBatal      *bool  `json:"is_batal"`
		CatatanBatal string `json:"catatan_batal"`
	}
	_ = c.BodyParser(&reqBody)

	var newBatal bool
	if reqBody.IsBatal != nil {
		newBatal = *reqBody.IsBatal
	} else if val := c.FormValue("is_batal"); val != "" {
		newBatal = val == "1" || strings.ToLower(val) == "true"
	} else {
		newBatal = !proker.IsBatal
	}

	catatanBatal := strings.TrimSpace(reqBody.CatatanBatal)
	if catatanBatal == "" {
		catatanBatal = strings.TrimSpace(c.FormValue("catatan_batal"))
	}

	proker.IsBatal = newBatal
	proker.CatatanBatal = catatanBatal
	proker.UpdatedAt = time.Now()

	if err := dbInstance.WithContext(c.Context()).Model(&domain.ProkerV3{}).Where("id = ?", proker.ID).Updates(map[string]interface{}{
		"is_batal":      newBatal,
		"catatan_batal": catatanBatal,
		"updated_at":    time.Now(),
	}).Error; err != nil {
		return commoninfra.HandleError(c, err)
	}

	// Sinkronisasi status batal ke seluruh aktivitas di bawah proker ini
	actUpdates := map[string]interface{}{
		"is_batal":   newBatal,
		"updated_at": time.Now(),
	}
	if newBatal {
		actUpdates["catatan_batal"] = catatanBatal
	} else {
		actUpdates["catatan_batal"] = ""
	}
	_ = dbInstance.WithContext(c.Context()).Model(&domain.AktivitasV2{}).Where("id_matrik_proker = ?", proker.ID).Updates(actUpdates).Error

	return c.JSON(fiber.Map{
		"message":       "Status pembatalan proker berhasil diperbarui",
		"uuid":          proker.UUID.String(),
		"is_batal":      proker.IsBatal,
		"catatan_batal": proker.CatatanBatal,
	})
}

// =======================================================
// PUT /proker_v3/aktivitas/:aktivitas_uuid/batal (Toggle Batal Aktivitas)
// =======================================================
func UpdateBatalAktivitasHandler(c *fiber.Ctx) error {
	rawActUuid := c.Params("aktivitas_uuid")
	if rawActUuid == "" {
		rawActUuid = c.Params("uuid")
	}
	if rawActUuid == "" {
		rawActUuid = c.Query("aktivitas_uuid")
	}
	if rawActUuid == "" {
		rawActUuid = c.Query("uuid")
	}
	parsedUuid, err := uuid.Parse(rawActUuid)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid aktivitas UUID"})
	}

	repo := repoProker
	if repo == nil && dbInstance != nil {
		repo = infra.NewProkerV3Repository(dbInstance)
	}
	act, err := repo.GetAktivitasEntityByUuid(c.Context(), parsedUuid)
	if err != nil || act == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Aktivitas tidak ditemukan"})
	}

	level, _ := commonpresentation.GetUserInfo(c)
	if level != "admin" {
		var fuID *int
		if act.IdFakultasUnit != nil {
			id := int(*act.IdFakultasUnit)
			fuID = &id
		} else if act.IdMatrikProker != nil {
			var mp domain.ProkerV3
			if err := dbInstance.WithContext(c.Context()).Select("id_fakultas_unit").Where("id = ?", *act.IdMatrikProker).First(&mp).Error; err == nil {
				fuID = mp.IdFakultasUnit
			}
		}
		if !canAccessProkerUnit(c, fuID, "", "") {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Akses ditolak: bukan unit anda"})
		}
	}

	var reqBody struct {
		IsBatal      *bool  `json:"is_batal"`
		CatatanBatal string `json:"catatan_batal"`
	}
	_ = c.BodyParser(&reqBody)

	var newBatal bool
	if reqBody.IsBatal != nil {
		newBatal = *reqBody.IsBatal
	} else if val := c.FormValue("is_batal"); val != "" {
		newBatal = val == "1" || strings.ToLower(val) == "true"
	} else {
		newBatal = !act.IsBatal
	}

	catatanBatal := strings.TrimSpace(reqBody.CatatanBatal)
	if catatanBatal == "" {
		catatanBatal = strings.TrimSpace(c.FormValue("catatan_batal"))
	}

	act.IsBatal = newBatal
	act.CatatanBatal = catatanBatal
	act.UpdatedAt = time.Now()

	if err := dbInstance.WithContext(c.Context()).Model(&domain.AktivitasV2{}).Where("id = ?", act.ID).Updates(map[string]interface{}{
		"is_batal":      newBatal,
		"catatan_batal": catatanBatal,
		"updated_at":    time.Now(),
	}).Error; err != nil {
		return commoninfra.HandleError(c, err)
	}

	return c.JSON(fiber.Map{
		"message":       "Status pembatalan aktivitas berhasil diperbarui",
		"uuid":          act.UUID.String(),
		"is_batal":      act.IsBatal,
		"catatan_batal": act.CatatanBatal,
	})
}

// =======================================================
// DELETE /proker_v3/:uuid (Delete matrix_proker & children)
// =======================================================
func DeleteProkerV3Handlerfunc(c *fiber.Ctx) error {
	rawUuid := c.Params("uuid")
	if rawUuid == "" {
		rawUuid = c.Query("uuid")
	}

	level, _ := commonpresentation.GetUserInfo(c)
	if level != "admin" {
		parsedUuid, err := uuid.Parse(rawUuid)
		if err == nil {
			query := getProker.GetProkerV3ByUuidQuery{UUID: parsedUuid}
			existing, err := mediatr.Send[getProker.GetProkerV3ByUuidQuery, *domain.ProkerV3Default](c.Context(), query)
			if err == nil && existing != nil {
				if !canAccessProkerUnit(c, existing.IdFakultasUnit, existing.UuidFakultasUnit, existing.NamaFakultasUnit) {
					return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Akses ditolak: bukan unit anda"})
				}
			}
		}
	}

	cmd := deleteProker.DeleteProkerV3Command{UUID: rawUuid}
	deletedUUID, err := mediatr.Send[deleteProker.DeleteProkerV3Command, string](c.Context(), cmd)
	if err != nil {
		return commoninfra.HandleError(c, err)
	}

	return commonpresentation.JsonUUID(c, deletedUUID)
}

// =======================================================
// GET /proker_v3/:uuid (Get matrix_proker detail)
// =======================================================
func GetProkerV3Handlerfunc(c *fiber.Ctx) error {
	rawUuid := c.Params("uuid")
	parsedUuid, err := uuid.Parse(rawUuid)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid UUID"})
	}

	level, _ := commonpresentation.GetUserInfo(c)

	query := getProker.GetProkerV3ByUuidQuery{UUID: parsedUuid}
	result, err := mediatr.Send[getProker.GetProkerV3ByUuidQuery, *domain.ProkerV3Default](c.Context(), query)
	if err != nil {
		return commoninfra.HandleError(c, err)
	}
	if result == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Program kerja tidak ditemukan"})
	}

	isVerifier := level == "upml" || level == "upmf" || level == "upms" || level == "admin"
	if !isVerifier && !canAccessProkerUnit(c, result.IdFakultasUnit, result.UuidFakultasUnit, result.NamaFakultasUnit) {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Akses ditolak: data milik unit lain"})
	}

	return c.JSON(result)
}

// =======================================================
// GET /proker_v3s & /proker_v3 (Tabel daftar matrix_proker)
// =======================================================
func GetAllProkerV3Handlerfunc(c *fiber.Ctx) error {
	search := c.Query("search")
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	mode := c.Query("mode", "paging")

	level, userFU := commonpresentation.GetUserInfo(c)

	var filters []commondomain.SearchFilter

	isVerifierRole := level == "upml" || level == "upmf" || level == "upms" || level == "admin" || level == "anggaran" || level == "keuangan" || level == "evaluator"
	var targetFU string
	if isVerifierRole {
		targetFU = strings.TrimSpace(c.Query("filter_fakultas_unit"))
		if targetFU == "" || targetFU == "all" || targetFU == "semua" {
			targetFU = strings.TrimSpace(c.Query("fakultas_unit"))
		}
		if targetFU == "" || targetFU == "all" || targetFU == "semua" {
			targetFU = strings.TrimSpace(c.Query("id_fakultas_unit"))
		}
		if targetFU == "" || targetFU == "all" || targetFU == "semua" {
			targetFU = strings.TrimSpace(c.Params("target"))
		}
		if (targetFU == "" || targetFU == "all" || targetFU == "semua") && userFU != "" && userFU != "0" && userFU != "all" && level != "admin" && level != "upms" && level != "anggaran" && level != "keuangan" && level != "evaluator" {
			targetFU = userFU
		}
	} else {
		if qFU := strings.TrimSpace(c.Query("filter_fakultas_unit")); qFU != "" && qFU != "all" && qFU != "semua" {
			targetFU = qFU
		} else if userFU != "" && userFU != "0" && userFU != "all" {
			targetFU = userFU
		}
	}

	if targetFU != "" && targetFU != "all" && targetFU != "semua" && targetFU != "0" {
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

	filterProdi := strings.TrimSpace(c.Query("filter_prodi"))
	if filterProdi != "" && filterProdi != "all" && filterProdi != "semua" && filterProdi != "0" {
		filters = append(filters, commondomain.SearchFilter{
			Field:    "filter_prodi",
			Operator: "eq",
			Value:    &filterProdi,
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

	filterTab := strings.TrimSpace(c.Query("filter_tab"))
	if filterTab == "" {
		filterTab = strings.TrimSpace(c.Query("tab"))
	}
	if filterTab != "" && filterTab != "all" && filterTab != "semua" {
		tabVal := filterTab
		if !strings.Contains(filterTab, ":") {
			roleForTab := level
			if roleForTab == "keuangan" {
				roleForTab = "anggaran"
			}
			tabVal = fmt.Sprintf("%s:%s", roleForTab, filterTab)
		}
		filters = append(filters, commondomain.SearchFilter{
			Field:    "filter_tab",
			Operator: "eq",
			Value:    &tabVal,
		})
	}

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

			if val != "" && val != "null" && val != "all" {
				valCopy := val
				filters = append(filters, commondomain.SearchFilter{
					Field:    fld,
					Operator: op,
					Value:    &valCopy,
				})
			}
		}
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

	result, err := mediatr.Send[getAllProker.GetAllProkerV3Query, commondomain.Paged[domain.ProkerV3Default]](c.Context(), query)
	if err != nil {
		return commoninfra.HandleError(c, err)
	}

	return adapter.Send(c, result)
}

// =======================================================
// GET /proker_v3/:uuid/aktivitas (Daftar aktivitasv2)
// =======================================================
func GetAktivitasByMatrixHandler(c *fiber.Ctx) error {
	rawUuid := c.Params("uuid")
	parsedUuid, err := uuid.Parse(rawUuid)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid matrix UUID"})
	}

	result, err := mediatr.Send[getAktivitasByMatrix.GetAktivitasByMatrixQuery, *domain.AktivitasByMatrixResult](c.Context(), getAktivitasByMatrix.GetAktivitasByMatrixQuery{UUID: parsedUuid})
	if err != nil {
		return commoninfra.HandleError(c, err)
	}
	if result == nil || result.Proker == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Program kerja tidak ditemukan"})
	}

	level, _ := commonpresentation.GetUserInfo(c)
	isVerifier := level == "upml" || level == "upmf" || level == "upms" || level == "admin" || level == "evaluator" || level == "anggaran" || level == "keuangan"
	if !isVerifier && !canAccessProkerUnit(c, result.Proker.IdFakultasUnit, result.Proker.UuidFakultasUnit, result.Proker.NamaFakultasUnit) {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Akses ditolak: bukan unit anda"})
	}

	return c.JSON(fiber.Map{
		"proker": result.Proker,
		"data":   result.Data,
		"total":  result.Total,
	})
}

// =======================================================
// GET /proker_v3/aktivitas/:aktivitas_uuid (Detail aktivitas + timelines + files)
// =======================================================
func GetAktivitasDetailHandler(c *fiber.Ctx) error {
	rawUuid := c.Params("aktivitas_uuid")
	if rawUuid == "" {
		rawUuid = c.Params("uuid")
	}
	parsedUuid, err := uuid.Parse(rawUuid)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid aktivitas UUID"})
	}

	result, err := mediatr.Send[getAktivitasDetail.GetAktivitasDetailQuery, *domain.AktivitasDetailResult](c.Context(), getAktivitasDetail.GetAktivitasDetailQuery{UUID: parsedUuid})
	if err != nil {
		return commoninfra.HandleError(c, err)
	}

	return c.JSON(fiber.Map{
		"aktivitas": result.Aktivitas,
		"proker":    result.Proker,
		"timelines": result.Timelines,
		"data":      result.Data,
	})
}

// =======================================================
// POST /proker_v3/:uuid/aktivitas (Tambah aktivitasv2)
// =======================================================
func CreateAktivitasHandler(c *fiber.Ctx) error {
	rawUuid := c.Params("uuid")

	activeRole := strings.ToLower(strings.TrimSpace(c.Get("X-Active-Role")))
	if activeRole == "" {
		activeRole = strings.ToLower(strings.TrimSpace(c.Get("X-Level")))
	}
	level, _ := commonpresentation.GetUserInfo(c)
	if activeRole == "" {
		activeRole = strings.ToLower(level)
	}

	if activeRole != "user" && level != "admin" {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Akses ditolak: Hanya level user yang dapat menambah aktivitas",
		})
	}

	if level != "admin" {
		parsedUuid, err := uuid.Parse(rawUuid)
		if err == nil {
			query := getProker.GetProkerV3ByUuidQuery{UUID: parsedUuid}
			existing, err := mediatr.Send[getProker.GetProkerV3ByUuidQuery, *domain.ProkerV3Default](c.Context(), query)
			if err == nil && existing != nil {
				if !canAccessProkerUnit(c, existing.IdFakultasUnit, existing.UuidFakultasUnit, existing.NamaFakultasUnit) {
					return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
						"error": "Akses ditolak: Program kerja bukan milik unit yang ditugaskan",
					})
				}
			}
		}
	}

	var cmd createAktivitas.CreateAktivitasV3Command

	if val := c.FormValue("aktivitas"); val != "" {
		cmd.Aktivitas = strings.TrimSpace(val)
	}
	if val := c.FormValue("pic"); val != "" {
		cmd.PIC = strings.TrimSpace(val)
	}
	if val := c.FormValue("anggaran"); val != "" {
		if ang, err := strconv.ParseFloat(val, 64); err == nil {
			cmd.Anggaran = ang
		}
	}
	if val := c.FormValue("status"); val != "" {
		cmd.Status = strings.TrimSpace(val)
	}
	if val := c.FormValue("capaian"); val != "" {
		cmd.Capaian = strings.TrimSpace(val)
	}
	if val := c.FormValue("masalah"); val != "" {
		cmd.Masalah = strings.TrimSpace(val)
	}
	if val := c.FormValue("pemecahan_masalah"); val != "" {
		cmd.PemecahanMasalah = strings.TrimSpace(val)
	}
	if val := c.FormValue("rab"); val != "" {
		cmd.RAB = strings.TrimSpace(val)
	}
	cmd.MatrixUUID = rawUuid

	// Validasi pembatasan anggaran jika pagu anggaran unit untuk tahun tersebut tersedia (Case 2)
	var proposedAnggaran float64 = cmd.Anggaran
	if cmd.RAB != "" {
		var rabItems []struct {
			Total       float64 `json:"total"`
			QtySatuan   float64 `json:"qty_satuan"`
			QtyPaket    float64 `json:"qty_paket"`
			Sbu         float64 `json:"sbu"`
			HargaSatuan float64 `json:"harga_satuan"`
		}
		if errR := json.Unmarshal([]byte(cmd.RAB), &rabItems); errR == nil && len(rabItems) > 0 {
			var sumRab float64
			for _, item := range rabItems {
				if item.Total > 0 {
					sumRab += item.Total
				} else {
					sbuVal := item.Sbu
					if sbuVal == 0 {
						sbuVal = item.HargaSatuan
					}
					sumRab += item.QtySatuan * item.QtyPaket * sbuVal
				}
			}
			if sumRab > 0 {
				proposedAnggaran = sumRab
			}
		}
	}

	if proposedAnggaran > 0 && repoProker != nil {
		if prokerUid, errU := uuid.Parse(rawUuid); errU == nil {
			prokerDef, _ := repoProker.GetDefaultByUuid(c.Context(), prokerUid)
			if prokerDef != nil && prokerDef.HasPagu && prokerDef.PaguUnit > 0 {
				if prokerDef.TotalAnggaranUnit+proposedAnggaran > prokerDef.PaguUnit {
					sisa := prokerDef.PaguUnit - prokerDef.TotalAnggaranUnit
					if sisa < 0 {
						sisa = 0
					}
					return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
						"error": fmt.Sprintf("Total anggaran melebihi batas pagu anggaran unit tahun %s. Pagu unit: Rp %.0f, Total terpakai: Rp %.0f, Maksimal yang diizinkan untuk aktivitas ini: Rp %.0f", prokerDef.Tahun, prokerDef.PaguUnit, prokerDef.TotalAnggaranUnit, sisa),
					})
				}
			}
		}
	}

	act, err := mediatr.Send[createAktivitas.CreateAktivitasV3Command, *domain.AktivitasV2](c.Context(), cmd)
	if err != nil {
		return commoninfra.HandleError(c, err)
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Aktivitas berhasil ditambahkan",
		"uuid":    act.UUID.String(),
		"data":    act,
	})
}

// =======================================================
// PUT /proker_v3/aktivitas/:aktivitas_uuid (Update aktivitasv2)
// =======================================================
func UpdateAktivitasHandler(c *fiber.Ctx) error {
	rawUuid := c.Params("aktivitas_uuid")
	var cmd updateAktivitas.UpdateAktivitasV3Command

	var reqBody struct {
		Aktivitas        string   `json:"aktivitas"`
		PIC              string   `json:"pic"`
		Anggaran         *float64 `json:"anggaran"`
		RAB              string   `json:"rab"`
		Realisasi        *float64 `json:"realisasi"`
		Pencairan        *float64 `json:"pencairan"`
		Status           string   `json:"status"`
		Capaian          string   `json:"capaian"`
		Masalah          string   `json:"masalah"`
		PemecahanMasalah string   `json:"pemecahan_masalah"`
	}
	if err := c.BodyParser(&reqBody); err == nil {
		if reqBody.Aktivitas != "" {
			cmd.Aktivitas = strings.TrimSpace(reqBody.Aktivitas)
		}
		if reqBody.PIC != "" {
			cmd.PIC = strings.TrimSpace(reqBody.PIC)
		}
		if reqBody.Anggaran != nil && *reqBody.Anggaran >= 0 {
			cmd.Anggaran = *reqBody.Anggaran
		}
		if reqBody.RAB != "" {
			cmd.RAB = strings.TrimSpace(reqBody.RAB)
		}
		if reqBody.Realisasi != nil {
			cmd.Realisasi = reqBody.Realisasi
		}
		if reqBody.Pencairan != nil {
			cmd.Pencairan = reqBody.Pencairan
		}
		if reqBody.Status != "" {
			cmd.Status = strings.TrimSpace(reqBody.Status)
		}
		if reqBody.Capaian != "" {
			cmd.Capaian = strings.TrimSpace(reqBody.Capaian)
		}
		if reqBody.Masalah != "" {
			cmd.Masalah = strings.TrimSpace(reqBody.Masalah)
		}
		if reqBody.PemecahanMasalah != "" {
			cmd.PemecahanMasalah = strings.TrimSpace(reqBody.PemecahanMasalah)
		}
	}

	if val := c.FormValue("aktivitas"); val != "" {
		cmd.Aktivitas = strings.TrimSpace(val)
	}
	if val := c.FormValue("pic"); val != "" {
		cmd.PIC = strings.TrimSpace(val)
	}
	if val := c.FormValue("anggaran"); val != "" {
		if ang, err := strconv.ParseFloat(val, 64); err == nil {
			cmd.Anggaran = ang
		}
	}
	if val := c.FormValue("rab"); val != "" {
		cmd.RAB = strings.TrimSpace(val)
	}
	if val := c.FormValue("realisasi"); val != "" {
		if r, err := strconv.ParseFloat(val, 64); err == nil {
			cmd.Realisasi = &r
		}
	}
	if val := c.FormValue("pencairan"); val != "" {
		if p, err := strconv.ParseFloat(val, 64); err == nil {
			cmd.Pencairan = &p
		}
	}
	if val := c.FormValue("status"); val != "" {
		cmd.Status = strings.TrimSpace(val)
	}
	if val := c.FormValue("capaian"); val != "" {
		cmd.Capaian = strings.TrimSpace(val)
	}
	if val := c.FormValue("masalah"); val != "" {
		cmd.Masalah = strings.TrimSpace(val)
	}
	if val := c.FormValue("pemecahan_masalah"); val != "" {
		cmd.PemecahanMasalah = strings.TrimSpace(val)
	}
	cmd.UUID = rawUuid

	// Validasi pembatasan anggaran jika pagu anggaran unit untuk tahun tersebut tersedia (Case 2)
	var proposedAnggaran float64 = cmd.Anggaran
	if cmd.RAB != "" {
		var rabItems []struct {
			Total       float64 `json:"total"`
			QtySatuan   float64 `json:"qty_satuan"`
			QtyPaket    float64 `json:"qty_paket"`
			Sbu         float64 `json:"sbu"`
			HargaSatuan float64 `json:"harga_satuan"`
		}
		if errR := json.Unmarshal([]byte(cmd.RAB), &rabItems); errR == nil && len(rabItems) > 0 {
			var sumRab float64
			for _, item := range rabItems {
				if item.Total > 0 {
					sumRab += item.Total
				} else {
					sbuVal := item.Sbu
					if sbuVal == 0 {
						sbuVal = item.HargaSatuan
					}
					sumRab += item.QtySatuan * item.QtyPaket * sbuVal
				}
			}
			if sumRab > 0 {
				proposedAnggaran = sumRab
			}
		}
	}

	if proposedAnggaran > 0 && dbInstance != nil && repoProker != nil {
		if actUid, errA := uuid.Parse(rawUuid); errA == nil {
			var currentAct domain.AktivitasV2
			if errFind := dbInstance.WithContext(c.Context()).Where("uuid = ?", actUid).First(&currentAct).Error; errFind == nil && currentAct.IdMatrikProker != nil {
				var proker domain.ProkerV3
				if errP := dbInstance.WithContext(c.Context()).Where("id = ?", *currentAct.IdMatrikProker).First(&proker).Error; errP == nil {
					prokerDef, _ := repoProker.GetDefaultByUuid(c.Context(), proker.UUID)
					if prokerDef != nil && prokerDef.HasPagu && prokerDef.PaguUnit > 0 {
						oldActAnggaran := currentAct.Anggaran
						delta := proposedAnggaran - oldActAnggaran
						if delta > 0 && prokerDef.TotalAnggaranUnit+delta > prokerDef.PaguUnit {
							sisaAllowed := prokerDef.PaguUnit - (prokerDef.TotalAnggaranUnit - oldActAnggaran)
							if sisaAllowed < 0 {
								sisaAllowed = 0
							}
							return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
								"error": fmt.Sprintf("Total alokasi anggaran melebihi batas pagu anggaran unit tahun %s. Pagu unit: Rp %.0f, Total terpakai: Rp %.0f, Maksimal alokasi untuk aktivitas ini: Rp %.0f", prokerDef.Tahun, prokerDef.PaguUnit, prokerDef.TotalAnggaranUnit, sisaAllowed),
							})
						}
					}
				}
			}
		}
	}

	act, err := mediatr.Send[updateAktivitas.UpdateAktivitasV3Command, *domain.AktivitasV2](c.Context(), cmd)
	if err != nil {
		return commoninfra.HandleError(c, err)
	}

	return c.JSON(fiber.Map{
		"message": "Aktivitas berhasil diperbarui",
		"uuid":    rawUuid,
		"data":    act,
	})
}

// =======================================================
// DELETE /proker_v3/aktivitas/:aktivitas_uuid (Hapus aktivitasv2)
// =======================================================
func DeleteAktivitasHandler(c *fiber.Ctx) error {
	rawUuid := c.Params("aktivitas_uuid")

	level := strings.ToLower(strings.TrimSpace(c.Get("X-Level")))
	if level == "" {
		if u, ok := c.Locals("user").(*commonpresentation.Account); ok && u != nil {
			level = strings.ToLower(strings.TrimSpace(u.Level))
		}
	}
	if level != "admin" {
		if uid, errP := uuid.Parse(rawUuid); errP == nil && dbInstance != nil {
			var act domain.AktivitasV2
			if errAct := dbInstance.WithContext(c.Context()).Where("uuid = ?", uid).First(&act).Error; errAct == nil {
				sVerif := strings.ToLower(strings.TrimSpace(act.StatusVerifikator))
				sMain := strings.ToLower(strings.TrimSpace(act.Status))
				sLpm := strings.ToLower(strings.TrimSpace(act.StatusLpm))

				isApproved := sVerif == "terima_verifikator" ||
					sVerif == "disetujui" ||
					strings.Contains(sVerif, "diverifikasi") ||
					sLpm == "disetujui" ||
					sMain == "disetujui" ||
					sMain == "terima_anggaran"

				if !isApproved && act.IdMatrikProker != nil {
					var pMp domain.ProkerV3
					if errP := dbInstance.WithContext(c.Context()).Where("id = ?", *act.IdMatrikProker).First(&pMp).Error; errP == nil {
						pVerif := strings.ToLower(strings.TrimSpace(pMp.StatusVerifikator))
						pMain := strings.ToLower(strings.TrimSpace(pMp.Status))

						isParentApproved := !strings.Contains(pMain, "draft") &&
							!strings.Contains(pMain, "revisi") &&
							!strings.Contains(pMain, "ditolak") &&
							!strings.Contains(pVerif, "draft") &&
							!strings.Contains(pVerif, "revisi") &&
							!strings.Contains(pVerif, "ditolak") &&
							(pVerif == "terima_verifikator" ||
								pVerif == "disetujui" ||
								strings.Contains(pVerif, "diverifikasi") ||
								pMain == "disetujui" ||
								pMain == "terima_anggaran")

						if isParentApproved {
							isApproved = true
						}
					}
				}

				if isApproved {
					return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
						"error": "Aktivitas tidak dapat dihapus karena sudah disetujui oleh evaluator",
					})
				}
			}
		}
	}

	cmd := deleteAktivitas.DeleteAktivitasV3Command{UUID: rawUuid}
	deletedUUID, err := mediatr.Send[deleteAktivitas.DeleteAktivitasV3Command, string](c.Context(), cmd)
	if err != nil {
		return commoninfra.HandleError(c, err)
	}

	return c.JSON(fiber.Map{
		"message": "Aktivitas berhasil dihapus",
		"uuid":    deletedUUID,
	})
}

// =======================================================
// PUT /proker_v3/aktivitas/:aktivitas_uuid/status (Ubah status)
// =======================================================
func UpdateAktivitasStatusHandler(c *fiber.Ctx) error {
	rawUuid := c.Params("aktivitas_uuid")
	if rawUuid == "" {
		rawUuid = c.Params("uuid")
	}
	if rawUuid == "" {
		rawUuid = c.Query("uuid")
	}

	var cmd updateStatus.UpdateStatusV3Command

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
	cmd.UUID = rawUuid
	cmd.Level = level
	cmd.UserFU = userFU

	if u, ok := c.Locals("user").(*commonpresentation.Account); ok && u != nil {
		cmd.UserName = u.Name
		if cmd.UserName == "" {
			cmd.UserName = u.NidnUsername
		}
	}

	res, err := mediatr.Send[updateStatus.UpdateStatusV3Command, *domain.UpdateStatusResult](c.Context(), cmd)
	if err != nil {
		return commoninfra.HandleError(c, err)
	}

	return c.JSON(res)
}

// =======================================================
// GET /proker_v3/:uuid/history (Vertical Timeline Step History)
// =======================================================
func GetHistoryStateHandler(c *fiber.Ctx) error {
	rawUuid := c.Params("uuid")
	if rawUuid == "" {
		rawUuid = c.Params("aktivitas_uuid")
	}
	if rawUuid == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "UUID is required"})
	}

	history, err := mediatr.Send[getHistory.GetHistoryStateV3Query, []domain.HistoryState](c.Context(), getHistory.GetHistoryStateV3Query{UUID: rawUuid})
	if err != nil {
		return commoninfra.HandleError(c, err)
	}

	return c.JSON(fiber.Map{
		"uuid": rawUuid,
		"data": history,
	})
}

// =======================================================
// GET /proker_v3_options
// =======================================================
func GetProkerV3OptionsHandler(c *fiber.Ctx) error {
	result, err := mediatr.Send[getOptions.GetProkerV3OptionsQuery, *domain.ProkerOptionsResult](c.Context(), getOptions.GetProkerV3OptionsQuery{})
	if err != nil {
		return commoninfra.HandleError(c, err)
	}
	return c.JSON(result)
}

// =======================================================
// GET /proker_v3/coa (Daftar Master COA Akun dari yps_siakun.m_coa)
// =======================================================
func GetCoaHandler(c *fiber.Ctx) error {
	search := strings.TrimSpace(c.Query("search"))
	if search == "" {
		search = strings.TrimSpace(c.Query("q"))
	}

	query := dbInstance.Table("yps_siakun.m_coa").
		Select("DISTINCT TRIM(no_akun) as no_akun, TRIM(nama_akun_id) as nama_akun").
		Where("LENGTH(TRIM(no_akun)) >= 6").
		Where("nama_akun_id IS NOT NULL AND TRIM(nama_akun_id) != ''")

	if search != "" {
		s := "%" + strings.ToLower(search) + "%"
		query = query.Where("(LOWER(no_akun) LIKE ? OR LOWER(nama_akun_id) LIKE ?)", s, s)
	}

	type CoaItem struct {
		NoAkun   string `json:"no_akun"`
		NamaAkun string `json:"nama_akun"`
	}

	var items []CoaItem
	if err := query.Order("no_akun ASC").Find(&items).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
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
	indikatorParam := strings.TrimSpace(c.Query("id_indikator"))
	if indikatorParam == "" {
		indikatorParam = strings.TrimSpace(c.Query("indikator"))
	}
	unitParam := strings.TrimSpace(c.Query("id_fakultas_unit"))
	if unitParam == "" {
		unitParam = strings.TrimSpace(c.Query("fakultas_unit"))
	}

	result, err := mediatr.Send[getTargetUnpak.GetTargetUnpakQuery, *domain.TargetUnpakResult](c.Context(), getTargetUnpak.GetTargetUnpakQuery{
		Tahun:          tahun,
		IndikatorParam: indikatorParam,
		UnitParam:      unitParam,
	})
	if err != nil {
		return commoninfra.HandleError(c, err)
	}
	return c.JSON(result)
}

// =======================================================
// POST /proker_v3/aktivitas/:aktivitas_uuid/timeline
// =======================================================
func CreateTimelineHandler(c *fiber.Ctx) error {
	actUUIDStr := c.Params("aktivitas_uuid")
	if actUUIDStr == "" {
		actUUIDStr = c.Params("uuid")
	}

	var cmd createTimeline.CreateTimelineV3Command

	if val := c.FormValue("nama_aktivitas"); val != "" {
		cmd.NamaAktivitas = val
	}
	if val := c.FormValue("pic"); val != "" {
		cmd.PIC = val
	}
	if val := c.FormValue("titik_kontrol"); val != "" {
		cmd.TitikKontrol = val
	}
	valAnggaran := strings.TrimSpace(c.FormValue("anggaran"))
	if valAnggaran == "" {
		valAnggaran = strings.TrimSpace(c.FormValue("nominal"))
	}
	if valAnggaran != "" {
		if ang, err := strconv.ParseFloat(valAnggaran, 64); err == nil {
			cmd.Anggaran = ang
		}
	}
	if val := c.FormValue("no_aktivitas"); val != "" {
		if n, err := strconv.Atoi(val); err == nil {
			cmd.NoAktivitas = n
		}
	}
	if val := c.FormValue("jadwal_rk"); val != "" {
		cmd.JadwalRk = val
	}
	if val := c.FormValue("jadwal_r"); val != "" {
		cmd.JadwalR = val
	}
	if val := c.FormValue("capaian"); val != "" {
		cmd.Capaian = val
	}
	if val := c.FormValue("masalah"); val != "" {
		cmd.Masalah = val
	}
	if val := c.FormValue("pemecahan_masalah"); val != "" {
		cmd.PemecahanMasalah = val
	}
	if val := c.FormValue("status"); val != "" {
		cmd.Status = val
	}
	if val := c.FormValue("tanggal_rk_awal"); val != "" {
		cmd.TanggalRkAwal = val
	}
	if val := c.FormValue("tanggal_rk_akhir"); val != "" {
		cmd.TanggalRkAkhir = val
	}
	if val := c.FormValue("tanggal_r_awal"); val != "" {
		cmd.TanggalRAwal = val
	}
	if val := c.FormValue("tanggal_r_akhir"); val != "" {
		cmd.TanggalRAkhir = val
	}
	if val := c.FormValue("rab"); val != "" {
		cmd.RAB = val
	}
	cmd.AktivitasUUID = actUUIDStr

	tl, err := mediatr.Send[createTimeline.CreateTimelineV3Command, *domain.ProkerTimelineV3](c.Context(), cmd)
	if err != nil {
		return commoninfra.HandleError(c, err)
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Timeline aktivitas berhasil ditambahkan",
		"uuid":    tl.UUID.String(),
		"data":    tl,
	})
}

// =======================================================
// PUT /proker_v3/timeline/:timeline_uuid
// =======================================================
func UpdateTimelineHandler(c *fiber.Ctx) error {
	tlUUIDStr := c.Params("timeline_uuid")

	var cmd updateTimeline.UpdateTimelineV3Command
	if strings.Contains(c.Get("Content-Type"), "application/json") {
		_ = c.BodyParser(&cmd)
	}

	if val := c.FormValue("nama_aktivitas"); val != "" {
		cmd.NamaAktivitas = val
	}
	if val := c.FormValue("pic"); val != "" {
		cmd.PIC = val
	}
	if val := c.FormValue("titik_kontrol"); val != "" {
		cmd.TitikKontrol = val
	}
	valUpdateAnggaran := strings.TrimSpace(c.FormValue("anggaran"))
	if valUpdateAnggaran == "" {
		valUpdateAnggaran = strings.TrimSpace(c.FormValue("nominal"))
	}
	if valUpdateAnggaran != "" {
		if ang, err := strconv.ParseFloat(valUpdateAnggaran, 64); err == nil {
			cmd.Anggaran = &ang
		}
	}
	if val := c.FormValue("no_aktivitas"); val != "" {
		if n, err := strconv.Atoi(val); err == nil {
			cmd.NoAktivitas = n
		}
	}
	if val := c.FormValue("jadwal_rk"); val != "" {
		cmd.JadwalRk = val
	}
	if val := c.FormValue("jadwal_r"); val != "" {
		cmd.JadwalR = val
	}
	if val := c.FormValue("capaian"); val != "" {
		cmd.Capaian = val
	}
	if val := c.FormValue("masalah"); val != "" {
		cmd.Masalah = val
	}
	if val := c.FormValue("pemecahan_masalah"); val != "" {
		cmd.PemecahanMasalah = val
	}
	if val := c.FormValue("status"); val != "" {
		cmd.Status = val
	}
	if val := c.FormValue("tanggal_rk_awal"); val != "" {
		cmd.TanggalRkAwal = val
	}
	if val := c.FormValue("tanggal_rk_akhir"); val != "" {
		cmd.TanggalRkAkhir = val
	}
	if val := c.FormValue("tanggal_r_awal"); val != "" {
		cmd.TanggalRAwal = val
	}
	if val := c.FormValue("tanggal_r_akhir"); val != "" {
		cmd.TanggalRAkhir = val
	}
	if val := c.FormValue("rab"); val != "" {
		cmd.RAB = val
	}
	cmd.UUID = tlUUIDStr

	tl, err := mediatr.Send[updateTimeline.UpdateTimelineV3Command, *domain.ProkerTimelineV3](c.Context(), cmd)
	if err != nil {
		return commoninfra.HandleError(c, err)
	}

	return c.JSON(fiber.Map{
		"message": "Timeline aktivitas berhasil diperbarui",
		"uuid":    tlUUIDStr,
		"data":    tl,
	})
}

// =======================================================
// PUT /proker_v3/aktivitas/:aktivitas_uuid/timeline/reorder
// =======================================================
func ReorderTimelineHandler(c *fiber.Ctx) error {
	var items []domain.ReorderTimelineItem
	if jsonStr := c.FormValue("items"); jsonStr != "" {
		_ = json.Unmarshal([]byte(jsonStr), &items)
	}

	count, err := mediatr.Send[reorderTimeline.ReorderTimelineV3Command, int](c.Context(), reorderTimeline.ReorderTimelineV3Command{Items: items})
	if err != nil {
		return commoninfra.HandleError(c, err)
	}

	return c.JSON(fiber.Map{
		"message": "Urutan nomor urut tahapan timeline berhasil diperbarui",
		"count":   count,
	})
}

// =======================================================
// DELETE /proker_v3/timeline/:timeline_uuid
// =======================================================
func DeleteTimelineHandler(c *fiber.Ctx) error {
	tlUUIDStr := c.Params("timeline_uuid")
	deletedUUID, err := mediatr.Send[deleteTimeline.DeleteTimelineV3Command, string](c.Context(), deleteTimeline.DeleteTimelineV3Command{UUID: tlUUIDStr})
	if err != nil {
		return commoninfra.HandleError(c, err)
	}

	return c.JSON(fiber.Map{
		"message": "Timeline aktivitas berhasil dihapus",
		"uuid":    deletedUUID,
	})
}

// =======================================================
// DELETE /proker_v3/file/:file_uuid
// =======================================================
func DeleteTimelineFileHandler(c *fiber.Ctx) error {
	fileUUIDStr := c.Params("file_uuid")
	deletedUUID, err := mediatr.Send[deleteTimelineFile.DeleteTimelineFileV3Command, string](c.Context(), deleteTimelineFile.DeleteTimelineFileV3Command{FileUUID: fileUUIDStr})
	if err != nil {
		return commoninfra.HandleError(c, err)
	}

	return c.JSON(fiber.Map{
		"message": "File berhasil dihapus",
		"uuid":    deletedUUID,
	})
}

// =======================================================
// GET /proker_v3/:uuid/dokumen (Daftar matrix_proker_dokumen)
// =======================================================
func GetProkerDokumenHandler(c *fiber.Ctx) error {
	rawUuid := c.Params("uuid")
	parsedUuid, err := uuid.Parse(rawUuid)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid matrix UUID"})
	}

	result, err := mediatr.Send[getDokumen.GetDokumenV3Query, *domain.ProkerDokumenResult](c.Context(), getDokumen.GetDokumenV3Query{ProkerUUID: parsedUuid})
	if err != nil {
		return commoninfra.HandleError(c, err)
	}

	return c.JSON(fiber.Map{
		"proker":     result.Proker,
		"activities": result.Activities,
		"data":       result.Data,
	})
}

// =======================================================
// DELETE /proker_v3/dokumen/:doc_uuid (Hapus matrix_proker_dokumen)
// =======================================================
func DeleteProkerDokumenHandler(c *fiber.Ctx) error {
	rawUuid := c.Params("doc_uuid")
	if rawUuid == "" {
		rawUuid = c.Params("uuid")
	}

	deletedUUID, err := mediatr.Send[deleteDokumen.DeleteDokumenV3Command, string](c.Context(), deleteDokumen.DeleteDokumenV3Command{DocUUID: rawUuid})
	if err != nil {
		return commoninfra.HandleError(c, err)
	}

	return c.JSON(fiber.Map{
		"message": "Dokumen berhasil dihapus",
		"uuid":    deletedUUID,
	})
}

// =======================================================
// PUT /proker_v3/dokumen/:doc_uuid/verifikasi
// =======================================================
func UpdateVerifikasiDokumenHandler(c *fiber.Ctx) error {
	rawUuid := c.Params("doc_uuid")
	if rawUuid == "" {
		rawUuid = c.Params("uuid")
	}

	var cmd updateVerifikasiDokumen.UpdateVerifikasiDokumenV3Command

	if val := c.FormValue("status_verifikasi"); val != "" {
		cmd.StatusVerifikasi = val
	}
	if val := c.FormValue("status"); val != "" {
		cmd.Status = val
	}
	if val := c.FormValue("catatan"); val != "" {
		cmd.Catatan = val
	}
	if val := c.FormValue("catatan_penolakan"); val != "" {
		cmd.CatatanPenolakan = val
	}
	cmd.DocUUID = rawUuid

	doc, err := mediatr.Send[updateVerifikasiDokumen.UpdateVerifikasiDokumenV3Command, *domain.MatrixProkerDokumen](c.Context(), cmd)
	if err != nil {
		return commoninfra.HandleError(c, err)
	}

	return c.JSON(fiber.Map{
		"message":           "Status verifikasi dokumen berhasil diperbarui",
		"uuid":              rawUuid,
		"status_verifikasi": doc.StatusVerifikasi,
		"catatan":           doc.Catatan,
		"catatan_penolakan": doc.CatatanPenolakan,
		"data":              doc,
	})
}

// =======================================================
// MULTIPART / FORM-DATA UPLOAD HANDLERS
// =======================================================

func sanitizeFolderSegment(s string) string {
	s = strings.TrimSpace(s)
	if s == "" {
		return "umum"
	}
	replacer := strings.NewReplacer(
		"/", "-",
		"\\", "-",
		":", "-",
		"*", "",
		"?", "",
		"\"", "",
		"<", "",
		">", "",
		"|", "-",
		"\n", "",
		"\r", "",
		"\t", " ",
	)
	cleaned := strings.TrimSpace(replacer.Replace(s))
	for strings.Contains(cleaned, "  ") {
		cleaned = strings.ReplaceAll(cleaned, "  ", " ")
	}
	if len(cleaned) > 100 {
		cleaned = cleaned[:100]
	}
	if cleaned == "" {
		return "umum"
	}
	return cleaned
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

func resolveTimelineUploadDir(ctx context.Context, timeline *domain.ProkerTimelineV3) string {
	if dbInstance == nil || repoProker == nil {
		return filepath.Join("./upload/prokerv2", "umum", "2026", "standar", "indikator", "aktivitas", sanitizeFolderSegment(timeline.NamaAktivitas))
	}

	var act domain.AktivitasV2
	if err := dbInstance.WithContext(ctx).Where("id = ?", timeline.IdAktivitasV2).First(&act).Error; err != nil {
		return filepath.Join("./upload/prokerv2", "umum", "2026", "standar", "indikator", "aktivitas", sanitizeFolderSegment(timeline.NamaAktivitas))
	}

	var proker *domain.ProkerV3Default
	if act.IdMatrikProker != nil && *act.IdMatrikProker > 0 {
		var mp domain.ProkerV3
		if err := dbInstance.WithContext(ctx).Where("id = ?", *act.IdMatrikProker).First(&mp).Error; err == nil {
			proker, _ = repoProker.GetDefaultByUuid(ctx, mp.UUID)
		}
	}

	unitSegment := ""
	if proker != nil && strings.TrimSpace(proker.NamaFakultasUnit) != "" {
		unitSegment = proker.NamaFakultasUnit
	} else if proker != nil && proker.IdFakultasUnit != nil {
		unitSegment = fmt.Sprintf("%d", *proker.IdFakultasUnit)
	} else if act.IdFakultasUnit != nil {
		unitSegment = fmt.Sprintf("%d", *act.IdFakultasUnit)
	}
	namaUnit := sanitizeFolderSegment(unitSegment)

	tahun := "2026"
	if proker != nil && strings.TrimSpace(proker.Tahun) != "" {
		tahun = proker.Tahun
	}
	tahunFolder := sanitizeFolderSegment(tahun)

	standar := "standar"
	if proker != nil && strings.TrimSpace(proker.StandarSubStandar) != "" {
		standar = proker.StandarSubStandar
	}
	standarFolder := sanitizeFolderSegment(standar)

	indikator := "indikator"
	if proker != nil && strings.TrimSpace(proker.IndikatorKinerja) != "" {
		indikator = proker.IndikatorKinerja
	}
	indikatorFolder := sanitizeFolderSegment(indikator)

	aktivitasNama := "aktivitas"
	if strings.TrimSpace(act.Aktivitas) != "" {
		aktivitasNama = act.Aktivitas
	}
	aktivitasFolder := sanitizeFolderSegment(aktivitasNama)

	timelineNama := "tahapan"
	if strings.TrimSpace(timeline.NamaAktivitas) != "" {
		timelineNama = timeline.NamaAktivitas
	}
	timelineFolder := sanitizeFolderSegment(timelineNama)

	return filepath.Join("./upload/prokerv2", namaUnit, tahunFolder, standarFolder, indikatorFolder, aktivitasFolder, timelineFolder)
}

// =======================================================
// POST /proker_v3/:uuid/dokumen_referensi
// =======================================================
func UploadDokumenReferensiHandler(c *fiber.Ctx) error {
	rawUuid := c.Params("uuid")
	parsedUuid, err := uuid.Parse(rawUuid)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid UUID"})
	}

	var act *domain.AktivitasV2
	if repoProker != nil {
		act, _ = repoProker.GetAktivitasEntityByUuid(c.Context(), parsedUuid)
	}

	form, err := c.MultipartForm()
	if err != nil || form == nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Gagal membaca file upload: " + err.Error()})
	}

	fileHeaders := form.File["files"]
	if len(fileHeaders) == 0 {
		fileHeaders = form.File["file"]
	}
	if len(fileHeaders) == 0 {
		fileHeaders = form.File["files[]"]
	}
	if len(fileHeaders) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Tidak ada file yang dipilih untuk diupload"})
	}

	keterangan := strings.TrimSpace(c.FormValue("keterangan"))
	if keterangan == "" {
		keterangan = "Dokumen Referensi Program"
	}

	uploadDir := "./upload/proker_v3"
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal membuat direktori upload"})
	}

	var filesToInsert []domain.ProkerTimelineFileV3
	for _, fh := range fileHeaders {
		ext := filepath.Ext(fh.Filename)
		fileUUID := uuid.New()
		savedName := fmt.Sprintf("docref-%s%s", fileUUID.String(), ext)
		destPath := filepath.Join(uploadDir, savedName)

		if err := c.SaveFile(fh, destPath); err != nil {
			continue
		}

		var actIDPtr *uint
		if act != nil {
			actIDPtr = &act.ID
		}

		filesToInsert = append(filesToInsert, domain.ProkerTimelineFileV3{
			UUID:          fileUUID,
			IdAktivitasV2: actIDPtr,
			File:          savedName,
			NamaFile:      fh.Filename,
			FileURL:       "/upload/proker_v3/" + savedName,
			FileSize:      fh.Size,
			FileType:      fh.Header.Get("Content-Type"),
			Keterangan:    keterangan,
		})
	}

	if len(filesToInsert) > 0 && repoProker != nil {
		if err := repoProker.AddTimelineFiles(c.Context(), filesToInsert); err != nil {
			return commoninfra.HandleError(c, err)
		}
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Dokumen referensi berhasil diunggah",
		"count":   len(filesToInsert),
		"files":   filesToInsert,
	})
}

// =======================================================
// POST /proker_v3/timeline/:timeline_uuid/files
// =======================================================
func UploadTimelineFilesHandler(c *fiber.Ctx) error {
	timelineUUIDStr := c.Params("timeline_uuid")
	timelineUUID, err := uuid.Parse(timelineUUIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid timeline UUID"})
	}

	timeline, err := repoProker.GetTimelineByUuid(c.Context(), timelineUUID)
	if err != nil {
		return commoninfra.HandleError(c, err)
	}
	if timeline == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Timeline tidak ditemukan"})
	}

	form, err := c.MultipartForm()
	if err != nil || form == nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Gagal membaca file upload: " + err.Error()})
	}

	fileHeaders := form.File["files"]
	if len(fileHeaders) == 0 {
		fileHeaders = form.File["file"]
	}
	if len(fileHeaders) == 0 {
		fileHeaders = form.File["files[]"]
	}
	if len(fileHeaders) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Tidak ada file yang dipilih untuk diupload"})
	}

	keterangan := c.FormValue("keterangan")
	uploadDir := resolveTimelineUploadDir(c.Context(), timeline)
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal membuat direktori upload: " + err.Error()})
	}

	var filesToInsert []domain.ProkerTimelineFileV3
	for _, fh := range fileHeaders {
		ext := filepath.Ext(fh.Filename)
		fileUUID := uuid.New()
		savedName := fmt.Sprintf("proker-%s%s", fileUUID.String(), ext)
		destPath := filepath.Join(uploadDir, savedName)

		if err := c.SaveFile(fh, destPath); err != nil {
			continue
		}

		relURL := filepath.ToSlash(destPath)
		if !strings.HasPrefix(relURL, "/") {
			relURL = "/" + strings.TrimPrefix(relURL, "./")
		}

		filesToInsert = append(filesToInsert, domain.ProkerTimelineFileV3{
			UUID:                  fileUUID,
			IdAktivitasV2Timeline: &timeline.ID,
			IdAktivitasV2:         &timeline.IdAktivitasV2,
			File:                  destPath,
			NamaFile:              fh.Filename,
			FileURL:               relURL,
			FileSize:              fh.Size,
			FileType:              fh.Header.Get("Content-Type"),
			Keterangan:            keterangan,
		})
	}

	if len(filesToInsert) > 0 {
		if err := repoProker.AddTimelineFiles(c.Context(), filesToInsert); err != nil {
			return commoninfra.HandleError(c, err)
		}
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "File berhasil diunggah",
		"count":   len(filesToInsert),
		"files":   filesToInsert,
	})
}

// =======================================================
// POST /proker_v3/timeline/:timeline_uuid/chunk
// =======================================================
func UploadTimelineFileChunkHandler(c *fiber.Ctx) error {
	tlUUIDStr := c.Params("timeline_uuid")
	timelineUUID, err := uuid.Parse(tlUUIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid timeline UUID"})
	}

	timeline, err := repoProker.GetTimelineByUuid(c.Context(), timelineUUID)
	if err != nil {
		return commoninfra.HandleError(c, err)
	}
	if timeline == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Timeline tahapan kegiatan tidak ditemukan"})
	}

	uploadIDStr := strings.TrimSpace(c.FormValue("upload_id"))
	if uploadIDStr == "" {
		uploadIDStr = strings.TrimSpace(c.Get("X-Upload-ID"))
	}
	fileUUID, err := uuid.Parse(uploadIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid or missing upload_id (must be UUID)"})
	}

	chunkIndex, _ := strconv.Atoi(c.FormValue("chunk_index"))
	totalChunks, _ := strconv.Atoi(c.FormValue("total_chunks"))
	chunkOffset, _ := strconv.ParseInt(c.FormValue("chunk_offset"), 10, 64)
	totalSize, _ := strconv.ParseInt(c.FormValue("total_size"), 10, 64)
	fileName := strings.TrimSpace(c.FormValue("file_name"))
	keterangan := strings.TrimSpace(c.FormValue("keterangan"))

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

	partPath := filepath.Join(tmpDir, fmt.Sprintf("%s.part", fileUUID.String()))
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

	if _, err := io.Copy(partFile, src); err != nil {
		src.Close()
		partFile.Close()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal menulis binary chunk: " + err.Error()})
	}
	src.Close()
	partFile.Close()

	if totalChunks > 0 && chunkIndex >= totalChunks-1 {
		targetDir := resolveTimelineUploadDir(c.Context(), timeline)
		if err := os.MkdirAll(targetDir, 0755); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal membuat direktori upload: " + err.Error()})
		}

		ext := filepath.Ext(fileName)
		if ext == "" {
			ext = filepath.Ext(fileHeader.Filename)
		}
		savedFileName := fmt.Sprintf("proker-%s%s", fileUUID.String(), ext)
		destPath := filepath.Join(targetDir, savedFileName)

		if err := os.Rename(partPath, destPath); err != nil {
			if copyErr := copyFile(partPath, destPath); copyErr != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal memindahkan file bukti: " + copyErr.Error()})
			}
			_ = os.Remove(partPath)
		}

		relURL := filepath.ToSlash(destPath)
		if !strings.HasPrefix(relURL, "/") {
			relURL = "/" + strings.TrimPrefix(relURL, "./")
		}

		record := domain.ProkerTimelineFileV3{
			UUID:                  fileUUID,
			IdAktivitasV2Timeline: &timeline.ID,
			IdAktivitasV2:         &timeline.IdAktivitasV2,
			File:                  destPath,
			NamaFile:              fileName,
			FileURL:               relURL,
			FileSize:              totalSize,
			FileType:              fileHeader.Header.Get("Content-Type"),
			Keterangan:            keterangan,
		}

		if err := repoProker.AddTimelineFiles(c.Context(), []domain.ProkerTimelineFileV3{record}); err != nil {
			_ = os.Remove(destPath)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal menyimpan record bukti ke database: " + err.Error()})
		}

		return c.Status(fiber.StatusCreated).JSON(fiber.Map{
			"message":   "Berkas bukti berhasil diunggah secara utuh",
			"completed": true,
			"data":      record,
		})
	}

	return c.JSON(fiber.Map{
		"message":     "Chunk binary berhasil diterima",
		"completed":   false,
		"chunk_index": chunkIndex,
		"uploaded":    chunkOffset + fileHeader.Size,
	})
}

// =======================================================
// DELETE /proker_v3/timeline/:timeline_uuid/chunk/:upload_id
// =======================================================
func CancelTimelineFileChunkHandler(c *fiber.Ctx) error {
	uploadIDStr := strings.TrimSpace(c.Params("upload_id"))
	if uploadIDStr == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Missing upload_id"})
	}
	partPath := filepath.Join("./upload/tmp", fmt.Sprintf("%s.part", uploadIDStr))
	_ = os.Remove(partPath)
	return c.JSON(fiber.Map{"message": "Upload berkas bukti dibatalkan dan file temporary dihapus"})
}

// =======================================================
// POST /proker_v3/:uuid/dokumen/chunk
// =======================================================
func UploadProkerDokumenChunkHandler(c *fiber.Ctx) error {
	rawUuid := c.Params("uuid")
	parsedUuid, err := uuid.Parse(rawUuid)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid matrix UUID"})
	}

	proker, err := repoProker.GetDefaultByUuid(c.Context(), parsedUuid)
	if err != nil {
		return commoninfra.HandleError(c, err)
	}
	if proker == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Program kerja tidak ditemukan"})
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

	if _, err := io.Copy(partFile, src); err != nil {
		src.Close()
		partFile.Close()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal menulis binary chunk: " + err.Error()})
	}
	src.Close()
	partFile.Close()

	if totalChunks > 0 && chunkIndex >= totalChunks-1 {
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

		if err := repoProker.CreateDokumen(c.Context(), &doc); err != nil {
			_ = os.Remove(destPath)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal menyimpan record dokumen ke database: " + err.Error()})
		}

		return c.Status(fiber.StatusCreated).JSON(fiber.Map{
			"message":   "Dokumen berhasil diunggah secara utuh",
			"completed": true,
			"data":      doc,
		})
	}

	return c.JSON(fiber.Map{
		"message":     "Chunk binary berhasil diterima",
		"completed":   false,
		"chunk_index": chunkIndex,
		"uploaded":    chunkOffset + fileHeader.Size,
	})
}

// =======================================================
// DELETE /proker_v3/:uuid/dokumen/chunk/:upload_id
// =======================================================
func CancelProkerDokumenChunkHandler(c *fiber.Ctx) error {
	uploadIDStr := strings.TrimSpace(c.Params("upload_id"))
	if uploadIDStr == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Missing upload_id"})
	}
	partPath := filepath.Join("./upload/tmp", fmt.Sprintf("%s.part", uploadIDStr))
	_ = os.Remove(partPath)
	return c.JSON(fiber.Map{"message": "Upload chunk dibatalkan dan file temporary dihapus"})
}

// =======================================================
// POST /proker_v3/:uuid/dokumen (Upload ke matrix_proker_dokumen)
// =======================================================
func UploadProkerDokumenHandler(c *fiber.Ctx) error {
	rawUuid := c.Params("uuid")
	parsedUuid, err := uuid.Parse(rawUuid)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid matrix UUID"})
	}

	proker, err := repoProker.GetDefaultByUuid(c.Context(), parsedUuid)
	if err != nil {
		return commoninfra.HandleError(c, err)
	}
	if proker == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Program kerja tidak ditemukan"})
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

		if err := repoProker.CreateDokumen(c.Context(), &doc); err != nil {
			_ = os.Remove(destPath)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Gagal menyimpan record dokumen ke database: " + err.Error(),
			})
		}
		createdDocs = append(createdDocs, doc)
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Dokumen berhasil diunggah",
		"count":   len(createdDocs),
		"data":    createdDocs,
	})
}

// =======================================================
// ModuleProkerV3 route initialization
// =======================================================
func ModuleProkerV3(app *fiber.App, db *gorm.DB) {
	repoProker = infra.NewProkerV3Repository(db)
	dbInstance = db
	_ = db.AutoMigrate(&domain.JadwalProkerV3{}, &domain.JadwalAktivitasV3{})

	roles := []string{"admin", "auditee", "fakultas", "unit", "user", "upml", "upmf", "upms", "anggaran", "keuangan", "evaluator"}
	whoamiURL := "http://localhost:3000/whoami"

	// 1. Options & Target Unpak & COA
	app.Get("/proker_v3_options", commonpresentation.SmartCompress(), commonpresentation.JWTMiddleware(), GetProkerV3OptionsHandler)
	app.Get("/proker_v3/target_unpak", commonpresentation.SmartCompress(), commonpresentation.JWTMiddleware(), GetTargetUnpakHandler)
	app.Get("/proker_v3/coa", commonpresentation.SmartCompress(), commonpresentation.JWTMiddleware(), GetCoaHandler)

	// Jadwal Proker V3 CRUD & Status Check (Didaftarkan sebelum /proker_v3/:uuid agar tidak tertabrak route param :uuid)
	adminRoles := []string{"admin"}
	app.Get("/proker_v3/jadwal", commonpresentation.SmartCompress(), commonpresentation.JWTMiddleware(), GetAllJadwalProkerV3Handler)
	app.Get("/proker_v3/jadwal_status", commonpresentation.SmartCompress(), commonpresentation.JWTMiddleware(), CheckJadwalProkerV3Handler)
	app.Get("/proker_v3/jadwal/:uuid", commonpresentation.SmartCompress(), commonpresentation.JWTMiddleware(), GetJadwalProkerV3Handler)
	app.Post("/proker_v3/jadwal", commonpresentation.JWTMiddleware(), commonpresentation.RBACMiddleware(adminRoles, whoamiURL), CreateJadwalProkerV3Handler)
	app.Put("/proker_v3/jadwal/:uuid", commonpresentation.JWTMiddleware(), commonpresentation.RBACMiddleware(adminRoles, whoamiURL), UpdateJadwalProkerV3Handler)
	app.Delete("/proker_v3/jadwal/:uuid", commonpresentation.JWTMiddleware(), commonpresentation.RBACMiddleware(adminRoles, whoamiURL), DeleteJadwalProkerV3Handler)

	// Jadwal Aktivitas V3 CRUD & Status Check (Pengaturan Buka/Tutup Tambah Aktivitas Baru)
	app.Get("/proker_v3/jadwal_aktivitas", commonpresentation.SmartCompress(), commonpresentation.JWTMiddleware(), GetAllJadwalAktivitasV3Handler)
	app.Get("/proker_v3/jadwal_aktivitas_status", commonpresentation.SmartCompress(), commonpresentation.JWTMiddleware(), CheckJadwalAktivitasV3Handler)
	app.Get("/proker_v3/jadwal_aktivitas/:uuid", commonpresentation.SmartCompress(), commonpresentation.JWTMiddleware(), GetJadwalAktivitasV3Handler)
	app.Post("/proker_v3/jadwal_aktivitas", commonpresentation.JWTMiddleware(), commonpresentation.RBACMiddleware(adminRoles, whoamiURL), CreateJadwalAktivitasV3Handler)
	app.Put("/proker_v3/jadwal_aktivitas/:uuid", commonpresentation.JWTMiddleware(), commonpresentation.RBACMiddleware(adminRoles, whoamiURL), UpdateJadwalAktivitasV3Handler)
	app.Delete("/proker_v3/jadwal_aktivitas/:uuid", commonpresentation.JWTMiddleware(), commonpresentation.RBACMiddleware(adminRoles, whoamiURL), DeleteJadwalAktivitasV3Handler)

	// Pagu Anggaran CRUD
	app.Get("/pagu_anggaran", commonpresentation.SmartCompress(), commonpresentation.JWTMiddleware(), GetAllPaguAnggaranHandler)
	app.Post("/pagu_anggaran", commonpresentation.JWTMiddleware(), commonpresentation.RBACMiddleware(roles, whoamiURL), SavePaguAnggaranHandler)
	app.Delete("/pagu_anggaran/:id", commonpresentation.JWTMiddleware(), commonpresentation.RBACMiddleware(roles, whoamiURL), DeletePaguAnggaranHandler)

	// Dashboard Monitoring Proker, Aktivitas & RAB (Wakil Rektor 1-4 & Fakultas)
	app.Get("/proker_v3/dashboard_monitoring", commonpresentation.SmartCompress(), commonpresentation.JWTMiddleware(), GetDashboardProkerMonitoringHandler)

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

	// 3. Aktivitas (aktivitasv2) CRUD per Matrix Proker
	app.Get("/proker_v3/:uuid/aktivitas", commonpresentation.SmartCompress(), commonpresentation.JWTMiddleware(), commonpresentation.RBACMiddleware(roles, whoamiURL), GetAktivitasByMatrixHandler)
	app.Get("/proker_v3/aktivitas/:aktivitas_uuid", commonpresentation.SmartCompress(), commonpresentation.JWTMiddleware(), commonpresentation.RBACMiddleware(roles, whoamiURL), GetAktivitasDetailHandler)
	app.Get("/proker_v3/:uuid/aktivitas/:aktivitas_uuid", commonpresentation.SmartCompress(), commonpresentation.JWTMiddleware(), commonpresentation.RBACMiddleware(roles, whoamiURL), GetAktivitasDetailHandler)
	app.Post("/proker_v3/:uuid/aktivitas", commonpresentation.JWTMiddleware(), commonpresentation.RBACMiddleware(roles, whoamiURL), CreateAktivitasHandler)
	app.Put("/proker_v3/aktivitas/:aktivitas_uuid", commonpresentation.JWTMiddleware(), commonpresentation.RBACMiddleware(roles, whoamiURL), UpdateAktivitasHandler)
	app.Delete("/proker_v3/aktivitas/:aktivitas_uuid", commonpresentation.JWTMiddleware(), commonpresentation.RBACMiddleware(roles, whoamiURL), DeleteAktivitasHandler)
	app.Put("/proker_v3/aktivitas/:aktivitas_uuid/status", commonpresentation.JWTMiddleware(), commonpresentation.RBACMiddleware(roles, whoamiURL), UpdateAktivitasStatusHandler)
	app.Put("/proker_v3/aktivitas/:aktivitas_uuid/batal", commonpresentation.JWTMiddleware(), commonpresentation.RBACMiddleware(roles, whoamiURL), UpdateBatalAktivitasHandler)
	app.Patch("/proker_v3/aktivitas/:aktivitas_uuid/batal", commonpresentation.JWTMiddleware(), commonpresentation.RBACMiddleware(roles, whoamiURL), UpdateBatalAktivitasHandler)
	app.Put("/proker_v3/:uuid/aktivitas/:aktivitas_uuid/batal", commonpresentation.JWTMiddleware(), commonpresentation.RBACMiddleware(roles, whoamiURL), UpdateBatalAktivitasHandler)
	app.Patch("/proker_v3/:uuid/aktivitas/:aktivitas_uuid/batal", commonpresentation.JWTMiddleware(), commonpresentation.RBACMiddleware(roles, whoamiURL), UpdateBatalAktivitasHandler)
	app.Post("/proker_v3/aktivitas/:aktivitas_uuid/dokumen_referensi", commonpresentation.JWTMiddleware(), commonpresentation.RBACMiddleware(roles, whoamiURL), UploadDokumenReferensiHandler)

	// Backward-compat aliases & Verification
	app.Put("/proker_v3/:uuid/status", commonpresentation.JWTMiddleware(), commonpresentation.RBACMiddleware(roles, whoamiURL), UpdateAktivitasStatusHandler)
	app.Put("/proker_v3/:uuid/verifikasi", commonpresentation.JWTMiddleware(), commonpresentation.RBACMiddleware(roles, whoamiURL), UpdateAktivitasStatusHandler)
	app.Post("/proker_v3/:uuid/dokumen_referensi", commonpresentation.JWTMiddleware(), commonpresentation.RBACMiddleware(roles, whoamiURL), UploadDokumenReferensiHandler)

	// 4. Timeline CRUD & Upload files
	app.Post("/proker_v3/aktivitas/:aktivitas_uuid/timeline", commonpresentation.JWTMiddleware(), commonpresentation.RBACMiddleware(roles, whoamiURL), CreateTimelineHandler)
	app.Post("/proker_v3/:uuid/timeline", commonpresentation.JWTMiddleware(), commonpresentation.RBACMiddleware(roles, whoamiURL), CreateTimelineHandler)
	app.Put("/proker_v3/timeline/:timeline_uuid", commonpresentation.JWTMiddleware(), commonpresentation.RBACMiddleware(roles, whoamiURL), UpdateTimelineHandler)
	app.Put("/proker_v3/:uuid/timeline/:timeline_uuid", commonpresentation.JWTMiddleware(), commonpresentation.RBACMiddleware(roles, whoamiURL), UpdateTimelineHandler)
	app.Delete("/proker_v3/timeline/:timeline_uuid", commonpresentation.JWTMiddleware(), commonpresentation.RBACMiddleware(roles, whoamiURL), DeleteTimelineHandler)
	app.Delete("/proker_v3/:uuid/timeline/:timeline_uuid", commonpresentation.JWTMiddleware(), commonpresentation.RBACMiddleware(roles, whoamiURL), DeleteTimelineHandler)
	app.Put("/proker_v3/aktivitas/:aktivitas_uuid/timeline/reorder", commonpresentation.JWTMiddleware(), commonpresentation.RBACMiddleware(roles, whoamiURL), ReorderTimelineHandler)
	app.Put("/proker_v3/:uuid/timeline/reorder", commonpresentation.JWTMiddleware(), commonpresentation.RBACMiddleware(roles, whoamiURL), ReorderTimelineHandler)

	// Multiple file upload for timelines
	app.Post("/proker_v3/timeline/:timeline_uuid/files", commonpresentation.JWTMiddleware(), commonpresentation.RBACMiddleware(roles, whoamiURL), UploadTimelineFilesHandler)
	app.Post("/proker_v3/:uuid/timeline/:timeline_uuid/files", commonpresentation.JWTMiddleware(), commonpresentation.RBACMiddleware(roles, whoamiURL), UploadTimelineFilesHandler)
	app.Post("/proker_v3/timeline/:timeline_uuid/chunk", commonpresentation.JWTMiddleware(), commonpresentation.RBACMiddleware(roles, whoamiURL), UploadTimelineFileChunkHandler)
	app.Post("/proker_v3/:uuid/timeline/:timeline_uuid/chunk", commonpresentation.JWTMiddleware(), commonpresentation.RBACMiddleware(roles, whoamiURL), UploadTimelineFileChunkHandler)
	app.Delete("/proker_v3/timeline/:timeline_uuid/chunk/:upload_id", commonpresentation.JWTMiddleware(), commonpresentation.RBACMiddleware(roles, whoamiURL), CancelTimelineFileChunkHandler)
	app.Delete("/proker_v3/:uuid/timeline/:timeline_uuid/chunk/:upload_id", commonpresentation.JWTMiddleware(), commonpresentation.RBACMiddleware(roles, whoamiURL), CancelTimelineFileChunkHandler)
	app.Delete("/proker_v3/file/:file_uuid", commonpresentation.JWTMiddleware(), commonpresentation.RBACMiddleware(roles, whoamiURL), DeleteTimelineFileHandler)

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

// =======================================================
// 7. JADWAL PROKER V3 CRUD HANDLERS
// =======================================================

type JadwalProkerV3Row struct {
	domain.JadwalProkerV3
	NamaFakultasUnit string `json:"nama_fakultas_unit"`
	TypeFakultasUnit string `json:"type_fakultas_unit"`
}

func GetAllJadwalProkerV3Handler(c *fiber.Ctx) error {
	if dbInstance == nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Database tidak tersedia"})
	}

	search := strings.TrimSpace(c.Query("search"))
	tahun := strings.TrimSpace(c.Query("tahun"))
	filterUnit := strings.TrimSpace(c.Query("filter_unit"))

	var list []JadwalProkerV3Row
	query := dbInstance.WithContext(c.Context()).
		Table("jadwal_proker_v3 j").
		Select("j.*, COALESCE(fu.nama_fak_prod_unit, '') as nama_fakultas_unit, COALESCE(fu.type, '') as type_fakultas_unit").
		Joins("LEFT JOIN v_fakultas_unit fu ON j.id_fakultas_unit = fu.id")

	if tahun != "" && tahun != "all" {
		query = query.Where("j.tahun = ?", tahun)
	}
	if filterUnit == "global" {
		query = query.Where("j.id_fakultas_unit IS NULL OR j.id_fakultas_unit = 0")
	} else if filterUnit == "specific" {
		query = query.Where("j.id_fakultas_unit IS NOT NULL AND j.id_fakultas_unit > 0")
	} else if filterUnit != "" && filterUnit != "all" {
		query = query.Where("j.id_fakultas_unit = ?", filterUnit)
	}
	if search != "" {
		likeSearch := "%" + search + "%"
		query = query.Where("(j.tahun LIKE ? OR fu.nama_fak_prod_unit LIKE ?)", likeSearch, likeSearch)
	}

	if err := query.Order("j.tahun DESC, (j.id_fakultas_unit IS NULL) DESC, j.tanggal_mulai ASC").Find(&list).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	now := time.Now()
	res := make([]domain.JadwalProkerV3DTO, 0, len(list))
	for _, item := range list {
		startOfDay := time.Date(item.TanggalMulai.Year(), item.TanggalMulai.Month(), item.TanggalMulai.Day(), 0, 0, 0, 0, item.TanggalMulai.Location())
		endOfDay := time.Date(item.TanggalAkhir.Year(), item.TanggalAkhir.Month(), item.TanggalAkhir.Day(), 23, 59, 59, 999999999, item.TanggalAkhir.Location())

		status := "berjalan"
		statusLabel := "Sedang Dibuka"
		if now.Before(startOfDay) {
			status = "belum_dibuka"
			statusLabel = "Belum Dibuka"
		} else if now.After(endOfDay) {
			status = "ditutup"
			statusLabel = "Ditutup"
		}

		namaUnit := item.NamaFakultasUnit
		typeUnit := item.TypeFakultasUnit
		if item.IdFakultasUnit == nil || *item.IdFakultasUnit == 0 {
			namaUnit = "Semua Unit / Fakultas (Global)"
			typeUnit = "all"
		}

		res = append(res, domain.JadwalProkerV3DTO{
			ID:               item.ID,
			UUID:             item.UUID.String(),
			Tahun:            item.Tahun,
			TanggalMulai:     item.TanggalMulai.Format("2006-01-02"),
			TanggalAkhir:     item.TanggalAkhir.Format("2006-01-02"),
			IdFakultasUnit:   item.IdFakultasUnit,
			NamaFakultasUnit: namaUnit,
			TypeFakultasUnit: typeUnit,
			Status:           status,
			StatusLabel:      statusLabel,
			CreatedAt:        item.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt:        item.UpdatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	return c.JSON(fiber.Map{
		"data":  res,
		"total": len(res),
	})
}

func GetJadwalProkerV3Handler(c *fiber.Ctx) error {
	if dbInstance == nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Database tidak tersedia"})
	}

	rawUuid := c.Params("uuid")
	parsedUuid, err := uuid.Parse(rawUuid)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid UUID"})
	}

	var item JadwalProkerV3Row
	if err := dbInstance.WithContext(c.Context()).
		Table("jadwal_proker_v3 j").
		Select("j.*, COALESCE(fu.nama_fak_prod_unit, '') as nama_fakultas_unit, COALESCE(fu.type, '') as type_fakultas_unit").
		Joins("LEFT JOIN v_fakultas_unit fu ON j.id_fakultas_unit = fu.id").
		Where("j.uuid = ?", parsedUuid).First(&item).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Jadwal proker tidak ditemukan"})
	}

	now := time.Now()
	startOfDay := time.Date(item.TanggalMulai.Year(), item.TanggalMulai.Month(), item.TanggalMulai.Day(), 0, 0, 0, 0, item.TanggalMulai.Location())
	endOfDay := time.Date(item.TanggalAkhir.Year(), item.TanggalAkhir.Month(), item.TanggalAkhir.Day(), 23, 59, 59, 999999999, item.TanggalAkhir.Location())

	status := "berjalan"
	statusLabel := "Sedang Dibuka"
	if now.Before(startOfDay) {
		status = "belum_dibuka"
		statusLabel = "Belum Dibuka"
	} else if now.After(endOfDay) {
		status = "ditutup"
		statusLabel = "Ditutup"
	}

	namaUnit := item.NamaFakultasUnit
	typeUnit := item.TypeFakultasUnit
	if item.IdFakultasUnit == nil || *item.IdFakultasUnit == 0 {
		namaUnit = "Semua Unit / Fakultas (Global)"
		typeUnit = "all"
	}

	return c.JSON(domain.JadwalProkerV3DTO{
		ID:               item.ID,
		UUID:             item.UUID.String(),
		Tahun:            item.Tahun,
		TanggalMulai:     item.TanggalMulai.Format("2006-01-02"),
		TanggalAkhir:     item.TanggalAkhir.Format("2006-01-02"),
		IdFakultasUnit:   item.IdFakultasUnit,
		NamaFakultasUnit: namaUnit,
		TypeFakultasUnit: typeUnit,
		Status:           status,
		StatusLabel:      statusLabel,
		CreatedAt:        item.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:        item.UpdatedAt.Format("2006-01-02 15:04:05"),
	})
}

func CheckJadwalProkerV3Handler(c *fiber.Ctx) error {
	if dbInstance == nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Database tidak tersedia"})
	}

	tahun := strings.TrimSpace(c.Query("tahun"))
	if tahun == "" {
		tahun = "2027"
	}

	unitParam := strings.TrimSpace(c.Query("id_fakultas_unit"))
	if unitParam == "" {
		unitParam = strings.TrimSpace(c.Query("fakultas_unit"))
	}

	var targetUnitID uint64
	var targetUnitName string
	var parentFakultasID uint64

	if unitParam != "" && unitParam != "all" && unitParam != "0" {
		var fu struct {
			ID       uint64 `gorm:"column:id"`
			Nama     string `gorm:"column:nama_fak_prod_unit"`
			Type     string `gorm:"column:type"`
			Fakultas string `gorm:"column:fakultas"`
		}
		var qFu *gorm.DB
		if _, errU := uuid.Parse(unitParam); errU == nil {
			qFu = dbInstance.WithContext(c.Context()).Table("v_fakultas_unit").Where("uuid = ?", unitParam)
		} else {
			qFu = dbInstance.WithContext(c.Context()).Table("v_fakultas_unit").Where("id = ?", unitParam)
		}
		if errTake := qFu.Take(&fu).Error; errTake == nil {
			targetUnitID = fu.ID
			targetUnitName = fu.Nama

			if strings.EqualFold(fu.Type, "prodi") && fu.Fakultas != "" {
				var parentFu struct {
					ID uint64 `gorm:"column:id"`
				}
				if errP := dbInstance.WithContext(c.Context()).Table("v_fakultas_unit").
					Select("id").
					Where("LOWER(TRIM(nama_fak_prod_unit)) = ? AND LOWER(type) = 'fakultas'", strings.ToLower(strings.TrimSpace(fu.Fakultas))).
					Take(&parentFu).Error; errP == nil {
					parentFakultasID = parentFu.ID
				}
			}
		}
	}

	var item domain.JadwalProkerV3
	found := false
	isUnitSpecific := false

	// 1. Check specific unit schedule
	if targetUnitID > 0 {
		if errFind := dbInstance.WithContext(c.Context()).
			Where("tahun = ? AND id_fakultas_unit = ?", tahun, targetUnitID).
			Order("id DESC").First(&item).Error; errFind == nil {
			found = true
			isUnitSpecific = true
		}
	}

	// 2. If prodi and not found, check parent fakultas schedule
	if !found && parentFakultasID > 0 {
		if errFind := dbInstance.WithContext(c.Context()).
			Where("tahun = ? AND id_fakultas_unit = ?", tahun, parentFakultasID).
			Order("id DESC").First(&item).Error; errFind == nil {
			found = true
			isUnitSpecific = true
		}
	}

	// 3. Fallback to global schedule (id_fakultas_unit IS NULL OR id_fakultas_unit = 0)
	if !found {
		if errFind := dbInstance.WithContext(c.Context()).
			Where("tahun = ? AND (id_fakultas_unit IS NULL OR id_fakultas_unit = 0)", tahun).
			Order("id DESC").First(&item).Error; errFind == nil {
			found = true
			isUnitSpecific = false
		}
	}

	if !found {
		return c.JSON(fiber.Map{
			"has_schedule":     false,
			"is_open":          true,
			"status":           "tidak_ada_jadwal",
			"status_label":     "Tidak Ada Pembatasan Jadwal",
			"tahun":            tahun,
			"is_unit_specific": false,
			"message":          fmt.Sprintf("Belum ada jadwal khusus pengajuan program kerja untuk tahun %s", tahun),
		})
	}

	now := time.Now()
	startOfDay := time.Date(item.TanggalMulai.Year(), item.TanggalMulai.Month(), item.TanggalMulai.Day(), 0, 0, 0, 0, item.TanggalMulai.Location())
	endOfDay := time.Date(item.TanggalAkhir.Year(), item.TanggalAkhir.Month(), item.TanggalAkhir.Day(), 23, 59, 59, 999999999, item.TanggalAkhir.Location())

	formatTanggalIndo := func(t time.Time) string {
		bulan := []string{
			"", "Januari", "Februari", "Maret", "April", "Mei", "Juni",
			"Juli", "Agustus", "September", "Oktober", "November", "Desember",
		}
		m := int(t.Month())
		if m >= 1 && m <= 12 {
			return fmt.Sprintf("%02d %s %d", t.Day(), bulan[m], t.Year())
		}
		return t.Format("02-01-2006")
	}

	isOpen := true
	status := "berjalan"
	statusLabel := "Sedang Dibuka"

	unitSuffix := ""
	if isUnitSpecific && targetUnitName != "" {
		unitSuffix = fmt.Sprintf(" (%s)", targetUnitName)
	}

	message := fmt.Sprintf("Pengajuan program kerja tahun %s%s dibuka dari %s s/d %s", item.Tahun, unitSuffix, formatTanggalIndo(item.TanggalMulai), formatTanggalIndo(item.TanggalAkhir))

	if now.Before(startOfDay) {
		isOpen = false
		status = "belum_dibuka"
		statusLabel = "Belum Dibuka"
		message = fmt.Sprintf("Pengajuan program kerja tahun %s%s belum dibuka. Jadwal pengajuan: %s s/d %s", item.Tahun, unitSuffix, formatTanggalIndo(item.TanggalMulai), formatTanggalIndo(item.TanggalAkhir))
	} else if now.After(endOfDay) {
		isOpen = false
		status = "ditutup"
		statusLabel = "Ditutup"
		message = fmt.Sprintf("Pengajuan program kerja tahun %s%s telah ditutup pada %s", item.Tahun, unitSuffix, formatTanggalIndo(item.TanggalAkhir))
	}

	return c.JSON(fiber.Map{
		"has_schedule":       true,
		"is_open":            isOpen,
		"status":             status,
		"status_label":       statusLabel,
		"tahun":              item.Tahun,
		"tanggal_mulai":      item.TanggalMulai.Format("2006-01-02"),
		"tanggal_akhir":      item.TanggalAkhir.Format("2006-01-02"),
		"id_fakultas_unit":   item.IdFakultasUnit,
		"is_unit_specific":   isUnitSpecific,
		"nama_fakultas_unit": targetUnitName,
		"message":            message,
	})
}

func CreateJadwalProkerV3Handler(c *fiber.Ctx) error {
	if dbInstance == nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Database tidak tersedia"})
	}

	type FormInput struct {
		Tahun          string  `json:"tahun" form:"tahun"`
		TanggalMulai   string  `json:"tanggal_mulai" form:"tanggal_mulai"`
		TanggalAkhir   string  `json:"tanggal_akhir" form:"tanggal_akhir"`
		IdFakultasUnit *string `json:"id_fakultas_unit" form:"id_fakultas_unit"`
	}

	var req FormInput
	if strings.Contains(c.Get("Content-Type"), "application/json") {
		_ = c.BodyParser(&req)
	}
	if req.Tahun == "" {
		req.Tahun = c.FormValue("tahun")
	}
	if req.TanggalMulai == "" {
		req.TanggalMulai = c.FormValue("tanggal_mulai")
	}
	if req.TanggalAkhir == "" {
		req.TanggalAkhir = c.FormValue("tanggal_akhir")
	}

	req.Tahun = strings.TrimSpace(req.Tahun)
	req.TanggalMulai = strings.TrimSpace(req.TanggalMulai)
	req.TanggalAkhir = strings.TrimSpace(req.TanggalAkhir)

	if req.Tahun == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Tahun wajib diisi"})
	}
	if req.TanggalMulai == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Tanggal mulai wajib diisi"})
	}
	if req.TanggalAkhir == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Tanggal akhir wajib diisi"})
	}

	tMulai, errM := time.Parse("2006-01-02", req.TanggalMulai)
	if errM != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Format tanggal mulai harus YYYY-MM-DD"})
	}
	tAkhir, errA := time.Parse("2006-01-02", req.TanggalAkhir)
	if errA != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Format tanggal akhir harus YYYY-MM-DD"})
	}

	if tAkhir.Before(tMulai) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Tanggal akhir tidak boleh lebih awal dari tanggal mulai"})
	}

	var fuIdPtr *uint64
	rawFu := ""
	if req.IdFakultasUnit != nil {
		rawFu = strings.TrimSpace(*req.IdFakultasUnit)
	}
	if rawFu == "" {
		rawFu = strings.TrimSpace(c.FormValue("id_fakultas_unit"))
	}
	if rawFu != "" && rawFu != "0" && rawFu != "all" && rawFu != "null" {
		if val, err := strconv.ParseUint(rawFu, 10, 64); err == nil && val > 0 {
			fuIdPtr = &val
		}
	}

	newJadwal := domain.JadwalProkerV3{
		UUID:           uuid.New(),
		Tahun:          req.Tahun,
		TanggalMulai:   tMulai,
		TanggalAkhir:   tAkhir,
		IdFakultasUnit: fuIdPtr,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	if err := dbInstance.WithContext(c.Context()).Create(&newJadwal).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Jadwal program kerja berhasil ditambahkan",
		"uuid":    newJadwal.UUID.String(),
		"data":    newJadwal,
	})
}

func UpdateJadwalProkerV3Handler(c *fiber.Ctx) error {
	if dbInstance == nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Database tidak tersedia"})
	}

	rawUuid := c.Params("uuid")
	parsedUuid, err := uuid.Parse(rawUuid)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid UUID"})
	}

	var item domain.JadwalProkerV3
	if err := dbInstance.WithContext(c.Context()).Where("uuid = ?", parsedUuid).First(&item).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Jadwal proker tidak ditemukan"})
	}

	type FormInput struct {
		Tahun          string  `json:"tahun" form:"tahun"`
		TanggalMulai   string  `json:"tanggal_mulai" form:"tanggal_mulai"`
		TanggalAkhir   string  `json:"tanggal_akhir" form:"tanggal_akhir"`
		IdFakultasUnit *string `json:"id_fakultas_unit" form:"id_fakultas_unit"`
	}

	var req FormInput
	if strings.Contains(c.Get("Content-Type"), "application/json") {
		_ = c.BodyParser(&req)
	}
	if req.Tahun == "" {
		req.Tahun = c.FormValue("tahun")
	}
	if req.TanggalMulai == "" {
		req.TanggalMulai = c.FormValue("tanggal_mulai")
	}
	if req.TanggalAkhir == "" {
		req.TanggalAkhir = c.FormValue("tanggal_akhir")
	}

	if req.Tahun != "" {
		item.Tahun = strings.TrimSpace(req.Tahun)
	}
	if req.TanggalMulai != "" {
		tMulai, errM := time.Parse("2006-01-02", strings.TrimSpace(req.TanggalMulai))
		if errM != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Format tanggal mulai harus YYYY-MM-DD"})
		}
		item.TanggalMulai = tMulai
	}
	if req.TanggalAkhir != "" {
		tAkhir, errA := time.Parse("2006-01-02", strings.TrimSpace(req.TanggalAkhir))
		if errA != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Format tanggal akhir harus YYYY-MM-DD"})
		}
		item.TanggalAkhir = tAkhir
	}

	rawFu := ""
	hasFuUpdate := false
	if req.IdFakultasUnit != nil {
		rawFu = strings.TrimSpace(*req.IdFakultasUnit)
		hasFuUpdate = true
	} else if c.FormValue("id_fakultas_unit") != "" {
		rawFu = strings.TrimSpace(c.FormValue("id_fakultas_unit"))
		hasFuUpdate = true
	}
	if hasFuUpdate {
		if rawFu == "" || rawFu == "0" || rawFu == "all" || rawFu == "null" {
			item.IdFakultasUnit = nil
		} else if val, err := strconv.ParseUint(rawFu, 10, 64); err == nil && val > 0 {
			item.IdFakultasUnit = &val
		}
	}

	if item.TanggalAkhir.Before(item.TanggalMulai) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Tanggal akhir tidak boleh lebih awal dari tanggal mulai"})
	}

	item.UpdatedAt = time.Now()
	if err := dbInstance.WithContext(c.Context()).Save(&item).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"message": "Jadwal program kerja berhasil diperbarui",
		"uuid":    item.UUID.String(),
		"data":    item,
	})
}

func DeleteJadwalProkerV3Handler(c *fiber.Ctx) error {
	if dbInstance == nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Database tidak tersedia"})
	}

	rawUuid := c.Params("uuid")
	parsedUuid, err := uuid.Parse(rawUuid)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid UUID"})
	}

	if err := dbInstance.WithContext(c.Context()).Where("uuid = ?", parsedUuid).Delete(&domain.JadwalProkerV3{}).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"message": "Jadwal program kerja berhasil dihapus",
		"uuid":    rawUuid,
	})
}

// =======================================================
// JADWAL AKTIVITAS V3 CRUD & STATUS CHECK
// =======================================================

type JadwalAktivitasV3Row struct {
	domain.JadwalAktivitasV3
	NamaFakultasUnit string `json:"nama_fakultas_unit"`
	TypeFakultasUnit string `json:"type_fakultas_unit"`
}

func GetAllJadwalAktivitasV3Handler(c *fiber.Ctx) error {
	if dbInstance == nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Database tidak tersedia"})
	}

	search := strings.TrimSpace(c.Query("search"))
	tahun := strings.TrimSpace(c.Query("tahun"))
	filterUnit := strings.TrimSpace(c.Query("filter_unit"))

	var list []JadwalAktivitasV3Row
	query := dbInstance.WithContext(c.Context()).
		Table("jadwal_aktivitas_v3 j").
		Select("j.*, COALESCE(fu.nama_fak_prod_unit, '') as nama_fakultas_unit, COALESCE(fu.type, '') as type_fakultas_unit").
		Joins("LEFT JOIN v_fakultas_unit fu ON j.id_fakultas_unit = fu.id")

	if tahun != "" && tahun != "all" {
		query = query.Where("j.tahun = ?", tahun)
	}
	if filterUnit == "global" {
		query = query.Where("j.id_fakultas_unit IS NULL OR j.id_fakultas_unit = 0")
	} else if filterUnit == "specific" {
		query = query.Where("j.id_fakultas_unit IS NOT NULL AND j.id_fakultas_unit > 0")
	} else if filterUnit != "" && filterUnit != "all" {
		query = query.Where("j.id_fakultas_unit = ?", filterUnit)
	}
	if search != "" {
		likeSearch := "%" + search + "%"
		query = query.Where("(j.tahun LIKE ? OR fu.nama_fak_prod_unit LIKE ?)", likeSearch, likeSearch)
	}

	if err := query.Order("j.tahun DESC, (j.id_fakultas_unit IS NULL) DESC, j.tanggal_mulai ASC").Find(&list).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	now := time.Now()
	res := make([]domain.JadwalAktivitasV3DTO, 0, len(list))
	for _, item := range list {
		startOfDay := time.Date(item.TanggalMulai.Year(), item.TanggalMulai.Month(), item.TanggalMulai.Day(), 0, 0, 0, 0, item.TanggalMulai.Location())
		endOfDay := time.Date(item.TanggalAkhir.Year(), item.TanggalAkhir.Month(), item.TanggalAkhir.Day(), 23, 59, 59, 999999999, item.TanggalAkhir.Location())

		status := "berjalan"
		statusLabel := "Sedang Dibuka"
		if now.Before(startOfDay) {
			status = "belum_dibuka"
			statusLabel = "Belum Dibuka"
		} else if now.After(endOfDay) {
			status = "ditutup"
			statusLabel = "Ditutup"
		}

		namaUnit := item.NamaFakultasUnit
		typeUnit := item.TypeFakultasUnit
		if item.IdFakultasUnit == nil || *item.IdFakultasUnit == 0 {
			namaUnit = "Semua Unit / Fakultas (Global)"
			typeUnit = "all"
		}

		res = append(res, domain.JadwalAktivitasV3DTO{
			ID:               item.ID,
			UUID:             item.UUID.String(),
			Tahun:            item.Tahun,
			TanggalMulai:     item.TanggalMulai.Format("2006-01-02"),
			TanggalAkhir:     item.TanggalAkhir.Format("2006-01-02"),
			IdFakultasUnit:   item.IdFakultasUnit,
			NamaFakultasUnit: namaUnit,
			TypeFakultasUnit: typeUnit,
			Status:           status,
			StatusLabel:      statusLabel,
			CreatedAt:        item.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt:        item.UpdatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	return c.JSON(fiber.Map{
		"data":  res,
		"total": len(res),
	})
}

func GetJadwalAktivitasV3Handler(c *fiber.Ctx) error {
	if dbInstance == nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Database tidak tersedia"})
	}

	rawUuid := c.Params("uuid")
	parsedUuid, err := uuid.Parse(rawUuid)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid UUID"})
	}

	var item JadwalAktivitasV3Row
	if err := dbInstance.WithContext(c.Context()).
		Table("jadwal_aktivitas_v3 j").
		Select("j.*, COALESCE(fu.nama_fak_prod_unit, '') as nama_fakultas_unit, COALESCE(fu.type, '') as type_fakultas_unit").
		Joins("LEFT JOIN v_fakultas_unit fu ON j.id_fakultas_unit = fu.id").
		Where("j.uuid = ?", parsedUuid).First(&item).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Jadwal aktivitas tidak ditemukan"})
	}

	now := time.Now()
	startOfDay := time.Date(item.TanggalMulai.Year(), item.TanggalMulai.Month(), item.TanggalMulai.Day(), 0, 0, 0, 0, item.TanggalMulai.Location())
	endOfDay := time.Date(item.TanggalAkhir.Year(), item.TanggalAkhir.Month(), item.TanggalAkhir.Day(), 23, 59, 59, 999999999, item.TanggalAkhir.Location())

	status := "berjalan"
	statusLabel := "Sedang Dibuka"
	if now.Before(startOfDay) {
		status = "belum_dibuka"
		statusLabel = "Belum Dibuka"
	} else if now.After(endOfDay) {
		status = "ditutup"
		statusLabel = "Ditutup"
	}

	namaUnit := item.NamaFakultasUnit
	typeUnit := item.TypeFakultasUnit
	if item.IdFakultasUnit == nil || *item.IdFakultasUnit == 0 {
		namaUnit = "Semua Unit / Fakultas (Global)"
		typeUnit = "all"
	}

	return c.JSON(domain.JadwalAktivitasV3DTO{
		ID:               item.ID,
		UUID:             item.UUID.String(),
		Tahun:            item.Tahun,
		TanggalMulai:     item.TanggalMulai.Format("2006-01-02"),
		TanggalAkhir:     item.TanggalAkhir.Format("2006-01-02"),
		IdFakultasUnit:   item.IdFakultasUnit,
		NamaFakultasUnit: namaUnit,
		TypeFakultasUnit: typeUnit,
		Status:           status,
		StatusLabel:      statusLabel,
		CreatedAt:        item.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:        item.UpdatedAt.Format("2006-01-02 15:04:05"),
	})
}

func CheckJadwalAktivitasV3Handler(c *fiber.Ctx) error {
	if dbInstance == nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Database tidak tersedia"})
	}

	tahun := strings.TrimSpace(c.Query("tahun"))
	if tahun == "" {
		tahun = "2027"
	}

	unitParam := strings.TrimSpace(c.Query("id_fakultas_unit"))
	if unitParam == "" {
		unitParam = strings.TrimSpace(c.Query("fakultas_unit"))
	}

	var targetUnitID uint64
	var targetUnitName string
	var parentFakultasID uint64

	if unitParam != "" && unitParam != "all" && unitParam != "0" {
		var fu struct {
			ID       uint64 `gorm:"column:id"`
			Nama     string `gorm:"column:nama_fak_prod_unit"`
			Type     string `gorm:"column:type"`
			Fakultas string `gorm:"column:fakultas"`
		}
		var qFu *gorm.DB
		if _, errU := uuid.Parse(unitParam); errU == nil {
			qFu = dbInstance.WithContext(c.Context()).Table("v_fakultas_unit").Where("uuid = ?", unitParam)
		} else {
			qFu = dbInstance.WithContext(c.Context()).Table("v_fakultas_unit").Where("id = ?", unitParam)
		}
		if errTake := qFu.Take(&fu).Error; errTake == nil {
			targetUnitID = fu.ID
			targetUnitName = fu.Nama

			if strings.EqualFold(fu.Type, "prodi") && fu.Fakultas != "" {
				var parentFu struct {
					ID uint64 `gorm:"column:id"`
				}
				if errP := dbInstance.WithContext(c.Context()).Table("v_fakultas_unit").
					Select("id").
					Where("LOWER(TRIM(nama_fak_prod_unit)) = ? AND LOWER(type) = 'fakultas'", strings.ToLower(strings.TrimSpace(fu.Fakultas))).
					Take(&parentFu).Error; errP == nil {
					parentFakultasID = parentFu.ID
				}
			}
		}
	}

	if targetUnitName == "" && c.Query("nama_fakultas_unit") != "" {
		targetUnitName = strings.TrimSpace(c.Query("nama_fakultas_unit"))
	}

	var item domain.JadwalAktivitasV3
	found := false
	isUnitSpecific := false

	// 1. Check specific unit schedule
	if targetUnitID > 0 {
		if errFind := dbInstance.WithContext(c.Context()).
			Where("tahun = ? AND id_fakultas_unit = ?", tahun, targetUnitID).
			Order("id DESC").First(&item).Error; errFind == nil {
			found = true
			isUnitSpecific = true
		}
	}

	// 2. If prodi and not found, check parent fakultas schedule
	if !found && parentFakultasID > 0 {
		if errFind := dbInstance.WithContext(c.Context()).
			Where("tahun = ? AND id_fakultas_unit = ?", tahun, parentFakultasID).
			Order("id DESC").First(&item).Error; errFind == nil {
			found = true
			isUnitSpecific = true
		}
	}

	// 3. Fallback to global schedule (id_fakultas_unit IS NULL OR id_fakultas_unit = 0)
	if !found {
		if errFind := dbInstance.WithContext(c.Context()).
			Where("tahun = ? AND (id_fakultas_unit IS NULL OR id_fakultas_unit = 0)", tahun).
			Order("id DESC").First(&item).Error; errFind == nil {
			found = true
			isUnitSpecific = false
		}
	}

	if !found {
		unitSuffix := ""
		if targetUnitName != "" {
			unitSuffix = fmt.Sprintf(" (%s)", targetUnitName)
		}
		return c.JSON(fiber.Map{
			"has_schedule":       false,
			"is_open":            false,
			"status":             "ditutup",
			"status_label":       "Ditutup",
			"tahun":              tahun,
			"is_unit_specific":   false,
			"nama_fakultas_unit": targetUnitName,
			"message":            fmt.Sprintf("Jadwal pembuatan aktivitas kegiatan tahun %s%s belum dibuka atau belum diatur oleh admin", tahun, unitSuffix),
		})
	}

	now := time.Now()
	startOfDay := time.Date(item.TanggalMulai.Year(), item.TanggalMulai.Month(), item.TanggalMulai.Day(), 0, 0, 0, 0, item.TanggalMulai.Location())
	endOfDay := time.Date(item.TanggalAkhir.Year(), item.TanggalAkhir.Month(), item.TanggalAkhir.Day(), 23, 59, 59, 999999999, item.TanggalAkhir.Location())

	formatTanggalIndo := func(t time.Time) string {
		bulan := []string{
			"", "Januari", "Februari", "Maret", "April", "Mei", "Juni",
			"Juli", "Agustus", "September", "Oktober", "November", "Desember",
		}
		m := int(t.Month())
		if m >= 1 && m <= 12 {
			return fmt.Sprintf("%02d %s %d", t.Day(), bulan[m], t.Year())
		}
		return t.Format("02-01-2006")
	}

	isOpen := true
	status := "berjalan"
	statusLabel := "Sedang Dibuka"

	unitSuffix := ""
	if targetUnitName != "" {
		unitSuffix = fmt.Sprintf(" (%s)", targetUnitName)
	}

	message := fmt.Sprintf("Pembuatan aktivitas kegiatan tahun %s%s dibuka dari %s s/d %s", item.Tahun, unitSuffix, formatTanggalIndo(item.TanggalMulai), formatTanggalIndo(item.TanggalAkhir))

	if now.Before(startOfDay) {
		isOpen = false
		status = "belum_dibuka"
		statusLabel = "Belum Dibuka"
		message = fmt.Sprintf("Pembuatan aktivitas kegiatan tahun %s%s belum dibuka. Jadwal: %s s/d %s", item.Tahun, unitSuffix, formatTanggalIndo(item.TanggalMulai), formatTanggalIndo(item.TanggalAkhir))
	} else if now.After(endOfDay) {
		isOpen = false
		status = "ditutup"
		statusLabel = "Ditutup"
		message = fmt.Sprintf("Pembuatan aktivitas kegiatan tahun %s%s telah ditutup pada %s", item.Tahun, unitSuffix, formatTanggalIndo(item.TanggalAkhir))
	}

	return c.JSON(fiber.Map{
		"has_schedule":       true,
		"is_open":            isOpen,
		"status":             status,
		"status_label":       statusLabel,
		"tahun":              item.Tahun,
		"tanggal_mulai":      item.TanggalMulai.Format("2006-01-02"),
		"tanggal_akhir":      item.TanggalAkhir.Format("2006-01-02"),
		"id_fakultas_unit":   item.IdFakultasUnit,
		"is_unit_specific":   isUnitSpecific,
		"nama_fakultas_unit": targetUnitName,
		"message":            message,
	})
}

func CreateJadwalAktivitasV3Handler(c *fiber.Ctx) error {
	if dbInstance == nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Database tidak tersedia"})
	}

	type FormInput struct {
		Tahun          string  `json:"tahun" form:"tahun"`
		TanggalMulai   string  `json:"tanggal_mulai" form:"tanggal_mulai"`
		TanggalAkhir   string  `json:"tanggal_akhir" form:"tanggal_akhir"`
		IdFakultasUnit *string `json:"id_fakultas_unit" form:"id_fakultas_unit"`
	}

	var req FormInput
	if strings.Contains(c.Get("Content-Type"), "application/json") {
		_ = c.BodyParser(&req)
	}
	if req.Tahun == "" {
		req.Tahun = c.FormValue("tahun")
	}
	if req.TanggalMulai == "" {
		req.TanggalMulai = c.FormValue("tanggal_mulai")
	}
	if req.TanggalAkhir == "" {
		req.TanggalAkhir = c.FormValue("tanggal_akhir")
	}

	req.Tahun = strings.TrimSpace(req.Tahun)
	req.TanggalMulai = strings.TrimSpace(req.TanggalMulai)
	req.TanggalAkhir = strings.TrimSpace(req.TanggalAkhir)

	if req.Tahun == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Tahun wajib diisi"})
	}
	if req.TanggalMulai == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Tanggal mulai wajib diisi"})
	}
	if req.TanggalAkhir == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Tanggal akhir wajib diisi"})
	}

	tMulai, errM := time.Parse("2006-01-02", req.TanggalMulai)
	if errM != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Format tanggal mulai harus YYYY-MM-DD"})
	}
	tAkhir, errA := time.Parse("2006-01-02", req.TanggalAkhir)
	if errA != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Format tanggal akhir harus YYYY-MM-DD"})
	}

	if tAkhir.Before(tMulai) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Tanggal akhir tidak boleh lebih awal dari tanggal mulai"})
	}

	var fuIdPtr *uint64
	rawFu := ""
	if req.IdFakultasUnit != nil {
		rawFu = strings.TrimSpace(*req.IdFakultasUnit)
	}
	if rawFu == "" {
		rawFu = strings.TrimSpace(c.FormValue("id_fakultas_unit"))
	}
	if rawFu != "" && rawFu != "0" && rawFu != "all" && rawFu != "null" {
		if val, err := strconv.ParseUint(rawFu, 10, 64); err == nil && val > 0 {
			fuIdPtr = &val
		}
	}

	newJadwal := domain.JadwalAktivitasV3{
		UUID:           uuid.New(),
		Tahun:          req.Tahun,
		TanggalMulai:   tMulai,
		TanggalAkhir:   tAkhir,
		IdFakultasUnit: fuIdPtr,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	if err := dbInstance.WithContext(c.Context()).Create(&newJadwal).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Jadwal aktivitas berhasil ditambahkan",
		"uuid":    newJadwal.UUID.String(),
		"data":    newJadwal,
	})
}

func UpdateJadwalAktivitasV3Handler(c *fiber.Ctx) error {
	if dbInstance == nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Database tidak tersedia"})
	}

	rawUuid := c.Params("uuid")
	parsedUuid, err := uuid.Parse(rawUuid)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid UUID"})
	}

	var item domain.JadwalAktivitasV3
	if err := dbInstance.WithContext(c.Context()).Where("uuid = ?", parsedUuid).First(&item).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Jadwal aktivitas tidak ditemukan"})
	}

	type FormInput struct {
		Tahun          string  `json:"tahun" form:"tahun"`
		TanggalMulai   string  `json:"tanggal_mulai" form:"tanggal_mulai"`
		TanggalAkhir   string  `json:"tanggal_akhir" form:"tanggal_akhir"`
		IdFakultasUnit *string `json:"id_fakultas_unit" form:"id_fakultas_unit"`
	}

	var req FormInput
	if strings.Contains(c.Get("Content-Type"), "application/json") {
		_ = c.BodyParser(&req)
	}
	if req.Tahun == "" {
		req.Tahun = c.FormValue("tahun")
	}
	if req.TanggalMulai == "" {
		req.TanggalMulai = c.FormValue("tanggal_mulai")
	}
	if req.TanggalAkhir == "" {
		req.TanggalAkhir = c.FormValue("tanggal_akhir")
	}

	if req.Tahun != "" {
		item.Tahun = strings.TrimSpace(req.Tahun)
	}
	if req.TanggalMulai != "" {
		tMulai, errM := time.Parse("2006-01-02", strings.TrimSpace(req.TanggalMulai))
		if errM != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Format tanggal mulai harus YYYY-MM-DD"})
		}
		item.TanggalMulai = tMulai
	}
	if req.TanggalAkhir != "" {
		tAkhir, errA := time.Parse("2006-01-02", strings.TrimSpace(req.TanggalAkhir))
		if errA != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Format tanggal akhir harus YYYY-MM-DD"})
		}
		item.TanggalAkhir = tAkhir
	}

	rawFu := ""
	hasFuUpdate := false
	if req.IdFakultasUnit != nil {
		rawFu = strings.TrimSpace(*req.IdFakultasUnit)
		hasFuUpdate = true
	} else if c.FormValue("id_fakultas_unit") != "" {
		rawFu = strings.TrimSpace(c.FormValue("id_fakultas_unit"))
		hasFuUpdate = true
	}

	if hasFuUpdate {
		if rawFu == "" || rawFu == "0" || rawFu == "all" || rawFu == "null" {
			item.IdFakultasUnit = nil
		} else if val, err := strconv.ParseUint(rawFu, 10, 64); err == nil && val > 0 {
			item.IdFakultasUnit = &val
		}
	}

	if item.TanggalAkhir.Before(item.TanggalMulai) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Tanggal akhir tidak boleh lebih awal dari tanggal mulai"})
	}

	item.UpdatedAt = time.Now()
	if err := dbInstance.WithContext(c.Context()).Save(&item).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"message": "Jadwal aktivitas berhasil diperbarui",
		"uuid":    item.UUID.String(),
		"data":    item,
	})
}

func DeleteJadwalAktivitasV3Handler(c *fiber.Ctx) error {
	if dbInstance == nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Database tidak tersedia"})
	}

	rawUuid := c.Params("uuid")
	parsedUuid, err := uuid.Parse(rawUuid)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid UUID"})
	}

	if err := dbInstance.WithContext(c.Context()).Where("uuid = ?", parsedUuid).Delete(&domain.JadwalAktivitasV3{}).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"message": "Jadwal aktivitas berhasil dihapus",
		"uuid":    rawUuid,
	})
}

// =======================================================
// PAGU ANGGARAN HANDLERS
// =======================================================

func GetAllPaguAnggaranHandler(c *fiber.Ctx) error {
	if dbInstance == nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Database tidak tersedia"})
	}

	tahun := c.Query("tahun")
	idFu := c.Query("id_fakultas_unit")

	var list []domain.PaguAnggaran
	query := dbInstance.WithContext(c.Context()).Table("pagu_anggaran")
	if tahun != "" {
		query = query.Where("tahun = ?", tahun)
	}
	if idFu != "" && idFu != "0" {
		query = query.Where("id_fakultas_unit = ?", idFu)
	}

	if err := query.Order("id DESC").Find(&list).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"data": list,
	})
}

func SavePaguAnggaranHandler(c *fiber.Ctx) error {
	if dbInstance == nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Database tidak tersedia"})
	}

	var req struct {
		ID             uint     `json:"id"`
		IdFakultasUnit *int     `json:"id_fakultas_unit"`
		FakultasUnit   string   `json:"fakultas_unit"`
		Tahun          string   `json:"tahun"`
		Pagu           *float64 `json:"pagu"`
	}

	if strings.Contains(c.Get("Content-Type"), "application/json") {
		_ = c.BodyParser(&req)
	}

	if req.Tahun == "" {
		req.Tahun = c.FormValue("tahun")
	}
	if req.FakultasUnit == "" {
		req.FakultasUnit = c.FormValue("fakultas_unit")
	}
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

	if req.Tahun == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Tahun wajib diisi"})
	}

	nom := float64(0)
	if req.Pagu != nil {
		nom = *req.Pagu
	}

	if req.IdFakultasUnit != nil && *req.IdFakultasUnit > 0 && req.FakultasUnit == "" {
		var fuName string
		_ = dbInstance.WithContext(c.Context()).Table("v_fakultas_unit").
			Select("nama_fak_prod_unit").
			Where("id = ?", *req.IdFakultasUnit).
			Scan(&fuName).Error
		req.FakultasUnit = fuName
	}

	var item domain.PaguAnggaran
	if req.ID > 0 {
		if err := dbInstance.WithContext(c.Context()).Where("id = ?", req.ID).First(&item).Error; err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Data pagu tidak ditemukan"})
		}
		item.IdFakultasUnit = req.IdFakultasUnit
		item.FakultasUnit = req.FakultasUnit
		item.Tahun = req.Tahun
		item.Pagu = nom
		item.UpdatedAt = time.Now()
		if err := dbInstance.WithContext(c.Context()).Save(&item).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
	} else {
		var existing domain.PaguAnggaran
		checkQuery := dbInstance.WithContext(c.Context()).Where("tahun = ?", req.Tahun)
		if req.IdFakultasUnit != nil {
			checkQuery = checkQuery.Where("id_fakultas_unit = ?", *req.IdFakultasUnit)
		} else if req.FakultasUnit != "" {
			checkQuery = checkQuery.Where("fakultas_unit = ?", req.FakultasUnit)
		}
		if err := checkQuery.First(&existing).Error; err == nil {
			existing.Pagu = nom
			existing.UpdatedAt = time.Now()
			if req.FakultasUnit != "" {
				existing.FakultasUnit = req.FakultasUnit
			}
			if req.IdFakultasUnit != nil {
				existing.IdFakultasUnit = req.IdFakultasUnit
			}
			if err := dbInstance.WithContext(c.Context()).Save(&existing).Error; err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
			}
			item = existing
		} else {
			item = domain.PaguAnggaran{
				IdFakultasUnit: req.IdFakultasUnit,
				FakultasUnit:   req.FakultasUnit,
				Tahun:          req.Tahun,
				Pagu:           nom,
				CreatedAt:      time.Now(),
				UpdatedAt:      time.Now(),
			}
			if err := dbInstance.WithContext(c.Context()).Create(&item).Error; err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
			}
		}
	}

	return c.JSON(fiber.Map{
		"message": "Pagu anggaran berhasil disimpan",
		"data":    item,
	})
}

func DeletePaguAnggaranHandler(c *fiber.Ctx) error {
	if dbInstance == nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Database tidak tersedia"})
	}

	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "ID wajib diisi"})
	}

	if err := dbInstance.WithContext(c.Context()).Where("id = ?", id).Delete(&domain.PaguAnggaran{}).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"message": "Pagu anggaran berhasil dihapus",
		"id":      id,
	})
}

// =======================================================
// GET /proker_v3/unit_pagu_summary
// =======================================================
func GetUnitPaguSummaryHandler(c *fiber.Ctx) error {
	if dbInstance == nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Database tidak tersedia"})
	}

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
	if reqFu == "" || reqFu == "0" || reqFu == "all" {
		reqFu = userFU
	}

	var fuId int
	var fuName string

	if reqFu != "" {
		var fuRow struct {
			ID              int    `gorm:"column:id"`
			NamaFakProdUnit string `gorm:"column:nama_fak_prod_unit"`
		}
		_ = dbInstance.WithContext(c.Context()).Table("v_fakultas_unit").
			Where("id = ? OR nama_fak_prod_unit = ?", reqFu, reqFu).
			First(&fuRow).Error

		if fuRow.ID > 0 {
			fuId = fuRow.ID
			fuName = fuRow.NamaFakProdUnit
		} else if parsed, err := strconv.Atoi(reqFu); err == nil {
			fuId = parsed
		}
	}

	if fuId == 0 && userFU != "" {
		if parsed, err := strconv.Atoi(userFU); err == nil {
			fuId = parsed
		}
	}

	// 1. Get Pagu Anggaran
	var paguRow struct {
		Pagu float64
	}
	errPagu := dbInstance.WithContext(c.Context()).Table("pagu_anggaran").
		Select("pagu").
		Where("tahun = ? AND (id_fakultas_unit = ? OR (fakultas_unit IS NOT NULL AND fakultas_unit != '' AND (fakultas_unit = ? OR fakultas_unit = ?)))", tahun, fuId, strconv.Itoa(fuId), fuName).
		Order("id DESC").
		Limit(1).
		Scan(&paguRow).Error

	hasPagu := errPagu == nil && paguRow.Pagu > 0
	paguVal := float64(0)
	if hasPagu {
		paguVal = paguRow.Pagu
	}

	// 2. Query Prokers for the unit and year to calculate cumulative total_anggaran & total_realisasi
	var prokers []struct {
		ID             uint    `gorm:"column:id"`
		TotAnggaran    float64 `gorm:"column:tot_anggaran"`
		TotRealisasi   float64 `gorm:"column:tot_realisasi"`
		CountRealisasi int     `gorm:"column:count_realisasi"`
	}

	_ = dbInstance.WithContext(c.Context()).Table("matrix_proker mp").
		Select(`
			mp.id,
			COALESCE((SELECT SUM(COALESCE(r.total, r.qty_satuan * r.qty_paket * r.sbu, 0)) FROM aktivitasv2_rab r JOIN aktivitasv2 a ON a.id = r.id_aktivitasv2 WHERE a.id_matrik_proker = mp.id), 0) as tot_anggaran,
			COALESCE((SELECT SUM(CASE WHEN a.realisasi IS NOT NULL THEN a.realisasi ELSE COALESCE((SELECT SUM(COALESCE(r2.total, r2.qty_satuan * r2.qty_paket * r2.sbu, 0)) FROM aktivitasv2_rab r2 WHERE r2.id_aktivitasv2 = a.id), 0) END) FROM aktivitasv2 a WHERE a.id_matrik_proker = mp.id), 0) as tot_realisasi,
			COALESCE((SELECT COUNT(*) FROM aktivitasv2 a WHERE a.id_matrik_proker = mp.id AND a.realisasi IS NOT NULL), 0) as count_realisasi
		`).
		Joins("LEFT JOIN mata_program p ON p.id = mp.id_mata_program").
		Joins("LEFT JOIN master_tahun mt ON mt.id = p.id_master_tahun").
		Where("(mp.id_fakultas_unit = ? OR (SELECT fu2.nama_fak_prod_unit FROM v_fakultas_unit fu2 WHERE fu2.id = mp.id_fakultas_unit) = ?) AND COALESCE(mt.tahun, '') = ?", fuId, fuName, tahun).
		Scan(&prokers).Error

	totalAnggaran := float64(0)
	totalRealisasi := float64(0)
	hasAnyRealisasi := false

	for _, p := range prokers {
		totalAnggaran += p.TotAnggaran
		totalRealisasi += p.TotRealisasi
		if p.CountRealisasi > 0 {
			hasAnyRealisasi = true
		}
	}

	hasRealDiff := hasAnyRealisasi && (totalRealisasi != totalAnggaran)
	effectiveTotal := totalRealisasi
	if !hasAnyRealisasi {
		effectiveTotal = totalAnggaran
	}
	isMelebihiPagu := hasPagu && (effectiveTotal > paguVal)
	sisaPagu := float64(0)
	if hasPagu {
		sisaPagu = paguVal - effectiveTotal
	}

	return c.JSON(fiber.Map{
		"has_pagu":           hasPagu,
		"pagu":               paguVal,
		"total_anggaran":     totalAnggaran,
		"total_realisasi":    totalRealisasi,
		"has_realisasi":      hasAnyRealisasi,
		"has_real_diff":      hasRealDiff,
		"is_melebihi_pagu":   isMelebihiPagu,
		"sisa_pagu":          sisaPagu,
		"nama_fakultas_unit": fuName,
		"id_fakultas_unit":   fuId,
		"tahun":              tahun,
	})
}

// =======================================================
// DASHBOARD MONITORING FOR WAREK 1-4 & FAKULTAS
// =======================================================

type DashboardMonitoringUnitOption struct {
	ID   interface{} `json:"id"`
	Nama string      `json:"nama"`
	Type string      `json:"type"`
}

type DashboardMonitoringRabItem struct {
	ID         uint    `json:"id"`
	UUID       string  `json:"uuid"`
	NoAkun     string  `json:"no_akun"`
	NamaAkun   string  `json:"nama_akun"`
	Kategori   string  `json:"kategori"`
	Rincian    string  `json:"rincian"`
	Keterangan string  `json:"keterangan"`
	QtySatuan  float64 `json:"qty_satuan"`
	QtyPaket   float64 `json:"qty_paket"`
	Sbu        float64 `json:"sbu"`
	Total      float64 `json:"total"`
}

type DashboardMonitoringAktivitas struct {
	ID              uint                         `json:"id"`
	UUID            string                       `json:"uuid"`
	Aktivitas       string                       `json:"aktivitas"`
	PIC             string                       `json:"pic"`
	Capaian         string                       `json:"capaian"`
	Realisasi       *float64                     `json:"realisasi"`
	Pencairan       *float64                     `json:"pencairan"`
	Status          string                       `json:"status"`
	StatusAnggaran  string                       `json:"status_anggaran"`
	StatusUpml      string                       `json:"status_upml"`
	StatusUpmf      string                       `json:"status_upmf"`
	StatusUpms      string                       `json:"status_upms"`
	StatusLpm       string                       `json:"status_lpm"`
	IsBatal         bool                         `json:"is_batal"`
	CatatanBatal    string                       `json:"catatan_batal"`
	TotalRab        float64                      `json:"total_rab"`
	StatusPencairan string                       `json:"status_pencairan"`
	Rab             []DashboardMonitoringRabItem `json:"rab"`
}

type DashboardMonitoringProker struct {
	ID                uint                           `json:"id"`
	UUID              string                         `json:"uuid"`
	IdFakultasUnit    *int                           `json:"id_fakultas_unit"`
	NamaFakultasUnit  string                         `json:"nama_fakultas_unit"`
	TypeFakultasUnit  string                         `json:"type_fakultas_unit"`
	Jenjang           string                         `json:"jenjang"`
	NamaProgram       string                         `json:"nama_program"`
	Tahun             string                         `json:"tahun"`
	IndikatorKinerja  string                         `json:"indikator_kinerja"`
	TargetUnpak       string                         `json:"target_unpak"`
	TargetUnit        string                         `json:"target_unit"`
	Output            string                         `json:"output"`
	Outcome           string                         `json:"outcome"`
	PIC               string                         `json:"pic"`
	Status            string                         `json:"status"`
	StatusVerifikator string                         `json:"status_verifikator"`
	StatusAnggaran    string                         `json:"status_anggaran"`
	BebanYayasan      bool                           `json:"beban_yayasan"`
	IsBatal           bool                           `json:"is_batal"`
	CatatanBatal      string                         `json:"catatan_batal"`
	TotalAktivitas    int                            `json:"total_aktivitas"`
	TotalAnggaran     float64                        `json:"total_anggaran"`
	TotalPencairan    float64                        `json:"total_pencairan"`
	Aktivitas         []DashboardMonitoringAktivitas `json:"aktivitas"`
}

func GetDashboardProkerMonitoringHandler(c *fiber.Ctx) error {
	if dbInstance == nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Database tidak tersedia"})
	}

	sid := c.FormValue("sid")
	level, userFU := commonpresentation.GetUserInfo(c)

	var userRow struct {
		ID           uint   `gorm:"column:id"`
		UUID         string `gorm:"column:uuid"`
		NidnUsername string `gorm:"column:nidn_username"`
		Level        string `gorm:"column:level"`
		Name         string `gorm:"column:name"`
		FakultasUnit string `gorm:"column:fakultas_unit"`
	}

	if sid != "" {
		_ = dbInstance.WithContext(c.Context()).Table("users").Where("uuid = ?", sid).First(&userRow).Error
	}
	if userRow.ID == 0 && userFU != "" {
		_ = dbInstance.WithContext(c.Context()).Table("users").Where("fakultas_unit = ? AND level = ?", userFU, level).First(&userRow).Error
	}

	userName := strings.TrimSpace(userRow.NidnUsername)
	userLevel := strings.ToLower(strings.TrimSpace(userRow.Level))
	if level != "" {
		userLevel = strings.ToLower(strings.TrimSpace(level))
	}
	fuID := strings.TrimSpace(userRow.FakultasUnit)
	if fuID == "" {
		fuID = strings.TrimSpace(userFU)
	}

	isAdmin := userLevel == "admin" || strings.EqualFold(userName, "admin")
	isWarek := strings.Contains(strings.ToLower(userName), "wakil rektor") ||
		strings.Contains(strings.ToLower(userName), "warek") ||
		strings.Contains(userLevel, "wakil rektor") ||
		strings.Contains(userLevel, "warek") ||
		fuID == "65" || fuID == "66" || fuID == "67" || fuID == "122"

	var fuInfo struct {
		ID              int    `gorm:"column:id"`
		NamaFakProdUnit string `gorm:"column:nama_fak_prod_unit"`
		Type            string `gorm:"column:type"`
		KodeFakultas    string `gorm:"column:kode_fakultas"`
	}
	if fuID != "" {
		_ = dbInstance.WithContext(c.Context()).Table("v_fakultas_unit").
			Where("id = ? OR nama_fak_prod_unit = ?", fuID, fuID).
			First(&fuInfo).Error
	}

	isFakultas := fuInfo.Type == "fakultas" || userLevel == "fakultas"

	if !isAdmin && !isWarek && !isFakultas {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"allowed": false,
			"message": "Fitur monitoring ini hanya dapat diakses oleh Wakil Rektor 1-4 dan Tingkat Fakultas",
		})
	}

	userRoleType := "warek"
	if isFakultas && !isWarek && !isAdmin {
		userRoleType = "fakultas"
	} else if isAdmin {
		userRoleType = "admin"
	}

	// Options for unit search dropdown
	var unitOptions []DashboardMonitoringUnitOption
	if isAdmin || isWarek {
		unitOptions = append(unitOptions, DashboardMonitoringUnitOption{
			ID:   "all",
			Nama: "Semua Fakultas & Unit",
			Type: "all",
		})
		var allUnits []struct {
			ID              int    `gorm:"column:id"`
			NamaFakProdUnit string `gorm:"column:nama_fak_prod_unit"`
			Type            string `gorm:"column:type"`
			Jenjang         string `gorm:"column:jenjang"`
		}
		_ = dbInstance.WithContext(c.Context()).Table("v_fakultas_unit").
			Order("FIELD(type, 'fakultas', 'prodi', 'unit'), nama_fak_prod_unit ASC").
			Scan(&allUnits).Error

		for _, u := range allUnits {
			label := u.NamaFakProdUnit
			if u.Type == "prodi" && u.Jenjang != "" {
				label = fmt.Sprintf("%s (%s)", u.NamaFakProdUnit, strings.ToUpper(u.Jenjang))
			}
			unitOptions = append(unitOptions, DashboardMonitoringUnitOption{
				ID:   u.ID,
				Nama: label,
				Type: u.Type,
			})
		}
	} else if isFakultas {
		facultyName := fuInfo.NamaFakProdUnit
		if facultyName == "" {
			facultyName = "Fakultas"
		}
		unitOptions = append(unitOptions, DashboardMonitoringUnitOption{
			ID:   "all",
			Nama: fmt.Sprintf("Semua Prodi (%s)", facultyName),
			Type: "all",
		})
		var facUnits []struct {
			ID              int    `gorm:"column:id"`
			NamaFakProdUnit string `gorm:"column:nama_fak_prod_unit"`
			Type            string `gorm:"column:type"`
			Jenjang         string `gorm:"column:jenjang"`
		}
		_ = dbInstance.WithContext(c.Context()).Table("v_fakultas_unit").
			Where("kode_fakultas = ?", fuInfo.KodeFakultas).
			Order("FIELD(type, 'fakultas', 'prodi'), nama_fak_prod_unit ASC").
			Scan(&facUnits).Error

		for _, u := range facUnits {
			label := u.NamaFakProdUnit
			if u.Type == "prodi" && u.Jenjang != "" {
				label = fmt.Sprintf("%s (%s)", u.NamaFakProdUnit, strings.ToUpper(u.Jenjang))
			}
			unitOptions = append(unitOptions, DashboardMonitoringUnitOption{
				ID:   u.ID,
				Nama: label,
				Type: u.Type,
			})
		}
	}

	tahun := strings.TrimSpace(c.Query("tahun"))
	if tahun == "" {
		var activeTahun string
		_ = dbInstance.WithContext(c.Context()).Table("master_tahun").
			Where("status = 'active'").
			Order("tahun DESC").
			Pluck("tahun", &activeTahun).Error
		if activeTahun != "" {
			tahun = activeTahun
		} else {
			tahun = "2027"
		}
	}
	selectedUnit := strings.TrimSpace(c.Query("id_fakultas_unit"))
	search := strings.TrimSpace(c.Query("search"))

	query := dbInstance.WithContext(c.Context()).Table("matrix_proker mp").
		Select(`
			mp.id as id,
			mp.uuid as uuid,
			mp.id_fakultas_unit as id_fakultas_unit,
			COALESCE(fu.nama_fak_prod_unit, '') as nama_fakultas_unit,
			COALESCE(fu.type, '') as type_fakultas_unit,
			COALESCE(fu.jenjang, '') as jenjang,
			COALESCE(prog.mata_program, '') as nama_program,
			COALESCE(mt.tahun, '') as tahun,
			COALESCE(mir.indikator, '') as indikator_kinerja,
			COALESCE(mp.target_unpak, '') as target_unpak,
			COALESCE(mp.target_unit, '') as target_unit,
			COALESCE(mp.output, '') as output,
			COALESCE(mp.outcome, '') as outcome,
			COALESCE(NULLIF(mp.pic, ''), '') as pic,
			COALESCE(NULLIF(mp.status, ''), 'draft') as status,
			COALESCE(NULLIF(mp.status_verifikator, ''), 'draft') as status_verifikator,
			COALESCE(NULLIF(mp.status_anggaran, ''), 'menunggu_verifikator') as status_anggaran,
			COALESCE(mp.beban_yayasan, 0) as beban_yayasan,
			COALESCE(mp.is_batal, 0) as is_batal,
			COALESCE(mp.catatan_batal, '') as catatan_batal
		`).
		Joins("LEFT JOIN mata_program prog ON mp.id_mata_program = prog.id").
		Joins("LEFT JOIN master_tahun mt ON prog.id_master_tahun = mt.id").
		Joins("LEFT JOIN master_indikator_renstra mir ON mp.id_indikator = mir.id").
		Joins("LEFT JOIN v_fakultas_unit fu ON mp.id_fakultas_unit = fu.id").
		Where("COALESCE(mt.tahun, '') = ?", tahun)

	if selectedUnit != "" && selectedUnit != "all" {
		query = query.Where("mp.id_fakultas_unit = ?", selectedUnit)
	} else if isFakultas && !isAdmin && !isWarek && fuInfo.KodeFakultas != "" {
		query = query.Where("mp.id_fakultas_unit IN (SELECT id FROM v_fakultas_unit WHERE kode_fakultas = ?)", fuInfo.KodeFakultas)
	}

	if search != "" {
		searchPattern := "%" + search + "%"
		query = query.Where("(prog.mata_program LIKE ? OR fu.nama_fak_prod_unit LIKE ? OR mir.indikator LIKE ? OR mp.pic LIKE ?)", searchPattern, searchPattern, searchPattern, searchPattern)
	}

	var prokerRows []struct {
		Id                uint   `gorm:"column:id"`
		UUID              string `gorm:"column:uuid"`
		IdFakultasUnit    *int   `gorm:"column:id_fakultas_unit"`
		NamaFakultasUnit  string `gorm:"column:nama_fakultas_unit"`
		TypeFakultasUnit  string `gorm:"column:type_fakultas_unit"`
		Jenjang           string `gorm:"column:jenjang"`
		NamaProgram       string `gorm:"column:nama_program"`
		Tahun             string `gorm:"column:tahun"`
		IndikatorKinerja  string `gorm:"column:indikator_kinerja"`
		TargetUnpak       string `gorm:"column:target_unpak"`
		TargetUnit        string `gorm:"column:target_unit"`
		Output            string `gorm:"column:output"`
		Outcome           string `gorm:"column:outcome"`
		PIC               string `gorm:"column:pic"`
		Status            string `gorm:"column:status"`
		StatusVerifikator string `gorm:"column:status_verifikator"`
		StatusAnggaran    string `gorm:"column:status_anggaran"`
		BebanYayasan      int    `gorm:"column:beban_yayasan"`
		IsBatal           int    `gorm:"column:is_batal"`
		CatatanBatal      string `gorm:"column:catatan_batal"`
	}

	if err := query.Order("mp.id DESC").Scan(&prokerRows).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	var prokerIDs []uint
	for _, p := range prokerRows {
		prokerIDs = append(prokerIDs, p.Id)
	}

	// Fetch all activities in batch
	var rawAktivitas []struct {
		Id             uint     `gorm:"column:id"`
		UUID           string   `gorm:"column:uuid"`
		IdMatrikProker uint     `gorm:"column:id_matrik_proker"`
		Aktivitas      string   `gorm:"column:aktivitas"`
		PIC            string   `gorm:"column:pic"`
		Capaian        string   `gorm:"column:capaian"`
		Realisasi      *float64 `gorm:"column:realisasi"`
		Pencairan      *float64 `gorm:"column:pencairan"`
		Status         string   `gorm:"column:status"`
		StatusAnggaran string   `gorm:"column:status_anggaran"`
		StatusUpml     string   `gorm:"column:status_upml"`
		StatusUpmf     string   `gorm:"column:status_upmf"`
		StatusUpms     string   `gorm:"column:status_upms"`
		StatusLpm      string   `gorm:"column:status_lpm"`
		IsBatal        int      `gorm:"column:is_batal"`
		CatatanBatal   string   `gorm:"column:catatan_batal"`
	}

	var aktivitasIDs []uint
	if len(prokerIDs) > 0 {
		_ = dbInstance.WithContext(c.Context()).Table("aktivitasv2").
			Select("id, uuid, id_matrik_proker, aktivitas, pic, capaian, realisasi, pencairan, status, status_anggaran, status_upml, status_upmf, status_upms, status_lpm, is_batal, COALESCE(catatan_batal, '') as catatan_batal").
			Where("id_matrik_proker IN (?)", prokerIDs).
			Order("id ASC").
			Scan(&rawAktivitas).Error

		for _, a := range rawAktivitas {
			aktivitasIDs = append(aktivitasIDs, a.Id)
		}
	}

	// Fetch all RAB items in batch
	var rawRAB []struct {
		Id            uint    `gorm:"column:id"`
		UUID          string  `gorm:"column:uuid"`
		IdAktivitasv2 uint    `gorm:"column:id_aktivitasv2"`
		NoAkun        string  `gorm:"column:no_akun"`
		NamaAkun      string  `gorm:"column:nama_akun"`
		Kategori      string  `gorm:"column:kategori"`
		Rincian       string  `gorm:"column:rincian"`
		Keterangan    string  `gorm:"column:keterangan"`
		QtySatuan     float64 `gorm:"column:qty_satuan"`
		QtyPaket      float64 `gorm:"column:qty_paket"`
		Sbu           float64 `gorm:"column:sbu"`
		Total         float64 `gorm:"column:total"`
	}

	if len(aktivitasIDs) > 0 {
		_ = dbInstance.WithContext(c.Context()).Table("aktivitasv2_rab").
			Where("id_aktivitasv2 IN (?)", aktivitasIDs).
			Order("id ASC").
			Scan(&rawRAB).Error
	}

	// Map RAB to Aktivitas
	rabMap := make(map[uint][]DashboardMonitoringRabItem)
	rabTotalMap := make(map[uint]float64)
	for _, r := range rawRAB {
		tot := r.Total
		if tot == 0 && r.QtySatuan > 0 && r.QtyPaket > 0 && r.Sbu > 0 {
			tot = r.QtySatuan * r.QtyPaket * r.Sbu
		}
		item := DashboardMonitoringRabItem{
			ID:         r.Id,
			UUID:       r.UUID,
			NoAkun:     r.NoAkun,
			NamaAkun:   r.NamaAkun,
			Kategori:   r.Kategori,
			Rincian:    r.Rincian,
			Keterangan: r.Keterangan,
			QtySatuan:  r.QtySatuan,
			QtyPaket:   r.QtyPaket,
			Sbu:        r.Sbu,
			Total:      tot,
		}
		rabMap[r.IdAktivitasv2] = append(rabMap[r.IdAktivitasv2], item)
		rabTotalMap[r.IdAktivitasv2] += tot
	}

	// Map Aktivitas to Proker
	aktMap := make(map[uint][]DashboardMonitoringAktivitas)
	prokerAnggaranMap := make(map[uint]float64)
	prokerPencairanMap := make(map[uint]float64)

	for _, a := range rawAktivitas {
		rabItems := rabMap[a.Id]
		if rabItems == nil {
			rabItems = []DashboardMonitoringRabItem{}
		}
		totRab := rabTotalMap[a.Id]

		// Status pencairan calculation:
		// Jika ada realisasi dari tim anggaran (hasil revisi), bandingkan cair dengan realisasi.
		// Jika belum ada data realisasi karena tidak ada revisi, bandingkan cair dengan pengajuan (totRab).
		targetBanding := totRab
		if a.Realisasi != nil && *a.Realisasi >= 0 {
			targetBanding = *a.Realisasi
		}

		statusPencairan := "belum_dicairkan"
		if a.Pencairan != nil && *a.Pencairan > 0 {
			pencairanVal := *a.Pencairan
			if pencairanVal == targetBanding {
				statusPencairan = "cair"
			} else if pencairanVal < targetBanding {
				statusPencairan = "tidak_penuh"
			} else {
				statusPencairan = "cair_kelebihan"
			}
			prokerPencairanMap[a.IdMatrikProker] += pencairanVal
		}

		prokerAnggaranMap[a.IdMatrikProker] += totRab

		aktItem := DashboardMonitoringAktivitas{
			ID:              a.Id,
			UUID:            a.UUID,
			Aktivitas:       a.Aktivitas,
			PIC:             a.PIC,
			Capaian:         a.Capaian,
			Realisasi:       a.Realisasi,
			Pencairan:       a.Pencairan,
			Status:          a.Status,
			StatusAnggaran:  a.StatusAnggaran,
			StatusUpml:      a.StatusUpml,
			StatusUpmf:      a.StatusUpmf,
			StatusUpms:      a.StatusUpms,
			StatusLpm:       a.StatusLpm,
			IsBatal:         a.IsBatal == 1,
			CatatanBatal:    a.CatatanBatal,
			TotalRab:        totRab,
			StatusPencairan: statusPencairan,
			Rab:             rabItems,
		}
		aktMap[a.IdMatrikProker] = append(aktMap[a.IdMatrikProker], aktItem)
	}

	var prokerList []DashboardMonitoringProker
	var globalTotalProker int
	var globalTotalAktivitas int
	var globalTotalAnggaran float64
	var globalTotalPencairan float64

	for _, p := range prokerRows {
		activities := aktMap[p.Id]
		if activities == nil {
			activities = []DashboardMonitoringAktivitas{}
		}
		totAnggaran := prokerAnggaranMap[p.Id]
		totPencairan := prokerPencairanMap[p.Id]

		unitLabel := p.NamaFakultasUnit
		if p.TypeFakultasUnit == "prodi" && p.Jenjang != "" {
			unitLabel = fmt.Sprintf("%s (%s)", p.NamaFakultasUnit, strings.ToUpper(p.Jenjang))
		}

		prokerItem := DashboardMonitoringProker{
			ID:                p.Id,
			UUID:              p.UUID,
			IdFakultasUnit:    p.IdFakultasUnit,
			NamaFakultasUnit:  unitLabel,
			TypeFakultasUnit:  p.TypeFakultasUnit,
			Jenjang:           p.Jenjang,
			NamaProgram:       p.NamaProgram,
			Tahun:             p.Tahun,
			IndikatorKinerja:  p.IndikatorKinerja,
			TargetUnpak:       p.TargetUnpak,
			TargetUnit:        p.TargetUnit,
			Output:            p.Output,
			Outcome:           p.Outcome,
			PIC:               p.PIC,
			Status:            p.Status,
			StatusVerifikator: p.StatusVerifikator,
			StatusAnggaran:    p.StatusAnggaran,
			BebanYayasan:      p.BebanYayasan == 1,
			IsBatal:           p.IsBatal == 1,
			CatatanBatal:      p.CatatanBatal,
			TotalAktivitas:    len(activities),
			TotalAnggaran:     totAnggaran,
			TotalPencairan:    totPencairan,
			Aktivitas:         activities,
		}

		prokerList = append(prokerList, prokerItem)
		globalTotalProker++
		globalTotalAktivitas += len(activities)
		globalTotalAnggaran += totAnggaran
		globalTotalPencairan += totPencairan
	}

	if prokerList == nil {
		prokerList = []DashboardMonitoringProker{}
	}
	if unitOptions == nil {
		unitOptions = []DashboardMonitoringUnitOption{}
	}

	return c.JSON(fiber.Map{
		"allowed":      true,
		"role_type":    userRoleType,
		"user_name":    userName,
		"unit_options": unitOptions,
		"data":         prokerList,
		"summary": fiber.Map{
			"total_proker":    globalTotalProker,
			"total_aktivitas": globalTotalAktivitas,
			"total_anggaran":  globalTotalAnggaran,
			"total_pencairan": globalTotalPencairan,
		},
	})
}

