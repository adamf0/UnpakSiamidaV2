package presentation

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/mehdihadeli/go-mediatr"
	"gorm.io/gorm"

	commondomain "UnpakSiamida/common/domain"
	commoninfra "UnpakSiamida/common/infrastructure"
	commonpresentation "UnpakSiamida/common/presentation"
	CreateProker "UnpakSiamida/modules/prokerv3/application/CreateProkerV3"
	DeleteProker "UnpakSiamida/modules/prokerv3/application/DeleteProkerV3"
	GetAllProker "UnpakSiamida/modules/prokerv3/application/GetAllProkerV3"
	GetProker "UnpakSiamida/modules/prokerv3/application/GetProkerV3"
	UpdateProker "UnpakSiamida/modules/prokerv3/application/UpdateProkerV3"
	domain "UnpakSiamida/modules/prokerv3/domain"
	infra "UnpakSiamida/modules/prokerv3/infrastructure"
)

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

func matchUnit(itemUnitID *int, itemUnitUUID string, userFU string) bool {
	if userFU == "" {
		return false
	}
	if itemUnitUUID != "" && strings.EqualFold(strings.TrimSpace(itemUnitUUID), strings.TrimSpace(userFU)) {
		return true
	}
	if itemUnitID != nil && strconv.Itoa(*itemUnitID) == strings.TrimSpace(userFU) {
		return true
	}
	return false
}

// =======================================================
// POST /proker_v3 (Simpan ke matrix_proker)
// =======================================================
func CreateProkerV3Handlerfunc(c *fiber.Ctx) error {
	level, userFU := getUserInfo(c)

	var cmd CreateProker.CreateProkerV3Command
	if err := c.BodyParser(&cmd); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid JSON format: " + err.Error()})
	}

	// Scoping non-admin
	if level != "admin" {
		if cmd.FakultasUuid == "" && cmd.IdFakultasUnit == nil {
			cmd.FakultasUuid = userFU
		}
	}

	uuidResult, err := mediatr.Send[CreateProker.CreateProkerV3Command, string](context.Background(), cmd)
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

	var cmd UpdateProker.UpdateProkerV3Command
	if err := c.BodyParser(&cmd); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid JSON format: " + err.Error()})
	}
	cmd.UUID = rawUuid

	level, userFU := getUserInfo(c)
	if level != "admin" {
		parsedUuid, err := uuid.Parse(rawUuid)
		if err == nil {
			query := GetProker.GetProkerV3ByUuidQuery{UUID: parsedUuid}
			existing, err := mediatr.Send[GetProker.GetProkerV3ByUuidQuery, *domain.ProkerV3Default](context.Background(), query)
			if err == nil && existing != nil {
				if !matchUnit(existing.IdFakultasUnit, existing.UuidFakultasUnit, userFU) {
					return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Akses ditolak: bukan unit anda"})
				}
			}
		}
	}

	uuidResult, err := mediatr.Send[UpdateProker.UpdateProkerV3Command, string](context.Background(), cmd)
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

	level, userFU := getUserInfo(c)
	if level != "admin" {
		parsedUuid, err := uuid.Parse(rawUuid)
		if err == nil {
			query := GetProker.GetProkerV3ByUuidQuery{UUID: parsedUuid}
			existing, err := mediatr.Send[GetProker.GetProkerV3ByUuidQuery, *domain.ProkerV3Default](context.Background(), query)
			if err == nil && existing != nil {
				if !matchUnit(existing.IdFakultasUnit, existing.UuidFakultasUnit, userFU) {
					return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Akses ditolak: bukan unit anda"})
				}
			}
		}
	}

	cmd := DeleteProker.DeleteProkerV3Command{UUID: rawUuid}
	deletedUUID, err := mediatr.Send[DeleteProker.DeleteProkerV3Command, string](context.Background(), cmd)
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

	level, userFU := getUserInfo(c)

	query := GetProker.GetProkerV3ByUuidQuery{UUID: parsedUuid}
	result, err := mediatr.Send[GetProker.GetProkerV3ByUuidQuery, *domain.ProkerV3Default](context.Background(), query)
	if err != nil {
		return commoninfra.HandleError(c, err)
	}
	if result == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Program kerja tidak ditemukan"})
	}

	if level != "admin" && !matchUnit(result.IdFakultasUnit, result.UuidFakultasUnit, userFU) {
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

	level, userFU := getUserInfo(c)

	var filters []commondomain.SearchFilter

	var targetFU string
	if level != "admin" {
		if userFU != "" && userFU != "0" && userFU != "all" {
			targetFU = userFU
		} else if qFU := strings.TrimSpace(c.Query("filter_fakultas_unit")); qFU != "" && qFU != "all" && qFU != "semua" {
			targetFU = qFU
		}
	} else {
		targetFU = strings.TrimSpace(c.Params("target"))
		if targetFU == "" || targetFU == "all" || targetFU == "semua" {
			targetFU = strings.TrimSpace(c.Query("filter_fakultas_unit"))
		}
		if targetFU == "" || targetFU == "all" || targetFU == "semua" {
			targetFU = strings.TrimSpace(c.Query("fakultas_unit"))
		}
		if targetFU == "" || targetFU == "all" || targetFU == "semua" {
			targetFU = strings.TrimSpace(c.Query("id_fakultas_unit"))
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

	query := GetAllProker.GetAllProkerV3Query{
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

	result, err := mediatr.Send[GetAllProker.GetAllProkerV3Query, commondomain.Paged[domain.ProkerV3Default]](context.Background(), query)
	if err != nil {
		return commoninfra.HandleError(c, err)
	}

	return adapter.Send(c, result)
}

// =======================================================
// GET /proker_v3/:uuid/aktivitas (Daftar aktivitasv2)
// =======================================================
func GetAktivitasByMatrixHandler(db *gorm.DB) fiber.Handler {
	repo := infra.NewProkerV3Repository(db)

	return func(c *fiber.Ctx) error {
		rawUuid := c.Params("uuid")
		parsedUuid, err := uuid.Parse(rawUuid)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid matrix UUID"})
		}

		// Get proker detail
		proker, err := repo.GetDefaultByUuid(c.Context(), parsedUuid)
		if err != nil {
			return commoninfra.HandleError(c, err)
		}
		if proker == nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Program kerja tidak ditemukan"})
		}

		level, userFU := getUserInfo(c)
		if level != "admin" && !matchUnit(proker.IdFakultasUnit, proker.UuidFakultasUnit, userFU) {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Akses ditolak: bukan unit anda"})
		}

		// Get activities for this matrix_proker
		activities, err := repo.GetAktivitasByMatrikId(c.Context(), proker.ID)
		if err != nil {
			return commoninfra.HandleError(c, err)
		}

		return c.JSON(fiber.Map{
			"proker": proker,
			"data":   activities,
			"total":  len(activities),
		})
	}
}

// =======================================================
// POST /proker_v3/:uuid/aktivitas (Tambah aktivitasv2)
// =======================================================
func CreateAktivitasHandler(db *gorm.DB) fiber.Handler {
	repo := infra.NewProkerV3Repository(db)

	return func(c *fiber.Ctx) error {
		rawUuid := c.Params("uuid")
		parsedUuid, err := uuid.Parse(rawUuid)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid matrix UUID"})
		}

		proker, err := repo.GetDefaultByUuid(c.Context(), parsedUuid)
		if err != nil {
			return commoninfra.HandleError(c, err)
		}
		if proker == nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Program kerja tidak ditemukan"})
		}

		level, userFU := getUserInfo(c)
		if level != "admin" && !matchUnit(proker.IdFakultasUnit, proker.UuidFakultasUnit, userFU) {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Akses ditolak: bukan unit anda"})
		}

		var input struct {
			Aktivitas        string                           `json:"aktivitas"`
			PIC              string                           `json:"pic"`
			Anggaran         float64                          `json:"anggaran"`
			LuaranProduk     string                           `json:"luaran_produk"`
			DokumenReferensi string                           `json:"dokumen_referensi"`
			Status           string                           `json:"status"`
			Capaian          string                           `json:"capaian"`
			Masalah          string                           `json:"masalah"`
			PemecahanMasalah string                           `json:"pemecahan_masalah"`
			Timelines        []CreateProker.TimelineItemInput `json:"timelines"`
		}

		if err := c.BodyParser(&input); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid JSON format: " + err.Error()})
		}

		namaAktivitas := strings.TrimSpace(input.Aktivitas)
		if namaAktivitas == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Nama aktivitas / kegiatan wajib diisi"})
		}

		statusStr := strings.ToLower(strings.TrimSpace(input.Status))
		if statusStr == "" {
			statusStr = "draft"
		}

		var fuIDUint64 *uint64
		if proker.IdFakultasUnit != nil {
			val := uint64(*proker.IdFakultasUnit)
			fuIDUint64 = &val
		}
		matrikIDUint64 := uint64(proker.ID)

		act := &domain.AktivitasV2{
			UUID:             uuid.New(),
			IdMatrikProker:   &matrikIDUint64,
			IdFakultasUnit:   fuIDUint64,
			Aktivitas:        namaAktivitas,
			PIC:              strings.TrimSpace(input.PIC),
			Capaian:          strings.TrimSpace(input.Capaian),
			Masalah:          strings.TrimSpace(input.Masalah),
			PemecahanMasalah: strings.TrimSpace(input.PemecahanMasalah),
			Anggaran:         input.Anggaran,
			Status:           statusStr,
		}

		timelines := make([]domain.ProkerTimelineV3, 0, len(input.Timelines))
		for i, tInput := range input.Timelines {
			tlUUID := uuid.New()
			if strings.TrimSpace(tInput.UUID) != "" {
				if parsed, err := uuid.Parse(tInput.UUID); err == nil {
					tlUUID = parsed
				}
			}

			noAktivitas := tInput.NoAktivitas
			if noAktivitas <= 0 {
				noAktivitas = i + 1
			}

			rkAwal := parseDateFlexible(tInput.TanggalRkAwal)
			if rkAwal == nil {
				def, _ := time.Parse("2006-01-02", "2026-01-01")
				rkAwal = &def
			}

			rkAkhir := parseDateFlexible(tInput.TanggalRkAkhir)
			if rkAkhir == nil {
				def, _ := time.Parse("2006-01-02", "2026-12-31")
				rkAkhir = &def
			}

			rAwal := parseDateFlexible(tInput.TanggalRAwal)
			rAkhir := parseDateFlexible(tInput.TanggalRAkhir)

			titikKontrol := strings.TrimSpace(tInput.TitikKontrol)
			if titikKontrol == "" {
				titikKontrol = "RK"
			}

			tlStatus := strings.TrimSpace(tInput.Status)
			if tlStatus == "" {
				tlStatus = "Belum Terlaksana"
			}

			tNama := strings.TrimSpace(tInput.NamaAktivitas)
			if tNama == "" {
				tNama = namaAktivitas
			}

			timelines = append(timelines, domain.ProkerTimelineV3{
				UUID:             tlUUID,
				NoAktivitas:      noAktivitas,
				NamaAktivitas:    tNama,
				PIC:              strings.TrimSpace(tInput.PIC),
				TitikKontrol:     titikKontrol,
				TanggalRkAwal:    *rkAwal,
				TanggalRkAkhir:   *rkAkhir,
				TanggalRAwal:     rAwal,
				TanggalRAkhir:    rAkhir,
				JadwalRk:         strings.TrimSpace(tInput.JadwalRk),
				JadwalR:          strings.TrimSpace(tInput.JadwalR),
				Capaian:          strings.TrimSpace(tInput.Capaian),
				Masalah:          strings.TrimSpace(tInput.Masalah),
				PemecahanMasalah: strings.TrimSpace(tInput.PemecahanMasalah),
				Status:           tlStatus,
			})
		}

		if err := repo.CreateAktivitas(c.Context(), act, timelines); err != nil {
			return commoninfra.HandleError(c, err)
		}

		return c.Status(fiber.StatusCreated).JSON(fiber.Map{
			"message": "Aktivitas berhasil ditambahkan",
			"uuid":    act.UUID.String(),
			"data":    act,
		})
	}
}

