package presentation

import (
	"context"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/mehdihadeli/go-mediatr"

	commondomain "UnpakSiamida/common/domain"
	"UnpakSiamida/common/helper"
	commoninfra "UnpakSiamida/common/infrastructure"
	commonpresentation "UnpakSiamida/common/presentation"
	CreateSwot "UnpakSiamida/modules/swot/application/CreateSwot"
	DeleteSwot "UnpakSiamida/modules/swot/application/DeleteSwot"
	GetAllSwots "UnpakSiamida/modules/swot/application/GetAllSwots"
	GetSwot "UnpakSiamida/modules/swot/application/GetSwot"
	UpdateSwot "UnpakSiamida/modules/swot/application/UpdateSwot"
	domainSwot "UnpakSiamida/modules/swot/domain"
)


// =======================================================
// POST /swot
// =======================================================
func CreateSwotHandlerfunc(c *fiber.Ctx) error {
	level, userFU := commonpresentation.GetUserInfo(c)

	var cmd CreateSwot.CreateSwotCommand

	cmd.Kategori = strings.TrimSpace(c.FormValue("kategori"))
	cmd.Faktor = strings.TrimSpace(c.FormValue("faktor"))
	cmd.Deskripsi = strings.TrimSpace(c.FormValue("deskripsi"))
	cmd.Strategi = strings.TrimSpace(c.FormValue("strategi"))
	cmd.TahunUuid = strings.TrimSpace(c.FormValue("tahun_uuid"))
	if cmd.TahunUuid == "" {
		cmd.TahunUuid = strings.TrimSpace(c.FormValue("id_master_tahun"))
	}
	cmd.IdMasterTahun = helper.ParseOptionalInt(c.FormValue("id_master_tahun"))
	cmd.Tahun = strings.TrimSpace(c.FormValue("tahun"))

	if val := c.FormValue("bobot"); val != "" {
		cmd.Bobot = helper.ParseOptionalFloat(val)
	}
	if val := c.FormValue("rating"); val != "" {
		cmd.Rating = helper.ParseOptionalInt(val)
	}
	if val := c.FormValue("skor"); val != "" {
		cmd.Skor = helper.ParseOptionalFloat(val)
	}

	// Scoping: Non-admin (auditee) dipaksa ke unit miliknya sendiri
	if level != "admin" {
		cmd.FakultasUuid = userFU
		cmd.IdFakultasUnit = nil
	} else {
		cmd.FakultasUuid = strings.TrimSpace(c.FormValue("fakultas_unit"))
		if cmd.FakultasUuid == "" {
			cmd.FakultasUuid = strings.TrimSpace(c.FormValue("id_fakultas_unit"))
		}
		cmd.IdFakultasUnit = helper.ParseOptionalInt(c.FormValue("id_fakultas_unit"))
	}

	createdUUID, err := mediatr.Send[CreateSwot.CreateSwotCommand, string](context.Background(), cmd)
	if err != nil {
		return commoninfra.HandleError(c, err)
	}

	return commonpresentation.JsonUUID(c, createdUUID)
}

