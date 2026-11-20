package presentation

import (
	"context"
	"encoding/json"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/mehdihadeli/go-mediatr"

	commondomain "UnpakSiamida/common/domain"
	"UnpakSiamida/common/helper"
	commoninfra "UnpakSiamida/common/infrastructure"
	commonpresentation "UnpakSiamida/common/presentation"

	CreateUser "UnpakSiamida/modules/user/application/CreateUser"
	DeleteUser "UnpakSiamida/modules/user/application/DeleteUser"
	GetAllUserOptions "UnpakSiamida/modules/user/application/GetAllUserOptions"
	GetAllUsers "UnpakSiamida/modules/user/application/GetAllUsers"
	GetUser "UnpakSiamida/modules/user/application/GetUser"
	MergeUser "UnpakSiamida/modules/user/application/MergeUser"
	SendAccountEmail "UnpakSiamida/modules/user/application/SendAccountEmail"
	SetupUuidUser "UnpakSiamida/modules/user/application/SetupUuidUser"
	UpdateUser "UnpakSiamida/modules/user/application/UpdateUser"
	userdomain "UnpakSiamida/modules/user/domain"

	"gorm.io/gorm"
)

// =======================================================
// POST /user
// =======================================================

// CreateUserHandler godoc
// @Summary Create new User
// @Tags User
// @Param name formData string true "Nama"
// @Param username formData string true "Username"
// @Param password formData string true "Password"
// @Param email formData string true "Email"
// @Param level formData string true "Level User"
// @Param FakultasUnit formData string false "Fakultas / Unit UUID"
// @Produce json
// @Success 200 {object} map[string]string "uuid of created user"
// @Failure 400 {object} commoninfra.ResponseError
// @Failure 404 {object} commoninfra.ResponseError
// @Failure 409 {object} commoninfra.ResponseError
// @Failure 500 {object} commoninfra.ResponseError
// @Router /user [post]
func CreateUserHandler(c *fiber.Ctx) error {
	name := strings.TrimSpace(c.FormValue("name"))
	username := strings.TrimSpace(c.FormValue("username"))
	password := strings.TrimSpace(c.FormValue("password"))
	email := strings.TrimSpace(c.FormValue("email"))
	level := strings.ToLower(strings.TrimSpace(c.FormValue("level")))
	if level != "admin" {
		level = "user"
	}

	fu := strings.TrimSpace(c.FormValue("FakultasUnit"))
	if fu == "" {
		fu = strings.TrimSpace(c.FormValue("uuid_fakultas_unit"))
	}

	cmd := CreateUser.CreateUserCommand{
		Name:             name,
		Username:         username,
		Password:         password,
		Email:            email,
		Level:            level,
		UuidFakultasUnit: helper.StrPtr(fu),
	}

	uuid, err := mediatr.Send[
		CreateUser.CreateUserCommand,
		string,
	](context.Background(), cmd)

	if err != nil {
		return commoninfra.HandleError(c, err)
	}

	return commonpresentation.JsonUUID(c, uuid)
}

// =======================================================
// PUT /user/{uuid}
// =======================================================

