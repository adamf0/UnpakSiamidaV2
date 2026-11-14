package ReorderTimelineV3

import domain "UnpakSiamida/modules/aktivitasv2/domain"

type ReorderTimelineV3Command struct {
	Items []domain.ReorderTimelineItem `json:"items" form:"items"`
}