// =======================================================
// PUT /swot/:uuid
// =======================================================
func UpdateSwotHandlerfunc(c *fiber.Ctx) error {
	rawUuid := c.Params("uuid")
	parsedUuid, err := uuid.Parse(rawUuid)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid UUID"})
	}

	level, userFU := commonpresentation.GetUserInfo(c)

	// Scoping: non-admin harus memiliki item ini
	if level != "admin" {
		existing, err := mediatr.Send[GetSwot.GetSwotByUuidQuery, *domainSwot.SwotDefault](
			context.Background(),
			GetSwot.GetSwotByUuidQuery{UUID: parsedUuid},
		)
		if err != nil {
			return commoninfra.HandleError(c, err)
		}
		if existing == nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "SWOT not found"})
		}
		if !helper.MatchUnit(existing.IdFakultasUnit, existing.UuidFakultasUnit, userFU) {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Akses ditolak: bukan unit anda"})
		}
	}

	var cmd UpdateSwot.UpdateSwotCommand
	cmd.UUID = rawUuid
	cmd.Kategori = strings.TrimSpace(c.FormValue("kategori"))
	cmd.Faktor = strings.TrimSpace(c.FormValue("faktor"))
	cmd.Deskripsi = strings.TrimSpace(c.FormValue("deskripsi"))
	cmd.Strategi = strings.TrimSpace(c.FormValue("strategi"))
	cmd.TahunUuid = strings.TrimSpace(c.FormValue("tahun_uuid"))
	if cmd.TahunUuid == "" {
		cmd.TahunUuid = strings.TrimSpace(c.FormValue("id_master_tahun"))
	}
	cmd.IdMasterTahun = helper.ParseOptionalInt(c.FormValue("id_master_tahun"))
	cmd.Tahun = strings.TrimSpace(c.FormValue("tahun"))

	if val := c.FormValue("bobot"); val != "" {
		cmd.Bobot = helper.ParseOptionalFloat(val)
	}
	if val := c.FormValue("rating"); val != "" {
		cmd.Rating = helper.ParseOptionalInt(val)
	}
	if val := c.FormValue("skor"); val != "" {
		cmd.Skor = helper.ParseOptionalFloat(val)
	}

	if level != "admin" {
		cmd.FakultasUuid = userFU
		cmd.IdFakultasUnit = nil
	} else {
		cmd.FakultasUuid = strings.TrimSpace(c.FormValue("fakultas_unit"))
		if cmd.FakultasUuid == "" {
			cmd.FakultasUuid = strings.TrimSpace(c.FormValue("id_fakultas_unit"))
		}
		cmd.IdFakultasUnit = helper.ParseOptionalInt(c.FormValue("id_fakultas_unit"))
	}

	updatedUUID, err := mediatr.Send[UpdateSwot.UpdateSwotCommand, string](context.Background(), cmd)
	if err != nil {
		return commoninfra.HandleError(c, err)
	}

	return c.JSON(fiber.Map{"uuid": updatedUUID})
}

// =======================================================
// DELETE /swot/:uuid
// =======================================================
func DeleteSwotHandlerfunc(c *fiber.Ctx) error {
	rawUuid := c.Params("uuid")
	parsedUuid, err := uuid.Parse(rawUuid)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid UUID"})
	}

	level, userFU := commonpresentation.GetUserInfo(c)

	if level != "admin" {
		existing, err := mediatr.Send[GetSwot.GetSwotByUuidQuery, *domainSwot.SwotDefault](
			context.Background(),
			GetSwot.GetSwotByUuidQuery{UUID: parsedUuid},
		)
		if err != nil {
			return commoninfra.HandleError(c, err)
		}
		if existing == nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "SWOT not found"})
		}
		if !helper.MatchUnit(existing.IdFakultasUnit, existing.UuidFakultasUnit, userFU) {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Akses ditolak: bukan unit anda"})
		}
	}

	cmd := DeleteSwot.DeleteSwotCommand{UUID: rawUuid}
	deletedUUID, err := mediatr.Send[DeleteSwot.DeleteSwotCommand, string](context.Background(), cmd)
	if err != nil {
		return commoninfra.HandleError(c, err)
	}

	return c.JSON(fiber.Map{"uuid": deletedUUID})
}

// =======================================================
// GET /swot/:uuid
// =======================================================
func GetSwotHandlerfunc(c *fiber.Ctx) error {
	rawUuid := c.Params("uuid")
	parsedUuid, err := uuid.Parse(rawUuid)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid UUID"})
	}

	item, err := mediatr.Send[GetSwot.GetSwotByUuidQuery, *domainSwot.SwotDefault](
		context.Background(),
		GetSwot.GetSwotByUuidQuery{UUID: parsedUuid},
	)
	if err != nil {
		return commoninfra.HandleError(c, err)
	}
	if item == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "SWOT not found"})
	}

	level, userFU := commonpresentation.GetUserInfo(c)
	if level != "admin" && level != "evaluator" {
		if !helper.MatchUnit(item.IdFakultasUnit, item.UuidFakultasUnit, userFU) {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Akses ditolak: bukan unit anda"})
		}
	}

	return c.JSON(item)
}

