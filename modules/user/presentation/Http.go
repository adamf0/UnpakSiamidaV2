package presentation

import (
    "context"
    "github.com/gofiber/fiber/v2"
    "github.com/mehdihadeli/go-mediatr"
    "strings"

    // "UnpakSiamida/common/domain"
    commoninfra "UnpakSiamida/common/infrastructure"
    commondomain "UnpakSiamida/common/domain"
    userdomain "UnpakSiamida/modules/user/domain"
    CreateUser "UnpakSiamida/modules/user/application/CreateUser"
    UpdateUser "UnpakSiamida/modules/user/application/UpdateUser"
    DeleteUser "UnpakSiamida/modules/user/application/DeleteUser"
    GetUser "UnpakSiamida/modules/user/application/GetUser"
    GetAllUsers "UnpakSiamida/modules/user/application/GetAllUsers"
)

func ModuleUser(app *fiber.App) {

    // ------------------------------------------------------
    // CREATE USER (POST /user)
    // ------------------------------------------------------
    app.Post("/user", func(c *fiber.Ctx) error {

        name := c.FormValue("name")
        username := c.FormValue("username")
        password := c.FormValue("password")
        email := c.FormValue("email")
        level := c.FormValue("level")
        fakultasUnit := c.FormValue("FakultasUnit") // opsional

        var fakultasUnitPtr *string
		if fakultasUnit != "" {
			fakultasUnitPtr = &fakultasUnit
		}

        cmd := CreateUser.CreateUserCommand{
            Name:         name,
            Username:     username,
            Password:     password,
            Email:        email,
            Level:        level,
            FakultasUnit: fakultasUnitPtr,
        }

        uuid, err := mediatr.Send[CreateUser.CreateUserCommand, string](context.Background(), cmd)
        if err != nil {
            return commoninfra.HandleError(c, err)
        }

        return c.JSON(fiber.Map{"uuid": uuid})
    })

    // ------------------------------------------------------
    // UPDATE USER (PUT /user/:uuid)
    // ------------------------------------------------------
    app.Put("/user/:uuid", func(c *fiber.Ctx) error {

        uuid := c.Params("uuid")

        name := c.FormValue("name")
        username := c.FormValue("username")
        password := c.FormValue("password")
        email := c.FormValue("email")
        level := c.FormValue("level")
        fakultasUnit := c.FormValue("FakultasUnit")

        var passwordPtr *string
		if password != "" {
			passwordPtr = &password
		}

        var fakultasUnitPtr *string
		if fakultasUnit != "" {
			fakultasUnitPtr = &fakultasUnit
		}

        cmd := UpdateUser.UpdateUserCommand{
            Uuid:         uuid,
            Name:         name,
            Username:     username,
            Password:     passwordPtr,
            Email:        email,
            Level:        level,
            FakultasUnit: fakultasUnitPtr,
        }

        updatedID, err := mediatr.Send[UpdateUser.UpdateUserCommand, string](context.Background(), cmd)
        if err != nil {
            return commoninfra.HandleError(c, err)
        }

        return c.JSON(fiber.Map{"uuid": updatedID})
    })

    // ------------------------------------------------------ 
    // DELETE USER (DELETE /user/:uuid)
    // ------------------------------------------------------
    app.Delete("/user/:uuid", func(c *fiber.Ctx) error {

        uuid := c.Params("uuid")

        cmd := DeleteUser.DeleteUserCommand{
            Uuid: uuid,
        }

        deletedID, err := mediatr.Send[DeleteUser.DeleteUserCommand, string](context.Background(), cmd)
        if err != nil {
            return commoninfra.HandleError(c, err)
        }

        return c.JSON(fiber.Map{"uuid": deletedID})
    })

    // ------------------------------------------------------ 
    // Get USER
    // ------------------------------------------------------
    app.Get("/user/:uuid", func(c *fiber.Ctx) error {
        uuid := c.Params("uuid")

        query := GetUser.GetUserByUuidQuery{
            Uuid: uuid,
        }

        user, err := mediatr.Send[GetUser.GetUserByUuidQuery, *userdomain.User](context.Background(), query)
        if err != nil {
            return commoninfra.HandleError(c, err)
        }

        if user == nil {
            return c.Status(404).JSON(fiber.Map{"error": "User not found"})
        }

        return c.JSON(user)
    })

    // ------------------------------------------------------ 
    // Get All USER
    // ------------------------------------------------------
    app.Get("/users", func(c *fiber.Ctx) error {
        mode := c.Query("mode", "paging") // default mode = paging
        page := c.QueryInt("page", 1)
        limit := c.QueryInt("limit", 10)
        search := c.Query("search", "")

        // Parse filters
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

                var valuePtr *string
                if rawValue != "" && rawValue != "null" {
                    valuePtr = &rawValue
                }
                filters = append(filters, commondomain.SearchFilter{
                    Field:    field,
                    Operator: op,
                    Value:    valuePtr,
                })
            }
        }

        query := GetAllUsers.GetAllUsersQuery{
            Search:        search,
            SearchFilters: filters,
        }

        // Pilih adapter sesuai mode
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

        // Ambil data
        users, err := mediatr.Send[GetAllUsers.GetAllUsersQuery, userdomain.PagedUsers](context.Background(), query)
        if err != nil {
            return commoninfra.HandleError(c, err)
        }

        return adapter.Send(c, users)
    })
}

