package ReorderTimelineV3

import domain "UnpakSiamida/modules/prokerv3/domain"

type ReorderTimelineV3Command struct {
	Items []domain.ReorderTimelineItem `json:"items" form:"items"`
}