// UpdateUserHandler godoc
// @Summary Update existing User
// @Tags User
// @Param uuid path string true "User UUID" format(uuid)
// @Param name formData string true "Nama"
// @Param username formData string true "Username"
// @Param password formData string false "Password (optional)"
// @Param email formData string true "Email"
// @Param level formData string true "Level User"
// @Param FakultasUnit formData string false "Fakultas / Unit UUID"
// @Produce json
// @Success 200 {object} map[string]string "uuid of updated user"
// @Failure 400 {object} commoninfra.ResponseError
// @Failure 404 {object} commoninfra.ResponseError
// @Failure 409 {object} commoninfra.ResponseError
// @Failure 500 {object} commoninfra.ResponseError
// @Router /user/{uuid} [put]
func UpdateUserHandler(c *fiber.Ctx) error {
	name := strings.TrimSpace(c.FormValue("name"))
	username := strings.TrimSpace(c.FormValue("username"))
	password := strings.TrimSpace(c.FormValue("password"))
	email := strings.TrimSpace(c.FormValue("email"))
	level := strings.ToLower(strings.TrimSpace(c.FormValue("level")))
	if level != "admin" {
		level = "user"
	}

	fu := strings.TrimSpace(c.FormValue("FakultasUnit"))
	if fu == "" {
		fu = strings.TrimSpace(c.FormValue("uuid_fakultas_unit"))
	}

	var passwordPtr *string
	if password != "" {
		passwordPtr = &password
	}

	cmd := UpdateUser.UpdateUserCommand{
		Uuid:             c.Params("uuid"),
		Name:             name,
		Username:         username,
		Password:         passwordPtr,
		Email:            email,
		Level:            level,
		UuidFakultasUnit: helper.StrPtr(fu),
	}

	updatedID, err := mediatr.Send[
		UpdateUser.UpdateUserCommand,
		string,
	](context.Background(), cmd)

	if err != nil {
		return commoninfra.HandleError(c, err)
	}

	return c.JSON(fiber.Map{"uuid": updatedID})
}

// =======================================================
// DELETE /user/{uuid}
// =======================================================

// DeleteUserHandler godoc
// @Summary Delete User
// @Tags User
// @Param uuid path string true "User UUID" format(uuid)
// @Produce json
// @Success 200 {object} map[string]string "uuid of deleted user"
// @Failure 400 {object} commoninfra.ResponseError
// @Failure 404 {object} commoninfra.ResponseError
// @Failure 409 {object} commoninfra.ResponseError
// @Failure 500 {object} commoninfra.ResponseError
//
//	@Router /user/{uuid} [delete]
func DeleteUserHandler(c *fiber.Ctx) error {

	cmd := DeleteUser.DeleteUserCommand{
		Uuid: c.Params("uuid"),
	}

	deletedID, err := mediatr.Send[
		DeleteUser.DeleteUserCommand,
		string,
	](context.Background(), cmd)

	if err != nil {
		return commoninfra.HandleError(c, err)
	}

	return c.JSON(fiber.Map{"uuid": deletedID})
}

// =======================================================
// GET /user/{uuid}
// =======================================================

// GetUserHandler godoc
// @Summary Get User by UUID
// @Tags User
// @Param uuid path string true "User UUID" format(uuid)
// @Produce json
// @Success 200 {object} userdomain.User
// @Failure 400 {object} commoninfra.ResponseError
// @Failure 404 {object} commoninfra.ResponseError
// @Failure 409 {object} commoninfra.ResponseError
// @Failure 500 {object} commoninfra.ResponseError
//
//	@Router /user/{uuid} [get]
func GetUserHandler(c *fiber.Ctx) error {

	query := GetUser.GetUserByUuidQuery{
		Uuid: c.Params("uuid"),
	}

	user, err := mediatr.Send[
		GetUser.GetUserByUuidQuery,
		*userdomain.User,
	](context.Background(), query)

	if err != nil {
		return commoninfra.HandleError(c, err)
	}

	if user == nil {
		return c.Status(404).JSON(fiber.Map{"error": "User not found"})
	}

	return c.JSON(user)
}

// =======================================================
// GET /users
// =======================================================

// GetAllUsersHandler godoc
// @Summary Get All Users
// @Tags User
// @Param mode query string false "paging | all | ndjson | sse"
// @Param page query int false "Page number"
// @Param limit query int false "Limit per page"
// @Param search query string false "Search keyword"
// @Param filters query string false "Search filters (field:op:value;...)"
// @Produce json
// @Success 200 {object} commondomain.Paged[userdomain.User]
// @Router /users [get]
func GetAllUsersHandler(c *fiber.Ctx) error {

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

	query := GetAllUsers.GetAllUsersQuery{
		Search:        search,
		SearchFilters: filters,
	}

	var adapter commonpresentation.OutputAdapter[userdomain.User]
	switch mode {
	case "all":
		adapter = &commonpresentation.AllAdapter[userdomain.User]{}
	case "ndjson":
		adapter = &commonpresentation.NDJSONAdapter[userdomain.User]{}
	case "sse":
		adapter = &commonpresentation.SSEAdapter[userdomain.User]{}
	default:
		query.Page = &page
		query.Limit = &limit
		adapter = &commonpresentation.PagingAdapter[userdomain.User]{}
	}

	result, err := mediatr.Send[
		GetAllUsers.GetAllUsersQuery,
		commondomain.Paged[userdomain.User],
	](context.Background(), query)

	if err != nil {
		return commoninfra.HandleError(c, err)
	}

	return adapter.Send(c, result)
}

