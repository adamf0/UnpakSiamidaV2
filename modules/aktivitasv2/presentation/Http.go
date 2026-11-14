package presentation

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	commoninfra "UnpakSiamida/common/infrastructure"
	commonpresentation "UnpakSiamida/common/presentation"
	createAktivitas "UnpakSiamida/modules/aktivitasv2/application/CreateAktivitasV3"
	createTimeline "UnpakSiamida/modules/aktivitasv2/application/CreateTimelineV3"
	deleteAktivitas "UnpakSiamida/modules/aktivitasv2/application/DeleteAktivitasV3"
	deleteTimelineFile "UnpakSiamida/modules/aktivitasv2/application/DeleteTimelineFileV3"
	deleteTimeline "UnpakSiamida/modules/aktivitasv2/application/DeleteTimelineV3"
	getAktivitasByMatrix "UnpakSiamida/modules/aktivitasv2/application/GetAktivitasByMatrix"
	getAktivitasDetail "UnpakSiamida/modules/aktivitasv2/application/GetAktivitasDetail"
	reorderTimeline "UnpakSiamida/modules/aktivitasv2/application/ReorderTimelineV3"
	updateAktivitas "UnpakSiamida/modules/aktivitasv2/application/UpdateAktivitasV3"
	updateTimeline "UnpakSiamida/modules/aktivitasv2/application/UpdateTimelineV3"
	domain "UnpakSiamida/modules/aktivitasv2/domain"
	infra "UnpakSiamida/modules/aktivitasv2/infrastructure"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/mehdihadeli/go-mediatr"
	"gorm.io/gorm"
)

var (
	repoAktivitas domain.IAktivitasV2Repository
	aktivitasDb   *gorm.DB
)

// Helper pengecekan hak akses unit aktivitas
func canAccessAktivitasProkerUnit(c *fiber.Ctx, itemUnitID *int, itemUnitUUID string, itemUnitName string, parentFakultas ...string) bool {
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

	// 1. Direct match
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
	if len(assignedUnitIDs) == 0 && userFU != "" && (itemUnitUUID == userFU || (itemUnitID != nil && strconv.Itoa(*itemUnitID) == userFU)) {
		return true
	}

	// 2. Strict scoping untuk UPML, UPMS, dan USER biasa
	if level == "upml" || level == "upms" || level == "user" {
		return false
	}

	// 3. UPMF: Akses 1 tipe fakultas
	if level == "upmf" {
		if pFaculty != "" {
			for _, aname := range assignedUnitNames {
				if strings.EqualFold(pFaculty, aname) {
					return true
				}
			}
		}

		if aktivitasDb != nil && len(assignedUnitIDs) > 0 {
			var matchCount int64
			var q *gorm.DB
			if itemUnitID != nil && *itemUnitID > 0 {
				q = aktivitasDb.Table("v_fakultas_unit fu").
					Where("fu.id = ? AND (fu.id IN (?) OR (fu.fakultas IS NOT NULL AND fu.fakultas != '' AND LOWER(TRIM(fu.fakultas)) IN (SELECT LOWER(TRIM(COALESCE(nama_fak_prod_unit, ''))) FROM v_fakultas_unit WHERE id IN (?) AND type = 'fakultas')))", *itemUnitID, assignedUnitIDs, assignedUnitIDs)
			} else if itemUnitUUID != "" {
				q = aktivitasDb.Table("v_fakultas_unit fu").
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

// Helper sanitize string folder segment
func sanitizeFolderSegment(s string) string {
	s = strings.TrimSpace(s)
	if s == "" {
		return "general"
	}
	replacer := strings.NewReplacer(
		"/", "_", "\\", "_", ":", "_", "*", "_", "?", "_",
		"\"", "_", "<", "_", ">", "_", "|", "_", " ", "_",
	)
	clean := replacer.Replace(s)
	for strings.Contains(clean, "__") {
		clean = strings.ReplaceAll(clean, "__", "_")
	}
	clean = strings.Trim(clean, "._")
	if clean == "" {
		return "general"
	}
	return clean
}

// GET /proker_v3/:uuid/aktivitas
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

	var pUnitID *int
	if result.Proker.IdFakultasUnit != nil {
		val := int(*result.Proker.IdFakultasUnit)
		pUnitID = &val
	}
	if !canAccessAktivitasProkerUnit(c, pUnitID, result.Proker.UuidFakultasUnit, result.Proker.NamaFakultasUnit, result.Proker.ParentFakultas) {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"code":  "FORBIDDEN",
			"error": "Akses ditolak: Anda tidak memiliki hak akses untuk melihat aktivitas program kerja unit ini",
		})
	}

	return c.JSON(fiber.Map{
		"proker": result.Proker,
		"data":   result.Data,
		"total":  result.Total,
	})
}

// GET /proker_v3/aktivitas/:aktivitas_uuid
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
	if result == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Aktivitas tidak ditemukan"})
	}

	if result.Proker != nil {
		var pUnitID *int
		if result.Proker.IdFakultasUnit != nil {
			val := int(*result.Proker.IdFakultasUnit)
			pUnitID = &val
		}
		if !canAccessAktivitasProkerUnit(c, pUnitID, result.Proker.UuidFakultasUnit, result.Proker.NamaFakultasUnit, result.Proker.ParentFakultas) {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"code":  "FORBIDDEN",
				"error": "Akses ditolak: Anda tidak memiliki hak akses untuk melihat aktivitas program kerja unit ini",
			})
		}
	}


	return c.JSON(fiber.Map{
		"aktivitas": result.Aktivitas,
		"proker":    result.Proker,
		"timelines": result.Timelines,
		"data":      result.Data,
	})
}

