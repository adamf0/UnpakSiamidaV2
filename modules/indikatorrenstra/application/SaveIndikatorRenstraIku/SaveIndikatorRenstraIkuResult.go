package application

type SaveIndikatorRenstraIkuResult struct {
	Message string   `json:"message"`
	Count   int      `json:"count"`
	IkuIds  []uint64 `json:"iku_ids"`
}
