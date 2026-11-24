package presentation

import (
	"fmt"
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

type mataProgramV2RequestPayload struct {
	MataProgram        string `json:"mata_program"`
	NamaProgram        string `json:"nama_program"`
	Mataprogram        string `json:"mataprogram"`
	Tahun              any    `json:"tahun"`
	IdMasterTahun      any    `json:"id_master_tahun"`
	TemplateRenstras   []any  `json:"template_renstras"`
	TemplateRenstraIDs []any  `json:"template_renstra_ids"`
	Templates          []any  `json:"templates"`
}

func parseMataProgramV2Input(c *fiber.Ctx) (mataProgram string, tahun string, templateIDs []uint, hasTemplates bool) {
	seenTemplates := make(map[uint]bool)

	addTemplateID := func(raw any) {
		if raw == nil {
			return
		}
		switch v := raw.(type) {
		case uint:
			if v > 0 && !seenTemplates[v] {
				seenTemplates[v] = true
				templateIDs = append(templateIDs, v)
			}
		case int:
			if v > 0 && !seenTemplates[uint(v)] {
				seenTemplates[uint(v)] = true
				templateIDs = append(templateIDs, uint(v))
			}
		case int64:
			if v > 0 && !seenTemplates[uint(v)] {
				seenTemplates[uint(v)] = true
				templateIDs = append(templateIDs, uint(v))
			}
		case float64:
			if v > 0 && !seenTemplates[uint(v)] {
				seenTemplates[uint(v)] = true
				templateIDs = append(templateIDs, uint(v))
			}
		case string:
			trimmed := strings.TrimSpace(v)
			if trimmed == "" {
				return
			}
			for _, part := range strings.Split(trimmed, ",") {
				val := strings.TrimSpace(part)
				if n, err := strconv.Atoi(val); err == nil && n > 0 {
					u := uint(n)
					if !seenTemplates[u] {
						seenTemplates[u] = true
						templateIDs = append(templateIDs, u)
					}
				}
			}
		}
	}

	// 1. Cek JSON body payload jika ada
	var payload mataProgramV2RequestPayload
	if err := c.BodyParser(&payload); err == nil {
		mataProgram = firstNonEmpty(payload.MataProgram, payload.NamaProgram, payload.Mataprogram)
		if payload.Tahun != nil {
			tahun = strings.TrimSpace(fmt.Sprintf("%v", payload.Tahun))
		}
		if tahun == "" && payload.IdMasterTahun != nil {
			tahun = strings.TrimSpace(fmt.Sprintf("%v", payload.IdMasterTahun))
		}

		if payload.TemplateRenstras != nil || payload.TemplateRenstraIDs != nil || payload.Templates != nil {
			hasTemplates = true
		}
		for _, t := range payload.TemplateRenstras {
			addTemplateID(t)
		}
		for _, t := range payload.TemplateRenstraIDs {
			addTemplateID(t)
		}
		for _, t := range payload.Templates {
			addTemplateID(t)
		}
	}

	// 2. Fallback baca dari form-urlencoded / multipart jika mataProgram atau templateIDs masih kosong
	if mataProgram == "" {
		mataProgram = firstNonEmpty(
			c.FormValue("mata_program"),
			c.FormValue("mataprogram"),
			c.FormValue("nama_program"),
		)
	}
	if tahun == "" {
		tahun = firstNonEmpty(c.FormValue("tahun"), c.FormValue("id_master_tahun"))
	}

	for _, key := range []string{"template_renstras", "template_renstra_ids", "template_renstra_id", "templates"} {
		if v := c.FormValue(key); v != "" {
			hasTemplates = true
			addTemplateID(v)
		}
	}

	// 3. Fallback jika dikirim via multipart/form-data
	if form, err := c.MultipartForm(); err == nil && form != nil {
		for _, key := range []string{
			"template_renstras[]", "template_renstras",
			"template_renstra_ids[]", "template_renstra_ids",
			"templates[]", "templates",
		} {
			if vals, ok := form.Value[key]; ok {
				hasTemplates = true
				for _, v := range vals {
					addTemplateID(v)
				}
			}
		}
	}

	return
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
	mataProgram, tahunStr, templateIDs, _ := parseMataProgramV2Input(c)

	res, err := mediatr.Send[createCmd.CreateMataProgramV2Command, *createCmd.CreateMataProgramV2Result](
		c.Context(),
		createCmd.CreateMataProgramV2Command{
			MataProgram:        mataProgram,
			Tahun:              tahunStr,
			TemplateRenstraIDs: templateIDs,
		},
	)
	if err != nil {
		return commoninfra.HandleError(c, err)
	}
	return c.JSON(res)
}

func UpdateMataProgramV2Handler(c *fiber.Ctx) error {
	mataProgram, tahunStr, templateIDs, hasTemplates := parseMataProgramV2Input(c)

	res, err := mediatr.Send[updateCmd.UpdateMataProgramV2Command, *updateCmd.UpdateMataProgramV2Result](
		c.Context(),
		updateCmd.UpdateMataProgramV2Command{
			UUID:               strings.TrimSpace(c.Params("uuid")),
			MataProgram:        mataProgram,
			Tahun:              tahunStr,
			TemplateRenstraIDs: templateIDs,
			HasTemplates:       hasTemplates,
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
			UUID: strings.TrimSpace(c.Params("uuid")),
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
			UUID: strings.TrimSpace(c.Params("uuid")),
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