// =======================================================
// PUT /proker_v3/aktivitas/:aktivitas_uuid (Update aktivitasv2)
// =======================================================
func UpdateAktivitasHandler(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		rawUuid := c.Params("aktivitas_uuid")
		actUUID, err := uuid.Parse(rawUuid)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid aktivitas UUID"})
		}

		var act domain.AktivitasV2
		if err := db.Where("uuid = ?", actUUID).First(&act).Error; err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Aktivitas tidak ditemukan"})
		}

		var input struct {
			Aktivitas        string   `json:"aktivitas"`
			PIC              string   `json:"pic"`
			Anggaran         float64  `json:"anggaran"`
			Realisasi        *float64 `json:"realisasi"`
			LuaranProduk     string   `json:"luaran_produk"`
			DokumenReferensi string   `json:"dokumen_referensi"`
			Status           string   `json:"status"`
			Capaian          string   `json:"capaian"`
			Masalah          string   `json:"masalah"`
			PemecahanMasalah string   `json:"pemecahan_masalah"`
		}

		if err := c.BodyParser(&input); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid JSON format: " + err.Error()})
		}

		if input.Aktivitas != "" {
			act.Aktivitas = strings.TrimSpace(input.Aktivitas)
		}
		if input.PIC != "" {
			act.PIC = strings.TrimSpace(input.PIC)
		}
		act.Anggaran = input.Anggaran
		if input.Realisasi != nil {
			act.Realisasi = input.Realisasi
		}
		if input.Status != "" {
			act.Status = strings.TrimSpace(input.Status)
		}
		if input.Capaian != "" {
			act.Capaian = strings.TrimSpace(input.Capaian)
		}
		if input.Masalah != "" {
			act.Masalah = strings.TrimSpace(input.Masalah)
		}
		if input.PemecahanMasalah != "" {
			act.PemecahanMasalah = strings.TrimSpace(input.PemecahanMasalah)
		}
		act.UpdatedAt = time.Now()

		if err := db.Save(&act).Error; err != nil {
			return commoninfra.HandleError(c, err)
		}

		return c.JSON(fiber.Map{
			"message": "Aktivitas berhasil diperbarui",
			"uuid":    rawUuid,
			"data":    act,
		})
	}
}

