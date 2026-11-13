package presentation

import (
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"

	commoninfra "UnpakSiamida/common/infrastructure"
	commonpresentation "UnpakSiamida/common/presentation"
	domainuser "UnpakSiamida/modules/user/domain"
)

var userRoleRepo domainuser.IUserRepository

func SetUserRoleRepo(repo domainuser.IUserRepository) {
	userRoleRepo = repo
}

// GetAllUserRoleHandler godoc
// @Summary Get list of user roles (userv3_role)
// @Tags User
// @Produce json
// @Router /user-role [get]
func GetAllUserRoleHandler(c *fiber.Ctx) error {
	tahun := strings.TrimSpace(c.Query("tahun"))
	status := strings.TrimSpace(c.Query("status"))
	level := strings.ToLower(strings.TrimSpace(c.Query("level")))
	search := strings.TrimSpace(c.Query("search"))
	userUuid := strings.TrimSpace(c.Query("user_uuid"))
	fuUuid := strings.TrimSpace(c.Query("uuid_fakultas_unit"))

	if userRoleRepo == nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "User repository not initialized"})
	}

	results, err := userRoleRepo.GetAllUserRole(c.Context(), tahun, status, level, search, userUuid, fuUuid)
	if err != nil {
		return commoninfra.HandleError(c, err)
	}

	return c.JSON(results)
}

// CreateUserRoleHandler godoc
// @Summary Create user role in userv3_role
// @Tags User
// @Accept multipart/form-data
// @Produce json
// @Router /user-role [post]
func CreateUserRoleHandler(c *fiber.Ctx) error {
	tahun := strings.TrimSpace(c.FormValue("tahun"))
	if tahun == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Tahun role wajib diisi"})
	}

	level := strings.ToLower(strings.TrimSpace(c.FormValue("level")))
	if level == "" {
		level = strings.ToLower(strings.TrimSpace(c.FormValue("role")))
	}
	if level == "users" {
		level = "user"
	}
	if level == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Level / Hak Akses wajib diisi"})
	}

	status := strings.ToLower(strings.TrimSpace(c.FormValue("status")))
	if status == "" {
		status = "active"
	}

	var userID uint
	if val := c.FormValue("id_user"); val != "" {
		if id, err := strconv.ParseUint(val, 10, 32); err == nil {
			userID = uint(id)
		}
	} else if val := strings.TrimSpace(c.FormValue("user_uuid")); val != "" {
		parsedUuid, err := uuid.Parse(val)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "UUID User tidak valid"})
		}
		u, err := userRoleRepo.GetByUuid(c.Context(), parsedUuid)
		if err != nil || u == nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "User tidak ditemukan"})
		}
		userID = u.ID
	}

	if userID == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "User wajib dipilih"})
	}

	var fuID *uint
	fuUuid := strings.TrimSpace(c.FormValue("uuid_fakultas_unit"))
	if fuUuid == "" {
		fuUuid = strings.TrimSpace(c.FormValue("fakultas_unit"))
	}
	if val := c.FormValue("id_fakultas_unit"); val != "" {
		if id, err := strconv.ParseUint(val, 10, 32); err == nil && id > 0 {
			uid := uint(id)
			fuID = &uid
		}
	}
	if fuID == nil && fuUuid != "" {
		resolvedID, _, _, err := userRoleRepo.GetFakultasUnitInfo(c.Context(), nil, fuUuid)
		if err == nil && resolvedID > 0 {
			fuID = &resolvedID
		}
	}

	newRole := &domainuser.UserRole{
		UUID:           uuid.New(),
		IdUser:         userID,
		IdFakultasUnit: fuID,
		Tahun:          tahun,
		Level:          level,
		Status:         status,
	}

	if err := userRoleRepo.CreateUserRole(c.Context(), newRole); err != nil {
		return commoninfra.HandleError(c, err)
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Role user berhasil ditambahkan",
		"data":    newRole,
	})
}

