package presentation

import (
	"context"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/mehdihadeli/go-mediatr"

	commoninfra "UnpakSiamida/common/infrastructure"
	commonpresentation "UnpakSiamida/common/presentation"
	commondomain "UnpakSiamida/common/domain"

	userdomain "UnpakSiamida/modules/user/domain"
	CreateUser "UnpakSiamida/modules/user/application/CreateUser"
	UpdateUser "UnpakSiamida/modules/user/application/UpdateUser"
	DeleteUser "UnpakSiamida/modules/user/application/DeleteUser"
	GetUser "UnpakSiamida/modules/user/application/GetUser"
	GetAllUsers "UnpakSiamida/modules/user/application/GetAllUsers"
    SetupUuidUser "UnpakSiamida/modules/user/application/SetupUuidUser"
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
// @Failure 400 {object} commondomain.Error
// @Router /user [post]
func CreateUserHandler(c *fiber.Ctx) error {

	fakultasUnit := c.FormValue("FakultasUnit")
	var fakultasUnitPtr *string
	if fakultasUnit != "" {
		fakultasUnitPtr = &fakultasUnit
	}

	cmd := CreateUser.CreateUserCommand{
		Name:         c.FormValue("name"),
		Username:     c.FormValue("username"),
		Password:     c.FormValue("password"),
		Email:        c.FormValue("email"),
		Level:        c.FormValue("level"),
		FakultasUnit: fakultasUnitPtr,
	}

	uuid, err := mediatr.Send[
		CreateUser.CreateUserCommand,
		string,
	](context.Background(), cmd)

	if err != nil {
		return commoninfra.HandleError(c, err)
	}

	return c.JSON(fiber.Map{"uuid": uuid})
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
// @Failure 400 {object} commondomain.Error
// @Router /user/{uuid} [put]
func UpdateUserHandler(c *fiber.Ctx) error {

	password := c.FormValue("password")
	var passwordPtr *string
	if password != "" {
		passwordPtr = &password
	}

	fakultasUnit := c.FormValue("FakultasUnit")
	var fakultasUnitPtr *string
	if fakultasUnit != "" {
		fakultasUnitPtr = &fakultasUnit
	}

	cmd := UpdateUser.UpdateUserCommand{
		Uuid:         c.Params("uuid"),
		Name:         c.FormValue("name"),
		Username:     c.FormValue("username"),
		Password:     passwordPtr,
		Email:        c.FormValue("email"),
		Level:        c.FormValue("level"),
		FakultasUnit: fakultasUnitPtr,
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
// @Failure 404 {object} commondomain.Error
// @Router /user/{uuid} [delete]
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
// @Failure 404 {object} commondomain.Error
// @Router /user/{uuid} [get]
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
// @Success 200 {object} userdomain.PagedUsers
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

	var adapter OutputAdapter
	switch mode {
	case "all":
		adapter = &AllAdapter{}
	case "ndjson":
		adapter = &NDJSONAdapter{}
	case "sse":
		adapter = &SSEAdapter{}
	default:
		query.Page = &page
		query.Limit = &limit
		adapter = &PagingAdapter{}
	}

	result, err := mediatr.Send[
		GetAllUsers.GetAllUsersQuery,
		userdomain.PagedUsers,
	](context.Background(), query)

	if err != nil {
		return commoninfra.HandleError(c, err)
	}

	return adapter.Send(c, result)
}

func SetupUuidUsersHandlerfunc(c *fiber.Ctx) error {
    cmd := SetupUuidUser.SetupUuidUserCommand{}

    message, err := mediatr.Send[SetupUuidUser.SetupUuidUserCommand, string](context.Background(), cmd)
    if err != nil {
        return commoninfra.HandleError(c, err)
    }

    return c.JSON(fiber.Map{"message": message})
}

func ModuleUser(app *fiber.App) {
	admin := []string{"admin"}
	whoamiURL := "http://localhost:3000/whoami"

    app.Get("/user/setupuuid", commonpresentation.JWTMiddleware(), commonpresentation.RBACMiddleware(admin, whoamiURL), SetupUuidUsersHandlerfunc)

	app.Post("/user", commonpresentation.JWTMiddleware(), commonpresentation.RBACMiddleware(admin, whoamiURL), CreateUserHandler)
	app.Put("/user/:uuid", commonpresentation.JWTMiddleware(), commonpresentation.RBACMiddleware(admin, whoamiURL), UpdateUserHandler)
	app.Delete("/user/:uuid", commonpresentation.JWTMiddleware(), commonpresentation.RBACMiddleware(admin, whoamiURL), DeleteUserHandler)
	app.Get("/user/:uuid", commonpresentation.JWTMiddleware(), GetUserHandler)
	app.Get("/users", commonpresentation.JWTMiddleware(), GetAllUsersHandler)
}
