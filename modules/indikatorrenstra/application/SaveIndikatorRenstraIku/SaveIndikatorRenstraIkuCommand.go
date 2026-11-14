package application

type SaveIndikatorRenstraIkuCommand struct {
	Identifier string   `json:"identifier"`
	IkuIds     []uint64 `json:"iku_ids"`
}
