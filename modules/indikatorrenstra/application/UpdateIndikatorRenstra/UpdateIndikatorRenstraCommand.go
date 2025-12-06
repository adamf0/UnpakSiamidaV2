package application

type UpdateIndikatorRenstraCommand struct {
    Uuid     	 string
	StandarRenstra  string
	Indikator       string
	Parent          *string
	Tahun           string
	TipeTarget      string
	Operator        *string //opsional
}
