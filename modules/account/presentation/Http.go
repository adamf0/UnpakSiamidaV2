package presentation

import (
	"context"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/mehdihadeli/go-mediatr"

	commoninfra "UnpakSiamida/common/infrastructure"
	commonpresentation "UnpakSiamida/common/presentation"
	login "UnpakSiamida/modules/account/application/Login"
	updateprofile "UnpakSiamida/modules/account/application/UpdateProfile"
	who "UnpakSiamida/modules/account/application/Whoami"
	domainaccount "UnpakSiamida/modules/account/domain"
)

// =======================================================
// POST /account
// =======================================================

// LoginHandler godoc
// @Summary Create new Account
// @Tags Account
// @Param username formData string true "Username"
// @Param password formData string true "Password"
// @Produce json
// @Success 200 {object} map[string]string "jwt"
// @Failure 400 {object} commoninfra.ResponseError
// @Failure 404 {object} commoninfra.ResponseError
// @Failure 409 {object} commoninfra.ResponseError
// @Failure 500 {object} commoninfra.ResponseError
// @Router /account [post]
func LoginHandlerfunc(c *fiber.Ctx) error {
	username := c.FormValue("username")
	password := c.FormValue("password")

	cmd := login.LoginCommand{
		Username: username,
		Password: password,
	}

	result, err := mediatr.Send[login.LoginCommand, *domainaccount.LoginResult](context.Background(), cmd)
	if err != nil {
		return commoninfra.HandleError(c, err)
	}

	return c.JSON(fiber.Map{
		// "user_id":       result.UserID,
		"access_token":  result.AccessToken,
		"refresh_token": result.RefreshToken,
	})
}

func WhoAmIHandler(c *fiber.Ctx) error {
	userID, _ := c.Locals("sid").(string)
	if userID == "" {
		userID = c.FormValue("sid")
	}

	cmd := who.WhoamiCommand{
		SID: userID,
	}
	result, err := mediatr.Send[who.WhoamiCommand, *domainaccount.Account](context.Background(), cmd)
	if err != nil {
		return commoninfra.HandleError(c, err)
	}

	return c.JSON(result)
}

func UpdateProfileHandler(c *fiber.Ctx) error {
	userID, _ := c.Locals("sid").(string)
	if strings.TrimSpace(userID) == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Sesi login tidak valid",
		})
	}

	type updateProfileBody struct {
		Username        string `json:"username" form:"username"`
		NidnUsername    string `json:"nidn_username" form:"nidn_username"`
		Name            string `json:"name" form:"name"`
		Email           string `json:"email" form:"email"`
		CurrentPassword string `json:"current_password" form:"current_password"`
		NewPassword     string `json:"new_password" form:"new_password"`
		Password        string `json:"password" form:"password"`
		ConfirmPassword string `json:"confirm_password" form:"confirm_password"`
	}

	var body updateProfileBody
	_ = c.BodyParser(&body)

	username := body.Username
	if strings.TrimSpace(username) == "" {
		username = body.NidnUsername
	}
	if strings.TrimSpace(username) == "" {
		username = c.FormValue("username")
	}
	if strings.TrimSpace(username) == "" {
		username = c.FormValue("nidn_username")
	}

	name := body.Name
	if strings.TrimSpace(name) == "" {
		name = c.FormValue("name")
	}

	email := body.Email
	if email == "" {
		email = c.FormValue("email")
	}

	currentPassword := body.CurrentPassword
	if currentPassword == "" {
		currentPassword = c.FormValue("current_password")
	}

	newPassword := body.NewPassword
	if strings.TrimSpace(newPassword) == "" {
		newPassword = body.Password
	}
	if strings.TrimSpace(newPassword) == "" {
		newPassword = c.FormValue("new_password")
	}
	if strings.TrimSpace(newPassword) == "" {
		newPassword = c.FormValue("password")
	}

	confirmPassword := body.ConfirmPassword
	if strings.TrimSpace(confirmPassword) == "" {
		confirmPassword = c.FormValue("confirm_password")
	}

	cmd := updateprofile.UpdateProfileCommand{
		SID:             strings.TrimSpace(userID),
		Username:        username,
		Name:            name,
		Email:           email,
		CurrentPassword: currentPassword,
		NewPassword:     newPassword,
		ConfirmPassword: confirmPassword,
	}

	result, err := mediatr.Send[updateprofile.UpdateProfileCommand, *domainaccount.Account](context.Background(), cmd)
	if err != nil {
		return commoninfra.HandleError(c, err)
	}

	return c.JSON(fiber.Map{
		"message": "Profil akun berhasil diperbarui",
		"data":    result,
	})
}

func ModuleAccount(app *fiber.App) {
	app.Post("/login", LoginHandlerfunc)
	app.Get("/whoami", commonpresentation.SmartCompress(), commonpresentation.JWTMiddleware(), WhoAmIHandler)
	app.Put("/account/profile", commonpresentation.JWTMiddleware(), UpdateProfileHandler)
	app.Patch("/account/profile", commonpresentation.JWTMiddleware(), UpdateProfileHandler)
}
