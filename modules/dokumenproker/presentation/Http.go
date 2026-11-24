package presentation

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/mehdihadeli/go-mediatr"
	"gorm.io/gorm"

	commondomain "UnpakSiamida/common/domain"
	"UnpakSiamida/common/helper"
	commoninfra "UnpakSiamida/common/infrastructure"
	commonpresentation "UnpakSiamida/common/presentation"
	DokumenProkerdomain "UnpakSiamida/modules/dokumenproker/domain"

	CreateDokumenProker "UnpakSiamida/modules/dokumenproker/application/CreateDokumenProker"
	DeleteDokumenProker "UnpakSiamida/modules/dokumenproker/application/DeleteDokumenProker"
	GetAllDokumenProkers "UnpakSiamida/modules/dokumenproker/application/GetAllDokumenProkers"
	GetDokumenProkerDefault "UnpakSiamida/modules/dokumenproker/application/GetDokumenProkerDefault"
	SetupUuidDokumenProker "UnpakSiamida/modules/dokumenproker/application/SetupUuidDokumenProker"
	UpdateDokumenProker "UnpakSiamida/modules/dokumenproker/application/UpdateDokumenProker"
)

var dokumenProkerDb *gorm.DB

func sanitizeFolderSegment(s string) string {
	s = strings.TrimSpace(s)
	if s == "" {
		return "LAINNYA"
	}
	invalidChars := []string{"/", "\\", ":", "*", "?", "\"", "<", ">", "|", "\n", "\r", "\t"}
	for _, char := range invalidChars {
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

func getDokumenJenisFolder(jenisDok string) string {
	j := strings.ToLower(strings.TrimSpace(jenisDok))
	switch {
	case strings.Contains(j, "sop"):
		return "sop"
	case strings.Contains(j, "proposal") || strings.Contains(j, "tor"):
		return "proposal_tor"
	case strings.Contains(j, "laporan"):
		return "laporan"
	case strings.Contains(j, "pendukung") || strings.Contains(j, "document") || strings.Contains(j, "dokumen"):
		return "dokumen_pendukung"
	case strings.Contains(j, "sk"):
		return "sk"
	default:
		return sanitizeFolderSegment(strings.ToLower(jenisDok))
	}
}

func getDokumenProkerUploadFolder(fakultasUnitVal, mataProgramVal, jenisDokVal string) string {
	namaUnit := "unit"
	tahun := "2027"
	namaProker := "proker"
	jenisFolder := getDokumenJenisFolder(jenisDokVal)

	if dokumenProkerDb != nil {
		mpStr := strings.TrimSpace(mataProgramVal)
		if mpStr != "" {
			var row struct {
				MataProgram  string
				Tahun        string
				FakultasUnit string
			}
			_ = dokumenProkerDb.Table("matrix_proker mpr").
				Select("COALESCE(mp.mata_program, '') as mata_program, COALESCE(mt.tahun, '') as tahun, COALESCE(vfu.nama_fak_prod_unit, '') as fakultas_unit").
				Joins("LEFT JOIN mata_programv2 mp ON mpr.id_mata_program = mp.id").
				Joins("LEFT JOIN master_tahun mt ON mp.id_master_tahun = mt.id").
				Joins("LEFT JOIN v_fakultas_unit vfu ON mpr.id_fakultas_unit = vfu.id").
				Where("mpr.uuid = ? OR mpr.id = ?", mpStr, mpStr).
				Take(&row).Error

			if row.MataProgram == "" {
				_ = dokumenProkerDb.Table("mata_programv2 mp").
					Select("mp.mata_program, COALESCE(mt.tahun, '') as tahun").
					Joins("LEFT JOIN master_tahun mt ON mp.id_master_tahun = mt.id").
					Where("mp.uuid = ? OR mp.id = ?", mpStr, mpStr).
					Take(&row).Error
			}

			if row.MataProgram == "" {
				_ = dokumenProkerDb.Table("mata_program mp").
					Select("mp.mata_program, mt.tahun, vfu.nama_fak_prod_unit as fakultas_unit").
					Joins("LEFT JOIN master_tahun mt ON mp.id_master_tahun = mt.id").
					Joins("LEFT JOIN v_fakultas_unit vfu ON mp.id_fakultas_unit = vfu.id").
					Where("mp.uuid = ? OR mp.id = ?", mpStr, mpStr).
					Take(&row).Error
			}

			if row.MataProgram != "" {
				namaProker = row.MataProgram
			}
			if row.Tahun != "" {
				tahun = row.Tahun
			}
			if row.FakultasUnit != "" {
				namaUnit = row.FakultasUnit
			}
		}

		fuStr := strings.TrimSpace(fakultasUnitVal)
		if fuStr != "" && (namaUnit == "unit" || namaUnit == "") {
			var uRow struct {
				NamaFakProdUnit string
			}
			_ = dokumenProkerDb.Table("v_fakultas_unit").
				Select("nama_fak_prod_unit").
				Where("uuid = ? OR id = ?", fuStr, fuStr).
				Take(&uRow).Error
			if uRow.NamaFakProdUnit != "" {
				namaUnit = uRow.NamaFakProdUnit
			}
		}
	}

	return filepath.Join(
		"./upload/prokerv3",
		sanitizeFolderSegment(namaUnit),
		sanitizeFolderSegment(tahun),
		sanitizeFolderSegment(namaProker),
		jenisFolder,
	)
}

// =======================================================
// POST /dokumenproker & /pelaporandokumen
// =======================================================

func CreateDokumenProkerHandlerfunc(c *fiber.Ctx) error {
	fakultas_unit := c.FormValue("fakultas_unit")
	mata_program := c.FormValue("mata_program")
	if mata_program == "" {
		mata_program = c.FormValue("id_mata_program")
	}
	jenis_dokumen := c.FormValue("jenis_dokumen")
	status := c.FormValue("status")
	if status == "" {
		status = c.FormValue("status_verifikasi")
	}
	catatan := helper.StrPtr(c.FormValue("catatan"))

	fileName := strings.TrimSpace(c.FormValue("file"))
	fileHeader, errFile := c.FormFile("file")
	if errFile == nil && fileHeader != nil {
		uploadDir := getDokumenProkerUploadFolder(fakultas_unit, mata_program, jenis_dokumen)
		if err := helper.EnsureDir775(uploadDir); err == nil {
			ext := filepath.Ext(fileHeader.Filename)
			if ext == "" {
				ext = ".pdf"
			}
			savedName := fmt.Sprintf("dokumen-%s%s", uuid.NewString(), ext)
			dest := filepath.Join(uploadDir, savedName)
			if err := c.SaveFile(fileHeader, dest); err == nil {
				_ = helper.EnsureFile664(dest)
				relURL := filepath.ToSlash(dest)
				if !strings.HasPrefix(relURL, "/") {
					relURL = "/" + strings.TrimPrefix(relURL, "./")
				}
				fileName = relURL
			}
		}
	}

	cmd := CreateDokumenProker.CreateDokumenProkerCommand{
		FakultasUuid:    fakultas_unit,
		MataProgramUuid: mata_program,
		JenisDokumen:    jenis_dokumen,
		File:            fileName,
		Status:          status,
		Catatan:         catatan,
	}

	uuidResult, err := mediatr.Send[CreateDokumenProker.CreateDokumenProkerCommand, string](context.Background(), cmd)
	if err != nil {
		return commoninfra.HandleError(c, err)
	}

	return commonpresentation.JsonUUID(c, uuidResult)
}

// =======================================================
// PUT /dokumenproker/{uuid} & /pelaporandokumen/{uuid}
// =======================================================

func UpdateDokumenProkerHandlerfunc(c *fiber.Ctx) error {
	uuidParam := c.Params("uuid")
	if uuidParam == "" {
		uuidParam = c.FormValue("uuid")
	}
	if uuidParam == "" {
		uuidParam = c.FormValue("id")
	}

	fakultas_unit := c.FormValue("fakultas_unit")
	mata_program := c.FormValue("mata_program")
	if mata_program == "" {
		mata_program = c.FormValue("id_mata_program")
	}
	jenis_dokumen := c.FormValue("jenis_dokumen")
	status := c.FormValue("status")
	if status == "" {
		status = c.FormValue("status_verifikasi")
	}
	catatan := helper.StrPtr(c.FormValue("catatan"))

	fileName := strings.TrimSpace(c.FormValue("file"))
	fileHeader, errFile := c.FormFile("file")
	if errFile == nil && fileHeader != nil {
		uploadDir := getDokumenProkerUploadFolder(fakultas_unit, mata_program, jenis_dokumen)
		if err := helper.EnsureDir775(uploadDir); err == nil {
			ext := filepath.Ext(fileHeader.Filename)
			if ext == "" {
				ext = ".pdf"
			}
			savedName := fmt.Sprintf("dokumen-%s%s", uuid.NewString(), ext)
			dest := filepath.Join(uploadDir, savedName)
			if err := c.SaveFile(fileHeader, dest); err == nil {
				_ = helper.EnsureFile664(dest)
				relURL := filepath.ToSlash(dest)
				if !strings.HasPrefix(relURL, "/") {
					relURL = "/" + strings.TrimPrefix(relURL, "./")
				}
				fileName = relURL
			}
		}
	}

	cmd := UpdateDokumenProker.UpdateDokumenProkerCommand{
		Uuid:            uuidParam,
		FakultasUuid:    fakultas_unit,
		MataProgramUuid: mata_program,
		JenisDokumen:    jenis_dokumen,
		File:            fileName,
		Status:          status,
		Catatan:         catatan,
	}

	updatedID, err := mediatr.Send[UpdateDokumenProker.UpdateDokumenProkerCommand, string](context.Background(), cmd)
	if err != nil {
		return commoninfra.HandleError(c, err)
	}

	return c.JSON(fiber.Map{"uuid": updatedID})
}

// =======================================================
// PUT/POST Verifikasi
// =======================================================

func VerifikasiDokumenProkerHandlerfunc(c *fiber.Ctx) error {
	uuidParam := c.Params("uuid")
	if uuidParam == "" {
		uuidParam = c.FormValue("uuid")
	}
	if uuidParam == "" {
		uuidParam = c.FormValue("id")
	}

	status := c.FormValue("status_verifikasi")
	if status == "" {
		status = c.FormValue("status")
	}
	catatan := helper.StrPtr(c.FormValue("catatan"))

	if strings.ToLower(strings.TrimSpace(status)) == "gagal_terverifikasi" {
		if catatan == nil || strings.TrimSpace(*catatan) == "" {
			return c.Status(400).JSON(fiber.Map{
				"code":    "DokumenProker.CatatanRequired",
				"message": "Catatan wajib diisi jika status gagal terverifikasi",
			})
		}
	}

	cmd := UpdateDokumenProker.UpdateDokumenProkerCommand{
		Uuid:    uuidParam,
		Status:  status,
		Catatan: catatan,
	}

	updatedID, err := mediatr.Send[UpdateDokumenProker.UpdateDokumenProkerCommand, string](context.Background(), cmd)
	if err != nil {
		return commoninfra.HandleError(c, err)
	}

	return c.JSON(fiber.Map{"uuid": updatedID, "message": "berhasil verifikasi"})
}

// =======================================================
// DELETE /dokumenproker/{uuid} & /pelaporandokumen/{uuid}
// =======================================================

func DeleteDokumenProkerHandlerfunc(c *fiber.Ctx) error {
	uuidParam := c.Params("uuid")
	if uuidParam == "" {
		uuidParam = c.FormValue("uuid")
	}

	cmd := DeleteDokumenProker.DeleteDokumenProkerCommand{
		Uuid: uuidParam,
	}

	deletedID, err := mediatr.Send[DeleteDokumenProker.DeleteDokumenProkerCommand, string](context.Background(), cmd)
	if err != nil {
		return commoninfra.HandleError(c, err)
	}

	return c.JSON(fiber.Map{"uuid": deletedID})
}

// =======================================================
// GET /DokumenProker/{uuid} & /PelaporanDokumen/{uuid}
// =======================================================

func GetDokumenProkerHandlerfunc(c *fiber.Ctx) error {
	uuidParam := c.Params("uuid")

	query := GetDokumenProkerDefault.GetDokumenProkerDefaultByUuidQuery{Uuid: uuidParam}

	dokumenProker, err := mediatr.Send[
		GetDokumenProkerDefault.GetDokumenProkerDefaultByUuidQuery,
		*DokumenProkerdomain.DokumenProkerDefault,
	](context.Background(), query)

	if err != nil {
		return commoninfra.HandleError(c, err)
	}

	if dokumenProker == nil {
		return c.Status(404).JSON(fiber.Map{"error": "DokumenProker not found"})
	}

	return c.JSON(dokumenProker)
}

// =======================================================
// GET /DokumenProkers & /PelaporanDokumens
// =======================================================

func GetAllDokumenProkersHandlerfunc(c *fiber.Ctx) error {
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

	fakultasUnit := strings.TrimSpace(c.FormValue("fakultas_unit"))
	if fakultasUnit == "" {
		fakultasUnit = strings.TrimSpace(c.FormValue("target"))
	}
	if fakultasUnit == "" {
		if fu, ok := c.Locals("fakultas_unit").(string); ok {
			fakultasUnit = fu
		}
	}

	// For non-admin (e.g. auditee), restrict to their own fakultas_unit forwarded from whoami or header
	if (currentRole == "auditee" || (currentRole != "admin" && currentRole != "")) && fakultasUnit != "" && fakultasUnit != "0" {
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

	filterStatus := strings.TrimSpace(c.Query("filter_status"))
	if filterStatus == "" {
		filterStatus = strings.TrimSpace(c.Query("status_verifikasi"))
	}
	if filterStatus == "" {
		filterStatus = strings.TrimSpace(c.Query("status"))
	}
	if filterStatus != "" && filterStatus != "all" && filterStatus != "semua" {
		filters = append(filters, commondomain.SearchFilter{
			Field:    "status_verifikasi",
			Operator: "eq",
			Value:    &filterStatus,
		})
	}

	sid := strings.TrimSpace(c.FormValue("sid"))
	if sid != "" {
		filters = append(filters, commondomain.SearchFilter{
			Field:    "sid",
			Operator: "eq",
			Value:    &sid,
		})
	}

	query := GetAllDokumenProkers.GetAllDokumenProkersQuery{
		Search:        search,
		SearchFilters: filters,
	}

	var adapter commonpresentation.OutputAdapter[DokumenProkerdomain.DokumenProkerDefault]
	switch mode {
	case "all":
		adapter = &commonpresentation.AllAdapter[DokumenProkerdomain.DokumenProkerDefault]{}
	case "ndjson":
		adapter = &commonpresentation.NDJSONAdapter[DokumenProkerdomain.DokumenProkerDefault]{}
	case "sse":
		adapter = &commonpresentation.SSEAdapter[DokumenProkerdomain.DokumenProkerDefault]{}
	default:
		query.Page = &page
		query.Limit = &limit
		adapter = &commonpresentation.PagingAdapter[DokumenProkerdomain.DokumenProkerDefault]{}
	}

	result, err := mediatr.Send[
		GetAllDokumenProkers.GetAllDokumenProkersQuery,
		commondomain.Paged[DokumenProkerdomain.DokumenProkerDefault],
	](context.Background(), query)

	if err != nil {
		return commoninfra.HandleError(c, err)
	}

	return adapter.Send(c, result)
}

func SetupUuidDokumenProkersHandlerfunc(c *fiber.Ctx) error {
	cmd := SetupUuidDokumenProker.SetupUuidDokumenProkerCommand{}

	message, err := mediatr.Send[SetupUuidDokumenProker.SetupUuidDokumenProkerCommand, string](context.Background(), cmd)
	if err != nil {
		return commoninfra.HandleError(c, err)
	}

	return c.JSON(fiber.Map{"message": message})
}

func ModuleDokumenProker(app *fiber.App, dbs ...*gorm.DB) {
	if len(dbs) > 0 {
		dokumenProkerDb = dbs[0]
	}
	roles := []string{"admin", "auditee", "auditor1", "auditor2", "fakultas", "unit", "user", "upml", "upmf", "upms"}
	whoamiURL := "http://localhost:3000/whoami"

	app.Get("/dokumenproker/setupuuid", SetupUuidDokumenProkersHandlerfunc)
	app.Get("/pelaporandokumen/setupuuid", SetupUuidDokumenProkersHandlerfunc)

	app.Post("/dokumenproker", commonpresentation.JWTMiddleware(), commonpresentation.RBACMiddleware(roles, whoamiURL), CreateDokumenProkerHandlerfunc)
	app.Post("/pelaporandokumen", commonpresentation.JWTMiddleware(), commonpresentation.RBACMiddleware(roles, whoamiURL), CreateDokumenProkerHandlerfunc)

	app.Put("/dokumenproker/:uuid", commonpresentation.JWTMiddleware(), commonpresentation.RBACMiddleware(roles, whoamiURL), UpdateDokumenProkerHandlerfunc)
	app.Put("/pelaporandokumen/:uuid", commonpresentation.JWTMiddleware(), commonpresentation.RBACMiddleware(roles, whoamiURL), UpdateDokumenProkerHandlerfunc)

	app.Put("/dokumenproker/:uuid/verifikasi", commonpresentation.JWTMiddleware(), commonpresentation.RBACMiddleware(roles, whoamiURL), VerifikasiDokumenProkerHandlerfunc)
	app.Put("/pelaporandokumen/:uuid/verifikasi", commonpresentation.JWTMiddleware(), commonpresentation.RBACMiddleware(roles, whoamiURL), VerifikasiDokumenProkerHandlerfunc)
	app.Post("/pelaporandokumen/verifikasi", commonpresentation.JWTMiddleware(), commonpresentation.RBACMiddleware(roles, whoamiURL), VerifikasiDokumenProkerHandlerfunc)

	app.Delete("/dokumenproker/:uuid", commonpresentation.JWTMiddleware(), commonpresentation.RBACMiddleware(roles, whoamiURL), DeleteDokumenProkerHandlerfunc)
	app.Delete("/pelaporandokumen/:uuid", commonpresentation.JWTMiddleware(), commonpresentation.RBACMiddleware(roles, whoamiURL), DeleteDokumenProkerHandlerfunc)

	app.Get("/DokumenProker/:uuid", commonpresentation.SmartCompress(), commonpresentation.JWTMiddleware(), commonpresentation.RBACMiddleware(roles, whoamiURL), GetDokumenProkerHandlerfunc)
	app.Get("/dokumenproker/:uuid", commonpresentation.SmartCompress(), commonpresentation.JWTMiddleware(), commonpresentation.RBACMiddleware(roles, whoamiURL), GetDokumenProkerHandlerfunc)
	app.Get("/PelaporanDokumen/:uuid", commonpresentation.SmartCompress(), commonpresentation.JWTMiddleware(), commonpresentation.RBACMiddleware(roles, whoamiURL), GetDokumenProkerHandlerfunc)
	app.Get("/pelaporandokumen/:uuid", commonpresentation.SmartCompress(), commonpresentation.JWTMiddleware(), commonpresentation.RBACMiddleware(roles, whoamiURL), GetDokumenProkerHandlerfunc)

	app.Get("/DokumenProkers", commonpresentation.SmartCompress(), commonpresentation.JWTMiddleware(), commonpresentation.RBACMiddleware(roles, whoamiURL), GetAllDokumenProkersHandlerfunc)
	app.Get("/dokumenprokers", commonpresentation.SmartCompress(), commonpresentation.JWTMiddleware(), commonpresentation.RBACMiddleware(roles, whoamiURL), GetAllDokumenProkersHandlerfunc)
	app.Get("/DokumenProkers/target/:target", commonpresentation.SmartCompress(), commonpresentation.JWTMiddleware(), commonpresentation.RBACMiddleware(roles, whoamiURL), GetAllDokumenProkersHandlerfunc)
	app.Get("/dokumenprokers/target/:target", commonpresentation.SmartCompress(), commonpresentation.JWTMiddleware(), commonpresentation.RBACMiddleware(roles, whoamiURL), GetAllDokumenProkersHandlerfunc)

	app.Get("/PelaporanDokumens", commonpresentation.SmartCompress(), commonpresentation.JWTMiddleware(), commonpresentation.RBACMiddleware(roles, whoamiURL), GetAllDokumenProkersHandlerfunc)
	app.Get("/pelaporandokumens", commonpresentation.SmartCompress(), commonpresentation.JWTMiddleware(), commonpresentation.RBACMiddleware(roles, whoamiURL), GetAllDokumenProkersHandlerfunc)
	app.Get("/PelaporanDokumens/target/:target", commonpresentation.SmartCompress(), commonpresentation.JWTMiddleware(), commonpresentation.RBACMiddleware(roles, whoamiURL), GetAllDokumenProkersHandlerfunc)
	app.Get("/pelaporandokumens/target/:target", commonpresentation.SmartCompress(), commonpresentation.JWTMiddleware(), commonpresentation.RBACMiddleware(roles, whoamiURL), GetAllDokumenProkersHandlerfunc)
}