// =======================================================
// GET /swots
// =======================================================
func GetAllSwotsHandlerfunc(c *fiber.Ctx) error {
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

	level, userFU := commonpresentation.GetUserInfo(c)

	var targetFU string
	if level != "admin" && level != "evaluator" {
		// Auditee: mutlak dibatasi hanya unit miliknya
		if userFU != "" && userFU != "0" {
			targetFU = userFU
		} else {
			targetFU = "NONE"
		}
	} else {
		// Admin / Evaluator: tampilkan semua atau filter pilihan admin
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
		if strings.Contains(targetFU, "-") {
			field = "fakultasunituuid"
		}
		filters = append(filters, commondomain.SearchFilter{
			Field:    field,
			Operator: "eq",
			Value:    &targetFU,
		})
	}

	// Filter Tahun
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

	// Filter Kategori
	filterKategori := strings.TrimSpace(c.Query("filter_kategori"))
	if filterKategori == "" {
		filterKategori = strings.TrimSpace(c.Query("kategori"))
	}
	if filterKategori != "" && filterKategori != "all" && filterKategori != "semua" {
		filters = append(filters, commondomain.SearchFilter{
			Field:    "kategori",
			Operator: "eq",
			Value:    &filterKategori,
		})
	}

	query := GetAllSwots.GetAllSwotsQuery{
		Search:        search,
		SearchFilters: filters,
	}

	var adapter commonpresentation.OutputAdapter[domainSwot.SwotDefault]
	switch mode {
	case "all":
		adapter = &commonpresentation.AllAdapter[domainSwot.SwotDefault]{}
	case "ndjson":
		adapter = &commonpresentation.NDJSONAdapter[domainSwot.SwotDefault]{}
	case "sse":
		adapter = &commonpresentation.SSEAdapter[domainSwot.SwotDefault]{}
	default:
		query.Page = &page
		query.Limit = &limit
		adapter = &commonpresentation.PagingAdapter[domainSwot.SwotDefault]{}
	}

	result, err := mediatr.Send[
		GetAllSwots.GetAllSwotsQuery,
		commondomain.Paged[domainSwot.SwotDefault],
	](context.Background(), query)

	if err != nil {
		return commoninfra.HandleError(c, err)
	}

	return adapter.Send(c, result)
}

func ModuleSwot(app *fiber.App) {
	roles := []string{"admin", "auditee", "fakultas", "unit", "upml", "upmf", "upms", "evaluator"}
	whoamiURL := "http://localhost:3000/whoami"

	app.Post("/swot", commonpresentation.JWTMiddleware(), commonpresentation.RBACMiddleware(roles, whoamiURL), CreateSwotHandlerfunc)
	app.Put("/swot/:uuid", commonpresentation.JWTMiddleware(), commonpresentation.RBACMiddleware(roles, whoamiURL), UpdateSwotHandlerfunc)
	app.Delete("/swot/:uuid", commonpresentation.JWTMiddleware(), commonpresentation.RBACMiddleware(roles, whoamiURL), DeleteSwotHandlerfunc)

	app.Get("/swot/:uuid", commonpresentation.SmartCompress(), commonpresentation.JWTMiddleware(), commonpresentation.RBACMiddleware(roles, whoamiURL), GetSwotHandlerfunc)
	app.Get("/Swot/:uuid", commonpresentation.SmartCompress(), commonpresentation.JWTMiddleware(), commonpresentation.RBACMiddleware(roles, whoamiURL), GetSwotHandlerfunc)

	app.Get("/swots", commonpresentation.SmartCompress(), commonpresentation.JWTMiddleware(), commonpresentation.RBACMiddleware(roles, whoamiURL), GetAllSwotsHandlerfunc)
	app.Get("/Swots", commonpresentation.SmartCompress(), commonpresentation.JWTMiddleware(), commonpresentation.RBACMiddleware(roles, whoamiURL), GetAllSwotsHandlerfunc)

	app.Get("/swots/target/:target", commonpresentation.SmartCompress(), commonpresentation.JWTMiddleware(), commonpresentation.RBACMiddleware(roles, whoamiURL), GetAllSwotsHandlerfunc)
	app.Get("/Swots/target/:target", commonpresentation.SmartCompress(), commonpresentation.JWTMiddleware(), commonpresentation.RBACMiddleware(roles, whoamiURL), GetAllSwotsHandlerfunc)
}