// UpdateUserRoleHandler godoc
// @Summary Update user role in userv3_role
// @Tags User
// @Accept multipart/form-data
// @Produce json
// @Router /user-role/{uuid} [put]
func UpdateUserRoleHandler(c *fiber.Ctx) error {
	rawUUID := strings.TrimSpace(c.Params("uuid"))
	parsedUUID, err := uuid.Parse(rawUUID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "UUID role tidak valid"})
	}

	existing, err := userRoleRepo.GetUserRoleByUuid(c.Context(), parsedUUID)
	if err != nil || existing == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "Data role tidak ditemukan"})
	}

	updates := make(map[string]interface{})

	if val := strings.TrimSpace(c.FormValue("tahun")); val != "" {
		updates["tahun"] = val
	}

	if val := strings.ToLower(strings.TrimSpace(c.FormValue("level"))); val != "" {
		if val == "users" {
			val = "user"
		}
		updates["level"] = val
	} else if val := strings.ToLower(strings.TrimSpace(c.FormValue("role"))); val != "" {
		if val == "users" {
			val = "user"
		}
		updates["level"] = val
	}

	if val := strings.ToLower(strings.TrimSpace(c.FormValue("status"))); val != "" {
		updates["status"] = val
	}

	fuUuid := strings.TrimSpace(c.FormValue("uuid_fakultas_unit"))
	if fuUuid == "" {
		fuUuid = strings.TrimSpace(c.FormValue("fakultas_unit"))
	}
	if val := c.FormValue("id_fakultas_unit"); val != "" {
		if id, err := strconv.ParseUint(val, 10, 32); err == nil {
			if id > 0 {
				uid := uint(id)
				updates["id_fakultas_unit"] = &uid
			} else {
				updates["id_fakultas_unit"] = nil
			}
		}
	} else if fuUuid != "" {
		resolvedID, _, _, err := userRoleRepo.GetFakultasUnitInfo(c.Context(), nil, fuUuid)
		if err == nil && resolvedID > 0 {
			updates["id_fakultas_unit"] = &resolvedID
		} else {
			updates["id_fakultas_unit"] = nil
		}
	} else if _, exists := c.GetReqHeaders()["X-Clear-Unit"]; exists || c.FormValue("clear_unit") == "1" {
		updates["id_fakultas_unit"] = nil
	}

	if len(updates) > 0 {
		if err := userRoleRepo.UpdateUserRole(c.Context(), existing.ID, updates); err != nil {
			return commoninfra.HandleError(c, err)
		}
	}

	return c.JSON(fiber.Map{
		"message": "Role user berhasil diperbarui",
	})
}

// DeleteUserRoleHandler godoc
// @Summary Delete user role from userv3_role
// @Tags User
// @Produce json
// @Router /user-role/{uuid} [delete]
func DeleteUserRoleHandler(c *fiber.Ctx) error {
	rawUUID := strings.TrimSpace(c.Params("uuid"))
	parsedUUID, err := uuid.Parse(rawUUID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "UUID role tidak valid"})
	}

	if err := userRoleRepo.DeleteUserRole(c.Context(), parsedUUID); err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "Role tidak ditemukan"})
		}
		return commoninfra.HandleError(c, err)
	}

	return c.JSON(fiber.Map{
		"message": "Role user berhasil dihapus",
	})
}

func ModuleUserRole(app *fiber.App, repo domainuser.IUserRepository) {
	SetUserRoleRepo(repo)
	admin := []string{"admin"}
	whoamiURL := "http://localhost:3000/whoami"

	app.Get("/user-role", commonpresentation.SmartCompress(), commonpresentation.JWTMiddleware(), GetAllUserRoleHandler)
	app.Post("/user-role", commonpresentation.JWTMiddleware(), commonpresentation.RBACMiddleware(admin, whoamiURL), CreateUserRoleHandler)
	app.Put("/user-role/:uuid", commonpresentation.JWTMiddleware(), commonpresentation.RBACMiddleware(admin, whoamiURL), UpdateUserRoleHandler)
	app.Delete("/user-role/:uuid", commonpresentation.JWTMiddleware(), commonpresentation.RBACMiddleware(admin, whoamiURL), DeleteUserRoleHandler)
}
