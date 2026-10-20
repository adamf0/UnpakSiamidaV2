package presentation

import (
	"strings"
	"time"

	commoninfra "UnpakSiamida/common/infrastructure"
	domainuser "UnpakSiamida/modules/user/domain"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserAbacPayload struct {
	IdUser           *uint  `json:"id_user" form:"id_user"`
	UserUuid         string `json:"user_uuid" form:"user_uuid"`
	Tahun            string `json:"tahun" form:"tahun"`
	IdFakultasUnit   *uint  `json:"id_fakultas_unit" form:"id_fakultas_unit"`
	UuidFakultasUnit string `json:"uuid_fakultas_unit" form:"uuid_fakultas_unit"`
	Role             string `json:"role" form:"role"`
	Status           string `json:"status" form:"status"`
}

// GetAllUserAbacHandler godoc
// @Summary Get list of user ABAC assignments
// @Tags User
// @Produce json
// @Router /user-abac [get]
func GetAllUserAbacHandler(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		tahun := strings.TrimSpace(c.Query("tahun"))
		status := strings.TrimSpace(c.Query("status"))
		role := strings.ToLower(strings.TrimSpace(c.Query("role")))
		search := strings.TrimSpace(c.Query("search"))
		userUuid := strings.TrimSpace(c.Query("user_uuid"))
		fuUuid := strings.TrimSpace(c.Query("uuid_fakultas_unit"))

		query := db.WithContext(c.Context()).Table("user_abac").
			Select(`
				user_abac.id,
				user_abac.uuid,
				user_abac.id_user,
				user_abac.tahun,
				user_abac.id_fakultas_unit,
				user_abac.role,
				user_abac.status,
				user_abac.created_at,
				user_abac.updated_at,
				u.name as user_name,
				u.nidn_username as user_username,
				u.email as user_email,
				u.uuid as user_uuid,
				u.level as user_level,
				fu.nama_fak_prod_unit as nama_fakultas_unit,
				fu.uuid as uuid_fakultas_unit,
				fu.type as type_fakultas_unit,
				fu.fakultas as fakultas_unit_parent,
				CASE
					WHEN user_abac.role IS NOT NULL AND user_abac.role != '' THEN LOWER(user_abac.role)
					WHEN fu.type = 'unit' THEN 'upml'
					WHEN LOWER(fu.nama_fak_prod_unit) LIKE '%pascasarjana%' OR LOWER(fu.nama_fak_prod_unit) LIKE '%vokasi%' OR LOWER(fu.nama_fak_prod_unit) LIKE '%sekolah%' THEN 'upms'
					ELSE 'upmf'
				END AS computed_role
			`).
			Joins("LEFT JOIN users u ON user_abac.id_user = u.id").
			Joins("LEFT JOIN v_fakultas_unit fu ON user_abac.id_fakultas_unit = fu.id")

		if tahun != "" && tahun != "all" {
			query = query.Where("user_abac.tahun = ?", tahun)
		}

		if status != "" && status != "all" {
			query = query.Where("user_abac.status = ?", status)
		}

		if userUuid != "" {
			query = query.Where("u.uuid = ?", userUuid)
		}

		if fuUuid != "" {
			query = query.Where("fu.uuid = ?", fuUuid)
		}

		if search != "" {
			like := "%" + strings.ToLower(search) + "%"
			query = query.Where("LOWER(u.name) LIKE ? OR LOWER(u.nidn_username) LIKE ? OR LOWER(fu.nama_fak_prod_unit) LIKE ?", like, like, like)
		}

		var results []domainuser.UserAbac
		if err := query.Order("user_abac.id DESC").Find(&results).Error; err != nil {
			return commoninfra.HandleError(c, err)
		}

		// Filter role if specified
		if role != "" && role != "all" {
			filtered := make([]domainuser.UserAbac, 0, len(results))
			for _, item := range results {
				r := strings.ToLower(item.ComputedRole)
				if r == "" {
					r = strings.ToLower(item.Role)
				}
				if r == role {
					filtered = append(filtered, item)
				}
			}
			results = filtered
		}

		return c.JSON(results)
	}
}

