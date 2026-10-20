package presentation

import (
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/mehdihadeli/go-mediatr"

	commoninfra "UnpakSiamida/common/infrastructure"
	commonpresentation "UnpakSiamida/common/presentation"
	createAbac "UnpakSiamida/modules/user/application/CreateUserAbac"
	deleteAbac "UnpakSiamida/modules/user/application/DeleteUserAbac"
	getAllAbac "UnpakSiamida/modules/user/application/GetAllUserAbac"
	updateAbac "UnpakSiamida/modules/user/application/UpdateUserAbac"
	domainuser "UnpakSiamida/modules/user/domain"
)

// GetAllUserAbacHandler godoc
// @Summary Get list of user ABAC assignments
// @Tags User
// @Produce json
// @Router /user-abac [get]
func GetAllUserAbacHandler(c *fiber.Ctx) error {
	query := getAllAbac.GetAllUserAbacQuery{
		Tahun:    strings.TrimSpace(c.Query("tahun")),
		Status:   strings.TrimSpace(c.Query("status")),
		Role:     strings.ToLower(strings.TrimSpace(c.Query("role"))),
		Search:   strings.TrimSpace(c.Query("search")),
		UserUuid: strings.TrimSpace(c.Query("user_uuid")),
		FuUuid:   strings.TrimSpace(c.Query("uuid_fakultas_unit")),
	}

	results, err := mediatr.Send[getAllAbac.GetAllUserAbacQuery, []domainuser.UserAbac](c.Context(), query)
	if err != nil {
		return commoninfra.HandleError(c, err)
	}

	return c.JSON(results)
}

// CreateUserAbacHandler godoc
// @Summary Create user ABAC assignment
// @Tags User
// @Accept multipart/form-data
// @Produce json
// @Router /user-abac [post]
func CreateUserAbacHandler(c *fiber.Ctx) error {
	cmd := createAbac.CreateUserAbacCommand{}

	// Explicit form-data parsing
	if val := c.FormValue("id_user"); val != "" {
		if id, err := strconv.ParseUint(val, 10, 32); err == nil {
			uid := uint(id)
			cmd.IdUser = &uid
		}
	}
	if val := c.FormValue("user_uuid"); val != "" {
		cmd.UserUuid = strings.TrimSpace(val)
	}
	if val := c.FormValue("tahun"); val != "" {
		cmd.Tahun = strings.TrimSpace(val)
	}
	if val := c.FormValue("id_fakultas_unit"); val != "" {
		if id, err := strconv.ParseUint(val, 10, 32); err == nil {
			uid := uint(id)
			cmd.IdFakultasUnit = &uid
		}
	}
	if val := c.FormValue("uuid_fakultas_unit"); val != "" {
		cmd.UuidFakultasUnit = strings.TrimSpace(val)
	}
	if val := c.FormValue("role"); val != "" {
		cmd.Role = strings.TrimSpace(val)
	}
	if val := c.FormValue("status"); val != "" {
		cmd.Status = strings.TrimSpace(val)
	}

	res, err := mediatr.Send[createAbac.CreateUserAbacCommand, *createAbac.CreateUserAbacResult](c.Context(), cmd)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(res)
}

// UpdateUserAbacHandler godoc
// @Summary Update user ABAC assignment
// @Tags User
// @Accept multipart/form-data
// @Produce json
// @Router /user-abac/{uuid} [put]
func UpdateUserAbacHandler(c *fiber.Ctx) error {
	rawUUID := strings.TrimSpace(c.Params("uuid"))
	if rawUUID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "UUID penugasan wajib disertakan"})
	}

	cmd := updateAbac.UpdateUserAbacCommand{
		UUID: rawUUID,
	}

	// Explicit form-data parsing
	if val := c.FormValue("tahun"); val != "" {
		cmd.Tahun = strings.TrimSpace(val)
	}
	if val := c.FormValue("id_fakultas_unit"); val != "" {
		if id, err := strconv.ParseUint(val, 10, 32); err == nil {
			uid := uint(id)
			cmd.IdFakultasUnit = &uid
		}
	}
	if val := c.FormValue("uuid_fakultas_unit"); val != "" {
		cmd.UuidFakultasUnit = strings.TrimSpace(val)
	}
	if val := c.FormValue("role"); val != "" {
		cmd.Role = strings.TrimSpace(val)
	}
	if val := c.FormValue("status"); val != "" {
		cmd.Status = strings.TrimSpace(val)
	}

	msg, err := mediatr.Send[updateAbac.UpdateUserAbacCommand, string](c.Context(), cmd)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": msg,
	})
}

// DeleteUserAbacHandler godoc
// @Summary Delete user ABAC assignment
// @Tags User
// @Produce json
// @Router /user-abac/{uuid} [delete]
func DeleteUserAbacHandler(c *fiber.Ctx) error {
	rawUUID := strings.TrimSpace(c.Params("uuid"))
	if rawUUID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "UUID penugasan wajib disertakan"})
	}

	cmd := deleteAbac.DeleteUserAbacCommand{
		UUID: rawUUID,
	}

	msg, err := mediatr.Send[deleteAbac.DeleteUserAbacCommand, string](c.Context(), cmd)
	if err != nil {
		return commoninfra.HandleError(c, err)
	}

	return c.JSON(fiber.Map{
		"message": msg,
	})
}

func ModuleUserAbac(app *fiber.App) {
	admin := []string{"admin"}
	whoamiURL := "http://localhost:3000/whoami"

	app.Get("/user-abac", commonpresentation.SmartCompress(), commonpresentation.JWTMiddleware(), GetAllUserAbacHandler)
	app.Post("/user-abac", commonpresentation.JWTMiddleware(), commonpresentation.RBACMiddleware(admin, whoamiURL), CreateUserAbacHandler)
	app.Put("/user-abac/:uuid", commonpresentation.JWTMiddleware(), commonpresentation.RBACMiddleware(admin, whoamiURL), UpdateUserAbacHandler)
	app.Delete("/user-abac/:uuid", commonpresentation.JWTMiddleware(), commonpresentation.RBACMiddleware(admin, whoamiURL), DeleteUserAbacHandler)
}