// =======================================================
// DELETE /proker_v3/aktivitas/:aktivitas_uuid (Hapus aktivitasv2)
// =======================================================
func DeleteAktivitasHandler(db *gorm.DB) fiber.Handler {
	repo := infra.NewProkerV3Repository(db)

	return func(c *fiber.Ctx) error {
		rawUuid := c.Params("aktivitas_uuid")
		actUUID, err := uuid.Parse(rawUuid)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid aktivitas UUID"})
		}

		if err := repo.DeleteAktivitas(c.Context(), actUUID); err != nil {
			return commoninfra.HandleError(c, err)
		}

		return c.JSON(fiber.Map{
			"message": "Aktivitas berhasil dihapus",
			"uuid":    rawUuid,
		})
	}
}

// =======================================================
// PUT /proker_v3/aktivitas/:aktivitas_uuid/status (Ubah status aktivitasv2)
// =======================================================
func UpdateAktivitasStatusHandler(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		rawUuid := c.Params("aktivitas_uuid")
		if rawUuid == "" {
			rawUuid = c.Params("uuid")
		}
		if rawUuid == "" {
			rawUuid = c.Query("uuid")
		}
		parsedUuid, err := uuid.Parse(rawUuid)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid UUID"})
		}

		var body struct {
			Status string `json:"status"`
		}
		if err := c.BodyParser(&body); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid JSON format"})
		}
		newStatus := strings.ToLower(strings.TrimSpace(body.Status))
		if newStatus == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Status tidak boleh kosong"})
		}

		level, userFU := getUserInfo(c)

		var act domain.AktivitasV2
		if err := db.Where("uuid = ?", parsedUuid).First(&act).Error; err != nil {
			// Check if it is a matrix_proker
			var mp domain.ProkerV3
			if errMp := db.Where("uuid = ?", parsedUuid).First(&mp).Error; errMp == nil {
				if level != "admin" {
					if !matchUnit(mp.IdFakultasUnit, "", userFU) {
						return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Akses ditolak: bukan data unit anda"})
					}
					if newStatus != "draft" && newStatus != "diajukan" {
						return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Sebagai auditee, status hanya boleh Draf atau Diajukan"})
					}
				}
				_ = db.Model(&domain.AktivitasV2{}).Where("id_matrik_proker = ?", mp.ID).Update("status", newStatus)
				return c.JSON(fiber.Map{
					"message": "Status program kerja berhasil diubah",
					"status":  newStatus,
					"uuid":    rawUuid,
				})
			}
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Data aktivitas/program kerja tidak ditemukan"})
		}

		if level != "admin" {
			var fuID *int
			if act.IdFakultasUnit != nil {
				val := int(*act.IdFakultasUnit)
				fuID = &val
			}
			if !matchUnit(fuID, "", userFU) {
				return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Akses ditolak: bukan data unit anda"})
			}
			if newStatus != "draft" && newStatus != "diajukan" {
				return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Sebagai auditee, status hanya boleh Draf atau Diajukan"})
			}
		}

		if err := db.Model(&act).Update("status", newStatus).Error; err != nil {
			return commoninfra.HandleError(c, err)
		}

		return c.JSON(fiber.Map{
			"message": "Status aktivitas berhasil diubah",
			"status":  newStatus,
			"uuid":    rawUuid,
		})
	}
}

// =======================================================
// POST /proker_v3/timeline/:timeline_uuid/files (MULTIPLE UPLOAD)
// =======================================================
func UploadTimelineFilesHandler(db *gorm.DB) fiber.Handler {
	repo := infra.NewProkerV3Repository(db)

	return func(c *fiber.Ctx) error {
		timelineUUIDStr := c.Params("timeline_uuid")
		timelineUUID, err := uuid.Parse(timelineUUIDStr)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid timeline UUID"})
		}

		timeline, err := repo.GetTimelineByUuid(c.Context(), timelineUUID)
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
		uploadDir := resolveTimelineUploadDir(db, repo, c.Context(), timeline)
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
			if err := repo.AddTimelineFiles(c.Context(), filesToInsert); err != nil {
				return commoninfra.HandleError(c, err)
			}
		}

		return c.Status(fiber.StatusCreated).JSON(fiber.Map{
			"message": "File berhasil diunggah",
			"count":   len(filesToInsert),
			"files":   filesToInsert,
		})
	}
}

// =======================================================
// DELETE /proker_v3/file/:file_uuid
// =======================================================
func DeleteTimelineFileHandler(db *gorm.DB) fiber.Handler {
	repo := infra.NewProkerV3Repository(db)

	return func(c *fiber.Ctx) error {
		fileUUIDStr := c.Params("file_uuid")
		fileUUID, err := uuid.Parse(fileUUIDStr)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid file UUID"})
		}

		fileObj, err := repo.GetTimelineFileByUuid(c.Context(), fileUUID)
		if err != nil {
			return commoninfra.HandleError(c, err)
		}
		if fileObj == nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "File tidak ditemukan"})
		}

		if err := repo.DeleteTimelineFile(c.Context(), fileUUID); err != nil {
			return commoninfra.HandleError(c, err)
		}

		if fileObj.File != "" {
			_ = os.Remove(fileObj.File)
			if fileObj.FileURL != "" {
				_ = os.Remove("." + fileObj.FileURL)
			}
			_ = os.Remove(filepath.Join("./upload/proker_v3", fileObj.File))
		}

		return c.JSON(fiber.Map{
			"message": "File berhasil dihapus",
			"uuid":    fileUUIDStr,
		})
	}
}

func parseDateFlexible(val string) *time.Time {
	val = strings.TrimSpace(val)
	if val == "" {
		return nil
	}
	if t, err := time.Parse("2006-01-02", val); err == nil {
		return &t
	}
	if t, err := time.Parse(time.RFC3339, val); err == nil {
		return &t
	}
	return nil
}

// =======================================================
// GET /proker_v3_options
// =======================================================
func GetProkerV3OptionsHandler(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var units []struct {
			ID   int    `json:"id"`
			UUID string `json:"uuid"`
			Nama string `json:"nama"`
		}
		_ = db.Table("v_fakultas_unit").
			Select("id, uuid, nama_fak_prod_unit as nama").
			Order("nama_fak_prod_unit ASC").
			Scan(&units).Error

		var tahuns []struct {
			ID    int    `json:"id"`
			UUID  string `json:"uuid"`
			Tahun string `json:"tahun"`
		}
		_ = db.Table("master_tahun").
			Select("id, uuid, tahun").
			Order("tahun DESC").
			Scan(&tahuns).Error

		var standars []struct {
			ID   int    `json:"id"`
			UUID string `json:"uuid"`
			Nama string `json:"nama"`
		}
		_ = db.Table("master_standar_renstra").
			Select("id, uuid, nama").
			Order("nama ASC").
			Scan(&standars).Error

		var indikators []struct {
			ID              int    `json:"id"`
			UUID            string `json:"uuid"`
			IdMasterStandar int    `json:"id_master_standar"`
			Indikator       string `json:"indikator"`
			Tahun           int    `json:"tahun"`
		}
		_ = db.Table("master_indikator_renstra").
			Select("id, uuid, id_master_standar, indikator, tahun").
			Order("indikator ASC").
			Scan(&indikators).Error

		return c.JSON(fiber.Map{
			"units":      units,
			"tahuns":     tahuns,
			"standars":   standars,
			"indikators": indikators,
		})
	}
}