// CreateUserAbacHandler godoc
// @Summary Create user ABAC assignment
// @Tags User
// @Accept json
// @Produce json
// @Router /user-abac [post]
func CreateUserAbacHandler(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var p UserAbacPayload
		_ = c.BodyParser(&p)

		if p.UserUuid == "" {
			p.UserUuid = c.FormValue("user_uuid")
		}
		if p.Tahun == "" {
			p.Tahun = c.FormValue("tahun")
		}
		if p.UuidFakultasUnit == "" {
			p.UuidFakultasUnit = c.FormValue("uuid_fakultas_unit")
		}
		if p.Role == "" {
			p.Role = c.FormValue("role")
		}
		if p.Status == "" {
			p.Status = c.FormValue("status")
		}

		p.Tahun = strings.TrimSpace(p.Tahun)
		if p.Tahun == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "Tahun penugasan wajib diisi",
			})
		}

		// Resolve User ID
		var userID uint
		if p.IdUser != nil && *p.IdUser > 0 {
			userID = *p.IdUser
		} else if strings.TrimSpace(p.UserUuid) != "" {
			var u domainuser.User
			if err := db.WithContext(c.Context()).Table("users").Select("id").Where("uuid = ?", strings.TrimSpace(p.UserUuid)).First(&u).Error; err != nil {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"message": "User yang dipilih tidak ditemukan dalam sistem",
				})
			}
			userID = u.ID
		} else {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "Pengguna wajib dipilih",
			})
		}

		// Resolve Fakultas Unit ID
		var fuID uint
		var fuType string
		var fuNama string
		if p.IdFakultasUnit != nil && *p.IdFakultasUnit > 0 {
			fuID = *p.IdFakultasUnit
			type FuRow struct {
				Type string `gorm:"column:type"`
				Nama string `gorm:"column:nama_fak_prod_unit"`
			}
			var row FuRow
			_ = db.WithContext(c.Context()).Table("v_fakultas_unit").Select("type, nama_fak_prod_unit").Where("id = ?", fuID).First(&row).Error
			fuType = row.Type
			fuNama = row.Nama
		} else if strings.TrimSpace(p.UuidFakultasUnit) != "" {
			type FuRow struct {
				ID   uint   `gorm:"column:id"`
				Type string `gorm:"column:type"`
				Nama string `gorm:"column:nama_fak_prod_unit"`
			}
			var row FuRow
			if err := db.WithContext(c.Context()).Table("v_fakultas_unit").Select("id, type, nama_fak_prod_unit").Where("uuid = ?", strings.TrimSpace(p.UuidFakultasUnit)).First(&row).Error; err != nil {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"message": "Fakultas/Unit yang dipilih tidak ditemukan",
				})
			}
			fuID = row.ID
			fuType = row.Type
			fuNama = row.Nama
		} else {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "Target Fakultas / Unit / Lembaga wajib dipilih",
			})
		}

		// Rule: "kalau yg saya pilih fakultas maka hanya bisa pilih upml dan upmf, kalau bukan fakultas maka bisa pilih upml upms"
		isFakultas := fuType == "fakultas" || strings.Contains(strings.ToLower(fuNama), "(fakultas)") || strings.HasPrefix(strings.ToLower(fuNama), "fakultas")

		role := strings.ToLower(strings.TrimSpace(p.Role))
		if role == "" || role == "auto" {
			if isFakultas {
				role = "upmf"
			} else {
				lowNama := strings.ToLower(fuNama)
				if strings.Contains(lowNama, "pascasarjana") || strings.Contains(lowNama, "vokasi") || strings.Contains(lowNama, "sekolah") {
					role = "upms"
				} else {
					role = "upml"
				}
			}
		} else {
			if isFakultas && role == "upms" {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"message": "Target Fakultas hanya dapat memilih role UPML atau UPMF",
				})
			}
			if !isFakultas && role == "upmf" {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"message": "Target non-Fakultas hanya dapat memilih role UPML atau UPMS",
				})
			}
		}

		status := strings.ToLower(strings.TrimSpace(p.Status))
		if status == "" {
			status = "active"
		}

		// Check duplicate
		var count int64
		_ = db.WithContext(c.Context()).Table("user_abac").
			Where("id_user = ? AND tahun = ? AND id_fakultas_unit = ?", userID, p.Tahun, fuID).
			Count(&count).Error

		if count > 0 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "Penugasan untuk user, tahun, dan target unit ini sudah ada",
			})
		}

		newUUID := uuid.New()
		now := time.Now()
		abac := domainuser.UserAbac{
			UUID:           newUUID,
			IdUser:         userID,
			Tahun:          p.Tahun,
			IdFakultasUnit: fuID,
			Role:           role,
			Status:         status,
			CreatedAt:      now,
			UpdatedAt:      now,
		}

		if err := db.WithContext(c.Context()).Table("user_abac").Create(&abac).Error; err != nil {
			return commoninfra.HandleError(c, err)
		}

		return c.Status(fiber.StatusCreated).JSON(fiber.Map{
			"message": "Penugasan role berhasil disimpan",
			"uuid":    newUUID.String(),
			"id":      abac.ID,
			"role":    role,
		})
	}
}

