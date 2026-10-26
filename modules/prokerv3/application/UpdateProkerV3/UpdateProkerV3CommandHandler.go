package UpdateProkerV3

import (
	domain "UnpakSiamida/modules/prokerv3/domain"
	"context"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UpdateProkerV3CommandHandler struct {
	Repo domain.IProkerV3Repository
	DB   *gorm.DB
}

func (h *UpdateProkerV3CommandHandler) Handle(
	ctx context.Context,
	cmd UpdateProkerV3Command,
) (string, error) {
	if err := cmd.Validate(); err != nil {
		return "", err
	}

	prokerUUID, err := uuid.Parse(strings.TrimSpace(cmd.UUID))
	if err != nil {
		return "", domain.ProkerNotFound()
	}

	existing, err := h.Repo.GetByUuid(ctx, prokerUUID)
	if err != nil {
		return "", err
	}
	if existing == nil {
		return "", domain.ProkerNotFound()
	}

	// 1. Resolve Fakultas Unit
	var fuID *int = cmd.IdFakultasUnit
	if fuID == nil && strings.TrimSpace(cmd.FakultasUuid) != "" {
		trimmed := strings.TrimSpace(cmd.FakultasUuid)
		if parsed, err := uuid.Parse(trimmed); err == nil {
			var fu struct {
				ID int
			}
			if err := h.DB.WithContext(ctx).Table("v_fakultas_unit").Select("id").Where("uuid = ?", parsed).Take(&fu).Error; err == nil {
				fuID = &fu.ID
			}
		} else if id, err := strconv.Atoi(trimmed); err == nil && id > 0 {
			fuID = &id
		}
	}
	if fuID != nil {
		existing.IdFakultasUnit = fuID
	}

	// 2. Resolve Mata Program & Update NamaProgram
	namaProg := strings.TrimSpace(cmd.NamaProgram)

	// If explicit id_mata_program provided and different from existing, switch to it
	if cmd.IdMataProgram != nil && *cmd.IdMataProgram > 0 {
		if existing.IdMataProgram == nil || *existing.IdMataProgram != *cmd.IdMataProgram {
			existing.IdMataProgram = cmd.IdMataProgram
		}
	}

	if existing.IdMataProgram != nil && *existing.IdMataProgram > 0 {
		var currentMp struct {
			ID          int
			MataProgram string
		}
		errMp := h.DB.WithContext(ctx).Table("mata_program").
			Select("id, mata_program").
			Where("id = ?", *existing.IdMataProgram).
			Take(&currentMp).Error

		if errMp == nil && currentMp.ID > 0 {
			if namaProg != "" && currentMp.MataProgram != namaProg {
				// Name changed. Check if other prokers share this mata_program
				var otherUsageCount int64
				h.DB.WithContext(ctx).Table("matrix_proker").
					Where("id_mata_program = ? AND id != ?", currentMp.ID, existing.ID).
					Count(&otherUsageCount)

				if otherUsageCount == 0 {
					// Only this proker uses it, update in-place
					updateFields := map[string]interface{}{
						"mata_program": namaProg,
					}
					if fuID != nil && *fuID > 0 {
						updateFields["id_fakultas_unit"] = *fuID
					}
					_ = h.DB.WithContext(ctx).Table("mata_program").
						Where("id = ?", currentMp.ID).
						Updates(updateFields).Error
				} else {
					// Other prokers share this record, find or create another mata_program with new name
					tahunStr := strings.TrimSpace(cmd.Tahun)
					if tahunStr == "" {
						tahunStr = "2026"
					}
					var mt struct {
						ID int
					}
					_ = h.DB.WithContext(ctx).Table("master_tahun").Select("id").Where("tahun = ?", tahunStr).Take(&mt).Error
					if mt.ID == 0 {
						mt.ID = 1
					}

					var existingWithNewName struct {
						ID int
					}
					errFind := h.DB.WithContext(ctx).Table("mata_program").
						Select("id").
						Where("mata_program = ? AND id_master_tahun = ?", namaProg, mt.ID).
						Take(&existingWithNewName).Error

					if errFind == nil && existingWithNewName.ID > 0 {
						existing.IdMataProgram = &existingWithNewName.ID
					} else {

						newMpUUID := uuid.New()
						_ = h.DB.WithContext(ctx).Exec(
							"INSERT INTO mata_program (uuid, id_master_tahun, id_fakultas_unit, mata_program) VALUES (?, ?, ?, ?)",
							newMpUUID, mt.ID, fuID, namaProg,
						).Error
						var createdMp struct {
							ID int
						}
						_ = h.DB.WithContext(ctx).Table("mata_program").Select("id").Where("uuid = ?", newMpUUID).Take(&createdMp).Error
						if createdMp.ID > 0 {
							existing.IdMataProgram = &createdMp.ID
						}
					}
				}
			} else if fuID != nil && *fuID > 0 {
				// Name identical, ensure id_fakultas_unit is synced if empty
				_ = h.DB.WithContext(ctx).Table("mata_program").
					Where("id = ? AND (id_fakultas_unit IS NULL OR id_fakultas_unit = 0)", currentMp.ID).
					Update("id_fakultas_unit", *fuID).Error
			}
		}
	} else if namaProg != "" {
		// No IdMataProgram on existing proker, find by name or create
		var existingMp struct {
			ID int
		}
		err := h.DB.WithContext(ctx).Table("mata_program").Select("id").Where("mata_program = ?", namaProg).Take(&existingMp).Error
		if err == nil && existingMp.ID > 0 {
			existing.IdMataProgram = &existingMp.ID
		} else {
			tahunStr := strings.TrimSpace(cmd.Tahun)
			if tahunStr == "" {
				tahunStr = "2026"
			}
			var mt struct {
				ID int
			}
			_ = h.DB.WithContext(ctx).Table("master_tahun").Select("id").Where("tahun = ?", tahunStr).Take(&mt).Error
			if mt.ID == 0 {
				mt.ID = 1
			}

			newMpUUID := uuid.New()
			_ = h.DB.WithContext(ctx).Exec(
				"INSERT INTO mata_program (uuid, id_master_tahun, id_fakultas_unit, mata_program) VALUES (?, ?, ?, ?)",
				newMpUUID, mt.ID, fuID, namaProg,
			).Error
			var createdMp struct {
				ID int
			}
			_ = h.DB.WithContext(ctx).Table("mata_program").Select("id").Where("uuid = ?", newMpUUID).Take(&createdMp).Error
			if createdMp.ID > 0 {
				existing.IdMataProgram = &createdMp.ID
			}
		}
	}

	// 3. Resolve Indikator ID
	if cmd.IdIndikator != nil && *cmd.IdIndikator > 0 {
		existing.IdIndikator = cmd.IdIndikator
	} else if strings.TrimSpace(cmd.IndikatorKinerja) != "" {
		indKinerja := strings.TrimSpace(cmd.IndikatorKinerja)
		var mir struct {
			ID int
		}
		if _, err := uuid.Parse(indKinerja); err == nil {
			_ = h.DB.WithContext(ctx).Table("master_indikator_renstra").Select("id").Where("uuid = ?", indKinerja).Take(&mir).Error
		} else {
			_ = h.DB.WithContext(ctx).Table("master_indikator_renstra").Select("id").Where("indikator = ?", indKinerja).Take(&mir).Error
		}
		if mir.ID > 0 {
			existing.IdIndikator = &mir.ID
		}
	}

	if cmd.TargetUnpak != "" {
		existing.TargetUnpak = strings.TrimSpace(cmd.TargetUnpak)
	}
	if cmd.CapaianUnitLalu != "" {
		existing.CapaianTerdahulu = strings.TrimSpace(cmd.CapaianUnitLalu)
	}
	if cmd.TargetUnit != "" {
		existing.TargetUnit = strings.TrimSpace(cmd.TargetUnit)
	}
	if cmd.Output != "" {
		existing.Output = strings.TrimSpace(cmd.Output)
	}
	if cmd.Outcome != "" {
		existing.Outcome = strings.TrimSpace(cmd.Outcome)
	}

	// 4. Update PIC in matrix_proker & sync to aktivitasv2
	existing.PIC = strings.TrimSpace(cmd.PIC)
	_ = h.DB.WithContext(ctx).Table("aktivitasv2").
		Where("id_matrik_proker = ?", existing.ID).
		Update("PIC", strings.TrimSpace(cmd.PIC)).Error

	if err := h.Repo.Update(ctx, existing); err != nil {
		return "", err
	}

	return existing.UUID.String(), nil
}