// =======================================================
// GET /proker_v3/target_unpak
// =======================================================
func GetTargetUnpakHandler(db *gorm.DB) fiber.Handler {
	type TemplateRow struct {
		ID        int    `gorm:"column:id"`
		Target    string `gorm:"column:target"`
		TargetMin string `gorm:"column:target_min"`
		TargetMax string `gorm:"column:target_max"`
		Satuan    string `gorm:"column:satuan"`
	}

	formatTargetRow := func(t TemplateRow) string {
		var targetVal string
		target := strings.TrimSpace(t.Target)
		min := strings.TrimSpace(t.TargetMin)
		max := strings.TrimSpace(t.TargetMax)
		satuan := strings.TrimSpace(t.Satuan)

		if target != "" {
			targetVal = target
		} else if min != "" && max != "" {
			targetVal = fmt.Sprintf("%s - %s", min, max)
		} else if min != "" {
			targetVal = fmt.Sprintf(">= %s", min)
		} else if max != "" {
			targetVal = fmt.Sprintf("<= %s", max)
		}

		if targetVal != "" && satuan != "" {
			return fmt.Sprintf("%s %s", targetVal, satuan)
		}
		return targetVal
	}

	return func(c *fiber.Ctx) error {
		tahun := strings.TrimSpace(c.Query("tahun"))
		indikatorParam := strings.TrimSpace(c.Query("id_indikator"))
		if indikatorParam == "" {
			indikatorParam = strings.TrimSpace(c.Query("indikator"))
		}
		unitParam := strings.TrimSpace(c.Query("id_fakultas_unit"))
		if unitParam == "" {
			unitParam = strings.TrimSpace(c.Query("fakultas_unit"))
		}

		if indikatorParam == "" {
			return c.JSON(fiber.Map{
				"found":        false,
				"target":       "-",
				"target_min":   "",
				"target_max":   "",
				"satuan":       "",
				"target_unpak": "-",
			})
		}

		indikatorID, err := strconv.Atoi(indikatorParam)
		if err != nil {
			var indRow struct{ ID int }
			if err := db.Table("master_indikator_renstra").Select("id").Where("uuid = ?", indikatorParam).Take(&indRow).Error; err == nil {
				indikatorID = indRow.ID
			}
		}

		var unitID int
		if unitParam != "" {
			if id, err := strconv.Atoi(unitParam); err == nil {
				unitID = id
			} else {
				var uRow struct{ ID int }
				if err := db.Table("v_fakultas_unit").Select("id").Where("uuid = ?", unitParam).Take(&uRow).Error; err == nil {
					unitID = uRow.ID
				}
			}
		}

		candidateIDs := []int{indikatorID}
		if indikatorID > 0 {
			var indRow struct {
				Indikator string `gorm:"column:indikator"`
			}
			if err := db.Table("master_indikator_renstra").Select("indikator").Where("id = ?", indikatorID).First(&indRow).Error; err == nil && indRow.Indikator != "" {
				var sameIDs []int
				_ = db.Table("master_indikator_renstra").Where("indikator = ?", indRow.Indikator).Pluck("id", &sameIDs).Error
				for _, sid := range sameIDs {
					if sid != indikatorID {
						candidateIDs = append(candidateIDs, sid)
					}
				}
			}
		}

		var tpl TemplateRow
		found := false

		if tahun != "" && unitID > 0 && len(candidateIDs) > 0 {
			if err := db.Table("template_renstra").
				Where("indikator IN ? AND fakultas_unit = ? AND tahun = ?", candidateIDs, unitID, tahun).
				Order("id DESC").First(&tpl).Error; err == nil {
				if formatTargetRow(tpl) != "" {
					found = true
				}
			}
		}

		if !found && unitID > 0 && len(candidateIDs) > 0 {
			if err := db.Table("template_renstra").
				Where("indikator IN ? AND fakultas_unit = ?", candidateIDs, unitID).
				Order("tahun DESC, id DESC").First(&tpl).Error; err == nil {
				if formatTargetRow(tpl) != "" {
					found = true
				}
			}
		}

		if !found && tahun != "" && len(candidateIDs) > 0 {
			if err := db.Table("template_renstra").
				Where("indikator IN ? AND tahun = ?", candidateIDs, tahun).
				Order("id DESC").First(&tpl).Error; err == nil {
				if formatTargetRow(tpl) != "" {
					found = true
				}
			}
		}

		if !found && len(candidateIDs) > 0 {
			if err := db.Table("template_renstra").
				Where("indikator IN ?", candidateIDs).
				Order("tahun DESC, id DESC").First(&tpl).Error; err == nil {
				if formatTargetRow(tpl) != "" {
					found = true
				}
			}
		}

		formatted := formatTargetRow(tpl)
		if formatted == "" {
			formatted = "-"
		}

		return c.JSON(fiber.Map{
			"found":        found,
			"target":       tpl.Target,
			"target_min":   tpl.TargetMin,
			"target_max":   tpl.TargetMax,
			"satuan":       tpl.Satuan,
			"target_unpak": formatted,
		})
	}
}

