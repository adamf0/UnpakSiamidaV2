package CreateMataProgramV2

type CreateMataProgramV2Command struct {
	MataProgram        string
	Tahun              string
	TemplateRenstraIDs []uint
}

type CreateMataProgramV2Result struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	UUID    string `json:"uuid"`
	Count   int    `json:"count"`
}