// =======================================================
// GET /user-options
// =======================================================

// GetAllUserOptionsHandler godoc
// @Summary Get All User Options
// @Tags User
// @Produce json
// @Success 200 {object} userdomain.UserOptions
// @Router /user-options [get]
func GetAllUserOptionsHandler(c *fiber.Ctx) error {
	query := GetAllUserOptions.GetAllUserOptionsQuery{}
	result, err := mediatr.Send[
		GetAllUserOptions.GetAllUserOptionsQuery,
		[]userdomain.UserOptions,
	](c.Context(), query)

	if err != nil {
		return commoninfra.HandleError(c, err)
	}

	return c.JSON(result)
}

func SetupUuidUsersHandlerfunc(c *fiber.Ctx) error {
	cmd := SetupUuidUser.SetupUuidUserCommand{}

	message, err := mediatr.Send[SetupUuidUser.SetupUuidUserCommand, string](context.Background(), cmd)
	if err != nil {
		return commoninfra.HandleError(c, err)
	}

	return c.JSON(fiber.Map{"message": message})
}

// =======================================================
// POST /users/merge
// =======================================================

// MergeUserHandler godoc
// @Summary Merge multiple secondary users into a primary user
// @Tags User
// @Accept multipart/form-data
// @Produce json
// @Success 200 {object} map[string]string "message"
// @Router /users/merge [post]
func MergeUserHandler(c *fiber.Ctx) error {
	primaryUUID := strings.TrimSpace(c.FormValue("primary_uuid"))
	var primaryID *uint
	if val := c.FormValue("primary_id"); val != "" {
		if id, err := strconv.ParseUint(val, 10, 32); err == nil {
			uid := uint(id)
			primaryID = &uid
		}
	}

	var secondaryUUIDs []string
	if val := c.FormValue("secondary_uuids"); val != "" {
		if err := json.Unmarshal([]byte(val), &secondaryUUIDs); err != nil {
			for _, s := range strings.Split(val, ",") {
				if trimmed := strings.TrimSpace(s); trimmed != "" {
					secondaryUUIDs = append(secondaryUUIDs, trimmed)
				}
			}
		}
	}

	var secondaryIDs []uint
	if val := c.FormValue("secondary_ids"); val != "" {
		_ = json.Unmarshal([]byte(val), &secondaryIDs)
	}

	cmd := MergeUser.MergeUserCommand{
		PrimaryUUID:    primaryUUID,
		SecondaryUUIDs: secondaryUUIDs,
		PrimaryID:      primaryID,
		SecondaryIDs:   secondaryIDs,
	}

	message, err := mediatr.Send[MergeUser.MergeUserCommand, string](c.Context(), cmd)
	if err != nil {
		return commoninfra.HandleError(c, err)
	}

	return c.JSON(fiber.Map{
		"message": message,
	})
}

// =======================================================
// POST /user/:uuid/send-email & POST /users/send-email
// =======================================================

