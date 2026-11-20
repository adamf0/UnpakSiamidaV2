package presentation

import (
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/mehdihadeli/go-mediatr"

	commoninfra "UnpakSiamida/common/infrastructure"
	commonpresentation "UnpakSiamida/common/presentation"
	createCmd "UnpakSiamida/modules/mataprogramv2/application/CreateMataProgramV2"
	deleteCmd "UnpakSiamida/modules/mataprogramv2/application/DeleteMataProgramV2"
	getAllQuery "UnpakSiamida/modules/mataprogramv2/application/GetAllMataProgramV2"
	getByUuidQuery "UnpakSiamida/modules/mataprogramv2/application/GetMataProgramV2"
	restoreCmd "UnpakSiamida/modules/mataprogramv2/application/RestoreMataProgramV2"
	updateCmd "UnpakSiamida/modules/mataprogramv2/application/UpdateMataProgramV2"
	domain "UnpakSiamida/modules/mataprogramv2/domain"
)

func firstNonEmpty(values ...string) string {
	for _, v := range values {
		if trimmed := strings.TrimSpace(v); trimmed != "" {
			return trimmed
		}
	}
	return ""
}

func extractMultipartTargets(c *fiber.Ctx) []string {
	var targets []string
	if form, err := c.MultipartForm(); err == nil && form != nil {
		for _, key := range []string{"targets[]", "targets", "fakultas_units[]", "fakultas_units"} {
			if vals, ok := form.Value[key]; ok {
				targets = append(targets, vals...)
			}
		}
	}
	return targets
}

func extractTemplateRenstraIDs(c *fiber.Ctx) []uint {
	var ids []uint
	seen := make(map[uint]bool)

	addID := func(raw string) {
		trimmed := strings.TrimSpace(raw)
		if trimmed == "" {
			return
		}
		for _, part := range strings.Split(trimmed, ",") {
			val := strings.TrimSpace(part)
			if n, err := strconv.Atoi(val); err == nil && n > 0 {
				u := uint(n)
				if !seen[u] {
					seen[u] = true
					ids = append(ids, u)
				}
			}
		}
	}

	// 1. Check multipart form
	if form, err := c.MultipartForm(); err == nil && form != nil {
		for _, key := range []string{
			"template_renstras[]", "template_renstras",
			"template_renstra_ids[]", "template_renstra_ids",
			"templates[]", "templates",
		} {
			if vals, ok := form.Value[key]; ok {
				for _, v := range vals {
					addID(v)
				}
			}
		}
	}

	// 2. Check form values
	for _, key := range []string{"template_renstras", "template_renstra_ids", "template_renstra_id"} {
		if v := c.FormValue(key); v != "" {
			addID(v)
		}
	}

	// 3. Check JSON payload if applicable
	type jsonPayload struct {
		TemplateRenstraIDs []uint   `json:"template_renstra_ids"`
		TemplateRenstras   []uint   `json:"template_renstras"`
		TemplateIDsStr     []string `json:"template_ids"`
	}
	var payload jsonPayload
	if err := c.BodyParser(&payload); err == nil {
		for _, id := range payload.TemplateRenstraIDs {
			if id > 0 && !seen[id] {
				seen[id] = true
				ids = append(ids, id)
			}
		}
		for _, id := range payload.TemplateRenstras {
			if id > 0 && !seen[id] {
				seen[id] = true
				ids = append(ids, id)
			}
		}
		for _, str := range payload.TemplateIDsStr {
			addID(str)
		}
	}

	return ids
}

func GetAllMataProgramV2Handler(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page"))
	limit, _ := strconv.Atoi(c.Query("limit"))

	res, err := mediatr.Send[getAllQuery.GetAllMataProgramV2Query, *getAllQuery.GetAllMataProgramV2Result](
		c.Context(),
		getAllQuery.GetAllMataProgramV2Query{
			Tahun:       firstNonEmpty(c.Query("tahun"), c.Query("filter_tahun")),
			Unit:        firstNonEmpty(c.Query("id_fakultas_unit"), c.Query("target"), c.Query("fakultas_unit")),
			Search:      firstNonEmpty(c.Query("search"), c.Query("q")),
			Mode:        strings.TrimSpace(c.Query("mode")),
			Ungrouped:   c.Query("ungrouped") == "true",
			ShowDeleted: c.Query("show_deleted") == "true" || c.Query("deleted") == "true" || c.Query("is_deleted") == "1",
			FiltersRaw:  strings.TrimSpace(c.Query("filters")),
			Page:        page,
			Limit:       limit,
		},
	)
	if err != nil {
		return commoninfra.HandleError(c, err)
	}

	if res.IsFlat {
		return c.JSON(res.Flat)
	}
	return c.JSON(res.Paged)
}

func GetMataProgramV2ByUuidHandler(c *fiber.Ctx) error {
	res, err := mediatr.Send[getByUuidQuery.GetMataProgramV2ByUuidQuery, *domain.MataProgramV2Detail](
		c.Context(),
		getByUuidQuery.GetMataProgramV2ByUuidQuery{
			UUID: strings.TrimSpace(c.Params("uuid")),
		},
	)
	if err != nil {
		return commoninfra.HandleError(c, err)
	}
	return c.JSON(res)
}