// =======================================================
// POST /proker_v3/:uuid/dokumen_referensi
// =======================================================
func UploadDokumenReferensiHandler(db *gorm.DB) fiber.Handler {
	repo := infra.NewProkerV3Repository(db)

	return func(c *fiber.Ctx) error {
		rawUuid := c.Params("uuid")
		parsedUuid, err := uuid.Parse(rawUuid)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid UUID"})
		}

		// Check if uuid is an aktivitasv2 or matrix_proker
		var act domain.AktivitasV2
		isAct := false
		if err := db.Where("uuid = ?", parsedUuid).First(&act).Error; err == nil {
			isAct = true
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
		var uploadedNames []string
		for _, fh := range fileHeaders {
			ext := filepath.Ext(fh.Filename)
			fileUUID := uuid.New()
			savedName := fmt.Sprintf("docref-%s%s", fileUUID.String(), ext)
			destPath := filepath.Join(uploadDir, savedName)

			if err := c.SaveFile(fh, destPath); err != nil {
				continue
			}

			var actIDPtr *uint
			if isAct {
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
			uploadedNames = append(uploadedNames, fh.Filename)
		}

		if len(filesToInsert) > 0 {
			if err := repo.AddTimelineFiles(c.Context(), filesToInsert); err != nil {
				return commoninfra.HandleError(c, err)
			}
		}

		return c.Status(fiber.StatusCreated).JSON(fiber.Map{
			"message": "Dokumen referensi berhasil diunggah",
			"count":   len(filesToInsert),
			"files":   filesToInsert,
		})
	}
}

// =======================================================
// POST /proker_v3/aktivitas/:aktivitas_uuid/timeline
// =======================================================
func CreateTimelineHandler(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		actUUIDStr := c.Params("aktivitas_uuid")
		if actUUIDStr == "" {
			actUUIDStr = c.Params("uuid")
		}
		actUUID, err := uuid.Parse(actUUIDStr)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid aktivitas UUID"})
		}

		var act domain.AktivitasV2
		if err := db.Where("uuid = ?", actUUID).First(&act).Error; err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Aktivitas tidak ditemukan"})
		}

		_ = db.Exec("ALTER TABLE aktivitasv2_timeline MODIFY COLUMN titik_kontrol TEXT NULL;").Error

		var input CreateProker.TimelineItemInput
		if err := c.BodyParser(&input); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid JSON format"})
		}

		rkAwal := parseDateFlexible(input.TanggalRkAwal)
		if rkAwal == nil {
			def, _ := time.Parse("2006-01-02", "2026-01-01")
			rkAwal = &def
		}

		rkAkhir := parseDateFlexible(input.TanggalRkAkhir)
		if rkAkhir == nil {
			def, _ := time.Parse("2006-01-02", "2026-12-31")
			rkAkhir = &def
		}

		noAktivitas := input.NoAktivitas
		if noAktivitas <= 0 {
			var count int64
			_ = db.Table("aktivitasv2_timeline").Where("id_aktivitasv2 = ?", act.ID).Count(&count)
			noAktivitas = int(count + 1)
		}

		titikKontrol := strings.TrimSpace(input.TitikKontrol)
		if titikKontrol == "" {
			titikKontrol = "RK"
		}

		status := strings.TrimSpace(input.Status)
		if status == "" {
			status = "Belum Terlaksana"
		}

		namaAkt := strings.TrimSpace(input.NamaAktivitas)
		if namaAkt == "" {
			namaAkt = act.Aktivitas
		}

		pic := strings.TrimSpace(input.PIC)
		if pic == "" {
			pic = act.PIC
		}

		newTL := domain.ProkerTimelineV3{
			UUID:             uuid.New(),
			IdAktivitasV2:    act.ID,
			NoAktivitas:      noAktivitas,
			NamaAktivitas:    namaAkt,
			PIC:              pic,
			TitikKontrol:     titikKontrol,
			TanggalRkAwal:    *rkAwal,
			TanggalRkAkhir:   *rkAkhir,
			TanggalRAwal:     parseDateFlexible(input.TanggalRAwal),
			TanggalRAkhir:    parseDateFlexible(input.TanggalRAkhir),
			JadwalRk:         strings.TrimSpace(input.JadwalRk),
			JadwalR:          strings.TrimSpace(input.JadwalR),
			Capaian:          strings.TrimSpace(input.Capaian),
			Masalah:          strings.TrimSpace(input.Masalah),
			PemecahanMasalah: strings.TrimSpace(input.PemecahanMasalah),
			Status:           status,
		}

		if err := db.Create(&newTL).Error; err != nil {
			return commoninfra.HandleError(c, err)
		}

		return c.Status(fiber.StatusCreated).JSON(fiber.Map{
			"message": "Timeline aktivitas berhasil ditambahkan",
			"uuid":    newTL.UUID.String(),
			"data":    newTL,
		})
	}
}

// =======================================================
// PUT /proker_v3/timeline/:timeline_uuid
// =======================================================
func UpdateTimelineHandler(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		tlUUIDStr := c.Params("timeline_uuid")
		tlUUID, err := uuid.Parse(tlUUIDStr)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid timeline UUID"})
		}

		var tl domain.ProkerTimelineV3
		if err := db.Where("uuid = ?", tlUUID).First(&tl).Error; err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Timeline aktivitas tidak ditemukan"})
		}

		var input CreateProker.TimelineItemInput
		if err := c.BodyParser(&input); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid JSON format"})
		}

		if input.NamaAktivitas != "" {
			tl.NamaAktivitas = strings.TrimSpace(input.NamaAktivitas)
		}
		if input.PIC != "" {
			tl.PIC = strings.TrimSpace(input.PIC)
		}
		if input.TitikKontrol != "" {
			tl.TitikKontrol = strings.TrimSpace(input.TitikKontrol)
		}
		if input.NoAktivitas > 0 {
			tl.NoAktivitas = input.NoAktivitas
		}
		if input.JadwalRk != "" {
			tl.JadwalRk = strings.TrimSpace(input.JadwalRk)
		}
		if input.JadwalR != "" {
			tl.JadwalR = strings.TrimSpace(input.JadwalR)
		}
		if input.Capaian != "" {
			tl.Capaian = strings.TrimSpace(input.Capaian)
		}
		if input.Masalah != "" {
			tl.Masalah = strings.TrimSpace(input.Masalah)
		}
		if input.PemecahanMasalah != "" {
			tl.PemecahanMasalah = strings.TrimSpace(input.PemecahanMasalah)
		}
		if input.Status != "" {
			tl.Status = strings.TrimSpace(input.Status)
		}
		if t := parseDateFlexible(input.TanggalRkAwal); t != nil {
			tl.TanggalRkAwal = *t
		}
		if t := parseDateFlexible(input.TanggalRkAkhir); t != nil {
			tl.TanggalRkAkhir = *t
		}
		if t := parseDateFlexible(input.TanggalRAwal); t != nil {
			tl.TanggalRAwal = t
		}
		if t := parseDateFlexible(input.TanggalRAkhir); t != nil {
			tl.TanggalRAkhir = t
		}

		if err := db.Save(&tl).Error; err != nil {
			return commoninfra.HandleError(c, err)
		}

		return c.JSON(fiber.Map{
			"message": "Timeline aktivitas berhasil diperbarui",
			"uuid":    tlUUIDStr,
			"data":    tl,
		})
	}
}

// =======================================================
// PUT /proker_v3/aktivitas/:aktivitas_uuid/timeline/reorder
// =======================================================
func ReorderTimelineHandler(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var items []struct {
			UUID        string `json:"uuid"`
			NoAktivitas int64  `json:"no_aktivitas"`
		}
		if err := c.BodyParser(&items); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Format JSON tidak valid: " + err.Error()})
		}

		for _, item := range items {
			if item.UUID != "" {
				_ = db.Table("aktivitasv2_timeline").
					Where("uuid = ?", item.UUID).
					Update("no_aktivitas", item.NoAktivitas).Error
			}
		}

		return c.JSON(fiber.Map{
			"message": "Urutan nomor urut tahapan timeline berhasil diperbarui",
			"count":   len(items),
		})
	}
}

// =======================================================
// DELETE /proker_v3/timeline/:timeline_uuid
// =======================================================
func DeleteTimelineHandler(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		tlUUIDStr := c.Params("timeline_uuid")
		tlUUID, err := uuid.Parse(tlUUIDStr)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid timeline UUID"})
		}

		var tl domain.ProkerTimelineV3
		if err := db.Where("uuid = ?", tlUUID).First(&tl).Error; err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Timeline aktivitas tidak ditemukan"})
		}

		if err := db.Delete(&tl).Error; err != nil {
			return commoninfra.HandleError(c, err)
		}

		return c.JSON(fiber.Map{
			"message": "Timeline aktivitas berhasil dihapus",
			"uuid":    tlUUIDStr,
		})
	}
}

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