// POST /proker_v3/:uuid/aktivitas
func CreateAktivitasHandler(c *fiber.Ctx) error {
	rawUuid := c.Params("uuid")

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

// PUT /proker_v3/aktivitas/:aktivitas_uuid
func UpdateAktivitasHandler(c *fiber.Ctx) error {
	rawUuid := c.Params("aktivitas_uuid")
	if rawUuid == "" {
		rawUuid = c.Params("uuid")
	}

	var cmd updateAktivitas.UpdateAktivitasV3Command
	cmd.UUID = rawUuid
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
	if val := c.FormValue("realisasi"); val != "" {
		if realVal, err := strconv.ParseFloat(val, 64); err == nil {
			cmd.Realisasi = &realVal
		}
	}
	if c.Request().PostArgs().Has("catatan_revisi_anggaran") || c.FormValue("catatan_revisi_anggaran") != "" {
		val := strings.TrimSpace(c.FormValue("catatan_revisi_anggaran"))
		cmd.CatatanRevisiAnggaran = &val
	}
	if val := c.FormValue("pencairan"); val != "" {
		if pencairanVal, err := strconv.ParseFloat(val, 64); err == nil {
			cmd.Pencairan = &pencairanVal
		}
	}
	if val := c.FormValue("status"); val != "" {
		cmd.Status = strings.TrimSpace(val)
	}

	act, err := mediatr.Send[updateAktivitas.UpdateAktivitasV3Command, *domain.AktivitasV2](c.Context(), cmd)
	if err != nil {
		return commoninfra.HandleError(c, err)
	}


	return c.JSON(fiber.Map{
		"message": "Aktivitas berhasil diperbarui",
		"uuid":    act.UUID.String(),
		"data":    act,
	})
}

// DELETE /proker_v3/aktivitas/:aktivitas_uuid
func DeleteAktivitasHandler(c *fiber.Ctx) error {
	rawUuid := c.Params("aktivitas_uuid")
	if rawUuid == "" {
		rawUuid = c.Params("uuid")
	}

	msg, err := mediatr.Send[deleteAktivitas.DeleteAktivitasV3Command, string](c.Context(), deleteAktivitas.DeleteAktivitasV3Command{UUID: rawUuid})
	if err != nil {
		return commoninfra.HandleError(c, err)
	}


	return c.JSON(fiber.Map{
		"message": msg,
		"uuid":    rawUuid,
	})
}

// POST /proker_v3/aktivitas/:aktivitas_uuid/timeline
func CreateTimelineHandler(c *fiber.Ctx) error {
	rawActUuid := c.Params("aktivitas_uuid")
	if rawActUuid == "" {
		rawActUuid = c.Params("uuid")
	}

	var cmd createTimeline.CreateTimelineV3Command
	cmd.AktivitasUUID = rawActUuid
	cmd.NamaAktivitas = strings.TrimSpace(c.FormValue("nama_aktivitas"))
	cmd.PIC = strings.TrimSpace(c.FormValue("pic"))
	cmd.TitikKontrol = strings.TrimSpace(c.FormValue("titik_kontrol"))
	cmd.JadwalRk = strings.TrimSpace(c.FormValue("jadwal_rk"))
	cmd.JadwalR = strings.TrimSpace(c.FormValue("jadwal_r"))
	cmd.Capaian = strings.TrimSpace(c.FormValue("capaian"))
	cmd.Masalah = strings.TrimSpace(c.FormValue("masalah"))
	cmd.PemecahanMasalah = strings.TrimSpace(c.FormValue("pemecahan_masalah"))
	cmd.Status = strings.TrimSpace(c.FormValue("status"))
	cmd.TanggalRkAwal = strings.TrimSpace(c.FormValue("tanggal_rk_awal"))
	cmd.TanggalRkAkhir = strings.TrimSpace(c.FormValue("tanggal_rk_akhir"))
	cmd.TanggalRAwal = strings.TrimSpace(c.FormValue("tanggal_r_awal"))
	cmd.TanggalRAkhir = strings.TrimSpace(c.FormValue("tanggal_r_akhir"))

	if val := c.FormValue("anggaran"); val != "" {
		if ang, err := strconv.ParseFloat(val, 64); err == nil {
			cmd.Anggaran = ang
		}
	}
	if val := c.FormValue("no_aktivitas"); val != "" {
		if no, err := strconv.Atoi(val); err == nil {
			cmd.NoAktivitas = no
		}
	}

	tl, err := mediatr.Send[createTimeline.CreateTimelineV3Command, *domain.ProkerTimelineV3](c.Context(), cmd)
	if err != nil {
		return commoninfra.HandleError(c, err)
	}


	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Timeline berhasil ditambahkan",
		"uuid":    tl.UUID.String(),
		"data":    tl,
	})
}

