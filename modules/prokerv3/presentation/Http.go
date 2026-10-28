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
	cmd.MatrixUUID = rawUuid

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
		Realisasi        *float64 `json:"realisasi"`
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
		if reqBody.Anggaran != nil && *reqBody.Anggaran > 0 {
			cmd.Anggaran = *reqBody.Anggaran
		}
		if reqBody.Realisasi != nil {
			cmd.Realisasi = reqBody.Realisasi
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
	if val := c.FormValue("realisasi"); val != "" {
		if r, err := strconv.ParseFloat(val, 64); err == nil {
			cmd.Realisasi = &r
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

	roles := []string{"admin", "auditee", "fakultas", "unit", "user", "upml", "upmf", "upms", "anggaran", "keuangan", "evaluator"}
	whoamiURL := "http://localhost:3000/whoami"

	// 1. Options & Target Unpak
	app.Get("/proker_v3_options", commonpresentation.SmartCompress(), commonpresentation.JWTMiddleware(), GetProkerV3OptionsHandler)
	app.Get("/proker_v3/target_unpak", commonpresentation.SmartCompress(), commonpresentation.JWTMiddleware(), GetTargetUnpakHandler)

	// Jadwal Proker V3 CRUD & Status Check (Didaftarkan sebelum /proker_v3/:uuid agar tidak tertabrak route param :uuid)
	adminRoles := []string{"admin"}
	app.Get("/proker_v3/jadwal", commonpresentation.SmartCompress(), commonpresentation.JWTMiddleware(), GetAllJadwalProkerV3Handler)
	app.Get("/proker_v3/jadwal_status", commonpresentation.SmartCompress(), commonpresentation.JWTMiddleware(), CheckJadwalProkerV3Handler)
	app.Get("/proker_v3/jadwal/:uuid", commonpresentation.SmartCompress(), commonpresentation.JWTMiddleware(), GetJadwalProkerV3Handler)
	app.Post("/proker_v3/jadwal", commonpresentation.JWTMiddleware(), commonpresentation.RBACMiddleware(adminRoles, whoamiURL), CreateJadwalProkerV3Handler)
	app.Put("/proker_v3/jadwal/:uuid", commonpresentation.JWTMiddleware(), commonpresentation.RBACMiddleware(adminRoles, whoamiURL), UpdateJadwalProkerV3Handler)
	app.Delete("/proker_v3/jadwal/:uuid", commonpresentation.JWTMiddleware(), commonpresentation.RBACMiddleware(adminRoles, whoamiURL), DeleteJadwalProkerV3Handler)

	// 2. Program Kerja (matrix_proker) CRUD
	app.Post("/proker_v3", commonpresentation.JWTMiddleware(), commonpresentation.RBACMiddleware(roles, whoamiURL), CreateProkerV3Handlerfunc)
	app.Get("/proker_v3/:uuid", commonpresentation.SmartCompress(), commonpresentation.JWTMiddleware(), commonpresentation.RBACMiddleware(roles, whoamiURL), GetProkerV3Handlerfunc)
	app.Put("/proker_v3/:uuid", commonpresentation.JWTMiddleware(), commonpresentation.RBACMiddleware(roles, whoamiURL), UpdateProkerV3Handlerfunc)
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

func GetAllJadwalProkerV3Handler(c *fiber.Ctx) error {
	if dbInstance == nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Database tidak tersedia"})
	}

	search := strings.TrimSpace(c.Query("search"))
	tahun := strings.TrimSpace(c.Query("tahun"))

	var list []domain.JadwalProkerV3
	query := dbInstance.WithContext(c.Context())

	if tahun != "" && tahun != "all" {
		query = query.Where("tahun = ?", tahun)
	}
	if search != "" {
		query = query.Where("tahun LIKE ?", "%"+search+"%")
	}

	if err := query.Order("tahun DESC, tanggal_mulai ASC").Find(&list).Error; err != nil {
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

		res = append(res, domain.JadwalProkerV3DTO{
			ID:           item.ID,
			UUID:         item.UUID.String(),
			Tahun:        item.Tahun,
			TanggalMulai: item.TanggalMulai.Format("2006-01-02"),
			TanggalAkhir: item.TanggalAkhir.Format("2006-01-02"),
			Status:       status,
			StatusLabel:  statusLabel,
			CreatedAt:    item.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt:    item.UpdatedAt.Format("2006-01-02 15:04:05"),
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

	var item domain.JadwalProkerV3
	if err := dbInstance.WithContext(c.Context()).Where("uuid = ?", parsedUuid).First(&item).Error; err != nil {
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

	return c.JSON(domain.JadwalProkerV3DTO{
		ID:           item.ID,
		UUID:         item.UUID.String(),
		Tahun:        item.Tahun,
		TanggalMulai: item.TanggalMulai.Format("2006-01-02"),
		TanggalAkhir: item.TanggalAkhir.Format("2006-01-02"),
		Status:       status,
		StatusLabel:  statusLabel,
		CreatedAt:    item.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:    item.UpdatedAt.Format("2006-01-02 15:04:05"),
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

	var item domain.JadwalProkerV3
	err := dbInstance.WithContext(c.Context()).Where("tahun = ?", tahun).Order("id DESC").First(&item).Error
	if err != nil {
		return c.JSON(fiber.Map{
			"has_schedule": false,
			"is_open":      true,
			"status":       "tidak_ada_jadwal",
			"status_label": "Tidak Ada Pembatasan Jadwal",
			"tahun":        tahun,
			"message":      fmt.Sprintf("Belum ada jadwal khusus pengajuan program kerja untuk tahun %s", tahun),
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
	message := fmt.Sprintf("Pengajuan program kerja tahun %s dibuka dari %s s/d %s", item.Tahun, formatTanggalIndo(item.TanggalMulai), formatTanggalIndo(item.TanggalAkhir))

	if now.Before(startOfDay) {
		isOpen = false
		status = "belum_dibuka"
		statusLabel = "Belum Dibuka"
		message = fmt.Sprintf("Pengajuan program kerja tahun %s belum dibuka. Jadwal pengajuan: %s s/d %s", item.Tahun, formatTanggalIndo(item.TanggalMulai), formatTanggalIndo(item.TanggalAkhir))
	} else if now.After(endOfDay) {
		isOpen = false
		status = "ditutup"
		statusLabel = "Ditutup"
		message = fmt.Sprintf("Pengajuan program kerja tahun %s telah ditutup pada %s", item.Tahun, formatTanggalIndo(item.TanggalAkhir))
	}

	return c.JSON(fiber.Map{
		"has_schedule":  true,
		"is_open":       isOpen,
		"status":        status,
		"status_label":  statusLabel,
		"tahun":         item.Tahun,
		"tanggal_mulai": item.TanggalMulai.Format("2006-01-02"),
		"tanggal_akhir": item.TanggalAkhir.Format("2006-01-02"),
		"message":       message,
	})
}

func CreateJadwalProkerV3Handler(c *fiber.Ctx) error {
	if dbInstance == nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Database tidak tersedia"})
	}

	type FormInput struct {
		Tahun        string `json:"tahun" form:"tahun"`
		TanggalMulai string `json:"tanggal_mulai" form:"tanggal_mulai"`
		TanggalAkhir string `json:"tanggal_akhir" form:"tanggal_akhir"`
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

	newJadwal := domain.JadwalProkerV3{
		UUID:         uuid.New(),
		Tahun:        req.Tahun,
		TanggalMulai: tMulai,
		TanggalAkhir: tAkhir,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
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
		Tahun        string `json:"tahun" form:"tahun"`
		TanggalMulai string `json:"tanggal_mulai" form:"tanggal_mulai"`
		TanggalAkhir string `json:"tanggal_akhir" form:"tanggal_akhir"`
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