// resolveTimelineUploadDir resolves destination folder:
// ./upload/prokerv2/[id_fakultas_unit]/[tahun]/[standar]/[indikator]/[aktivitasv2.nama_aktivitas]/[aktivitasv2_timeline.nama_aktivitas]
func resolveTimelineUploadDir(db *gorm.DB, repo domain.IProkerV3Repository, ctx context.Context, timeline *domain.ProkerTimelineV3) string {
	var act domain.AktivitasV2
	if err := db.WithContext(ctx).Where("id = ?", timeline.IdAktivitasV2).First(&act).Error; err != nil {
		return filepath.Join("./upload/prokerv2", "umum", "2026", "standar", "indikator", "aktivitas", sanitizeFolderSegment(timeline.NamaAktivitas))
	}

	var proker *domain.ProkerV3Default
	if act.IdMatrikProker != nil && *act.IdMatrikProker > 0 {
		var mp domain.ProkerV3
		if err := db.WithContext(ctx).Where("id = ?", *act.IdMatrikProker).First(&mp).Error; err == nil {
			proker, _ = repo.GetDefaultByUuid(ctx, mp.UUID)
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
// GET /proker_v3/aktivitas/:aktivitas_uuid (Detail aktivitas + timelines + files)
// =======================================================
func GetAktivitasDetailHandler(db *gorm.DB) fiber.Handler {
	repo := infra.NewProkerV3Repository(db)

	return func(c *fiber.Ctx) error {
		rawUuid := c.Params("aktivitas_uuid")
		if rawUuid == "" {
			rawUuid = c.Params("uuid")
		}
		parsedUuid, err := uuid.Parse(rawUuid)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid aktivitas UUID"})
		}

		act, err := repo.GetAktivitasByUuid(c.Context(), parsedUuid)
		if err != nil {
			return commoninfra.HandleError(c, err)
		}
		if act == nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Aktivitas tidak ditemukan"})
		}

		var proker *domain.ProkerV3Default
		if act.IdMatrikProker != nil && *act.IdMatrikProker > 0 {
			var mp domain.ProkerV3
			if err := db.Where("id = ?", *act.IdMatrikProker).First(&mp).Error; err == nil {
				proker, _ = repo.GetDefaultByUuid(c.Context(), mp.UUID)
			}
		}

		return c.JSON(fiber.Map{
			"aktivitas": act,
			"proker":    proker,
			"timelines": act.Timelines,
			"data":      act.Timelines,
		})
	}
}

var ensureDokumenOnce sync.Once

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

func ensureDokumenTable(ctx context.Context, db *gorm.DB) {
	ensureDokumenOnce.Do(func() {
		_ = db.AutoMigrate(&domain.MatrixProkerDokumen{})
		_ = db.WithContext(ctx).Exec(`
			CREATE TABLE IF NOT EXISTS matrix_proker_dokumen (
				id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
				uuid VARCHAR(36) NOT NULL UNIQUE,
				id_matrik_proker BIGINT UNSIGNED NOT NULL,
				jenis_dokumen VARCHAR(100) NOT NULL,
				nama_file VARCHAR(255) NOT NULL,
				file VARCHAR(500) NOT NULL,
				file_url VARCHAR(500) NOT NULL,
				file_size BIGINT DEFAULT 0,
				file_type VARCHAR(100) DEFAULT '',
				status_verifikasi VARCHAR(100) DEFAULT 'belum_terverifikasi',
				catatan TEXT NULL,
				catatan_penolakan TEXT NULL,
				keterangan TEXT NULL,
				created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
				updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
				deleted_at DATETIME NULL,
				KEY idx_matrik_proker (id_matrik_proker),
				KEY idx_deleted_at (deleted_at)
			) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
		`).Error
		_ = db.WithContext(ctx).Exec(`ALTER TABLE matrix_proker_dokumen ADD COLUMN catatan TEXT NULL;`).Error
		_ = db.WithContext(ctx).Exec(`ALTER TABLE matrix_proker_dokumen ADD COLUMN catatan_penolakan TEXT NULL;`).Error
		_ = db.WithContext(ctx).Exec(`ALTER TABLE matrix_proker_dokumen ADD COLUMN keterangan TEXT NULL;`).Error
		_ = db.WithContext(ctx).Exec(`ALTER TABLE aktivitasv2 ADD COLUMN realisasi DOUBLE NULL;`).Error
		_ = db.WithContext(ctx).Exec(`
			ALTER TABLE matrix_proker_dokumen 
			MODIFY COLUMN status_verifikasi VARCHAR(100) NULL DEFAULT 'belum_terverifikasi';
		`).Error
		_ = db.WithContext(ctx).Exec(`
			UPDATE matrix_proker_dokumen 
			SET catatan = keterangan 
			WHERE (catatan IS NULL OR catatan = '') AND (keterangan IS NOT NULL AND keterangan != '');
		`).Error
		_ = db.WithContext(ctx).Exec(`
			UPDATE matrix_proker_dokumen 
			SET keterangan = catatan 
			WHERE (keterangan IS NULL OR keterangan = '') AND (catatan IS NOT NULL AND catatan != '');
		`).Error
	})
}

// =======================================================
// GET /proker_v3/:uuid/dokumen (Daftar matrix_proker_dokumen)
// =======================================================
func GetProkerDokumenHandler(db *gorm.DB) fiber.Handler {
	repo := infra.NewProkerV3Repository(db)

	return func(c *fiber.Ctx) error {
		ensureDokumenTable(c.Context(), db)

		rawUuid := c.Params("uuid")
		parsedUuid, err := uuid.Parse(rawUuid)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid matrix UUID"})
		}

		proker, err := repo.GetDefaultByUuid(c.Context(), parsedUuid)
		if err != nil {
			return commoninfra.HandleError(c, err)
		}
		if proker == nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Program kerja tidak ditemukan"})
		}

		var docs []domain.MatrixProkerDokumen
		if err := db.Unscoped().Where("id_matrik_proker = ?", proker.ID).
			Order("id DESC").
			Find(&docs).Error; err != nil {
			return commoninfra.HandleError(c, err)
		}

		activities, _ := repo.GetAktivitasByMatrikId(c.Context(), proker.ID)

		return c.JSON(fiber.Map{
			"proker":     proker,
			"activities": activities,
			"data":       docs,
		})
	}
}

// =======================================================
// POST /proker_v3/:uuid/dokumen/chunk (Chunked binary upload)
// =======================================================
func UploadProkerDokumenChunkHandler(db *gorm.DB) fiber.Handler {
	repo := infra.NewProkerV3Repository(db)

	return func(c *fiber.Ctx) error {
		ensureDokumenTable(c.Context(), db)

		rawUuid := c.Params("uuid")
		parsedUuid, err := uuid.Parse(rawUuid)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid matrix UUID"})
		}

		proker, err := repo.GetDefaultByUuid(c.Context(), parsedUuid)
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

		// If this is the last chunk, finalize file and save record
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
			// Gunakan UUID dari tabel matrix_proker_dokumen sebagai nama file
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

			if err := db.Create(&doc).Error; err != nil {
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
}

// =======================================================
// DELETE /proker_v3/:uuid/dokumen/chunk/:upload_id (Cancel chunk upload)
// =======================================================
func CancelProkerDokumenChunkHandler(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		uploadIDStr := strings.TrimSpace(c.Params("upload_id"))
		if uploadIDStr == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Missing upload_id"})
		}
		partPath := filepath.Join("./upload/tmp", fmt.Sprintf("%s.part", uploadIDStr))
		_ = os.Remove(partPath)
		return c.JSON(fiber.Map{"message": "Upload chunk dibatalkan dan file temporary dihapus"})
	}
}

