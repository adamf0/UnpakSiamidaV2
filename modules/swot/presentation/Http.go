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
	CreateSwot "UnpakSiamida/modules/swot/application/CreateSwot"
	DeleteSwot "UnpakSiamida/modules/swot/application/DeleteSwot"
	GetAllSwots "UnpakSiamida/modules/swot/application/GetAllSwots"
	GetSwot "UnpakSiamida/modules/swot/application/GetSwot"
	UpdateSwot "UnpakSiamida/modules/swot/application/UpdateSwot"
	domainSwot "UnpakSiamida/modules/swot/domain"
)

func parseOptionalFloat(val string) *float64 {
	val = strings.TrimSpace(val)
	if val == "" {
		return nil
	}
	if f, err := strconv.ParseFloat(val, 64); err == nil {
		return &f
	}
	return nil
}

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
// POST /swot
// =======================================================
func CreateSwotHandlerfunc(c *fiber.Ctx) error {
	level, userFU := getUserInfo(c)

	var req struct {
		FakultasUuid   string   `json:"fakultas_unit"`
		IdFakultasUnit *int     `json:"id_fakultas_unit"`
		TahunUuid      string   `json:"tahun_uuid"`
		IdMasterTahun  *int     `json:"id_master_tahun"`
		Tahun          string   `json:"tahun"`
		Kategori       string   `json:"kategori"`
		Faktor         string   `json:"faktor"`
		Deskripsi      string   `json:"deskripsi"`
		Strategi       string   `json:"strategi"`
		Bobot          *float64 `json:"bobot"`
		Rating         *int     `json:"rating"`
		Skor           *float64 `json:"skor"`
	}

	_ = c.BodyParser(&req)

	if req.Kategori == "" {
		req.Kategori = c.FormValue("kategori")
	}
	if req.Faktor == "" {
		req.Faktor = c.FormValue("faktor")
	}
	if req.Deskripsi == "" {
		req.Deskripsi = c.FormValue("deskripsi")
	}
	if req.Strategi == "" {
		req.Strategi = c.FormValue("strategi")
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
	if req.Bobot == nil && c.FormValue("bobot") != "" {
		req.Bobot = parseOptionalFloat(c.FormValue("bobot"))
	}
	if req.Rating == nil && c.FormValue("rating") != "" {
		req.Rating = parseOptionalInt(c.FormValue("rating"))
	}
	if req.Skor == nil && c.FormValue("skor") != "" {
		req.Skor = parseOptionalFloat(c.FormValue("skor"))
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

	cmd := CreateSwot.CreateSwotCommand{
		FakultasUuid:   req.FakultasUuid,
		IdFakultasUnit: req.IdFakultasUnit,
		TahunUuid:      req.TahunUuid,
		IdMasterTahun:  req.IdMasterTahun,
		Tahun:          req.Tahun,
		Kategori:       req.Kategori,
		Faktor:         req.Faktor,
		Deskripsi:      req.Deskripsi,
		Strategi:       req.Strategi,
		Bobot:          req.Bobot,
		Rating:         req.Rating,
		Skor:           req.Skor,
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

	level, userFU := getUserInfo(c)

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
		if !matchUnit(existing.IdFakultasUnit, existing.UuidFakultasUnit, userFU) {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Akses ditolak: bukan unit anda"})
		}
	}

	var req struct {
		FakultasUuid   string   `json:"fakultas_unit"`
		IdFakultasUnit *int     `json:"id_fakultas_unit"`
		TahunUuid      string   `json:"tahun_uuid"`
		IdMasterTahun  *int     `json:"id_master_tahun"`
		Tahun          string   `json:"tahun"`
		Kategori       string   `json:"kategori"`
		Faktor         string   `json:"faktor"`
		Deskripsi      string   `json:"deskripsi"`
		Strategi       string   `json:"strategi"`
		Bobot          *float64 `json:"bobot"`
		Rating         *int     `json:"rating"`
		Skor           *float64 `json:"skor"`
	}

	_ = c.BodyParser(&req)

	if req.Kategori == "" {
		req.Kategori = c.FormValue("kategori")
	}
	if req.Faktor == "" {
		req.Faktor = c.FormValue("faktor")
	}
	if req.Deskripsi == "" {
		req.Deskripsi = c.FormValue("deskripsi")
	}
	if req.Strategi == "" {
		req.Strategi = c.FormValue("strategi")
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
	if req.Bobot == nil && c.FormValue("bobot") != "" {
		req.Bobot = parseOptionalFloat(c.FormValue("bobot"))
	}
	if req.Rating == nil && c.FormValue("rating") != "" {
		req.Rating = parseOptionalInt(c.FormValue("rating"))
	}
	if req.Skor == nil && c.FormValue("skor") != "" {
		req.Skor = parseOptionalFloat(c.FormValue("skor"))
	}

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

	cmd := UpdateSwot.UpdateSwotCommand{
		UUID:           rawUuid,
		FakultasUuid:   req.FakultasUuid,
		IdFakultasUnit: req.IdFakultasUnit,
		TahunUuid:      req.TahunUuid,
		IdMasterTahun:  req.IdMasterTahun,
		Tahun:          req.Tahun,
		Kategori:       req.Kategori,
		Faktor:         req.Faktor,
		Deskripsi:      req.Deskripsi,
		Strategi:       req.Strategi,
		Bobot:          req.Bobot,
		Rating:         req.Rating,
		Skor:           req.Skor,
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

	level, userFU := getUserInfo(c)

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
		if !matchUnit(existing.IdFakultasUnit, existing.UuidFakultasUnit, userFU) {
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

	level, userFU := getUserInfo(c)
	if level != "admin" {
		if !matchUnit(item.IdFakultasUnit, item.UuidFakultasUnit, userFU) {
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

	level, userFU := getUserInfo(c)

	var targetFU string
	if level != "admin" {
		// Auditee: mutlak dibatasi hanya unit miliknya
		if userFU != "" && userFU != "0" {
			targetFU = userFU
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
	roles := []string{"admin", "auditee", "fakultas", "unit", "upml", "upmf", "upms"}
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
