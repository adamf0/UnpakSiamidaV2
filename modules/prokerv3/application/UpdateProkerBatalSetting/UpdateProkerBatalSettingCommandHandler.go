package UpdateProkerBatalSetting

import (
	"context"
	"fmt"
	"strings"
	"time"

	domain "UnpakSiamida/modules/prokerv3/domain"
)

type UpdateProkerBatalSettingCommandHandler struct {
	Repo domain.IProkerV3Repository
}

func (h *UpdateProkerBatalSettingCommandHandler) Handle(ctx context.Context, cmd UpdateProkerBatalSettingCommand) (*UpdateProkerBatalSettingResult, error) {
	if h.Repo == nil {
		return nil, domain.RepoNotInitialized()
	}

	if !cmd.IsAdmin {
		return nil, domain.AdminOnlyProkerBatal()
	}

	currentVal, _ := h.Repo.GetProkerBatalSettingRaw(ctx)
	cfg := domain.ParseProkerBatalConfig(currentVal)
	msg := ""

	action := strings.ToLower(strings.TrimSpace(cmd.Action))
	deleteID := cmd.DeleteUnitID
	if deleteID == nil && (action == "delete_unit" || action == "remove_unit") {
		deleteID = cmd.UnitID
	}

	if deleteID != nil && *deleteID > 0 {
		filtered := make([]domain.UnitBatalSchedule, 0, len(cfg.Units))
		for _, u := range cfg.Units {
			if u.IdFakultasUnit != *deleteID {
				filtered = append(filtered, u)
			}
		}
		cfg.Units = filtered
		msg = "Jadwal khusus fakultas/unit berhasil dihapus (kembali mengikuti jadwal global)"
	} else if action == "reset_units" {
		cfg.Units = []domain.UnitBatalSchedule{}
		msg = "Semua jadwal khusus fakultas/unit berhasil direset"
	} else if action == "reset_all" {
		cfg.GlobalDate = ""
		cfg.Units = []domain.UnitBatalSchedule{}
		msg = "Pengaturan tanggal aktif Proker Batal berhasil direset total"
	} else if len(cmd.TargetUnits) > 0 {
		rawDate := ""
		if cmd.Tanggal != nil {
			rawDate = *cmd.Tanggal
		} else if cmd.ActiveDate != nil {
			rawDate = *cmd.ActiveDate
		} else if cmd.ProkerBatalAktifTanggal != nil {
			rawDate = *cmd.ProkerBatalAktifTanggal
		}
		rawDate = strings.TrimSpace(rawDate)
		if rawDate == "" {
			return nil, domain.TanggalAktifTargetRequired()
		}
		parsed, err := time.ParseInLocation("2006-01-02", rawDate, time.Local)
		if err != nil {
			return nil, domain.DateFormatInvalid()
		}
		normalizedDate := parsed.Format("2006-01-02")
		nowStr := time.Now().Format("2006-01-02 15:04:05")

		namesMap, _ := h.Repo.GetFakultasUnitNamesMap(ctx, cmd.TargetUnits)

		for _, tid := range cmd.TargetUnits {
			if tid <= 0 {
				continue
			}
			foundIdx := -1
			for i := range cfg.Units {
				if cfg.Units[i].IdFakultasUnit == tid {
					foundIdx = i
					break
				}
			}
			unitName := namesMap[tid]
			if foundIdx >= 0 {
				cfg.Units[foundIdx].Tanggal = normalizedDate
				cfg.Units[foundIdx].UpdatedAt = nowStr
				if unitName != "" {
					cfg.Units[foundIdx].NamaFakultasUnit = unitName
				}
			} else {
				if unitName == "" {
					unitName = fmt.Sprintf("Unit #%d", tid)
				}
				cfg.Units = append(cfg.Units, domain.UnitBatalSchedule{
					IdFakultasUnit:   tid,
					NamaFakultasUnit: unitName,
					Tanggal:          normalizedDate,
					UpdatedAt:        nowStr,
				})
			}
		}
		msg = fmt.Sprintf("Jadwal Proker Batal untuk %d fakultas/unit berhasil disimpan", len(cmd.TargetUnits))
	} else if cmd.Units != nil {
		validUnits := make([]domain.UnitBatalSchedule, 0, len(*cmd.Units))
		for _, u := range *cmd.Units {
			if u.IdFakultasUnit <= 0 {
				continue
			}
			parsed, err := time.ParseInLocation("2006-01-02", strings.TrimSpace(u.Tanggal), time.Local)
			if err != nil {
				return nil, fmt.Errorf("Format tanggal tidak valid untuk unit #%d (harus YYYY-MM-DD)", u.IdFakultasUnit)
			}
			u.Tanggal = parsed.Format("2006-01-02")
			validUnits = append(validUnits, u)
		}
		cfg.Units = validUnits
		msg = "Jadwal khusus fakultas/unit berhasil diperbarui"
	} else {
		rawDate := ""
		if cmd.GlobalDate != nil {
			rawDate = *cmd.GlobalDate
		} else if cmd.ProkerBatalAktifTanggal != nil {
			rawDate = *cmd.ProkerBatalAktifTanggal
		} else if cmd.Tanggal != nil {
			rawDate = *cmd.Tanggal
		} else if cmd.ActiveDate != nil {
			rawDate = *cmd.ActiveDate
		}
		rawDate = strings.TrimSpace(rawDate)

		if rawDate == "" || action == "reset_global" {
			cfg.GlobalDate = ""
			msg = "Tanggal aktif Proker Batal (Global) berhasil direset"
		} else {
			parsed, err := time.ParseInLocation("2006-01-02", rawDate, time.Local)
			if err != nil {
				return nil, domain.DateFormatInvalid()
			}
			cfg.GlobalDate = parsed.Format("2006-01-02")
			msg = "Tanggal aktif Proker Batal (Global) berhasil disimpan"
		}
	}

	serialized, err := domain.SerializeProkerBatalConfig(cfg)
	if err != nil {
		return nil, err
	}

	var updatedByPtr *int
	if strings.TrimSpace(cmd.ActorSID) != "" {
		uid, err := h.Repo.GetUserIDByUUID(ctx, cmd.ActorSID)
		if err == nil && uid != nil && *uid > 0 {
			updatedByPtr = uid
		}
	}

	if err := h.Repo.SaveProkerBatalSetting(ctx, serialized, updatedByPtr); err != nil {
		return nil, err
	}

	status, err := h.Repo.GetProkerBatalSettingStatus(ctx, time.Now(), nil)
	if err != nil {
		return nil, err
	}

	return &UpdateProkerBatalSettingResult{
		Message:                 msg,
		SettingKey:              status.SettingKey,
		ProkerBatalAktifTanggal: status.ActiveDate,
		GlobalDate:              status.GlobalDate,
		HasDate:                 status.HasDate,
		IsActive:                status.IsActive,
		Status:                  status.Status,
		StatusLabel:             status.StatusText,
		UpdatedAt:               status.UpdatedAt,
		DetailMessage:           status.Message,
		Units:                   status.Units,
	}, nil
}