// =======================================================
// POST /proker_v3/timeline/:timeline_uuid/chunk (Timeline file chunk upload)
// =======================================================
func UploadTimelineFileChunkHandler(db *gorm.DB) fiber.Handler {
	repo := infra.NewProkerV3Repository(db)

	return func(c *fiber.Ctx) error {
		tlUUIDStr := c.Params("timeline_uuid")
		timelineUUID, err := uuid.Parse(tlUUIDStr)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid timeline UUID"})
		}

		timeline, err := repo.GetTimelineByUuid(c.Context(), timelineUUID)
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

		// If this is the last chunk, finalize file and save record
		if totalChunks > 0 && chunkIndex >= totalChunks-1 {
			targetDir := resolveTimelineUploadDir(db, repo, c.Context(), timeline)
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

			if err := db.Create(&record).Error; err != nil {
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
}

// =======================================================
// DELETE /proker_v3/timeline/:timeline_uuid/chunk/:upload_id
// =======================================================
func CancelTimelineFileChunkHandler(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		uploadIDStr := strings.TrimSpace(c.Params("upload_id"))
		if uploadIDStr == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Missing upload_id"})
		}
		partPath := filepath.Join("./upload/tmp", fmt.Sprintf("%s.part", uploadIDStr))
		_ = os.Remove(partPath)
		return c.JSON(fiber.Map{"message": "Upload berkas bukti dibatalkan dan file temporary dihapus"})
	}
}

// =======================================================
// POST /proker_v3/:uuid/dokumen (Upload ke matrix_proker_dokumen)
// Path: ./upload/prokerv2/[nama unit]/[tahun]/[standar]/[indikator]/[jenis file]
// =======================================================
func UploadProkerDokumenHandler(db *gorm.DB) fiber.Handler {
	repo := infra.NewProkerV3Repository(db)

	return func(c *fiber.Ctx) error {
		ensureDokumenTable(c.Context(), db)

		rawUuid := c.Params("uuid")
		parsedUuid, err := uuid.Parse(rawUuid)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid matrix UUID"})
		}

		proker, err := repo.GetDefaultByUuid(c.Context(), parsedUuid)
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
			// Gunakan UUID sebagai nama file di server
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

			if err := db.Create(&doc).Error; err != nil {
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
}

// =======================================================
// DELETE /proker_v3/dokumen/:doc_uuid (Hapus matrix_proker_dokumen)
// =======================================================
func DeleteProkerDokumenHandler(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		rawUuid := c.Params("doc_uuid")
		if rawUuid == "" {
			rawUuid = c.Params("uuid")
		}
		docUUID, err := uuid.Parse(rawUuid)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid document UUID"})
		}

		var doc domain.MatrixProkerDokumen
		if err := db.Unscoped().Where("uuid = ?", docUUID).First(&doc).Error; err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Dokumen tidak ditemukan"})
		}

		filePath := doc.File
		if err := db.Unscoped().Delete(&doc).Error; err != nil {
			return commoninfra.HandleError(c, err)
		}

		if filePath != "" {
			_ = os.Remove(filePath)
		}

		return c.JSON(fiber.Map{
			"message": "Dokumen berhasil dihapus",
			"uuid":    rawUuid,
		})
	}
}

// =======================================================
// PUT /proker_v3/dokumen/:doc_uuid/verifikasi
// =======================================================
func UpdateVerifikasiDokumenHandler(db *gorm.DB) fiber.Handler {
	type VerifikasiInput struct {
		StatusVerifikasi string `json:"status_verifikasi"`
		Catatan          string `json:"catatan"`
		CatatanPenolakan string `json:"catatan_penolakan"`
	}

	return func(c *fiber.Ctx) error {
		rawUuid := c.Params("doc_uuid")
		if rawUuid == "" {
			rawUuid = c.Params("uuid")
		}
		docUUID, err := uuid.Parse(rawUuid)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid document UUID"})
		}

		var doc domain.MatrixProkerDokumen
		if err := db.Where("uuid = ?", docUUID).First(&doc).Error; err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Dokumen tidak ditemukan"})
		}

		var input VerifikasiInput
		if err := c.BodyParser(&input); err != nil || input.StatusVerifikasi == "" {
			input.StatusVerifikasi = c.FormValue("status_verifikasi")
			if input.StatusVerifikasi == "" {
				input.StatusVerifikasi = c.FormValue("status")
			}
			if input.Catatan == "" {
				input.Catatan = c.FormValue("catatan")
			}
			if input.CatatanPenolakan == "" {
				input.CatatanPenolakan = c.FormValue("catatan_penolakan")
			}
		}

		status := strings.TrimSpace(input.StatusVerifikasi)
		if status == "" {
			status = "terverifikasi"
		}

		doc.StatusVerifikasi = status
		// Simpan alasan penolakan pada kolom CatatanPenolakan tanpa menimpa doc.Catatan (catatan pengunggah)
		if status == "ditolak" || status == "gagal_terverifikasi" {
			alasan := strings.TrimSpace(input.CatatanPenolakan)
			if alasan == "" {
				alasan = strings.TrimSpace(input.Catatan)
			}
			doc.CatatanPenolakan = alasan
		} else if status == "terverifikasi" {
			doc.CatatanPenolakan = ""
		}
		doc.UpdatedAt = time.Now()

		if err := db.Save(&doc).Error; err != nil {
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
}

