package application

type CreateIndikatorRenstraCommand struct {
    StandarRenstra  string
	Indikator       string
	Parent          *string
	Tahun           string
	TipeTarget      string
	Operator        *string //opsional
}
