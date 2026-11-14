package UpdateProkerBatalSetting

import domain "UnpakSiamida/modules/prokerv3/domain"

type UpdateProkerBatalSettingCommand struct {
	ProkerBatalAktifTanggal *string                     `json:"proker_batal_aktif_tanggal" form:"proker_batal_aktif_tanggal"`
	GlobalDate              *string                     `json:"global_date" form:"global_date"`
	Tanggal                 *string                     `json:"tanggal" form:"tanggal"`
	ActiveDate              *string                     `json:"active_date" form:"active_date"`
	Action                  string                      `json:"action" form:"action"`
	TargetUnits             []int                       `json:"target_units"`
	UnitID                  *int                        `json:"unit_id" form:"unit_id"`
	DeleteUnitID            *int                        `json:"delete_unit_id" form:"delete_unit_id"`
	Units                   *[]domain.UnitBatalSchedule `json:"units"`
	ActorSID                string                      `json:"actor_sid"`
	IsAdmin                 bool                        `json:"is_admin"`
}

type UpdateProkerBatalSettingResult struct {
	Message                 string                    `json:"message"`
	SettingKey              string                    `json:"setting_key"`
	ProkerBatalAktifTanggal string                    `json:"proker_batal_aktif_tanggal"`
	GlobalDate              string                    `json:"global_date"`
	HasDate                 bool                      `json:"has_date"`
	IsActive                bool                      `json:"is_active"`
	Status                  string                    `json:"status"`
	StatusLabel             string                    `json:"status_label"`
	UpdatedAt               string                    `json:"updated_at"`
	DetailMessage           string                    `json:"detail_message"`
	Units                   []domain.UnitBatalStatus  `json:"units"`
}
