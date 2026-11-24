package application

import (
	domainuser "UnpakSiamida/modules/user/domain"
)

type GetEmailHistoriesQuery struct {
	Search string
	Status string
	Page   int
	Limit  int
}

type GetEmailHistoriesResult struct {
	Data  []domainuser.HistoryEmail `json:"data"`
	Total int64                     `json:"total"`
	Page  int                       `json:"page"`
	Limit int                       `json:"limit"`
}
