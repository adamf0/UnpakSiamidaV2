package CreateProkerV3

import (
	domain "UnpakSiamida/modules/prokerv3/domain"
	"context"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CreateProkerV3CommandHandler struct {
	Repo domain.IProkerV3Repository
	DB   *gorm.DB
}

func parseDateFlexible(val string) *time.Time {
	val = strings.TrimSpace(val)
	if val == "" {
		return nil
	}
	if t, err := time.Parse("2006-01-02", val); err == nil {
		return &t
	}
	if t, err := time.Parse(time.RFC3339, val); err == nil {
		return &t
	}
	return nil
}

func (h *CreateProkerV3CommandHandler) Handle(
	ctx context.Context,
	cmd CreateProkerV3Command,
) (string, error) {
	if err := cmd.Validate(); err != nil {
		return "", err
	}

	// 1. Resolve Fakultas / Unit ID
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

	// 2. Resolve or Create Mata Program
	var mpID *int = cmd.IdMataProgram
	if mpID == nil || *mpID <= 0 {
		namaProg := strings.TrimSpace(cmd.NamaProgram)
		if namaProg != "" {
			var existingMp struct {
				ID int
			}
			err := h.DB.WithContext(ctx).Table("mata_program").
				Select("id").
				Where("mata_program = ?", namaProg).
				Take(&existingMp).Error

			if err == nil && existingMp.ID > 0 {
				mpID = &existingMp.ID
			} else {
				// Find id_master_tahun
				tahunStr := strings.TrimSpace(cmd.Tahun)
				if tahunStr == "" {
					tahunStr = "2026"
				}
				var mt struct {
					ID int
				}
				_ = h.DB.WithContext(ctx).Table("master_tahun").
					Select("id").
					Where("tahun = ?", tahunStr).
					Take(&mt).Error
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
						mpID = &createdMp.ID
					}
				}
			}
		}
	}

	// 3. Resolve Indikator ID
	var indID *int = cmd.IdIndikator
	if indID == nil || *indID <= 0 {
		indKinerja := strings.TrimSpace(cmd.IndikatorKinerja)
		if indKinerja != "" {
			var mir struct {
				ID int
			}
			// Search by UUID or indikator text
			if _, err := uuid.Parse(indKinerja); err == nil {
				_ = h.DB.WithContext(ctx).Table("master_indikator_renstra").Select("id").Where("uuid = ?", indKinerja).Take(&mir).Error
			} else {
				_ = h.DB.WithContext(ctx).Table("master_indikator_renstra").Select("id").Where("indikator = ?", indKinerja).Take(&mir).Error
			}
			if mir.ID > 0 {
				indID = &mir.ID
			}
		}
	}

	// 4. Create matrix_proker record
	matrixProker := &domain.ProkerV3{
		UUID:             uuid.New(),
		IdFakultasUnit:   fuID,
		IdMataProgram:    mpID,
		IdIndikator:      indID,
		TargetUnpak:      strings.TrimSpace(cmd.TargetUnpak),
		CapaianTerdahulu: strings.TrimSpace(cmd.CapaianUnitLalu),
		TargetUnit:       strings.TrimSpace(cmd.TargetUnit),
		Output:           strings.TrimSpace(cmd.Output),
		Outcome:          strings.TrimSpace(cmd.Outcome),
	}

	if err := h.Repo.Create(ctx, matrixProker); err != nil {
		return "", err
	}

	return matrixProker.UUID.String(), nil
}
