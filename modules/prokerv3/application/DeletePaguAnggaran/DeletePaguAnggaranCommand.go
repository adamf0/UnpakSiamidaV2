package DeletePaguAnggaran

type DeletePaguAnggaranCommand struct {
	ID string
}

type DeletePaguAnggaranResult struct {
	Message string `json:"message"`
	ID      string `json:"id"`
}
