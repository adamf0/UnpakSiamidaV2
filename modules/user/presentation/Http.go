package presentation

import (
	"context"
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
	SetupUuidUser "UnpakSiamida/modules/user/application/SetupUuidUser"
	UpdateUser "UnpakSiamida/modules/user/application/UpdateUser"
	userdomain "UnpakSiamida/modules/user/domain"

	"gorm.io/gorm"
)

// =======================================================
// POST /user
// =======================================================

type UserPayload struct {
	Name             string `json:"name" form:"name"`
	Username         string `json:"username" form:"username"`
	Password         string `json:"password" form:"password"`
	Email            string `json:"email" form:"email"`
	Level            string `json:"level" form:"level"`
	FakultasUnit     string `json:"FakultasUnit" form:"FakultasUnit"`
	UuidFakultasUnit string `json:"uuid_fakultas_unit" form:"uuid_fakultas_unit"`
}

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
	var p UserPayload
	_ = c.BodyParser(&p)

	name := p.Name
	if name == "" {
		name = c.FormValue("name")
	}
	username := p.Username
	if username == "" {
		username = c.FormValue("username")
	}
	password := p.Password
	if password == "" {
		password = c.FormValue("password")
	}
	email := p.Email
	if email == "" {
		email = c.FormValue("email")
	}
	level := p.Level
	if level == "" {
		level = c.FormValue("level")
	}

	fu := p.FakultasUnit
	if fu == "" {
		fu = p.UuidFakultasUnit
	}
	if fu == "" {
		fu = c.FormValue("FakultasUnit")
	}
	if fu == "" {
		fu = c.FormValue("uuid_fakultas_unit")
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
	var p UserPayload
	_ = c.BodyParser(&p)

	name := p.Name
	if name == "" {
		name = c.FormValue("name")
	}
	username := p.Username
	if username == "" {
		username = c.FormValue("username")
	}
	password := p.Password
	if password == "" {
		password = c.FormValue("password")
	}
	email := p.Email
	if email == "" {
		email = c.FormValue("email")
	}
	level := p.Level
	if level == "" {
		level = c.FormValue("level")
	}

	fu := p.FakultasUnit
	if fu == "" {
		fu = p.UuidFakultasUnit
	}
	if fu == "" {
		fu = c.FormValue("FakultasUnit")
	}
	if fu == "" {
		fu = c.FormValue("uuid_fakultas_unit")
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

type MergeUserPayload struct {
	PrimaryUUID    string   `json:"primary_uuid"`
	SecondaryUUIDs []string `json:"secondary_uuids"`
	PrimaryID      *uint    `json:"primary_id"`
	SecondaryIDs   []uint   `json:"secondary_ids"`
}

// MergeUserHandler godoc
// @Summary Merge multiple secondary users into a primary user
// @Tags User
// @Accept json
// @Produce json
// @Success 200 {object} map[string]string "message"
// @Router /users/merge [post]
func MergeUserHandler(c *fiber.Ctx) error {
	var p MergeUserPayload
	if err := c.BodyParser(&p); err != nil {
		return commoninfra.HandleError(c, err)
	}

	cmd := MergeUser.MergeUserCommand{
		PrimaryUUID:    p.PrimaryUUID,
		SecondaryUUIDs: p.SecondaryUUIDs,
		PrimaryID:      p.PrimaryID,
		SecondaryIDs:   p.SecondaryIDs,
	}

	message, err := mediatr.Send[MergeUser.MergeUserCommand, string](c.Context(), cmd)
	if err != nil {
		return commoninfra.HandleError(c, err)
	}

	return c.JSON(fiber.Map{
		"message": message,
	})
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

	if len(dbs) > 0 && dbs[0] != nil {
		db := dbs[0]
		_ = db.Exec("ALTER TABLE `users` MODIFY COLUMN `level` VARCHAR(100) NOT NULL DEFAULT 'user'")
		app.Get("/user-abac", commonpresentation.SmartCompress(), commonpresentation.JWTMiddleware(), GetAllUserAbacHandler(db))
		app.Post("/user-abac", commonpresentation.JWTMiddleware(), commonpresentation.RBACMiddleware(admin, whoamiURL), CreateUserAbacHandler(db))
		app.Put("/user-abac/:uuid", commonpresentation.JWTMiddleware(), commonpresentation.RBACMiddleware(admin, whoamiURL), UpdateUserAbacHandler(db))
		app.Delete("/user-abac/:uuid", commonpresentation.JWTMiddleware(), commonpresentation.RBACMiddleware(admin, whoamiURL), DeleteUserAbacHandler(db))
	}
}

