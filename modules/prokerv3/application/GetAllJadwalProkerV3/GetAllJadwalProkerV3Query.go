package GetAllJadwalProkerV3

import (
	domain "UnpakSiamida/modules/prokerv3/domain"
)

type GetAllJadwalProkerV3Query struct {
	Filter domain.JadwalFilter
}

type GetAllJadwalProkerV3Result struct {
	Data  []domain.JadwalProkerV3DTO `json:"data"`
	Total int                        `json:"total"`
}
