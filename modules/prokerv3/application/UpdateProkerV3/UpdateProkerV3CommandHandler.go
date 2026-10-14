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

	// 2. Resolve Mata Program
	if cmd.IdMataProgram != nil && *cmd.IdMataProgram > 0 {
		existing.IdMataProgram = cmd.IdMataProgram
	} else if strings.TrimSpace(cmd.NamaProgram) != "" {
		namaProg := strings.TrimSpace(cmd.NamaProgram)
		var existingMp struct {
			ID int
		}
		err := h.DB.WithContext(ctx).Table("mata_program").Select("id").Where("mata_program = ?", namaProg).Take(&existingMp).Error
		if err == nil && existingMp.ID > 0 {
			existing.IdMataProgram = &existingMp.ID
		} else {
			// Find id_master_tahun
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
			res := h.DB.WithContext(ctx).Exec(
				"INSERT INTO mata_program (uuid, id_master_tahun, mata_program) VALUES (?, ?, ?)",
				newMpUUID, mt.ID, namaProg,
			)
			if res.Error == nil {
				var createdMp struct {
					ID int
				}
				_ = h.DB.WithContext(ctx).Table("mata_program").Select("id").Where("uuid = ?", newMpUUID).Take(&createdMp).Error
				if createdMp.ID > 0 {
					existing.IdMataProgram = &createdMp.ID
				}
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

	if err := h.Repo.Update(ctx, existing); err != nil {
		return "", err
	}

	return existing.UUID.String(), nil
}