func CreateMataProgramV2Handler(c *fiber.Ctx) error {
	mataProgram := firstNonEmpty(
		c.FormValue("mata_program"),
		c.FormValue("mataprogram"),
		c.FormValue("nama_program"),
	)
	tahunStr := firstNonEmpty(c.FormValue("tahun"), c.FormValue("id_master_tahun"))

	templateIDs := extractTemplateRenstraIDs(c)
	targets := extractMultipartTargets(c)
	if single := firstNonEmpty(c.FormValue("id_fakultas_unit"), c.FormValue("fakultas_unit"), c.FormValue("target")); single != "" {
		targets = append(targets, single)
	}

	res, err := mediatr.Send[createCmd.CreateMataProgramV2Command, *createCmd.CreateMataProgramV2Result](
		c.Context(),
		createCmd.CreateMataProgramV2Command{
			MataProgram:        mataProgram,
			Tahun:              tahunStr,
			Targets:            targets,
			TemplateRenstraIDs: templateIDs,
		},
	)
	if err != nil {
		return commoninfra.HandleError(c, err)
	}
	return c.JSON(res)
}

func UpdateMataProgramV2Handler(c *fiber.Ctx) error {
	templateIDs := extractTemplateRenstraIDs(c)
	res, err := mediatr.Send[updateCmd.UpdateMataProgramV2Command, *updateCmd.UpdateMataProgramV2Result](
		c.Context(),
		updateCmd.UpdateMataProgramV2Command{
			UUID:               strings.TrimSpace(c.Params("uuid")),
			Scope:              strings.TrimSpace(c.Query("scope")),
			MataProgram:        firstNonEmpty(c.FormValue("mata_program"), c.FormValue("mataprogram"), c.FormValue("nama_program")),
			Tahun:              strings.TrimSpace(c.FormValue("tahun")),
			Targets:            extractMultipartTargets(c),
			TemplateRenstraIDs: templateIDs,
			IdFakultasUnit:     strings.TrimSpace(c.FormValue("id_fakultas_unit")),
		},
	)
	if err != nil {
		return commoninfra.HandleError(c, err)
	}
	return c.JSON(res)
}

func DeleteMataProgramV2Handler(c *fiber.Ctx) error {
	res, err := mediatr.Send[deleteCmd.DeleteMataProgramV2Command, *deleteCmd.DeleteMataProgramV2Result](
		c.Context(),
		deleteCmd.DeleteMataProgramV2Command{
			UUID:  strings.TrimSpace(c.Params("uuid")),
			Scope: strings.TrimSpace(c.Query("scope")),
		},
	)
	if err != nil {
		return commoninfra.HandleError(c, err)
	}
	return c.JSON(res)
}

func RestoreMataProgramV2Handler(c *fiber.Ctx) error {
	res, err := mediatr.Send[restoreCmd.RestoreMataProgramV2Command, *restoreCmd.RestoreMataProgramV2Result](
		c.Context(),
		restoreCmd.RestoreMataProgramV2Command{
			UUID:  strings.TrimSpace(c.Params("uuid")),
			Scope: strings.TrimSpace(c.Query("scope")),
		},
	)
	if err != nil {
		return commoninfra.HandleError(c, err)
	}
	return c.JSON(res)
}

func ModuleMataProgramV2(app *fiber.App) {
	roles := []string{"admin", "auditee", "fakultas", "unit", "user", "upml", "upmf", "upms", "anggaran", "keuangan", "evaluator"}
	whoamiURL := "http://localhost:3000/whoami"

	// GET ALL / SEARCH
	app.Get("/mataprogramv2s", commonpresentation.SmartCompress(), commonpresentation.JWTMiddleware(), GetAllMataProgramV2Handler)
	app.Get("/mataprogramv2", commonpresentation.SmartCompress(), commonpresentation.JWTMiddleware(), GetAllMataProgramV2Handler)

	// GET SINGLE BY UUID
	app.Get("/mataprogramv2/:uuid", commonpresentation.SmartCompress(), commonpresentation.JWTMiddleware(), GetMataProgramV2ByUuidHandler)

	// CREATE / BATCH CREATE
	app.Post("/mataprogramv2", commonpresentation.JWTMiddleware(), commonpresentation.RBACMiddleware(roles, whoamiURL), CreateMataProgramV2Handler)

	// UPDATE BY UUID
	app.Put("/mataprogramv2/:uuid", commonpresentation.JWTMiddleware(), commonpresentation.RBACMiddleware(roles, whoamiURL), UpdateMataProgramV2Handler)

	// DELETE BY UUID
	app.Delete("/mataprogramv2/:uuid", commonpresentation.JWTMiddleware(), commonpresentation.RBACMiddleware(roles, whoamiURL), DeleteMataProgramV2Handler)

	// RESTORE BY UUID
	app.Post("/mataprogramv2/:uuid/restore", commonpresentation.JWTMiddleware(), commonpresentation.RBACMiddleware(roles, whoamiURL), RestoreMataProgramV2Handler)
	app.Put("/mataprogramv2/:uuid/restore", commonpresentation.JWTMiddleware(), commonpresentation.RBACMiddleware(roles, whoamiURL), RestoreMataProgramV2Handler)
}