// PUT /proker_v3/timeline/:timeline_uuid
func UpdateTimelineHandler(c *fiber.Ctx) error {
	rawTlUuid := c.Params("timeline_uuid")

	var cmd updateTimeline.UpdateTimelineV3Command
	cmd.UUID = rawTlUuid
	cmd.NamaAktivitas = strings.TrimSpace(c.FormValue("nama_aktivitas"))
	cmd.PIC = strings.TrimSpace(c.FormValue("pic"))
	cmd.TitikKontrol = strings.TrimSpace(c.FormValue("titik_kontrol"))
	cmd.JadwalRk = strings.TrimSpace(c.FormValue("jadwal_rk"))
	cmd.JadwalR = strings.TrimSpace(c.FormValue("jadwal_r"))
	cmd.Capaian = strings.TrimSpace(c.FormValue("capaian"))
	cmd.Masalah = strings.TrimSpace(c.FormValue("masalah"))
	cmd.PemecahanMasalah = strings.TrimSpace(c.FormValue("pemecahan_masalah"))
	cmd.Status = strings.TrimSpace(c.FormValue("status"))
	cmd.TanggalRkAwal = strings.TrimSpace(c.FormValue("tanggal_rk_awal"))
	cmd.TanggalRkAkhir = strings.TrimSpace(c.FormValue("tanggal_rk_akhir"))
	cmd.TanggalRAwal = strings.TrimSpace(c.FormValue("tanggal_r_awal"))
	cmd.TanggalRAkhir = strings.TrimSpace(c.FormValue("tanggal_r_akhir"))

	if val := c.FormValue("anggaran"); val != "" {
		if ang, err := strconv.ParseFloat(val, 64); err == nil {
			cmd.Anggaran = &ang
		}
	}
	if val := c.FormValue("no_aktivitas"); val != "" {
		if no, err := strconv.Atoi(val); err == nil {
			cmd.NoAktivitas = no
		}
	}

	tl, err := mediatr.Send[updateTimeline.UpdateTimelineV3Command, *domain.ProkerTimelineV3](c.Context(), cmd)
	if err != nil {
		return commoninfra.HandleError(c, err)
	}


	return c.JSON(fiber.Map{
		"message": "Timeline berhasil diperbarui",
		"uuid":    tl.UUID.String(),
		"data":    tl,
	})
}

