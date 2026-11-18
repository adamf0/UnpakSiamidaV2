package CreateMataProgramV2

type CreateMataProgramV2Command struct {
	MataProgram string
	Tahun       string
	Targets     []string
}

type CreateMataProgramV2Result struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	UUID    string `json:"uuid"`
	Count   int    `json:"count"`
}
