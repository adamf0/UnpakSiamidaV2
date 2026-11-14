package GetProkerBatalSetting

import domain "UnpakSiamida/modules/prokerv3/domain"

type GetProkerBatalSettingQuery struct {
	TargetUnitID   *int
	TargetUnitUUID string
	UserUnit       string
	UserLevel      string
}

type GetProkerBatalSettingResult = domain.ProkerBatalSettingStatus