// POST /proker_v3/timeline/reorder
func ReorderTimelineHandler(c *fiber.Ctx) error {
	type reorderBody struct {
		Items []domain.ReorderTimelineItem `json:"items"`
	}
	var body reorderBody
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid body"})
	}

	msg, err := mediatr.Send[reorderTimeline.ReorderTimelineV3Command, string](c.Context(), reorderTimeline.ReorderTimelineV3Command{Items: body.Items})
	if err != nil {
		return commoninfra.HandleError(c, err)
	}


	return c.JSON(fiber.Map{"message": msg})
}

// DELETE /proker_v3/timeline/:timeline_uuid
func DeleteTimelineHandler(c *fiber.Ctx) error {
	rawUuid := c.Params("timeline_uuid")
	msg, err := mediatr.Send[deleteTimeline.DeleteTimelineV3Command, string](c.Context(), deleteTimeline.DeleteTimelineV3Command{UUID: rawUuid})
	if err != nil {
		return commoninfra.HandleError(c, err)
	}

	return c.JSON(fiber.Map{
		"message": msg,
		"uuid":    rawUuid,
	})
}

// DELETE /proker_v3/file/:file_uuid
func DeleteTimelineFileHandler(c *fiber.Ctx) error {
	rawUuid := c.Params("file_uuid")
	msg, err := mediatr.Send[deleteTimelineFile.DeleteTimelineFileV3Command, string](c.Context(), deleteTimelineFile.DeleteTimelineFileV3Command{FileUUID: rawUuid})
	if err != nil {
		return commoninfra.HandleError(c, err)
	}

	return c.JSON(fiber.Map{
		"message": msg,
		"uuid":    rawUuid,
	})
}

// UploadTimelineFilesHandler handles file uploads for timeline
func UploadTimelineFilesHandler(c *fiber.Ctx) error {
	tlUuidStr := c.Params("timeline_uuid")
	tlUuid, err := uuid.Parse(tlUuidStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid timeline UUID"})
	}

	if repoAktivitas == nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Repository not initialized"})
	}

	timeline, err := repoAktivitas.GetTimelineByUuid(c.Context(), tlUuid)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Timeline not found"})
	}

	form, err := c.MultipartForm()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Failed to parse multipart form"})
	}

	files := form.File["files"]
	if len(files) == 0 {
		files = form.File["file"]
	}
	if len(files) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "No files provided"})
	}

	uploadDir := filepath.Join("upload", "timeline_files", fmt.Sprintf("%d", timeline.IdAktivitasV2))
	_ = os.MkdirAll(uploadDir, 0755)

	var savedFiles []domain.ProkerTimelineFileV3
	for _, f := range files {
		fileUid := uuid.New()
		ext := filepath.Ext(f.Filename)
		cleanName := strings.TrimSuffix(f.Filename, ext)
		savedFilename := fmt.Sprintf("%s_%s%s", sanitizeFolderSegment(cleanName), fileUid.String()[:8], ext)
		destPath := filepath.Join(uploadDir, savedFilename)

		if err := c.SaveFile(f, destPath); err != nil {
			continue
		}

		savedFiles = append(savedFiles, domain.ProkerTimelineFileV3{
			UUID:                  fileUid,
			IdAktivitasV2Timeline: &timeline.ID,
			IdAktivitasV2:         &timeline.IdAktivitasV2,
			File:                  destPath,
			NamaFile:              f.Filename,
			FileURL:               "/" + filepath.ToSlash(destPath),
			FilePath:              destPath,
			FileSize:              f.Size,
			FileType:              f.Header.Get("Content-Type"),
			CreatedAt:             time.Now(),
			UpdatedAt:             time.Now(),
		})
	}

	if len(savedFiles) > 0 {
		_ = repoAktivitas.AddTimelineFiles(c.Context(), savedFiles)
	}

	return c.JSON(fiber.Map{
		"message": fmt.Sprintf("Berhasil upload %d file", len(savedFiles)),
		"files":   savedFiles,
	})
}