// UpdateUserAbacHandler godoc
// @Summary Update user ABAC assignment
// @Tags User
// @Accept json
// @Produce json
// @Router /user-abac/{uuid} [put]
func UpdateUserAbacHandler(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		rawUUID := strings.TrimSpace(c.Params("uuid"))
		if rawUUID == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "UUID penugasan wajib disertakan"})
		}

		var p UserAbacPayload
		_ = c.BodyParser(&p)

		if p.Tahun == "" {
			p.Tahun = c.FormValue("tahun")
		}
		if p.Role == "" {
			p.Role = c.FormValue("role")
		}
		if p.Status == "" {
			p.Status = c.FormValue("status")
		}
		if p.UuidFakultasUnit == "" {
			p.UuidFakultasUnit = c.FormValue("uuid_fakultas_unit")
		}

		var existing domainuser.UserAbac
		if err := db.WithContext(c.Context()).Table("user_abac").Where("uuid = ?", rawUUID).First(&existing).Error; err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "Data penugasan tidak ditemukan"})
		}

		updates := map[string]interface{}{
			"updated_at": time.Now(),
		}

		if p.Status != "" {
			updates["status"] = strings.ToLower(strings.TrimSpace(p.Status))
		}
		if p.Tahun != "" {
			updates["tahun"] = strings.TrimSpace(p.Tahun)
		}
		targetFuID := existing.IdFakultasUnit
		if p.IdFakultasUnit != nil && *p.IdFakultasUnit > 0 {
			targetFuID = *p.IdFakultasUnit
			updates["id_fakultas_unit"] = *p.IdFakultasUnit
		} else if strings.TrimSpace(p.UuidFakultasUnit) != "" {
			type FuRow struct {
				ID uint `gorm:"column:id"`
			}
			var row FuRow
			if err := db.WithContext(c.Context()).Table("v_fakultas_unit").Select("id").Where("uuid = ?", strings.TrimSpace(p.UuidFakultasUnit)).First(&row).Error; err == nil {
				targetFuID = row.ID
				updates["id_fakultas_unit"] = row.ID
			}
		}

		if p.Role != "" {
			reqRole := strings.ToLower(strings.TrimSpace(p.Role))
			if reqRole != "auto" {
				type FuRow struct {
					Type string `gorm:"column:type"`
					Nama string `gorm:"column:nama_fak_prod_unit"`
				}
				var fuRow FuRow
				_ = db.WithContext(c.Context()).Table("v_fakultas_unit").Select("type, nama_fak_prod_unit").Where("id = ?", targetFuID).First(&fuRow).Error

				isFak := fuRow.Type == "fakultas" || strings.Contains(strings.ToLower(fuRow.Nama), "(fakultas)") || strings.HasPrefix(strings.ToLower(fuRow.Nama), "fakultas")
				if isFak && reqRole == "upms" {
					return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
						"message": "Target Fakultas hanya dapat memilih role UPML atau UPMF",
					})
				}
				if !isFak && reqRole == "upmf" {
					return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
						"message": "Target non-Fakultas hanya dapat memilih role UPML atau UPMS",
					})
				}
				updates["role"] = reqRole
			}
		}

		if err := db.WithContext(c.Context()).Table("user_abac").Where("id = ?", existing.ID).Updates(updates).Error; err != nil {
			return commoninfra.HandleError(c, err)
		}

		return c.JSON(fiber.Map{
			"message": "Penugasan role berhasil diperbarui",
		})
	}
}

// DeleteUserAbacHandler godoc
// @Summary Delete user ABAC assignment
// @Tags User
// @Produce json
// @Router /user-abac/{uuid} [delete]
func DeleteUserAbacHandler(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		rawUUID := strings.TrimSpace(c.Params("uuid"))
		if rawUUID == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "UUID penugasan wajib disertakan"})
		}

		result := db.WithContext(c.Context()).Table("user_abac").Where("uuid = ?", rawUUID).Delete(&domainuser.UserAbac{})
		if result.Error != nil {
			return commoninfra.HandleError(c, result.Error)
		}
		if result.RowsAffected == 0 {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "Data penugasan tidak ditemukan"})
		}

		return c.JSON(fiber.Map{
			"message": "Penugasan role berhasil dihapus",
		})
	}
}