// SendAccountEmailHandler godoc
// @Summary Send account credentials email via Outbox Pattern
// @Tags User
// @Param uuid path string false "User UUID"
// @Param user_uuids formData []string false "List of User UUIDs for batch"
// @Produce json
// @Success 200 {object} SendAccountEmail.SendAccountEmailResult
// @Router /user/{uuid}/send-email [post]
// @Router /users/send-email [post]
func SendAccountEmailHandler(c *fiber.Ctx) error {
	var userUUIDs []string

	// Ambil dari path param jika ada (/user/:uuid/send-email)
	if paramUuid := strings.TrimSpace(c.Params("uuid")); paramUuid != "" {
		userUUIDs = append(userUUIDs, paramUuid)
	}

	// Ambil dari JSON body jika ada
	if len(c.Body()) > 0 {
		var body struct {
			UserUUIDs []string `json:"user_uuids"`
			UUIDs     []string `json:"uuids"`
			UUID      string   `json:"uuid"`
		}
		if err := json.Unmarshal(c.Body(), &body); err == nil {
			if len(body.UserUUIDs) > 0 {
				userUUIDs = append(userUUIDs, body.UserUUIDs...)
			}
			if len(body.UUIDs) > 0 {
				userUUIDs = append(userUUIDs, body.UUIDs...)
			}
			if body.UUID != "" {
				userUUIDs = append(userUUIDs, body.UUID)
			}
		}
	}

	// Ambil dari Form value jika ada
	if val := c.FormValue("user_uuids"); val != "" {
		var formArr []string
		if err := json.Unmarshal([]byte(val), &formArr); err == nil {
			userUUIDs = append(userUUIDs, formArr...)
		} else {
			for _, s := range strings.Split(val, ",") {
				if trimmed := strings.TrimSpace(s); trimmed != "" {
					userUUIDs = append(userUUIDs, trimmed)
				}
			}
		}
	}
	if val := c.FormValue("uuid"); val != "" {
		userUUIDs = append(userUUIDs, strings.TrimSpace(val))
	}

	// Deduplikasi UUID
	seen := make(map[string]bool)
	var uniqueUUIDs []string
	for _, u := range userUUIDs {
		uClean := strings.TrimSpace(u)
		if uClean != "" && !seen[uClean] {
			seen[uClean] = true
			uniqueUUIDs = append(uniqueUUIDs, uClean)
		}
	}

	if len(uniqueUUIDs) == 0 {
		return commoninfra.HandleError(c, userdomain.EmptyData())
	}

	cmd := SendAccountEmail.SendAccountEmailCommand{
		UserUUIDs: uniqueUUIDs,
	}

	res, err := mediatr.Send[SendAccountEmail.SendAccountEmailCommand, *SendAccountEmail.SendAccountEmailResult](c.Context(), cmd)
	if err != nil {
		return commoninfra.HandleError(c, err)
	}

	return c.JSON(res)
}

func ModuleUser(app *fiber.App, dbs ...*gorm.DB) {
	admin := []string{"admin"}
	whoamiURL := "http://localhost:3000/whoami"

	app.Get("/user/setupuuid", SetupUuidUsersHandlerfunc)

	app.Post("/user", commonpresentation.JWTMiddleware(), commonpresentation.RBACMiddleware(admin, whoamiURL), CreateUserHandler)
	app.Put("/user/:uuid", commonpresentation.JWTMiddleware(), commonpresentation.RBACMiddleware(admin, whoamiURL), UpdateUserHandler)
	app.Delete("/user/:uuid", commonpresentation.JWTMiddleware(), commonpresentation.RBACMiddleware(admin, whoamiURL), DeleteUserHandler)
	app.Get("/user/:uuid", commonpresentation.SmartCompress(), commonpresentation.JWTMiddleware(), GetUserHandler)
	app.Get("/users", commonpresentation.SmartCompress(), commonpresentation.JWTMiddleware(), GetAllUsersHandler)

	app.Get("/user-options", commonpresentation.SmartCompress(), commonpresentation.JWTMiddleware(), GetAllUserOptionsHandler)

	app.Post("/users/merge", commonpresentation.JWTMiddleware(), commonpresentation.RBACMiddleware(admin, whoamiURL), MergeUserHandler)
	app.Post("/user/merge", commonpresentation.JWTMiddleware(), commonpresentation.RBACMiddleware(admin, whoamiURL), MergeUserHandler)

	app.Post("/user/:uuid/send-email", commonpresentation.JWTMiddleware(), commonpresentation.RBACMiddleware(admin, whoamiURL), SendAccountEmailHandler)
	app.Post("/users/send-email", commonpresentation.JWTMiddleware(), commonpresentation.RBACMiddleware(admin, whoamiURL), SendAccountEmailHandler)

	ModuleUserAbac(app)
}