// UpdateAktivitasStatusHandler
func UpdateAktivitasStatusHandler(c *fiber.Ctx) error {
	rawUuid := c.Params("aktivitas_uuid")
	if rawUuid == "" {
		rawUuid = c.Params("uuid")
	}
	if rawUuid == "" {
		rawUuid = c.Query("uuid")
	}

	catatanPenolakan := strings.TrimSpace(c.FormValue("catatan_penolakan"))
	if catatanPenolakan == "" {
		catatanPenolakan = strings.TrimSpace(c.FormValue("alasan"))
	}
	if catatanPenolakan == "" {
		catatanPenolakan = strings.TrimSpace(c.FormValue("catatan"))
	}

	_, _ = commonpresentation.GetUserInfo(c)
	userName := ""
	if u, ok := c.Locals("user").(*commonpresentation.Account); ok && u != nil {
		userName = u.Name
		if userName == "" {
			userName = u.NidnUsername
		}
	}

	parsedUuid, err := uuid.Parse(rawUuid)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid aktivitas UUID"})
	}

	req := domain.UpdateStatusRequest{
		UUID:             parsedUuid,
		EntityType:       "aktivitas",
		Status:           c.FormValue("status"),
		Role:             c.FormValue("role"),
		CatatanPenolakan: catatanPenolakan,
		UserName:         userName,
	}

	if repoAktivitas == nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Repository not initialized"})
	}

	res, err := repoAktivitas.UpdateStatus(c.Context(), req)
	if err != nil {
		return commoninfra.HandleError(c, err)
	}


	return c.JSON(res)
}

// UpdateBatalAktivitasHandler
func UpdateBatalAktivitasHandler(c *fiber.Ctx) error {
	rawActUuid := c.Params("aktivitas_uuid")
	if rawActUuid == "" {
		rawActUuid = c.Params("uuid")
	}
	parsedUuid, err := uuid.Parse(rawActUuid)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid aktivitas UUID"})
	}

	if repoAktivitas == nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Repository not initialized"})
	}

	act, err := repoAktivitas.GetAktivitasEntityByUuid(c.Context(), parsedUuid)
	if err != nil || act == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Aktivitas tidak ditemukan"})
	}

	var newBatal bool
	if val := strings.TrimSpace(c.FormValue("is_batal")); val != "" {
		newBatal = val == "1" || strings.ToLower(val) == "true"
	} else {
		newBatal = !act.IsBatal
	}

	catatanBatal := strings.TrimSpace(c.FormValue("catatan_batal"))
	act.IsBatal = newBatal
	act.CatatanBatal = catatanBatal
	act.UpdatedAt = time.Now()

	if err := repoAktivitas.UpdateAktivitas(c.Context(), act); err != nil {
		return commoninfra.HandleError(c, err)
	}

	return c.JSON(fiber.Map{
		"message":       "Status pembatalan aktivitas berhasil diperbarui",
		"is_batal":      act.IsBatal,
		"catatan_batal": act.CatatanBatal,
	})
}

// UploadDokumenReferensiHandler
func UploadDokumenReferensiHandler(c *fiber.Ctx) error {
	rawUuid := c.Params("aktivitas_uuid")
	if rawUuid == "" {
		rawUuid = c.Params("uuid")
	}
	parsedUuid, err := uuid.Parse(rawUuid)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid UUID"})
	}

	if repoAktivitas == nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Repository not initialized"})
	}

	act, _ := repoAktivitas.GetAktivitasEntityByUuid(c.Context(), parsedUuid)

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
	_ = os.MkdirAll(uploadDir, 0755)

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

	if len(filesToInsert) > 0 {
		_ = repoAktivitas.AddTimelineFiles(c.Context(), filesToInsert)
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Dokumen referensi berhasil diunggah",
		"count":   len(filesToInsert),
		"files":   filesToInsert,
	})
}