func ModuleProkerV3(app *fiber.App, db *gorm.DB) {
	roles := []string{"admin", "auditee", "fakultas", "unit"}
	whoamiURL := "http://localhost:3000/whoami"

	// 1. Options & Target Unpak
	app.Get("/proker_v3_options", commonpresentation.SmartCompress(), commonpresentation.JWTMiddleware(), GetProkerV3OptionsHandler(db))
	app.Get("/proker_v3/target_unpak", commonpresentation.SmartCompress(), commonpresentation.JWTMiddleware(), GetTargetUnpakHandler(db))

	// 2. Program Kerja (matrix_proker) CRUD
	app.Post("/proker_v3", commonpresentation.JWTMiddleware(), commonpresentation.RBACMiddleware(roles, whoamiURL), CreateProkerV3Handlerfunc)
	app.Get("/proker_v3/:uuid", commonpresentation.SmartCompress(), commonpresentation.JWTMiddleware(), commonpresentation.RBACMiddleware(roles, whoamiURL), GetProkerV3Handlerfunc)
	app.Put("/proker_v3/:uuid", commonpresentation.JWTMiddleware(), commonpresentation.RBACMiddleware(roles, whoamiURL), UpdateProkerV3Handlerfunc)
	app.Delete("/proker_v3/:uuid", commonpresentation.JWTMiddleware(), commonpresentation.RBACMiddleware(roles, whoamiURL), DeleteProkerV3Handlerfunc)

	app.Get("/proker_v3", commonpresentation.SmartCompress(), commonpresentation.JWTMiddleware(), commonpresentation.RBACMiddleware(roles, whoamiURL), GetAllProkerV3Handlerfunc)
	app.Get("/proker_v3s", commonpresentation.SmartCompress(), commonpresentation.JWTMiddleware(), commonpresentation.RBACMiddleware(roles, whoamiURL), GetAllProkerV3Handlerfunc)

	// 3. Aktivitas (aktivitasv2) CRUD per Matrix Proker
	app.Get("/proker_v3/:uuid/aktivitas", commonpresentation.SmartCompress(), commonpresentation.JWTMiddleware(), commonpresentation.RBACMiddleware(roles, whoamiURL), GetAktivitasByMatrixHandler(db))
	app.Get("/proker_v3/aktivitas/:aktivitas_uuid", commonpresentation.SmartCompress(), commonpresentation.JWTMiddleware(), commonpresentation.RBACMiddleware(roles, whoamiURL), GetAktivitasDetailHandler(db))
	app.Get("/proker_v3/:uuid/aktivitas/:aktivitas_uuid", commonpresentation.SmartCompress(), commonpresentation.JWTMiddleware(), commonpresentation.RBACMiddleware(roles, whoamiURL), GetAktivitasDetailHandler(db))
	app.Post("/proker_v3/:uuid/aktivitas", commonpresentation.JWTMiddleware(), commonpresentation.RBACMiddleware(roles, whoamiURL), CreateAktivitasHandler(db))
	app.Put("/proker_v3/aktivitas/:aktivitas_uuid", commonpresentation.JWTMiddleware(), commonpresentation.RBACMiddleware(roles, whoamiURL), UpdateAktivitasHandler(db))
	app.Delete("/proker_v3/aktivitas/:aktivitas_uuid", commonpresentation.JWTMiddleware(), commonpresentation.RBACMiddleware(roles, whoamiURL), DeleteAktivitasHandler(db))
	app.Put("/proker_v3/aktivitas/:aktivitas_uuid/status", commonpresentation.JWTMiddleware(), commonpresentation.RBACMiddleware(roles, whoamiURL), UpdateAktivitasStatusHandler(db))
	app.Post("/proker_v3/aktivitas/:aktivitas_uuid/dokumen_referensi", commonpresentation.JWTMiddleware(), commonpresentation.RBACMiddleware(roles, whoamiURL), UploadDokumenReferensiHandler(db))

	// Backward-compat aliases
	app.Put("/proker_v3/:uuid/status", commonpresentation.JWTMiddleware(), commonpresentation.RBACMiddleware(roles, whoamiURL), UpdateAktivitasStatusHandler(db))
	app.Post("/proker_v3/:uuid/dokumen_referensi", commonpresentation.JWTMiddleware(), commonpresentation.RBACMiddleware(roles, whoamiURL), UploadDokumenReferensiHandler(db))

	// 4. Timeline CRUD & Upload files
	app.Post("/proker_v3/aktivitas/:aktivitas_uuid/timeline", commonpresentation.JWTMiddleware(), commonpresentation.RBACMiddleware(roles, whoamiURL), CreateTimelineHandler(db))
	app.Post("/proker_v3/:uuid/timeline", commonpresentation.JWTMiddleware(), commonpresentation.RBACMiddleware(roles, whoamiURL), CreateTimelineHandler(db))
	app.Put("/proker_v3/timeline/:timeline_uuid", commonpresentation.JWTMiddleware(), commonpresentation.RBACMiddleware(roles, whoamiURL), UpdateTimelineHandler(db))
	app.Put("/proker_v3/:uuid/timeline/:timeline_uuid", commonpresentation.JWTMiddleware(), commonpresentation.RBACMiddleware(roles, whoamiURL), UpdateTimelineHandler(db))
	app.Delete("/proker_v3/timeline/:timeline_uuid", commonpresentation.JWTMiddleware(), commonpresentation.RBACMiddleware(roles, whoamiURL), DeleteTimelineHandler(db))
	app.Delete("/proker_v3/:uuid/timeline/:timeline_uuid", commonpresentation.JWTMiddleware(), commonpresentation.RBACMiddleware(roles, whoamiURL), DeleteTimelineHandler(db))
	app.Put("/proker_v3/aktivitas/:aktivitas_uuid/timeline/reorder", commonpresentation.JWTMiddleware(), commonpresentation.RBACMiddleware(roles, whoamiURL), ReorderTimelineHandler(db))
	app.Put("/proker_v3/:uuid/timeline/reorder", commonpresentation.JWTMiddleware(), commonpresentation.RBACMiddleware(roles, whoamiURL), ReorderTimelineHandler(db))

	// Multiple file upload for timelines
	app.Post("/proker_v3/timeline/:timeline_uuid/files", commonpresentation.JWTMiddleware(), commonpresentation.RBACMiddleware(roles, whoamiURL), UploadTimelineFilesHandler(db))
	app.Post("/proker_v3/:uuid/timeline/:timeline_uuid/files", commonpresentation.JWTMiddleware(), commonpresentation.RBACMiddleware(roles, whoamiURL), UploadTimelineFilesHandler(db))
	app.Post("/proker_v3/timeline/:timeline_uuid/chunk", commonpresentation.JWTMiddleware(), commonpresentation.RBACMiddleware(roles, whoamiURL), UploadTimelineFileChunkHandler(db))
	app.Post("/proker_v3/:uuid/timeline/:timeline_uuid/chunk", commonpresentation.JWTMiddleware(), commonpresentation.RBACMiddleware(roles, whoamiURL), UploadTimelineFileChunkHandler(db))
	app.Delete("/proker_v3/timeline/:timeline_uuid/chunk/:upload_id", commonpresentation.JWTMiddleware(), commonpresentation.RBACMiddleware(roles, whoamiURL), CancelTimelineFileChunkHandler(db))
	app.Delete("/proker_v3/:uuid/timeline/:timeline_uuid/chunk/:upload_id", commonpresentation.JWTMiddleware(), commonpresentation.RBACMiddleware(roles, whoamiURL), CancelTimelineFileChunkHandler(db))
	app.Delete("/proker_v3/file/:file_uuid", commonpresentation.JWTMiddleware(), commonpresentation.RBACMiddleware(roles, whoamiURL), DeleteTimelineFileHandler(db))

	// 5. Matrix Proker Dokumen (Upload ke matrix_proker_dokumen)
	app.Get("/proker_v3/:uuid/dokumen", commonpresentation.SmartCompress(), commonpresentation.JWTMiddleware(), commonpresentation.RBACMiddleware(roles, whoamiURL), GetProkerDokumenHandler(db))
	app.Post("/proker_v3/:uuid/dokumen/chunk", commonpresentation.JWTMiddleware(), commonpresentation.RBACMiddleware(roles, whoamiURL), UploadProkerDokumenChunkHandler(db))
	app.Delete("/proker_v3/:uuid/dokumen/chunk/:upload_id", commonpresentation.JWTMiddleware(), commonpresentation.RBACMiddleware(roles, whoamiURL), CancelProkerDokumenChunkHandler(db))
	app.Post("/proker_v3/:uuid/dokumen", commonpresentation.JWTMiddleware(), commonpresentation.RBACMiddleware(roles, whoamiURL), UploadProkerDokumenHandler(db))
	app.Delete("/proker_v3/dokumen/:doc_uuid", commonpresentation.JWTMiddleware(), commonpresentation.RBACMiddleware(roles, whoamiURL), DeleteProkerDokumenHandler(db))
	app.Put("/proker_v3/dokumen/:doc_uuid/verifikasi", commonpresentation.JWTMiddleware(), commonpresentation.RBACMiddleware(roles, whoamiURL), UpdateVerifikasiDokumenHandler(db))
	app.Put("/proker_v3/dokumen/:doc_uuid/status", commonpresentation.JWTMiddleware(), commonpresentation.RBACMiddleware(roles, whoamiURL), UpdateVerifikasiDokumenHandler(db))
	app.Put("/proker_v3/:uuid/dokumen/:doc_uuid/verifikasi", commonpresentation.JWTMiddleware(), commonpresentation.RBACMiddleware(roles, whoamiURL), UpdateVerifikasiDokumenHandler(db))
}