// UploadTimelineFileChunkHandler
func UploadTimelineFileChunkHandler(c *fiber.Ctx) error {
	tlUUIDStr := c.Params("timeline_uuid")
	timelineUUID, err := uuid.Parse(tlUUIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid timeline UUID"})
	}

	if repoAktivitas == nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Repository not initialized"})
	}

	timeline, err := repoAktivitas.GetTimelineByUuid(c.Context(), timelineUUID)
	if err != nil || timeline == nil {
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
	totalSizeRaw := strings.TrimSpace(c.FormValue("total_size"))
	var totalSize int64
	if totalSizeRaw != "" {
		parsedTotalSize, parseErr := strconv.ParseInt(totalSizeRaw, 10, 64)
		if parseErr != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Ukuran total file tidak valid"})
		}
		totalSize = parsedTotalSize
	}
	fileName := strings.TrimSpace(c.FormValue("file_name"))
	keterangan := strings.TrimSpace(c.FormValue("keterangan"))

	fileHeader, err := c.FormFile("chunk")
	if err != nil {
		fileHeader, err = c.FormFile("file")
	}
	if err != nil || fileHeader == nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Tidak ada binary chunk yang dikirim"})
	}

	effectiveName := fileName
	if effectiveName == "" {
		effectiveName = fileHeader.Filename
	}
	expectedSize := totalSize
	if expectedSize <= 0 {
		expectedSize = fileHeader.Size
	}

	tmpDir := "./upload/tmp"
	partPath := filepath.Join(tmpDir, fmt.Sprintf("%s.part", fileUUID.String()))

	if err := ValidateTimelineUploadFile(effectiveName, expectedSize); err != nil {
		_ = os.Remove(partPath)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	if chunkOffset < 0 || chunkOffset+fileHeader.Size > MaxTimelineUploadBytes {
		_ = os.Remove(partPath)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Ukuran file melebihi batas maksimal 150 MB"})
	}

	_ = os.MkdirAll(tmpDir, 0755)

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
		if stat, statErr := os.Stat(partPath); statErr == nil {
			if stat.Size() > MaxTimelineUploadBytes {
				_ = os.Remove(partPath)
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Ukuran file melebihi batas maksimal 150 MB"})
			}
			if totalSize <= 0 {
				totalSize = stat.Size()
			}
		}

		targetDir := "./upload/proker_v3"
		_ = os.MkdirAll(targetDir, 0755)

		ext := filepath.Ext(effectiveName)
		savedFileName := fmt.Sprintf("proker-%s%s", fileUUID.String(), ext)
		destPath := filepath.Join(targetDir, savedFileName)

		if err := os.Rename(partPath, destPath); err != nil {
			_ = os.Remove(partPath)
		}

		relURL := "/upload/proker_v3/" + savedFileName

		record := domain.ProkerTimelineFileV3{
			UUID:                  fileUUID,
			IdAktivitasV2Timeline: &timeline.ID,
			IdAktivitasV2:         &timeline.IdAktivitasV2,
			File:                  destPath,
			NamaFile:              ClampRunes(effectiveName, MaxTimelineNamaFileRunes),
			FileURL:               relURL,
			FileSize:              totalSize,
			FileType:              ClampRunes(fileHeader.Header.Get("Content-Type"), 100),
			Keterangan:            keterangan,
		}

		_ = repoAktivitas.AddTimelineFiles(c.Context(), []domain.ProkerTimelineFileV3{record})

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

// CancelTimelineFileChunkHandler
func CancelTimelineFileChunkHandler(c *fiber.Ctx) error {
	uploadIDStr := strings.TrimSpace(c.Params("upload_id"))
	fileUUID, err := uuid.Parse(uploadIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid upload_id (must be UUID)"})
	}
	partPath := filepath.Join("./upload/tmp", fmt.Sprintf("%s.part", fileUUID.String()))
	_ = os.Remove(partPath)
	return c.JSON(fiber.Map{"message": "Upload berkas bukti dibatalkan dan file temporary dihapus"})
}

func ModuleAktivitasV2(app *fiber.App, dbs ...*gorm.DB) {
	if len(dbs) > 0 && dbs[0] != nil {
		aktivitasDb = dbs[0]
		repoAktivitas = infra.NewAktivitasV2Repository(dbs[0])
	}

	roles := []string{"admin", "user", "upml", "upmf", "upms", "evaluator", "anggaran", "keuangan"}
	whoamiURL := "http://localhost:3000/whoami"

	// Aktivitas routes
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
	app.Post("/proker_v3/:uuid/dokumen_referensi", commonpresentation.JWTMiddleware(), commonpresentation.RBACMiddleware(roles, whoamiURL), UploadDokumenReferensiHandler)

	// Timeline routes
	app.Post("/proker_v3/aktivitas/:aktivitas_uuid/timeline", commonpresentation.JWTMiddleware(), commonpresentation.RBACMiddleware(roles, whoamiURL), CreateTimelineHandler)
	app.Post("/proker_v3/:uuid/timeline", commonpresentation.JWTMiddleware(), commonpresentation.RBACMiddleware(roles, whoamiURL), CreateTimelineHandler)
	app.Put("/proker_v3/timeline/:timeline_uuid", commonpresentation.JWTMiddleware(), commonpresentation.RBACMiddleware(roles, whoamiURL), UpdateTimelineHandler)
	app.Put("/proker_v3/:uuid/timeline/:timeline_uuid", commonpresentation.JWTMiddleware(), commonpresentation.RBACMiddleware(roles, whoamiURL), UpdateTimelineHandler)
	app.Delete("/proker_v3/timeline/:timeline_uuid", commonpresentation.JWTMiddleware(), commonpresentation.RBACMiddleware(roles, whoamiURL), DeleteTimelineHandler)
	app.Delete("/proker_v3/:uuid/timeline/:timeline_uuid", commonpresentation.JWTMiddleware(), commonpresentation.RBACMiddleware(roles, whoamiURL), DeleteTimelineHandler)
	app.Put("/proker_v3/aktivitas/:aktivitas_uuid/timeline/reorder", commonpresentation.JWTMiddleware(), commonpresentation.RBACMiddleware(roles, whoamiURL), ReorderTimelineHandler)
	app.Put("/proker_v3/:uuid/timeline/reorder", commonpresentation.JWTMiddleware(), commonpresentation.RBACMiddleware(roles, whoamiURL), ReorderTimelineHandler)

	// Files
	app.Post("/proker_v3/timeline/:timeline_uuid/files", commonpresentation.JWTMiddleware(), commonpresentation.RBACMiddleware(roles, whoamiURL), UploadTimelineFilesHandler)
	app.Post("/proker_v3/:uuid/timeline/:timeline_uuid/files", commonpresentation.JWTMiddleware(), commonpresentation.RBACMiddleware(roles, whoamiURL), UploadTimelineFilesHandler)
	app.Post("/proker_v3/timeline/:timeline_uuid/chunk", commonpresentation.JWTMiddleware(), commonpresentation.RBACMiddleware(roles, whoamiURL), UploadTimelineFileChunkHandler)
	app.Post("/proker_v3/:uuid/timeline/:timeline_uuid/chunk", commonpresentation.JWTMiddleware(), commonpresentation.RBACMiddleware(roles, whoamiURL), UploadTimelineFileChunkHandler)
	app.Delete("/proker_v3/timeline/:timeline_uuid/chunk/:upload_id", commonpresentation.JWTMiddleware(), commonpresentation.RBACMiddleware(roles, whoamiURL), CancelTimelineFileChunkHandler)
	app.Delete("/proker_v3/:uuid/timeline/:timeline_uuid/chunk/:upload_id", commonpresentation.JWTMiddleware(), commonpresentation.RBACMiddleware(roles, whoamiURL), CancelTimelineFileChunkHandler)
	app.Delete("/proker_v3/file/:file_uuid", commonpresentation.JWTMiddleware(), commonpresentation.RBACMiddleware(roles, whoamiURL), DeleteTimelineFileHandler)
}
